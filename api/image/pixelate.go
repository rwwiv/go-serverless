package image

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"net/http"
	"strconv"

	uimage "github.com/rwwiv/go-serverless/util/uimage"
)

func PixelateHandler(w http.ResponseWriter, r *http.Request) {
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

	pixelSize, err := strconv.Atoi(r.PostFormValue("pixel_size"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	processedImg := pixelate(img, pixelSize)

	outputBuf := new(bytes.Buffer)
	err = uimage.Encode(outputBuf, processedImg, contentType)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(outputBuf.Bytes())
}

func pixelate(img image.Image, pixelSize int) *image.RGBA {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	pixelatedImg := image.NewRGBA(bounds)

	// if pixelSize is 1, just copy the image
	if pixelSize <= 1 {
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				pixel := img.At(x, y)
				col := color.RGBAModel.Convert(pixel).(color.RGBA)
				pixelatedImg.Set(x, y, col)
			}
		}
		return pixelatedImg
	}

	var sizeFillsImg bool
	// clamp pixelSize to image size
	if width > height && pixelSize > width {
		pixelSize = width
		sizeFillsImg = true
	} else if height >= width && pixelSize > height {
		pixelSize = height
		sizeFillsImg = true
	}

	// if pixelSize is as large as the image, just fill the image with the average color
	if sizeFillsImg {
		r, g, b := meanAvgColor(img)
		col := color.RGBA{r, g, b, 255}
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				pixelatedImg.Set(x, y, col)
			}
		}
		return pixelatedImg
	}

	centerX := width / 2
	centerY := height / 2

	// start at the top left corner of the image
	startX := centerX % pixelSize
	if startX != 0 {
		startX -= pixelSize
	}
	startY := centerY % pixelSize
	if startY != 0 {
		startY -= pixelSize
	}

	for x := startX; x < width; x += pixelSize {
		for y := startY; y < height; y += pixelSize {
			rect := image.Rect(x, y, x+pixelSize, y+pixelSize)

			// clamp rect to image bounds
			if rect.Min.X < 0 {
				rect.Min.X = 0
			}
			if rect.Min.Y < 0 {
				rect.Min.Y = 0
			}
			if rect.Max.X > width {
				rect.Max.X = width
			}
			if rect.Max.Y > height {
				rect.Max.Y = height
			}

			r, g, b := meanAvgColorOverRect(img, rect)
			col := color.RGBA{r, g, b, 255}

			// fill rect with average color
			for x2 := rect.Min.X; x2 <= rect.Max.X; x2++ {
				for y2 := rect.Min.Y; y2 <= rect.Max.Y; y2++ {
					pixelatedImg.Set(x2, y2, col)
				}
			}
		}
	}

	return pixelatedImg
}

func meanAvgColor(img image.Image) (red, green, blue uint8) {
	bounds := img.Bounds()
	return meanAvgColorOverRect(img, bounds)
}

func meanAvgColorOverRect(img image.Image, rect image.Rectangle) (red, green, blue uint8) {
	var redSum float64
	var greenSum float64
	var blueSum float64

	// sum up all the colors in the rect
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		for y := rect.Min.Y; y <= rect.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			redSum += float64(r/256) * float64(r/256)
			greenSum += float64(g/256) * float64(g/256)
			blueSum += float64(b/256) * float64(b/256)
		}
	}

	rectArea := float64((rect.Dx() + 1) * (rect.Dy() + 1))

	// average the sums
	red = uint8(math.Round(math.Sqrt(redSum / rectArea)))
	green = uint8(math.Round(math.Sqrt(greenSum / rectArea)))
	blue = uint8(math.Round(math.Sqrt(blueSum / rectArea)))

	return
}
