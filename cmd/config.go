/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	hostsUrl  = "hosts_url"
	hostsPath = "hosts_path"
	startTag  = "start_tag"
	endTag    = "end_tag"
)

var file, url, path string

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := viper.WriteConfig(); err != nil {
				log.Fatalln(err)
			}
		}()
		viper.Set(hostsUrl, url)
		viper.Set(hostsPath, path)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&url, "url", "u", "https://raw.githubusercontent.com/ittuann/GitHub-IP-hosts/refs/heads/main/hosts", "")
	configCmd.Flags().StringVarP(&path, "path", "p", "/etc/hosts", "")
}
