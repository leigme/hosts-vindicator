package cmd

/*
//Copyright © 2025 leig HERE <leigme@gmail.com>
*/
import (
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

type config struct {
	skip      bool   `yaml:"skip"`
	replace   bool   `yaml:"replace"`
	hostsUrl  string `yaml:"hosts_url"`
	hostsPath string `yaml:"hosts_path"`
	startTag  string `yaml:"start_tag"`
	endTag    string `yaml:"end_tag"`
	file      string `yaml:"file"`
}

type hostsTemplate struct {
	header string
	start  string
	lines  []string
	footer string
	end    string
}

const (
	hostsUrl        = "hosts_url"
	hostsPath       = "hosts_path"
	startTag        = "start_tag"
	endTag          = "end_tag"
	hostsTmp        = "hosts.tmp"
	defaultStartTag = "# GitHub IP hosts Start"
	defaultEndTag   = "# GitHub IP hosts End"
)

// configCmd represents the config command
var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("args len: %d\n", len(args))
			defer func() {
				if err := viper.WriteConfig(); err != nil {
					log.Fatalln(err)
				}
			}()
			viper.Set(hostsUrl, conf.hostsUrl)
			viper.Set(hostsPath, conf.hostsPath)
		},
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&conf.hostsUrl, "url", "u", "https://raw.githubusercontent.com/ittuann/GitHub-IP-hosts/refs/heads/main/hosts", "")
	configCmd.Flags().StringVarP(&conf.hostsPath, "path", "p", "/etc/hosts", "")
}
