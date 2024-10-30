package ml

import (
	"time"

	"math/rand"
)

// Fungsi untuk split data menjadi training dan testing
func SplitTrainTest(dataset []FeatureProvider, trainRatio float64) (train []FeatureProvider, test []FeatureProvider) {
	// Buat copy dataset dan shuffle
	shuffled := make([]FeatureProvider, len(dataset))
	copy(shuffled, dataset)

	// Shuffle menggunakan Fisher-Yates algorithm
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	// Hitung split point
	splitPoint := int(float64(len(dataset)) * trainRatio)

	// Split data
	train = shuffled[:splitPoint]
	test = shuffled[splitPoint:]

	return
}
