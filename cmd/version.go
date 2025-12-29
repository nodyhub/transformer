/*
Copyright © 2025 Jan Harrie <jan@nody.cc>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	Commit    = ""
	BuildTime = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommit: %s\nBuildTime: %s\n", Version, Commit, BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
