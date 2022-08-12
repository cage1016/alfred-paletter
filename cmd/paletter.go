/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"fmt"

	"github.com/Baldomo/paletter"
	aw "github.com/deanishe/awgo"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-paletter/alfred"
	"github.com/cage1016/alfred-paletter/lib"
)

var (
	GaryIcon = &aw.Icon{Value: "gray.png"}
)

// paletterCmd represents the paletter command
var paletterCmd = &cobra.Command{
	Use:   "paletter",
	Short: "Extract color from an image",
	Run: func(cmd *cobra.Command, args []string) {
		q := args[0]

		var path string
		var err error
		if lib.ReUrl.Match([]byte(q)) {
			// http, https base image url
			q, err = lib.Download(q, wf.DataDir())
			if err != nil {
				wf.NewItem(fmt.Sprintf("`%s` download fail", q)).Subtitle("Try a different query?").Icon(GaryIcon)
				wf.SendFeedback()
				return
			}
		} else if lib.ReB64.Match([]byte(q)) {
			// base64 image data
			q, err = lib.DecodeBase64(q, wf.DataDir())
			if err != nil {
				wf.NewItem(fmt.Sprintf("`%s` is %v", q, err)).Subtitle("Support png, jpeg, gif, webp, bmp, tiff. Try a different query?").Icon(GaryIcon)
				wf.SendFeedback()
				return
			}
		}

		// local file
		img, err := paletter.OpenImage(q)
		if err != nil {
			wf.NewItem(fmt.Sprintf("%s", err.Error())).Subtitle("Support png, jpeg, gif, webp, bmp, tiff. Try a different query?").Icon(GaryIcon)
		} else {
			obs := paletter.ImageToObservation(img)
			cs, _ := paletter.CalculatePalette(obs, alfred.GetNumberOfColor(wf))
			colors := paletter.ColorsFromClusters(cs)

			uniColors := lib.Unique(colors)
			for i, c := range uniColors {
				wf.Add(1)
				path = fmt.Sprintf("%s/%d.png", wf.DataDir(), i)
				go lib.GenPng(wf, path, c.Color)

				ni := wf.NewItem(c.Hex).
					Subtitle(fmt.Sprintf("^ ⌥, ↩ Copy %s", c.Hex)).
					Valid(true).
					Arg(c.Hex).
					Icon(&aw.Icon{Value: path}).
					Var("action", "copy").
					Quicklook(q)

				if alfred.GetCopyAllSeparate(wf) {
					ni.Opt().
						Subtitle(fmt.Sprintf("↩ Copy %s separate", uniColors.HexsString())).
						Arg(uniColors.Hexs()...).
						Var("action", "copy_all_separate")
				} else {
					ni.Opt().
						Subtitle(fmt.Sprintf("Copy %s", uniColors.HexsString())).
						Arg(uniColors.HexsString()).
						Var("action", "copy")
				}

				ni.Ctrl().
					Subtitle(fmt.Sprintf("↩ Search Google for '%s'", c.Hex)).
					Arg(c.Hex).
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
