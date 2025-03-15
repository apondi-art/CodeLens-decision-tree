package model

import (
	"fmt"
	"math"
	"sync"
)

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

// CountClassInstances counts instances per target class.
func (d *Dataset) CountClassInstances() map[string]int {
	// Return cached result if available
	if d.TargetOccurrence != nil && d.TargetColumn != "" {
		return d.TargetOccurrence
	}

	// Estimate capacity based on typical number of classes
	estimatedClasses := 10
	if d.TargetOccurrence != nil {
		estimatedClasses = len(d.TargetOccurrence)
	}

	classCounts := make(map[string]int, estimatedClasses)

	for _, instance := range d.RowInstances {
		// Skip if the target column is missing or has a nil value
		if value, exists := instance[d.TargetColumn]; exists && value != nil {
			if class, ok := value.(string); ok {
				classCounts[class]++
			}
		}
	}

	// Cache the result for future use
	if d.TargetColumn != "" {
		d.TargetOccurrence = classCounts
	}

	return classCounts
}

// GetUniqueValues returns all unique values for a given attribute
func (d *Dataset) GetUniqueValues(attr string) []interface{} {
	// For small datasets, just use the standard approach
	if len(d.RowInstances) < 1000 {
		uniqueValues := make(map[interface{}]bool)

		// Iterate through all instances and collect unique values
		for _, instance := range d.RowInstances {
			if value, exists := instance[attr]; exists && value != nil {
				uniqueValues[value] = true
			}
		}

		// Convert the map keys to a slice with pre-allocation
		result := make([]interface{}, 0, len(uniqueValues))
		for value := range uniqueValues {
			result = append(result, value)
		}

		return result
	}

	// For large datasets, use concurrent processing
	const chunkSize = 500
	numWorkers := (len(d.RowInstances) + chunkSize - 1) / chunkSize

	// Create chunks of work
	var wg sync.WaitGroup
	resultChan := make(chan map[interface{}]bool, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(startIdx, endIdx int) {
			defer wg.Done()
			localUnique := make(map[interface{}]bool)

			// Process chunk
			for j := startIdx; j < endIdx && j < len(d.RowInstances); j++ {
				if value, exists := d.RowInstances[j][attr]; exists && value != nil {
					localUnique[value] = true
				}
			}

			resultChan <- localUnique
		}(i*chunkSize, (i+1)*chunkSize)
	}

	// Close channel when all workers complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Merge results
	mergedUnique := make(map[interface{}]bool)
	for localUnique := range resultChan {
		for value := range localUnique {
			mergedUnique[value] = true
		}
	}

	// Convert to slice
	result := make([]interface{}, 0, len(mergedUnique))
	for value := range mergedUnique {
		result = append(result, value)
	}

	return result
}

// GetNumericValues extracts all the numeric values for a specified attribute.
// It returns a slice of float64 values for the attribute, or an error if the attribute isn't found or the data isn't numeric.
func (d *Dataset) GetNumericValues(attr string) ([]float64, error) {
	// Create a slice to store numeric values
	var numericValues []float64

	// Iterate over each instance in the dataset
	for _, instance := range d.RowInstances {
		// Check if the attribute exists in the instance
		attrValue, exists := instance[attr]
		if !exists {
			return nil, fmt.Errorf("attribute '%s' not found in the dataset", attr)
		}

		// Ensure the attribute is of type float64 (numeric)
		if numVal, ok := attrValue.(float64); ok {
			numericValues = append(numericValues, numVal)
		} else {
			// If the value isn't numeric, we skip this instance
			continue
		}
	}

	// Return the extracted numeric values
	return numericValues, nil
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
	// Use cached target occurrence if available
	var classCounts map[string]int
	if d.TargetOccurrence != nil {
		classCounts = d.TargetOccurrence
	} else {
		classCounts = d.CountClassInstances()
	}

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

// SplitByNumericThreshold splits the dataset based on a numerical threshold.
func (d *Dataset) SplitByNumericThreshold(attr string, threshold float64) (map[interface{}]*Dataset, error) {
	// Pre-allocate with expected size
	subsets := make(map[interface{}]*Dataset, 2)

	// Count instances for each subset to pre-allocate
	leCount := 0
	gtCount := 0

	for _, instance := range d.RowInstances {
		value, ok := instance[attr].(float64)
		if !ok {
			continue
		}

		if value <= threshold {
			leCount++
		} else {
			gtCount++
		}
	}

	// Create subsets with proper capacity
	lessThanOrEqual := &Dataset{
		RowInstances:     make([]map[string]interface{}, 0, leCount),
		TargetColumn:     d.TargetColumn,
		ColumnAttributes: d.ColumnAttributes,
		ColumnNames:      d.ColumnNames,
		NonTargetColumns: d.NonTargetColumns,
	}

	greaterThan := &Dataset{
		RowInstances:     make([]map[string]interface{}, 0, gtCount),
		TargetColumn:     d.TargetColumn,
		ColumnAttributes: d.ColumnAttributes,
		ColumnNames:      d.ColumnNames,
		NonTargetColumns: d.NonTargetColumns,
	}

	// Add instances to the appropriate subset
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

	// Set total rows for each subset
	lessThanOrEqual.TotalRows = len(lessThanOrEqual.RowInstances)
	greaterThan.TotalRows = len(greaterThan.RowInstances)

	subsets["<="] = lessThanOrEqual
	subsets[">"] = greaterThan

	return subsets, nil
}

// SplitByCategoricalValue splits the dataset based on the unique values of a categorical attribute.
func (d *Dataset) SplitByCategoricalValue(attr string) (map[interface{}]*Dataset, error) {
	// Get unique values first to determine how many subsets we'll need
	uniqueValues := d.GetUniqueValues(attr)
	if len(uniqueValues) == 0 {
		return nil, fmt.Errorf("no unique values found for attribute: %s", attr)
	}

	// For smaller datasets or few unique values, use the direct approach
	if len(d.RowInstances) < 1000 || len(uniqueValues) < 5 {
		// Pre-allocate the map with the correct capacity
		subsets := make(map[interface{}]*Dataset, len(uniqueValues))

		// Count instances per value to pre-allocate arrays
		valueCounts := make(map[interface{}]int, len(uniqueValues))
		for _, instance := range d.RowInstances {
			if value, exists := instance[attr]; exists && value != nil {
				valueCounts[value]++
			}
		}

		// Create a subset for each unique value with appropriate capacity
		for _, value := range uniqueValues {
			count := valueCounts[value]
			subset := &Dataset{
				RowInstances:     make([]map[string]interface{}, 0, count),
				TargetColumn:     d.TargetColumn,
				ColumnAttributes: d.ColumnAttributes,
				ColumnNames:      d.ColumnNames,
				NonTargetColumns: d.NonTargetColumns,
			}

			// Add instances with matching value
			for _, instance := range d.RowInstances {
				if instance[attr] == value {
					subset.RowInstances = append(subset.RowInstances, instance)
				}
			}

			subset.TotalRows = len(subset.RowInstances)
			subsets[value] = subset
		}

		return subsets, nil
	}

	// For larger datasets with many values, use concurrent approach
	subsets := make(map[interface{}]*Dataset, len(uniqueValues))
	var mutex sync.Mutex

	// Create empty datasets for each value
	for _, value := range uniqueValues {
		subsets[value] = &Dataset{
			RowInstances:     make([]map[string]interface{}, 0),
			TargetColumn:     d.TargetColumn,
			ColumnAttributes: d.ColumnAttributes,
			ColumnNames:      d.ColumnNames,
			NonTargetColumns: d.NonTargetColumns,
		}
	}

	// Process in parallel chunks
	const chunkSize = 500
	numWorkers := (len(d.RowInstances) + chunkSize - 1) / chunkSize
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(startIdx, endIdx int) {
			defer wg.Done()

			// Local maps to reduce lock contention
			localSubsets := make(map[interface{}][]map[string]interface{}, len(uniqueValues))
			for _, value := range uniqueValues {
				localSubsets[value] = make([]map[string]interface{}, 0)
			}

			// Process chunk
			for j := startIdx; j < endIdx && j < len(d.RowInstances); j++ {
				instance := d.RowInstances[j]
				if value, exists := instance[attr]; exists && value != nil {
					if _, found := localSubsets[value]; found {
						localSubsets[value] = append(localSubsets[value], instance)
					}
				}
			}

			// Merge results with lock
			mutex.Lock()
			defer mutex.Unlock()
			for value, instances := range localSubsets {
				if len(instances) > 0 {
					subsets[value].RowInstances = append(subsets[value].RowInstances, instances...)
				}
			}
		}(i*chunkSize, (i+1)*chunkSize)
	}

	wg.Wait()

	// Set total rows for each subset
	for _, subset := range subsets {
		subset.TotalRows = len(subset.RowInstances)
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
