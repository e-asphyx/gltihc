// +build js,wasm
package main

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"syscall/js"

	"github.com/e-asphyx/gltihc/engine"
	log "github.com/sirupsen/logrus"
)

func processImageFunc(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return nil
	}

	// Get image
	srcObj := args[0]
	src := make([]byte, srcObj.Length())
	js.CopyBytesToGo(src, srcObj)

	// Get options
	o := args[1]
	opt := engine.Options{
		MinIterations:  o.Get("minIterations").Int(),
		MaxIterations:  o.Get("maxIterations").Int(),
		BlockSize:      o.Get("blockSize").Int(),
		MinSegmentSize: o.Get("minSegmentSize").Float(),
		MaxSegmentSize: o.Get("maxSegmentSize").Float(),
		MinFilters:     o.Get("minFilters").Int(),
		MaxFilters:     o.Get("maxFilters").Int(),
	}

	if prop := o.Get("filters"); prop.Type() == js.TypeObject {
		opt.Filters = make([]string, prop.Length())
		for i := range opt.Filters {
			opt.Filters[i] = prop.Index(i).String()
		}
	}

	if prop := o.Get("ops"); prop.Type() == js.TypeObject {
		opt.Ops = make([]string, prop.Length())
		for i := range opt.Ops {
			opt.Ops[i] = prop.Index(i).String()
		}
	}

	reader := bytes.NewReader(src)
	sourceImg, _, err := image.Decode(reader)
	if err != nil {
		log.Error(err)
		return err.Error()
	}

	resImg, err := opt.Apply(sourceImg)
	if err != nil {
		log.Error(err)
		return err.Error()
	}

	var outBuf bytes.Buffer
	if err := jpeg.Encode(&outBuf, resImg, &jpeg.Options{Quality: 95}); err != nil {
		log.Error(err)
		return err.Error()
	}

	dstObj := js.Global().Get("Uint8Array").New(outBuf.Len())
	js.CopyBytesToJS(dstObj, outBuf.Bytes())

	return dstObj
}

func main() {
	js.Global().Set("_gltihcProcessImage", js.FuncOf(processImageFunc))

	if v := js.Global().Get("_gltihcInitDone"); v.Type() == js.TypeFunction {
		v.Invoke()
	}

	select {}
}
