package PartHelper

import "github.com/skip2/go-qrcode"

func InchesToPixels(inches float64) int {
	dpi := 300.0 // pixels per inch
	return int(inches * dpi)
}

func GetQRCode(partName string) []byte {
	// create a qr code that is 120x120 pixels
	bytes, err := qrcode.Encode(partName, qrcode.Medium, InchesToPixels(.4))
	if err != nil {
		panic(err)
	}
	return bytes
}
