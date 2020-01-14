package engine

import (
	"errors"
	"fmt"
	"image"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
)

type Options struct {
	MinIterations  int
	MaxIterations  int
	BlockSize      int
	MinSegmentSize float64
	MaxSegmentSize float64
	MinFilters     int
	MaxFilters     int
	Filters        []string
	Ops            []string
	Threads        int
}

var (
	ErrOptions       = errors.New("invalid engine options")
	ErrImageTooSmall = errors.New("image is too small")
)

func copyImage(dst, src *image.NRGBA64) {
	dst.Rect = src.Rect
	dst.Stride = src.Stride
	copy(dst.Pix, src.Pix)
}

func clearStripe(img *image.NRGBA64, y0, y1 int) {
	stripe := (*[(1<<31 - 1) >> 3]uint64)(unsafe.Pointer(&img.Pix[0]))[((y0 * img.Stride) >> 3):((y1 * img.Stride) >> 3):((y1 * img.Stride) >> 3)]
	for i := range stripe {
		stripe[i] = 0
	}
}

func (opt *Options) Apply(img image.Image) (image.Image, error) {
	if opt.BlockSize <= 0 ||
		opt.MinSegmentSize > 1 || opt.MaxSegmentSize > 1 ||
		opt.MinSegmentSize < 0 || opt.MaxSegmentSize < opt.MinSegmentSize ||
		opt.MinFilters <= 0 || opt.MaxFilters < opt.MinFilters ||
		opt.MinIterations < 0 || opt.MaxIterations < opt.MinIterations {
		return nil, ErrOptions
	}

	if opt.Filters != nil {
		for _, f := range opt.Filters {
			if GetFilterID(f) < 0 {
				return nil, fmt.Errorf("unknown filter: %s", f)
			}
		}
	}

	if opt.Ops != nil {
		for _, o := range opt.Ops {
			if GetOpID(o) == nil {
				return nil, fmt.Errorf("unknown op: %s", o)
			}
		}
	}

	threadsNum := opt.Threads
	if threadsNum <= 0 {
		threadsNum = runtime.GOMAXPROCS(0)
	}
	log.Tracef("threadsNum: %d", threadsNum)

	imageW := img.Bounds().Dx()
	imageH := img.Bounds().Dy()

	src := image.NewNRGBA64(image.Rect(0, 0, imageW, imageH))
	dst := image.NewNRGBA64(image.Rect(0, 0, imageW, imageH))
	tmp0 := image.NewNRGBA64(image.Rect(0, 0, imageW, imageH))
	tmp1 := image.NewNRGBA64(image.Rect(0, 0, imageW, imageH))
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)

	iterations := opt.MinIterations + rand.Intn(opt.MaxIterations-opt.MinIterations+1)
	for itn := 0; itn < iterations; itn++ {
		// Copy back
		copyImage(src, dst)

		blocksX := imageW / opt.BlockSize
		blocksY := imageH / opt.BlockSize
		blocks := blocksX * blocksY
		if blocks == 0 {
			return nil, ErrImageTooSmall
		}

		if float64(blocks)*opt.MinSegmentSize < 1 {
			return nil, ErrImageTooSmall
		}

		p := opt.MinSegmentSize + rand.Float64()*(opt.MaxSegmentSize-opt.MinSegmentSize)
		segBlocks := int(float64(blocks) * p)
		if segBlocks == 0 {
			continue
		}
		segStart := rand.Intn(blocks - segBlocks + 1)

		var segShift int
		if rand.Intn(2) == 1 {
			// Apply shift
			segShift = rand.Intn(blocks)
		}

		blocksPerThread := (segBlocks + threadsNum - 1) / threadsNum

		// Clear intermediate images
		stripeY0 := (segStart / blocksX) * opt.BlockSize
		stripeY1 := ((segStart+segBlocks)/blocksX + 1) * opt.BlockSize
		clearStripe(tmp0, stripeY0, stripeY1)
		clearStripe(tmp1, stripeY0, stripeY1)

		filtersNum := opt.MinFilters + rand.Intn(opt.MaxFilters-opt.MinFilters+1)
		filters := make([]Filter, filtersNum)

		fo := FilterOptions{
			BlockSize: opt.BlockSize,
			Reference: src,
		}
		for i := range filters {
			var n int
			if opt.Filters != nil {
				n = rand.Intn(len(opt.Filters))
				n = GetFilterID(opt.Filters[n])
			} else {
				n = rand.Intn(FilterNumFilters)
			}
			if filters[i] = NewRandomizedFilter(n, &fo); filters[i] == nil {
				return nil, ErrOptions
			}
		}

		ops := make([]Operation, filtersNum)
		for i := 0; i < filtersNum; i++ {
			if i < filtersNum-1 {
				ops[i] = GetOp(OpReplace)
			} else if opt.Ops != nil {
				opn := rand.Intn(len(opt.Ops))
				ops[i] = GetOpID(opt.Ops[opn])
			} else {
				opn := rand.Intn(OpNumOps)
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
				ss, dd *image.NRGBA64
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

			var wg sync.WaitGroup
			delta := blocksPerThread
			for tlen, tstart := segBlocks, segStart; tlen > 0; tlen, tstart = tlen-delta, tstart+delta {
				if delta > tlen {
					delta = tlen
				}
				wg.Add(1)
				worker := func(b, ln int) {
					log.Tracef("block: %d..%d, filter: %v", b, b+ln, filters[fc])
					// Apply block by block
					for ; ln > 0; b, ln = b+1, ln-1 {
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
					wg.Done()
				}
				if blocksPerThread == segBlocks {
					worker(tstart, delta)
				} else {
					go worker(tstart, delta)
				}
			}
			wg.Wait()

			tmp0, tmp1 = tmp1, tmp0
		}
	}

	// Convert to 8bpp
	ret := image.NewNRGBA(image.Rect(0, 0, dst.Bounds().Dx(), dst.Bounds().Dy()))
	draw.Draw(ret, ret.Bounds(), dst, dst.Bounds().Min, draw.Src)

	return ret, nil
}
