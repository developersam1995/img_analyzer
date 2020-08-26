package main

import (
	"image"
	"io"
	"os"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

//This function takes an imagefilename and upscale factor, returns an in memory reader for the upscaled image
func upscale(fileName string, factor int) (i io.Reader, err error) {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return
	}

	c, _, err := image.DecodeConfig(f)
	if err != nil {
		return
	}

	img, err := imgio.Open(fileName)
	if err != nil {
		return
	}

	resized := transform.Resize(img, c.Width*factor, c.Height*factor, transform.Linear)
	pngEncoder := imgio.PNGEncoder()

	r, w := io.Pipe()
	go func() {
		pngEncoder(w, resized)
		w.Close()
	}()
	return r, nil
}
