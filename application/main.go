package main

import (
	"ecu"
	"fmt"
	"gen"
	"ml"
)

func main() {

	// generateData("./data/large_dataset2.csv")
	// generateModel("./data/large_dataset.csv", "./data/model.json")
	getPredictionAccuration("./data/model.json", "./data/large_dataset2.csv")

}

func generateData(rawFile string) {
	generator := gen.NewGenerator()
	data := generator.GenerateData(1600, 400)
	err := gen.SaveToCSV(data, rawFile)
	if err != nil {
		panic(err)
	}
}

func generateModel(fileTraining, fileModel string) {

	// Load data
	dataset, err := ml.LoadDataFromCSV(fileTraining, ecu.CreateECUData)
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

	tree.SaveModel(fileModel)

	tree.PrintTree()
}

func getPredictionAccuration(fileModel, fileData string) {
	tree, err := ml.LoadModel(fileModel)
	if err != nil {
		panic(err)
	}

	dataset, err := ml.LoadDataFromCSV(fileData, ecu.CreateECUData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Accuracy: %.2f%%\n", tree.GetPredictionAccuration(dataset))
}
