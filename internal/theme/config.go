package theme

import (
	"fmt"

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

// All lists every built-in theme variant.
func All() []Config {
	return []Config{MonochromePurple()}
}

// Find returns the built-in theme with the given slug.
func Find(slug string) (*Config, error) {
	for i := range All() {
		if All()[i].Slug == slug {
			c := All()[i]
			return &c, nil
		}
	}
	return nil, fmt.Errorf("unknown theme %q", slug)
}

// MonochromePurple is the purple monochromatic variant.
func MonochromePurple() Config {
	ring := color.StandardRing(32)
	return Config{
		Slug: "monochrome-purple",
		Name: "Monochrome-Purple",
		Props: color.PaletteProps{
			FG: color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Scalar(2),
				Hue:    color.Scalar(290),
			},
			BG: color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Scalar(2),
				Hue:    color.Scalar(290),
			},
			CoShades: &color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Scalar(16),
				Hue:    color.Scalar(290),
			},
			Accent: &color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Span(42, 36),
				Hue:    color.Scalar(277),
			},
			CoAccent: &color.Grade{
				Luma:   color.Span(6, 90),
				Chroma: color.Span(42, 36),
				Hue:    color.Scalar(245),
			},
			Ring: &color.Grade{
				Chroma: color.Span(32, 56),
			},
			Red:        &color.Grade{Hue: color.Scalar(ring[color.KeyRed])},
			Orange:     &color.Grade{Hue: color.Scalar(ring[color.KeyOrange])},
			Yellow:     &color.Grade{Hue: color.Scalar(ring[color.KeyYellow])},
			Chartreuse: &color.Grade{Hue: color.Scalar(ring[color.KeyChartreuse])},
			Green:      &color.Grade{Hue: color.Scalar(ring[color.KeyGreen])},
			Spring:     &color.Grade{Hue: color.Scalar(ring[color.KeySpring])},
			Cyan:       &color.Grade{Hue: color.Scalar(ring[color.KeyCyan])},
			Azure:      &color.Grade{Hue: color.Scalar(ring[color.KeyAzure])},
			Blue:       &color.Grade{Hue: color.Scalar(ring[color.KeyBlue])},
			Violet:     &color.Grade{Hue: color.Scalar(ring[color.KeyViolet])},
			Magenta:    &color.Grade{Hue: color.Scalar(ring[color.KeyMagenta])},
			Rose:       &color.Grade{Hue: color.Scalar(ring[color.KeyRose])},
		},
		BracketGrades: []color.PaletteKey{
			color.KeyAzure, color.KeySpring, color.KeyChartreuse,
			color.KeyOrange, color.KeyRose, color.KeyViolet,
		},
	}
}
