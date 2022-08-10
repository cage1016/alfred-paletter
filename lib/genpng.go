package lib

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	aw "github.com/deanishe/awgo"
)

func GenPng(wf *aw.Workflow, path string, color color.Color) {
	defer wf.Done()

	newImgRect := image.Rect(0, 0, 1, 1)
	newImg := image.NewRGBA(newImgRect)
	draw.Src.Draw(newImg, newImgRect, &image.Uniform{color}, image.Point{})

	outFile, err := os.Create(path)
	if err != nil {
		log.Fatalf("Could not create file: %s", err)
	}
	defer outFile.Close()
	err = png.Encode(outFile, newImg)
	if err != nil {
		log.Fatalf("Could not encode image: %s", err)
	}
}
