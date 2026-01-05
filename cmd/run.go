package cmd

import (
	"fmt"
	"os"

	"github.com/nodyhub/transformer"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <step file> (args...)",
	Short: "Execute the steps defined in the specified step file.",
	Long: `The run command allows you to execute a series of steps defined in a step file.
You can pass additional arguments after '--' which can be used within the steps.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stepFile := args[0]
		stepArgs := []string{}
		if len(args) > 1 {
			stepArgs = args[1:]
		}

		// read the step file
		data, err := os.ReadFile(stepFile)
		if err != nil {
			return fmt.Errorf("error reading step file: %w", err)
		}

		err = transformer.Run(cmd.Context(), string(data), stepArgs)
		if err != nil {
			return fmt.Errorf("error running steps: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
