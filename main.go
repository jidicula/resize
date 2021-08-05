package main

import (
	"fmt"
	"image"
	"os"

	flag "github.com/spf13/pflag"

	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
)

var usage = `Usage: resize [options...] <input PNG image> [<output PNG image>]

resize is a tool for resizing a PNG image by a provided factor.

If no factor is specified, defaults to factor 1. If no output filename is
provided, resized image is written to ./out.png

Examples:

    $ resize --factor=2 testdata/2.png
    $ resize -f=2 testdata/1.png new.png
    $ resize -f 2 testdata/1.png new.png

Options:
  --factor (-f) Resize factor.
  --help (-h) Prints this message.
`

var factor int

func init() {
	flag.IntVarP(&factor, "factor", "f", 1, "Resize factor.")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputFile := flag.Arg(0)

	// Set output name
	outputFile := "out.png"
	if flag.Arg(1) != "" {
		outputFile = flag.Arg(1)
	}

	// Decode image from file
	reader, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	bounds := m.Bounds()

	maxY := bounds.Max.Y
	maxX := bounds.Max.X

	newMaxY := factor * maxY
	newMaxX := factor * maxX

	newImage := image.NewRGBA(image.Rect(0, 0, newMaxX, newMaxY))

	// TODO: functionalize into both single thread and concurrent implementations
	for y := 0; y < newMaxY; y++ {
		for x := 0; x < newMaxX; x++ {
			r, g, b, a := m.At(x/factor, y/factor).RGBA()
			newImage.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	defer f.Close()

	err = png.Encode(f, newImage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
