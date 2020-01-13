package engine

import (
	"image"
	"testing"
)

func BenchmarkEngine(b *testing.B) {
	img := image.NewNRGBA64(image.Rect(0, 0, 1000, 1000))

	opt := Options{
		MinIterations:  10,
		MaxIterations:  10,
		BlockSize:      16,
		MinSegmentSize: 1,
		MaxSegmentSize: 1,
		MinFilters:     4,
		MaxFilters:     4,
		Threads:        1,
	}

	for i := 0; i < b.N; i++ {
		_, err := opt.Apply(img)
		if err != nil {
			b.Fatal(err)
		}
	}
}
