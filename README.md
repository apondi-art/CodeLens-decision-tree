#   Fast and Scalable Decision Tree

## Introduction
Jam is a fast and scalable C4.5 Decision Tree classifier implemented in Go. It provides a complete toolkit for building, training, pruning, and using decision trees for classification tasks. Designed for high performance, it efficiently handles large datasets using parallelization and memory optimization techniques, making it ideal for real-world machine learning applications.

## why decision tree
Decision trees are a fundamental machine-learning technique used in spam detection, fraud analysis, medical diagnosis, and recommendation systems. This project allows users to:

 Train a decision tree model on large datasets quickly.
    Make predictions on new data with high accuracy.
    Handle both categorical and numerical features efficiently.
    Store trained models in JSON format for easy portability.


### Features
- Complete C4.5 Decision Tree Implementation – Fully implements the C4.5 algorithm for classification tasks.
- Supports Categorical & Numerical Attributes – Handles both types of features seamlessly.
- Efficient Missing Value Handling – Automatically manages incomplete data.
- Post-Pruning to Prevent Overfitting – Improves model generalization.
- Parallelized Tree Construction – Speeds up training using concurrency.
- Memory-Efficient Data Structures – Optimized for large datasets with minimal memory usage.
- Easy-to-Use CLI – Simple command-line interface for training and predictions.
- JSON-Based Model Storage – Supports model serialization and deserialization.
- Comprehensive Evaluation Metrics – Provides accuracy, precision, recall, and F1-score for model assessment.

## Table of Contents

- [Installation and SetUp](#installation-and-setup)
- [Usage Instructions](#usage)
- [Algorithm Details](#algorithm-details)
- [Performance Optimizations](#performance-optimizations)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Contributors](#contributors)
- [Contributing Guidelines](#contributing-guidelines)
- [License](#license)


## Installation and SetUp

### Prerequisites
- Go 1.18 or later
- Git

### Step-by-Step Installation

#### Clone the repository
```bash
git clone https://learn.zone01kisumu.ke/git/quochieng/CodeLens-decision-tree
```

#### Navigate to the project directory
```bash
cd CodeLens-decision-tree
```

### Build the CLI tool
- run the script to create abinary file that in a directory named bin/
```go
chmod +x build.sh
```
### Create a binary
```go
./build.sh
```
### Usage
- Training a Model
```bash
./bin/dt -c train -i cmd/dt/test.csv -t <TargetAtribut> -o model.dt
```

- Making Predictions
```bash
./bin/dt -c predict -i cmd/dt/test.csv -m model.dt -o predictions.csv
```


## Algorithm Details


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
#### Decision Tree Basics
A decision tree is a hierarchical structure that makes sequential decisions based on feature values to arrive at a classification. In this project:

1. Each internal node represents a test on an attribute
2. Each branch represents an outcome of that test
3. Each leaf node represents a class prediction

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

## Error Handling
This implementation incorporates several error handling techniques to ensure robustness and smooth execution:

    Missing or Nil Attributes: The system checks for missing or nil values in required attributes and returns descriptive error messages indicating the specific missing attribute.

    Type Conversion Errors: Before converting data types, the system verifies if the conversion is possible, and handles errors by providing fallback results to prevent crashes.

    Invalid Data or Attribute Types: The system checks for valid data types before processing, skipping invalid data and ensuring smooth operation without interruptions.

    Empty Datasets: Operations on empty datasets are handled by early checks, which return default values (e.g., 0 for entropy) to avoid errors.

    Invalid Splits or Thresholds: The system verifies that valid splits exist for attributes and returns an error if no valid splits are found, preventing further erroneous processing.

    Concurrency Issues: Mutexes are used to ensure safe and synchronized access to shared data during parallel processing, preventing race conditions.

    Unknown Attributes: Descriptive error messages are returned when an unknown attribute is encountered, making it easy to identify and fix dataset issues.

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

## Testing
```bash
# Run all tests
go test ./...

# Run benchmarks
go test -bench=. ./test/benchmark
```
## Contributors

We appreciate the efforts of our contributors. Connect with them on LinkedIn and GitHub:

- [![Stephen Kisengese](https://img.shields.io/badge/LinkedIn-Stephen%20Kisengese-blue?style=flat&logo=linkedin)](https://www.linkedin.com/in/skisenge) 
  [![GitHub](https://img.shields.io/badge/GitHub-skisenge-black?style=flat&logo=github)](https://learn.zone01kisumu.ke/git/skisenge)

- [![Quinter Ochieng](https://img.shields.io/badge/LinkedIn-Quinter%20Ochieng-blue?style=flat&logo=linkedin)](https://www.linkedin.com/in/quinterochieng) 
  [![GitHub](https://img.shields.io/badge/GitHub-QuinterOchieng-black?style=flat&logo=github)](https://learn.zone01kisumu.ke/git/QuinterOchieng)

- [![Samuel Omulo](https://img.shields.io/badge/LinkedIn-Samuel%20Omulo-blue?style=flat&logo=linkedin)](https://www.linkedin.com/in/samuel-omulo-634694261) 
  [![GitHub](https://img.shields.io/badge/GitHub-somulo1-black?style=flat&logo=github)](https://github.com/somulo1)

- [![Antony Musumba](https://img.shields.io/badge/LinkedIn-Antony%20Musumba-blue?style=flat&logo=linkedin)](https://www.linkedin.com/in/antonymusumba) 
  [![GitHub](https://img.shields.io/badge/GitHub-antomusumba-black?style=flat&logo=github)](https://learn.zone01kisumu.ke/git/antomusumba)

### Contributing Guidelines

We welcome contributions from everyone! To get started, please follow these steps:

1. **Fork the repository**: Click the "Fork" button at the top right of the repository page to create your copy.
2. **Create a feature branch**: Use a descriptive name for your branch, e.g., `feature/add-new-feature`.
3. **Make your changes**: Implement your feature or fix the bug, and ensure your code adheres to our coding standards.
4. **Write tests for your changes**: Add tests to verify your changes work as expected.
5. **Run the test suite**: Ensure all tests pass before submitting your changes. Use the command:
6. **Create a pull request**: Submit your changes for review.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.