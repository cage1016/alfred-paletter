package alfred

import (
	aw "github.com/deanishe/awgo"
)

const (
	NumberOfColor     = "number_of_color"
	ColorsHexWithHash = "colors_hex_with_hash"
	CopyAllSeparate   = "copy_all_separate"
)

func GetNumberOfColor(wf *aw.Workflow) int {
	return wf.Config.GetInt(NumberOfColor, 3)
}

func GetColorsHexWithHash(wf *aw.Workflow) bool {
	return wf.Config.GetBool(ColorsHexWithHash)
}

func GetCopyAllSeparate(wf *aw.Workflow) bool {
	return wf.Config.GetBool(CopyAllSeparate)
}
