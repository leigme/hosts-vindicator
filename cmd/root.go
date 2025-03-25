/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	conf    string
	rootCmd = &cobra.Command{
		Use:   "hosts-vindicator",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.Println("execute: " + conf)
	createConfigDir()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&conf, "conf", "", "")
}

func createConfigDir() {
	configDir := filepath.Dir(configPath())
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			log.Fatalln(err)
		}
	}
}

func configPath() string {
	if !strings.EqualFold(conf, "") {
		return conf
	}
	var homeDir string
	var err error
	if homeDir, err = os.UserHomeDir(); err != nil {
		homeDir = "."
	}
	configPath := filepath.Join(homeDir, ".config", "hosts-vindicator", "hv.yaml")
	return configPath
}

func initConfig() {
	log.Println(configPath())
	viper.SetConfigName(filepath.Base(configPath()))
	viper.SetConfigFile(configPath())
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
}
