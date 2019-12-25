package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/e-asphyx/gltihc/engine"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		opt     engine.Options
		prefix  string
		copies  int
		debug   bool
		filters string
		ops     string
	)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <input>\n\nOptions:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nFilters:\n  color, gray, src, rgba, seta, ycc, prgb, prgba, pycc, cop, ctoa, mix, qrgba, qycc, inv, invrgba, invycc, gs, rasp\n\nOps:\n  cmp, src, add, addrgbm, addyccm, mulrgb, mulycc, xorrgb, xorycc\n")
	}

	flag.BoolVar(&debug, "debug", false, "Debug")
	flag.IntVar(&copies, "copies", 1, "Copies")
	flag.StringVar(&prefix, "prefix", "", "Prefix")
	flag.IntVar(&opt.Iterations, "iter", 1, "Iterations")
	flag.IntVar(&opt.BlockSize, "bs", 16, "BlockSize")
	flag.Float64Var(&opt.MinSegmentSize, "min-segment-size", 0.01, "Minimum segment size relative to image size")
	flag.Float64Var(&opt.MaxSegmentSize, "max-segment-size", 0.2, "Maximun segment size relative to image size")
	flag.IntVar(&opt.MinFilters, "min-filters", 1, "Minimum filters number in a chain")
	flag.IntVar(&opt.MaxFilters, "max-filters", 1, "Maximun filters number in a chain")
	flag.StringVar(&filters, "filters", "", "Allowed filters")
	flag.StringVar(&ops, "ops", "", "Allowed ops")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if filters != "" {
		opt.Filters = strings.Split(filters, ",")
	}

	if ops != "" {
		opt.Ops = strings.Split(ops, ",")
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
		log.Debugf("copy: %d", c)

		res, err := engine.Apply(source, &opt)
		if err != nil {
			log.Fatal(err)
		}

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
