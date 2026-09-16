package bngsdk

import (
	"bytes"
	"encoding/binary"
	"io"
	"log/slog"
	"net"
	"testing"
)

func BenchmarkUpdate(b *testing.B) {
	// Silence logging output so slog calls don't pollute benchmark stats
	slogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Initialize the SDK with port 0 to bind to an OS-assigned ephemeral port
	sdk, err := NewBngSDK(Options{
		Logger:           slogger,
		SourceType:       UDPData,
		ImportUDPAddress: "127.0.0.1",
		ImportUDPPort:    0,
	})
	if err != nil {
		b.Fatalf("Failed to initialize SDK: %v", err)
	}
	defer sdk.Close()

	// Access the underlying reader connection to determine the dynamically bound port
	ogReader, ok := sdk.reader.(*OgUDPReader)
	if !ok || ogReader.udpConnection == nil || ogReader.udpConnection.connection == nil {
		b.Fatalf("Failed to retrieve underlying UDP connection")
	}

	serverAddr := ogReader.udpConnection.connection.LocalAddr().(*net.UDPAddr)

	// Dial the UDP socket as a client to send test data
	clientConn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		b.Fatalf("Failed to dial UDP server: %v", err)
	}
	defer clientConn.Close()

	// Pre-serialize a dummy Outgauge struct matching the required byte layout
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
		// Feed a packet into the network transport socket
		_, err := clientConn.Write(packetBytes)
		if err != nil {
			b.Fatalf("Failed to write to UDP socket: %v", err)
		}

		// Run the main API loop method
		_, err = sdk.Update()
		if err != nil {
			b.Fatalf("Update failed at iteration %d: %v", i, err)
		}
	}
}
