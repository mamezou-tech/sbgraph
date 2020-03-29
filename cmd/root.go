package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Config store settings
type Config struct {
	WorkDir string `yaml:"workdir"`
}

var cfgFile string

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sbv",
	Short: "A CLI to analize Scrapbox project.",
	Long:  LongUsage(`
		scrapbox-vix (sbv) is a CLI to fetch data and visualize Scrapbox projects.
		  fetch page data (JSON format)
		  aggregate user activities (pages created, views of created page, etc.)
		  generate graph data (GraphViz dot file)
	`),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	CheckErr(err)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sbv.yaml)")
	rootCmd.PersistentFlags().StringP("workdir", "d", "_work", "working directory")
	viper.BindPFlag("workdir", rootCmd.PersistentFlags().Lookup("workdir"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".sbv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		if err := viper.Unmarshal(&config); err == nil {
			fmt.Printf("config: %#v\n", config)
		}
	}
}
