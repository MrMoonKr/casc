package common

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/pkg/errors"
)

//ParseCSV returns an error if not all columns are found.
func ParseCSV(r io.Reader, columns ...string) ([]map[string]string, error) {
	columnsCheck := map[string]struct{}{}
	for _, name := range columns {
		columnsCheck[name] = struct{}{}
	}
	csv := csv.NewReader(r)
	csv.Comma = '|'
	csv.Comment = '#'
	lines, err := csv.ReadAll()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(lines) == 0 {
		return nil, errors.New("csv has no lines")
	}
	names := lines[0]
	data := lines[1:]
	lookup := []map[string]string{}
	for _, line := range data {
		l := map[string]string{}
		for i, name := range names {
			n := strings.Split(name, "!")
			delete(columnsCheck, n[0])
			l[n[0]] = line[i]
		}
		lookup = append(lookup, l)
	}
	if len(columnsCheck) > 0 {
		var notFound []string
		for col := range columnsCheck {
			notFound = append(notFound, col)
		}
		return nil, errors.Errorf("colums not found in the csv: %s", strings.Join(notFound, ", "))
	}
	return lookup, nil
}
