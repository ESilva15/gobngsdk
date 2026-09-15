package bngsdk

type SocketWriter struct {
	udpConnection *UDPTransport
}

func NewSocketWriter(ip string, port int) (*SocketWriter, error) {
	conn, err := NewUDPWriter(ip, port)
	if err != nil {
		return nil, err
	}

	return &SocketWriter{
		udpConnection: conn,
	}, nil
}

func (sw *SocketWriter) Write(data []byte) (int, error) {
	// slog.Debug("Writing", "data", data)
	return sw.udpConnection.Write(data)
}
