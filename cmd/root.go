/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "hosts-vindicator",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", defaultPath(), "")
}

func configDir() {
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		if !os.IsExist(err) {
			log.Fatalln(err)
		}
	}
}

func defaultPath() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".config", execName(), "hv.yaml")
	}
	return ""
}

func execName() string {
	if execPath, err := os.Executable(); err == nil {
		return filepath.Base(execPath)
	}
	return "hosts-vindicator"
}

func initConfig() {
	configDir()
	viper.SetConfigName(filepath.Base(file))
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")
}

func loadConfig() {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
}
