package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of configuration",
	Long: LongUsage(
		`Show status of configuration

		  sbgraph status

		Show config file used and current settings.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		doStatus(cmd)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func doStatus(cmd *cobra.Command) {
	fmt.Printf("config file: %s\n", viper.ConfigFileUsed())
	fmt.Printf("config: %#v\n", config)
}
