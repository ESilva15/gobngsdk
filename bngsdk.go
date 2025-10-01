// Package bngsdk defines an API to interact with the BeamNG outgauge data in
// Go
package bngsdk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	// UDPIP is the IP of the UDP outgauge server
	UDPIP = "127.0.0.1"
	// UDPPort is the port of the UDP outgauge server
	UDPPort = 4444
)

type BNGSDK struct {
	Addr     *net.UDPAddr
	Conn     *net.UDPConn
	Buffer   []byte
	Data     *Outgauge
	DataDict map[string]any
}

func createUDPConnection(ip string, port int) (*net.UDPConn, *net.UDPAddr, error) {
	// Define the IP address and port to listen on
	addr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	// Create a UDP socket
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, nil, err
	}

	return conn, addr, nil
}

// ReadData will read new data from the UDP server
func (sdk *BNGSDK) ReadData() error {
	// Receive data from the socket
	n, _, err := sdk.Conn.ReadFromUDP(sdk.Buffer)
	if err != nil {
		fmt.Println("Error reading from UDP:", err)
		return err
	}

	// Check if enough data was received to fill our struct
	if n < binary.Size(Outgauge{}) {
		fmt.Println("Received packet too small for Outgauge struct")
		return err
	}

	// Read the binary data into the struct
	reader := bytes.NewReader(sdk.Buffer[:n])
	if err := binary.Read(reader, binary.LittleEndian, sdk.Data); err != nil {
		fmt.Println("Error decoding UDP packet:", err)
		return err
	}

	// Update the local map
	sdk.DataDict = sdk.Data.ToMap()

	// this means there's new data
	return nil
}

func (sdk *BNGSDK) Close() {
	_ = sdk.Conn.Close()
}

// Init initializes a BeamNG SDK struct
func Init(ip string, port int) (BNGSDK, error) {
	var err error
	sdk := BNGSDK{}

	// Create the connection to the OutGauge server
	sdk.Conn, sdk.Addr, err = createUDPConnection(ip, port)
	if err != nil {
		return BNGSDK{}, err
	}

	// Initiate the data variables
	sdk.Buffer = make([]byte, 1024)
	sdk.Data = &Outgauge{}
	sdk.DataDict = sdk.Data.ToMap()

	return sdk, nil
}
