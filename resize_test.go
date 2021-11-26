package main

import (
	"image"
	"image/color"
	"testing"
)

var whitePixel = color.RGBA{
	R: uint8(255),
	G: uint8(255),
	B: uint8(255),
	A: uint8(255),
}

// newImage creates a new white, solid Image.RGBA for testing.
func newImage(size int) image.RGBA {
	r := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: size,
			Y: size,
		},
	}
	im := image.NewRGBA(r)
	for x := 0; x <= size; x++ {
		for y := 0; y <= size; y++ {
			im.SetRGBA(x, y, whitePixel)
		}
	}
	return *im
}

func TestConcurrentResize(t *testing.T) {
	size := 3
	inputImage := newImage(size)
	concurrentResize(size*2, size*2, &inputImage)
}

func BenchmarkConcurrentResize(b *testing.B) {
	size := 300
	inputImage := newImage(size)
	for i := 0; i < b.N; i++ {
		concurrentResize(size*2, size*2, &inputImage)
	}
}

func TestSingleThreadResize(t *testing.T) {
	size := 300
	inputImage := newImage(size)
	singleThreadResize(size*2, size*2, &inputImage)
}

func BenchmarkSingleThreadResize(b *testing.B) {
	size := 300
	inputImage := newImage(size)
	for i := 0; i < b.N; i++ {
		singleThreadResize(size*2, size*2, &inputImage)
	}
}
