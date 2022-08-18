/*
Copyright © 2022 KAI CHU CHUNG <cage.chung@gmail.com> (https://kaichu.io)

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cage1016/alfred-paletter/lib"
)

// clearPngCmd represents the clearPng command
var clearPngCmd = &cobra.Command{
	Use:   "clearPng",
	Short: "Remove temporary png files",
	Run: func(cmd *cobra.Command, args []string) {
		lib.RemoveAllPng(fmt.Sprintf("%s/%s", wf.DataDir(), args[0]))
	},
}

func init() {
	rootCmd.AddCommand(clearPngCmd)
}
