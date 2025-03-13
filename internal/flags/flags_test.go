package flags

import (
	"flag"
	"os"
	"testing"
)

// TestParseFlags_NoFlags ensures an error occurs when no flags are provided
func TestParseFlags_NoFlags(t *testing.T) {
	os.Args = []string{"cmd"} // No flags provided
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	expectedError := "missing required flag: -c (train or predict). Example: -c train"

	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got: %v", expectedError, err)
	}
}

// TestParseFlags_MultipleMissingTrainFlags ensures errors for multiple missing flags in "train"
func TestParseFlags_MultipleMissingTrainFlags(t *testing.T) {
	os.Args = []string{"cmd", "-c", "train", "-i", "data.csv"} // Missing -t and -o
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	expectedError := "missing required flags for train: -t <target_column>, -o <output.dt>\nExample usage:\n  Train: ./dt -c train -i data.csv -t class -o model.dt\n  Predict: ./dt -c predict -i test.csv -m model.dt -o predictions.csv"

	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got: %v", expectedError, err)
	}
}

// TestParseFlags_MultipleMissingPredictFlags ensures errors for multiple missing flags in "predict"
func TestParseFlags_MultipleMissingPredictFlags(t *testing.T) {
	os.Args = []string{"cmd", "-c", "predict", "-i", "test.csv"} // Missing -m and -o
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	expectedError := "missing required flags for predict: -m <model.dt>, -o <predictions.csv>\nExample usage:\n  Train: ./dt -c train -i data.csv -t class -o model.dt\n  Predict: ./dt -c predict -i test.csv -m model.dt -o predictions.csv"

	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got: %v", expectedError, err)
	}
}

// TestParseFlags_EmptyValues ensures empty values for flags are treated as missing
func TestParseFlags_EmptyValues(t *testing.T) {
	os.Args = []string{"cmd", "-c", "", "-i", "", "-t", "", "-o", ""}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	_, err := ParseFlags()
	expectedError := "missing required flag: -c (train or predict). Example: -c train"

	if err == nil || err.Error() != expectedError {
		t.Errorf("expected error %q, got: %v", expectedError, err)
	}
}

// TestParseFlags_UnexpectedExtraFlags ensures unexpected flags are ignored
func TestParseFlags_UnexpectedExtraFlags(t *testing.T) {
	os.Args = []string{"cmd", "-c", "train", "-i", "data.csv", "-t", "class", "-o", "model.dt", "-x", "unexpected"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError) // ContinueOnError prevents os.Exit

	_, err := ParseFlags()
	if err == nil {
		t.Errorf("expected error due to unrecognized flag, but got nil")
	}
}
