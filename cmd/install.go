package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/metruzanca/one-hue-theme/internal/gen"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

// VS Code extension identity used when assembling an extension on install.
// The <publisher>.<name> combination is what VS Code keys the extension by.
const (
	extPublisher = "metruzanca"
	extName      = "theme-monochrome"
	extVersion   = "0.1.0"
	extDisplay   = "Monochrome Theme"
	extLicense   = "MIT"
)

var (
	installEditors   string // --editors: comma list; empty prompts interactively
	installVscodeDir string
	installZedDir    string
)

// editor bundles one installable editor with its theme target.
type editor struct {
	name   string // flag/prompt value
	label  string // display name
	cli    string // binary that signals the editor is available
	target gen.Target
}

var editors = []editor{
	{"vscode", "VS Code", "code", gen.Vscode()},
	{"zed", "Zed", "zed", gen.Zed()},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Build and install theme files for one or more editors",
	Long: `Build the theme variants and copy them into each selected editor's
theme directory. The editor selection is presented as an interactive
multi-select; pass --editors to install non-interactively.

VS Code themes install as a local extension under ~/.vscode/extensions.
Zed themes install under ~/.config/zed/themes.`,
	Example: "  one-hue install\n" +
		"  one-hue install --editors vscode,zed\n" +
		"  one-hue install --editors zed",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstall()
	},
}

func init() {
	installCmd.Flags().StringVar(&installEditors, "editors", "",
		"comma-separated editors to install (vscode,zed); prompts when empty")
	installCmd.Flags().StringVar(&installVscodeDir, "vscode-dir", "",
		"override the VS Code extensions directory")
	installCmd.Flags().StringVar(&installZedDir, "zed-dir", "",
		"override the Zed themes directory")
	rootCmd.AddCommand(installCmd)
}

func runInstall() error {
	selected, err := resolveEditors()
	if err != nil {
		return err
	}
	if len(selected) == 0 {
		return fmt.Errorf("no editors selected")
	}

	configs, err := theme.Load(configFlag)
	if err != nil {
		return err
	}
	built := make([]*theme.Built, len(configs))
	for i := range configs {
		built[i] = theme.Build(configs[i])
	}

	for _, e := range selected {
		files, err := e.target.Render(built)
		if err != nil {
			return fmt.Errorf("render %s: %w", e.name, err)
		}
		if err := installFor(e, files, built); err != nil {
			return fmt.Errorf("install %s: %w", e.name, err)
		}
	}
	return nil
}

func resolveEditors() ([]editor, error) {
	if installEditors != "" {
		var out []editor
		for _, name := range strings.Split(installEditors, ",") {
			name = strings.TrimSpace(strings.ToLower(name))
			e, ok := findEditor(name)
			if !ok {
				return nil, fmt.Errorf("unknown editor %q (available: %s)", name, editorNames())
			}
			out = append(out, e)
		}
		return out, nil
	}

	var picked []string
	for _, e := range editors {
		if _, err := exec.LookPath(e.cli); err == nil {
			picked = append(picked, e.name)
		}
	}
	opts := make([]huh.Option[string], 0, len(editors))
	for _, e := range editors {
		opts = append(opts, huh.NewOption(e.label, e.name))
	}
	if err := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title("Install theme for").
			Options(opts...).
			Value(&picked),
	)).Run(); err != nil {
		return nil, fmt.Errorf("interactive selection failed (use --editors vscode,zed instead): %w", err)
	}

	var out []editor
	for _, name := range picked {
		if e, ok := findEditor(name); ok {
			out = append(out, e)
		}
	}
	return out, nil
}

func findEditor(name string) (editor, bool) {
	for _, e := range editors {
		if e.name == name {
			return e, true
		}
	}
	return editor{}, false
}

func editorNames() string {
	names := make([]string, len(editors))
	for i, e := range editors {
		names[i] = e.name
	}
	return strings.Join(names, ", ")
}

func installFor(e editor, files []gen.File, built []*theme.Built) error {
	switch e.name {
	case "vscode":
		return installVscode(files, built)
	case "zed":
		return installZed(files)
	}
	return fmt.Errorf("no installer for editor %q", e.name)
}

// installVscode lays the generated files out as a local extension: the
// manifest plus the themes/ tree under <extensions>/<publisher>.<name>.
// VS Code scans this directory on window reload.
func installVscode(files []gen.File, built []*theme.Built) error {
	dir := installVscodeDir
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dir = filepath.Join(home, ".vscode", "extensions")
	}
	extDir := filepath.Join(dir, extPublisher+"."+extName)

	if err := os.RemoveAll(filepath.Join(extDir, "themes")); err != nil {
		return err
	}
	for _, f := range files {
		p := filepath.Join(extDir, "themes", f.Path)
		if err := writeFile(p, f.Content); err != nil {
			return err
		}
		fmt.Printf("installed %s\n", p)
	}

	manifest, err := vscodeManifest(built)
	if err != nil {
		return err
	}
	p := filepath.Join(extDir, "package.json")
	if err := writeFile(p, manifest); err != nil {
		return err
	}
	fmt.Printf("installed %s\n", p)
	fmt.Println("Reload the VS Code window (Developer: Reload Window) to enable the theme.")
	return nil
}

// installZed copies each generated file into the Zed themes directory.
func installZed(files []gen.File) error {
	dir := installZedDir
	if dir == "" {
		var err error
		dir, err = zedThemesDir()
		if err != nil {
			return err
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, f := range files {
		p := filepath.Join(dir, filepath.Base(f.Path))
		if err := writeFile(p, f.Content); err != nil {
			return err
		}
		fmt.Printf("installed %s\n", p)
	}
	fmt.Println("The theme appears in the theme selector (cmd-k cmd-t) on next Zed start.")
	return nil
}

func writeFile(p string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}

// zedThemesDir returns the per-user Zed themes directory. Zed reads local
// themes from ~/.config/zed/themes on Linux and macOS, %APPDATA%\Zed\themes
// on Windows.
func zedThemesDir() (string, error) {
	if runtime.GOOS == "windows" {
		cfg, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(cfg, "Zed", "themes"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "zed", "themes"), nil
}

// vscodeManifest assembles the extension manifest contributing every built
// theme, mirroring the repo-root package.json used for vsce packaging.
func vscodeManifest(built []*theme.Built) ([]byte, error) {
	themes := make([]gen.Value, 0, len(built))
	for _, b := range built {
		uiTheme := "vs-dark"
		if !b.Palette.IsDark() {
			uiTheme = "vs-light"
		}
		themes = append(themes, []gen.Field{
			{Key: "label", Val: b.Config.Name},
			{Key: "uiTheme", Val: uiTheme},
			{Key: "path", Val: "./themes/vscode/" + b.Config.Slug + ".json"},
		})
	}
	return gen.Encode([]gen.Field{
		{Key: "name", Val: extName},
		{Key: "displayName", Val: extDisplay},
		{Key: "version", Val: extVersion},
		{Key: "description", Val: "One-hue themes for VS Code, generated by the one-hue CLI."},
		{Key: "categories", Val: []gen.Value{"Themes"}},
		{Key: "publisher", Val: extPublisher},
		{Key: "license", Val: extLicense},
		{Key: "contributes", Val: []gen.Field{
			{Key: "themes", Val: themes},
		}},
		{Key: "engines", Val: []gen.Field{
			{Key: "vscode", Val: "^1.33.0"},
		}},
	}), nil
}
