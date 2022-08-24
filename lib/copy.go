package lib

import (
	"fmt"
	"io"
	"os"
)

func Copy(src, dataDir string) (string, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %s", err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %s", err)
	}
	defer source.Close()

	dst := fmt.Sprintf("%s/d.tiff", dataDir)
	destination, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %s", err)
	}

	return dst, nil
}
