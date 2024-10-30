package ecu

import (
	"fmt"
	"ml"
	"strconv"
)

// Struktur data untuk ECU
type ECUData struct {
	RPM      int
	Gear     int
	Speed    int
	IsAttack bool // true jika status=1
}

func (e ECUData) IsAnomaly() bool { return e.IsAttack }

func (e ECUData) GetFeatureValue(feature int) int {
	switch feature {
	case 0:
		return e.RPM
	case 1:
		return e.Gear
	case 2:
		return e.Speed
	default:
		return 0
	}
}

func (c ECUData) GetFeatureName(feature int) string {
	switch feature {
	case 0:
		return "rpm"
	case 1:
		return "gear"
	case 2:
		return "speed"
	default:
		return ""
	}
}

func (e ECUData) GetFeatureCount() int {
	return 3 // RPM, Gear, Speed
}

func CreateECUData(values []string) (ml.FeatureProvider, error) {

	// if len(values) != 4 {
	// 	return nil, fmt.Errorf("expected 4 values, got %d", len(values))
	// }

	rpm, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, fmt.Errorf("invalid RPM value: %v", err)
	}

	gear, err := strconv.Atoi(values[1])
	if err != nil {
		return nil, fmt.Errorf("invalid Gear value: %v", err)
	}

	speed, err := strconv.Atoi(values[2])
	if err != nil {
		return nil, fmt.Errorf("invalid Speed value: %v", err)
	}

	status, err := strconv.Atoi(values[3])
	if err != nil {
		return nil, fmt.Errorf("invalid Status value: %v", err)
	}

	return ECUData{
		RPM:      rpm,
		Gear:     gear,
		Speed:    speed,
		IsAttack: status == 1,
	}, nil
}

func GetSequentialAnomalyDetector() *ml.WindowDetector {

	detector := ml.NewWindowDetector(3)

	configs := CreateECUConfigs()
	for _, config := range configs {
		detector.AddFeatureConfig(config)
	}

	// Gunakan
	// isAnomaly, err := detector.AddData(yourData)
	// if err != nil {
	// 	panic(err)
	// }

	return detector
}
