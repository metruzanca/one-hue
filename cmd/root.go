// Package cmd implements the one-hue command line interface.
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	targetFlag = "all"
	themeFlag  = ""
	outFlag    = "themes"
	printFlag  = false
)

// rootCmd is the default build command: running `one-hue` with no subcommand
// builds every theme variant for every editor.
var rootCmd = &cobra.Command{
	Use:   "one-hue",
	Short: "Generate static theme files for code editors",
	Long: `One-hue generates static theme files for code editors from a single,
shared color model. It builds the configured theme variants and writes the
results under the output directory, one subdirectory per editor.

Supported editors (targets): vscode, zed.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBuild()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&targetFlag, "target", "all",
		"editor target to render: vscode, zed, or all")
	rootCmd.PersistentFlags().StringVar(&themeFlag, "theme", "",
		"theme variant slug; empty builds all variants")
	rootCmd.PersistentFlags().StringVar(&outFlag, "out", "themes",
		"output directory, files go under <out>/<target>/")
	rootCmd.PersistentFlags().BoolVar(&printFlag, "print", false,
		"print the palette grades as terminal color blocks")

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(listCmd)
}

// Execute runs the CLI and returns the first error encountered, if any.
func Execute() error {
	return rootCmd.Execute()
}
