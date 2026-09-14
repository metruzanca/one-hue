package gen

import (
	"github.com/metruzanca/one-hue-theme/internal/color"
	"github.com/metruzanca/one-hue-theme/internal/theme"
)

// vscodeTarget renders one theme file per built theme in VSCode's color-theme
// JSON format.
type vscodeTarget struct{}

// Vscode returns the VSCode output target.
func Vscode() Target { return vscodeTarget{} }

func (vscodeTarget) Name() string { return "vscode" }

func (v vscodeTarget) Render(all []*theme.Built) ([]File, error) {
	var files []File
	for _, b := range all {
		files = append(files, File{
			Path:    "vscode/" + b.Config.Slug + ".json",
			Content: Encode(v.buildTheme(b)),
		})
	}
	return files, nil
}

func (v vscodeTarget) buildTheme(b *theme.Built) Value {
	pl := b.Palette
	fg := pl.Get(color.KeyFG)
	bg := pl.Get(color.KeyBG)
	red := pl.Get(color.KeyRed)
	orange := pl.Get(color.KeyOrange)
	yellow := pl.Get(color.KeyYellow)
	green := pl.Get(color.KeyGreen)
	cyan := pl.Get(color.KeyCyan)
	blue := pl.Get(color.KeyBlue)
	accent := pl.Get(color.KeyAccent)

	primaryBg := bg[0]
	primaryFg := fg[8]
	border := fg[3]
	r := b.Roles
	abs := pl.AbsGrade

	return []Field{
		{"$schema", "vscode://schemas/color-theme"},
		{"name", b.Config.Name},
		{"colors", []Field{
			{"focusBorder", accent[5].Hex()},
			{"foreground", primaryFg.Hex()},
			{"errorForeground", red[6].Hex()},
			{"editor.background", primaryBg.Hex()},
			{"editor.foreground", primaryFg.Hex()},
			{"editor.lineHighlightBackground", bg[1].Hex()},
			{"editorCursor.foreground", accent[10].Hex()},
			{"editorLineNumber.foreground", fg[5].Hex()},
			{"editorCodeLens.foreground", fg[5].Hex()},
			{"editorActiveLineNumber.foreground", fg[10].Hex()},

			{"editor.selectionBackground", accent[7].Alpha(theme.SelectionAlpha).HexAA()},
			{"editor.inactiveSelectionBackground", bg[7].Alpha(theme.SelectionAlpha).HexAA()},
			{"editor.selectionHighlightBackground", accent[4].Alpha(theme.HighlightAlpha).HexAA()},
			{"editor.wordHighlightBackground", bg[5].Alpha(theme.HighlightAlpha).HexAA()},
			{"editor.wordHighlightStrongBackground", accent[5].Alpha(theme.HighlightAlpha).HexAA()},
			{"editor.findMatchBackground", orange[9].Alpha(theme.SelectionAlpha).HexAA()},
			{"editor.findMatchHighlightBackground", orange[7].Alpha(theme.HighlightAlpha).HexAA()},
			{"editor.findRangeHighlightBackground", bg[4].Alpha(theme.HighlightAlpha).HexAA()},
			{"editorLink.activeForeground", accent[10].Hex()},
			{"editorBracketMatch.background", bg[4].Alpha(theme.HighlightAlpha).HexAA()},
			{"editorBracketMatch.border", fg[4].Alpha(theme.HighlightBorderAlpha).HexAA()},

			{"minimap.findMatchHighlight", orange[7].Alpha(theme.HighlightAlpha).HexAA()},
			{"minimap.selectionHighlight", accent[7].Alpha(theme.SelectionAlpha).HexAA()},

			{"editorError.foreground", red[8].Hex()},
			{"editorWarning.foreground", yellow[8].Hex()},
			{"editorInfo.foreground", blue[8].Hex()},
			{"editorHint.foreground", blue[8].Hex()},

			{"notification.errorBackground", red[4].Hex()},
			{"notification.errorForeground", red[8].Hex()},
			{"notification.infoBackground", blue[4].Hex()},
			{"notification.infoForeground", blue[8].Hex()},
			{"notification.warningBackground", yellow[4].Hex()},
			{"notification.warningForeground", yellow[8].Hex()},

			{"debugToolBar.background", bg[2].HexAA()},
			{"editorWidget.background", bg[1].HexAA()},
			{"editorSuggestWidget.background", primaryBg.HexAA()},

			{"editorGroup.border", fg[1].Hex()},
			{"editorGroupHeader.noTabsBackground", primaryBg.Hex()},
			{"editorGroupHeader.tabsBackground", primaryBg.Hex()},
			{"tab.border", primaryBg.Hex()},
			{"tab.inactiveBackground", primaryBg.Hex()},
			{"tab.inactiveForeground", fg[7].Hex()},
			{"tab.activeBackground", primaryBg.Hex()},
			{"tab.activeForeground", fg[10].Hex()},
			{"tab.hoverBackground", bg[1].Hex()},
			{"tab.hoverBorder", accent[5].Hex()},
			{"tab.activeBorder", accent[8].Hex()},

			{"peekView.border", fg[2].Hex()},
			{"peekViewTitle.background", bg[1].HexAA()},
			{"peekViewEditor.background", bg[3].Alpha(1.0 / 6.0).HexAA()},

			{"scrollbar.shadow", color.RGB(0, 0, 0).Alpha(0.1).HexAA()},
			{"scrollbarSlider.background", bg[10].Alpha(0.075).HexAA()},
			{"scrollbarSlider.activeBackground", bg[10].Alpha(0.15).HexAA()},
			{"scrollbarSlider.hoverBackground", bg[10].Alpha(0.15).HexAA()},

			{"editorOverviewRuler.border", "#00000000"},
			{"editorGutter.modifiedBackground", cyan[4].Hex()},
			{"editorGutter.addedBackground", green[4].Hex()},
			{"editorGutter.deletedBackground", red[4].Hex()},
			{"editorOverviewRuler.modifiedForeground", cyan[10].Alpha(0.2).HexAA()},
			{"editorOverviewRuler.addedForeground", cyan[10].Alpha(0.2).HexAA()},
			{"editorOverviewRuler.deletedForeground", cyan[10].Alpha(0.2).HexAA()},
			{"editorOverviewRuler.infoForeground", blue[5].Alpha(0.9).HexAA()},
			{"editorOverviewRuler.warningForeground", yellow[5].Alpha(0.9).HexAA()},
			{"editorOverviewRuler.errorForeground", red[5].Alpha(0.9).HexAA()},
			{"editorOverviewRuler.findMatchForeground", orange[5].Alpha(0.9).HexAA()},
			{"editorOverviewRuler.bracketMatchForeground", fg[10].Alpha(0.3).HexAA()},
			{"editorOverviewRuler.wordHighlightForeground", fg[10].Alpha(0.3).HexAA()},
			{"editorOverviewRuler.wordHighlightStrongForeground", fg[10].Alpha(0.3).HexAA()},
			{"diffEditor.removedTextBackground", red[5].Alpha(0.15).HexAA()},
			{"diffEditor.insertedTextBackground", green[5].Alpha(0.1).HexAA()},

			{"editorBracketHighlight.foreground1", b.BracketColors[0].Hex()},
			{"editorBracketHighlight.foreground2", b.BracketColors[1].Hex()},
			{"editorBracketHighlight.foreground3", b.BracketColors[2].Hex()},
			{"editorBracketHighlight.foreground4", b.BracketColors[3].Hex()},
			{"editorBracketHighlight.foreground5", b.BracketColors[4].Hex()},
			{"editorBracketHighlight.foreground6", b.BracketColors[5].Hex()},

			{"textLink.foreground", accent[7].Hex()},
			{"textLink.activeForeground", accent[8].Hex()},

			{"sideBarTitle.foreground", fg[10].Hex()},
			{"sideBar.background", bg[1].Hex()},
			{"sideBarSectionHeader.background", bg[3].Hex()},

			{"list.highlightForeground", accent[6].Hex()},
			{"list.hoverBackground", bg[3].Hex()},
			{"list.inactiveSelectionBackground", bg[5].Alpha(0.3).HexAA()},
			{"list.activeSelectionBackground", accent[5].Alpha(0.25).HexAA()},
			{"list.focusBackground", accent[5].Alpha(0.25).HexAA()},
			{"list.inactiveSelectionForeground", fg[10].Hex()},
			{"list.activeSelectionForeground", fg[10].Hex()},
			{"list.focusForeground", fg[10].Hex()},
			{"list.errorForeground", red[8].Hex()},
			{"list.warningForeground", yellow[8].Hex()},

			{"gitDecoration.modifiedResourceForeground", cyan[8].Hex()},
			{"gitDecoration.addedResourceForeground", green[8].Hex()},
			{"gitDecoration.untrackedResourceForeground", green[8].Hex()},
			{"gitDecoration.deletedResourceForeground", red[8].Hex()},
			{"gitDecoration.conflictingResourceForeground", red[8].Hex()},

			{"dropdown.background", bg[1].Hex()},
			{"dropdown.border", border.Hex()},
			{"dropdown.foreground", fg[8].Hex()},

			{"input.background", bg[0].Hex()},
			{"input.border", border.Hex()},
			{"input.foreground", fg[8].Hex()},
			{"input.placeholderForeground", fg[5].Hex()},
			{"inputOption.activeBorder", accent[6].Hex()},

			{"button.background", accent[abs(6)].Hex()},
			{"button.foreground", fg[abs(10)].Hex()},
			{"button.hoverBackground", accent[abs(7)].Hex()},
			{"badge.background", accent[abs(4)].Hex()},
			{"badge.foreground", fg[abs(10)].Hex()},
			{"activityBarBadge.background", accent[abs(4)].Hex()},
			{"activityBarBadge.foreground", fg[abs(10)].Hex()},
			{"activityBar.background", bg[3].Hex()},
			{"activityBar.foreground", fg[9].Hex()},

			{"statusBar.background", bg[2].Hex()},
			{"statusBar.foreground", primaryFg.Hex()},
			{"statusBar.noFolderBackground", bg[2].Hex()},
			{"statusBar.noFolderForeground", primaryFg.Hex()},
			{"statusBar.debuggingBackground", bg[2].Hex()},
			{"statusBar.debuggingForeground", primaryFg.Hex()},

			{"panel.border", border.Hex()},
			{"panelTitle.activeBorder", accent[5].Hex()},

			{"titleBar.activeBackground", bg[abs(2)].Hex()},
			{"titleBar.activeForeground", fg[abs(9)].Hex()},
			{"titleBar.inactiveBackground", bg[abs(3)].Hex()},
			{"titleBar.inactiveForeground", fg[abs(8)].Hex()},

			{"breadcrumb.foreground", r.Comment.Color.Hex()},
			{"breadcrumb.focusForeground", r.Quote.Color.Hex()},
			{"breadcrumb.activeSelectionForeground", r.Quote.Color.Hex()},
			{"breadcrumbPicker.background", bg[1].Hex()},
		}},
		{"tokenColors", v.tokenColors(b)},
		{"semanticHighlighting", true},
	}
}

func (v vscodeTarget) tokenColors(b *theme.Built) Value {
	r := b.Roles

	scopeStr := func(s string) Value { return s }
	scopeList := func(s ...string) Value {
		out := make([]Value, len(s))
		for i, v := range s {
			out[i] = v
		}
		return out
	}
	settings := func(role theme.Role, extra ...Field) []Field {
		var s []Field
		if role.Font != theme.FontNormal {
			s = append(s, Field{"fontStyle", string(role.Font)})
		}
		s = append(s, Field{"foreground", role.Color.Hex()})
		return append(s, extra...)
	}

	return []Value{
		[]Field{ // Identifier
			{"name", "Identifier"},
			{"scope", scopeList(
				"variable",
				"meta.definition.variable.name",
				"support.variable",
				"variable.other.readwrite",
				"variable.other.constant",
				"variable.other.readonly",
			)},
			{"settings", settings(r.Identifier)},
		},
		[]Field{ // Parameter
			{"name", "Parameter"},
			{"scope", scopeList("variable.parameter", "variable.other.local", "variable.other.constant")},
			{"settings", settings(r.Parameter)},
		},
		[]Field{ // Member Access
			{"name", "Member Access"},
			{"scope", scopeList("variable.other.property", "variable.other.constant.property")},
			{"settings", settings(r.Property)},
		},
		[]Field{ // Object keys, TS grammar specific
			{"name", "Object keys, TS grammar specific"},
			{"scope", scopeList("meta.object-literal.key", "meta.object-literal.key entity.name.function")},
			{"settings", settings(r.Identifier)},
		},
		[]Field{ // Comment
			{"name", "Comment"},
			{"scope", scopeList("comment", "punctuation.comment", "punctuation.definition.comment")},
			{"settings", settings(r.Comment)},
		},
		[]Field{ // Operator
			{"name", "Operator"},
			{"scope", scopeList("keyword.operator")},
			{"settings", settings(r.Operator)},
		},
		[]Field{ // Punctuation
			{"name", "Punctuation"},
			{"scope", scopeList(
				"punctuation",
				"delimiter",
				"bracket",
				"brace",
				"paren",
				"delimiter.tag",
				"punctuation.tag",
				"tag.html",
				"tag.xml",
				"meta.property-value punctuation.separator.key-value",
				"punctuation.definition.metadata.md",
				"string.link.md",
				"meta.brace",
			)},
			{"settings", settings(r.Punctuation)},
		},
		[]Field{ // JavaScript string interpolation ${}
			{"name", "JavaScript string interpolation ${}"},
			{"scope", scopeList(
				"punctuation.definition.template-expression.begin.js",
				"punctuation.definition.template-expression.begin.ts",
				"punctuation.definition.template-expression.end.ts",
				"punctuation.definition.template-expression.end.js",
				"punctuation.section.embedded.begin.metatag.php",
				"punctuation.section.embedded.end.metatag.php",
				"punctuation.definition.typeparameters",
			)},
			{"settings", settings(r.WeakOperator)},
		},
		[]Field{ // String
			{"name", "String"},
			{"scope", scopeList(
				"string",
				"meta.property-value.string",
				"support.constant.property-value.string",
				"meta.structure.dictionary.value.json string.quoted.double.json",
				"meta.structure.dictionary.json string.quoted.double.json",
				"meta.preprocessor string",
			)},
			{"settings", settings(r.Quote)},
		},
		[]Field{ // template expression: reset styles
			{"scope", scopeList("meta.template.expression.ts")},
			{"settings", []Field{{"fontStyle", ""}}},
		},
		[]Field{ // Primitive Literals
			{"name", "Primitive Literals"},
			{"scope", scopeList(
				"constant.numeric",
				"constant.dec.numeric",
				"constant.hex.numeric",
				"meta.property-value.numeric",
				"support.constant.property-value.numeric",
				"meta.property-value.color",
				"support.constant.property-value.color",
				"constant.language",
			)},
			{"settings", settings(r.Literal)},
		},
		[]Field{ // Namespace
			{"name", "Namespace"},
			{"scope", scopeList("entity.name.namespace", "entity.name.type.module")},
			{"settings", settings(r.Namespace)},
		},
		[]Field{ // Type name
			{"name", "Type name"},
			{"scope", scopeList("entity.name.type")},
			{"settings", settings(r.TypeName)},
		},
		[]Field{ // Built-in type name
			{"name", "Built-in type name"},
			{"scope", scopeList("support.type")},
			{"settings", settings(r.BuiltInType)},
		},
		[]Field{ // Type parameter
			{"name", "Type parameter"},
			{"scope", scopeList("entity.name.type.parameter")},
			{"settings", settings(r.TypeParameter)},
		},
		[]Field{ // User names
			{"name", "User names"},
			{"scope", scopeList(
				"constant.character",
				"constant.other",
				"entity.name",
				"entity.name.class",
				"entity.name.function",
				"entity.other.inherited-class",
				"entity.other.attribute-name",
				"entity.other.attribute-name",
				"entity.other.attribute-name.html",
				"support.type.property-name",
				"string.key",
				"entity.name.tag.table",
				"meta.structure.dictionary.json string.quoted.double.json",
			)},
			{"settings", settings(r.Declare)},
		},
		[]Field{ // Methods
			{"name", "Methods"},
			{"scope", scopeList("entity.name.function.member")},
			{"settings", settings(r.Method)},
		},
		[]Field{ // Weak Keyword
			{"name", "Weak Keyword"},
			{"scope", scopeList("variable.language.this")},
			{"settings", settings(r.WeakKeyword)},
		},
		[]Field{ // Keyword
			{"name", "Keyword"},
			{"scope", scopeList(
				"keyword",
				"meta.property-value.keyword",
				"support.constant.property-value.keyword",
				"meta.preprocessor.keyword",
				"keyword.other.use",
				"keyword.other.function.use",
				"keyword.other.namespace",
				"keyword.other.new",
				"keyword.other.special-method",
				"keyword.other.unit",
				"keyword.other.use-as",
			)},
			{"settings", settings(r.Keyword)},
		},
		[]Field{ // Storage
			{"name", "Storage"},
			{"scope", scopeList(
				"storage",
				"storage.type",
				"storage.type.ts",
				"storage.type.var.ts",
				"storage.type.js",
				"storage.type.var.js",
				"storage.type.const.ts",
				"storage.type.let.ts",
				"storage.type.let.js",
				"storage.type.const.js",
				"entity.name.tag",
			)},
			{"settings", settings(r.Storage)},
		},
		[]Field{ // Pointer, access, etc
			{"name", "Pointer, access, etc"},
			{"scope", scopeList("meta.ptr", "meta.pointer", "meta.address", "meta.array.cxx")},
			{"settings", []Field{}},
		},
		[]Field{ // Preprocessor
			{"name", "Preprocessor"},
			{"scope", scopeStr("meta.preprocessor")},
			{"settings", settings(r.Declare)},
		},
		[]Field{ // Library
			{"name", "Library"},
			{"scope", scopeList("support.class", "support.function", "support.constant")},
			{"settings", settings(r.Library)},
		},
		[]Field{ // Invalid
			{"name", "Invalid"},
			{"scope", scopeStr("invalid")},
			{"settings", settings(r.Invalid)},
		},
		[]Field{ // Invalid deprecated
			{"name", "Invalid deprecated"},
			{"scope", scopeList("invalid.deprecated")},
			{"settings", settings(r.Invalid)},
		},
		[]Field{ // Markdown Title Hash
			{"name", "Markdown Title Hash"},
			{"scope", scopeList("punctuation.definition.heading.md", "entity.name.type.md", "beginning.punctuation")},
			{"settings", settings(r.Declare)},
		},
		[]Field{ // Markdown titles
			{"name", "Markdown titles"},
			{"scope", scopeList("markup.heading", "entity.name.section")},
			{"settings", settings(r.Keyword)},
		},
		[]Field{ // Markdown Raw
			{"name", "Markdown Raw"},
			{"scope", scopeList("markup.raw", "markup.inline.raw", "markup.fenced", "markup.fenced_code")},
			{"settings", settings(r.Quote)},
		},
		[]Field{ // Markdown link
			{"name", "Markdown link"},
			{"scope", scopeList(
				"markup.link",
				"string.other.link.title",
				"string.other.link.description",
				"meta.link.inline",
				"meta.image.inline",
			)},
			{"settings", settings(r.Declare)},
		},
		[]Field{ // Makefile Variables
			{"name", "Makefile Variables"},
			{"scope", scopeList("variable.language.makefile", "variable.other.makefile")},
			{"settings", settings(r.Declare)},
		},
		[]Field{ // markup.italic
			{"scope", scopeList("markup.italic")},
			{"settings", []Field{{"fontStyle", "italic"}}},
		},
		[]Field{ // markup.bold
			{"scope", scopeList("markup.bold")},
			{"settings", []Field{{"fontStyle", "bold"}}},
		},
		[]Field{ // CSS Class
			{"name", "CSS Class"},
			{"scope", scopeList("entity.other.attribute-name.class.css")},
			{"settings", settings(r.Library)},
		},
		[]Field{ // CSS Tag name
			{"name", "CSS Tag name"},
			{"scope", scopeList("entity.name.tag.css")},
			{"settings", settings(r.Keyword)},
		},
		[]Field{ // CSS Property
			{"name", "CSS Property"},
			{"scope", scopeList("meta.property-name.css")},
			{"settings", settings(r.Declare)},
		},
	}
}

var _ Target = vscodeTarget{}
