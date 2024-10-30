package ml

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"math/rand"
)

// type DataWithAnomalyField interface {
// 	IsAnomaly() bool
// }

type ECUData struct {
	isAttack bool // true jika status=1
}

func (e ECUData) IsAnomaly() bool {
	return e.isAttack
}

// Struktur data untuk tree
type Node struct {
	Feature    int // 0=RPM, 1=Gear, 2=Speed
	Threshold  int
	Left       *Node
	Right      *Node
	IsLeaf     bool
	Prediction bool
}

// Load data dari CSV
func LoadDataFromCSV(filename string) ([]ECUData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var dataset []ECUData
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		rpm, _ := strconv.Atoi(record[0])
		gear, _ := strconv.Atoi(record[1])
		speed, _ := strconv.Atoi(record[2])
		status, _ := strconv.Atoi(record[3])

		data := ECUData{
			RPM:      rpm,
			Gear:     gear,
			Speed:    speed,
			IsAttack: status == 1,
		}
		dataset = append(dataset, data)
	}

	return dataset, nil
}

// Fungsi untuk mendapatkan nilai feature berdasarkan index
func getFeatureValue(data ECUData, feature int) int {
	switch feature {
	case 0:
		return data.RPM
	case 1:
		return data.Gear
	case 2:
		return data.Speed
	default:
		return 0
	}
}

// Hitung proporsi attack dalam dataset
func calculateAttackProportion(dataset []ECUData) float64 {
	attackCount := 0
	for _, data := range dataset {
		if data.IsAttack {
			attackCount++
		}
	}
	return float64(attackCount) / float64(len(dataset))
}

// Fungsi untuk membangun tree
func BuildTree(dataset []ECUData, depth int, maxDepth int) *Node {
	// Log indentasi untuk visualisasi kedalaman
	indent := strings.Repeat("  ", depth)

	fmt.Printf("%sBuildTree: depth=%d, dataset size=%d\n", indent, depth, len(dataset))

	// Base case: jika dataset kosong
	if len(dataset) == 0 {
		fmt.Printf("%s└─ Empty dataset, returning leaf node (prediction=false)\n", indent)
		return &Node{
			IsLeaf:     true,
			Prediction: false,
		}
	}

	// Base case: jika sudah mencapai max depth atau semua data punya label sama
	attackProp := calculateAttackProportion(dataset)
	fmt.Printf("%s├─ Attack proportion: %.2f%%\n", indent, attackProp*100)

	if depth >= maxDepth {
		fmt.Printf("%s└─ Reached max depth, returning leaf node (prediction=%v)\n",
			indent, attackProp >= 0.5)
		return &Node{
			IsLeaf:     true,
			Prediction: attackProp >= 0.5,
		}
	}

	if attackProp == 0 || attackProp == 1 {
		fmt.Printf("%s└─ Pure node (all %s), returning leaf node\n",
			indent, map[bool]string{true: "attacks", false: "normal"}[attackProp == 1])
		return &Node{
			IsLeaf:     true,
			Prediction: attackProp >= 0.5,
		}
	}

	// Cari split terbaik
	bestFeature, bestThreshold, bestGain := findBestSplit(dataset)
	fmt.Printf("%s├─ Best split: feature=%s, threshold=%d, gain=%.4f\n",
		indent,
		map[int]string{0: "RPM", 1: "Gear", 2: "Speed"}[bestFeature],
		bestThreshold,
		bestGain)

	if bestGain == 0 {
		fmt.Printf("%s└─ No gain from splitting, returning leaf node (prediction=%v)\n",
			indent, attackProp >= 0.5)
		return &Node{
			IsLeaf:     true,
			Prediction: attackProp >= 0.5,
		}
	}

	// Split dataset
	leftData, rightData := splitDataset(dataset, bestFeature, bestThreshold)
	fmt.Printf("%s├─ Split result: left=%d samples, right=%d samples\n",
		indent, len(leftData), len(rightData))

	// Buat node
	node := &Node{
		Feature:   bestFeature,
		Threshold: bestThreshold,
		IsLeaf:    false,
	}

	// Rekursif untuk left dan right child
	fmt.Printf("%s├─ Building left subtree...\n", indent)
	node.Left = BuildTree(leftData, depth+1, maxDepth)

	fmt.Printf("%s└─ Building right subtree...\n", indent)
	node.Right = BuildTree(rightData, depth+1, maxDepth)

	return node
}

// Fungsi untuk mencari split terbaik
func findBestSplit(dataset []ECUData) (bestFeature int, bestThreshold int, bestGain float64) {
	bestGain = 0

	// Untuk setiap feature (RPM, Gear, Speed)
	for feature := 0; feature < 3; feature++ {
		// Cari nilai unik untuk feature ini
		values := make(map[int]bool)
		for _, data := range dataset {
			values[getFeatureValue(data, feature)] = true
		}

		// Untuk setiap nilai possible threshold
		for threshold := range values {
			gain := calculateInformationGain(dataset, feature, threshold)
			if gain > bestGain {
				bestGain = gain
				bestFeature = feature
				bestThreshold = threshold
			}
		}
	}

	return
}

// Fungsi untuk menghitung information gain
func calculateInformationGain(dataset []ECUData, feature int, threshold int) float64 {
	parentEntropy := calculateEntropy(dataset)

	leftData, rightData := splitDataset(dataset, feature, threshold)

	// Hitung weighted entropy setelah split
	leftWeight := float64(len(leftData)) / float64(len(dataset))
	rightWeight := float64(len(rightData)) / float64(len(dataset))

	leftEntropy := calculateEntropy(leftData)
	rightEntropy := calculateEntropy(rightData)

	weightedEntropy := leftWeight*leftEntropy + rightWeight*rightEntropy

	return parentEntropy - weightedEntropy
}

// Fungsi untuk menghitung entropy
func calculateEntropy(dataset []ECUData) float64 {
	if len(dataset) == 0 {
		return 0
	}

	attackProp := calculateAttackProportion(dataset)
	if attackProp == 0 || attackProp == 1 {
		return 0
	}

	return -attackProp*log2(attackProp) - (1-attackProp)*log2(1-attackProp)
}

func log2(x float64) float64 {
	return math.Log(x) / math.Log(2)
}

// Fungsi untuk split dataset
func splitDataset(dataset []ECUData, feature int, threshold int) (left []ECUData, right []ECUData) {
	for _, data := range dataset {
		if getFeatureValue(data, feature) <= threshold {
			left = append(left, data)
		} else {
			right = append(right, data)
		}
	}
	return
}

// Fungsi untuk prediksi
func (node *Node) Predict(data ECUData) bool {
	if node.IsLeaf {
		return node.Prediction
	}

	if getFeatureValue(data, node.Feature) <= node.Threshold {
		return node.Left.Predict(data)
	}
	return node.Right.Predict(data)
}

// Fungsi untuk split data menjadi training dan testing
func SplitTrainTest(dataset []ECUData, trainRatio float64) (train []ECUData, test []ECUData) {
	// Buat copy dataset dan shuffle
	shuffled := make([]ECUData, len(dataset))
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

// Fungsi untuk save model ke file
func SaveModel(node *Node, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(node)
}

// Fungsi untuk load model dari file
func LoadModel(filename string) (*Node, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root Node
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&root); err != nil {
		return nil, err
	}

	return &root, nil
}

func PrintTree(node *Node, prefix string, isLeft bool) {
	if node == nil {
		return
	}

	// Karakter untuk visualisasi tree
	var connector string
	if isLeft {
		connector = "├──"
	} else {
		connector = "└──"
	}

	if node.IsLeaf {
		fmt.Printf("%s%s [Leaf: %v]\n", prefix, connector,
			map[bool]string{true: "Attack", false: "Normal"}[node.Prediction])
		return
	}

	fmt.Printf("%s%s [%s <= %d]\n", prefix, connector,
		map[int]string{0: "RPM", 1: "Gear", 2: "Speed"}[node.Feature],
		node.Threshold)

	// Tentukan prefix untuk child nodes
	newPrefix := prefix
	if isLeft {
		newPrefix += "│   "
	} else {
		newPrefix += "    "
	}

	PrintTree(node.Left, newPrefix, true)
	PrintTree(node.Right, newPrefix, false)
}

func GetPredictionAccuration(testData []ECUData, tree *Node) float64 {

	correct := 0
	for _, data := range testData {
		prediction := tree.Predict(data)
		if prediction == data.IsAttack {
			correct++
		}
	}

	accuracy := float64(correct) / float64(len(testData))
	return accuracy * 100
}
