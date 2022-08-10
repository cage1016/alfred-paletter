package lib

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

func Download(url, path string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch image: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch %s: %s", url, resp.Status)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file: %s", err)
	}

	return nil
}

func DecodeBase64(str, dataDir string) (string, error) {
	coI := strings.Index(string(str), ",")
	rawImage := string(str)[coI+1:]

	unBased, _ := base64.StdEncoding.DecodeString(string(rawImage))
	res := bytes.NewReader(unBased)
	switch strings.TrimSuffix(str[5:coI], ";base64") {
	case "image/png":
		pngI, err := png.Decode(res)
		if err != nil {
			return "", fmt.Errorf("failed to decode image: %s", err)
		}

		path := fmt.Sprintf("%s/d.png", dataDir)
		file, err := os.Create(path)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %s", err)
		}
		defer file.Close()

		err = png.Encode(file, pngI)
		if err != nil {
			return "", fmt.Errorf("failed to encode png: %s", err)
		}
		return path, nil
	case "image/jpeg":
		jpgI, err := jpeg.Decode(res)
		if err != nil {
			return "", fmt.Errorf("failed to decode jpeg: %s", err)
		}

		path := fmt.Sprintf("%s/d.jpeg", dataDir)
		file, err := os.Create(path)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %s", err)
		}
		defer file.Close()

		err = jpeg.Encode(file, jpgI, nil)
		if err != nil {
			return "", fmt.Errorf("failed to encode jpeg: %s", err)
		}
		return path, nil
	case "image/gif":
		gifI, err := gif.Decode(res)
		if err != nil {
			return "", fmt.Errorf("failed to decode gif: %s", err)
		}

		path := fmt.Sprintf("%s/d.gif", dataDir)
		file, err := os.Create(path)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %s", err)
		}

		err = gif.Encode(file, gifI, nil)
		if err != nil {
			return "", fmt.Errorf("failed to encode gif: %s", err)
		}

		return path, nil
	}
	return "", fmt.Errorf("unsupported image type")
}
