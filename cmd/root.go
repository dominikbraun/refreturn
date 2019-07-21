package cmd

import (
	"fmt"
	"os"

	"github.com/dominikbraun/refreturn/core"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "refreturn",
	Short: `Find functions that return a reference and cause allocations.`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Run(".")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
