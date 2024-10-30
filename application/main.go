package main

import (
	"fmt"
	"gen"
	"ml"
)

func main() {

	// generateData("./data/large_dataset_2.csv")

	// runMLForTraining("./data/large_dataset.csv", "./data/model.json")

	tree, err := ml.LoadModel("./data/model.json")
	if err != nil {
		panic(err)
	}

	dataset, err := ml.LoadDataFromCSV("./data/large_dataset_2.csv")
	if err != nil {
		panic(err)
	}

	accuracy := ml.GetPredictionAccuration(dataset, tree)
	fmt.Printf("Accuracy: %.2f%%\n", accuracy)

}

func generateData(rawFile string) {
	generator := gen.NewGenerator()
	data := generator.GenerateData(800, 200)
	gen.SaveToCSV(data, rawFile)
}

func runMLForTraining(fileTraining, fileModel string) {

	// Load data
	dataset, err := ml.LoadDataFromCSV(fileTraining)
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
