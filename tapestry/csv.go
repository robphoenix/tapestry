package tapestry

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
)

// CSVData extracts data from csv data files
func CSVData(f string) ([]map[string]string, error) {
	c, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", f, err)
	}
	defer c.Close()
	r := csv.NewReader(bufio.NewReader(c))
	headers, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read headers %s: %v", f, err)
	}
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records %s: %v", f, err)
	}
	var data []map[string]string
	for _, record := range records {
		d := make(map[string]string, len(headers))
		for i := 0; i < len(headers); i++ {
			d[headers[i]] = record[i]
		}
		data = append(data, d)
	}
	return data, nil
}
