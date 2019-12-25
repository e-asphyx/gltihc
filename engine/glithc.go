package engine

import (
	"errors"
	"fmt"
	"image"
	"math/rand"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
)

type Options struct {
	Iterations     int
	BlockSize      int
	MinSegmentSize float64
	MaxSegmentSize float64
	MinFilters     int
	MaxFilters     int
}

var (
	ErrOptions       = errors.New("invalid engine options")
	ErrImageTooSmall = errors.New("image is too small")
)

func Apply(img image.Image, opt *Options) (image.Image, error) {
	if opt.BlockSize <= 0 || opt.Iterations <= 0 ||
		opt.MinSegmentSize > 1 || opt.MaxSegmentSize > 1 ||
		opt.MinSegmentSize < 0 || opt.MaxSegmentSize < opt.MinSegmentSize ||
		opt.MinFilters <= 0 || opt.MaxFilters < opt.MinFilters {
		return nil, ErrOptions
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	imageW := img.Bounds().Dx()
	imageH := img.Bounds().Dy()

	src := image.NewRGBA64(image.Rect(0, 0, imageW, imageH))
	dst := image.NewRGBA64(image.Rect(0, 0, imageW, imageH))
	tmp0 := image.NewRGBA64(image.Rect(0, 0, imageW, imageH))
	tmp1 := image.NewRGBA64(image.Rect(0, 0, imageW, imageH))
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)

	for itn := 0; itn < opt.Iterations; itn++ {
		// Copy back
		draw.Draw(src, src.Bounds(), dst, dst.Bounds().Min, draw.Src)

		blocksX := imageW / opt.BlockSize
		blocksY := imageH / opt.BlockSize
		blocks := blocksX * blocksY
		if blocks == 0 {
			return nil, ErrImageTooSmall
		}

		if float64(blocks)*opt.MinSegmentSize < 1 {
			return nil, ErrImageTooSmall
		}

		p := opt.MinSegmentSize + rnd.Float64()*(opt.MaxSegmentSize-opt.MinSegmentSize)
		segBlocks := int(float64(blocks) * p)
		segStart := rnd.Intn(blocks - segBlocks + 1)

		var segShift int
		if rnd.Intn(2) == 1 {
			// Apply shift
			segShift = rnd.Intn(blocks)
		}

		// Clear intermediate images
		draw.Draw(tmp0, tmp0.Bounds(), image.Transparent, tmp0.Bounds().Min, draw.Src)
		draw.Draw(tmp1, tmp1.Bounds(), image.Transparent, tmp1.Bounds().Min, draw.Src)

		filtersNum := opt.MinFilters + rnd.Intn(opt.MaxFilters-opt.MinFilters+1)
		filters := make([]Filter, filtersNum)

		fo := FilterOptions{
			BlockSize: opt.BlockSize,
			Reference: src,
		}
		for i := range filters {
			n := rnd.Intn(FilterNumFilters)
			filters[i] = NewRandomizedFilter(n, &fo, rnd)
		}

		ops := make([]Operation, filtersNum)
		for i := 0; i < filtersNum; i++ {
			if i < filtersNum-1 {
				ops[i] = GetOp(OpReplace)
			} else {
				opn := rnd.Intn(OpNumOps)
				ops[i] = GetOp(opn)
			}

		}

		if log.IsLevelEnabled(log.DebugLevel) {
			fs := make([]string, len(filters))
			for i, f := range filters {
				fs[i] = fmt.Sprintf("{%v,%v}", f, ops[i])
			}
			log.Debugf("iter: %d, shift: %d, filters: [%s]", itn, segShift, strings.Join(fs, ","))
		}

		for fc := 0; fc < filtersNum; fc++ {
			var (
				ss, dd draw.Image
			)
			if fc == 0 {
				ss = src
			} else {
				ss = tmp0
			}
			if fc < filtersNum-1 {
				dd = tmp1
			} else {
				dd = dst
			}

			// Apply block by block
			for b := segStart; b < segStart+segBlocks; b++ {
				sb := b
				if fc == 0 {
					// Apply shift
					sb = (b + segShift) % blocks
				}

				dx, dy := (b%blocksX)*opt.BlockSize, (b/blocksX)*opt.BlockSize
				dr := image.Rect(dx, dy, dx+opt.BlockSize, dy+opt.BlockSize)
				sp := image.Point{(sb % blocksX) * opt.BlockSize, (sb / blocksX) * opt.BlockSize}
				filters[fc].Apply(dd, dr, ss, sp, ops[fc])
			}
			tmp0, tmp1 = tmp1, tmp0
		}
	}

	ret := image.NewRGBA(image.Rect(0, 0, dst.Bounds().Dx(), dst.Bounds().Dy()))
	draw.Draw(ret, ret.Bounds(), dst, dst.Bounds().Min, draw.Src)

	return ret, nil
}
