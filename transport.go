package bngsdk

import (
	"errors"
	"log/slog"
	"net"
	"time"
	"unsafe"
)

// ErrReaderDisabled = errors.New("reader is currently turned off")
var (
	ErrNoData   = errors.New("no new data available")
	readTimeout = (time.Second / 60) * 5 // N missed frames at 60fps
)

type UDPTransport struct {
	address    *net.UDPAddr
	connection *net.UDPConn
	dataChan   chan []byte
}

func NewUDPReader(ip string, port int) (*UDPTransport, error) {
	// Define the IP address and port to listen on
	addr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	udpT := UDPTransport{
		address:    addr,
		connection: conn,
		dataChan:   make(chan []byte, 1),
	}

	go udpT.udpSink()

	return &udpT, nil
}

func NewUDPWriter(ip string, port int) (*UDPTransport, error) {
	// Define the IP address and port to listen on
	addr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &UDPTransport{
		address:    addr,
		connection: conn,
	}, nil
}

// udpSink is a loop that will consume the most recent packets on the port
// instead of letting them pile up
func (ut *UDPTransport) udpSink() {
	defer close(ut.dataChan)

	for {
		ut.connection.SetReadDeadline(time.Now().Add(readTimeout))

		buf := make([]byte, unsafe.Sizeof(Outgauge{}))
		nBytes, _, err := ut.connection.ReadFromUDP(buf)
		if err != nil {
			return
		}

		select {
		case ut.dataChan <- buf[:nBytes]:
		default:
			<-ut.dataChan
			ut.dataChan <- buf[:nBytes]
		}
	}
}

func (ut *UDPTransport) Write(data []byte) (int, error) {
	return ut.connection.Write(data)
}

func (ut *UDPTransport) Read(buffer []byte) (int, error) {
	data, ok := <-ut.dataChan
	if !ok {
		slog.Error(ErrNoData.Error())
		return 0, ErrNoData
	}

	return copy(buffer, data), nil
}

func (ut *UDPTransport) Close() error {
	if ut.connection != nil {
		return ut.connection.Close()
	}

	return nil
}
