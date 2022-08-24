/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"fmt"
	"image/color"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/hisamafahri/coco"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-paletter/fs"
	"github.com/cage1016/alfred-paletter/lib"
)

// generatorCmd represents the generator command
var generatorCmd = &cobra.Command{
	Use:   "generator",
	Short: "Images generator",
	Run:   runGeneratorCmd,
}

func runGeneratorCmd(cmd *cobra.Command, args []string) {
	t, _ := cmd.Flags().GetString("type")

	av := aw.NewArgVars()
	av.Var("far", "bar") // magic var

	switch t {
	case "svg":
		dfs := fs.NewDefaultFs(wf.DataDir())
		err := dfs.WriteFile("palette.svg", args[0], true)
		if err != nil {
			av.Var("err", err.Error())
		}
		av.Arg(fmt.Sprintf("%s/palette.svg", wf.DataDir()))
	case "image":
		arg := strings.Split(args[0], ",")
		hexs, ct := arg[0], arg[1]

		m := map[color.Color]string{}
		colors := []color.Color{}
		for _, hex := range strings.Split(hexs, " ") {
			c, _ := colorful.Hex(hex)
			colors = append(colors, c)
			r, g, b := c.RGB255()
			switch ct {
			case "rgb":
				m[c] = fmt.Sprintf("%d ,%d ,%d", r, g, b)
			case "hex":
				m[c] = hex
			case "cmyk":
				cmyk := coco.Rgb2Cmyk(float64(r), float64(g), float64(b))
				m[c] = fmt.Sprintf("%0.f ,%0.f ,%0.f ,%0.f", cmyk[0], cmyk[1], cmyk[2], cmyk[3])
			case "hsl":
				hsl := coco.Rgb2Hsl(float64(r), float64(g), float64(b))
				m[c] = fmt.Sprintf("%0.f ,%0.f ,%0.f", hsl[0], hsl[1], hsl[2])
			case "hsv":
				hsv := coco.Rgb2Hsv(float64(r), float64(g), float64(b))
				m[c] = fmt.Sprintf("%0.f ,%0.f ,%0.f", hsv[0], hsv[1], hsv[2])
			case "hwb":
				hwb := coco.Rgb2Hwb(float64(r), float64(g), float64(b))
				m[c] = fmt.Sprintf("%0.f ,%0.f ,%0.f", hwb[0], hwb[1], hwb[2])
			case "lab":
				lab := coco.Rgb2Lab(float64(r), float64(g), float64(b))
				m[c] = fmt.Sprintf("%0.f ,%0.f ,%0.f", lab[0], lab[1], lab[2])
			case "xyz":
				xyz := coco.Rgb2Xyz(float64(r), float64(g), float64(b))
				m[c] = fmt.Sprintf("%0.f ,%0.f ,%0.f", xyz[0], xyz[1], xyz[2])
			case "hcg":
				hcg := coco.Rgb2Hcg(float64(r), float64(g), float64(b))
				m[c] = fmt.Sprintf("%0.f ,%0.f ,%0.f", hcg[0], hcg[1], hcg[2])
			}
		}

		wf.Add(1)
		path := fmt.Sprintf("%s/%s.png", wf.DataDir(), ct)
		go lib.GenPng2(wf, path, colors, m, ct)
		av.Arg(path)
	case "ase":
		wf.Add(1)
		path := fmt.Sprintf("%s/palette.ase", wf.DataDir())
		go lib.GenerateASE(wf, path, strings.Split(args[0], " "))
		av.Arg(path)
	default:
	}

	if err := av.Send(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(generatorCmd)
	generatorCmd.PersistentFlags().StringP("type", "t", "", "type of generator")
}
