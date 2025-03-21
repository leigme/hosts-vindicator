/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	defer func() {
		if err := viper.WriteConfig(); err != nil {
			log.Println(err)
		}
	}()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	createConfigDir()
	initConfig()
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
	var homeDir string
	var err error
	if homeDir, err = os.UserHomeDir(); err != nil {
		homeDir = "."
	}
	configPath := filepath.Join(homeDir, ".config", "hosts-vindicator", "hv.yaml")
	return configPath
}

func initConfig() {
	viper.SetConfigName(filepath.Base(configPath()))
	viper.SetConfigFile(configPath())
	viper.SetConfigType("yaml")
}
