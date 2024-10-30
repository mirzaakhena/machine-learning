package ml

import (
	"encoding/json"
	"os"
)

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
