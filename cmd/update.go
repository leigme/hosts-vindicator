package cmd

/*
//Copyright © 2025 leig HERE <leigme@gmail.com>
*/
import (
	"bufio"
	_ "embed"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	//go:embed template/hosts.tpl
	hostsTpl  string
	stc, etc  string
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			data := make(map[string]any)
			loadConfig()
			url := viper.GetString(hostsUrl)
			stc = viper.GetString(startTag)
			if strings.EqualFold("", stc) {
				stc = defaultStartTag
			}
			data["StartTag"] = stc
			etc = viper.GetString(endTag)
			if strings.EqualFold("", etc) {
				etc = defaultEndTag
			}
			data["EndTag"] = etc
			if !strings.EqualFold("", url) && !conf.skip {
				downloadTmp(url, tmpPath())
			}
			readTmp(data)
			if !conf.replace {
				readHosts(data)
			}
			writeHosts(data)
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

func tmpPath() string {
	return filepath.Join(filepath.Dir(conf.file), hostsTmp)
}

func downloadTmp(url, fileName string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	bar := progressbar.DefaultBytes(
		-1, "downloading: ")
	out, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err = out.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

func readTmp(data map[string]any) {
	contents := make([]string, 0)
	tmp, err := os.Open(tmpPath())
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err = tmp.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	br := bufio.NewScanner(tmp)
	for br.Scan() {
		line := br.Text()
		if strings.EqualFold(line, stc) || strings.EqualFold(line, etc) {
			continue
		}
		contents = append(contents, line)
	}
	data["Contents"] = contents
}

func readHosts(data map[string]any) {
	hp := viper.GetString(hostsPath)
	hosts, err := os.Open(hp)
	defer func() {
		if err = hosts.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalln(err)
		}
		if err = os.MkdirAll(filepath.Dir(hp), 0666); err != nil {
			log.Fatalln(err)
		} else {
			if hosts, err = os.Create(viper.GetString(hostsPath)); err != nil {
				log.Fatalln(err)
			}
		}
		return
	}
	headers := make([]string, 0)
	footers := make([]string, 0)
	scanner := bufio.NewScanner(hosts)
	before := true
	after := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.EqualFold("", line) {
			continue
		}
		if strings.Contains(line, stc) {
			before = false
			continue
		}
		if before && !after {
			headers = append(headers, line)
			continue
		}
		if strings.Contains(line, etc) {
			after = true
			continue
		}
		if !before && after {
			footers = append(footers, line)
			continue
		}
	}
	data["Headers"] = headers
	data["Footers"] = footers
}

func writeHosts(data map[string]any) {
	writeFileByTemplate(viper.GetString(hostsPath), hostsTpl, data)
}

func writeFileByTemplate(filePath string, tpl string, data map[string]any) {
	t, err := template.New("").Parse(tpl)
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	bw := bufio.NewWriter(f)
	if err = t.Execute(bw, data); err != nil {
		log.Fatalln(err)
	}
	if _, err = bw.WriteString("\n"); err != nil {
		log.Fatalln(err)
	}
	if err = bw.Flush(); err != nil {
		log.Fatalln(err)
	}
}
