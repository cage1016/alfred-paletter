/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-paletter/lib"
	"github.com/cage1016/alfred-paletter/sqlite"
)

// aCmd represents the a command
var aCmd = &cobra.Command{
	Use:   "paletterfi",
	Short: "Find Clipboard Image and open with Paletter",
	Run: func(cmd *cobra.Command, args []string) {
		CheckForUpdate()

		var q, r string
		if len(args) > 0 {
			q, r = args[0], ""
			match := lib.ReNumberOfColor.FindStringSubmatch(q)
			if len(match) > 0 {
				index := lib.ReNumberOfColor.FindAllStringSubmatchIndex(q, -1)[0][0]
				r = q[index+1:]
			}
		}

		repo := sqlite.New(db)
		var items []*lib.ClipBoard
		var err error
		if r != "" {
			nc, _ := strconv.Atoi(strings.TrimSpace(r))
			if nc == 0 {
				nc = 10
			}
			items, err = repo.List(context.Background(), nc)
		} else {
			items, err = repo.List(context.Background(), 10)
		}

		if err != nil {
			wf.NewItem(err.Error()).Subtitle("Try a different query?").Icon(GaryIcon)
		} else {
			b := filepath.Join(cfg.DbConfig.Home, cfg.DbConfig.Path, cfg.DbConfig.Name) + ".data"
			for _, item := range items {
				wf.NewItem(item.Item).
					Subtitle(item.DateTime).
					Valid(true).
					Arg(item.DataHash).
					Quicklook(filepath.Join(b, item.DataHash)).
					Icon(&aw.Icon{Value: filepath.Join(b, item.DataHash)})
			}
		}

		wf.SendFeedback()
	},
}

func init() {
	rootCmd.AddCommand(aCmd)
}
