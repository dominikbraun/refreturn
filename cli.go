package main

import (
	"github.com/spf13/cobra"
)

func rootCommand(version string) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:     "refreturn <PATH>",
		Version: version,
		Short:   `Find functions that return a reference and cause allocations.`,
		Long: `refreturn finds all Go functions in a directory tree that return a reference
and cause a potential unnecessary heap allocation.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return Run(path)
		},
	}

	return rootCommand
}
