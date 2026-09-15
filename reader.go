package bngsdk

import "os"

type BngImporter interface {
	Reset() error
	Next([]byte) (int, error)
}

func NewBinaryImporter(file string) (BngImporter, error) {
	bin, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return OgBinReader(bin), nil
}

func NewSocketImporter(address string, port int) (BngImporter, error) {
	reader, err := NewOgUDPReader(address, port)
	if err != nil {
		return nil, err
	}

	return reader, nil
}
