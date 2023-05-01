package uimage

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func Decode(file io.Reader, contentType string) (image.Image, error) {
	var img image.Image
	var err error

	switch contentType {
	case "image/png":
		img, err = png.Decode(file)
		if err != nil {
			panic(err)
		}
	case "image/jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			panic(err)
		}
	default:
		return nil, errors.New("unsupported")
	}

	return img, nil
}
