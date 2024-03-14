/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

risk_level.go

Risk level values
*/
package sandbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"strings"
)

type RiskLevel int

const (
	RiskLevelUnknown RiskLevel = iota
	RiskLevelNotReady
	RiskLevelUnsupported
	RiskLevelNoRisk
	RiskLevelLow
	RiskLevelMedium
	RiskLevelHigh
	RiskLevelError
)

func RiskLevelThreat(r RiskLevel) bool {
	return r == RiskLevelLow || r == RiskLevelMedium || r == RiskLevelHigh
}

var RiskLevelString = [...]string{
	"Unknown",
	"Not Ready",
	"Not Analyzed",
	"No Risk",
	"Low Risk",
	"Medium Risk",
	"High Risk",
	"Error",
}

func (r RiskLevel) String() string {
	return RiskLevelString[r]
}

var RiskLevelColor = [...]color.Color{
	color.RGBA{0, 0, 0, 255},
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

var ErrUnknownRiskLevel = errors.New("unknown risk level")

// UnmarshalJSON implements the Unmarshaler interface of the json package for RiskLevel.
func (r *RiskLevel) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	for i, s := range RiskLevelString {
		if strings.EqualFold(s, v) {
			*r = RiskLevel(i)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrUnknownRiskLevel, v)
}

// MarshalJSON implements the Marshaler interface of the json package for RiskLevel.
func (r RiskLevel) MarshalJSON() ([]byte, error) {
	if r < 0 || r >= RiskLevelError {
		return nil, ErrUnknownRiskLevel
	}
	return []byte(fmt.Sprintf("\"%s\"", r.String())), nil
}
