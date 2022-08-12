package lib

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

var m = map[string]Handler{
	"image/png":  {"png", png.Decode, png.Encode},
	"image/jpeg": {"jpeg", jpeg.Decode, func(w io.Writer, m image.Image) error { return jpeg.Encode(w, m, nil) }},
	"image/gif":  {"gif", gif.Decode, func(w io.Writer, m image.Image) error { return gif.Encode(w, m, nil) }},
	"image/webp": {"webp", webp.Decode, png.Encode},
	"image/bmp":  {"bmp", bmp.Decode, bmp.Encode},
	"image/tiff": {"tiff", tiff.Decode, func(w io.Writer, m image.Image) error { return tiff.Encode(w, m, nil) }},
}

type Handler struct {
	Ext    string
	Decode func(r io.Reader) (image.Image, error)
	Encode func(w io.Writer, m image.Image) error
}

func DecodeBase64(str, dataDir string) (string, error) {
	coI := strings.Index(string(str), ",")
	rawImage := string(str)[coI+1:]

	unBased, _ := base64.StdEncoding.DecodeString(string(rawImage))
	res := bytes.NewReader(unBased)

	if h, ok := m[strings.TrimSuffix(str[5:coI], ";base64")]; ok {
		img, err := h.Decode(res)
		if err != nil {
			return "", fmt.Errorf("failed to decode image: %s", err)
		}

		path := fmt.Sprintf("%s/d.%s", dataDir, h.Ext)
		file, err := os.Create(path)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %s", err)
		}
		defer file.Close()

		err = h.Encode(file, img)
		if err != nil {
			return "", fmt.Errorf("failed to encode image: %s", err)
		}
		return path, nil
	}
	return "", fmt.Errorf("unsupported image type")
}
