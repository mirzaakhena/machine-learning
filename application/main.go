package main

import (
	"bufio"
	"ecu"
	"fmt"
	"gen"
	"ml"
	"os"
	"strconv"
	"strings"
)

func main() {

	// generate to create model
	// generateData("./data/large_dataset.csv")
	// generateModel("./data/large_dataset.csv", "./data/model.json")

	// generate to predict data
	// generateData("./data/large_dataset2.csv")
	// getPredictionAccuration("./data/model.json", "./data/large_dataset2.csv")

	sequenceAnomalyDetection()
}

func sequenceAnomalyDetection() {
	detector := ecu.GetSequentialAnomalyDetector()

	anomalyCount := 0
	CSVFileReader("./data/data_sequential.csv", func(index int, data string) {

		str := ParseCSVLine(data)

		rpm, err := strconv.Atoi(str[2])
		if err != nil {
			return
		}

		gear, err := strconv.Atoi(str[3])
		if err != nil {
			return
		}

		speed, err := strconv.Atoi(str[4])
		if err != nil {
			return
		}

		ecuData := ecu.ECUData{
			RPM:   rpm,
			Gear:  gear,
			Speed: speed,
		}

		isAnomaly, err := detector.AddData(ecuData)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			return
		}

		if isAnomaly {
			fmt.Printf("has anomaly at line %d\n", index)
			anomalyCount++
		}

	})

	if anomalyCount == 0 {
		fmt.Printf("no anomaly found")
	}
}

type individualDetector struct{}

func (individualDetector) generateData(destinationFileData string) {
	generator := gen.NewGenerator()
	data := generator.GenerateData(1600, 400)
	err := gen.SaveToCSV(data, destinationFileData)
	if err != nil {
		panic(err)
	}
}

func (individualDetector) generateModel(sourceFileTraining, destinationFileModel string) {

	// Load data
	dataset, err := ml.LoadDataFromCSV(sourceFileTraining, ecu.CreateECUData)
	if err != nil {
		panic(err)
	}

	// Split untuk training dan testing
	trainData, testData := ml.SplitTrainTest(dataset, 0.8)

	// Build tree
	tree := trainData.BuildTree(0, 5)

	// Test accuracy
	accuracy := tree.GetPredictionAccuration(testData)
	fmt.Printf("Accuracy: %.2f%%\n", accuracy*100)

	tree.SaveModel(destinationFileModel)

	tree.PrintTree()
}

func (individualDetector) getPredictionAccuration(sourceFileModel, sourceFileData string) {
	tree, err := ml.LoadModel(sourceFileModel)
	if err != nil {
		panic(err)
	}

	dataset, err := ml.LoadDataFromCSV(sourceFileData, ecu.CreateECUData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Accuracy: %.2f%%\n", tree.GetPredictionAccuration(dataset))
}

func CSVFileReader(filename string, readline func(index int, data string)) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	index := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			readline(index, line)
			index++
		}
	}

	return scanner.Err()
}

func ParseCSVLine(line string) []string {
	// Split dengan koma
	parts := strings.Split(line, ",")

	// Jika kurang dari 5 field, return semua
	if len(parts) <= 5 {
		return parts
	}

	// Ambil 5 field pertama
	result := parts[:5]

	// Gabungkan sisa field sebagai field terakhir
	remainder := strings.Join(parts[5:], ",")
	result = append(result, strings.TrimSpace(remainder))

	return result
}
