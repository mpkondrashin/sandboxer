package task

import (
	"image/color"
)

type RiskLevel int

const (
	RiskLevelUnknown RiskLevel = iota
	RiskLevelUnsupported
	RiskLevelNoRisk
	RiskLevelLow
	RiskLevelMedium
	RiskLevelHigh
	RiskLevelError
)

var RiskLevelString = [...]string{
	"Unknown",
	"Unsupported",
	"NoRisk",
	"Low",
	"Medium",
	"High",
	"Error",
}

func (r RiskLevel) String() string {
	return RiskLevelString[r]
}

var RiskLevelColor = [...]color.Color{
	color.RGBA{0, 0, 0, 255},
	color.RGBA{158, 158, 158, 255},
	color.RGBA{0, 180, 0, 255},
	color.RGBA{255, 153, 0, 255},
	color.RGBA{230, 102, 0, 255},
	color.RGBA{204, 51, 0, 255},
	color.RGBA{255, 0, 0, 255},
}

func (r RiskLevel) Color() color.Color {
	return RiskLevelColor[r]
}
