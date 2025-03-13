package model

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
	return []interface{}{}
}

// GetNumericValues returns all values for a numeric attribute as floats
func (d *Dataset) GetNumericValues(attr string) ([]float64, error) {
	return []float64{}, nil
}

// CalculateClassEntropy calculates the entropy of the target attribute
func (d *Dataset) CalculateClassEntropy() float64 {
	return 0
}

// IsPure returns true if all instances belong to the same class
func (d *Dataset) IsPure() bool {
	return false
}

// GetMajorityClass returns the most frequent class in the dataset
func (d *Dataset) GetMajorityClass() string {
	return ""
}

// SplitByNumericThreshold splits dataset based on numeric attribute threshold
func (d *Dataset) SplitByNumericThreshold(attr string, threshold float64) (map[string]*Dataset, error) {
	return map[string]*Dataset{}, nil
}

// SplitByCategoricalValue splits dataset based on categorical attribute values
func (d *Dataset) SplitByCategoricalValue(attr string) (map[interface{}]*Dataset, error) {
	return map[interface{}]*Dataset{}, nil
}

// SplitWithMissingValues splits dataset handling missing values using weights
func (d *Dataset) SplitWithMissingValues(attr string, numericThreshold *float64) (map[interface{}]*Dataset, error) {
	return map[interface{}]*Dataset{}, nil
}

// ApplyPreprocessing applies any necessary preprocessing steps to the dataset
func (d *Dataset) ApplyPreprocessing() error {
	return nil
}
