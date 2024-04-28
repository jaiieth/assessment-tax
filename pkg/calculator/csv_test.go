package calculator_test

import (
	"os"
	"testing"

	calc "github.com/jaiieth/assessment-tax/pkg/calculator"
	"github.com/stretchr/testify/assert"
)

func NewCSV(s string) (*os.File, error) {
	// Create a temporary CSV file with expected headers and non-empty values
	tempFile, err := os.CreateTemp("", "taxes.csv")
	if err != nil {
		return nil, err
	}
	_, err = tempFile.WriteString(s)

	if err != nil {
		return nil, err
	}
	// Flush the buffer to ensure data is written immediately
	err = tempFile.Sync()
	if err != nil {
		return nil, err
	}

	// Move file pointer to the beginning of the file
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

// Valid CSV file with expected headers and non-empty values returns nil error
func TestValidateValidCSVFile(t *testing.T) {
	csvData := `totalIncome,wht,donation
      1000,200,50
      2000,400,100
      3000,600,150`

	tempFile, err := NewCSV(csvData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	defer os.Remove(tempFile.Name())

	taxCSV := calc.TaxCSVInstance{
		File: tempFile,
	}

	err = taxCSV.Validate()

	assert.NoError(t, err, "Expected nil error, but got: %v", err)
}

// Invalid CSV file with missing headers returns error
func TestValidateInvalidCSVFile(t *testing.T) {
	csvData := "Invalid CSV"

	tempFile, err := NewCSV(csvData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	defer os.Remove(tempFile.Name())

	taxCSV := calc.TaxCSVInstance{
		File: tempFile,
	}

	err = taxCSV.Validate()

	assert.Error(t, err, "Expected non-nil error, but got nil")
}

// Valid CSV file with extra headers and non-empty values returns error
func TestValidateValidCSVFileWithExtraHeaders(t *testing.T) {
	csvData := `totalIncome,wht,donation,extraHeader,
  1000,200,50`

	tempFile, err := NewCSV(csvData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer os.Remove(tempFile.Name())

	taxCSV := calc.TaxCSVInstance{
		File: tempFile,
	}

	err = taxCSV.Validate()
	assert.Error(t, err, "Expected non-nil error, but got nil")
}

// Valid CSV file with expected headers and empty values should return an error
func TestValidateValidCSVFileWithExpectedHeadersAndEmptyValues(t *testing.T) {
	csvData := `totalIncome,wht,donation
                  ,,`
	tempFile, err := NewCSV(csvData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer os.Remove(tempFile.Name())

	taxCSV := calc.TaxCSVInstance{
		File: tempFile,
	}

	err = taxCSV.Validate()
	assert.Error(t, err, "Expected non-nil error, but got nil")
}

// validate non-numeric values should return an error
func TestUnmarshalCSVWithNonNumericValues(t *testing.T) {
	csvData := `totalIncome,wht,donation
                  100,200,300
                  abc,def,ghi`

	tempFile, err := NewCSV(csvData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer os.Remove(tempFile.Name())

	taxCSV := calc.TaxCSVInstance{
		File: tempFile,
	}

	err = taxCSV.Unmarshal(&[]calc.TaxCSV{})
	assert.Error(t, err, "Expected non-nil error, but got nil")
}

func TestUnmarshalValidCSV(t *testing.T) {
	csvData := `totalIncome,wht,donation
                  100,200,300
                  200,0,0`

	tempFile, err := NewCSV(csvData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	defer os.Remove(tempFile.Name())

	taxCSV := calc.TaxCSVInstance{
		File: tempFile,
	}

	err = taxCSV.Unmarshal(&[]calc.TaxCSV{})
	assert.NoError(t, err, "Expected nil error, but got: %v", err)
}
