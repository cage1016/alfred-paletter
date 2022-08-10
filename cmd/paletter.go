/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"fmt"

	"regexp"
	"strings"

	"github.com/Baldomo/paletter"
	aw "github.com/deanishe/awgo"
	color "github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-paletter/alfred"
	"github.com/cage1016/alfred-paletter/lib"
)

const (
	hexPatternHash = "#%02x%02x%02x"
	hexPattern     = "%02x%02x%02x"
)

var (
	GaryIcon = &aw.Icon{Value: "gray.png"}

	reUrl = regexp.MustCompile(`(?m)^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
	reB64 = regexp.MustCompile(`(data:image\/[^;]+;base64,.*?)`)
)

func Hex(wf *aw.Workflow, col color.Color) string {
	if alfred.GetColorsHexWithHash(wf) {
		return fmt.Sprintf(hexPatternHash, uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
	}
	return fmt.Sprintf(hexPattern, uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
}

func unique(intSlice []color.Color) []color.Color {
	keys := make(map[string]bool)
	list := []color.Color{}
	for _, entry := range intSlice {
		if _, value := keys[entry.Hex()]; !value {
			keys[entry.Hex()] = true
			list = append(list, entry)
		}
	}
	return list
}

// paletterCmd represents the paletter command
var paletterCmd = &cobra.Command{
	Use:   "paletter",
	Short: "Extract color from an image",
	Run: func(cmd *cobra.Command, args []string) {
		q := args[0]

		var path string
		if reUrl.Match([]byte(q)) {
			path = fmt.Sprintf("%s/d.png", wf.DataDir())
			lib.Download(q, path)
			q = path
		} else if reB64.Match([]byte(q)) {
			path, err := lib.DecodeBase64(q, wf.DataDir())
			if err != nil {
				wf.NewItem(fmt.Sprintf("`%s` is %v", q, err)).Subtitle("Try a different query?").Icon(GaryIcon)
				wf.SendFeedback()
				return
			}
			q = path
		}

		img, err := paletter.OpenImage(q)
		if err != nil {
			wf.NewItem(fmt.Sprintf("`%s` is invalid", q)).Subtitle("Try a different query?").Icon(GaryIcon)
		} else {
			obs := paletter.ImageToObservation(img)
			cs, _ := paletter.CalculatePalette(obs, alfred.GetNumberOfColor(wf))
			colors := paletter.ColorsFromClusters(cs)

			hexs := make([]string, len(colors))
			for i, c := range colors {
				hexs[i] = Hex(wf, c)
			}

			uniColors := unique(colors)
			for i, c := range uniColors {
				wf.Add(1)
				path = fmt.Sprintf("%s/%d.png", wf.DataDir(), i)
				go lib.GenPng(wf, path, c)

				ni := wf.NewItem(hexs[i]).
					Subtitle(fmt.Sprintf("^ ⌥ ↩, %s", hexs[i])).
					Valid(true).
					Arg(hexs[i]).
					Icon(&aw.Icon{Value: path}).
					Var("action", "copy").
					Quicklook(q)

				if alfred.GetCopyAllSeparate(wf) {
					ni.Opt().
						Subtitle(fmt.Sprintf("↩ Copy %s separate", strings.Join(hexs, " "))).
						Arg(hexs...).
						Var("action", "copy_all_separate")
				} else {
					ni.Opt().
						Subtitle(fmt.Sprintf("Copy %s", strings.Join(hexs, " "))).
						Arg(strings.Join(hexs, " ")).
						Var("action", "copy")
				}

				ni.Ctrl().
					Subtitle(fmt.Sprintf("↩ Search Google for '%s'", hexs[i])).
					Arg(hexs[i]).
					Var("action", "search")
			}

			if len(uniColors) == 0 {
				wf.NewItem("No color found").Subtitle("Try a different query?").Icon(GaryIcon)
			}
		}

		wf.SendFeedback()
	},
}

func init() {
	rootCmd.AddCommand(paletterCmd)
}
