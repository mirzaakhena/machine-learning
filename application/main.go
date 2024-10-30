package main

import (
	"ecu"
	"fmt"
	"gen"
	"ml"
)

func main() {

	// generateData("./data2/large_dataset2.csv")

	// runMLForTraining("./data2/large_dataset.csv", "./data2/model.json")

	getPredictionAccuration("./data2/model.json", "./data2/large_dataset2.csv")

}

func generateData(rawFile string) {
	generator := gen.NewGenerator()
	data := generator.GenerateData(1600, 400)
	err := gen.SaveToCSV(data, rawFile)
	if err != nil {
		panic(err)
	}
}

func runMLForTraining(fileTraining, fileModel string) {

	// Load data
	dataset, err := ml.LoadDataFromCSV(fileTraining, ecu.CreateECUData)
	if err != nil {
		panic(err)
	}

	// Split untuk training dan testing
	trainData, testData := ml.SplitTrainTest(dataset, 0.8)

	// Build tree
	tree := ml.BuildTree(trainData, 0, 5)

	// Test accuracy
	accuracy := ml.GetPredictionAccuration(testData, tree)
	fmt.Printf("Accuracy: %.2f%%\n", accuracy*100)

	ml.SaveModel(tree, fileModel)

	ml.PrintTree(tree, "", true)
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
