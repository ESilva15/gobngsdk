package bngsdk

import (
	"encoding/binary"
	"os"
)

type OgBinWriter struct {
	file       *os.File
	totalBytes int64
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

func (ogw *OgBinWriter) Close() error {
	if ogw.file != nil {
		return ogw.file.Close()
	}

	return nil
}

func (ogw *OgBinWriter) Write(data []byte) (int, error) {
	err := binary.Write(ogw.file, binary.LittleEndian, data)
	if err != nil {
		return 0, err
	}

	nBytes := len(data)
	ogw.totalBytes += int64(nBytes)

	return nBytes, nil
}

func (ogw *OgBinWriter) GetTotalWritten() int64 {
	return ogw.totalBytes
}
