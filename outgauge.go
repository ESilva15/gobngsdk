package bngsdk

import "encoding/binary"

var outgaugeSize = binary.Size(Outgauge{})

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
