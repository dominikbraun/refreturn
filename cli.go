package main

import (
	"github.com/spf13/cobra"
)

const (
	shortHelp = `Find functions that return a reference and cause allocations.`
	longHelp  = `refreturn finds all Go functions in a directory tree that return a
reference and cause a potential unnecessary heap allocation.`
)

// rootCommand generates the top-level `refreturn` command.
func rootCommand(version string) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:     "refreturn <PATH>",
		Version: version,
		Short:   shortHelp,
		Long:    longHelp,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			return Run(path)
		},
	}

	return rootCommand
}
