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

	WriteToPNG(partImage, "test_get_qrqode.png")

}
func TestMakeQRGrid(t *testing.T) {
	// Create some parts
	var parts []*PartsDB.Part
	for i := 0; i < 150; i++ {
		name := fmt.Sprintf("part%d", i)
		parts = append(parts, &PartsDB.Part{
			Name:   name,
			QrCode: GetQRCode(name),
			Count:  1,
		})
	}

	images := MakeQRGrid(parts)

	for i, image := range images {
		WriteToPNG(image, fmt.Sprintf("test_qrcode_grid_%d.png", i))
	}
}
