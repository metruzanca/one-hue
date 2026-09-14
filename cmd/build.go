package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/metruzanca/one-hue-theme/internal/gen"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build static theme files for one or all editors",
	Example: "  one-hue build\n" +
		"  one-hue build --target zed\n" +
		"  one-hue build --theme monochrome-purple --target vscode\n" +
		"  one-hue build --print",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBuild()
	},
}

func runBuild() error {
	targets := gen.Targets
	if !strings.EqualFold(targetFlag, "all") {
		t, err := gen.FindTarget(targetFlag)
		if err != nil {
			return err
		}
		targets = []gen.Target{t}
	}

	configs := theme.All()
	if themeFlag != "" {
		match, err := theme.Find(themeFlag)
		if err != nil {
			return err
		}
		configs = []theme.Config{*match}
	}

	var files []gen.File
	for _, cfg := range configs {
		built := theme.Build(cfg)
		if printFlag {
			fmt.Printf("%s grades:\n", built.Config.Name)
			built.Palette.Print()
			fmt.Println()
		}
		for _, t := range targets {
			f, err := t.Render([]*theme.Built{built})
			if err != nil {
				return fmt.Errorf("%s: %w", t.Name(), err)
			}
			files = append(files, f...)
		}
	}

	if err := gen.WriteFiles(outFlag, files); err != nil {
		return err
	}
	for _, f := range files {
		fmt.Printf("wrote %s/%s\n", outFlag, f.Path)
	}
	return nil
}