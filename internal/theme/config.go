package theme

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/metruzanca/one-hue-theme/internal/color"
)

// Config describes one theme variant: the grade configuration plus how the
// bracket-highlight colors are picked from the ring.
type Config struct {
	// Slug is the output filename stem, e.g. "monochrome-purple".
	Slug string
	// Name is the display name in the editor, e.g. "Monochrome-Purple".
	Name string
	// Props feeds the palette builder.
	Props color.PaletteProps
	// BracketGrades picks the six bracket-highlight colors from the ring,
	// in order of nesting depth.
	BracketGrades []color.PaletteKey
}

// tomlFile is the shape of the themes configuration file.
type tomlFile struct {
	Themes []tomlTheme `toml:"theme"`
}

// tomlTheme is one theme variant: a name and an accent color. The slug is
// derived from the name and the rest of the color model is derived from the
// accent by fromAccent.
type tomlTheme struct {
	Name   string `toml:"name"`
	Accent string `toml:"accent"`
}

// Load reads a themes TOML file and turns each entry into a Config, deriving
// the whole monochromatic + accent model from the accent color.
func Load(path string) ([]Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read themes config: %w", err)
	}
	var f tomlFile
	if err := toml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse themes config %s: %w", path, err)
	}
	if len(f.Themes) == 0 {
		return nil, fmt.Errorf("themes config %s defines no themes", path)
	}

	seen := make(map[string]bool, len(f.Themes))
	configs := make([]Config, 0, len(f.Themes))
	for _, t := range f.Themes {
		if t.Name == "" {
			return nil, fmt.Errorf("themes config %s: theme is missing a name", path)
		}
		slug := Slugify(t.Name)
		if slug == "" {
			return nil, fmt.Errorf("themes config %s: theme name %q produces an empty slug", path, t.Name)
		}
		if seen[slug] {
			return nil, fmt.Errorf("themes config %s: duplicate theme slug %q", path, slug)
		}
		seen[slug] = true

		accent, err := color.Hex(t.Accent)
		if err != nil {
			return nil, fmt.Errorf("themes config %s: theme %q accent: %w", path, t.Name, err)
		}
		configs = append(configs, fromAccent(slug, t.Name, accent))
	}
	return configs, nil
}

// Slugify turns a theme name into its filename stem: lowercased, with runs of
// non-alphanumeric characters collapsed to single hyphens.
func Slugify(name string) string {
	var b strings.Builder
	prevHyphen := true
	for _, r := range strings.ToLower(name) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevHyphen = false
		case !prevHyphen:
			b.WriteByte('-')
			prevHyphen = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// Find returns the config with the given slug.
func Find(slug string, configs []Config) (*Config, error) {
	for i := range configs {
		if configs[i].Slug == slug {
			c := configs[i]
			return &c, nil
		}
	}
	return nil, fmt.Errorf("unknown theme %q", slug)
}

// fromAccent derives one theme's full grade set from a single accent color.
// The monochromatic model lives here: the accent hue drives the accent and
// co-accent grades and tints the neutral grays, while the semantic ring stays
// fixed so error/warning/git colors do not rotate with the accent.
func fromAccent(slug, name string, accent color.Color) Config {
	h := accent.Hue()
	return Config{
		Slug: slug,
		Name: name,
		Props: color.PaletteProps{
			FG: color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Scalar(2),
				Hue:    color.Scalar(h),
			},
			BG: color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Scalar(2),
				Hue:    color.Scalar(h),
			},
			CoShades: &color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Scalar(16),
				Hue:    color.Scalar(h),
			},
			Accent: &color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Span(42, 36),
				Hue:    color.Scalar(h),
			},
			CoAccent: &color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Span(42, 36),
				Hue:    color.Scalar(normalizeHue(h - 32)),
			},
			Ring: &color.Grade{
				Chroma: color.Span(32, 56),
			},
			Red:        &color.Grade{Hue: color.Scalar(ringHue(color.KeyRed))},
			Orange:     &color.Grade{Hue: color.Scalar(ringHue(color.KeyOrange))},
			Yellow:     &color.Grade{Hue: color.Scalar(ringHue(color.KeyYellow))},
			Chartreuse: &color.Grade{Hue: color.Scalar(ringHue(color.KeyChartreuse))},
			Green:      &color.Grade{Hue: color.Scalar(ringHue(color.KeyGreen))},
			Spring:     &color.Grade{Hue: color.Scalar(ringHue(color.KeySpring))},
			Cyan:       &color.Grade{Hue: color.Scalar(ringHue(color.KeyCyan))},
			Azure:      &color.Grade{Hue: color.Scalar(ringHue(color.KeyAzure))},
			Blue:       &color.Grade{Hue: color.Scalar(ringHue(color.KeyBlue))},
			Violet:     &color.Grade{Hue: color.Scalar(ringHue(color.KeyViolet))},
			Magenta:    &color.Grade{Hue: color.Scalar(ringHue(color.KeyMagenta))},
			Rose:       &color.Grade{Hue: color.Scalar(ringHue(color.KeyRose))},
		},
		BracketGrades: []color.PaletteKey{
			color.KeyAzure, color.KeySpring, color.KeyChartreuse,
			color.KeyOrange, color.KeyRose, color.KeyViolet,
		},
	}
}

// ringBase anchors the semantic ring colors independent of the accent, so the
// 12 ring colors stay in their canonical positions across every theme.
const ringBase = 32.0

// ringHue returns the fixed LCH hue of one semantic ring color.
func ringHue(k color.PaletteKey) float64 {
	off := color.StandardRing(ringBase)
	return off[k]
}

// normalizeHue wraps an angle into 0..360.
func normalizeHue(h float64) float64 {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	return h
}