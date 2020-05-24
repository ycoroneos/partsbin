package main

import "testing"

func TestMakeQRGrid(t *testing.T) {
	// Create some parts
	name := "firstPart"
	parts := []*Part{
		&Part{
			row:    0,
			col:    0,
			name:   name,
			qrCode: getQRCode(name),
			count:  1,
		},
		&Part{
			row:    0,
			col:    1,
			name:   name,
			qrCode: getQRCode(name),
			count:  1,
		},
	}

	image := makeQRGrid(parts)
	// TODO: write image to file

}
