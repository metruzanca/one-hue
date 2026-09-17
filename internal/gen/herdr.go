package gen

import (
	"fmt"
	"strings"

	"github.com/metruzanca/one-hue-theme/internal/color"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

// herdrTarget renders one TOML snippet per built theme. Herdr keeps a single
// active theme in its config.toml, so each file is a self-contained
// [theme.custom] override block the user can paste or the install command
// merges into ~/.config/herdr/config.toml.
type herdrTarget struct{}

// Herdr returns the Herdr output target.
func Herdr() Target { return herdrTarget{} }

func (herdrTarget) Name() string { return "herdr" }

func (h herdrTarget) Render(all []*theme.Built) ([]File, error) {
	var files []File
	for _, b := range all {
		files = append(files, File{
			Path:    "herdr/" + b.Config.Slug + ".toml",
			Content: h.buildSnippet(b),
		})
	}
	return files, nil
}

// buildSnippet renders the [theme.custom] override block for one theme.
func (h herdrTarget) buildSnippet(b *theme.Built) []byte {
	pl := b.Palette
	fg := pl.Get(color.KeyFG)
	bg := pl.Get(color.KeyBG)
	accent := pl.Get(color.KeyAccent)
	red := pl.Get(color.KeyRed)
	orange := pl.Get(color.KeyOrange)
	yellow := pl.Get(color.KeyYellow)
	green := pl.Get(color.KeyGreen)
	cyan := pl.Get(color.KeyCyan)
	blue := pl.Get(color.KeyBlue)

	var s strings.Builder
	s.WriteString("[theme.custom]\n")
	write := func(key, hex string) {
		fmt.Fprintf(&s, "%s = %q\n", key, hex)
	}
	write("accent", accent[7].Hex())
	write("panel_bg", bg[1].Hex())
	write("sidebar_bg", bg[0].Hex())
	write("active_row_bg", bg[2].Hex())
	write("selection_bg", bg[3].Hex())
	write("surface0", bg[1].Hex())
	write("surface1", bg[2].Hex())
	write("surface_dim", bg[0].Hex())
	write("overlay0", bg[3].Hex())
	write("overlay1", bg[4].Hex())
	write("text", fg[8].Hex())
	write("subtext0", fg[6].Hex())
	write("mauve", accent[7].Hex())
	write("green", green[7].Hex())
	write("yellow", yellow[7].Hex())
	write("red", red[7].Hex())
	write("blue", blue[7].Hex())
	write("teal", cyan[7].Hex())
	write("peach", orange[7].Hex())
	return []byte(s.String())
}

var _ Target = herdrTarget{}
