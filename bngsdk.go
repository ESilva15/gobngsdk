// Package bngsdk defines an API to interact with the BeamNG outgauge data in Go
package bngsdk

import (
	"encoding/binary"
	"io"
	"log/slog"
	"math"
	"unsafe"
)

const (
	// DefaultUDPIP is the IP of the UDP outgauge server
	DefaultUDPIP = "127.0.0.1"
	// DefaultUDPPort is the port of the UDP outgauge server
	DefaultUDPPort = 4444
)

type TelemetryContainer int

const (
	BinaryFile TelemetryContainer = iota
	UDPData    TelemetryContainer = iota
)

type Options struct {
	Logger *slog.Logger
	// Import settings
	SourceType       TelemetryContainer // type of source data
	BinSourcePath    string             // Path to source binary
	ImportUDPAddress string
	ImportUDPPort    int
	Loop             bool // Wheter to loop when we reach the end of file
	// Export settings
	ExportUDPAddress string
	ExportUDPPort    int
	ExportDataType   TelemetryContainer // export type of telemetry: store .bin or replay in UDP
	ExportDataPath   string             // path where to export the data
	ExportData       bool               // whether to export the telemetry data
}

type BeamNGSDK struct {
	Opts   Options
	reader BngImporter
	writer BngExporter
	Data   Outgauge
	buffer []byte
}

func NewBngSDK(opts Options) (*BeamNGSDK, error) {
	sdk := BeamNGSDK{
		Opts:   opts,
		buffer: make([]byte, unsafe.Sizeof(Outgauge{})),
	}

	// Set the passed logger as the default logger
	slog.SetDefault(sdk.Opts.Logger)

	var err error

	err = sdk.openReader()
	if err != nil {
		slog.Error("failed to open reader", "err", err)
		return nil, err
	}

	err = sdk.openWriter()
	if err != nil {
		slog.Error("failed to open writer", "err", err)
		return nil, err
	}

	return &sdk, nil
}

func (sdk *BeamNGSDK) Close() error {
	sdk.buffer = nil

	if sdk.writer != nil {
		sdk.writer.Close()
	}

	if sdk.reader != nil {
		sdk.reader.Close()
	}

	return nil
}

func (sdk *BeamNGSDK) openReader() error {
	switch sdk.Opts.SourceType {
	case BinaryFile:
		reader, err := NewBinaryImporter(sdk.Opts.BinSourcePath)
		if err != nil {
			slog.Error("failed to create BinaryImporter",
				"path", sdk.Opts.BinSourcePath, "err", err)
			return err
		}

		slog.Info(
			"created BinaryImporter",
			"path", sdk.Opts.BinSourcePath,
		)
		sdk.reader = reader
	case UDPData:
		reader, err := NewSocketImporter(sdk.Opts.ImportUDPAddress, sdk.Opts.ImportUDPPort)
		if err != nil {
			slog.Error(
				"failed to create SocketImporter",
				"address", sdk.Opts.ImportUDPAddress, "port", sdk.Opts.ImportUDPPort, "err", err,
			)
			return err
		}

		slog.Info(
			"created SocketImporter",
			"address", sdk.Opts.ImportUDPAddress, "port", sdk.Opts.ImportUDPPort,
		)
		sdk.reader = reader
	}

	return nil
}

func (sdk *BeamNGSDK) openWriter() error {
	// if the user didn't request data export we don't need a writer
	if !sdk.Opts.ExportData {
		return nil
	}

	switch sdk.Opts.ExportDataType {
	case BinaryFile:
		writer, err := NewBinaryExporter(sdk.Opts.ExportDataPath)
		if err != nil {
			slog.Error("failed to create BinaryExporter",
				"path", sdk.Opts.ExportDataPath, "err", err)
			return err
		}

		slog.Info(
			"created BinaryExporter",
			"path", sdk.Opts.ExportDataPath,
		)
		sdk.writer = writer
	case UDPData:
		writer, err := NewSocketExporter(sdk.Opts.ExportUDPAddress, sdk.Opts.ExportUDPPort)
		if err != nil {
			slog.Error(
				"failed to create SocketExporter",
				"address", sdk.Opts.ExportUDPAddress, "port", sdk.Opts.ExportUDPPort, "err", err,
			)
			return err
		}

		slog.Info(
			"created SocketExporter",
			"address", sdk.Opts.ExportUDPAddress, "port", sdk.Opts.ExportUDPPort,
		)
		sdk.writer = writer
	}

	return nil
}

func (sdk *BeamNGSDK) Update() (int, error) {
	var nBytes int
	var err error

	nBytes, err = sdk.reader.Next(sdk.buffer)
	// We check this first because we want to know if we need to loop
	if err == io.EOF {
		if sdk.Opts.Loop {
			err = sdk.reader.Reset()
			if err != nil {
				return 0, err
			}
		}

		// NOTE: could we make the reset return the next piece of data?
		// We update to the start of the file since we had to reset
		nBytes, err = sdk.reader.Next(sdk.buffer)
	}
	if err != nil {
		return 0, err
	}

	if sdk.Opts.ExportData {
		sdk.writer.Write(sdk.buffer)
	}

	return nBytes, sdk.parseData(sdk.buffer)
}

func ParseData(ogData *Outgauge, buffer []byte) error {
	ogData.Time = binary.LittleEndian.Uint32(buffer[0:4])
	copy(ogData.Car[:], buffer[4:8])
	ogData.Flags = binary.LittleEndian.Uint16(buffer[8:10])
	ogData.Gear = int8(buffer[10])
	ogData.Plid = int8(buffer[11])
	ogData.Speed = math.Float32frombits(binary.LittleEndian.Uint32(buffer[12:16]))
	ogData.RPM = math.Float32frombits(binary.LittleEndian.Uint32(buffer[16:20]))
	ogData.Turbo = math.Float32frombits(binary.LittleEndian.Uint32(buffer[20:24]))
	ogData.EngTemp = math.Float32frombits(binary.LittleEndian.Uint32(buffer[24:28]))
	ogData.Fuel = math.Float32frombits(binary.LittleEndian.Uint32(buffer[28:32]))
	ogData.OilPressure = math.Float32frombits(binary.LittleEndian.Uint32(buffer[32:36]))
	ogData.OilTemp = math.Float32frombits(binary.LittleEndian.Uint32(buffer[36:40]))
	ogData.DashLights = binary.LittleEndian.Uint32(buffer[40:44])
	ogData.ShowLights = binary.LittleEndian.Uint32(buffer[44:48])
	ogData.Throttle = math.Float32frombits(binary.LittleEndian.Uint32(buffer[48:52]))
	ogData.Brake = math.Float32frombits(binary.LittleEndian.Uint32(buffer[52:56]))
	ogData.Clutch = math.Float32frombits(binary.LittleEndian.Uint32(buffer[56:60]))
	copy(ogData.Display1[:], buffer[60:76])
	copy(ogData.Display2[:], buffer[76:92])
	ogData.ID = int32(binary.LittleEndian.Uint32(buffer[92:96]))

	return nil
}

func (sdk *BeamNGSDK) parseData(buffer []byte) error {
	return ParseData(&sdk.Data, buffer)
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
