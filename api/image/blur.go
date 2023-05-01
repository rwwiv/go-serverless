package api

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/esimov/stackblur-go"
	"github.com/rwwiv/go-serverless/util"
)

func BlurHandler(w http.ResponseWriter, r *http.Request) {
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

	blurStrength, err := strconv.Atoi(r.PostFormValue("blur_strength"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	processedImg, err := stackblur.Process(img, uint32(blurStrength))
	if err != nil {
		panic(err)
	}

	outputBuf := new(bytes.Buffer)
	err = util.Encode(outputBuf, processedImg, contentType)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(outputBuf.Bytes())
}
