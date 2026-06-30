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
	// DefaultUDPIP is the IP of the UDP outgauge server
	DefaultUDPIP = "127.0.0.1"
	// DefaultUDPPort is the port of the UDP outgauge server
	DefaultUDPPort = 4444
)

type BeamNGSDK struct {
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
func (sdk *BeamNGSDK) ReadData() error {
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
	// NOTE: create a reader for this struct and then reset it with the data from here
	// instead of creating this one everytime
	reader := bytes.NewReader(sdk.Buffer[:n])
	if err := binary.Read(reader, binary.LittleEndian, sdk.Data); err != nil {
		fmt.Println("Error decoding UDP packet:", err)
		return err
	}

	// Update the local map
	// NOTE: maybe stop doing this here and make the user request this when he
	// explicitly wants it
	sdk.DataDict = sdk.Data.ToMap()

	// this means there's new data
	return nil
}

func (sdk *BeamNGSDK) Close() error {
	return sdk.Conn.Close()
}

// Init initializes a BeamNG SDK struct
// NOTE: Change this to output a *BeamNGSDK
func Init(ip string, port int) (BeamNGSDK, error) {
	var err error
	sdk := BeamNGSDK{}

	// Create the connection to the OutGauge server
	sdk.Conn, sdk.Addr, err = createUDPConnection(ip, port)
	if err != nil {
		return BeamNGSDK{}, err
	}

	// Initiate the data variables
	sdk.Buffer = make([]byte, 1024)
	sdk.Data = &Outgauge{}
	sdk.DataDict = sdk.Data.ToMap()

	return sdk, nil
}

// SDK utilities

// ShowLights - functions to check if a given dash light is on [START]

// ShiftLight returns:
// true if the shift light is on
// false if the shift light is off
func (sdk *BeamNGSDK) ShiftLight() bool {
	return sdk.Data.ShowLights&DL_SHIFT != 0
}

// HighBeam returns:
// true if the high beams are on
// false if the high beams are off
func (sdk *BeamNGSDK) HighBeam() bool {
	return sdk.Data.ShowLights&DL_FULLBEAM != 0
}

// Handbrake returns:
// true if the handbrake is pulled
// false if the handbrake is down
func (sdk *BeamNGSDK) Handbrake() bool {
	return sdk.Data.ShowLights&DL_HANDBRAKE != 0
}

// Pitspeed returns:
// true if the pit speed limiter is on
// false if the pit speed limiter is off
// NOTE: this may not be used in BeamNG.drive, haven't checked yet
func (sdk *BeamNGSDK) Pitspeed() bool {
	return sdk.Data.ShowLights&DL_HANDBRAKE != 0
}

// TractionControl returns:
// true if traction control is on
// false if traction control is off
func (sdk *BeamNGSDK) TractionControl() bool {
	return sdk.Data.ShowLights&DL_TC != 0
}

// LeftIndicator returns:
// true if the left indicator is on
// false if the left indicator is off
func (sdk *BeamNGSDK) LeftIndicator() bool {
	return sdk.Data.ShowLights&DL_SIGNAL_L != 0
}

// RightIndicator returns:
// true if the right indicator is on
// false if the right indicator is off
func (sdk *BeamNGSDK) RightIndicator() bool {
	return sdk.Data.ShowLights&DL_SIGNAL_R != 0
}

// AnyIndicator returns:
// true if the any indicator is on
// false if the all indicators are off
func (sdk *BeamNGSDK) AnyIndicator() bool {
	return sdk.Data.ShowLights&DL_SIGNAL_ANY != 0
}

// OilLight returns
// true if the oil light is on
// false if the oil light is off
func (sdk *BeamNGSDK) OilLight() bool {
	return sdk.Data.ShowLights&DL_OILWARN != 0
}

// BatteryLight returns
// true if the battery light is on
// false if the battery light is off
func (sdk *BeamNGSDK) BatteryLight() bool {
	return sdk.Data.ShowLights&DL_BATTERY != 0
}

// ABS returns
// true if the ABS is engaged
// false if the ABS isn't engaged
func (sdk *BeamNGSDK) ABS() bool {
	return sdk.Data.ShowLights&DL_ABS != 0
}

// ShowLights - functions to check if a given dash light is on [END]

// DashLights - functions to check if a given dash light is provided [START]

// HasShiftLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasShiftLight() bool {
	return sdk.Data.DashLights&DL_SHIFT != 0
}

// HasHighBeamLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasHighBeamLight() bool {
	return sdk.Data.DashLights&DL_FULLBEAM != 0
}

// HasHandbrakeLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasHandbrakeLight() bool {
	return sdk.Data.DashLights&DL_HANDBRAKE != 0
}

// HasPitspeed returns:
// true if its available
// false if its unavailable
// NOTE: this may not be used in BeamNG.drive, haven't checked yet
func (sdk *BeamNGSDK) HasPitspeed() bool {
	return sdk.Data.DashLights&DL_HANDBRAKE != 0
}

// HasTractionControlLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasTractionControlLight() bool {
	return sdk.Data.DashLights&DL_TC != 0
}

// HasLeftIndicatorLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasLeftIndicatorLight() bool {
	return sdk.Data.DashLights&DL_SIGNAL_L != 0
}

// HasRightIndicatorLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasRightIndicatorLight() bool {
	return sdk.Data.DashLights&DL_SIGNAL_R != 0
}

// HasAnyIndicatorLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasAnyIndicatorLight() bool {
	return sdk.Data.DashLights&DL_SIGNAL_ANY != 0
}

// HasOilLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasOilLight() bool {
	return sdk.Data.DashLights&DL_OILWARN != 0
}

// HasBatteryLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasBatteryLight() bool {
	return sdk.Data.DashLights&DL_BATTERY != 0
}

// HasABSLight returns:
// true if its available
// false if its unavailable
func (sdk *BeamNGSDK) HasABSLight() bool {
	return sdk.Data.DashLights&DL_ABS != 0
}

// DashLights - functions to check if a given dash light is provided [END]
