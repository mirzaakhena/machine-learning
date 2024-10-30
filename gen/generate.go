package gen

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
)

// VehicleData represents a single row of vehicle data
type VehicleData struct {
	RPM         int
	Gear        int
	Speed       int
	Status      int
	Description string
}

// Generator contains methods for generating vehicle data
type Generator struct {
	data map[string]bool // Used to ensure uniqueness
}

// NewGenerator creates a new Generator instance
func NewGenerator() *Generator {
	return &Generator{
		data: make(map[string]bool),
	}
}

// generateNormalCase generates a single normal case
func (g *Generator) generateNormalCase() *VehicleData {
	gear := rand.Intn(5) + 1 // 1-5

	// Calculate appropriate speed range based on gear
	minSpeed := max(5, (gear-1)*15)
	maxSpeed := min(gear*40, 200)
	speed := rand.Intn(maxSpeed-minSpeed+1) + minSpeed

	// Calculate appropriate RPM range
	minRPM := max(800, 1000*speed/(gear*40))
	maxRPM := 5500
	rpm := rand.Intn(maxRPM-minRPM+1) + minRPM

	return &VehicleData{
		RPM:         rpm,
		Gear:        gear,
		Speed:       speed,
		Status:      0,
		Description: "normal",
	}
}

// generateAnomalyCase generates a single anomaly case
func (g *Generator) generateAnomalyCase() *VehicleData {
	anomalyType := rand.Intn(5) + 1
	var data VehicleData

	switch anomalyType {
	case 1: // Over-revving
		data = VehicleData{
			RPM:         rand.Intn(1500) + 5501, // 5501-7000
			Gear:        rand.Intn(3) + 1,       // 1-3
			Speed:       rand.Intn(100) + 20,    // 20-120
			Status:      1,
			Description: "Over-revving",
		}
	case 2: // Stalling
		data = VehicleData{
			RPM:         rand.Intn(400) + 400, // 400-799
			Gear:        rand.Intn(5) + 1,     // 1-5
			Speed:       rand.Intn(21) + 10,   // 10-30
			Status:      1,
			Description: "Stalling",
		}
	case 3: // Gear-speed mismatch
		gear := rand.Intn(2) + 1 // 1-2
		data = VehicleData{
			RPM:         rand.Intn(2000) + 2000, // 2000-4000
			Gear:        gear,
			Speed:       rand.Intn(71) + 50 + gear*50, // 50-120
			Status:      1,
			Description: "Gear-speed mismatch",
		}
	case 4: // Neutral with speed
		data = VehicleData{
			RPM:         rand.Intn(2200) + 800, // 800-3000
			Gear:        0,                     // Neutral
			Speed:       rand.Intn(61) + 20,    // 20-80
			Status:      1,
			Description: "Neutral with speed",
		}
	case 5: // RPM too low for speed-gear
		gear := rand.Intn(3) + 3 // 3-5
		data = VehicleData{
			RPM:         rand.Intn(701) + 800, // 800-1500
			Gear:        gear,
			Speed:       rand.Intn(16) + gear*30, // Higher speed for gear
			Status:      1,
			Description: "RPM too low for speed-gear",
		}
	}

	return &data
}

// GenerateData generates n rows of vehicle data
func (g *Generator) GenerateData(normalCount, anomalyCount int) []VehicleData {
	result := make([]VehicleData, 0, normalCount+anomalyCount)

	// Generate normal cases
	for len(result) < normalCount {
		data := g.generateNormalCase()
		key := fmt.Sprintf("%d,%d,%d", data.RPM, data.Gear, data.Speed)
		if !g.data[key] {
			g.data[key] = true
			result = append(result, *data)
		}
	}

	// Generate anomaly cases
	for len(result) < normalCount+anomalyCount {
		data := g.generateAnomalyCase()
		key := fmt.Sprintf("%d,%d,%d", data.RPM, data.Gear, data.Speed)
		if !g.data[key] {
			g.data[key] = true
			result = append(result, *data)
		}
	}

	// Shuffle the result
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

// SaveToCSV saves the generated data to a CSV file
func SaveToCSV(data []VehicleData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"rpm", "gear", "speed", "status", "description"}); err != nil {
		return err
	}

	// Write data
	for _, row := range data {
		if err := writer.Write([]string{
			fmt.Sprintf("%d", row.RPM),
			fmt.Sprintf("%d", row.Gear),
			fmt.Sprintf("%d", row.Speed),
			fmt.Sprintf("%d", row.Status),
			row.Description,
		}); err != nil {
			return err
		}
	}

	return nil
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
