package bngsdk

import (
	"encoding/binary"
	"os"
)

type OgBinWriter struct {
	file *os.File
}

func NewOgBinWriter(path string) (*OgBinWriter, error) {
	outputFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &OgBinWriter{
		file: outputFile,
	}, nil
}

func (ogw *OgBinWriter) Write(data []byte) (int, error) {
	err := binary.Write(ogw.file, binary.LittleEndian, data)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}
