package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
)

// CheckErr will print error message to stderr and exit with code 1
func CheckErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

// LongUsage is helper for print long usage
func LongUsage(s string) string {
	if len(s) == 0 {
		return s
	}
	return heredoc.Doc(s)
}

// CheckProject will check project name is not empty
func CheckProject(projectName string) {
	if projectName == "" {
		fmt.Fprintln(os.Stderr, "Project name not set. execute `sbgraph project -p <project name>`")
		os.Exit(1)
	}
}

// Contains will check if a string is present in a slice ignoring case
func Contains(s []string, str string) bool {
	for _, v := range s {
		if strings.ToLower(v) == strings.ToLower(str) {
			return true
		}
	}
	return false
}
