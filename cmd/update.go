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
	download         bool
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
				if download {
					downloadTmp(url, tmpPath())
				}
				readHosts()
				writeHosts()
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&download, "download", "d", true, "")
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
	oldHosts, err := os.Open(viper.GetString(hostsPath))
	if err != nil {
		log.Fatalln(err)
	}
	defer oldHosts.Close()
	scanner := bufio.NewScanner(oldHosts)
	skip := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.EqualFold(line, startTag) {
			skip = true
			continue
		}
		if !skip {
			headers = append(headers, line)
		}
		if strings.EqualFold(line, endTag) {
			skip = false
			continue
		}
		if !skip {
			footers = append(footers, line)
		}
	}
}

func writeHosts() {
	newHosts, err := os.Open(viper.GetString(hostsPath))
	if err != nil {
		log.Fatalln(err)
	}
	defer newHosts.Close()
	bw := bufio.NewWriter(newHosts)
	for _, header := range headers {
		bw.WriteString(header)
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
		if _, err = bw.WriteString(line); err != nil {
			log.Fatalln(err)
		}
	}
	for _, footer := range footers {
		if _, err = bw.WriteString(footer); err != nil {
			log.Fatalln(err)
		}
	}
	if err = bw.Flush(); err != nil {
		log.Fatalln(err)
	}
}
