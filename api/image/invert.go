package image

import (
	"bytes"
	"image"
	"image/color"
	"net/http"

	uimage "github.com/rwwiv/go-serverless/util/image"
)

func InvertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // max 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := r.PostFormValue("content_type")

	img, err := uimage.Decode(file, contentType)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	processedImg := invert(img)

	outputBuf := new(bytes.Buffer)
	err = uimage.Encode(outputBuf, processedImg, contentType)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(outputBuf.Bytes())
}

func invert(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	imgSet := image.NewRGBA(bounds)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			col := color.RGBA{255 - uint8(r/256), 255 - uint8(g/256), 255 - uint8(b/256), uint8(a / 256)}
			imgSet.Set(x, y, col)
		}
	}
	return imgSet
}
