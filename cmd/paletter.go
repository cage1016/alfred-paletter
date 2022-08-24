/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Baldomo/paletter"
	aw "github.com/deanishe/awgo"
	"github.com/google/uuid"
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
		CheckForUpdate()

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
		switch true {
		case lib.ReUrl.Match([]byte(q)): // http, https base image url
			q, err = lib.Download(q, wf.DataDir())
			if err != nil {
				wf.NewItem(fmt.Sprintf("`%s` download fail", q)).Subtitle("Try a different query?").Icon(GaryIcon)
				wf.SendFeedback()
				return
			}
		case lib.ReB64.Match([]byte(q)): // base64 image data
			q, err = lib.DecodeBase64(q, wf.DataDir())
			if err != nil {
				wf.NewItem(fmt.Sprintf("`%s` is %v", q, err)).Subtitle("Support png, jpeg, gif, webp, bmp, tiff. Try a different query?").Icon(GaryIcon)
				wf.SendFeedback()
				return
			}
		case lib.ClipBoardTiff.Match([]byte(q)): // alfred clipboard image
			b := filepath.Join(cfg.DbConfig.Home, cfg.DbConfig.Path, cfg.DbConfig.Name) + ".data"
			q, err = lib.Copy(filepath.Join(b, q), wf.DataDir())
			if err != nil {
				wf.NewItem(fmt.Sprintf("`%s` is %v", q, err)).Subtitle("Clipboard file not found. Try a different query?").Icon(GaryIcon)
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
			if r != "" {
				nc, _ := strconv.Atoi(strings.TrimSpace(r))
				if nc == 0 {
					nc = alfred.GetNumberOfColor(wf)
				}
				cs, _ = paletter.CalculatePalette(obs, nc)
			} else {
				cs, _ = paletter.CalculatePalette(obs, alfred.GetNumberOfColor(wf))
			}
			colors := paletter.ColorsFromClusters(cs)

			uniColors := lib.Unique(colors)
			for i, c := range uniColors {
				wf.Add(1)
				path = fmt.Sprintf("%s/cmd-paletter/%v.png", wf.DataDir(), uuid.New().String())
				go lib.GenPng(wf, path, c.Color)

				ni := wf.NewItem(c.Hex).
					Subtitle("⌘+L ^ ⌥, ↩ Export Palette").
					Arg(uniColors.HexsString()).
					Valid(true).
					UID(strconv.Itoa(i)).
					Quicklook(q).
					Copytext(c.Hex).
					Largetype(c.Hex).
					Icon(&aw.Icon{Value: path}).
					Var("action", "export")

				ni.Opt().
					Subtitle("↩ Copy single Palette").
					Valid(true).
					Arg(c.Hex).
					Var("action", "copy_palette")

				ni.Ctrl().
					Subtitle("↩ Copy multiple Palette").
					Valid(true).
					Arg(uniColors.HexsString()).
					Var("action", "copy_palette")
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
