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
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-paletter/lib"
)

var keys = []string{
	"rgb",
	"hex",
	"cmyk",
	"hsl",
	"hsv",
	"hwb",
	"lab",
	"xyz",
	"hcg",
}

type Opt struct {
	Subtitle string
	Arg      []string
}

type CustomItem struct {
	Text      string
	Subtitle  string
	Arg       string
	LargeType string
	Opt       Opt
}

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy single or multiple color palette",
	Run:   runcopyCmd,
}

func runcopyCmd(cmd *cobra.Command, args []string) {
	CheckForUpdate()
	lib.RemoveAllPng(fmt.Sprintf("%s/cmd-paletter", wf.DataDir()))

	hexs := strings.Split(args[0], " ")

	m := map[string][]string{}
	colors := []color.Color{}
	for _, hex := range hexs {
		c, _ := colorful.Hex(hex)
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
	}

	// wf.NewItem
	nis := map[string]CustomItem{}
	for k, v := range m {
		nis[k] = CustomItem{
			Text: strings.Join(v, " "),
			Subtitle: func(b bool) string {
				if b {
					return fmt.Sprintf("⌘+L ⌥, ↩ Copy %s", changecase.Constant(k))
				}
				return fmt.Sprintf("⌘+L, ↩ Copy %s", changecase.Constant(k))
			}(len(v) > 1),
			Arg:       strings.Join(v, "\n"),
			LargeType: strings.Join(v, "\n"),
			Opt: Opt{
				Subtitle: fmt.Sprintf("↩ Copy %s separate", changecase.Constant(k)),
				Arg:      v,
			},
		}
	}

	wf.Add(1)
	path := fmt.Sprintf("%s/cmd-copy/%v.png", wf.DataDir(), uuid.New().String())
	go lib.GenPng(wf, path, colors...)

	wf.Configure(aw.SuppressUIDs(true))
	sort.Strings(keys)
	for _, k := range keys {
		v := nis[k]
		ni := wf.NewItem(v.Text).
			Subtitle(v.Subtitle).
			Arg(v.Arg).
			Valid(true).
			UID(k).
			Largetype(v.LargeType).
			Icon(&aw.Icon{Value: path}).
			Var("action", "copy")

		if len(v.Opt.Arg) > 1 {
			ni.Opt().
				Subtitle(v.Opt.Subtitle).
				Arg(v.Opt.Arg...).
				Valid(true).
				Var("action", "copy_all_separate")
		}
	}

	wf.SendFeedback()
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
