package theme

import (
	"github.com/metruzanca/one-hue-theme/internal/color"
)

// Alpha fractions shared by the renderers.
const (
	SelectionAlpha       = 5.0 / 16.0
	HighlightAlpha       = 4.0 / 16.0
	HighlightBorderAlpha = 12.0 / 16.0
)

// FontStyle is how a token may be emphasized.
type FontStyle string

// Font styles.
const (
	FontNormal FontStyle = ""
	FontItalic FontStyle = "italic"
	FontBold   FontStyle = "bold"
)

// Role is one derived token color with optional emphasis.
type Role struct {
	Color color.Color
	Font  FontStyle
}

// Roles holds every syntax role the renderers map into their own scope sets.
// The values mirror the token settings of the original generator.
type Roles struct {
	Identifier    Role
	Parameter     Role
	Property      Role
	Comment       Role
	Operator      Role
	WeakOperator  Role
	Literal       Role
	Quote         Role
	Declare       Role
	Method        Role
	TypeName      Role
	Namespace     Role
	BuiltInType   Role
	TypeParameter Role
	Punctuation   Role
	WeakKeyword   Role
	Keyword       Role
	Storage       Role
	Library       Role
	Invalid       Role
}

// Built is a fully built theme: its palette plus the derived syntax roles.
type Built struct {
	Config  Config
	Palette *color.Palette
	Roles   Roles
	// BracketColors are the six bracket-highlight colors in nesting order.
	BracketColors []color.Color
}

// Build turns a config into a Built theme.
func Build(cfg Config) *Built {
	pl := color.NewPalette(cfg.Props)
	fg := pl.Get(color.KeyFG)
	coShades := pl.Get(color.KeyCoShades)
	accent := pl.Get(color.KeyAccent)
	coAccent := pl.Get(color.KeyCoAccent)
	red := pl.Get(color.KeyRed)

	brackets := make([]color.Color, 0, len(cfg.BracketGrades))
	for _, k := range cfg.BracketGrades {
		brackets = append(brackets, pl.Get(k)[6])
	}

	return &Built{
		Config:  cfg,
		Palette: pl,
		Roles: Roles{
			Identifier:    Role{Color: fg[8]},
			Parameter:     Role{Color: fg[8], Font: FontItalic},
			Property:      Role{Color: coShades[7]},
			Comment:       Role{Color: fg[5]},
			Operator:      Role{Color: accent[10]},
			WeakOperator:  Role{Color: coShades[6]},
			Literal:       Role{Color: accent[9]},
			Quote:         Role{Color: accent[7], Font: FontItalic},
			Declare:       Role{Color: fg[9]},
			Method:        Role{Color: fg[9]},
			TypeName:      Role{Color: coAccent[8]},
			Namespace:     Role{Color: coShades[9]},
			BuiltInType:   Role{Color: coAccent[9]},
			TypeParameter: Role{Color: coAccent[8], Font: FontItalic},
			Punctuation:   Role{Color: fg[5]},
			WeakKeyword:   Role{Color: fg[10]},
			Keyword:       Role{Color: fg[10], Font: FontBold},
			Storage:       Role{Color: fg[10], Font: FontBold},
			Library:       Role{Color: fg[9]},
			Invalid:       Role{Color: red[7]},
		},
		BracketColors: brackets,
	}
}
