package color

import "math"

const delta = 6.0 / 29.0

// D65 normalization point.
const xn = 0.95047
const yn = 1.0
const zn = 1.08883

// Lab2RGB converts CIELAB channels (L 0..100, a/b signed) to 8-bit RGB.
func Lab2RGB(l, a, b float64) (r, g, bb float64) {
	y := (l + 16) / 116
	x := a/500 + y
	z := y - b/200

	x = xn * labFN(x)
	y = yn * labFN(y)
	z = zn * labFN(z)

	r = x*3.2406 + y*-1.5372 + z*-0.4986
	g = x*-0.9689 + y*1.8758 + z*0.0415
	bb = x*0.0557 + y*-0.204 + z*1.057

	return rgbGamma(r) * 255, rgbGamma(g) * 255, rgbGamma(bb) * 255
}

func rgbGamma(n float64) float64 {
	var r1 float64
	if n > 0.0031308 {
		r1 = 1.055*math.Pow(n, 1/2.4) - 0.055
	} else {
		r1 = 12.92 * n
	}
	return math.Max(0, math.Min(1, r1))
}

func labFN(x float64) float64 {
	if x > delta {
		return x * x * x
	}
	return 3 * delta * delta * (x - 4.0/29.0)
}

// RGB2Lab converts 8-bit RGB channels to CIELAB (D65).
func RGB2Lab(r, g, b float64) (l, a, bb float64) {
	r = rgbDeGamma(r / 255)
	g = rgbDeGamma(g / 255)
	b = rgbDeGamma(b / 255)

	x := labF((r*0.4124 + g*0.3576 + b*0.1805) / xn)
	y := labF((r*0.2126 + g*0.7152 + b*0.0722) / yn)
	z := labF((r*0.0193 + g*0.1192 + b*0.9505) / zn)

	return 116*y - 16, 500 * (x - y), 200 * (y - z)
}

func rgbDeGamma(n float64) float64 {
	if n > 0.04045 {
		return math.Pow((n+0.055)/1.055, 2.4)
	}
	return n / 12.92
}

func labF(x float64) float64 {
	if x > delta*delta*delta {
		return math.Pow(x, 1.0/3.0)
	}
	return x/(3*delta*delta) + 4.0/29.0
}
