package ecu

import (
	"math"
	"ml"
)

type ECUComparator struct {
	currentGear       int
	maxRPMPerGear     map[int]int
	speedRangePerGear map[int][2]int
}

func NewECUComparator() *ECUComparator {
	return &ECUComparator{
		maxRPMPerGear: map[int]int{
			1: 4000,
			2: 3500,
			3: 3000,
			4: 2500,
			5: 2000,
		},
		speedRangePerGear: map[int][2]int{
			1: {0, 20},
			2: {15, 40},
			3: {30, 70},
			4: {50, 100},
			5: {70, 150},
		},
	}
}

func (e *ECUComparator) SetCurrentGear(gear int) {
	e.currentGear = gear
}

type RPMComparator struct {
	*ECUComparator
}

func (c *RPMComparator) Compare(prev, current float64) float64 {
	// Check absolute change (max 1000 per second)
	if math.Abs(current-prev) > 1000 {
		return 1.0
	}

	// Check idle range (800-1000)
	if current < 800 && c.currentGear > 0 {
		return 1.0
	}

	// Check max RPM for current gear
	if c.currentGear > 0 && current > float64(c.maxRPMPerGear[c.currentGear]) {
		return 1.0
	}

	// Check gear shift rules
	if c.currentGear > 1 && current < 1500 {
		return 1.0 // Should downshift
	}

	shiftUpRPM := map[int]float64{
		1: 3000,
		2: 2800,
		3: 2500,
		4: 2200,
	}
	if rpm, exists := shiftUpRPM[c.currentGear]; exists && current > rpm {
		return 1.0 // Should upshift
	}

	return 0.0
}

type GearComparator struct{}

func (c *GearComparator) Compare(prev, current float64) float64 {
	// Can only change by 1 at a time
	if math.Abs(current-prev) > 1 {
		return 1.0
	}
	// Must be between 0-5
	if current < 0 || current > 5 {
		return 1.0
	}
	return 0.0
}

type SpeedComparator struct {
	*ECUComparator
}

func (c *SpeedComparator) Compare(prev, current float64) float64 {
	// Check max change (5 km/h per second)
	if math.Abs(current-prev) > 5 {
		return 1.0
	}

	// Check speed range for current gear
	if c.currentGear > 0 {
		speedRange := c.speedRangePerGear[c.currentGear]
		if current < float64(speedRange[0]) || current > float64(speedRange[1]) {
			return 1.0
		}
	}

	return 0.0
}

// Function to create ECU configs
func CreateECUConfigs() []ml.FeatureConfig {
	ecuComp := NewECUComparator()

	return []ml.FeatureConfig{
		{
			Name:       "rpm",
			Threshold:  0.5,
			Comparator: &RPMComparator{ECUComparator: ecuComp},
		},
		{
			Name:       "gear",
			Threshold:  0.5,
			Comparator: &GearComparator{},
		},
		{
			Name:       "speed",
			Threshold:  0.5,
			Comparator: &SpeedComparator{ECUComparator: ecuComp},
		},
	}
}
