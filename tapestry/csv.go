package tapestry

import (
	"bufio"
	"encoding/csv"
	"os"
)

// CSVData extracts data from csv data files
func CSVData(f string) ([]map[string]string, error) {
	c, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(bufio.NewReader(c))
	headers, err := r.Read()
	if err != nil {
		return nil, err
	}
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var data []map[string]string
	for _, record := range records {
		d := make(map[string]string, len(headers))
		for i := 0; i < len(headers); i++ {
			d[headers[i]] = r[i]
		}
		data = append(data, d)
	}
	return data
}
