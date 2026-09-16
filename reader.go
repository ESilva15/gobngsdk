package bngsdk

type BngImporter interface {
	Reset() error
	Next([]byte) (int, error)
	GetTotalRead() int64
	Close() error
}

func NewBinaryImporter(file string) (BngImporter, error) {
	return NewOgBinReader(file), nil
}

func NewSocketImporter(address string, port int) (BngImporter, error) {
	reader, err := NewOgUDPReader(address, port)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
