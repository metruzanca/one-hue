package gen

import (
	"github.com/metruzanca/one-hue-theme/internal/color"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

// zedTarget renders all themes into one Zed theme-family file.
type zedTarget struct{}

// Zed returns the Zed editor output target.
func Zed() Target { return zedTarget{} }

func (zedTarget) Name() string { return "zed" }

func (z zedTarget) Render(all []*theme.Built) ([]File, error) {
	themes := make([]Value, 0, len(all))
	for _, b := range all {
		themes = append(themes, z.buildTheme(b))
	}
	return []File{{
		Path: "zed/one-hue.json",
		Content: Encode([]Field{
			{"$schema", "https://zed.dev/schema/themes/v0.2.0.json"},
			{"name", "One Hue Theme"},
			{"author", "metruzanca"},
			{"themes", themes},
		}),
	}}, nil
}

func (z zedTarget) buildTheme(b *theme.Built) Value {
	pl := b.Palette
	fg := pl.Get(color.KeyFG)
	bg := pl.Get(color.KeyBG)
	red := pl.Get(color.KeyRed)
	orange := pl.Get(color.KeyOrange)
	yellow := pl.Get(color.KeyYellow)
	green := pl.Get(color.KeyGreen)
	cyan := pl.Get(color.KeyCyan)
	blue := pl.Get(color.KeyBlue)
	violet := pl.Get(color.KeyViolet)
	magenta := pl.Get(color.KeyMagenta)
	accent := pl.Get(color.KeyAccent)

	primaryBg := bg[0]
	primaryFg := fg[8]
	border := fg[3]

	// hex is the Zed 8-digit form: #RRGGBBAA.
	hex := func(c color.Color) string { return c.HexAA() }
	// ahex renders a color with the given opacity multiplier.
	ahex := func(c color.Color, a float64) string { return c.Alpha(a).HexAA() }

	appearance := "dark"
	if !pl.IsDark() {
		appearance = "light"
	}

	return []Field{
		{"name", b.Config.Name},
		{"appearance", appearance},
		{"style", []Field{
			// Surfaces & borders
			{"border", hex(border)},
			{"border.variant", hex(fg[2])},
			{"border.focused", hex(accent[5])},
			{"border.selected", hex(accent[4])},
			{"border.transparent", "#00000000"},
			{"border.disabled", hex(fg[2])},
			{"elevated_surface.background", hex(bg[2])},
			{"surface.background", hex(bg[1])},
			{"background", hex(primaryBg)},
			{"element.background", hex(bg[1])},
			{"element.hover", hex(bg[2])},
			{"element.active", hex(bg[3])},
			{"element.selected", hex(bg[3])},
			{"element.disabled", hex(bg[1])},
			{"drop_target.background", ahex(accent[7], 0.5)},
			{"ghost_element.background", "#00000000"},
			{"ghost_element.hover", hex(bg[2])},
			{"ghost_element.active", hex(bg[3])},
			{"ghost_element.selected", hex(bg[3])},
			{"ghost_element.disabled", hex(bg[1])},

			// Text & icons
			{"text", hex(primaryFg)},
			{"text.muted", hex(fg[5])},
			{"text.placeholder", hex(fg[4])},
			{"text.disabled", hex(fg[3])},
			{"text.accent", hex(accent[7])},
			{"icon", hex(primaryFg)},
			{"icon.muted", hex(fg[5])},
			{"icon.disabled", hex(fg[3])},
			{"icon.placeholder", hex(fg[5])},
			{"icon.accent", hex(accent[7])},

			// Chrome
			{"status_bar.background", hex(bg[2])},
			{"title_bar.background", hex(bg[2])},
			{"title_bar.inactive_background", hex(bg[3])},
			{"toolbar.background", hex(bg[1])},
			{"tab_bar.background", hex(bg[1])},
			{"tab.inactive_background", hex(primaryBg)},
			{"tab.active_background", hex(bg[1])},
			{"search.match_background", ahex(orange[7], 0.25)},
			{"search.active_match_background", ahex(orange[9], 0.3)},
			{"panel.background", hex(bg[1])},
			{"panel.focused_border", nil},
			{"pane.focused_border", nil},
			{"scrollbar.thumb.background", ahex(bg[10], 0.075)},
			{"scrollbar.thumb.hover_background", ahex(bg[10], 0.15)},
			{"scrollbar.thumb.border", "#00000000"},
			{"scrollbar.track.background", "#00000000"},
			{"scrollbar.track.border", hex(bg[2])},

			// Editor
			{"editor.foreground", hex(primaryFg)},
			{"editor.background", hex(primaryBg)},
			{"editor.gutter.background", hex(primaryBg)},
			{"editor.subheader.background", hex(bg[1])},
			{"editor.active_line.background", ahex(bg[1], 0.75)},
			{"editor.highlighted_line.background", hex(bg[1])},
			{"editor.line_number", hex(fg[5])},
			{"editor.active_line_number", hex(fg[10])},
			{"editor.hover_line_number", hex(fg[8])},
			{"editor.invisible", hex(fg[3])},
			{"editor.wrap_guide", ahex(fg[1], 0.05)},
			{"editor.active_wrap_guide", ahex(fg[2], 0.1)},
			{"editor.document_highlight.read_background", ahex(accent[4], theme.HighlightAlpha)},
			{"editor.document_highlight.write_background", ahex(accent[5], 0.4)},
			{"editor.document_highlight.bracket_background", ahex(bg[4], theme.HighlightAlpha)},

			// Terminal
			{"terminal.background", hex(primaryBg)},
			{"terminal.foreground", hex(primaryFg)},
			{"terminal.bright_foreground", hex(fg[10])},
			{"terminal.dim_foreground", hex(fg[5])},
			{"terminal.ansi.black", hex(bg[0])},
			{"terminal.ansi.bright_black", hex(bg[8])},
			{"terminal.ansi.dim_black", hex(bg[2])},
			{"terminal.ansi.red", hex(red[7])},
			{"terminal.ansi.bright_red", hex(red[9])},
			{"terminal.ansi.dim_red", hex(red[5])},
			{"terminal.ansi.green", hex(green[7])},
			{"terminal.ansi.bright_green", hex(green[9])},
			{"terminal.ansi.dim_green", hex(green[5])},
			{"terminal.ansi.yellow", hex(yellow[7])},
			{"terminal.ansi.bright_yellow", hex(yellow[9])},
			{"terminal.ansi.dim_yellow", hex(yellow[5])},
			{"terminal.ansi.blue", hex(blue[7])},
			{"terminal.ansi.bright_blue", hex(blue[9])},
			{"terminal.ansi.dim_blue", hex(blue[5])},
			{"terminal.ansi.magenta", hex(magenta[7])},
			{"terminal.ansi.bright_magenta", hex(magenta[9])},
			{"terminal.ansi.dim_magenta", hex(magenta[5])},
			{"terminal.ansi.cyan", hex(cyan[7])},
			{"terminal.ansi.bright_cyan", hex(cyan[9])},
			{"terminal.ansi.dim_cyan", hex(cyan[5])},
			{"terminal.ansi.white", hex(fg[8])},
			{"terminal.ansi.bright_white", hex(fg[10])},
			{"terminal.ansi.dim_white", hex(fg[5])},

			{"link_text.hover", hex(accent[8])},

			// Version control
			{"version_control.added", hex(green[6])},
			{"version_control.modified", hex(yellow[6])},
			{"version_control.word_added", ahex(green[7], 0.35)},
			{"version_control.word_deleted", ahex(red[7], 0.8)},
			{"version_control.deleted", hex(red[6])},
			{"version_control.conflict_marker.ours", ahex(green[7], 0.1)},
			{"version_control.conflict_marker.theirs", ahex(blue[7], 0.1)},

			// Status colors
			{"conflict", hex(red[6])},
			{"conflict.background", ahex(red[4], 0.1)},
			{"conflict.border", hex(red[3])},
			{"created", hex(green[6])},
			{"created.background", ahex(green[4], 0.1)},
			{"created.border", hex(green[3])},
			{"deleted", hex(red[6])},
			{"deleted.background", ahex(red[4], 0.1)},
			{"deleted.border", hex(red[3])},
			{"error", hex(red[6])},
			{"error.background", ahex(red[4], 0.1)},
			{"error.border", hex(red[3])},
			{"hidden", hex(fg[5])},
			{"hidden.background", ahex(fg[5], 0.1)},
			{"hidden.border", hex(fg[2])},
			{"hint", hex(blue[6])},
			{"hint.background", ahex(blue[4], 0.1)},
			{"hint.border", hex(blue[3])},
			{"ignored", hex(fg[5])},
			{"ignored.background", ahex(fg[5], 0.1)},
			{"ignored.border", hex(fg[2])},
			{"info", hex(blue[6])},
			{"info.background", ahex(blue[4], 0.1)},
			{"info.border", hex(blue[3])},
			{"modified", hex(cyan[6])},
			{"modified.background", ahex(cyan[4], 0.1)},
			{"modified.border", hex(cyan[3])},
			{"predictive", hex(fg[5])},
			{"predictive.background", ahex(fg[5], 0.1)},
			{"predictive.border", hex(fg[2])},
			{"renamed", hex(blue[6])},
			{"renamed.background", ahex(blue[4], 0.1)},
			{"renamed.border", hex(blue[3])},
			{"success", hex(green[6])},
			{"success.background", ahex(green[4], 0.1)},
			{"success.border", hex(green[3])},
			{"unreachable", hex(fg[5])},
			{"unreachable.background", ahex(fg[5], 0.1)},
			{"unreachable.border", hex(fg[2])},
			{"warning", hex(yellow[6])},
			{"warning.background", ahex(yellow[4], 0.1)},
			{"warning.border", hex(yellow[3])},

			// Collaboration cursors
			{"players", z.players(accent, red, yellow, violet, green, cyan, magenta, orange)},

			// Syntax
			{"syntax", z.syntax(b, fg, red, green, blue)},
		}},
	}
}

func (z zedTarget) players(accent, red, yellow, violet, green, cyan, magenta, orange []color.Color) Value {
	palettes := [][]color.Color{accent, red, yellow, violet, green, cyan, magenta, orange}
	players := make([]Value, 0, len(palettes))
	for _, p := range palettes {
		c := p[7]
		players = append(players, []Field{
			{"cursor", c.HexAA()},
			{"background", c.HexAA()},
			{"selection", c.Alpha(0.24).HexAA()},
		})
	}
	return players
}

// syntax builds the Zed syntax map from the theme roles plus a few semantic
// colors for diffs.
func (z zedTarget) syntax(b *theme.Built, fg, red, green, blue []color.Color) Value {
	r := b.Roles

	// node takes a role and returns its [color, font_style, font_weight] as
	// positional values fit for emit.
	type syn struct {
		color  string
		bold   bool
		it     bool
		weight any
	}
	role := func(rl theme.Role) syn {
		s := syn{color: rl.Color.Hex()}
		switch rl.Font {
		case theme.FontBold:
			s.bold = true
		case theme.FontItalic:
			s.it = true
		}
		return s
	}
	emit := func(s syn) []Field {
		out := []Field{{"color", s.color}}
		if s.it {
			out = append(out, Field{"font_style", "italic"})
		}
		if s.bold {
			out = append(out, Field{"font_weight", 700})
		}
		return out
	}
	// plain colors without a role
	plain := func(c color.Color) []Field { return []Field{{"color", c.Hex()}} }

	kind := syn{color: fg[8].Hex()}
	kw := role(r.Keyword)
	st := role(r.Storage)
	op := role(r.Operator)
	weakOp := role(r.WeakOperator)
	com := role(r.Comment)
	lit := role(r.Literal)
	quote := role(r.Quote)
	decl := role(r.Declare)
	typeName := r.TypeName
	builtIn := r.BuiltInType
	param := role(r.Parameter)
	prop := role(r.Property)
	ident := role(r.Identifier)
	punct := role(r.Punctuation)

	return []Field{
		{"attribute", emit(decl)},
		{"boolean", emit(lit)},
		{"comment", emit(com)},
		{"comment.doc", emit(com)},
		{"constant", emit(lit)},
		{"constructor", emit(decl)},
		{"embedded", emit(ident)},
		{"emphasis", emit(ident)},
		{"emphasis.strong", emit(syn{color: kind.color, bold: true})},
		{"enum", emit(role(r.TypeName))},
		{"function", emit(decl)},
		{"hint", emit(role(theme.Role{Color: blue[6]}))},
		{"keyword", emit(kw)},
		{"label", emit(decl)},
		{"link_text", emit(role(theme.Role{Color: fg[9], Font: theme.FontItalic}))},
		{"link_uri", emit(lit)},
		{"namespace", emit(role(r.Namespace))},
		{"number", emit(lit)},
		{"operator", emit(op)},
		{"predictive", emit(syn{color: com.color, it: true})},
		{"preproc", emit(decl)},
		{"primary", emit(ident)},
		{"property", emit(prop)},
		{"punctuation", emit(punct)},
		{"punctuation.bracket", emit(punct)},
		{"punctuation.delimiter", emit(punct)},
		{"punctuation.list_marker", emit(punct)},
		{"punctuation.markup", emit(punct)},
		{"punctuation.special", emit(role(r.Operator))},
		{"selector", emit(decl)},
		{"selector.pseudo", emit(role(r.TypeName))},
		{"string", emit(quote)},
		{"string.escape", emit(weakOp)},
		{"string.regex", emit(lit)},
		{"string.special", emit(lit)},
		{"string.special.symbol", emit(lit)},
		{"tag", emit(st)},
		{"text.literal", emit(quote)},
		{"title", emit(kw)},
		{"type", emit(role(typeName))},
		{"type.builtin", emit(role(builtIn))},
		{"variable", emit(ident)},
		{"variable.parameter", emit(param)},
		{"variable.special", emit(lit)},
		{"variant", emit(role(typeName))},
		{"diff.plus", plain(green[8])},
		{"diff.minus", plain(red[8])},
	}
}
