package cmd

import (
	"github.com/spf13/cobra"

	"github.com/metruzanca/one-hue-theme/internal/serve"
)

var serveAddr string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the theme preview web app",
	Long: `Serve a local web app for previewing, creating and editing theme
variants. The app shows a live CodeMirror preview of each theme and lets you
change a theme's name and accent color; on save it rewrites themes.toml and
regenerates every editor target.`,
	Example: "  one-hue serve\n" +
		"  one-hue serve --addr :9000\n" +
		"  one-hue serve --config themes.toml --out themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return serve.Serve(serve.Options{
			Addr:       serveAddr,
			ConfigPath: configFlag,
			OutDir:     outFlag,
		})
	},
}

func init() {
	serveCmd.Flags().StringVar(&serveAddr, "addr", ":8080",
		"listen address, e.g. :8080 or 127.0.0.1:9000")
	rootCmd.AddCommand(serveCmd)
}
