package engine

import "image/color"

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

type Operation interface {
	Apply(dst, src color.Color) color.Color
	String() string
}

type opCompose struct{}

func (o opCompose) Apply(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	a := sa + da*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (o opCompose) String() string { return "comp" }

type opReplace struct{}

func (o opReplace) Apply(dst, src color.Color) color.Color { return src }

func (o opReplace) String() string { return "rep" }

type opAdd struct{}

func (o opAdd) Apply(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	a := sa + da*(0xffff-sa)/0xffff
	if sa != 0 {
		sr = (sr * 0xffff) / sa
		sg = (sg * 0xffff) / sa
		sb = (sb * 0xffff) / sa
	}
	if da != 0 {
		sr += (dr * 0xffff) / da
		sg += (dg * 0xffff) / da
		sb += (db * 0xffff) / da
	}
	if sr > 0xffff {
		sr = 0xffff
	}
	if sg > 0xffff {
		sg = 0xffff
	}
	if sb > 0xffff {
		sb = 0xffff
	}
	// Premultiply
	sr = (sr * a) / 0xffff
	sg = (sg * a) / 0xffff
	sb = (sb * a) / 0xffff
	// Compose
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (o opAdd) String() string { return "add" }

type opAddRGBMod struct{}

func (o opAddRGBMod) Apply(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	a := sa + da*(0xffff-sa)/0xffff
	if sa != 0 {
		sr = (sr * 0xffff) / sa
		sg = (sg * 0xffff) / sa
		sb = (sb * 0xffff) / sa
	}
	if da != 0 {
		sr += (dr * 0xffff) / da
		sg += (dg * 0xffff) / da
		sb += (db * 0xffff) / da
	}
	sr &= 0xffff
	sg &= 0xffff
	sb &= 0xffff
	// Premultiply
	sr = (sr * a) / 0xffff
	sg = (sg * a) / 0xffff
	sb = (sb * a) / 0xffff
	// Compose
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (o opAddRGBMod) String() string { return "addrgbm" }

type opAddYCCMod struct{}

func (o opAddYCCMod) Apply(dst, src color.Color) color.Color {
	s := color.NYCbCrAModel.Convert(src).(color.NYCbCrA)
	d := color.NYCbCrAModel.Convert(dst).(color.NYCbCrA)
	s.Y = (s.Y + d.Y) & 0xff
	s.Cb = uint8((int32(s.Cb) + int32(d.Cb) - 128) & 0xff)
	s.Cr = uint8((int32(s.Cr) + int32(d.Cr) - 128) & 0xff)
	// Compose
	sr, sg, sb, sa := s.RGBA()
	dr, dg, db, da := dst.RGBA()
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	a := sa + da*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (o opAddYCCMod) String() string { return "addyccm" }

type opMulRGB struct{}

func (o opMulRGB) Apply(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	a := sa + da*(0xffff-sa)/0xffff
	if sa != 0 && da != 0 {
		sr = (((sr * 0xffff) / sa) * ((dr * 0xffff) / da)) / 0xffff
		sg = (((sg * 0xffff) / sa) * ((dg * 0xffff) / da)) / 0xffff
		sb = (((sb * 0xffff) / sa) * ((db * 0xffff) / da)) / 0xffff
	} else {
		sr, sg, sb = 0, 0, 0
	}
	// Premultiply
	sr = (sr * a) / 0xffff
	sg = (sg * a) / 0xffff
	sb = (sb * a) / 0xffff
	// Compose
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (o opMulRGB) String() string { return "mulrgb" }

type opMulYCC struct{}

func (o opMulYCC) Apply(dst, src color.Color) color.Color {
	s := color.NYCbCrAModel.Convert(src).(color.NYCbCrA)
	d := color.NYCbCrAModel.Convert(dst).(color.NYCbCrA)
	s.Y = uint8(uint32(s.Y) * uint32(d.Y) / 0xff)
	s.Cb = uint8((int32(s.Cb)-128)*(int32(d.Cb)-128)/0xff + 128)
	s.Cr = uint8((int32(s.Cr)-128)*(int32(d.Cr)-128)/0xff + 128)
	// Compose
	sr, sg, sb, sa := s.RGBA()
	dr, dg, db, da := dst.RGBA()
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	a := sa + da*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (o opMulYCC) String() string { return "mulycc" }

type opXorRGB struct{}

func (o opXorRGB) Apply(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	a := sa + da*(0xffff-sa)/0xffff
	if sa != 0 {
		sr = ((sr * 0xffff) / sa)
		sg = ((sg * 0xffff) / sa)
		sb = ((sb * 0xffff) / sa)
	}
	if da != 0 {
		sr ^= ((dr * 0xffff) / da)
		sg ^= ((dg * 0xffff) / da)
		sb ^= ((db * 0xffff) / da)
	}
	// Premultiply
	sr = (sr * a) / 0xffff
	sg = (sg * a) / 0xffff
	sb = (sb * a) / 0xffff
	// Compose
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func (o opXorRGB) String() string { return "xorrgb" }

type opXorYCC struct{}

func (o opXorYCC) Apply(dst, src color.Color) color.Color {
	s := color.NYCbCrAModel.Convert(src).(color.NYCbCrA)
	d := color.NYCbCrAModel.Convert(dst).(color.NYCbCrA)
	s.Y = s.Y ^ d.Y
	s.Cb = uint8(((int32(s.Cb) - 128) ^ (int32(d.Cb) - 128)) + 128)
	s.Cr = uint8(((int32(s.Cr) - 128) ^ (int32(d.Cr) - 128)) + 128)
	// Compose
	sr, sg, sb, sa := s.RGBA()
	dr, dg, db, da := dst.RGBA()
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	a := sa + da*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
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

func GetOpByName(op string) Operation {
	return opsNamesTable[op]
}
