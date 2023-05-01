package api

import (
	"bytes"
	"image"
	"image/color"
	"net/http"

	"github.com/rwwiv/go-serverless/util"
)

func GrayScaleHandler(w http.ResponseWriter, r *http.Request) {
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

	img, err := util.Decode(file, contentType)
	if err != nil {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	processedImg := grayscale(img)

	outputBuf := new(bytes.Buffer)
	err = util.Encode(outputBuf, processedImg, contentType)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(outputBuf.Bytes())
}

func grayscale(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	imgSet := image.NewRGBA(bounds)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			px := color.Gray{uint8(lum / 256)}
			imgSet.Set(x, y, px)
		}
	}
	return imgSet
}
