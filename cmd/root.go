// Package cmd implements the one-hue command line interface.
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configFlag = "themes.toml"
	outFlag    = "themes"
	printFlag  = false
)

// rootCmd lists the subcommands; running `one-hue` with no arguments prints
// its help.
var rootCmd = &cobra.Command{
	Use:   "one-hue",
	Short: "Generate static theme files for code editors",
	Long: `One-hue generates static theme files for code editors from a single,
shared color model. It builds the configured theme variants and writes the
results under the output directory, one subdirectory per editor.

Supported editors (targets): vscode, zed, opencode, herdr.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "themes.toml",
		"TOML file defining theme variants (slug, name, accent)")
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
