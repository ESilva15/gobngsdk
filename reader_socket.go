package bngsdk

import (
	"errors"
)

type OgUDPReader struct {
	udpConnection *UDPTransport
}

var ErrInvalidOutgaugeData = errors.New("data is of different size than outgauge")

func NewOgUDPReader(ip string, port int) (*OgUDPReader, error) {
	conn, err := NewUDPReader(ip, port)
	if err != nil {
		return nil, err
	}

	return &OgUDPReader{
		udpConnection: conn,
	}, nil
}

func (ogr *OgUDPReader) Close() error {
	if ogr.udpConnection != nil {
		return ogr.udpConnection.Close()
	}

	return nil
}

func (ogr *OgUDPReader) Reset() error {
	// NOTE: what to implement here?
	return nil
}

func (ogr *OgUDPReader) Next(buffer []byte) (int, error) {
	nBytes, err := ogr.udpConnection.Read(buffer)
	if err != nil {
		return 0, err
	}

	// Check if enough data was received to fill our struct
	if int64(nBytes) != OutgaugeSize {
		return 0, ErrInvalidOutgaugeData
	}

	return nBytes, nil
}

func (ogr *OgUDPReader) GetTotalRead() int64 {
	return ogr.udpConnection.GetTotalBytes()
}
