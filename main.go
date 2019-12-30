package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/e-asphyx/gltihc/engine"
	log "github.com/sirupsen/logrus"
)

type preset struct {
	filters []string
	ops     []string
}

var presets = map[string]preset{
	"tame": {
		filters: []string{
			"color",
			"gray",
			"src",
			"rgba",
			"seta",
			"ycc",
			"prgb",
			"prgba",
			"pycc",
			"copy",
			"ctoa",
			"mix",
			"quant",
			"qrgba",
			"qycca",
			"qy",
			"inv",
			"gs",
			"rasp",
		},
		ops: []string{
			"cmp",
			"src",
			"add",
			"mulrgb",
			"mulycc",
		},
	},
	"nocolorshift": {
		filters: []string{
			"gray",
			"src",
			"seta",
			"ctoa",
			"quant",
			"qy",
			"inv",
			"gs",
			"rasp",
		},
		ops: []string{
			"cmp",
			"src",
			"add",
			"mulrgb",
			"mulycc",
			"xorrgb",
			"xorycc",
		},
	},
}

func main() {
	var (
		opt     engine.Options
		prefix  string
		copies  int
		debug   bool
		filters string
		ops     string
		preset  string
		dir     string
	)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <input...>\n\nOptions:\n", path.Base(os.Args[0]))
		flag.PrintDefaults()

		p := make([]string, 0, len(presets))
		for name := range presets {
			p = append(p, name)
		}
		sort.Strings(p)

		fmt.Fprintf(flag.CommandLine.Output(),
			"\nFilters:\n  %s\n\nOperations:\n  %s\n\nPresets:\n  %s\n",
			strings.Join(engine.FilterNames(), ", "),
			strings.Join(engine.OpNames(), ", "),
			strings.Join(p, ", "))
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
	flag.StringVar(&preset, "preset", "", "Preset")
	flag.StringVar(&dir, "dir", "", "Output directory")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if preset != "" {
		p, ok := presets[preset]
		if !ok {
			log.Fatalf("Unknown preset `%s'", preset)
		}
		opt.Filters = p.filters
		opt.Ops = p.ops
	} else {
		if filters != "" {
			opt.Filters = strings.Split(filters, ",")
		}
		if ops != "" {
			opt.Ops = strings.Split(ops, ",")
		}
	}

	if dir != "" {
		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Fatal(err)
		}
	}

	for _, infile := range flag.Args() {
		pfx := prefix
		if pfx == "" {
			b := path.Base(infile)
			if i := strings.IndexByte(b, '.'); i > 0 {
				pfx = b[:i] + "_"
			}
		}

		log.Printf("processing: %s", infile)

		reader, err := os.Open(infile)
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

			res, err := opt.Apply(source)
			if err != nil {
				log.Fatal(err)
			}

			outname := fmt.Sprintf("%s%0*d.png", pfx, ln, c)
			if dir != "" {
				outname = filepath.Join(dir, outname)
			}

			log.Printf("writing: %s", outname)
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
}
