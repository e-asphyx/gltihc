package engine

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sort"
)

const (
	FilterColor = iota
	FilterGray
	FilterSource
	FilterSetRGBAComp
	FilterSetA
	FilterSetYCCComp
	FilterPermRGB
	FilterPermRGBA
	FilterPermYCC
	FilterCopyComp
	FilterCToA
	FilterMix
	FilterQuant
	FilterQuantRGBA
	FilterQuantYCCA
	FilterQuantY
	FilterInv
	FilterInvRGBAComp
	FilterInvA
	FilterInvYCCComp
	FilterGrayscale
	FilterBitRasp
	FilterNumFilters
)

type Filter interface {
	Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation)
	String() string
}

type FilterOptions struct {
	Reference *image.NRGBA64
	BlockSize int
}

type filterConstructor func(opt *FilterOptions) Filter

type filterColor color.NRGBA

func (f filterColor) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	sr := (uint32(f.R) << 8) | uint32(f.R)
	sg := (uint32(f.G) << 8) | uint32(f.G)
	sb := (uint32(f.B) << 8) | uint32(f.B)
	sa := (uint32(f.A) << 8) | uint32(f.A)
	s := NRGBA{sr, sg, sb, sa}

	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			c := op.Apply(NRGBA{dr, dg, db, da}, s)

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
	}
}

func (f filterColor) String() string {
	return fmt.Sprintf("color:[%d,%d,%d,%d]", f.R, f.G, f.B, f.A)
}

func newFilterColor(opt *FilterOptions) Filter {
	a := uint32(rand.Intn(256))
	r := uint32(rand.Intn(256))
	g := uint32(rand.Intn(256))
	b := uint32(rand.Intn(256))
	return filterColor(color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
}

func newFilterGray(opt *FilterOptions) Filter {
	a := uint32(rand.Intn(256))
	v := uint32(rand.Intn(256))
	return filterColor(color.NRGBA{uint8(v), uint8(v), uint8(v), uint8(a)})
}

type filterSetRGBAComp struct {
	c, v uint8
}

func (f filterSetRGBAComp) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	val := (uint32(f.v) << 8) | uint32(f.v)

	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			v := NRGBA{sr, sg, sb, sa}
			v[f.c] = val

			c := op.Apply(NRGBA{dr, dg, db, da}, v)

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterSetRGBAComp) String() string {
	return fmt.Sprintf("rgba:{%d:%d}", f.c, f.v)
}

func newFilterSetRGBAComp(opt *FilterOptions) Filter {
	return filterSetRGBAComp{uint8(rand.Intn(4)), uint8(rand.Intn(256))}
}

func newFilterSetA(opt *FilterOptions) Filter {
	return filterSetRGBAComp{3, uint8(rand.Intn(256))}
}

type filterSource struct{}

func (f filterSource) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{sr, sg, sb, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterSource) String() string {
	return "src"
}

func newFilterSource(opt *FilterOptions) Filter {
	return filterSource{}
}

type filterSetYCCComp struct {
	c, v uint8
}

func (f filterSetYCCComp) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			sY, sCb, sCr := color.RGBToYCbCr(uint8(sr>>8), uint8(sg>>8), uint8(sb>>8))
			v := [3]uint8{sY, sCb, sCr}
			v[f.c] = f.v

			r8, g8, b8 := color.YCbCrToRGB(v[0], v[1], v[2])
			sr = (uint32(r8) << 8) | uint32(r8)
			sg = (uint32(g8) << 8) | uint32(g8)
			sb = (uint32(b8) << 8) | uint32(b8)

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{sr, sg, sb, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterSetYCCComp) String() string {
	return fmt.Sprintf("ycc:{%d:%d}", f.c, f.v)
}

func newFilterSetYCCComp(opt *FilterOptions) Filter {
	return filterSetYCCComp{uint8(rand.Intn(3)), uint8(rand.Intn(256))}
}

type filterPermRGBA [4]int

func (f filterPermRGBA) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			v := NRGBA{sr, sg, sb, sa}
			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{v[f[0]], v[f[1]], v[f[2]], v[f[3]]})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterPermRGBA) String() string {
	return fmt.Sprintf("prgba:[%d,%d,%d,%d]", f[0], f[1], f[2], f[3])
}

func newFilterPermRGBA(opt *FilterOptions) Filter {
	p := rand.Perm(4)
	return filterPermRGBA{p[0], p[1], p[2], p[3]}
}

func newFilterPermRGB(opt *FilterOptions) Filter {
	p := rand.Perm(3)
	return filterPermRGBA{p[0], p[1], p[2], 3}
}

type filterCopyComp struct {
	d uint8
	s uint8
}

func (f filterCopyComp) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			v := NRGBA{sr, sg, sb, sa}
			v[f.d] = v[f.s]
			c := op.Apply(NRGBA{dr, dg, db, da}, v)

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterCopyComp) String() string {
	return fmt.Sprintf("cc:[%d:%d]", f.d, f.s)
}

func newFilterCopyComp(opt *FilterOptions) Filter {
	p := rand.Perm(4)
	return filterCopyComp{uint8(p[0]), uint8(p[1])}
}

func newFilterCToA(opt *FilterOptions) Filter {
	p := rand.Intn(3)
	return filterCopyComp{3, uint8(p)}
}

type filterPermYCC [3]int

func (f filterPermYCC) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			sY, sCb, sCr := color.RGBToYCbCr(uint8(sr>>8), uint8(sg>>8), uint8(sb>>8))
			v := [3]uint8{sY, sCb, sCr}

			r8, g8, b8 := color.YCbCrToRGB(v[f[0]], v[f[1]], v[f[2]])
			sr = (uint32(r8) << 8) | uint32(r8)
			sg = (uint32(g8) << 8) | uint32(g8)
			sb = (uint32(b8) << 8) | uint32(b8)

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{sr, sg, sb, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterPermYCC) String() string {
	return fmt.Sprintf("pycc:[%d,%d,%d]", f[0], f[1], f[2])
}

func newFilterPermYCC(opt *FilterOptions) Filter {
	p := rand.Perm(3)
	return filterPermYCC{p[0], p[1], p[2]}
}

type filterMix [3][3]float64

func (f filterMix) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			rr := int32(f[0][0]*float64(sr) + f[0][1]*float64(sg) + f[0][2]*float64(sb))
			gg := int32(f[1][0]*float64(sr) + f[1][1]*float64(sg) + f[1][2]*float64(sb))
			bb := int32(f[2][0]*float64(sr) + f[2][1]*float64(sg) + f[2][2]*float64(sb))

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

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{uint32(rr), uint32(gg), uint32(bb), sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
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

func newFilterMix(opt *FilterOptions) Filter {
	return filterMix{
		[3]float64{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0},
		[3]float64{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0},
		[3]float64{2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0, 2.0*rand.Float64() - 1.0},
	}
}

type filterQuantRGBA [4]uint8

func (f filterQuantRGBA) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	m := [4]uint32{1 << (f[0] + 8), 1 << (f[1] + 8), 1 << (f[2] + 8), 1 << (f[3] + 8)}

	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			sr = (sr + (m[0] >> 1)) &^ (m[0] - 1)
			sg = (sg + (m[1] >> 1)) &^ (m[1] - 1)
			sb = (sb + (m[2] >> 1)) &^ (m[2] - 1)
			sa = (sa + (m[3] >> 1)) &^ (m[3] - 1)
			if sr > 0xffff {
				sr = 0xffff
			}
			if sg > 0xffff {
				sg = 0xffff
			}
			if sb > 0xffff {
				sb = 0xffff
			}
			if sa > 0xffff {
				sa = 0xffff
			}

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{sr, sg, sb, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterQuantRGBA) String() string {
	return fmt.Sprintf("qrgba:[%d,%d,%d,%d]", f[0], f[1], f[2], f[3])
}

func newFilterQuantRGBA(opt *FilterOptions) Filter {
	return filterQuantRGBA{uint8(rand.Intn(8)), uint8(rand.Intn(8)), uint8(rand.Intn(8)), uint8(rand.Intn(8))}
}

func newFilterQuant(opt *FilterOptions) Filter {
	n := uint8(rand.Intn(8))
	return filterQuantRGBA{n, n, n, uint8(rand.Intn(8))}
}

type filterQuantYCCA [4]uint8

func (f filterQuantYCCA) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	m := [4]uint32{1 << f[0], 1 << f[1], 1 << f[2], 1 << f[3]}

	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			sY, sCb, sCr := color.RGBToYCbCr(uint8(sr>>8), uint8(sg>>8), uint8(sb>>8))
			yy := (uint32(sY) + (m[0] >> 1)) &^ (m[0] - 1)
			cb := (uint32(sCb) + (m[1] >> 1)) &^ (m[1] - 1)
			cr := (uint32(sCr) + (m[2] >> 1)) &^ (m[2] - 1)
			aa := (uint32(sa>>8) + (m[3] >> 1)) &^ (m[3] - 1)
			if yy > 0xff {
				yy = 0xff
			}
			if cb > 0xff {
				cb = 0xff
			}
			if cr > 0xff {
				cr = 0xff
			}
			if aa > 0xff {
				aa = 0xff
			}

			r8, g8, b8 := color.YCbCrToRGB(uint8(yy), uint8(cb), uint8(cr))
			sr = (uint32(r8) << 8) | uint32(r8)
			sg = (uint32(g8) << 8) | uint32(g8)
			sb = (uint32(b8) << 8) | uint32(b8)
			sa = (uint32(aa) << 8) | uint32(aa)

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{sr, sg, sb, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterQuantYCCA) String() string {
	return fmt.Sprintf("qycc:[%d,%d,%d,%d]", f[0], f[1], f[2], f[3])
}

func newFilterQuantYCCA(opt *FilterOptions) Filter {
	return filterQuantYCCA{uint8(rand.Intn(8)), uint8(rand.Intn(8)), uint8(rand.Intn(8)), uint8(rand.Intn(8))}
}

func newFilterQuantY(opt *FilterOptions) Filter {
	return filterQuantYCCA{uint8(rand.Intn(8)), 0, 0, 0}
}

type filterInv struct{}

func (f filterInv) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{0xffff - sr, 0xffff - sg, 0xffff - sb, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterInv) String() string {
	return "inv"
}

func newFilterInv(opt *FilterOptions) Filter { return filterInv{} }

type filterInvRGBAComp uint8

func (f filterInvRGBAComp) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			v := NRGBA{sr, sg, sb, sa}
			v[f] = 0xffff - v[f]

			c := op.Apply(NRGBA{dr, dg, db, da}, v)

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterInvRGBAComp) String() string {
	return fmt.Sprintf("irgba:[%d]", f)
}

func newFilterInvRGBAComp(opt *FilterOptions) Filter {
	return filterInvRGBAComp(rand.Intn(4))
}

func newFilterInvA(opt *FilterOptions) Filter {
	return filterInvRGBAComp(3)
}

type filterInvYCCComp uint8

func (f filterInvYCCComp) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			sY, sCb, sCr := color.RGBToYCbCr(uint8(sr>>8), uint8(sg>>8), uint8(sb>>8))
			v := [3]uint32{uint32(sY), uint32(sCb), uint32(sCr)}
			v[f] = 0xff - v[f]
			if f > 0 {
				// for color components zero = 128
				v[f]++
				if v[f] > 0xff {
					v[f] = 0xff
				}
			}

			r8, g8, b8 := color.YCbCrToRGB(uint8(v[0]), uint8(v[1]), uint8(v[2]))
			sr = (uint32(r8) << 8) | uint32(r8)
			sg = (uint32(g8) << 8) | uint32(g8)
			sb = (uint32(b8) << 8) | uint32(b8)

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{sr, sg, sb, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterInvYCCComp) String() string {
	return fmt.Sprintf("iycc:[%d]", f)
}

func newFilterInvYCCComp(opt *FilterOptions) Filter {
	return filterInvYCCComp(rand.Intn(3))
}

type filterGrayscale struct{}

func (f filterGrayscale) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			yy := (19595*sr + 38470*sg + 7471*sb + 1<<15) >> 16

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{yy, yy, yy, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterGrayscale) String() string {
	return "gs"
}

func newFilterGrayscale(opt *FilterOptions) Filter { return filterGrayscale{} }

type filterBitRasp struct {
	mode  uint8
	op    uint8
	ror   uint8
	mask  uint8
	alpha uint8
}

func (f filterBitRasp) Apply(dst *image.NRGBA64, dr image.Rectangle, src *image.NRGBA64, sp image.Point, op Operation) {
	m := uint32(f.mask)
	si := sp.Y*src.Stride + (sp.X << 3)
	di := dr.Min.Y*dst.Stride + (dr.Min.X << 3)
	w := dr.Dx()
	for y := dr.Min.Y; y < dr.Max.Y; y++ {
		dstSpan := dst.Pix[di : di+(w<<3) : di+(w<<3)]
		srcSpan := src.Pix[si : si+(w<<3) : si+(w<<3)]
		var i int
		for x := dr.Min.X; x < dr.Max.X; x++ {
			sr := (uint32(srcSpan[i+0]) << 8) | uint32(srcSpan[i+1])
			sg := (uint32(srcSpan[i+2]) << 8) | uint32(srcSpan[i+3])
			sb := (uint32(srcSpan[i+4]) << 8) | uint32(srcSpan[i+5])
			sa := (uint32(srcSpan[i+6]) << 8) | uint32(srcSpan[i+7])

			dr := (uint32(dstSpan[i+0]) << 8) | uint32(dstSpan[i+1])
			dg := (uint32(dstSpan[i+2]) << 8) | uint32(dstSpan[i+3])
			db := (uint32(dstSpan[i+4]) << 8) | uint32(dstSpan[i+5])
			da := (uint32(dstSpan[i+6]) << 8) | uint32(dstSpan[i+7])

			var mix uint32
			switch f.mode {
			case 0:
				mix = uint32(x)
			case 1:
				mix = uint32(y)
			case 2:
				mix = uint32(y + x)
			case 3:
				mix = uint32(y - x)
			case 4:
				mix = uint32(y) | uint32(x)
			case 5:
				mix = uint32(y) & uint32(x)
			default:
				mix = uint32(y) ^ uint32(x)
			}

			mix &= 0xff
			mix = (mix >> f.ror) | ((mix << (8 - f.ror)) & 0xff)
			mix = (mix & m) << 8

			aa := sa
			switch f.op {
			case 0:
				sr &= mix
				sg &= mix
				sb &= mix
				aa &= mix
			case 1:
				sr ^= mix
				sg ^= mix
				sb ^= mix
				aa ^= mix
			case 2:
				sr |= mix
				sg |= mix
				sb |= mix
				aa |= mix
			default:
				sr = sr&^m | mix
				sg = sg&^m | mix
				sb = sb&^m | mix
				aa = aa&^m | mix
			}
			if f.alpha == 1 {
				sa = aa
			}

			c := op.Apply(NRGBA{dr, dg, db, da}, NRGBA{sr, sg, sb, sa})

			dstSpan[i+0] = uint8(c[0] >> 8)
			dstSpan[i+1] = uint8(c[0])
			dstSpan[i+2] = uint8(c[1] >> 8)
			dstSpan[i+3] = uint8(c[1])
			dstSpan[i+4] = uint8(c[2] >> 8)
			dstSpan[i+5] = uint8(c[2])
			dstSpan[i+6] = uint8(c[3] >> 8)
			dstSpan[i+7] = uint8(c[3])
			i += 8
		}
		di += dst.Stride
		si += src.Stride
	}
}

func (f filterBitRasp) String() string {
	return fmt.Sprintf("rasp:{m:%d,op:%d,mask:%d,a:%d,r:%d}", f.mode, f.op, f.mask, f.alpha, f.ror)
}

func newFilterBitRasp(opt *FilterOptions) Filter {
	ret := filterBitRasp{
		mode:  uint8(rand.Intn(7)),
		op:    uint8(rand.Intn(4)),
		alpha: uint8(rand.Intn(2)),
		ror:   uint8(rand.Intn(8)),
	}

	if rand.Intn(2) == 1 {
		bits := 2 + rand.Intn(7)
		ret.mask = (1 << bits) - 1
	} else {
		ret.mask = uint8(rand.Intn(256))
	}
	return ret
}

var filtersTable = []filterConstructor{
	FilterColor:       newFilterColor,
	FilterGray:        newFilterGray,
	FilterSource:      newFilterSource,
	FilterSetRGBAComp: newFilterSetRGBAComp,
	FilterSetA:        newFilterSetA,
	FilterSetYCCComp:  newFilterSetYCCComp,
	FilterPermRGB:     newFilterPermRGB,
	FilterPermRGBA:    newFilterPermRGBA,
	FilterPermYCC:     newFilterPermYCC,
	FilterCopyComp:    newFilterCopyComp,
	FilterCToA:        newFilterCToA,
	FilterMix:         newFilterMix,
	FilterQuant:       newFilterQuant,
	FilterQuantRGBA:   newFilterQuantRGBA,
	FilterQuantYCCA:   newFilterQuantYCCA,
	FilterQuantY:      newFilterQuantY,
	FilterInv:         newFilterInv,
	FilterInvRGBAComp: newFilterInvRGBAComp,
	FilterInvA:        newFilterInvA,
	FilterInvYCCComp:  newFilterInvYCCComp,
	FilterGrayscale:   newFilterGrayscale,
	FilterBitRasp:     newFilterBitRasp,
}

func NewRandomizedFilter(f int, opt *FilterOptions) Filter {
	if f < len(filtersTable) {
		return filtersTable[f](opt)
	}
	return nil
}

var filterNames = map[string]int{
	"color":   FilterColor,
	"gray":    FilterGray,
	"src":     FilterSource,
	"rgba":    FilterSetRGBAComp,
	"seta":    FilterSetA,
	"ycc":     FilterSetYCCComp,
	"prgb":    FilterPermRGB,
	"prgba":   FilterPermRGBA,
	"pycc":    FilterPermYCC,
	"copy":    FilterCopyComp,
	"ctoa":    FilterCToA,
	"mix":     FilterMix,
	"quant":   FilterQuant,
	"qrgba":   FilterQuantRGBA,
	"qycca":   FilterQuantYCCA,
	"qy":      FilterQuantY,
	"inv":     FilterInv,
	"invrgba": FilterInvRGBAComp,
	"inva":    FilterInvA,
	"invycc":  FilterInvYCCComp,
	"gs":      FilterGrayscale,
	"rasp":    FilterBitRasp,
}

func GetFilterID(name string) int {
	if id, ok := filterNames[name]; ok {
		return id
	}
	return -1
}

func FilterNames() []string {
	ret := make([]string, 0, len(filterNames))
	for name := range filterNames {
		ret = append(ret, name)
	}
	sort.Strings(ret)
	return ret
}
