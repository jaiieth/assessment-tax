package calculator

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gocarina/gocsv"
)

type TaxCSVInstance struct {
	File multipart.File
}

func (ti TaxCSVInstance) Validate() error {
	rows, err := gocsv.LazyCSVReader(ti.File).ReadAll()
	if err != nil {
		return fmt.Errorf("wrong csv format")
	}

	header := rows[0]
	expectedHeaders := []string{"totalIncome", "wht", "donation"}

	for i, h := range header {
		if h != expectedHeaders[i] {
			return fmt.Errorf("wrong csv format")
		}
	}

	for _, row := range rows {

		for _, value := range row {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("wrong csv format")
			}
		}
	}
	// Rewind to the beginning of csv, So the `t.File` can be read again
	ti.File.Seek(0, 0)
	return nil
}

func (t TaxCSVInstance) Unmarshal(s interface{}) error {
	if err := gocsv.UnmarshalMultipartFile(&t.File, s); err != nil {
		return err
	}

	return nil
}
