package PartsDB

import "strings"

type Part struct {
	Row         int // used only in json
	Column      int // used only in json
	Depth       int //used only in json
	QrCode      []byte
	Count       uint
	Name        string
	Raw_barcode string
	Initialized bool
}

type RCIndex struct {
	Row int
	Col int
}

type RCDIndex struct {
	Row   int
	Col   int
	Depth int
}

type RCDCabinetIndex struct {
	Row     int
	Col     int
	Depth   int
	Cabinet int
}

type PartsDB interface {
	FindFuzzyPart(name string) []Part
	AddPart(name, raw_barcode string, amount uint) bool
	//IndexRC(row, col uint) (Part, bool)
	//IndexRCD(row, col, depth int) (Part, bool)
	IndexMeta(i interface{}) (Part, bool)
	GetAllActiveParts() []Part
	Save()
	Reload()
	//GetCabinetName() string
	GetUniqueName() string
	Show(p Part) string
}

// some type expansion for the fuzzy find module
type fuzzyparts []Part

func (fz fuzzyparts) String(i int) string {
	name := strings.ToLower(fz[i].Name)
	return name
}

func (fz fuzzyparts) Len() int {
	return len(fz)
}
