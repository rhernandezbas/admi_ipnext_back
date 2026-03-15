package reporte

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
)

// ToCSV converts a slice of structs to CSV bytes using field names as headers.
func ToCSV(records interface{}) ([]byte, error) {
	v := reflect.ValueOf(records)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("records must be a slice")
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if v.Len() == 0 {
		w.Flush()
		return buf.Bytes(), nil
	}

	elem := v.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	t := elem.Type()

	headers := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		headers[i] = t.Field(i).Name
	}
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for i := 0; i < v.Len(); i++ {
		row := v.Index(i)
		if row.Kind() == reflect.Ptr {
			row = row.Elem()
		}
		record := make([]string, row.NumField())
		for j := 0; j < row.NumField(); j++ {
			record[j] = fmt.Sprintf("%v", row.Field(j).Interface())
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}
