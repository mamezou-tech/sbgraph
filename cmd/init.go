package cmd

import (
	"fmt"

	"github.com/kondoumh/scrapbox-viz/pkg/file"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize working directory",
	Long: LongUsage(`
		Initalize working directory for Scrapbox project.

		  sbv init
		
		if 'workdir' exists in $HOME/.sbv.yaml or set by global -d(--workdir) flag, it will be created.
		`),
	Run: func(cmd *cobra.Command, args []string) {
		doInit(cmd)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func doInit(cmd *cobra.Command) {
	fmt.Printf("Check and create workdir : %s", config.WorkDir)
	err := file.CreateDir(config.WorkDir)
	CheckErr(err)
}
