package lib

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"

	aw "github.com/deanishe/awgo"
	"github.com/golang/freetype/truetype"
	changecase "github.com/ku/go-change-case"
	"github.com/postprime/graphics-go/graphics"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	template "github.com/cage1016/alfred-paletter/templates"
)

const (
	border            = 0
	paletteHeightPerc = 0.5
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GenPng(wf *aw.Workflow, path string, colors ...color.Color) {
	defer wf.Done()

	newImgRect := image.Rect(0, 0, len(colors), 1)
	newImg := image.NewRGBA(newImgRect)

	for i, c := range colors {
		draw.Src.Draw(newImg, image.Rect(i, 0, 1+i, 1), &image.Uniform{c}, image.Point{})
	}

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

func RemoveAllPng(dir string) {
	pngs, _ := filepath.Glob(filepath.Join(dir, "*"))
	for _, f := range pngs {
		os.Remove(f)
	}
}

// Calculates width of a palette color and total remainder (gap between palette end and original image end)
func calcPaletteWidth(srcWidth int, nColors int) (width int, totalRemainder int) {
	calc := float64(srcWidth-(border*(nColors-1))) / float64(nColors)
	_, remainder := math.Modf(calc)
	return int(math.Floor(calc)), int(math.Round(remainder * float64(nColors)))
}

func colorRects(srcWidth int, srcHeight int, nColors int) []image.Rectangle {
	var ret []image.Rectangle

	// Calculate width and height
	width, remainder := calcPaletteWidth(srcWidth, nColors)
	height := srcHeight - 100

	// Vertical offset is 2 border widths plus the height of the original image
	yOffset := (border * 2) + 0
	for i := 0; i < nColors; i++ {
		xOffset := i*(width+border) + border

		// Cheat by stretching the last color to be flush with the original image
		if i == nColors-1 {
			width += remainder
		}

		ret = append(ret, image.Rect(xOffset, yOffset, width+xOffset, height+yOffset))
	}

	return ret
}

// Luminance computes the luminance (~ brightness) of the given color. Range: 0.0 for black to 1.0 for white.
func Luminance(col color.Color) float64 {
	r, g, b, _ := col.RGBA()
	return (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) / float64(0xffff)
}

func reverseColor(c color.Color) *image.Uniform {
	if Luminance(c) > 0.5 {
		return image.Black
	}
	return image.White
}

func GenPng2(wf *aw.Workflow, path string, colors []color.Color, m map[color.Color]string, ct string) {
	defer wf.Done()

	te := template.NewEngine()
	f, err := truetype.Parse([]byte(te.MustAssetString("fonts/luxisr.ttf")))
	if err != nil {
		log.Println(err)
		return
	}

	newImgRect := image.Rect(0, 0, 1600, 1200)
	newImg := image.NewRGBA(newImgRect)

	// Fill background with white
	draw.Src.Draw(newImg, newImgRect, image.White, image.Point{})

	// Draw palette colors
	for i, rec := range colorRects(1600, 1200, len(colors)) {
		draw.Src.Draw(newImg, rec, &image.Uniform{colors[i]}, image.Point{})
	}

	dst := image.NewRGBA(image.Rect(0, 0, 1200, 1600))
	err = graphics.Rotate(dst, newImg, &graphics.RotateOptions{math.Pi / 2})
	if err != nil {
		log.Fatal(err)
	}

	// Draw Color text
	{
		var size float64 = func() (size float64) {
			switch w := 1600 / len(colors); {
			case w >= 320:
				size = float64(w) / float64(8)
			case w >= 160:
				size = float64(w) / float64(4)
			case w >= 80:
				size = float64(w) / float64(3)
			case w >= 20:
				size = float64(w) / float64(2)
			default:
				size = float64(w)
			}

			if size > 48 {
				size = 48
			}
			return
		}()
		var dpi float64 = 72

		for i, rec := range colorRects(1600, 1200, len(colors)) {
			d := font.Drawer{
				Dst:  dst,
				Src:  reverseColor(colors[i]),
				Face: truetype.NewFace(f, &truetype.Options{Size: size, DPI: dpi}),
				Dot:  fixed.P(180, rec.Min.X+(rec.Max.X-rec.Min.X)/2+int(size)/2),
			}
			d.DrawString(m[colors[i]])
		}

		err = graphics.Rotate(newImg, dst, &graphics.RotateOptions{-math.Pi / 2})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Draw text
	{
		d := font.Drawer{
			Dst:  newImg,
			Src:  image.Black,
			Face: truetype.NewFace(f, &truetype.Options{Size: 25, DPI: 72}),
			Dot:  fixed.P(1100, 1160),
		}
		d.DrawString("https://github.com/cage1016/alfred-paletter")

		d = font.Drawer{
			Dst:  newImg,
			Src:  image.Black,
			Face: truetype.NewFace(f, &truetype.Options{Size: 25, DPI: 72}),
			Dot:  fixed.P(50, 1160),
		}
		d.DrawString(changecase.Constant(ct))
	}

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
