/*
Copyright © 2026 Jan Harrie <jan@nody.cc>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// dependenciesCmd represents the dependencies command
var dependenciesCmd = &cobra.Command{
	Use:   "modules",
	Short: "Manage modules for Transformer",
}

func init() {
	rootCmd.AddCommand(dependenciesCmd)
}
