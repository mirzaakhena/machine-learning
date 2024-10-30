# ECU Anomaly Detection

A Go-based machine learning system for detecting anomalies in Electronic Control Unit (ECU) data using decision trees. This project implements a custom decision tree algorithm to identify various vehicle anomalies based on RPM, gear, and speed readings.

## Features

- Custom decision tree implementation for anomaly detection
- Data generation for training and testing
- Model persistence (save/load capabilities)
- Support for various ECU anomaly types:
  - Over-revving
  - Stalling
  - Gear-speed mismatch
  - Neutral with speed
  - RPM too low for speed-gear

## Project Structure

- `generate.go`: Data generation utilities for creating synthetic ECU data
- `ecu.go`: ECU data model and feature provider implementation
- `ml.go`: Core machine learning implementation (decision tree algorithm)
- `main.go`: Application entry point and example usage

## Installation

```bash
# Clone the repository
git clone [repository-url]

# Change into project directory
cd ecu-anomaly-detection

# Install dependencies
go mod tidy
```

## Usage

The project can be used in two main ways:

### 1. Generate Training Data and Create Model

```go
func main() {
    // Generate training data
    generateData("./data/large_dataset.csv")
    
    // Create and save the model
    generateModel("./data/large_dataset.csv", "./data/model.json")
}
```

### 2. Use Existing Model for Predictions

```go
func main() {
    // Generate test data
    generateData("./data/large_dataset2.csv")
    
    // Get prediction accuracy using existing model
    getPredictionAccuration("./data/model.json", "./data/large_dataset2.csv")
}
```

## Data Format

The ECU data contains the following features:

- `rpm`: Engine RPM
- `gear`: Current gear position (0-5, where 0 is neutral)
- `speed`: Vehicle speed in km/h
- `status`: Anomaly status (0 for normal, 1 for anomaly)
- `description`: Type of anomaly or "normal"

## Model Parameters

The decision tree model uses the following parameters:

- Maximum depth: 5 levels
- Training/Testing split: 80/20
- Features considered: RPM, Gear, Speed

## Helper Functions

### Generate Data
```go
func generateData(rawFile string) {
    generator := gen.NewGenerator()
    data := generator.GenerateData(1600, 400) // Generate 1600 normal cases and 400 anomaly cases
    err := gen.SaveToCSV(data, rawFile)
    if err != nil {
        panic(err)
    }
}
```

### Generate Model
```go
func generateModel(fileTraining, fileModel string) {
    // Load and process training data
    dataset, err := ml.LoadDataFromCSV(fileTraining, ecu.CreateECUData)
    if err != nil {
        panic(err)
    }
    
    // Create and save model
    trainData, testData := ml.SplitTrainTest(dataset, 0.8)
    tree := trainData.BuildTree(0, 5)
    tree.SaveModel(fileModel)
}
```

### Get Prediction Accuracy
```go
func getPredictionAccuration(fileModel, fileData string) {
    tree, err := ml.LoadModel(fileModel)
    dataset, err := ml.LoadDataFromCSV(fileData, ecu.CreateECUData)
    fmt.Printf("Accuracy: %.2f%%\n", tree.GetPredictionAccuration(dataset))
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details