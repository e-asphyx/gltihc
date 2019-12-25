package engine

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

type Randn interface {
	ExpFloat64() float64
	Float32() float32
	Float64() float64
	Int() int
	Int31() int32
	Int31n(n int32) int32
	Int63() int64
	Int63n(n int64) int64
	Intn(n int) int
	NormFloat64() float64
	Perm(n int) []int
	Read(p []byte) (n int, err error)
	Seed(seed int64)
	Shuffle(n int, swap func(i, j int))
	Uint32() uint32
	Uint64() uint64
}

const (
	FilterColor = iota
	FilterSource
	FilterSetRGBAComp
	FilterSetYCCComp
	FilterPermRGB
	FilterPermRGBA
	FilterPermYCC
	FilterMix
	FilterQuantRGBA
	FilterQuantYCCA
	FilterInv
	FilterInvRGBAComp
	FilterInvYCCComp
	FilterGrayscale
	FilterBitRasp
	FilterNumFilters
)

type Filter interface {
	Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation)
	String() string
}

type FilterOptions struct {
	Reference image.Image
	BlockSize int
}

type filterConstructor func(opt *FilterOptions, r Randn) Filter

type filterColor color.RGBA

func (f filterColor) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), color.RGBA(f)))
		}
	}
}

func (f filterColor) String() string {
	return fmt.Sprintf("color:[%d,%d,%d,%d]", f.R, f.G, f.B, f.A)
}

func newFilterColor(opt *FilterOptions, rand Randn) Filter {
	a := uint32(rand.Intn(256))
	r := (uint32(rand.Intn(256)) * a) / 0xff
	g := (uint32(rand.Intn(256)) * a) / 0xff
	b := (uint32(rand.Intn(256)) * a) / 0xff
	return filterColor(color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
}

type filterSetRGBAComp struct {
	c, v uint8
}

func (f filterSetRGBAComp) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	val := (uint32(f.v) << 8) | uint32(f.v)
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			if a != 0 {
				r = (r * 0xffff) / a
				g = (g * 0xffff) / a
				b = (b * 0xffff) / a
			}
			v := [4]uint32{r, g, b, a}
			v[f.c] = val
			a = v[3]
			r = (v[0] * a) / 0xffff
			g = (v[1] * a) / 0xffff
			b = (v[2] * a) / 0xffff
			c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterSetRGBAComp) String() string {
	return fmt.Sprintf("rgba:{%d:%d}", f.c, f.v)
}

func newFilterSetRGBAComp(opt *FilterOptions, rand Randn) Filter {
	return filterSetRGBAComp{uint8(rand.Intn(4)), uint8(rand.Intn(256))}
}

type filterSource struct{}

func (f filterSource) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterSource) String() string {
	return "src"
}

func newFilterSource(opt *FilterOptions, rand Randn) Filter {
	return filterSource{}
}

type filterSetYCCComp struct {
	c, v uint8
}

func (f filterSetYCCComp) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			s := src.At(sp.X+x, sp.Y+y)
			ycc := color.NYCbCrAModel.Convert(s).(color.NYCbCrA)
			v := [3]uint8{ycc.Y, ycc.Cb, ycc.Cr}
			v[f.c] = f.v
			c := color.NYCbCrA{color.YCbCr{v[0], v[1], v[2]}, ycc.A}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterSetYCCComp) String() string {
	return fmt.Sprintf("ycc:{%d:%d}", f.c, f.v)
}

func newFilterSetYCCComp(opt *FilterOptions, rand Randn) Filter {
	return filterSetYCCComp{uint8(rand.Intn(3)), uint8(rand.Intn(256))}
}

type filterPermRGBA [4]int

func (f filterPermRGBA) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			if a != 0 {
				r = (r * 0xffff) / a
				g = (g * 0xffff) / a
				b = (b * 0xffff) / a
			}
			v := [4]uint32{r, g, b, a}
			a = v[f[3]]
			r = (v[f[0]] * a) / 0xffff
			g = (v[f[1]] * a) / 0xffff
			b = (v[f[2]] * a) / 0xffff
			c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterPermRGBA) String() string {
	return fmt.Sprintf("prgba:[%d,%d,%d,%d]", f[0], f[1], f[2], f[3])
}

func newFilterPermRGBA(opt *FilterOptions, rand Randn) Filter {
	p := rand.Perm(4)
	return filterPermRGBA{p[0], p[1], p[2], p[3]}
}

func newFilterPermRGB(opt *FilterOptions, rand Randn) Filter {
	p := rand.Perm(3)
	return filterPermRGBA{p[0], p[1], p[2], 3}
}

type filterPermYCC [3]int

func (f filterPermYCC) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			s := src.At(sp.X+x, sp.Y+y)
			ycc := color.NYCbCrAModel.Convert(s).(color.NYCbCrA)
			v := [3]uint8{ycc.Y, ycc.Cb, ycc.Cr}
			yy := v[f[0]]
			cb := v[f[1]]
			cr := v[f[2]]
			c := color.NYCbCrA{color.YCbCr{yy, cb, cr}, ycc.A}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterPermYCC) String() string {
	return fmt.Sprintf("pycc:[%d,%d,%d]", f[0], f[1], f[2])
}

func newFilterPermYCC(opt *FilterOptions, rand Randn) Filter {
	p := rand.Perm(3)
	return filterPermYCC{p[0], p[1], p[2]}
}

type filterMix [3][3]float64

func (f filterMix) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			rr := int32(f[0][0]*float64(r) + f[0][1]*float64(g) + f[0][2]*float64(b))
			gg := int32(f[1][0]*float64(r) + f[1][1]*float64(g) + f[1][2]*float64(b))
			bb := int32(f[2][0]*float64(r) + f[2][1]*float64(g) + f[2][2]*float64(b))
			if rr < 0 {
				rr = 0
			} else if rr > 0xffff {
				rr = 0xffff
			}
			if gg < 0 {
				gg = 0
			} else if gg > 0xffff {
				gg = 0xffff
			}
			if bb < 0 {
				bb = 0
			} else if bb > 0xffff {
				bb = 0xffff
			}
			c := color.RGBA64{uint16(rr), uint16(gg), uint16(bb), uint16(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterMix) String() string {
	return fmt.Sprintf(
		"mix:[%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f]",
		f[0][0], f[0][1], f[0][2],
		f[1][0], f[1][1], f[1][2],
		f[2][0], f[2][1], f[2][2],
	)
}

func newFilterMix(opt *FilterOptions, rand Randn) Filter {
	return filterMix{
		[3]float64{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0},
		[3]float64{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0},
		[3]float64{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0},
	}
}

type filterQuantRGBA [4]uint8

func (f filterQuantRGBA) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	m := [4]uint32{1 << (f[0] + 8), 1 << (f[1] + 8), 1 << (f[2] + 8), 1 << (f[3] + 8)}
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			if a != 0 {
				r = (r * 0xffff) / a
				g = (g * 0xffff) / a
				b = (b * 0xffff) / a
			}
			r = (r + (m[0] >> 1)) &^ (m[0] - 1)
			g = (g + (m[1] >> 1)) &^ (m[1] - 1)
			b = (b + (m[2] >> 1)) &^ (m[2] - 1)
			a = (a + (m[3] >> 1)) &^ (m[3] - 1)
			if r > 0xffff {
				r = 0xffff
			}
			if g > 0xffff {
				g = 0xffff
			}
			if b > 0xffff {
				b = 0xffff
			}
			if a > 0xffff {
				a = 0xffff
			}
			r = (r * a) / 0xffff
			g = (g * a) / 0xffff
			b = (b * a) / 0xffff
			c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterQuantRGBA) String() string {
	return fmt.Sprintf("qrgba:[%d,%d,%d,%d]", f[0], f[1], f[2], f[3])
}

func newFilterQuantRGBA(opt *FilterOptions, rand Randn) Filter {
	return filterQuantRGBA{uint8(rand.Intn(8)), uint8(rand.Intn(8)), uint8(rand.Intn(8)), uint8(rand.Intn(8))}
}

type filterQuantYCCA [4]uint8

func (f filterQuantYCCA) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	m := [4]uint32{1 << f[0], 1 << f[1], 1 << f[2], 1 << f[3]}
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			s := src.At(sp.X+x, sp.Y+y)
			ycc := color.NYCbCrAModel.Convert(s).(color.NYCbCrA)
			yy := (uint32(ycc.Y) + (m[0] >> 1)) &^ (m[0] - 1)
			cb := (uint32(ycc.Cb) + (m[1] >> 1)) &^ (m[1] - 1)
			cr := (uint32(ycc.Cr) + (m[2] >> 1)) &^ (m[2] - 1)
			a := (uint32(ycc.A) + (m[3] >> 1)) &^ (m[3] - 1)
			if yy > 0xff {
				yy = 0xff
			}
			if cb > 0xff {
				cb = 0xff
			}
			if cr > 0xff {
				cr = 0xff
			}
			if a > 0xff {
				a = 0xff
			}
			c := color.NYCbCrA{color.YCbCr{Y: uint8(yy), Cb: uint8(cb), Cr: uint8(cr)}, uint8(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterQuantYCCA) String() string {
	return fmt.Sprintf("qycc:[%d,%d,%d,%d]", f[0], f[1], f[2], f[3])
}

func newFilterQuantYCCA(opt *FilterOptions, rand Randn) Filter {
	return filterQuantYCCA{uint8(rand.Intn(8)), uint8(rand.Intn(8)), uint8(rand.Intn(8)), uint8(rand.Intn(8))}
}

type filterInv struct{}

func (f filterInv) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			r = a - r
			g = a - g
			b = a - b
			c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterInv) String() string {
	return "inv"
}

func newFilterInv(opt *FilterOptions, rand Randn) Filter { return filterInv{} }

type filterInvRGBAComp uint8

func (f filterInvRGBAComp) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			if a != 0 {
				r = (r * 0xffff) / a
				g = (g * 0xffff) / a
				b = (b * 0xffff) / a
			}
			v := [4]uint32{r, g, b, a}
			v[f] = 0xffff - v[f]
			a = v[3]
			r = (v[0] * a) / 0xffff
			g = (v[1] * a) / 0xffff
			b = (v[2] * a) / 0xffff
			c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterInvRGBAComp) String() string {
	return fmt.Sprintf("irgba:[%d]", f)
}

func newFilterInvRGBAComp(opt *FilterOptions, rand Randn) Filter {
	return filterInvRGBAComp(rand.Intn(4))
}

type filterInvYCCComp uint8

func (f filterInvYCCComp) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			s := src.At(sp.X+x, sp.Y+y)
			ycc := color.NYCbCrAModel.Convert(s).(color.NYCbCrA)
			v := [3]uint32{uint32(ycc.Y), uint32(ycc.Cb), uint32(ycc.Cr)}
			v[f] = 0xff - v[f]
			if f > 0 {
				// for color components zero = 128
				v[f]++
				if v[f] > 0xff {
					v[f] = 0xff
				}
			}
			c := color.NYCbCrA{color.YCbCr{uint8(v[0]), uint8(v[1]), uint8(v[2])}, ycc.A}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterInvYCCComp) String() string {
	return fmt.Sprintf("iycc:[%d]", f)
}

func newFilterInvYCCComp(opt *FilterOptions, rand Randn) Filter {
	return filterInvYCCComp(rand.Intn(3))
}

type filterGrayscale struct{}

func (f filterGrayscale) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			if a != 0 {
				r = (r * 0xffff) / a
				g = (g * 0xffff) / a
				b = (b * 0xffff) / a
			}
			yy := (19595*r + 38470*g + 7471*b + 1<<15) >> 16
			yy = (yy * a) / 0xffff
			c := color.RGBA64{uint16(yy), uint16(yy), uint16(yy), uint16(a)}
			dst.Set(dr.Min.X+x, dr.Min.Y+y, op.Apply(dst.At(dr.Min.X+x, dr.Min.Y+y), c))
		}
	}
}

func (f filterGrayscale) String() string {
	return "gs"
}

func newFilterGrayscale(opt *FilterOptions, rand Randn) Filter { return filterGrayscale{} }

type filterBitRasp struct {
	mode  uint8
	op    uint8
	bits  uint8
	mask  uint8
	alpha uint8
}

func (f filterBitRasp) Apply(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point, op Operation) {
	var m uint32
	if f.bits != 0 {
		m = (uint32(1) << f.bits) - 1
	} else {
		m = uint32(f.mask)
	}
	for y := 0; y < dr.Dy(); y++ {
		for x := 0; x < dr.Dx(); x++ {
			r, g, b, a := src.At(sp.X+x, sp.Y+y).RGBA()
			if a != 0 {
				r = (r * 0xffff) / a
				g = (g * 0xffff) / a
				b = (b * 0xffff) / a
			}
			dx, dy := dr.Min.X+x, dr.Min.Y+y
			var mix uint32
			switch f.mode {
			case 0:
				mix = uint32(dx)
			case 1:
				mix = uint32(dy)
			case 2:
				mix = uint32(dy + dx)
			case 3:
				mix = uint32(dy) | uint32(dx)
			case 4:
				mix = uint32(dy) & uint32(dx)
			default:
				mix = uint32(dy) ^ uint32(dx)
			}
			mix = (mix & m) << 8
			aa := a
			switch f.op {
			case 0:
				r &= mix
				g &= mix
				b &= mix
				aa &= mix
			case 1:
				r ^= mix
				g ^= mix
				b ^= mix
				aa ^= mix
			case 2:
				r |= mix
				g |= mix
				b |= mix
				aa |= mix
			default:
				r = r&^m | mix
				g = g&^m | mix
				b = b&^m | mix
				aa = aa&^m | mix
			}
			if f.alpha == 1 {
				a = aa
			}
			r = (r * a) / 0xffff
			g = (g * a) / 0xffff
			b = (b * a) / 0xffff
			c := color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
			dst.Set(dx, dy, op.Apply(dst.At(dx, dy), c))
		}
	}
}

func (f filterBitRasp) String() string {
	return fmt.Sprintf("rasp:{m:%d,op:%d,b:%d,mask:%d,a:%d}", f.mode, f.op, f.bits, f.mask, f.alpha)
}

func newFilterBitRasp(opt *FilterOptions, rand Randn) Filter {
	ret := filterBitRasp{
		mode:  uint8(rand.Intn(6)),
		op:    uint8(rand.Intn(4)),
		alpha: uint8(rand.Intn(2)),
	}
	if rand.Intn(2) == 1 {
		ret.bits = uint8(2 + rand.Intn(7))
	} else {
		ret.mask = uint8(rand.Intn(256))
	}
	return ret
}

var filtersTable = []filterConstructor{
	FilterColor:       newFilterColor,
	FilterSource:      newFilterSource,
	FilterSetRGBAComp: newFilterSetRGBAComp,
	FilterSetYCCComp:  newFilterSetYCCComp,
	FilterPermRGB:     newFilterPermRGB,
	FilterPermRGBA:    newFilterPermRGBA,
	FilterPermYCC:     newFilterPermYCC,
	FilterMix:         newFilterMix,
	FilterQuantRGBA:   newFilterQuantRGBA,
	FilterQuantYCCA:   newFilterQuantYCCA,
	FilterInv:         newFilterInv,
	FilterInvRGBAComp: newFilterInvRGBAComp,
	FilterInvYCCComp:  newFilterInvYCCComp,
	FilterGrayscale:   newFilterGrayscale,
	FilterBitRasp:     newFilterBitRasp,
}

func NewRandomizedFilter(f int, opt *FilterOptions, r Randn) Filter {
	if f < len(filtersTable) {
		return filtersTable[f](opt, r)
	}
	return nil
}
