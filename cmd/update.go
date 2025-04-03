/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"github.com/schollz/progressbar/v3"
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
	hostsTmp        = "hosts.tmp"
	defaultStartTag = "# GitHub IP hosts Start"
	defaultEndTag   = "# GitHub IP hosts End"
)

var (
	stc, etc         string
	skip, replace    bool
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
			stc = viper.GetString(startTag)
			if strings.EqualFold("", stc) {
				stc = defaultStartTag
			}
			etc = viper.GetString(endTag)
			if strings.EqualFold("", etc) {
				etc = defaultEndTag
			}
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
	bar := progressbar.DefaultBytes(
		resp.ContentLength, "downloading: ")
	out, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer out.Close()

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
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
		if err := os.MkdirAll(filepath.Dir(hp), 0666); err != nil {
			log.Fatalln(err)
		} else {
			if hosts, err = os.Create(viper.GetString(hostsPath)); err != nil {
				log.Fatalln(err)
			}
		}
		return
	}
	scanner := bufio.NewScanner(hosts)
	before := true
	after := false
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, startTag) {
			before = false
			continue
		}
		if before && !after {
			headers = append(headers, line)
		}
		if !before && after {
			footers = append(footers, line)
		}
		if strings.Contains(line, endTag) {
			after = true
			continue
		}
	}
}

func writeHosts() {
	tmpStat, err := os.Stat(tmpPath())
	if err != nil {
		log.Fatalln(err)
	}
	hosts, err := os.OpenFile(viper.GetString(hostsPath), os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer hosts.Close()

	bar := progressbar.DefaultBytes(tmpStat.Size(), "writing: ")
	hb := io.MultiWriter(hosts, bar)

	bw := bufio.NewWriter(hb)
	io.MultiWriter(bw, bar)
	for _, header := range headers {
		bw.WriteString(header + "\n")
	}
	bw.WriteString(stc + "\n")

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
	bw.WriteString(etc + "\n")
	for _, footer := range footers {
		if _, err = bw.WriteString(footer + "\n"); err != nil {
			log.Fatalln(err)
		}
	}
	if err = bw.Flush(); err != nil {
		log.Fatalln(err)
	}
}
