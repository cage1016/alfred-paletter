package lib

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func Download(url, dataDir string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %s", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch %s: %s", url, resp.Status)
	}

	buff := bytes.NewBuffer(nil)
	bodyBytes, err := ioutil.ReadAll(io.TeeReader(resp.Body, buff))
	if err != nil {
		return "", fmt.Errorf("failed to read image: %s", err)
	}

	_, format, err := image.DecodeConfig(buff)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/d.%s", dataDir, format)
	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err)
	}
	defer file.Close()

	_, err = io.Copy(file, (bytes.NewReader(bodyBytes)))
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %s", err)
	}

	return path, nil
}
