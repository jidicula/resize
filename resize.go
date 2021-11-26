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
  --concurrent (-c) Resize concurrently.
  --help (-h) Prints this message.

`

var factor int
var concurrent bool

func init() {
	flag.IntVarP(&factor, "factor", "f", 1, "Resize factor.")
	flag.BoolVarP(&concurrent, "concurrent", "c", false, "Resize concurrently.")
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
		os.Exit(2)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(3)
	}
	bounds := m.Bounds()

	maxY := bounds.Max.Y
	maxX := bounds.Max.X

	newMaxY := factor * maxY
	newMaxX := factor * maxX

	var newImage image.RGBA
	if concurrent {
		newImage = concurrentResize(newMaxY, newMaxX, m)
	} else {
		newImage = singleThreadResize(newMaxY, newMaxX, m)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(4)
	}
	defer f.Close()

	err = png.Encode(f, &newImage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(5)
	}
}

// singleThreadResize resizes an Image in a single thread.
func singleThreadResize(newMaxY int, newMaxX int, inputImage image.Image) image.RGBA {
	fmt.Printf("Resizing...\n")

	newImage := *image.NewRGBA(image.Rect(0, 0, newMaxX, newMaxY))
	for y := 0; y < newMaxY; y++ {
		for x := 0; x < newMaxX; x++ {
			r, g, b, a := inputImage.At(x/factor, y/factor).RGBA()
			newImage.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return newImage
}

// concurrentResize ...
func concurrentResize(newMaxY int, newMaxX int, inputImage image.Image) image.RGBA {
	fmt.Printf("Concurrently resizing...\n")
	newImage := *image.NewRGBA(image.Rect(0, 0, newMaxX, newMaxY))
	for y := 0; y < newMaxY; y++ {
		for x := 0; x < newMaxX; x++ {
			go func(x int, y int) {
				r, g, b, a := inputImage.At(x/factor, y/factor).RGBA()
				newImage.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			}(x, y)
		}
	}
	return newImage
}
