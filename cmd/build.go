package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/metruzanca/one-hue-theme/internal/gen"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static theme files for all editors",
	Example: "  one-hue build\n" +
		"  one-hue build --print",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBuild()
	},
}

func runBuild() error {
	configs, err := theme.Load(configFlag)
	if err != nil {
		return err
	}

	built := make([]*theme.Built, len(configs))
	for i, cfg := range configs {
		built[i] = theme.Build(cfg)
		if printFlag {
			fmt.Printf("%s grades:\n", built[i].Config.Name)
			built[i].Palette.Print()
			fmt.Println()
		}
	}

	files, err := gen.RenderAll(built)
	if err != nil {
		return err
	}

	if err := gen.WriteFiles(outFlag, files); err != nil {
		return err
	}
	for _, f := range files {
		fmt.Printf("wrote %s/%s\n", outFlag, f.Path)
	}
	return nil
}
