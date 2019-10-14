package cli

import (
	"fmt"
	"os"

	"github.com/dominikbraun/refreturn/core"

	"github.com/spf13/cobra"
)

// rootCmd depicts the main command for refreturn. At the
// moment, there are no subcommands added to it.
var rootCmd = &cobra.Command{
	Use:     "refreturn <directory>",
	Version: "0.1.0",
	Short:   `Find functions that return a reference and cause allocations.`,
	Long: `refreturn finds all Go functions in a directory tree that return a reference
and cause a potential unnecessary heap allocation.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing all functions that return a reference:")
		core.Run(args[0])
		fmt.Println("Finished.")
	},
}

// Execute runs the root command instance which will then
// trigger the execution of the actual logic.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
