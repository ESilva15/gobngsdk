package bngsdk

import (
	"encoding/binary"
	"errors"
	"math"
)

var (
	OutgaugeSize     = int64(binary.Size(Outgauge{}))
	ErrNotEnoughData = errors.New("buffer len is too short for Outgauge unpacking")
)

// More documentation at https://go.beamng.com/protocols.
// Or at BeamNG/lua/vehicle/protocols/outgauge.lua
// OG_x buts for flags
const (
	OG_SHIFT = 1     // key // N/A
	OG_CTRL  = 2     // key // N/A
	OG_TURBO = 8192  // show turbo gauge
	OG_KM    = 16384 // if not set - user prefers MILES
	OG_BAR   = 32768 // if not set - user prefers PSI
)

// DL_x Flags
const (
	DL_SHIFT      = 1 << 0  // shift light
	DL_FULLBEAM   = 1 << 1  // full beam
	DL_HANDBRAKE  = 1 << 2  // handbrake
	DL_PITSPEED   = 1 << 3  // pit speed limiter // N/A
	DL_TC         = 1 << 4  // tc active or switched off
	DL_SIGNAL_L   = 1 << 5  // left turn signal
	DL_SIGNAL_R   = 1 << 6  // right turn signal
	DL_SIGNAL_ANY = 1 << 7  // shared turn signal // N/A
	DL_OILWARN    = 1 << 8  // oil pressure warning
	DL_BATTERY    = 1 << 9  // battery warning
	DL_ABS        = 1 << 10 // abs active or switched off
	DL_SPARE      = 1 << 11 // N/A
)

// Outgauge describes the data served by the BeamNG outgauge server
type Outgauge struct {
	// data        []byte
	Time        uint32   // time in milliseconds (to check order)
	Car         [4]byte  // Car name
	Flags       uint16   // Info (see OG_x below)
	Gear        int8     // Reverse:0, Neutral:1, First:2...
	Plid        int8     // Unique ID of viewed player (0 = none)
	Speed       float32  // M/S
	RPM         float32  // RPM
	Turbo       float32  // BAR
	EngTemp     float32  // C
	Fuel        float32  // 0 to 1
	OilPressure float32  // BAR
	OilTemp     float32  // C
	DashLights  uint32   // Dash lights available (see DL_x below)
	ShowLights  uint32   // Dash lights currently switched on
	Throttle    float32  // 0 to 1
	Brake       float32  // 0 to 1
	Clutch      float32  // 0 to 1
	Display1    [16]byte // Usually Fuel
	Display2    [16]byte // Usually Settings
	ID          int32    // optional - only if OutGauge ID is specified
}

func (og *Outgauge) ParseData(buffer []byte) error {
	if int64(len(buffer)) < OutgaugeSize {
		return ErrNotEnoughData
	}

	og.Time = binary.LittleEndian.Uint32(buffer[0:4])
	copy(og.Car[:], buffer[4:8])
	og.Flags = binary.LittleEndian.Uint16(buffer[8:10])
	og.Gear = int8(buffer[10])
	og.Plid = int8(buffer[11])
	og.Speed = math.Float32frombits(binary.LittleEndian.Uint32(buffer[12:16]))
	og.RPM = math.Float32frombits(binary.LittleEndian.Uint32(buffer[16:20]))
	og.Turbo = math.Float32frombits(binary.LittleEndian.Uint32(buffer[20:24]))
	og.EngTemp = math.Float32frombits(binary.LittleEndian.Uint32(buffer[24:28]))
	og.Fuel = math.Float32frombits(binary.LittleEndian.Uint32(buffer[28:32]))
	og.OilPressure = math.Float32frombits(binary.LittleEndian.Uint32(buffer[32:36]))
	og.OilTemp = math.Float32frombits(binary.LittleEndian.Uint32(buffer[36:40]))
	og.DashLights = binary.LittleEndian.Uint32(buffer[40:44])
	og.ShowLights = binary.LittleEndian.Uint32(buffer[44:48])
	og.Throttle = math.Float32frombits(binary.LittleEndian.Uint32(buffer[48:52]))
	og.Brake = math.Float32frombits(binary.LittleEndian.Uint32(buffer[52:56]))
	og.Clutch = math.Float32frombits(binary.LittleEndian.Uint32(buffer[56:60]))
	copy(og.Display1[:], buffer[60:76])
	copy(og.Display2[:], buffer[76:92])
	og.ID = int32(binary.LittleEndian.Uint32(buffer[92:96]))

	return nil
}

// TODO: think about having direct memory access for the fields instead of
// parsing everything
// func NewOutgauge(buffer []byte) (*Outgauge, error) {
// 	if len(buffer) < outgaugeSize {
// 		return nil, ErrNotEnoughData
// 	}
//
// 	og := &Outgauge{
// 		data: make([]byte, outgaugeSize),
// 	}
//
// 	return og, nil
// }

// SDK utilities

// Values Getters [START]

// Values Getters [END]

// ShowLights - functions to check if a given dash light is on [START]

// ShiftLight reports whether the shift light is on
func (og *Outgauge) ShiftLight() bool {
	return og.ShowLights&DL_SHIFT != 0
}

// HighBeam reports whether high beams are on
func (og *Outgauge) HighBeam() bool {
	return og.ShowLights&DL_FULLBEAM != 0
}

// Handbrake reports whether the handbrake is pulled
func (og *Outgauge) Handbrake() bool {
	return og.ShowLights&DL_HANDBRAKE != 0
}

// Pitspeed reports whether the pit speed limiter is engaged
//
// NOTE: this may not be used in BeamNG.drive, haven't checked yet
func (og *Outgauge) Pitspeed() bool {
	return og.ShowLights&DL_PITSPEED != 0
}

// TractionControl reports wheter TC is engaged
func (og *Outgauge) TractionControl() bool {
	return og.ShowLights&DL_TC != 0
}

// LeftIndicator reports whether the left indicator is on
func (og *Outgauge) LeftIndicator() bool {
	return og.ShowLights&DL_SIGNAL_L != 0
}

// RightIndicator reports wheter the right indicator is on
func (og *Outgauge) RightIndicator() bool {
	return og.ShowLights&DL_SIGNAL_R != 0
}

// AnyIndicator reports whether any indicator is on
func (og *Outgauge) AnyIndicator() bool {
	return og.ShowLights&DL_SIGNAL_ANY != 0
}

// OilLight reports whether the oil warning light is on
func (og *Outgauge) OilLight() bool {
	return og.ShowLights&DL_OILWARN != 0
}

// BatteryLight reports whether the battery light is on
func (og *Outgauge) BatteryLight() bool {
	return og.ShowLights&DL_BATTERY != 0
}

// ABS reports whether the ABS light is on
func (og *Outgauge) ABS() bool {
	return og.ShowLights&DL_ABS != 0
}

// Spare reports whether spare was flipped
func (og *Outgauge) Spare() bool {
	return og.ShowLights&DL_SPARE != 0
}

// ShowLights - functions to check if a given dash light is on [END]

// DashLights - functions to check if a given dash light is provided [START]

// HasShiftLight reports whether a shift light is available
func (og *Outgauge) HasShiftLight() bool {
	return og.DashLights&DL_SHIFT != 0
}

// HasHighBeamLight reports whether a high beam light is available
func (og *Outgauge) HasHighBeamLight() bool {
	return og.DashLights&DL_FULLBEAM != 0
}

// HasHandbrakeLight reports wheter a handbrake light is available
func (og *Outgauge) HasHandbrakeLight() bool {
	return og.DashLights&DL_HANDBRAKE != 0
}

// HasPitspeed reports whether a pit speed limitr is available
// NOTE: this may not be used in BeamNG.drive, haven't checked yet
func (og *Outgauge) HasPitspeed() bool {
	return og.DashLights&DL_PITSPEED != 0
}

// HasTractionControlLight reports whether a traction control light is available
func (og *Outgauge) HasTractionControlLight() bool {
	return og.DashLights&DL_TC != 0
}

// HasLeftIndicatorLight reports whether a left indicator is available
func (og *Outgauge) HasLeftIndicatorLight() bool {
	return og.DashLights&DL_SIGNAL_L != 0
}

// HasRightIndicatorLight reports whether a right indicator is available
func (og *Outgauge) HasRightIndicatorLight() bool {
	return og.DashLights&DL_SIGNAL_R != 0
}

// HasAnyIndicatorLight reports whether an any indicator light is available
func (og *Outgauge) HasAnyIndicatorLight() bool {
	return og.DashLights&DL_SIGNAL_ANY != 0
}

// HasOilLight reports whether a oil light is available
func (og *Outgauge) HasOilLight() bool {
	return og.DashLights&DL_OILWARN != 0
}

// HasBatteryLight reports whether a battery light is available
func (og *Outgauge) HasBatteryLight() bool {
	return og.DashLights&DL_BATTERY != 0
}

// HasABSLight reports whether an ABS light is available
func (og *Outgauge) HasABSLight() bool {
	return og.DashLights&DL_ABS != 0
}

// HasSpare reports whether spare was flipped
func (og *Outgauge) HasSpare() bool {
	return og.DashLights&DL_SPARE != 0
}

// DashLights - functions to check if a given dash light is provided [END]

// Flags - functions to check if a given flag is ON [START]

// HasTurbo reports whether there's a turbo
func (og *Outgauge) HasTurbo() bool {
	return og.Flags&OG_TURBO != 0
}

// PrefersKm reports whether the user prefers kilometers:
//   - true is prefers Km
//   - false is prefers Mi
func (og *Outgauge) PrefersKm() bool {
	return og.Flags&OG_KM != 0
}

// PrefersBAR reports whether the user prefers BAR:
//   - true is prefers BAR
//   - false is prefers PSI
func (og *Outgauge) PrefersBAR() bool {
	return og.Flags&OG_BAR != 0
}

// Flags - functions to check if a given flag is ON [END]

// Data Retrieval [START]

// ToMap creates a map with the data in the Outgauge struct
func (sdk *BeamNGSDK) ToMap() map[string]any {
	return map[string]any{
		"Time":        sdk.data.Time,        // time in milliseconds (to check order)
		"Car":         sdk.data.Car,         // Car name
		"Flags":       sdk.data.Flags,       // Info (see OG_x below)
		"Gear":        sdk.data.Gear,        // Reverse:0, Neutral:1, First:2...
		"Plid":        sdk.data.Plid,        // Unique ID of viewed player (0 = none)
		"Speed":       sdk.data.Speed,       // M/S
		"RPM":         sdk.data.RPM,         // RPM
		"Turbo":       sdk.data.Turbo,       // BAR
		"EngTemp":     sdk.data.EngTemp,     // C
		"Fuel":        sdk.data.Fuel,        // 0 to 1
		"OilPressure": sdk.data.OilPressure, // BAR
		"OilTemp":     sdk.data.OilTemp,     // C
		"DashLights":  sdk.data.DashLights,  // Dash lights available (see DL_x below)
		"ShowLights":  sdk.data.ShowLights,  // Dash lights currently switched on
		"Throttle":    sdk.data.Throttle,    // 0 to 1
		"Brake":       sdk.data.Brake,       // 0 to 1
		"Clutch":      sdk.data.Clutch,      // 0 to 1
		"Display1":    sdk.data.Display1,    // Usually Fuel
		"Display2":    sdk.data.Display2,    // Usually Settings
		"ID":          sdk.data.ID,          // optional - only if OutGauge ID is specified
	}
}

// Data Retrieval [END]
