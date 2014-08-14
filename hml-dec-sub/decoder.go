package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	timefmt = "2006-01-02 15:04"
)

type Decoder struct {
	r    *csv.Reader
	init bool // whether Decoder has been initialized
}

func NewDecoder(r io.Reader) Decoder {
	return Decoder{r: csv.NewReader(bufio.NewReader(r))}
}

// ReadHeader reads the first row of the underlying csv file and
// makes sure it has the expected format
func (dec *Decoder) ReadHeader() error {
	var err error
	row, err := dec.r.Read()
	rt := reflect.TypeOf((*Row)(nil)).Elem()
	nmax := rt.NumField()
	if len(row) < nmax {
		nmax = len(row)
	}

	for i := 0; i < nmax; i++ {
		field := rt.Field(i)
		name := field.Tag.Get("hml")
		if name == "" {
			name = field.Name
		}
		row[i] = strings.TrimSpace(row[i])
		if name != row[i] {
			return fmt.Errorf("hml: field #%d. expected [%s]. got [%s]",
				i,
				name,
				row[i],
			)
		}
	}
	return err
}

func (dec *Decoder) Decode(data *Row) error {
	if !dec.init {
		if err := dec.ReadHeader(); err != nil {
			return err
		}
		dec.init = true
	}

	row, err := dec.r.Read()
	if err != nil {
		return err
	}
	// remove stray-\r
	row[len(row)-1] = strings.Replace(row[len(row)-1], "\r", "", -1)

	// fmt.Printf("row: %q\n", row)

	idx := 0
	data.SubmissionID, err = strconv.Atoi(row[idx])
	if err != nil {
		return err
	}

	idx++
	data.DateSubmittedUTC, err = timeParse(timefmt, row[idx])
	if err != nil {
		return err
	}

	idx++
	data.TeamID, err = strconv.Atoi(row[idx])
	if err != nil {
		return err
	}

	idx++
	data.TeamName = row[idx]

	idx++
	data.UserID, err = strconv.Atoi(row[idx])
	if err != nil {
		return err
	}

	idx++
	data.UserDisplayName = row[idx]

	idx++
	data.PublicScore, err = parseFloat(row[idx], 64)
	if err != nil {
		return err
	}

	idx++
	data.PrivateScore, err = parseFloat(row[idx], 64)
	if err != nil {
		return err
	}

	idx++
	data.IsSelected = strings.Contains(strings.ToLower(row[idx]), "true")

	idx++
	data.DateRescoredUTC, err = timeParse(timefmt, row[idx])
	if err != nil {
		fmt.Printf("%q -> %v\n", row[idx], err)
		return err
	}

	idx++
	data.PrevPublicScore, err = parseFloat(row[idx], 64)
	if err != nil {
		return err
	}

	idx++
	data.PrevPrivateScore, err = parseFloat(row[idx], 64)
	if err != nil {
		return err
	}

	return err
}

func timeParse(layout, value string) (time.Time, error) {
	var t time.Time
	var err error
	if value == "" {
		return t, err
	}

	return time.Parse(layout, value)
}

func parseFloat(s string, bitSize int) (float64, error) {
	var f float64
	var err error
	if s == "" {
		return f, err
	}

	return strconv.ParseFloat(s, bitSize)
}
