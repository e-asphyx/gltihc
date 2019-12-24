package glithc

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

type Operation func(dst, src color.Color) color.Color

func opCompose(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	a := sa + da*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func opReplace(dst, src color.Color) color.Color {
	return src
}

func opAdd(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	a := sa + da*(0xffff-sa)/0xffff
	sr = (sr*0xffff)/sa + (dr*0xffff)/da
	sg = (sg*0xffff)/sa + (dg*0xffff)/da
	sb = (sb*0xffff)/sa + (db*0xffff)/da
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

func opAddRGBMod(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	a := sa + da*(0xffff-sa)/0xffff
	sr = ((sr*0xffff)/sa + (dr*0xffff)/da) & 0xffff
	sg = ((sg*0xffff)/sa + (dg*0xffff)/da) & 0xffff
	sb = ((sb*0xffff)/sa + (db*0xffff)/da) & 0xffff
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

func opAddYCCMod(dst, src color.Color) color.Color {
	s := color.NYCbCrAModel.Convert(src).(color.NYCbCrA)
	d := color.NYCbCrAModel.Convert(dst).(color.NYCbCrA)
	s.Y = (s.Y + d.Y) & 0xff
	s.Cb = (s.Cb + d.Cb) & 0xff
	s.Cr = (s.Cr + d.Cr) & 0xff
	// Compose
	sr, sg, sb, sa := s.RGBA()
	dr, dg, db, da := dst.RGBA()
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	a := sa + da*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func opMulRGB(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	a := sa + da*(0xffff-sa)/0xffff
	sr = (((sr * 0xffff) / sa) * ((dr * 0xffff) / da)) / 0xffff
	sg = (((sg * 0xffff) / sa) * ((dg * 0xffff) / da)) / 0xffff
	sb = (((sb * 0xffff) / sa) * ((db * 0xffff) / da)) / 0xffff
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

func opMulYCC(dst, src color.Color) color.Color {
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

func opXorRGB(dst, src color.Color) color.Color {
	sr, sg, sb, sa := src.RGBA()
	dr, dg, db, da := dst.RGBA()
	a := sa + da*(0xffff-sa)/0xffff
	sr = ((sr * 0xffff) / sa) ^ ((dr * 0xffff) / da)
	sg = ((sg * 0xffff) / sa) ^ ((dg * 0xffff) / da)
	sb = ((sb * 0xffff) / sa) ^ ((db * 0xffff) / da)
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

func opXorYCC(dst, src color.Color) color.Color {
	s := color.NYCbCrAModel.Convert(src).(color.NYCbCrA)
	d := color.NYCbCrAModel.Convert(dst).(color.NYCbCrA)
	s.Y = s.Y ^ d.Y
	s.Cb = uint8((int32(s.Cb) - 128) ^ (int32(d.Cb) - 128) + 128)
	s.Cr = uint8((int32(s.Cr) - 128) ^ (int32(d.Cr) - 128) + 128)
	// Compose
	sr, sg, sb, sa := s.RGBA()
	dr, dg, db, da := dst.RGBA()
	r := sr + dr*(0xffff-sa)/0xffff
	g := sg + dg*(0xffff-sa)/0xffff
	b := sb + db*(0xffff-sa)/0xffff
	a := sa + da*(0xffff-sa)/0xffff
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

var ops = []Operation{
	OpCompose:   opCompose,
	OpReplace:   opReplace,
	OpAdd:       opAdd,
	OpAddRGBMod: opAddRGBMod,
	OpAddYCCMod: opAddYCCMod,
	OpMulRGB:    opMulRGB,
	OpMulYCC:    opMulYCC,
	OpXorRGB:    opXorRGB,
	OpXorYCC:    opXorYCC,
}

func GetOp(op int) Operation {
	if op < len(ops) {
		return ops[op]
	}
	return nil
}
