// Command one-hue generates static theme files for code editors from a
// shared color model.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/metruzanca/one-hue-theme/internal/gen"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

func usage() {
	fmt.Fprintf(os.Stderr, `one-hue — generate one-hue editor themes

Usage:
  one-hue build [flags]   build theme files (default command)
  one-hue list            list available themes and output targets

Build flags:
  --target <name|all>   editor to render (default: all)
  --theme  <slug>       theme variant to build (default: all)
  --out    <dir>        output directory, files go under <dir>/<target>/ (default: themes)
  --print               print the palette grades as terminal color blocks
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		fs := flag.NewFlagSet("build", flag.ExitOnError)
		target := fs.String("target", "all", "editor target: vscode, zed, or all")
		slug := fs.String("theme", "", "theme variant slug; empty builds all")
		out := fs.String("out", "themes", "output directory")
		printGrades := fs.Bool("print", false, "print palette grades as ANSI blocks")
		fs.Usage = func() {
			usage()
			fs.PrintDefaults()
		}
		if err := fs.Parse(os.Args[2:]); err != nil {
			os.Exit(2)
		}
		os.Exit(runBuild(*target, *slug, *out, *printGrades))
	case "list":
		os.Exit(runList())
	case "--help", "-h", "help":
		usage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "one-hue: unknown command %q\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func runBuild(targetName, slug, outDir string, printGrades bool) int {
	// Resolve requested targets.
	targets := gen.Targets
	if !strings.EqualFold(targetName, "all") {
		t, err := gen.FindTarget(targetName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		targets = []gen.Target{t}
	}

	// Resolve requested themes.
	configs := theme.All()
	if slug != "" {
		var match *theme.Config
		for i := range configs {
			if configs[i].Slug == slug {
				match = &configs[i]
				break
			}
		}
		if match == nil {
			fmt.Fprintf(os.Stderr, "one-hue: unknown theme %q\n", slug)
			return 1
		}
		configs = []theme.Config{*match}
	}

	var files []gen.File
	for _, cfg := range configs {
		built := theme.Build(cfg)
		if printGrades {
			fmt.Printf("%s grades:\n", built.Config.Name)
			built.Palette.Print()
			fmt.Println()
		}
		for _, t := range targets {
			f, err := t.Render([]*theme.Built{built})
			if err != nil {
				fmt.Fprintf(os.Stderr, "one-hue: %s: %v\n", t.Name(), err)
				return 1
			}
			files = append(files, f...)
		}
	}

	if err := gen.WriteFiles(outDir, files); err != nil {
		fmt.Fprintf(os.Stderr, "one-hue: %v\n", err)
		return 1
	}
	for _, f := range files {
		fmt.Printf("wrote %s/%s\n", outDir, f.Path)
	}
	return 0
}

func runList() int {
	fmt.Println("Themes:")
	for _, cfg := range theme.All() {
		fmt.Printf("  %-20s %s\n", cfg.Slug, cfg.Name)
	}
	fmt.Println("Targets:")
	for _, t := range gen.Targets {
		fmt.Printf("  %-8s %s\n", t.Name(), "editor theme renderer")
	}
	return 0
}
