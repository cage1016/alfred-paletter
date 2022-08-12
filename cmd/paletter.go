/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Baldomo/paletter"
	aw "github.com/deanishe/awgo"
	"github.com/muesli/clusters"
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
		q, r := args[0], ""
		match := lib.ReNumberOfColor.FindStringSubmatch(q)
		if len(match) > 0 {
			index := lib.ReNumberOfColor.FindAllStringSubmatchIndex(q, -1)[0][0]
			q, r = strings.TrimSpace(q[:index]), q[index+1:]
		} else {
			q = strings.TrimSpace(q)
		}

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
			wf.NewItem(err.Error()).Subtitle("Support png, jpeg, gif, webp, bmp, tiff. Try a different query?").Icon(GaryIcon)
		} else {
			var cs clusters.Clusters
			obs := paletter.ImageToObservation(img)
			if len(r) > 1 {
				nc, err := strconv.Atoi(r)
				if err != nil {
					wf.NewItem(err.Error()).Subtitle("Support png, jpeg, gif, webp, bmp, tiff. Try a different query?").Icon(GaryIcon)
					wf.SendFeedback()
					return
				}
				cs, _ = paletter.CalculatePalette(obs, nc)
			} else {
				cs, _ = paletter.CalculatePalette(obs, alfred.GetNumberOfColor(wf))
			}
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
