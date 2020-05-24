package main

import (
	"fmt"
	"github.com/skip2/go-qrcode"
)

type Part struct {
	row    uint
	column uint
	qrCode []byte
	count  uint
	name   string
}

func getQRCode(partName string) []byte {
	// create a qr code that is 120x120 pixels
	bytes, err := qrcode.Encode(partName, qrcode.Medium, 120)
	if err != nil {
		panic(err)
	}
	return bytes
}

func inchesToPixels(inches float) int {
	dpi := 300 // pixels per inch
	return inches * dpi
}

func makeQRGrid(parts []*Part) Image {
	xStride := 2         //inches
	yStride := .5        // inches
	stickerWidth := 1.75 //inches
	stickerHeight := .5  //inches
	colCount := 4
	rowCount := 20
	startingPointX, startingPointY := .3, .515 // coordinates of first sticker in inches

	for _, part := range parts {
		partImage, _, err := image.Decode(part.qrCode)
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	fmt.Println("vim-go")
	_ = qrcode.WriteFile("https://example.org", qrcode.Medium, 0, "qr.png")

}
