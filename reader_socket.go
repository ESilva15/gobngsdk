package bngsdk

import (
	"fmt"
	"log/slog"
)

type OgUDPReader struct {
	udpConnection *UDPTransport
}

func NewOgUDPReader(ip string, port int) (*OgUDPReader, error) {
	conn, err := NewUDPReader(ip, port)
	if err != nil {
		return nil, err
	}

	return &OgUDPReader{
		udpConnection: conn,
	}, nil
}

func (ogr *OgUDPReader) Reset() error {
	// NOTE: what to implement here?
	return nil
}

func (ogr *OgUDPReader) Next(buffer []byte) (int, error) {
	slog.Debug("Reading")

	nBytes, err := ogr.udpConnection.Read(buffer)
	if err != nil {
		return 0, err
	}

	// Check if enough data was received to fill our struct
	if nBytes < outgaugeSize {
		return 0, fmt.Errorf("received data is smaller than outgauge size")
	}

	return 0, nil
}
