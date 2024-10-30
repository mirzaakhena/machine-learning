package ml

import (
	"fmt"
	"math"
)

type SequentialProvider interface {
	HasValueCount
	GetFeatureName(feature int) string
}

// FeatureComparator mendefinisikan bagaimana membandingkan nilai feature
type FeatureComparator interface {
	// Compare membandingkan dua nilai dan mengembalikan tingkat perubahannya
	// Mengembalikan nilai antara 0-1 yang merepresentasikan tingkat perubahan
	Compare(prev, current float64) float64
}

// DefaultComparator menggunakan relative change
type DefaultComparator struct{}

func (c DefaultComparator) Compare(prev, current float64) float64 {
	if prev == 0 {
		if current == 0 {
			return 0
		}
		return 1.0 // 100% change
	}
	return math.Abs((current - prev) / prev)
}

// FeatureConfig mendefinisikan konfigurasi untuk setiap feature
type FeatureConfig struct {
	Name       string
	Threshold  float64
	Comparator FeatureComparator
}

// SequentialDetector adalah interface umum untuk deteksi anomali sequential
type SequentialDetector interface {
	AddData(data SequentialProvider) (bool, error)
	SetThreshold(featureName string, threshold float64) error
	AddFeatureConfig(config FeatureConfig) error
}

// WindowDetector implementasi sliding window untuk sequential detector
type WindowDetector struct {
	WindowSize     int
	History        []SequentialProvider
	FeatureConfigs map[string]FeatureConfig
}

// NewWindowDetector membuat instance baru WindowDetector
func NewWindowDetector(windowSize int) *WindowDetector {
	return &WindowDetector{
		WindowSize:     windowSize,
		History:        make([]SequentialProvider, 0),
		FeatureConfigs: make(map[string]FeatureConfig),
	}
}

// AddFeatureConfig menambahkan konfigurasi untuk feature baru
func (wd *WindowDetector) AddFeatureConfig(config FeatureConfig) error {
	if config.Comparator == nil {
		config.Comparator = DefaultComparator{}
	}
	wd.FeatureConfigs[config.Name] = config
	return nil
}

// SetThreshold mengubah threshold untuk feature tertentu
func (wd *WindowDetector) SetThreshold(featureName string, threshold float64) error {
	config, exists := wd.FeatureConfigs[featureName]
	if !exists {
		return fmt.Errorf("feature not found: %s", featureName)
	}
	config.Threshold = threshold
	wd.FeatureConfigs[featureName] = config
	return nil
}

// AddData menambahkan data baru dan mendeteksi anomali
func (wd *WindowDetector) AddData(data SequentialProvider) (bool, error) {
	// Tambahkan ke history
	wd.History = append(wd.History, data)

	// Jika belum cukup history, return false
	if len(wd.History) < wd.WindowSize {
		return false, nil
	}

	// Jaga ukuran window
	if len(wd.History) > wd.WindowSize {
		wd.History = wd.History[1:]
	}

	return wd.detectAnomaly()
}

func (wd *WindowDetector) detectAnomaly() (bool, error) {
	currentData := wd.History[len(wd.History)-1]
	prevData := wd.History[len(wd.History)-2]

	// Periksa setiap feature yang terdaftar
	for featureName, config := range wd.FeatureConfigs {
		featureIndex := -1
		// Cari index feature
		for i := 0; i < currentData.GetFeatureCount(); i++ {
			// Asumsikan nama feature disimpan di metadata atau cara lain
			if currentData.GetFeatureName(i) == featureName {
				featureIndex = i
				break
			}
		}

		if featureIndex == -1 {
			return false, fmt.Errorf("feature not found in data: %s", featureName)
		}

		// Bandingkan nilai
		change := config.Comparator.Compare(
			float64(prevData.GetFeatureValue(featureIndex)),
			float64(currentData.GetFeatureValue(featureIndex)),
		)

		if change > config.Threshold {
			return true, nil
		}
	}

	return false, nil
}

// Contoh penggunaan:
// func ExampleUsage() {
// 	// Buat detector
// 	detector := NewWindowDetector(3)

// 	// Tambahkan konfigurasi untuk setiap feature
// 	detector.AddFeatureConfig(FeatureConfig{
// 		Name:       "temperature",
// 		Threshold:  0.2, // 20% change
// 		Comparator: DefaultComparator{},
// 	})

// 	detector.AddFeatureConfig(FeatureConfig{
// 		Name:       "pressure",
// 		Threshold:  0.3, // 30% change
// 		Comparator: DefaultComparator{},
// 	})

// 	detector.AddFeatureConfig(FeatureConfig{
// 		Name:       "humidity",
// 		Threshold:  0.4,
// 		Comparator: CustomComparator{},
// 	})
// }

// // Custom comparator jika perlu
// type CustomComparator struct{}

// func (c CustomComparator) Compare(prev, current float64) float64 {
// 	// Implementasi custom
// 	return math.Abs(current-prev) / 100 // Contoh
// }
