/*
Copyright © 2026 Jan Harrie <jan@nody.cc>
*/
package cmd

import (
	"fmt"

	"github.com/nodyhub/transformer/registry"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List modules that are available",
	Run: func(cmd *cobra.Command, args []string) {
		for _, dep := range registry.ListModules() {
			fmt.Printf("%s\n", dep)
		}
	},
}

func init() {
	dependenciesCmd.AddCommand(listCmd)
}
