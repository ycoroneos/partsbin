package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ycoroneos/partsbin/PartViewer"
	"github.com/ycoroneos/partsbin/PartsDB"
)

func usage() {
	fmt.Printf("./partsdb <scan | print | find | convert | admin>\n")
}

func scannerapp(pdb PartsDB.PartsDB) {
	// loop which reads barcode inputs and decodes them
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("feed me barcode: ")
		rawBarcode, _ := scanner.ReadString('\n')
		fmt.Printf("raw barcode %+q\n", rawBarcode)
		barcode := strings.SplitAfter(rawBarcode, "-ND")[0]
		barcode = strings.SplitAfter(barcode, "~P")[1]
		//barcode = strings.SplitAfter(barcode, "[)>\x1b[20~06\x1b[19~P")[1]
		fmt.Printf("decoded %+q\n", barcode)
		success := pdb.AddPart(barcode, rawBarcode, 1)
		if success {
			pdb.Save()
			indexed_part := pdb.FindFuzzyPart(barcode)[0]
			fmt.Println(pdb.Show(indexed_part))
		} else {
			fmt.Printf("out of space!\n")
		}
	}

}

func printerapp(pdb PartsDB.PartsDB) {
	images := PartViewer.MakeQRGrid(pdb.GetAllActiveParts())
	for i, image := range images {
		PartViewer.WriteToPNG(image, fmt.Sprintf("%s_grid_%d.png", pdb.GetUniqueName(), i))
	}
}

func finderapp(pdb PartsDB.PartsDB) {
	// loop which reads barcode inputs and decodes them
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("tell me your wish, stalker: ")
		wish, _ := scanner.ReadString('\n')
		wish = strings.TrimSpace(wish)
		parts := pdb.FindFuzzyPart(wish)
		for _, p := range parts {
			fmt.Println(pdb.Show(p))
		}
	}

}

func converterapp(args []string) {

	usage := fmt.Sprintf("usage: <version_from> <version_to> <file_from> <file_to>\n")
	if len(args) != 4 {
		fmt.Printf("you said : %+v\n", args)
		panic(usage)
	}

	version_from := args[0]
	version_to := args[1]
	file_from := args[2]
	file_to := args[3]

	if strings.Compare(version_from, "v1") == 0 && strings.Compare(version_to, "v2") == 0 {
		PartsDB.ConvertFromPartsDBV1(file_from, file_to, 2)
	}

	fmt.Printf("finished converting, please check %s\n", file_to)

}

// one-off things you need to do with code
func adminapp(args []string) {

	usage := fmt.Sprintf("usage: <filename>\n")
	if len(args) != 1 {
		fmt.Printf("you said : %+v\n", args)
		panic(usage)
	}

	pdb := PartsDB.OpenPartsDBV2(args[0])
	pdb.AddCabinet("shelf1", 8, 8, 2)
	pdb.Save()

	fmt.Printf("finished adding cabinet, please check %s\n", args[0])

}

func main() {

	//pdb := PartsDB.MakeOrOpenPartsDB("parts_v1.json", "shelf0", 8, 8)
	// we've moved to PDBV2
	pdb := PartsDB.OpenPartsDBV2("parts_v2.json")

	cmdargs := os.Args[1:]

	fmt.Printf("%+v", cmdargs)

	if len(cmdargs) == 0 {
		usage()
	} else if strings.Compare(cmdargs[0], "scan") == 0 {
		fmt.Printf("launched partsdb in scanner mode\n")
		scannerapp(pdb)
	} else if strings.Compare(cmdargs[0], "print") == 0 {
		fmt.Printf("launched partsdb in printer mode\n")
		printerapp(pdb)
	} else if strings.Compare(cmdargs[0], "find") == 0 {
		fmt.Printf("launched partsdb in finder mode\n")
		finderapp(pdb)
	} else if strings.Compare(cmdargs[0], "convert") == 0 {
		fmt.Printf("converting databases\n")
		converterapp(cmdargs[1:])
	} else if strings.Compare(cmdargs[0], "admin") == 0 {
		fmt.Printf("administrating\n")
		adminapp(cmdargs[1:])
	}

}
