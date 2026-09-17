package gen

import (
	"github.com/metruzanca/one-hue-theme/internal/color"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

// opencodeTarget renders one theme file per built theme in OpenCode's custom
// theme JSON format. Each file is placed in the user's themes directory as-is.
type opencodeTarget struct{}

// Opencode returns the OpenCode output target.
func Opencode() Target { return opencodeTarget{} }

func (opencodeTarget) Name() string { return "opencode" }

func (o opencodeTarget) Render(all []*theme.Built) ([]File, error) {
	var files []File
	for _, b := range all {
		files = append(files, File{
			Path:    "opencode/" + b.Config.Slug + ".json",
			Content: Encode(o.buildTheme(b)),
		})
	}
	return files, nil
}

// buildTheme maps the one-hue palette and syntax roles onto OpenCode's theme
// keys. Values are plain hex: the theme is dark-only, so light mode keeps the
// same colors.
func (o opencodeTarget) buildTheme(b *theme.Built) Value {
	pl := b.Palette
	fg := pl.Get(color.KeyFG)
	bg := pl.Get(color.KeyBG)
	accent := pl.Get(color.KeyAccent)
	coAccent := pl.Get(color.KeyCoAccent)
	red := pl.Get(color.KeyRed)
	yellow := pl.Get(color.KeyYellow)
	green := pl.Get(color.KeyGreen)
	blue := pl.Get(color.KeyBlue)
	r := b.Roles

	return []Field{
		{"$schema", "https://opencode.ai/theme.json"},
		{"theme", []Field{
			{"primary", accent[7].Hex()},
			{"secondary", coAccent[7].Hex()},
			{"accent", accent[6].Hex()},
			{"error", red[7].Hex()},
			{"warning", yellow[7].Hex()},
			{"success", green[7].Hex()},
			{"info", blue[7].Hex()},

			{"text", fg[8].Hex()},
			{"textMuted", fg[5].Hex()},
			{"background", bg[0].Hex()},
			{"backgroundPanel", bg[1].Hex()},
			{"backgroundElement", bg[2].Hex()},
			{"border", fg[3].Hex()},
			{"borderActive", accent[5].Hex()},
			{"borderSubtle", fg[2].Hex()},

			{"diffAdded", green[8].Hex()},
			{"diffRemoved", red[8].Hex()},
			{"diffContext", fg[8].Hex()},
			{"diffHunkHeader", accent[8].Hex()},
			{"diffHighlightAdded", green[7].Hex()},
			{"diffHighlightRemoved", red[7].Hex()},
			{"diffAddedBg", green[4].Hex()},
			{"diffRemovedBg", red[4].Hex()},
			{"diffContextBg", bg[1].Hex()},
			{"diffLineNumber", fg[5].Hex()},
			{"diffAddedLineNumberBg", green[4].Hex()},
			{"diffRemovedLineNumberBg", red[4].Hex()},

			{"markdownText", fg[8].Hex()},
			{"markdownHeading", accent[8].Hex()},
			{"markdownLink", blue[8].Hex()},
			{"markdownLinkText", accent[7].Hex()},
			{"markdownCode", r.Quote.Color.Hex()},
			{"markdownBlockQuote", fg[5].Hex()},
			{"markdownEmph", accent[7].Hex()},
			{"markdownStrong", fg[10].Hex()},
			{"markdownHorizontalRule", fg[3].Hex()},
			{"markdownListItem", accent[7].Hex()},
			{"markdownListEnumeration", coAccent[8].Hex()},
			{"markdownImage", blue[8].Hex()},
			{"markdownImageText", fg[6].Hex()},
			{"markdownCodeBlock", fg[8].Hex()},

			{"syntaxComment", r.Comment.Color.Hex()},
			{"syntaxKeyword", r.Keyword.Color.Hex()},
			{"syntaxFunction", r.Method.Color.Hex()},
			{"syntaxVariable", r.Identifier.Color.Hex()},
			{"syntaxString", r.Quote.Color.Hex()},
			{"syntaxNumber", r.Literal.Color.Hex()},
			{"syntaxType", r.TypeName.Color.Hex()},
			{"syntaxOperator", r.Operator.Color.Hex()},
			{"syntaxPunctuation", r.Punctuation.Color.Hex()},
		}},
	}
}

var _ Target = opencodeTarget{}
