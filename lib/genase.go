package lib

import (
	"os"

	"github.com/ARolek/ase"
	aw "github.com/deanishe/awgo"
	"github.com/lucasb-eyer/go-colorful"
)

func GenerateASE(wf *aw.Workflow, path string, hexs []string) {
	defer wf.Done()

	ac := []ase.Color{}
	for _, hex := range hexs {
		c, _ := colorful.Hex(hex)
		r, g, b := c.RGB255()

		ac = append(ac, ase.Color{
			Name:   hex,
			Model:  "RGB",
			Values: []float32{float32(r) / 255.0, float32(g) / 255.0, float32(b) / 255.0},
			Type:   "Global",
		})
	}

	sampleAse := ase.ASE{}
	sampleAse.Colors = ac

	// Create the file to write the encoded ASE
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	// It’s idiomatic to defer a Close immediately after opening a file.
	defer f.Close()

	//	encode our ASE file
	ase.Encode(sampleAse, f)
}
