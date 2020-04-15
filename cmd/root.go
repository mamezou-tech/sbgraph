package cmd

import (
	"os"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Config store settings
type Config struct {
	WorkDir        string `yaml:"workdir"`
	CurrentProject string `yaml:"currentproject"`
}

var cfgFile string

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sbgraph",
	Short: "A CLI to analize Scrapbox project.",
	Long: LongUsage(`
		sbgraph is a CLI to fetch data and visualize Scrapbox projects.
		  Fetch page data (JSON format)
		  Aggregate user activities (pages created, views of created page, etc.)
		  Generate graph data (as Graphviz dot file)
	`),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	CheckErr(err)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sbgraph.yaml)")
	wd, err := os.Getwd()
	CheckErr(err)
	wkdir := filepath.FromSlash(wd + "/_work")
	rootCmd.PersistentFlags().StringP("workdir", "d", wkdir, "working directory")
	rootCmd.PersistentFlags().StringP("currentproject", "c", "", "current project")
	viper.BindPFlag("workdir", rootCmd.PersistentFlags().Lookup("workdir"))
	viper.BindPFlag("currentproject", rootCmd.PersistentFlags().Lookup("currentproject"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		CheckErr(err)
		confPath := filepath.FromSlash(home + "/.sbgraph.yaml")
		viper.SetConfigFile(confPath)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Read config file in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
		err := viper.Unmarshal(&config)
		CheckErr(err)
	}
	SaveConfigAs()
}

func SaveConfigAs() {
	home, err := homedir.Dir()
	CheckErr(err)
	confPath := filepath.FromSlash(home + "/.sbgraph.yaml")
	if err := viper.SafeWriteConfigAs(confPath); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfigAs(confPath)
			CheckErr(err)
		}
	}
}

func SaveConfig() {
	home, err := homedir.Dir()
	CheckErr(err)
	confPath := filepath.FromSlash(home + "/.sbgraph.yaml")
	viper.SetConfigFile(confPath)
	err = viper.WriteConfig()
	CheckErr(err)
	if err = viper.ReadInConfig(); err == nil {
		fmt.Println("reload config file:", viper.ConfigFileUsed())
		err := viper.Unmarshal(&config)
		CheckErr(err)
	}
	fmt.Printf("config update: %#v\n", config)
}
