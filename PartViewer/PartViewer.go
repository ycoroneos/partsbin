package PartViewer

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/skip2/go-qrcode"
	"golang.org/x/image/tiff"

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
	dpi := 72.0 // pixels per inch
	return int(inches * dpi)
}

func getPointFromPartIndex(index, rowCount, colCount, xStridePx, yStridePx, xOffsetPx, yOffsetPx int) image.Point {
	col := index / colCount
	row := index % colCount
	return image.Point{xOffsetPx + row*xStridePx, yOffsetPx + col*yStridePx}
}

func MakeQRGrid(parts []*PartsDB.Part) image.Image {
	xStride := inchesToPixels(2.05)
	yStride := inchesToPixels(.5)
	stickerWidth := inchesToPixels(1.75)
	stickerHeight := inchesToPixels(.5)
	colCount := 4
	rowCount := 20
	xOffsetPx, yOffsetPx := inchesToPixels(.3), inchesToPixels(.515) // coordinates of first sticker in inches

	//rectangle for the paper
	r := image.Rectangle{image.Point{0, 0}, image.Point{inchesToPixels(8.5), inchesToPixels(11)}}
	fmt.Printf("paper rect %+v", r)
	paper := image.NewRGBA(r)

	// distance between columns too little
	// rotated 90

	for i, part := range parts {
		partImage, _, err := image.Decode(bytes.NewReader(part.QrCode))
		if err != nil {
			panic(err)
		}
		//WriteToPostScript(partImage, fmt.Sprintf("%s.png", part.Name))

		// create a standard rectangle for the parts
		// partRectangle := image.Rectangle{Min: image.Point{90, 304}, Max: image.Point{stickerWidth + 90, stickerHeight + 304}}

		//// make sure QR code image bounds are within rectangle we define
		//if !partImage.Bounds().In(partRectangle) {
		//	panic("qrcode too big for slot")
		//}

		partLocation := getPointFromPartIndex(i, rowCount, colCount, xStride, yStride, xOffsetPx, yOffsetPx)
		partRectangle := image.Rectangle{Min: partLocation, Max: image.Point{stickerWidth + partLocation.X, stickerHeight + partLocation.Y}}
		// translatedPartRectangle := image.Rectangle{Min: partLocation, Max: partLocation.Add(image.Point{stickerWidth, stickerHeight})}
		//shiftedPartRectange := image.Rectangle{Min:}
		fmt.Printf("drawing part %d, rectangle %+v, location: %+v\n", i, partRectangle, partLocation)

		// draw.Draw(paper, translatedPartRectangle, partImage, partLocation, draw.Over)
		draw.Draw(paper, partRectangle, partImage, image.Point{0, 0}, draw.Over)
		//draw.Draw(paper, translatedPartRectangle, partImage, image.Point{0, 0}, draw.Over)
		//draw.Draw(paper, partRectangle, partImage, image.Point{-1 * partLocation.Y, -1 * partLocation.X}, draw.Over)
		//draw.Draw(paper, image.Rectangle{image.Point{-10, -100}, image.Point{-10 - stickerWidth, -100 - stickerHeight}}, partImage, image.Point{-10, -100}, draw.Over)
		//WriteToPostScript(paper, fmt.Sprintf("paper_%d.png", i))
	}

	return paper
}

func WriteToPostScript(image image.Image, filepath string) {
	// Write to file.
	outputFile, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	// encode to tiff
	err = tiff.Encode(outputFile, image, &tiff.Options{
		Compression: tiff.Uncompressed,
		Predictor:   false,
	})
	if err != nil {
		panic(err)
	}

	//// Encode to png.
	//err = png.Encode(outputFile, image)
	//if err != nil {
	//	panic(err)
	//}

}
