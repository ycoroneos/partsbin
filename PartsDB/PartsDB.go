package PartsDB

type Part struct {
	Row         uint //deprecated
	Column      uint //deprecated
	QrCode      []byte
	Count       uint
	Name        string
	Raw_barcode string
	Initialized bool
}

type PartsDB interface {
	FindFuzzyPart(name string) []Part
	AddPart(name, raw_barcode string, amount uint) bool
	IndexRC(row, col uint) (Part, bool)
	GetAllActiveParts() []Part
	Save()
	Reload()
	GetCabinetName() string
	Show(p Part) string
}
