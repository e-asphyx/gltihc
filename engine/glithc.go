package engine

import (
	"image"
	"math/rand"
	"time"

	"golang.org/x/image/draw"
)

type Options struct {
	Iterations     int
	BlockSize      int
	MinSegmentSize float64
	MaxSegmentSize float64
}

func Apply(img image.Image, opt *Options) image.Image {
	src := image.NewRGBA64(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	dst := image.NewRGBA64(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	mins := opt.MinSegmentSize
	maxs := opt.MaxSegmentSize

	if mins < 0 {
		mins = 0
	} else if mins > 1.0 {
		mins = 1.0
	}
	if maxs < 0 {
		maxs = 0
	} else if maxs > 1.0 {
		maxs = 1.0
	}

	if mins > maxs {
		mins, maxs = maxs, mins
	}

	for itn := 0; itn < opt.Iterations; itn++ {
		srcBounds := src.Bounds()
		dstBounds := dst.Bounds()
		// Copy back
		draw.Draw(src, srcBounds, dst, dstBounds.Min, draw.Src)

		blocksX := srcBounds.Dx() / opt.BlockSize
		blocksY := srcBounds.Dy() / opt.BlockSize
		blocks := blocksX * blocksY

		p := mins + rnd.Float64()*(maxs-mins)
		segBlocks := int(float64(blocks) * p)
		if segBlocks == 0 {
			continue
		}
		segStart := rnd.Intn(blocks - segBlocks + 1)

		var segShift int
		if rnd.Intn(2) == 1 {
			// Apply shift
			segShift = rnd.Intn(blocks*2) - blocks
		}

		filtern := rnd.Intn(FilterNumFilters)
		fo := FilterOptions{
			BlockSize: opt.BlockSize,
			Reference: src,
		}
		filter := NewRandomizedFilter(filtern, &fo, rnd)

		opn := rnd.Intn(OpNumOps)
		op := GetOp(opn)

		// Apply block by block
		for b := segStart; b < segStart+segBlocks; b++ {
			sb := (b + segShift) % blocks
			if sb < 0 {
				sb += blocks
			}
			dx, dy := (b%blocksX)*opt.BlockSize, (b/blocksX)*opt.BlockSize
			dr := image.Rect(dx+dstBounds.Min.X, dy+dstBounds.Min.Y, dx+opt.BlockSize+dstBounds.Min.X, dy+opt.BlockSize+dstBounds.Min.Y)
			sp := image.Point{(sb%blocksX)*opt.BlockSize + srcBounds.Min.X, (sb/blocksX)*opt.BlockSize + srcBounds.Min.Y}
			filter.Apply(dst, dr, src, sp, op)
		}
	}

	ret := image.NewRGBA(image.Rect(0, 0, dst.Bounds().Dx(), dst.Bounds().Dy()))
	draw.Draw(ret, ret.Bounds(), dst, dst.Bounds().Min, draw.Src)

	return ret
}
