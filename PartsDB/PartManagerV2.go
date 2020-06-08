package PartsDB

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/ycoroneos/partsbin/PartHelper"
)

type Cabinet struct {
	Nrows       int
	Ncols       int
	Ndepth      int
	Parts       [][][]Part // [row][col][depth]
	CabinetName string
}

type partmanagerV2 struct {
	Cabinets []Cabinet
	Savefile string
}

func (pm *partmanagerV2) FindFuzzyPart(name string) []Part {
	activeparts := pm.GetAllActiveParts()
	results := fuzzy.FindFrom(strings.ToLower(name), fuzzyparts(activeparts))
	partresults := make([]Part, 0)
	for _, i := range results {
		partresults = append(partresults, activeparts[i.Index])
	}
	return partresults
}

func (pm *partmanagerV2) AddPart(name, raw_barcode string, amount uint) bool {

	// scan through all bins in all cabinets. If we enounter the name then also
	// check that we have scanned a unique barcode. If so, add to the count.
	// Otherwise this is a new item so we put it in the first blank

	for _, cab := range pm.Cabinets {

		for row := 0; row < cab.Nrows; row++ {
			for col := 0; col < cab.Ncols; col++ {
				for depth := 0; depth < cab.Ndepth; depth++ {
					if !cab.Parts[row][col][depth].Initialized {

						cab.Parts[row][col][depth] = Part{
							QrCode:      PartHelper.GetQRCode(name),
							Count:       amount,
							Name:        name,
							Raw_barcode: raw_barcode,
							Initialized: true,
						}
						return true
					} else {
						if strings.Compare(cab.Parts[row][col][depth].Name, name) == 0 {
							if strings.Compare(cab.Parts[row][col][depth].Raw_barcode, raw_barcode) != 0 {
								cab.Parts[row][col][depth].Count += amount
							}
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (pm *partmanagerV2) nativeindex(row, col, depth, cab int) (Part, bool) {
	if cab < len(pm.Cabinets) && row < pm.Cabinets[cab].Nrows && col < pm.Cabinets[cab].Ncols && depth < pm.Cabinets[cab].Ndepth && pm.Cabinets[cab].Parts[row][col][depth].Initialized {
		return pm.Cabinets[cab].Parts[row][col][depth], true
	}

	return Part{}, false

}

func (pm *partmanagerV2) IndexMeta(i interface{}) (Part, bool) {
	switch v := i.(type) {
	case RCDCabinetIndex:
		return pm.nativeindex(v.Row, v.Col, v.Depth, v.Cabinet)
	default:
		fmt.Printf("unsupported indexing method %+v\n", v)
		return Part{}, false
	}
}

func (pm *partmanagerV2) GetAllActiveParts() []Part {
	parts := make([]Part, 0)
	for cnum, cab := range pm.Cabinets {
		for row := 0; row < cab.Nrows; row++ {
			for col := 0; col < cab.Ncols; col++ {
				for depth := 0; depth < cab.Ndepth; depth++ {
					index := RCDCabinetIndex{
						Row:     row,
						Col:     col,
						Depth:   depth,
						Cabinet: cnum,
					}
					part, valid := pm.IndexMeta(index)
					if valid {
						parts = append(parts, part)
					}
				}
			}
		}
	}
	return parts
}

func (pm *partmanagerV2) Save() {
	data, err := json.Marshal(pm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(pm.Savefile, data, 0644)
	if err != nil {
		panic(err)
	}
}

func (pm *partmanagerV2) Reload() {
	data, err := ioutil.ReadFile(pm.Savefile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, pm)
	if err != nil {
		panic(err)
	}

}

func (pm *partmanagerV2) GetUniqueName() string {
	uniqueName := ""
	for _, cab := range pm.Cabinets {
		uniqueName += fmt.Sprintf("_%s_", cab.CabinetName)
	}
	return uniqueName
}

func (pm *partmanagerV2) Show(p Part) string {
	for cnum, cab := range pm.Cabinets {
		for row := 0; row < cab.Nrows; row++ {
			for col := 0; col < cab.Ncols; col++ {
				for depth := 0; depth < cab.Ndepth; depth++ {
					index := RCDCabinetIndex{
						Row:     row,
						Col:     col,
						Depth:   depth,
						Cabinet: cnum,
					}
					ipart, success := pm.IndexMeta(index)
					if success && ipart.Initialized && (strings.Compare(ipart.Name, p.Name) == 0) {
						return fmt.Sprintf("row %d, col %d, depth %d, count %d, pn %s\n", row, col, depth, p.Count, p.Name)
					}
				}
			}
		}
	}
	return fmt.Sprintf("not found\n")
}

func OpenPartsDBV2(filename string) PartsDB {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		panic(err)
	}

	if info.IsDir() {
		panic(fmt.Sprintf("cannot open directory %s", filename))
	}

	pm := &partmanagerV2{
		Savefile: filename,
	}

	pm.Reload()

	return PartsDB(pm)
}

func MakePartsDBV2(filename string, cabinets []Cabinet) {
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		panic(fmt.Sprintf("file %s already exists", filename))
	}

	for _, cab := range cabinets {
		cab.Parts = make([][][]Part, cab.Nrows)
		for row := range cab.Parts {
			cab.Parts[row] = make([][]Part, cab.Nrows)
			for col := range cab.Parts[row] {
				cab.Parts[row][col] = make([]Part, cab.Ndepth)
				for depth := range cab.Parts[row][col] {
					cab.Parts[row][col][depth] = Part{
						Row:         row,
						Column:      col,
						Depth:       depth,
						Initialized: false,
					}
				}
			}
		}

	}

	pm := &partmanagerV2{
		Cabinets: cabinets,
		Savefile: filename,
	}

	pdb := PartsDB(pm)
	pdb.Save()
}

func ConvertFromPartsDBV1(fromFilename, toFilename string, depth int) {

	// open pdb v1
	info, err := os.Stat(fromFilename)
	if os.IsNotExist(err) {
		panic(err)
	}

	if info.IsDir() {
		panic("trying to open a directory")
	}

	pmv1 := &partmanager{
		Savefile: fromFilename,
	}

	pmv1.Reload()

	// pdbv1 only has a single cabinet
	cabinets := []Cabinet{
		Cabinet{
			Nrows:       int(pmv1.Nrows),
			Ncols:       int(pmv1.Ncols),
			Ndepth:      depth,
			CabinetName: pmv1.CabinetName,
		},
	}

	// make the parts array
	cabinets[0].Parts = make([][][]Part, cabinets[0].Nrows)
	for row := range cabinets[0].Parts {
		cabinets[0].Parts[row] = make([][]Part, cabinets[0].Nrows)
		for col := range cabinets[0].Parts[row] {
			cabinets[0].Parts[row][col] = make([]Part, cabinets[0].Ndepth)
			for depth := range cabinets[0].Parts[row][col] {
				cabinets[0].Parts[row][col][depth] = Part{
					Row:         row,
					Column:      col,
					Depth:       depth,
					Initialized: false,
				}
			}
		}
	}

	// all parts from pdb v1 will be at depth=0 in pdb v2
	for row := uint(0); row < pmv1.Nrows; row++ {
		for col := uint(0); col < pmv1.Ncols; col++ {
			depth := 0
			cabinets[0].Parts[row][col][depth] = pmv1.Parts[row][col]
		}
	}

	// finalize
	MakePartsDBV2(toFilename, cabinets)
}
