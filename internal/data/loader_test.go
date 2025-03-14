package data_test

import (
	"os"
	"testing"

	"CodeLens-decision-tree/internal/data"
)

func createTempCSV(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "testdata-*.csv")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func TestGenerateCSVData_ValidData(t *testing.T) {
	csvContent := "Feature1,Feature2,Target\n1.5,A,Yes\n2.3,B,No\n3.7,A,Yes\n"
	filePath, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("failed to create temp CSV: %v", err)
	}
	defer os.Remove(filePath)

	dataset, err := data.GenerateCSVData(filePath, "Target")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if dataset.TotalRows != 3 {
		t.Errorf("expected 3 rows, got: %d", dataset.TotalRows)
	}

	if len(dataset.ColumnNames) != 3 {
		t.Errorf("expected 3 columns, got: %d", len(dataset.ColumnNames))
	}
}

func TestGenerateCSVData_MissingValues(t *testing.T) {
	csvContent := "Feature1,Feature2,Target\n, ,Yes\n2.3,B,\n3.7,A,No\n"
	filePath, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("failed to create temp CSV: %v", err)
	}
	defer os.Remove(filePath)

	dataset, err := data.GenerateCSVData(filePath, "Target")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if dataset.ColumnAttributes["Feature1"].MissingCount != 1 {
		t.Errorf("expected 1 missing value in Feature1, got: %d", dataset.ColumnAttributes["Feature1"].MissingCount)
	}

	if dataset.ColumnAttributes["Target"].MissingCount != 1 {
		t.Errorf("expected 1 missing value in Target, got: %d", dataset.ColumnAttributes["Target"].MissingCount)
	}
}

func TestGenerateCSVData_InvalidTargetColumn(t *testing.T) {
	csvContent := "Feature1,Feature2,Output\n1.5,A,Yes\n2.3,B,No\n"
	filePath, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("failed to create temp CSV: %v", err)
	}
	defer os.Remove(filePath)

	_, err = data.GenerateCSVData(filePath, "Target")
	if err == nil {
		t.Errorf("expected error for missing target column, got nil")
	}
}

func TestGenerateCSVData_EmptyFile(t *testing.T) {
	filePath, err := createTempCSV("")
	if err != nil {
		t.Fatalf("failed to create temp CSV: %v", err)
	}
	defer os.Remove(filePath)

	_, err = data.GenerateCSVData(filePath, "Target")
	if err == nil {
		t.Errorf("expected error for empty CSV, got nil")
	}
}

func TestGenerateCSVData_InconsistentRowLength(t *testing.T) {
	csvContent := "Feature1,Feature2,Target\n1.5,A,Yes\n2.3,B\n3.7,A,No\n"
	filePath, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("failed to create temp CSV: %v", err)
	}
	defer os.Remove(filePath)

	dataset, err := data.GenerateCSVData(filePath, "Target")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if dataset.TotalRows != 2 {
		t.Errorf("expected 2 valid rows, got: %d", dataset.TotalRows)
	}
}

func TestGenerateCSVData_AllValuesMissing(t *testing.T) {
	csvContent := "Feature1,Feature2,Target\n,,\n,,\n,,\n"
	filePath, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("failed to create temp CSV: %v", err)
	}
	defer os.Remove(filePath)

	dataset, err := data.GenerateCSVData(filePath, "Target")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	for _, col := range dataset.ColumnNames {
		if dataset.ColumnAttributes[col].MissingCount != 3 {
			t.Errorf("expected all values missing in %s, got: %d", col, dataset.ColumnAttributes[col].MissingCount)
		}
	}
}

func TestGenerateCSVData_SingleRow(t *testing.T) {
	csvContent := "Feature1,Feature2,Target\n5.0,X,Yes\n"
	filePath, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("failed to create temp CSV: %v", err)
	}
	defer os.Remove(filePath)

	dataset, err := data.GenerateCSVData(filePath, "Target")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if dataset.TotalRows != 1 {
		t.Errorf("expected 1 row, got: %d", dataset.TotalRows)
	}
}
