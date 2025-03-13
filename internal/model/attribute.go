package model

type AttributeType int

const (
	Categorical AttributeType = iota
	Numerical
	Timestamp
)

// Split represents a decision split point for an attribute
type Split struct {
	Attribute      *Attribute      // The attribute being split on
	Type           string          // "categorical" or "numerical"
	Value          interface{}     // Threshold value for numerical splits
	CategoricalMap map[string]bool // Map of categorical values for categorical splits
	GainRatio      float64         // The gain ratio of this split
}

type Attribute struct {
	Name           string
	Type           AttributeType
	PossibleValues []interface{} // For categorical attributes
	Min            float64       // For numerical attributes
	Max            float64       // For numerical attributes
	MissingCount   int           // Count of missing values
	// Additional fields as needed
}

func (a *Attribute) CalculateGainRatio(dataset *Dataset) float64 {
	return 0.0
}

func (a *Attribute) FindBestSplit(dataset *Dataset) (Split, float64) {
	return Split{}, 0.0
}
