/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const hostsTmp = "hosts.tmp"

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if url := viper.GetString(hostsUrl); !strings.EqualFold("", url) {
			downloadTmp(url, filepath.Join(filepath.Dir(configPath()), hostsTmp))
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func downloadTmp(url, fileName string) {
	// https://github.com/ittuann/GitHub-IP-hosts/blob/main/hosts
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
