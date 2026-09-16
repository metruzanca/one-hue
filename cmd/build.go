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

	configs, err := theme.Load(configFlag)
	if err != nil {
		return err
	}
	if themeFlag != "" {
		match, err := theme.Find(themeFlag, configs)
		if err != nil {
			return err
		}
		configs = []theme.Config{*match}
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

	var files []gen.File
	for _, t := range targets {
		f, err := t.Render(built)
		if err != nil {
			return fmt.Errorf("%s: %w", t.Name(), err)
		}
		files = append(files, f...)
	}

	if err := gen.WriteFiles(outFlag, files); err != nil {
		return err
	}
	for _, f := range files {
		fmt.Printf("wrote %s/%s\n", outFlag, f.Path)
	}
	return nil
}
