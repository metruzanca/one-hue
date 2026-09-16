package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/metruzanca/one-hue-theme/internal/gen"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available theme variants and editor targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		configs, err := theme.Load(configFlag)
		if err != nil {
			return err
		}
		fmt.Println("Themes:")
		for _, cfg := range configs {
			fmt.Printf("  %-20s %s\n", cfg.Slug, cfg.Name)
		}
		fmt.Println("Targets:")
		for _, t := range gen.Targets {
			fmt.Printf("  %-8s editor theme renderer\n", t.Name())
		}
		return nil
	},
}
