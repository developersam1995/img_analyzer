package main

import (
	"image"
	"io"
	"os"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func upscale(fileName string) (i io.Reader, err error) {
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

	resized := transform.Resize(img, c.Width*4, c.Height*4, transform.Linear)
	pngEncoder := imgio.PNGEncoder()

	r, w := io.Pipe()
	go func() {
		pngEncoder(w, resized)
		w.Close()
	}()
	return r, nil
}
