package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path"

	"github.com/e-asphyx/gltihc/engine"
)

func main() {
	var (
		opt    engine.Options
		prefix string
		copies int
	)

	flag.IntVar(&copies, "copies", 1, "Copies")
	flag.StringVar(&prefix, "prefix", "", "Prefix")
	flag.IntVar(&opt.Iterations, "iter", 1, "Iterations")
	flag.IntVar(&opt.BlockSize, "bs", 16, "BlockSize")
	flag.Float64Var(&opt.MinSegmentSize, "min-segment-size", 0.01, "Minimum segment size relative to image size")
	flag.Float64Var(&opt.MaxSegmentSize, "max-segment-size", 0.2, "Maximun segment size relative to image size")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <input>\n\nOptions:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	reader, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	source, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	if err := reader.Close(); err != nil {
		log.Fatal(err)
	}

	var ln int
	for c := copies - 1; c > 0; c = c / 10 {
		ln++
	}
	for c := 0; c < copies; c++ {
		res := engine.Apply(source, &opt)

		outname := fmt.Sprintf("%s%0*d.png", prefix, ln, c)
		f, err := os.Create(outname)
		if err != nil {
			log.Fatal(err)
		}

		if err := png.Encode(f, res); err != nil {
			f.Close()
			log.Fatal(err)
		}

		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
