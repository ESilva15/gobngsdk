package bngsdk

import (
	"net"
)

type UDPTransport struct {
	address    *net.UDPAddr
	connection *net.UDPConn
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

	return &UDPTransport{
		address:    addr,
		connection: conn,
	}, nil
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

func (ut *UDPTransport) Write(data []byte) (int, error) {
	return ut.connection.Write(data)
}

func (ut *UDPTransport) Read(buffer []byte) (int, error) {
	n, _, err := ut.connection.ReadFromUDP(buffer)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (ut *UDPTransport) Close() error {
	if ut.connection != nil {
		return ut.connection.Close()
	}

	return nil
}
