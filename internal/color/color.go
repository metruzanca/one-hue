package color

import (
	"fmt"
	"math"
	"strings"
)

// Color is an sRGB color carrying its CIELAB coordinates and opacity.
// Red/Green/Blue are 0..255 floats, L 0..100, opacity 0..255.
type Color struct {
	R  float64
	G  float64
	B  float64
	L  float64
	A  float64
	Bv float64
	// Opacity is 0..255, matching the trailing alpha byte of #RRGGBBAA hex.
	Opacity float64
}

// C returns the chroma (distance from the L axis).
func (c Color) C() float64 {
	return math.Hypot(c.A, c.Bv)
}

// Hue returns the LCH hue angle in degrees, 0..360.
func (c Color) Hue() float64 {
	h := math.Atan2(c.Bv, c.A) / DEGREE
	if h < 0 {
		h += 360
	}
	return h
}

// DEGREE converts degrees to radians.
const DEGREE = math.Pi / 180

// Alpha returns a copy with multiplied opacity.
func (c Color) Alpha(o float64) Color {
	c.Opacity *= o
	return c
}

// RGB builds a Color from 8-bit channels.
func RGB(r, g, b float64) Color {
	l, a, bv := RGB2Lab(r, g, b)
	return Color{R: r, G: g, B: b, L: l, A: a, Bv: bv, Opacity: 0xff}
}

// Lab builds a Color from CIELAB channels.
func Lab(l, a, b float64) Color {
	r, g, bb := Lab2RGB(l, a, b)
	return Color{R: r, G: g, B: bb, L: l, A: a, Bv: b, Opacity: 0xff}
}

// LCH builds a Color from L, chroma and hue angle in degrees.
func LCH(l, c, hDeg float64) Color {
	h := hDeg * DEGREE
	a := c * math.Cos(h)
	b := c * math.Sin(h)
	r, g, bb := Lab2RGB(l, a, b)
	return Color{R: r, G: g, B: bb, L: l, A: a, Bv: b, Opacity: 0xff}
}

// Hex parses a #RRGGBB string.
func Hex(s string) (Color, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return Color{}, fmt.Errorf("expected #RRGGBB, got %q", s)
	}
	var r, g, b int
	if _, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b); err != nil {
		return Color{}, err
	}
	return RGB(float64(r), float64(g), float64(b)), nil
}

func lerp(a, s, b float64) float64 {
	return a + s*(b-a)
}

// Mix mixes this color into other in LAB space by scale 0..1.
func (c Color) Mix(scale float64, other Color) Color {
	l1 := lerp(c.L, scale, other.L)
	a1 := lerp(c.A, scale, other.A)
	b1 := lerp(c.Bv, scale, other.Bv)
	op := c.Opacity + (other.Opacity-c.Opacity)*scale
	r, g, bb := Lab2RGB(l1, a1, b1)
	return Color{R: r, G: g, B: bb, L: l1, A: a1, Bv: b1, Opacity: op}
}

func unGamma(x, gamma float64) float64 {
	return math.Pow(x/0xff, gamma)
}

func gamma(x, gamma float64) float64 {
	return math.Pow(x, 1/gamma) * 0xff
}

// MixRGB mixes this color into other in gamma-corrected RGB space.
func (c Color) MixRGB(scale float64, other Color, g ...float64) Color {
	gam := 2.2
	if len(g) > 0 {
		gam = g[0]
	}
	r1 := lerp(unGamma(c.R, gam), scale, unGamma(other.R, gam))
	g1 := lerp(unGamma(c.G, gam), scale, unGamma(other.G, gam))
	b1 := lerp(unGamma(c.B, gam), scale, unGamma(other.B, gam))
	op := c.Opacity + (other.Opacity-c.Opacity)*scale
	rr, gg, bb := gamma(r1, gam), gamma(g1, gam), gamma(b1, gam)
	l, a, b := RGB2Lab(rr, gg, bb)
	return Color{R: rr, G: gg, B: bb, L: l, A: a, Bv: b, Opacity: op}
}

// chan clamps to a 2-digit uppercase hex byte.
func chanByte(n float64) string {
	x := int(math.Min(0xff, math.Max(0, math.Round(n))))
	return fmt.Sprintf("%02X", x)
}

// Hex returns the #RRGGBB string.
func (c Color) Hex() string {
	return "#" + chanByte(c.R) + chanByte(c.G) + chanByte(c.B)
}

// HexAA returns the #RRGGBBAA string including opacity.
func (c Color) HexAA() string {
	return c.Hex() + chanByte(c.Opacity)
}

// ANSIBlock returns a false-color terminal block labeled with the raw hex,
// with foreground chosen by relative lightness.
func (c Color) ANSIBlock() string {
	fg := "30"
	if c.L < 50 {
		fg = "97"
	}
	return fmt.Sprintf(
		"\x1b[%s;48;2;%d;%d;%dm %s \x1b[0m",
		fg, int(c.R), int(c.G), int(c.B), c.HexRaw(),
	)
}

// HexRaw returns the RRGGBB string without the leading hash.
func (c Color) HexRaw() string {
	return chanByte(c.R) + chanByte(c.G) + chanByte(c.B)
}

// Grades interpolates between a and b across steps+1 colors using a power curve.
func Grades(steps int, a, b Color, power ...float64) []Color {
	p := 1.0
	if len(power) > 0 {
		p = power[0]
	}
	out := make([]Color, 0, steps+1)
	for s := 0; s <= steps; s++ {
		out = append(out, a.Mix(math.Pow(float64(s)/float64(steps), p), b))
	}
	return out
}
