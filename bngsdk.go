// Package bngsdk defines an API to interact with the BeamNG outgauge data in Go
package bngsdk

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

const (
	// DefaultUDPIP is the IP of the UDP outgauge server
	DefaultUDPIP = "127.0.0.1"
	// DefaultUDPPort is the port of the UDP outgauge server
	DefaultUDPPort = 4444
)

type BeamNGSDK struct {
	Addr   *net.UDPAddr
	Conn   *net.UDPConn
	Buffer []byte
	Data   Outgauge
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

// parseData will parse the bytes read from the socket into the Outgauge struct
func (sdk *BeamNGSDK) parseData() error {
	sdk.Data.Time = binary.LittleEndian.Uint32(sdk.Buffer[0:4])
	copy(sdk.Data.Car[:], sdk.Buffer[4:8])
	sdk.Data.Flags = binary.LittleEndian.Uint16(sdk.Buffer[8:10])
	sdk.Data.Gear = int8(sdk.Buffer[10])
	sdk.Data.Plid = int8(sdk.Buffer[11])
	sdk.Data.Speed = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[12:16]))
	sdk.Data.RPM = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[16:20]))
	sdk.Data.Turbo = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[20:24]))
	sdk.Data.EngTemp = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[24:28]))
	sdk.Data.Fuel = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[28:32]))
	sdk.Data.OilPressure = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[32:36]))
	sdk.Data.OilTemp = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[36:40]))
	sdk.Data.DashLights = binary.LittleEndian.Uint32(sdk.Buffer[40:44])
	sdk.Data.ShowLights = binary.LittleEndian.Uint32(sdk.Buffer[44:48])
	sdk.Data.Throttle = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[48:52]))
	sdk.Data.Brake = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[52:56]))
	sdk.Data.Throttle = math.Float32frombits(binary.LittleEndian.Uint32(sdk.Buffer[56:60]))
	copy(sdk.Data.Display1[:], sdk.Buffer[60:76])
	copy(sdk.Data.Display2[:], sdk.Buffer[76:92])
	sdk.Data.ID = int32(binary.LittleEndian.Uint32(sdk.Buffer[92:96]))

	return nil
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
	if n < outgaugeSize {
		fmt.Println("Received packet too small for Outgauge struct")
		return err
	}

	// Read the binary data into the struct
	err = sdk.parseData()
	if err != nil {
		fmt.Println("Error decoding UDP packet:", err)
		return err
	}

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

	return sdk, nil
}

// SDK utilities

// ShowLights - functions to check if a given dash light is on [START]

// ShiftLight reports whether the shift light is on
func (sdk *BeamNGSDK) ShiftLight() bool {
	return sdk.Data.ShowLights&DL_SHIFT != 0
}

// HighBeam reports whether high beams are on
func (sdk *BeamNGSDK) HighBeam() bool {
	return sdk.Data.ShowLights&DL_FULLBEAM != 0
}

// Handbrake reports whether the handbrake is pulled
func (sdk *BeamNGSDK) Handbrake() bool {
	return sdk.Data.ShowLights&DL_HANDBRAKE != 0
}

// Pitspeed reports whether the pit speed limiter is engaged
//
// NOTE: this may not be used in BeamNG.drive, haven't checked yet
func (sdk *BeamNGSDK) Pitspeed() bool {
	return sdk.Data.ShowLights&DL_PITSPEED != 0
}

// TractionControl reports wheter TC is engaged
func (sdk *BeamNGSDK) TractionControl() bool {
	return sdk.Data.ShowLights&DL_TC != 0
}

// LeftIndicator reports whether the left indicator is on
func (sdk *BeamNGSDK) LeftIndicator() bool {
	return sdk.Data.ShowLights&DL_SIGNAL_L != 0
}

// RightIndicator reports wheter the right indicator is on
func (sdk *BeamNGSDK) RightIndicator() bool {
	return sdk.Data.ShowLights&DL_SIGNAL_R != 0
}

// AnyIndicator reports whether any indicator is on
func (sdk *BeamNGSDK) AnyIndicator() bool {
	return sdk.Data.ShowLights&DL_SIGNAL_ANY != 0
}

// OilLight reports whether the oil warning light is on
func (sdk *BeamNGSDK) OilLight() bool {
	return sdk.Data.ShowLights&DL_OILWARN != 0
}

// BatteryLight reports whether the battery light is on
func (sdk *BeamNGSDK) BatteryLight() bool {
	return sdk.Data.ShowLights&DL_BATTERY != 0
}

// ABS reports whether the ABS light is on
func (sdk *BeamNGSDK) ABS() bool {
	return sdk.Data.ShowLights&DL_ABS != 0
}

// ShowLights - functions to check if a given dash light is on [END]

// DashLights - functions to check if a given dash light is provided [START]

// HasShiftLight reports whether a shift light is available
func (sdk *BeamNGSDK) HasShiftLight() bool {
	return sdk.Data.DashLights&DL_SHIFT != 0
}

// HasHighBeamLight reports whether a high beam light is available
func (sdk *BeamNGSDK) HasHighBeamLight() bool {
	return sdk.Data.DashLights&DL_FULLBEAM != 0
}

// HasHandbrakeLight reports wheter a handbrake light is available
func (sdk *BeamNGSDK) HasHandbrakeLight() bool {
	return sdk.Data.DashLights&DL_HANDBRAKE != 0
}

// HasPitspeed reports whether a pit speed limitr is available
// NOTE: this may not be used in BeamNG.drive, haven't checked yet
func (sdk *BeamNGSDK) HasPitspeed() bool {
	return sdk.Data.DashLights&DL_HANDBRAKE != 0
}

// HasTractionControlLight reports whether a traction control light is available
func (sdk *BeamNGSDK) HasTractionControlLight() bool {
	return sdk.Data.DashLights&DL_TC != 0
}

// HasLeftIndicatorLight reports whether a left indicator is available
func (sdk *BeamNGSDK) HasLeftIndicatorLight() bool {
	return sdk.Data.DashLights&DL_SIGNAL_L != 0
}

// HasRightIndicatorLight reports whether a right indicator is available
func (sdk *BeamNGSDK) HasRightIndicatorLight() bool {
	return sdk.Data.DashLights&DL_SIGNAL_R != 0
}

// HasAnyIndicatorLight reports whether an any indicator light is available
func (sdk *BeamNGSDK) HasAnyIndicatorLight() bool {
	return sdk.Data.DashLights&DL_SIGNAL_ANY != 0
}

// HasOilLight reports whether a oil light is available
func (sdk *BeamNGSDK) HasOilLight() bool {
	return sdk.Data.DashLights&DL_OILWARN != 0
}

// HasBatteryLight reports whether a battery light is available
func (sdk *BeamNGSDK) HasBatteryLight() bool {
	return sdk.Data.DashLights&DL_BATTERY != 0
}

// HasABSLight reports whether an ABS light is available
func (sdk *BeamNGSDK) HasABSLight() bool {
	return sdk.Data.DashLights&DL_ABS != 0
}

// DashLights - functions to check if a given dash light is provided [END]

// Flags - functions to check if a given flag is ON [START]

// HasTurbo reports whether there's a turbo
func (sdk *BeamNGSDK) HasTurbo() bool {
	return sdk.Data.Flags&OG_TURBO != 0
}

// PrefersKm reports whether the user prefers kilometers:
//   - true is prefers Km
//   - false is prefers Mi
func (sdk *BeamNGSDK) PrefersKm() bool {
	return sdk.Data.Flags&OG_KM != 0
}

// PrefersBAR reports whether the user prefers BAR:
//   - true is prefers BAR
//   - false is prefers PSI
func (sdk *BeamNGSDK) PrefersBAR() bool {
	return sdk.Data.Flags&OG_BAR != 0
}

// Flags - functions to check if a given flag is ON [END]

// Data Retrieval [START]

// ToMap creates a map with the data in the Outgauge struct
func (sdk *BeamNGSDK) ToMap() map[string]any {
	return map[string]any{
		"Time":        sdk.Data.Time,        // time in milliseconds (to check order)
		"Car":         sdk.Data.Car,         // Car name
		"Flags":       sdk.Data.Flags,       // Info (see OG_x below)
		"Gear":        sdk.Data.Gear,        // Reverse:0, Neutral:1, First:2...
		"Plid":        sdk.Data.Plid,        // Unique ID of viewed player (0 = none)
		"Speed":       sdk.Data.Speed,       // M/S
		"RPM":         sdk.Data.RPM,         // RPM
		"Turbo":       sdk.Data.Turbo,       // BAR
		"EngTemp":     sdk.Data.EngTemp,     // C
		"Fuel":        sdk.Data.Fuel,        // 0 to 1
		"OilPressure": sdk.Data.OilPressure, // BAR
		"OilTemp":     sdk.Data.OilTemp,     // C
		"DashLights":  sdk.Data.DashLights,  // Dash lights available (see DL_x below)
		"ShowLights":  sdk.Data.ShowLights,  // Dash lights currently switched on
		"Throttle":    sdk.Data.Throttle,    // 0 to 1
		"Brake":       sdk.Data.Brake,       // 0 to 1
		"Clutch":      sdk.Data.Clutch,      // 0 to 1
		"Display1":    sdk.Data.Display1,    // Usually Fuel
		"Display2":    sdk.Data.Display2,    // Usually Settings
		"ID":          sdk.Data.ID,          // optional - only if OutGauge ID is specified
	}
}

// Data Retrieval [END]
