package PartsDB

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Part struct {
	Row         uint
	Column      uint
	QrCode      []byte
	Count       uint
	Name        string
	Raw_barcode string
	Initialized bool
}

type partmanager struct {
	Nrows    uint
	Ncols    uint
	Parts    [][]Part
	Savefile string
}

type PartsDB interface {
	FindFuzzyPart(name string) []Part
	AddPart(name, raw_barcode string, amount uint)
	IndexRC(row, col uint) (Part, bool)
	Save()
	Reload()
}

func (pm *partmanager) FindFuzzyPart(name string) []Part {
	return []Part{Part{}}
}

func (pm *partmanager) AddPart(name, raw_barcode string, amount uint) {

	// scan through all columns. If we enounter the name then also
	// check that we have scanned a unique barcode. If so, add to the count.
	// Otherwise this is a new item so we put it in the first blank

	for row := uint(0); row < pm.Nrows; row++ {
		for col := uint(0); col < pm.Ncols; col++ {
			if !pm.Parts[row][col].Initialized {

				pm.Parts[row][col] = Part{
					Row:         row,
					Column:      col,
					Count:       amount,
					Name:        name,
					Raw_barcode: raw_barcode,
					Initialized: true,
				}
				return
			} else {
				if strings.Compare(pm.Parts[row][col].Name, name) == 0 {
					if strings.Compare(pm.Parts[row][col].Raw_barcode, raw_barcode) != 0 {
						pm.Parts[row][col].Count += amount
					}
					return
				}
			}

		}
	}
}

func (pm *partmanager) IndexRC(row, col uint) (Part, bool) {
	if row < pm.Nrows && col < pm.Ncols && pm.Parts[row][col].Initialized {
		return pm.Parts[row][col], true
	} else {
		return Part{}, false
	}
}

func (pm *partmanager) Save() {
	data, err := json.Marshal(pm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(pm.Savefile, data, 0644)
	if err != nil {
		panic(err)
	}
}

func (pm *partmanager) Reload() {
	data, err := ioutil.ReadFile(pm.Savefile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, pm)
	if err != nil {
		panic(err)
	}

}

func MakeOrOpenPartsDB(filename string, ncols, nrows uint) PartsDB {
	newfile := false
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		newfile = true
	}

	if !newfile && info.IsDir() {
		panic("trying to open a directory")
	}

	pm := &partmanager{
		Nrows:    nrows,
		Ncols:    ncols,
		Savefile: filename,
	}

	if newfile {
		// make the 2d parts array
		pm.Parts = make([][]Part, nrows)
		for i := range pm.Parts {
			pm.Parts[i] = make([]Part, ncols)
		}
		fmt.Printf("made new file %v\n", filename)
	} else {
		pm.Reload()
		fmt.Printf("loaded file %v : %+v\n", filename, pm)
	}

	return PartsDB(pm)

}
