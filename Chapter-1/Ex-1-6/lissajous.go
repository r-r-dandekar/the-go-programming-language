// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand/v2"
	"os"
)

var palette = []color.Color{color.Black, 
							color.RGBA{0x00, 0xFF, 0x00, 0xFF}, 
							color.RGBA{0x22, 0xDD, 0x00, 0xFF}, 
							color.RGBA{0x44, 0xBB, 0x00, 0xFF}, 
							color.RGBA{0x66, 0x99, 0x00, 0xFF}, 
							color.RGBA{0x99, 0x66, 0x00, 0xFF}, 
							color.RGBA{0xBB, 0x44, 0x00, 0xFF}, 
							color.RGBA{0xDD, 0x22, 0x00, 0xFF}, 
							color.RGBA{0xFF, 0x00, 0x00, 0xFF}}

const (
	blackIndex = 0
	numColors = 8
)

func main() {
	f, err := os.Create("out.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lissajous(f)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			var primaryIndex uint8
			if i%(numColors*2) < numColors {
				primaryIndex = uint8(1+(i%numColors))
			} else {
				primaryIndex = uint8(numColors - i%numColors)
			}
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), primaryIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
