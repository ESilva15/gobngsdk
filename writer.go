package bngsdk

type BngExporter interface {
	Write([]byte) (int, error)
}

func NewSocketExporter(address string, port int) (BngExporter, error) {
	writer, err := NewSocketWriter(address, port)
	if err != nil {
		return nil, err
	}

	return writer, nil
}

func NewBinaryExporter(path string) (BngExporter, error) {
	writer, err := NewOgBinWriter(path)
	if err != nil {
		return nil, err
	}

	return writer, nil
}
