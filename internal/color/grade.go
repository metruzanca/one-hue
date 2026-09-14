package color

import "math"

// Grade is a luma/chroma/hue gradient specification. A range value may be a
// scalar (constant) or a [start, end] pair. A zero Range is constant 0,
// matching the original where null/undefined grades evaluate to 0.
type Grade struct {
	// Power is the exponent applied to the progression. 0 means linear (1).
	Power  float64
	Luma   Range
	Chroma Range
	Hue    Range
}

// Range is a scalar (Start == End) or an interpolation [Start, End] span.
// The zero value is constant 0.
type Range struct {
	Start float64
	End   float64
}

// Scalar returns a constant Range.
func Scalar(v float64) Range {
	return Range{Start: v, End: v}
}

// Span returns an interpolated Range.
func Span(a, b float64) Range {
	return Range{Start: a, End: b}
}

func mix(a, s, b float64) float64 {
	return a + s*(b-a)
}

// GetColorGrades builds a steps+1 element array of colors across the grade,
// where index 0 is the dark end for dark themes (step 0 of the progression).
func GetColorGrades(steps int, grade Grade) []Color {
	power := grade.Power
	if power == 0 {
		power = 1
	}
	colors := make([]Color, 0, steps+1)
	for step := 0; step <= steps; step++ {
		scaleChroma := float64(step) / float64(steps)
		scaleLuma := math.Pow(scaleChroma, power)

		luma := mix(grade.Luma.Start, scaleLuma, grade.Luma.End)
		chroma := mix(grade.Chroma.Start, scaleChroma, grade.Chroma.End)
		hue := mix(grade.Hue.Start, scaleChroma, grade.Hue.End)
		colors = append(colors, LCH(luma, chroma, hue))
	}
	return colors
}
