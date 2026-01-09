/*
Copyright © 2026 Jan Harrie <jan@nody.cc>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/nodyhub/transformer/registry"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <module> [<module>...]",
	Short: "install a module that is available and not yet installed",
	Run: func(cmd *cobra.Command, args []string) {
		var hooks []string
		if len(args) > 0 {
			switch {
			case args[0] == "all":
				slog.Debug("installing all available modules")
				args = registry.ListModules()
			case strings.Contains(args[0], ","):
				args = strings.Split(args[0], ",")
			case len(args[0]) == 0:
				slog.Debug("no modules specified for installation")
				fmt.Println("please specify at least one module to install, or 'all' to install all available modules.")
				return
			}
			hooks = append(hooks, args...)

		} else {
			fmt.Println("please specify at least one module to install, or 'all' to install all available modules.")
			return
		}

		ctx := cmd.Context()
		for _, mod := range hooks {
			slog.Debug("installing module", "module", mod)
			err := registry.RunModuleHooks(ctx, mod)
			if err != nil {
				slog.Error("failed to install module", "module", mod, "error", err)
				os.Exit(1)
			} else {
				fmt.Printf("successfully installed module '%s'\n", mod)
			}
		}
	},
}

func init() {
	dependenciesCmd.AddCommand(installCmd)
}
