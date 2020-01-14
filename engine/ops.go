package engine

import (
	"image/color"
	"sort"
)

const (
	OpCompose = iota
	OpReplace
	OpAdd
	OpAddRGBMod
	OpAddYCCMod
	OpMulRGB
	OpMulYCC
	OpXorRGB
	OpXorYCC
	OpNumOps
)

type NRGBA [4]uint32

type Operation interface {
	Apply(dst, src NRGBA) NRGBA
	String() string
}

type opCompose struct{}

func (o opCompose) Apply(dst, src NRGBA) NRGBA {
	if src[3] == 0xffff {
		return src
	}
	sr := (src[0] * src[3]) / 0xffff
	sg := (src[1] * src[3]) / 0xffff
	sb := (src[2] * src[3]) / 0xffff

	dr := (dst[0] * dst[3]) / 0xffff
	dg := (dst[1] * dst[3]) / 0xffff
	db := (dst[2] * dst[3]) / 0xffff

	r := sr + dr*(0xffff-src[3])/0xffff
	g := sg + dg*(0xffff-src[3])/0xffff
	b := sb + db*(0xffff-src[3])/0xffff
	a := src[3] + dst[3]*(0xffff-src[3])/0xffff

	if a != 0 {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}

	return NRGBA{r, g, b, a}
}

func (o opCompose) String() string { return "comp" }

type opReplace struct{}

func (o opReplace) Apply(dst, src NRGBA) NRGBA { return src }

func (o opReplace) String() string { return "rep" }

type opAdd struct{}

func (o opAdd) Apply(dst, src NRGBA) NRGBA {
	sr := src[0] + dst[0]
	sg := src[1] + dst[1]
	sb := src[2] + dst[2]
	if sr > 0xffff {
		sr = 0xffff
	}
	if sg > 0xffff {
		sg = 0xffff
	}
	if sb > 0xffff {
		sb = 0xffff
	}

	if src[3] == 0xffff {
		return NRGBA{sr, sg, sb, src[3]}
	}

	// Compose
	sr = (sr * src[3]) / 0xffff
	sg = (sg * src[3]) / 0xffff
	sb = (sb * src[3]) / 0xffff

	dr := (dst[0] * dst[3]) / 0xffff
	dg := (dst[1] * dst[3]) / 0xffff
	db := (dst[2] * dst[3]) / 0xffff

	r := sr + dr*(0xffff-src[3])/0xffff
	g := sg + dg*(0xffff-src[3])/0xffff
	b := sb + db*(0xffff-src[3])/0xffff
	a := src[3] + dst[3]*(0xffff-src[3])/0xffff

	if a != 0 {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}

	return NRGBA{r, g, b, a}
}

func (o opAdd) String() string { return "add" }

type opAddRGBMod struct{}

func (o opAddRGBMod) Apply(dst, src NRGBA) NRGBA {
	sr := src[0] + dst[0]
	sg := src[1] + dst[1]
	sb := src[2] + dst[2]

	sr &= 0xffff
	sg &= 0xffff
	sb &= 0xffff

	if src[3] == 0xffff {
		return NRGBA{sr, sg, sb, src[3]}
	}

	// Compose
	sr = (sr * src[3]) / 0xffff
	sg = (sg * src[3]) / 0xffff
	sb = (sb * src[3]) / 0xffff

	dr := (dst[0] * dst[3]) / 0xffff
	dg := (dst[1] * dst[3]) / 0xffff
	db := (dst[2] * dst[3]) / 0xffff

	r := sr + dr*(0xffff-src[3])/0xffff
	g := sg + dg*(0xffff-src[3])/0xffff
	b := sb + db*(0xffff-src[3])/0xffff
	a := src[3] + dst[3]*(0xffff-src[3])/0xffff

	if a != 0 {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}

	return NRGBA{r, g, b, a}
}

func (o opAddRGBMod) String() string { return "addrgbm" }

type opAddYCCMod struct{}

func (o opAddYCCMod) Apply(dst, src NRGBA) NRGBA {
	sY, sCb, sCr := color.RGBToYCbCr(uint8(src[0]>>8), uint8(src[1]>>8), uint8(src[2]>>8))
	dY, dCb, dCr := color.RGBToYCbCr(uint8(dst[0]>>8), uint8(dst[1]>>8), uint8(dst[2]>>8))

	sY = (sY + dY) & 0xff
	sCb = uint8((int32(sCb) + int32(dCb) - 128) & 0xff)
	sCr = uint8((int32(sCr) + int32(dCr) - 128) & 0xff)

	r8, g8, b8 := color.YCbCrToRGB(sY, sCb, sCr)
	sr := (uint32(r8) << 8) | uint32(r8)
	sg := (uint32(g8) << 8) | uint32(g8)
	sb := (uint32(b8) << 8) | uint32(b8)

	if src[3] == 0xffff {
		return NRGBA{sr, sg, sb, src[3]}
	}

	// Compose
	sr = (sr * src[3]) / 0xffff
	sg = (sg * src[3]) / 0xffff
	sb = (sb * src[3]) / 0xffff

	dr := (dst[0] * dst[3]) / 0xffff
	dg := (dst[1] * dst[3]) / 0xffff
	db := (dst[2] * dst[3]) / 0xffff

	r := sr + dr*(0xffff-src[3])/0xffff
	g := sg + dg*(0xffff-src[3])/0xffff
	b := sb + db*(0xffff-src[3])/0xffff
	a := src[3] + dst[3]*(0xffff-src[3])/0xffff

	if a != 0 {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}

	return NRGBA{r, g, b, a}
}

func (o opAddYCCMod) String() string { return "addyccm" }

type opMulRGB struct{}

func (o opMulRGB) Apply(dst, src NRGBA) NRGBA {
	sr := (src[0] * dst[0]) / 0xffff
	sg := (src[1] * dst[1]) / 0xffff
	sb := (src[2] * dst[2]) / 0xffff
	if sr > 0xffff {
		sr = 0xffff
	}
	if sg > 0xffff {
		sg = 0xffff
	}
	if sb > 0xffff {
		sb = 0xffff
	}

	if src[3] == 0xffff {
		return NRGBA{sr, sg, sb, src[3]}
	}

	// Compose
	sr = (sr * src[3]) / 0xffff
	sg = (sg * src[3]) / 0xffff
	sb = (sb * src[3]) / 0xffff

	dr := (dst[0] * dst[3]) / 0xffff
	dg := (dst[1] * dst[3]) / 0xffff
	db := (dst[2] * dst[3]) / 0xffff

	r := sr + dr*(0xffff-src[3])/0xffff
	g := sg + dg*(0xffff-src[3])/0xffff
	b := sb + db*(0xffff-src[3])/0xffff
	a := src[3] + dst[3]*(0xffff-src[3])/0xffff

	if a != 0 {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}

	return NRGBA{r, g, b, a}
}

func (o opMulRGB) String() string { return "mulrgb" }

type opMulYCC struct{}

func (o opMulYCC) Apply(dst, src NRGBA) NRGBA {
	sY, sCb, sCr := color.RGBToYCbCr(uint8(src[0]>>8), uint8(src[1]>>8), uint8(src[2]>>8))
	dY, dCb, dCr := color.RGBToYCbCr(uint8(dst[0]>>8), uint8(dst[1]>>8), uint8(dst[2]>>8))

	sY = uint8(uint32(sY) * uint32(dY) / 0xff)
	sCb = uint8((int32(sCb)-128)*(int32(dCb)-128)/0xff + 128)
	sCr = uint8((int32(sCr)-128)*(int32(dCr)-128)/0xff + 128)

	r8, g8, b8 := color.YCbCrToRGB(sY, sCb, sCr)
	sr := (uint32(r8) << 8) | uint32(r8)
	sg := (uint32(g8) << 8) | uint32(g8)
	sb := (uint32(b8) << 8) | uint32(b8)

	if src[3] == 0xffff {
		return NRGBA{sr, sg, sb, src[3]}
	}

	// Compose
	sr = (sr * src[3]) / 0xffff
	sg = (sg * src[3]) / 0xffff
	sb = (sb * src[3]) / 0xffff

	dr := (dst[0] * dst[3]) / 0xffff
	dg := (dst[1] * dst[3]) / 0xffff
	db := (dst[2] * dst[3]) / 0xffff

	r := sr + dr*(0xffff-src[3])/0xffff
	g := sg + dg*(0xffff-src[3])/0xffff
	b := sb + db*(0xffff-src[3])/0xffff
	a := src[3] + dst[3]*(0xffff-src[3])/0xffff

	if a != 0 {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}

	return NRGBA{r, g, b, a}
}

func (o opMulYCC) String() string { return "mulycc" }

type opXorRGB struct{}

func (o opXorRGB) Apply(dst, src NRGBA) NRGBA {
	sr := src[0] ^ dst[0]
	sg := src[1] ^ dst[1]
	sb := src[2] ^ dst[2]

	if src[3] == 0xffff {
		return NRGBA{sr, sg, sb, src[3]}
	}

	// Compose
	sr = (sr * src[3]) / 0xffff
	sg = (sg * src[3]) / 0xffff
	sb = (sb * src[3]) / 0xffff

	dr := (dst[0] * dst[3]) / 0xffff
	dg := (dst[1] * dst[3]) / 0xffff
	db := (dst[2] * dst[3]) / 0xffff

	r := sr + dr*(0xffff-src[3])/0xffff
	g := sg + dg*(0xffff-src[3])/0xffff
	b := sb + db*(0xffff-src[3])/0xffff
	a := src[3] + dst[3]*(0xffff-src[3])/0xffff

	if a != 0 {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}

	return NRGBA{r, g, b, a}
}

func (o opXorRGB) String() string { return "xorrgb" }

type opXorYCC struct{}

func (o opXorYCC) Apply(dst, src NRGBA) NRGBA {
	sY, sCb, sCr := color.RGBToYCbCr(uint8(src[0]>>8), uint8(src[1]>>8), uint8(src[2]>>8))
	dY, dCb, dCr := color.RGBToYCbCr(uint8(dst[0]>>8), uint8(dst[1]>>8), uint8(dst[2]>>8))

	sY = sY ^ dY
	sCb = uint8(((int32(sCb) - 128) ^ (int32(dCb) - 128)) + 128)
	sCr = uint8(((int32(sCr) - 128) ^ (int32(dCr) - 128)) + 128)

	r8, g8, b8 := color.YCbCrToRGB(sY, sCb, sCr)
	sr := (uint32(r8) << 8) | uint32(r8)
	sg := (uint32(g8) << 8) | uint32(g8)
	sb := (uint32(b8) << 8) | uint32(b8)

	if src[3] == 0xffff {
		return NRGBA{sr, sg, sb, src[3]}
	}

	// Compose
	sr = (sr * src[3]) / 0xffff
	sg = (sg * src[3]) / 0xffff
	sb = (sb * src[3]) / 0xffff

	dr := (dst[0] * dst[3]) / 0xffff
	dg := (dst[1] * dst[3]) / 0xffff
	db := (dst[2] * dst[3]) / 0xffff

	r := sr + dr*(0xffff-src[3])/0xffff
	g := sg + dg*(0xffff-src[3])/0xffff
	b := sb + db*(0xffff-src[3])/0xffff
	a := src[3] + dst[3]*(0xffff-src[3])/0xffff

	if a != 0 {
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a
	}

	return NRGBA{r, g, b, a}
}

func (o opXorYCC) String() string { return "xorycc" }

var opsTable = []Operation{
	OpCompose:   opCompose{},
	OpReplace:   opReplace{},
	OpAdd:       opAdd{},
	OpAddRGBMod: opAddRGBMod{},
	OpAddYCCMod: opAddYCCMod{},
	OpMulRGB:    opMulRGB{},
	OpMulYCC:    opMulYCC{},
	OpXorRGB:    opXorRGB{},
	OpXorYCC:    opXorYCC{},
}

func GetOp(op int) Operation {
	if op < len(opsTable) {
		return opsTable[op]
	}
	return nil
}

var opsNamesTable = map[string]Operation{
	"cmp":     opCompose{},
	"src":     opReplace{},
	"add":     opAdd{},
	"addrgbm": opAddRGBMod{},
	"addyccm": opAddYCCMod{},
	"mulrgb":  opMulRGB{},
	"mulycc":  opMulYCC{},
	"xorrgb":  opXorRGB{},
	"xorycc":  opXorYCC{},
}

func GetOpID(op string) Operation {
	return opsNamesTable[op]
}

func OpNames() []string {
	ret := make([]string, 0, len(opsNamesTable))
	for name := range opsNamesTable {
		ret = append(ret, name)
	}
	sort.Strings(ret)
	return ret
}
