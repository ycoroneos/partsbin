package PartViewer

import (
	"bytes"
	"fmt"
	"image"
	"testing"

	"github.com/ycoroneos/partsbin/PartsDB"
)

func TestGetQRCode(t *testing.T) {
	// Create some parts
	codeBytes := GetQRCode("partName")
	partImage, _, err := image.Decode(bytes.NewReader(codeBytes))
	if err != nil {
		panic(err)
	}

	WriteToPostScript(partImage, "test_get_qrqode.png")

}
func TestMakeQRGrid(t *testing.T) {
	// Create some parts
	var parts []*PartsDB.Part
	for i := 0; i < 50; i++ {
		name := fmt.Sprintf("part%d", i)
		parts = append(parts, &PartsDB.Part{
			Name:   name,
			QrCode: GetQRCode(name),
			Count:  1,
		})
	}

	image := MakeQRGrid(parts)
	WriteToPostScript(image, "test_qrqode_grid.tiff")
}
