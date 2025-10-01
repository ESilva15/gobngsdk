package bngsdk

// More documentation at https://go.beamng.com/protocols.

/* TODO
   - OG_x Flags
   - DL_x Flags
*/

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

// ToMap creates a map with the data in the Outgauge struct
func (o *Outgauge) ToMap() map[string]any {
	return map[string]any{
		"Time":        o.Time,        // time in milliseconds (to check order)
		"Car":         o.Car,         // Car name
		"Flags":       o.Flags,       // Info (see OG_x below)
		"Gear":        o.Gear,        // Reverse:0, Neutral:1, First:2...
		"Plid":        o.Plid,        // Unique ID of viewed player (0 = none)
		"Speed":       o.Speed,       // M/S
		"RPM":         o.RPM,         // RPM
		"Turbo":       o.Turbo,       // BAR
		"EngTemp":     o.EngTemp,     // C
		"Fuel":        o.Fuel,        // 0 to 1
		"OilPressure": o.OilPressure, // BAR
		"OilTemp":     o.OilTemp,     // C
		"DashLights":  o.DashLights,  // Dash lights available (see DL_x below)
		"ShowLights":  o.ShowLights,  // Dash lights currently switched on
		"Throttle":    o.Throttle,    // 0 to 1
		"Brake":       o.Brake,       // 0 to 1
		"Clutch":      o.Clutch,      // 0 to 1
		"Display1":    o.Display1,    // Usually Fuel
		"Display2":    o.Display2,    // Usually Settings
		"ID":          o.ID,          // optional - only if OutGauge ID is specified
	}
}
