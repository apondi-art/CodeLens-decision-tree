# C4.5 Decision Tree Implementation

## Overview
This project is a scalable, high-performance implementation of the C4.5 decision tree algorithm in Go. It provides a complete toolkit for building, training, pruning, and using decision trees for classification tasks, with special focus on handling large datasets efficiently through parallelization and memory optimization techniques.

## Table of Contents
- [Features](#features)
- [Project Structure](#project-structure)
- [Installation](#installation)
- [Usage](#usage)
- [Algorithm Details](#algorithm-details)
- [Performance Optimizations](#performance-optimizations)
- [API Reference](#api-reference)
- [Development](#development)
- [Testing](#testing)
- [Contributing Guidelines](#contributing-guidelines)
- [License](#license)

## Features
- Full implementation of the C4.5 decision tree algorithm
- Support for both categorical and numerical attributes
- Efficient handling of missing values
- Post-pruning capabilities to prevent overfitting
- Parallelized tree building for improved performance
- Memory-efficient data structures for large datasets
- Command-line interface for training and prediction
- Model serialization and deserialization (JSON format)
- Comprehensive evaluation metrics

## Project Structure
```
decision-tree/
├── cmd/
│   └── dt/
│       ├── main.go                 # CLI entry point
│       └── main_test.go            # Test file for CLI commands
├── internal/
│   ├── data/
│   │   ├── loader.go               # Data loading and preprocessing
│   │   ├── parser.go               # CSV parsing logic
│   │   ├── validator.go            # Input validation
│   │   └── data_test.go            # Test file for data handling
│   ├── model/
│   │   ├── attribute.go            # Attribute representation
│   │   ├── dataset.go              # Dataset structure
│   │   ├── node.go                 # Decision tree node
│   │   ├── tree.go                 # Decision tree model
│   │   └── model_test.go            # Test file for model components
│   ├── algorithm/
│   │   ├── entropy.go              # Entropy and information gain calculations
│   │   ├── splitter.go             # Data splitting logic
│   │   ├── c45.go                  # Core C4.5 implementation
│   │   ├── pruner.go               # Tree pruning
│   │   └── algorithm_test.go        # Test file for algorithms
│   ├── serialization/
│   │   ├── json_serializer.go      # Model serialization to JSON
│   │   ├── deserializer.go         # Model loading from JSON
│   │   └── serialization_test.go     # Test file for serialization
│   └── parallel/
│       ├── worker_pool.go          # Parallelization utilities
│       └── parallel_test.go         # Test file for parallelization
├── pkg/
│   └── metrics/
│       ├── evaluation.go           # Performance metrics calculation
│       └── metrics_test.go          # Test file for metrics
└── README.md                       # Documentation
└── LICENSE                         # License file
```

## Installation

### Prerequisites
- Go 1.18 or later
- Git

### Step-by-Step Installation
```bash
# Clone the repository
https://learn.zone01kisumu.ke/git/quochieng/CodeLens-decision-tree

# Navigate to the project directory
CodeLens-decision-tree

# Build the CLI tool
go build -o bin/dt ./cmd/dt
```

## Usage

### Training a Model
```bash
# Train a model using a CSV file
./bin/dt train --input data.csv --target class --output model.json
```

### Making Predictions
```bash
# Use a trained model to make predictions
./bin/dt predict --input test.csv --model model.json --output predictions.csv
```

### Advanced Options
```bash
# Train with pruning enabled and parallel processing
./bin/dt train --input data.csv --target class --output model.json --prune --validate validation.csv --workers 4

# Train with custom parameters
./bin/dt train --input data.csv --target class --output model.json --max-depth 10 --min-samples 5
```

## Algorithm Details

### Decision Tree Basics
A decision tree is a hierarchical structure that makes sequential decisions based on feature values to arrive at a classification. In this project:

1. Each internal node represents a test on an attribute
2. Each branch represents an outcome of that test
3. Each leaf node represents a class prediction

### C4.5 Algorithm Specifics
This implementation follows the C4.5 algorithm with these key components:

#### 1. Entropy and Information Gain
```go
// Calculate entropy of a dataset based on class distribution
func CalculateEntropy(dataset *model.Dataset, targetAttr string) float64

// Calculate information gain for an attribute
func CalculateInformationGain(dataset *model.Dataset, attr *model.Attribute, targetAttr string) float64

// Calculate gain ratio (C4.5 specific)
func CalculateGainRatio(dataset *model.Dataset, attr *model.Attribute, targetAttr string) float64
```

#### 2. Data Splitting
```go
// Split dataset based on attribute values
func SplitDataset(dataset *model.Dataset, attr *model.Attribute, value interface{}) (*model.Dataset, error)

// Find best split point for numerical attributes
func FindBestNumericalSplit(dataset *model.Dataset, attr *model.Attribute, targetAttr string) (float64, float64)

// Handle missing values during prediction
func DistributeInstance(instance map[string]interface{}, children map[interface{}]*model.Node) map[*model.Node]float64
```

#### 3. Tree Building
```go
// Build decision tree recursively
func BuildTree(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string, depth int, maxDepth int) (*model.Node, error)

// Select best attribute for splitting
func SelectBestAttribute(dataset *model.Dataset, attributes []*model.Attribute, targetAttr string) (*model.Attribute, float64)

// Determine majority class for leaf nodes
func MajorityClass(dataset *model.Dataset, targetAttr string) string
```

#### 4. Tree Pruning
```go
// Prune the decision tree using a validation set
func PruneTree(root *model.Node, validationSet *model.Dataset, targetAttr string) *model.Node

// Estimate error for a node using a validation set
func EstimateError(node *model.Node, dataset *model.Dataset, targetAttr string) float64
```

## Performance Optimizations

### Memory Efficiency
- Use dataset pointers rather than copying data
- Implement efficient data structures (e.g., sparse representations for categorical data)

### Algorithmic Optimizations
- Cache entropy calculations for repeated operations
- Pre-sort numerical attributes to avoid redundant sorting
- Early stopping for low-gain splits

### Parallelization Strategy
- Parallelize node construction for different branches using worker pools
- Implement work-stealing for load balancing
- Use channels for coordination between workers

### Scaling for Large Datasets
- Batch processing for datasets that don't fit in memory
- Sampling techniques for initial attribute evaluation
- Support for distributed processing (optional advanced feature)

## API Reference

### Data Handling
```go
// Load data from CSV file
func LoadCSV(path string) (*model.Dataset, error)

// Infer data types from CSV content
func InferDataTypes(data [][]string) (map[string]string, error)

// Handle missing values in the dataset
func HandleMissingValues(dataset *model.Dataset) error
```

### Core Model Components
```go
// Attribute types
type AttributeType int
const (
    Categorical AttributeType = iota
    Numerical
    Timestamp
)

// Decision tree node
type Node struct {
    Attribute       *Attribute
    SplitValue      interface{}
    Children        map[interface{}]*Node
    IsLeaf          bool
    PredictedClass  string
}

// Decision tree
type DecisionTree struct {
    Root           *Node
    Attributes     []*Attribute
    TargetAttr     string
    AttributeTypes map[string]AttributeType
}
```

### Training and Prediction
```go
// Train a decision tree
func (t *DecisionTree) Train(dataset *Dataset) error

// Make a prediction for a single instance
func (t *DecisionTree) Predict(instance map[string]interface{}) string

// Make predictions for multiple instances
func (t *DecisionTree) BatchPredict(dataset *Dataset) []string
```

### Serialization
```go
// Save model to file
func SerializeTree(tree *model.DecisionTree, path string) error

// Load model from file
func DeserializeTree(path string) (*model.DecisionTree, error)
```

## Development
Collaborators
This project is maintained by the following contributors:


- [Stephen Kisengese](https://learn.zone01kisumu.ke/git/skisenge)
- [Quinter Ochieng](https://learn.zone01kisumu.ke/git/QuinterOchieng)
- [Samuel Omulo](https://learn.zone01kisumu.ke/git/somulo)
- [Antony Musumba](https://learn.zone01kisumu.ke/git/antomusumba)


    ### Contributing Guidelines

    We welcome contributions from everyone! To get started, please follow these steps:

1. **Fork the repository**: Click the "Fork" button at the top right of the repository page to create your copy.
2. **Create a feature branch**: Use a descriptive name for your branch, e.g., `feature/add-new-feature`.
3. **Make your changes**: Implement your feature or fix the bug, and ensure your code adheres to our coding standards.
4. **Write tests for your changes**: Add tests to verify your changes work as expected.
5. **Run the test suite**: Ensure all tests pass before submitting your changes. Use the command:
6. **Create a pull request**: Submit your changes for review.

   ```bash
   go test ./...

### Coding Standards
- Follow Go best practices and code style
- Write clear, concise comments
- Maintain high test coverage
- Optimize for performance and memory efficiency

## Testing
```bash
# Run all tests
go test ./...

# Run benchmarks
go test -bench=. ./test/benchmark
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.