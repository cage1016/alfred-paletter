package lib

import (
	"fmt"
	"regexp"
	"strings"

	aw "github.com/deanishe/awgo"
	color "github.com/lucasb-eyer/go-colorful"

	"github.com/cage1016/alfred-paletter/alfred"
)

const (
	hexPatternHash = "#%02x%02x%02x"
	hexPattern     = "%02x%02x%02x"
)

var (
	ReUrl = regexp.MustCompile(`(?m)^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
	ReB64 = regexp.MustCompile(`(data:image\/[^;]+;base64,.*?)`)
)

type HexColor struct {
	Hex   string
	Color color.Color
}

type HexColors []HexColor

func (hs HexColors) HexsString() string {
	var hexs []string
	for _, h := range hs {
		hexs = append(hexs, h.Hex)
	}
	return strings.Join(hexs, " ")
}

func (hs HexColors) Hexs() []string {
	var hexs []string
	for _, h := range hs {
		hexs = append(hexs, h.Hex)
	}
	return hexs
}

func Hex(wf *aw.Workflow, col color.Color) string {
	if alfred.GetColorsHexWithHash(wf) {
		return fmt.Sprintf(hexPatternHash, uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
	}
	return fmt.Sprintf(hexPattern, uint8(col.R*255.0+0.5), uint8(col.G*255.0+0.5), uint8(col.B*255.0+0.5))
}

func Unique(intSlice []color.Color) HexColors {
	keys := make(map[string]bool)
	list := []HexColor{}
	for _, entry := range intSlice {
		if _, value := keys[entry.Hex()]; !value {
			keys[entry.Hex()] = true
			list = append(list, HexColor{Hex: entry.Hex(), Color: entry})
		}
	}
	return HexColors(list)
}
