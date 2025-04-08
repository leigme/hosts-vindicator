package cmd

/*
//Copyright © 2025 leig HERE <leigme@gmail.com>
*/
import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var (
	conf    = config{}
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
	rootCmd.PersistentFlags().StringVarP(&conf.file, "file", "f", defaultPath(), "")
	rootCmd.PersistentFlags().BoolVarP(&conf.skip, "skip", "s", false, "")
	rootCmd.PersistentFlags().BoolVarP(&conf.replace, "replace", "r", false, "")
}

func configDir() {
	if err := os.MkdirAll(filepath.Dir(conf.file), os.ModePerm); err != nil {
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
	viper.SetConfigName(filepath.Base(conf.file))
	viper.SetConfigFile(conf.file)
	viper.SetConfigType("yaml")
}

func loadConfig() {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalln(err)
	}
}
