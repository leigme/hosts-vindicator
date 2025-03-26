/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	hostsTmp = "hosts.tmp"
	startTag = "# GitHub IP hosts Start"
	endTag   = "# GitHub IP hosts End"
)

var (
	skip             bool
	replace          bool
	headers, footers []string
	updateCmd        = &cobra.Command{
		Use:   "update",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			loadConfig()
			url := viper.GetString(hostsUrl)
			if !strings.EqualFold("", url) {
				if !skip {
					downloadTmp(url, tmpPath())
				}
				if !replace {
					readHosts()
				}
				writeHosts()
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&skip, "skip", "s", false, "")
	updateCmd.Flags().BoolVarP(&replace, "replace", "r", false, "")
}

func tmpPath() string {
	return filepath.Join(filepath.Dir(file), hostsTmp)
}

func downloadTmp(url, fileName string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

func readHosts() {
	hp := viper.GetString(hostsPath)
	hosts, err := os.Open(hp)
	defer hosts.Close()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalln(err)
		}
		if err := os.MkdirAll(filepath.Dir(hp), os.ModePerm); err != nil {
			log.Fatalln(err)
		} else {
			if hosts, err = os.Create(viper.GetString(hostsPath)); err != nil {
				log.Fatalln(err)
			}
		}
		return
	}
	scanner := bufio.NewScanner(hosts)
	ignore := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, startTag) {
			ignore = true
			continue
		}
		if !ignore {
			headers = append(headers, line)
			continue
		}
		if strings.Contains(line, endTag) {
			ignore = false
			continue
		}
		if !ignore {
			footers = append(footers, line)
		}
	}
}

func writeHosts() {
	hosts, err := os.OpenFile(viper.GetString(hostsPath), os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer hosts.Close()
	bw := bufio.NewWriter(hosts)
	if len(headers) == 0 {
		headers = append(headers, startTag)
	}
	for _, header := range headers {
		bw.WriteString(header + "\n")
	}
	tmp, err := os.Open(tmpPath())
	if err != nil {
		log.Fatalln(err)
	}
	defer tmp.Close()
	br := bufio.NewScanner(tmp)
	for br.Scan() {
		line := br.Text()
		if strings.EqualFold(line, startTag) || strings.EqualFold(line, endTag) {
			continue
		}
		if _, err = bw.WriteString(line + "\n"); err != nil {
			log.Fatalln(err)
		}
	}
	if len(footers) == 0 {
		footers = append(footers, endTag)
	}
	for _, footer := range footers {
		if _, err = bw.WriteString(footer + "\n"); err != nil {
			log.Fatalln(err)
		}
	}
	if err = bw.Flush(); err != nil {
		log.Fatalln(err)
	}
}
