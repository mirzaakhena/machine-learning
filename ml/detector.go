package ml

// Model interface
type AnomalyDetector interface {
	Predict(data ECUData) bool
}

// Implementasi model
type DecisionTreeModel struct {
	Root *Node
}

// Implement interface
func (dt *DecisionTreeModel) Predict(data ECUData) bool {
	return Predict(dt.Root, data)
}
