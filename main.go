package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ycoroneos/partsbin/PartsDB"
)

func main() {

	pdb := PartsDB.MakeOrOpenPartsDB("parts_v1.json", 8, 8)

	// loop which reads barcode inputs and decodes them
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("feed me barcode: ")
		rawBarcode, _ := scanner.ReadString('\n')
		barcode := strings.SplitAfter(rawBarcode, "-ND")[0]
		barcode = strings.SplitAfter(barcode, "[)>\x1b[20~06\x1b[19~P")[1]
		fmt.Printf("decoded %+q\n", barcode)
		pdb.AddPart(barcode, rawBarcode, 1)
		pdb.Save()
	}

}
