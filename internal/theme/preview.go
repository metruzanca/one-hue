package theme

import (
	"github.com/metruzanca/one-hue-theme/internal/color"
)

// RoleColor is one derived syntax role in a form the web preview can consume.
type RoleColor struct {
	// Hex is the #RRGGBB color.
	Hex string
	// Font is the emphasis: "", "italic" or "bold".
	Font string
}

// Preview is the JSON-able payload served to the web app: the editor chrome,
// the derived syntax roles and the full palette, all derived from the accent.
type Preview struct {
	Name    string               `json:"name"`
	Slug    string               `json:"slug"`
	Accent  string               `json:"accent"`
	Dark    bool                 `json:"dark"`
	Editor  PreviewEditor        `json:"editor"`
	Roles   map[string]RoleColor `json:"roles"`
	Palette map[string][]string  `json:"palette"`
}

// PreviewEditor holds the CodeMirror editor chrome colors.
type PreviewEditor struct {
	Background         string   `json:"background"`
	Foreground         string   `json:"foreground"`
	Cursor             string   `json:"cursor"`
	Selection          string   `json:"selection"`
	SelectionHighlight string   `json:"selectionHighlight"`
	SearchMatch        string   `json:"searchMatch"`
	BracketMatch       string   `json:"bracketMatch"`
	LineNumber         string   `json:"lineNumber"`
	ActiveLineNumber   string   `json:"activeLineNumber"`
	ActiveLine         string   `json:"activeLine"`
	Gutter             string   `json:"gutter"`
	Brackets           []string `json:"brackets"`
}

// PreviewOf builds a Preview payload from a built theme.
func PreviewOf(b *Built) Preview {
	pl := b.Palette
	fg := pl.Get(color.KeyFG)
	bg := pl.Get(color.KeyBG)
	accent := pl.Get(color.KeyAccent)
	orange := pl.Get(color.KeyOrange)

	r := b.Roles
	roles := map[string]RoleColor{}
	for key, role := range map[string]Role{
		"identifier": r.Identifier, "parameter": r.Parameter,
		"property": r.Property, "comment": r.Comment,
		"operator": r.Operator, "weakOperator": r.WeakOperator,
		"literal": r.Literal, "quote": r.Quote,
		"declare": r.Declare, "method": r.Method,
		"typeName": r.TypeName, "namespace": r.Namespace,
		"builtInType": r.BuiltInType, "typeParameter": r.TypeParameter,
		"punctuation": r.Punctuation, "weakKeyword": r.WeakKeyword,
		"keyword": r.Keyword, "storage": r.Storage,
		"library": r.Library, "invalid": r.Invalid,
	} {
		roles[key] = RoleColor{Hex: role.Color.Hex(), Font: string(role.Font)}
	}

	brackets := make([]string, 0, len(b.BracketColors))
	for _, c := range b.BracketColors {
		brackets = append(brackets, c.Hex())
	}

	palette := map[string][]string{}
	for _, k := range color.PaletteKeys {
		grade := make([]string, 0, len(pl.Get(k)))
		for _, c := range pl.Get(k) {
			grade = append(grade, c.Hex())
		}
		palette[string(k)] = grade
	}

	return Preview{
		Name:   b.Config.Name,
		Slug:   b.Config.Slug,
		Accent: b.Config.Accent,
		Dark:   pl.IsDark(),
		Editor: PreviewEditor{
			Background:         bg[0].Hex(),
			Foreground:         fg[8].Hex(),
			Cursor:             accent[10].Hex(),
			Selection:          accent[7].Alpha(SelectionAlpha).HexAA(),
			SelectionHighlight: accent[4].Alpha(HighlightAlpha).HexAA(),
			SearchMatch:        orange[9].Alpha(SelectionAlpha).HexAA(),
			BracketMatch:       bg[4].Alpha(HighlightAlpha).HexAA(),
			LineNumber:         fg[5].Hex(),
			ActiveLineNumber:   fg[10].Hex(),
			ActiveLine:         bg[1].Hex(),
			Gutter:             bg[0].Hex(),
			Brackets:           brackets,
		},
		Roles:   roles,
		Palette: palette,
	}
}
