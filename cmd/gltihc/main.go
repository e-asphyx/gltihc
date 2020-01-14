package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/e-asphyx/gltihc/engine"
	log "github.com/sirupsen/logrus"
	_ "golang.org/x/image/tiff"
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
			"inva",
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

var funcMap = template.FuncMap{
	"base": filepath.Base,
	"basename": func(n string) string {
		n = filepath.Base(n)
		if i := strings.LastIndexByte(n, '.'); i >= 0 {
			return n[:i]
		}
		return n
	},
	"ext": filepath.Ext,
}

type tplContext struct {
	Input       string
	InputCount  int
	NumInputs   int
	CopiesCount int
	NumCopies   int
}

func main() {
	var (
		opt      engine.Options
		format   string
		copies   int
		logLevel string
		filters  string
		ops      string
		preset   string
		dir      string
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

	flag.IntVar(&opt.MinIterations, "min-iterations", 10, "Minimum iterations number")
	flag.IntVar(&opt.MaxIterations, "max-iterations", 10, "Maximum iterations number")
	flag.IntVar(&opt.BlockSize, "bs", 16, "BlockSize")
	flag.Float64Var(&opt.MinSegmentSize, "min-segment-size", 0.01, "Minimum segment size relative to image size")
	flag.Float64Var(&opt.MaxSegmentSize, "max-segment-size", 0.2, "Maximum segment size relative to image size")
	flag.IntVar(&opt.MinFilters, "min-filters", 1, "Minimum filters number in a chain")
	flag.IntVar(&opt.MaxFilters, "max-filters", 1, "Maximum filters number in a chain")
	flag.IntVar(&opt.Threads, "threads", 0, "Number of threads")
	flag.StringVar(&logLevel, "log", "info", "Log level")
	flag.IntVar(&copies, "copies", 1, "Copies")
	flag.StringVar(&format, "fmt", "{{.Input | basename}}_{{printf \"%08d\" .CopiesCount}}.png", "Output file name format")
	flag.StringVar(&filters, "filters", "", "Allowed filters")
	flag.StringVar(&ops, "ops", "", "Allowed ops")
	flag.StringVar(&preset, "preset", "", "Preset")
	flag.StringVar(&dir, "dir", "", "Output directory")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if lv, err := log.ParseLevel(logLevel); err != nil {
		log.Fatal(err)
	} else {
		log.SetLevel(lv)
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

	outTpl, err := template.New("output").Funcs(funcMap).Parse(format)
	if err != nil {
		log.Fatal(err)
	}

	if dir != "" {
		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Fatal(err)
		}
	}

	rand.Seed(time.Now().UnixNano())

	inputs := flag.Args()
	for cnt, infile := range inputs {
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

			var outName strings.Builder
			err := outTpl.Execute(&outName, &tplContext{
				Input:       infile,
				InputCount:  cnt,
				NumInputs:   len(inputs),
				CopiesCount: c,
				NumCopies:   copies,
			})
			if err != nil {
				log.Fatal(err)
			}

			res, err := opt.Apply(source)
			if err != nil {
				log.Fatal(err)
			}

			name := outName.String()
			if dir != "" {
				name = filepath.Join(dir, name)
			}

			log.Printf("writing: %s", name)
			f, err := os.Create(name)
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
