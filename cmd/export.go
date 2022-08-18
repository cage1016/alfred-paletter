/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"fmt"
	"image/color"
	"sort"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/google/uuid"
	"github.com/hisamafahri/coco"
	changecase "github.com/ku/go-change-case"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut/palette"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-paletter/lib"
	template "github.com/cage1016/alfred-paletter/templates"
)

type Result struct {
	Name string
	Hex  string
	RGB  []float64
	CMYK []float64
	HSV  []float64
	HSL  []float64
	LAB  []float64
	HCG  []float64
	HWB  []float64
	XYZ  []float64
}

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Run:   runExportCmd,
}

func runExportCmd(cmd *cobra.Command, args []string) {
	CheckForUpdate()

	hexs := strings.Split(args[0], " ")
	results := []Result{}

	m := map[string][]string{}
	m2 := map[color.Color]string{}
	colors := []color.Color{}
	for _, hex := range hexs {
		c, _ := colorful.Hex(hex)
		m2[c] = hex
		colors = append(colors, c)
		r, g, b := c.RGB255()

		cmyk := coco.Rgb2Cmyk(float64(r), float64(g), float64(b))
		hsl := coco.Rgb2Hsl(float64(r), float64(g), float64(b))
		hsv := coco.Rgb2Hsv(float64(r), float64(g), float64(b))
		hwb := coco.Rgb2Hwb(float64(r), float64(g), float64(b))
		lab := coco.Rgb2Lab(float64(r), float64(g), float64(b))
		xyz := coco.Rgb2Xyz(float64(r), float64(g), float64(b))
		hcg := coco.Rgb2Hcg(float64(r), float64(g), float64(b))

		m["rgb"] = append(m["rgb"], fmt.Sprintf("%d, %d, %d", r, g, b))
		m["hex"] = append(m["hex"], hex)
		m["cmyk"] = append(m["cmyk"], fmt.Sprintf("%0.f, %0.f, %0.f, %0.f", cmyk[0], cmyk[1], cmyk[2], cmyk[3]))
		m["hsl"] = append(m["hsl"], fmt.Sprintf("%0.f, %0.f, %0.f", hsl[0], hsl[1], hsl[2]))
		m["hsv"] = append(m["hsv"], fmt.Sprintf("%0.f, %0.f, %0.f", hsv[0], hsv[1], hsv[2]))
		m["hwb"] = append(m["hwb"], fmt.Sprintf("%0.f, %0.f, %0.f", hwb[0], hwb[1], hwb[2]))
		m["lab"] = append(m["lab"], fmt.Sprintf("%0.f, %0.f, %0.f", lab[0], lab[1], lab[2]))
		m["xyz"] = append(m["xyz"], fmt.Sprintf("%0.f, %0.f, %0.f", xyz[0], xyz[1], xyz[2]))
		m["hcg"] = append(m["hcg"], fmt.Sprintf("%0.f, %0.f, %0.f", hcg[0], hcg[1], hcg[2]))

		cc, _ := palette.Wikipedia.Name(c)
		results = append(results, Result{
			Name: cc[0].Name,
			RGB:  []float64{float64(r), float64(g), float64(b)},
			Hex:  hex,
			CMYK: []float64{cmyk[0], cmyk[1], cmyk[2]},
			HSL:  []float64{hsl[0], hsl[1], hsl[2]},
			HSV:  []float64{hsv[0], hsv[1], hsv[2]},
			HWB:  []float64{hwb[0], hwb[1], hwb[2]},
			LAB:  []float64{lab[0], lab[1], lab[2]},
			XYZ:  []float64{xyz[0], xyz[1], xyz[2]},
			HCG:  []float64{hcg[0], hcg[1], hcg[2]},
		})
	}

	te := template.NewEngine()
	// Code
	code, _ := te.Execute("code", results)
	wf.NewItem(code).
		Subtitle("⌘+L, ↩ Copy Code").
		Arg(code).
		Valid(true).
		Largetype(code).
		Icon(WebCodingIcon).
		Var("action", "copy")

	// CSS
	css, _ := te.Execute("css", results)
	wf.NewItem(css).
		Subtitle("⌘+L, ↩ Copy CSS").
		Arg(css).
		Valid(true).
		Largetype(css).
		Icon(CssCodeIcon).
		Var("action", "copy")

	// SVG
	svg, _ := te.Execute("svg", results)
	wf.NewItem(svg).
		Subtitle("⌘+L, ↩ Browse SVG file").
		Arg(svg).
		Valid(true).
		Largetype(svg).
		Icon(SvgIcon).
		Var("action", "svc")

	// Images
	wf.Add(1)
	path := fmt.Sprintf("%s/cmd-export/%v.png", wf.DataDir(), uuid.New().String())
	go lib.GenPng2(wf, path, colors, m2, "hex")

	wf.Configure(aw.SuppressUIDs(true))
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		wf.NewItem(strings.Join(v, " ")).
			Subtitle(fmt.Sprintf("⌘+L, ↩ Browse %s file", changecase.Constant(k))).
			Arg(fmt.Sprintf("%s,%s", args[0], k)).
			Valid(true).
			Largetype(strings.Join(v, "\n")).
			Icon(&aw.Icon{Value: path}).
			Var("action", "image")
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
