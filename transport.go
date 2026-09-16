package bngsdk

import (
	"errors"
	"log/slog"
	"net"
	"sync"
	"time"
	"unsafe"
)

// ErrReaderDisabled = errors.New("reader is currently turned off")
var (
	ErrNoData   = errors.New("no new data available")
	readTimeout = (time.Second / 60) * 5 // N missed frames at 60fps
	packetPool  = sync.Pool{
		New: func() any {
			var b packetBuffer
			return &b
		},
	}
)

type packetBuffer [unsafe.Sizeof(Outgauge{})]byte

type frame struct {
	Buf *packetBuffer
	Len int
}

type UDPTransport struct {
	address    *net.UDPAddr
	connection *net.UDPConn
	dataChan   chan frame
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
		dataChan:   make(chan frame, 1),
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

		bufPtr := packetPool.Get().(*packetBuffer)

		nBytes, _, err := ut.connection.ReadFromUDP(bufPtr[:])
		if err != nil {
			return
		}

		frame := frame{
			Buf: bufPtr,
			Len: nBytes,
		}

		select {
		case ut.dataChan <- frame:
			// Packet sent successfuly
		default:
			select {
			case oldFrame := <-ut.dataChan:
				packetPool.Put(oldFrame.Buf)
			default:
			}
			ut.dataChan <- frame
		}
	}
}

func (ut *UDPTransport) Write(data []byte) (int, error) {
	return ut.connection.Write(data)
}

func (ut *UDPTransport) Read(buffer []byte) (int, error) {
	latestFrame, ok := <-ut.dataChan
	if !ok {
		slog.Error(ErrNoData.Error())
		return 0, ErrNoData
	}

	nBytes := copy(buffer, latestFrame.Buf[:latestFrame.Len])
	packetPool.Put(latestFrame.Buf)

	return nBytes, nil
}

func (ut *UDPTransport) Close() error {
	if ut.connection != nil {
		return ut.connection.Close()
	}

	return nil
}
