package model

import "math"

// Dataset represents a collection of instances with attributes
type Dataset struct {
	RowInstances     []map[string]interface{} // Data instances stored as key-value pair(map[column name]rowvalue)
	ColumnAttributes map[string]*Attribute    // Column metadata(column description)
	ColumnNames      []string                 // Ordered list of attribute names
	TargetColumn     string                   // Target column name
	TargetOccurrence map[string]int           // Frequency of each class in dataset
	TotalRows        int                      // Number of rows in dataset
	NonTargetColumns int                      // Number of attributes (excluding target)
}

// NewDataset creates a new dataset from raw data
func NewDataset(data [][]string, headers []string, attrTypes map[string]AttributeType) (*Dataset, error) {
	return &Dataset{}, nil
}

// Clone creates a deep copy of the dataset
func (d *Dataset) Clone() *Dataset {
	return &Dataset{}
}

// FilterByAttributeValue returns a subset of the dataset where attr has specified value
func (d *Dataset) FilterByAttributeValue(attr string, value interface{}) *Dataset {
	return &Dataset{}
}

// FilterByNumericCondition returns subset where numeric attr meets condition (>, <, etc.)
func (d *Dataset) FilterByNumericCondition(attr string, condition string, threshold float64) *Dataset {
	return &Dataset{}
}

// CountClassInstances counts instances per target class.
func (d *Dataset) CountClassInstances() map[string]int {
	classCounts := make(map[string]int)

	for _, instance := range d.RowInstances {
		// Skip if the target column is missing or has a nil value
		if value, exists := instance[d.TargetColumn]; exists && value != nil {
			class := value.(string)
			classCounts[class]++
		}
	}

	return classCounts
}

// GetUniqueValues returns all unique values for a given attribute
func (d *Dataset) GetUniqueValues(attr string) []interface{} {
	uniqueValues := make(map[interface{}]bool)

	// Iterate through all instances and collect unique values.
	for _, instance := range d.RowInstances {
		if value, exists := instance[attr]; exists {
			uniqueValues[value] = true
		}
	}

	// Convert the map keys to a slice.
	result := make([]interface{}, 0, len(uniqueValues))
	for value := range uniqueValues {
		result = append(result, value)
	}

	return result
}

// GetNumericValues returns all values for a numeric attribute as floats
func (d *Dataset) GetNumericValues(attr string) ([]float64, error) {
	return []float64{}, nil
}

// CalculateClassEntropy calculates the entropy of the target attribute
func (d *Dataset) CalculateClassEntropy() float64 {
	if d.TotalRows == 0 {
		return 0
	}

	// Count occurrences of each class
	classCounts := d.CountClassInstances()

	// Calculate entropy
	entropy := 0.0
	for _, count := range classCounts {
		if count > 0 {
			probability := float64(count) / float64(d.TotalRows)
			entropy -= probability * math.Log2(probability)
		}
	}

	return entropy
}

// IsPure returns true if all instances belong to the same class
func (d *Dataset) IsPure() bool {
	if d.TotalRows == 0 {
		return true
	}

	classCounts := d.CountClassInstances()
	return len(classCounts) == 1
}

// GetMajorityClass returns the most frequent class in the dataset
func (d *Dataset) GetMajorityClass() string {
	classCounts := d.CountClassInstances()

	if len(classCounts) == 0 {
		return ""
	}

	majorityClass := ""
	maxCount := 0

	for class, count := range classCounts {
		if count > maxCount {
			maxCount = count
			majorityClass = class
		}
	}

	return majorityClass
}

// func (d *Dataset) SplitByNumericThreshold(attr string, threshold float64) (map[string]*Dataset, error)

// SplitByNumericThreshold splits the dataset based on a numerical threshold.
func (d *Dataset) SplitByNumericThreshold(attr string, threshold float64) (map[interface{}]*Dataset, error) {
	subsets := make(map[interface{}]*Dataset)

	// Create subsets for values <= threshold and values > threshold.
	lessThanOrEqual := &Dataset{
		RowInstances: []map[string]interface{}{},
		TargetColumn: d.TargetColumn,
	}
	greaterThan := &Dataset{
		RowInstances: []map[string]interface{}{},
		TargetColumn: d.TargetColumn,
	}

	// Add instances to the appropriate subset.
	for _, instance := range d.RowInstances {
		value, ok := instance[attr].(float64)
		if !ok {
			continue
		}

		if value <= threshold {
			lessThanOrEqual.RowInstances = append(lessThanOrEqual.RowInstances, instance)
		} else {
			greaterThan.RowInstances = append(greaterThan.RowInstances, instance)
		}
	}

	subsets["<="] = lessThanOrEqual
	subsets[">"] = greaterThan

	return subsets, nil
}

// SplitByCategoricalValue splits the dataset based on the unique values of a categorical attribute.
func (d *Dataset) SplitByCategoricalValue(attr string) (map[interface{}]*Dataset, error) {
	subsets := make(map[interface{}]*Dataset)

	// Get all unique values for the attribute.
	uniqueValues := d.GetUniqueValues(attr)

	// Create a subset for each unique value.
	for _, value := range uniqueValues {
		subset := &Dataset{
			RowInstances: []map[string]interface{}{},
			TargetColumn: d.TargetColumn,
		}

		// Add instances with the current value to the subset.
		for _, instance := range d.RowInstances {
			if instance[attr] == value {
				subset.RowInstances = append(subset.RowInstances, instance)
			}
		}

		subsets[value] = subset
	}

	return subsets, nil
}

// SplitWithMissingValues splits dataset handling missing values using weights
func (d *Dataset) SplitWithMissingValues(attr string, numericThreshold *float64) (map[interface{}]*Dataset, error) {
	return map[interface{}]*Dataset{}, nil
}

// ApplyPreprocessing applies any necessary preprocessing steps to the dataset
func (d *Dataset) ApplyPreprocessing() error {
	return nil
}
