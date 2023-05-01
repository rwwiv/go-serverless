package image

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func Encode(w io.Writer, img image.Image, contentType string) error {
	var err error

	switch contentType {
	case "image/png":
		err = png.Encode(w, img)
		if err != nil {
			return err
		}
	case "image/jpeg":
		err = jpeg.Encode(w, img, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
