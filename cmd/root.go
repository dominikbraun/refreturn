package cmd

import (
	"fmt"
	"os"

	"github.com/dominikbraun/refreturn/core"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "refreturn <directory>",
	Short: `Find functions that return a reference and cause allocations.`,
	Long: `refreturn finds all Go functions in a directory tree that return a reference
and cause a potential unnecessary heap allocation.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		core.Run(args[0])
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
