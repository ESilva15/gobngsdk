package bngsdk

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
)

func BenchmarkReadData(b *testing.B) {
	// Spin up an UDP server
	sdk, err := Init("127.0.0.1", 0)
	if err != nil {
		b.Fatalf("Failed to initialize SDK: %v", err)
	}
	defer sdk.Close()

	// Retrieve the actual assigned UDP address
	localAddr := sdk.Conn.LocalAddr().(*net.UDPAddr)

	// Start a client to stream data
	clientConn, err := net.DialUDP("udp", nil, localAddr)
	if err != nil {
		b.Fatalf("Failed to dial local UDP socket: %v", err)
	}
	defer clientConn.Close()

	// Pre serialize some data
	dummyOutgauge := Outgauge{
		Time:  424242,
		Car:   [4]byte{'P', 'E', 'R', 'F'},
		Speed: 45.2,
		RPM:   3500.0,
	}
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, dummyOutgauge); err != nil {
		b.Fatalf("Failed to serialize dummy struct: %v", err)
	}
	packetBytes := buf.Bytes()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Feed a packet into the network buffer right before reading it
		_, err := clientConn.Write(packetBytes)
		if err != nil {
			b.Fatalf("Failed to write to UDP socket: %v", err)
		}

		// Execute the target function
		err = sdk.ReadData()
		if err != nil {
			b.Fatalf("ReadData failed at iteration %d: %v", i, err)
		}
	}
}
