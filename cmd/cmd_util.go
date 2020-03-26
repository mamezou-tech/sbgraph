package cmd

import (
	"fmt"
	"os"
	"github.com/MakeNowJust/heredoc/v2"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func LongUsage(s string) string {
	if len(s) == 0 {
		return s
	}
	return heredoc.Doc(s)
}
