package bngsdk

import (
	"io"
	"os"
	"unsafe"
)

type GobReader struct {
	TotalRead int64
	File      *os.File
	Buf       []byte
}

func NewOgBinReader(fp string) *GobReader {
	bin, err := os.Open(fp)
	if err != nil {
		return nil
	}

	return &GobReader{
		TotalRead: 0,
		File:      bin,
		Buf:       make([]byte, unsafe.Sizeof(Outgauge{})),
	}
}

func (g *GobReader) Close() error {
	if g.File != nil {
		return g.File.Close()
	}

	return nil
}

func (g *GobReader) Reset() error {
	_, err := g.File.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	g.TotalRead = 0

	return nil
}

func (g *GobReader) Next(buffer []byte) (int, error) {
	nBytes, err := io.ReadFull(g.File, buffer)
	if err != nil {
		return 0, err
	}

	pos, _ := g.File.Seek(0, io.SeekCurrent)
	g.TotalRead = pos

	return nBytes, nil
}
