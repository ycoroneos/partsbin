package PartViewer

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/skip2/go-qrcode"

	"github.com/ycoroneos/partsbin/PartsDB"
)

func GetQRCode(partName string) []byte {
	// create a qr code that is 120x120 pixels
	bytes, err := qrcode.Encode(partName, qrcode.Medium, inchesToPixels(.4))
	if err != nil {
		panic(err)
	}
	return bytes
}

func inchesToPixels(inches float64) int {
	dpi := 300.0 // pixels per inch
	return int(inches * dpi)
}

func getPointFromPartIndex(index, rowCount, colCount, xStridePx, yStridePx, xOffsetPx, yOffsetPx int) image.Point {
	col := index / colCount
	row := index % colCount
	return image.Point{xOffsetPx + row*xStridePx, yOffsetPx + col*yStridePx}
}

func MakeQRGrid(parts []*PartsDB.Part) []image.Image {
	xStride := inchesToPixels(2.05)
	yStride := inchesToPixels(.5)
	stickerWidth := inchesToPixels(1.75)
	stickerHeight := inchesToPixels(.5)
	colCount := 4
	rowCount := 20
	xOffsetPx, yOffsetPx := inchesToPixels(.3), inchesToPixels(.515) // coordinates of first sticker in inches

	//rectangle for the paper
	r := image.Rectangle{image.Point{0, 0}, image.Point{inchesToPixels(8.5), inchesToPixels(11)}}
	paper := image.NewRGBA(r)

	pages := make([]image.Image, 0)

	partsperpage := colCount * rowCount

	min := func(a, b int) int {
		if a <= b {
			return a
		}
		return b
	}

	for i := 0; i < len(parts); i += partsperpage {
		thispageparts := parts[i:min(i+partsperpage, len(parts))]
		for i, part := range thispageparts {
			partImage, _, err := image.Decode(bytes.NewReader(part.QrCode))
			if err != nil {
				panic(err)
			}

			partLocation := getPointFromPartIndex(i, rowCount, colCount, xStride, yStride, xOffsetPx, yOffsetPx)
			partRectangle := image.Rectangle{Min: partLocation, Max: image.Point{stickerWidth + partLocation.X, stickerHeight + partLocation.Y}}

			draw.Draw(paper, partRectangle, partImage, image.Point{0, 0}, draw.Over)
		}
		pages = append(pages, paper)
		paper = image.NewRGBA(r)
	}

	return pages
}

func WriteToPNG(image image.Image, filepath string) {
	// Write to file.
	outputFile, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// Encode to png.
	err = png.Encode(outputFile, image)
	if err != nil {
		panic(err)
	}

}
