package excel

import (
	"encoding/csv"
	"fmt"
	"os"
)

func Export(filePath string, data [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to CSV: %w", err)
		}
	}

	return nil
}
