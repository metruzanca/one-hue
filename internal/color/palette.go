package color

import "math"

// GRADES is the number of gradient intervals: every palette array has
// GRADES+1 entries, index 0 being the darkest in dark themes.
const GRADES = 10

// PaletteKey identifies one of the derived color arrays.
type PaletteKey string

// Palette keys.
const (
	KeyFG         PaletteKey = "fg"
	KeyBG         PaletteKey = "bg"
	KeyCoShades   PaletteKey = "coShades"
	KeyAccent     PaletteKey = "accent"
	KeyCoAccent   PaletteKey = "coAccent"
	KeyRed        PaletteKey = "red"
	KeyOrange     PaletteKey = "orange"
	KeyYellow     PaletteKey = "yellow"
	KeyChartreuse PaletteKey = "chartreuse"
	KeyGreen      PaletteKey = "green"
	KeySpring     PaletteKey = "spring"
	KeyCyan       PaletteKey = "cyan"
	KeyAzure      PaletteKey = "azure"
	KeyBlue       PaletteKey = "blue"
	KeyViolet     PaletteKey = "violet"
	KeyMagenta    PaletteKey = "magenta"
	KeyRose       PaletteKey = "rose"
)

// PaletteKeys lists the derived color array names in canonical order.
var PaletteKeys = []PaletteKey{
	KeyFG, KeyBG, KeyCoShades, KeyAccent, KeyCoAccent,
	KeyRed, KeyOrange, KeyYellow, KeyChartreuse, KeyGreen, KeySpring,
	KeyCyan, KeyAzure, KeyBlue, KeyViolet, KeyMagenta, KeyRose,
}

// PaletteProps is the input configuration for building a Palette.
type PaletteProps struct {
	// Shades
	FG       Grade
	BG       Grade
	CoShades *Grade

	// Accent
	Accent   *Grade
	CoAccent *Grade

	// Color ring
	Ring *Grade

	// 12 ring colors uniformly distributed on the LCH cylinder
	Red        *Grade
	Orange     *Grade
	Yellow     *Grade
	Chartreuse *Grade
	Green      *Grade
	Spring     *Grade
	Cyan       *Grade
	Azure      *Grade
	Blue       *Grade
	Violet     *Grade
	Magenta    *Grade
	Rose       *Grade
}

// Palette holds the derived 11-step color arrays.
type Palette struct {
	arrays map[PaletteKey][]Color
	isDark bool
}

// NewPalette builds a Palette from the props, merging shades and accent
// inheritance the same way the original generator does:
//
//	coShades  = fg  + coShades
//	accent    = fg  + accent
//	coAccent  = fg  + coAccent
//	ring color = fg + accent + ring + <color>
func NewPalette(p PaletteProps) *Palette {
	fg := p.FG

	pl := &Palette{
		arrays: map[PaletteKey][]Color{
			KeyFG:       GetColorGrades(GRADES, p.FG),
			KeyBG:       GetColorGrades(GRADES, p.BG),
			KeyCoShades: GetColorGrades(GRADES, merge(fg, p.CoShades)),
			KeyAccent:   GetColorGrades(GRADES, merge(fg, p.Accent)),
			KeyCoAccent: GetColorGrades(GRADES, merge(fg, p.CoAccent)),
		},
	}

	ringCommon := fg
	ringCommon = merge(ringCommon, p.Accent)
	ringCommon = merge(ringCommon, p.Ring)

	pl.arrays[KeyRed] = GetColorGrades(GRADES, merge(ringCommon, p.Red))
	pl.arrays[KeyOrange] = GetColorGrades(GRADES, merge(ringCommon, p.Orange))
	pl.arrays[KeyYellow] = GetColorGrades(GRADES, merge(ringCommon, p.Yellow))
	pl.arrays[KeyChartreuse] = GetColorGrades(GRADES, merge(ringCommon, p.Chartreuse))
	pl.arrays[KeyGreen] = GetColorGrades(GRADES, merge(ringCommon, p.Green))
	pl.arrays[KeySpring] = GetColorGrades(GRADES, merge(ringCommon, p.Spring))
	pl.arrays[KeyCyan] = GetColorGrades(GRADES, merge(ringCommon, p.Cyan))
	pl.arrays[KeyAzure] = GetColorGrades(GRADES, merge(ringCommon, p.Azure))
	pl.arrays[KeyBlue] = GetColorGrades(GRADES, merge(ringCommon, p.Blue))
	pl.arrays[KeyViolet] = GetColorGrades(GRADES, merge(ringCommon, p.Violet))
	pl.arrays[KeyMagenta] = GetColorGrades(GRADES, merge(ringCommon, p.Magenta))
	pl.arrays[KeyRose] = GetColorGrades(GRADES, merge(ringCommon, p.Rose))

	pl.isDark = pl.arrays[KeyFG][0].L < pl.arrays[KeyFG][GRADES].L
	return pl
}

// merge applies an override grade on top of a base grade field by field.
// A zero field in the override leaves the base value untouched.
func merge(base Grade, override *Grade) Grade {
	if override == nil {
		return base
	}
	if override.Power != 0 {
		base.Power = override.Power
	}
	if override.Luma != (Range{}) {
		base.Luma = override.Luma
	}
	if override.Chroma != (Range{}) {
		base.Chroma = override.Chroma
	}
	if override.Hue != (Range{}) {
		base.Hue = override.Hue
	}
	return base
}

// Get returns one of the derived 11-step arrays by key.
func (p *Palette) Get(k PaletteKey) []Color {
	return p.arrays[k]
}

// IsDark reports whether this palette is a dark theme (fg brightens upward).
func (p *Palette) IsDark() bool {
	return p.isDark
}

// AbsGrade returns the palette index for light-theme-inverted contexts.
// Dark themes use the grade as-is; light themes reverse it.
func (p *Palette) AbsGrade(grade int) int {
	if p.isDark {
		return grade
	}
	return GRADES - grade
}

// Print renders every grade as ANSI color blocks, matching the original
// generator's debug output.
func (p *Palette) Print() {
	for _, k := range PaletteKeys {
		line := ""
		for _, c := range p.arrays[k] {
			line += c.ANSIBlock()
		}
		println(string(k), ":", line)
	}
}

// StandardRing returns hue offsets for the 12 ring colors, 30 degrees apart,
// starting just above the given reference angle.
func StandardRing(angle float64) map[PaletteKey]float64 {
	angle = math.Mod(angle, 30) + 30
	return map[PaletteKey]float64{
		KeyRed:        angle + 0,
		KeyOrange:     angle + 30,
		KeyYellow:     angle + 60,
		KeyChartreuse: angle + 90,
		KeyGreen:      angle + 120,
		KeySpring:     angle + 150,
		KeyCyan:       angle + 180,
		KeyAzure:      angle + 210,
		KeyBlue:       angle + 240,
		KeyViolet:     angle + 270,
		KeyMagenta:    angle + 300,
		KeyRose:       angle + 330,
	}
}
