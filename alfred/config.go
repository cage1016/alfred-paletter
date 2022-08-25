package alfred

import (
	aw "github.com/deanishe/awgo"
)

const (
	NumberOfColors          = "number_of_colors"
	NumberOfClipboardImages = "number_of_clipboard_images"
)

func GetNumberOfColors(wf *aw.Workflow) int {
	return wf.Config.GetInt(NumberOfColors, 7)
}

func GetNumberOfClipboard(wf *aw.Workflow) int {
	return wf.Config.GetInt(NumberOfClipboardImages, 10)
}
