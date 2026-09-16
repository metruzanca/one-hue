// Package gen turns built themes into editor-specific static files.
package gen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/metruzanca/one-hue-theme/internal/theme"
)

// File is one generated output file, with Path relative to the themes dir.
type File struct {
	Path    string
	Content []byte
}

// Target renders built themes for one editor. Adding a new editor means
// implementing this interface in this package and registering it in Targets.
//
// Render receives every requested theme at once so a target can emit a single
// family file (as Zed does) or one file per theme (as VSCode does).
type Target interface {
	Name() string
	Render(all []*theme.Built) ([]File, error)
}

// Targets lists every supported output format in the order used by the CLI.
var Targets = []Target{
	Vscode(),
	Zed(),
}

// RenderAll renders every theme for every target.
func RenderAll(all []*theme.Built) ([]File, error) {
	var files []File
	for _, t := range Targets {
		f, err := t.Render(all)
		if err != nil {
			return nil, fmt.Errorf("target %s: %w", t.Name(), err)
		}
		files = append(files, f...)
	}
	return files, nil
}

// WriteFiles writes generated files under the themes directory, creating
// subdirectories as needed.
func WriteFiles(themeDir string, files []File) error {
	for _, f := range files {
		p := filepath.Join(themeDir, f.Path)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, f.Content, 0o644); err != nil {
			return err
		}
	}
	return nil
}
