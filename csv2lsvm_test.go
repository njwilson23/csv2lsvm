package main

import (
	"bufio"
	"fmt"
	"testing"
)

type mockWriter struct {
	content []byte
}

func (w *mockWriter) Write(p []byte) (int, error) {
	w.content = append(w.content, p...)
	return len(p), nil
}

func (w *mockWriter) String() string {
	return string(w.content)
}

func TestRowToString(t *testing.T) {
	var row Row
	var s string

	row = Row{Empty: false, Schema: []int{0, 1, 2}, Features: []float64{1.5, 2.5, 3.5}, Label: 10}
	s = row.ToString(2)
	if s != "10.00 0:1.50 1:2.50 2:3.50\n" {
		t.Fail()
	}

	s = row.ToString(3)
	if s != "10.000 0:1.500 1:2.500 2:3.500\n" {
		t.Fail()
	}

	row = Row{Empty: false, Schema: []int{1, 2}, Features: []float64{2.5, 3.5}, Label: -10}
	s = row.ToString(2)
	if s != "-10.00 1:2.50 2:3.50\n" {
		t.Fail()
	}
}

func TestWriteLibSVM(t *testing.T) {
	section := Section{[]Row{
		Row{Empty: false, Schema: []int{0, 1, 2}, Features: []float64{1.5, 2.5, 3.5}, Label: 10},
		Row{Empty: false, Schema: []int{0, 1, 2}, Features: []float64{2.5, 3.5, 1.5}, Label: 2.1},
		Row{Empty: false, Schema: []int{0, 2}, Features: []float64{1.5, 2.5}, Label: -4}}}

	writer := &mockWriter{[]byte{}}
	buffer := bufio.NewWriter(writer)
	err := section.WriteLibSVM(buffer, 2)
	if err != nil {
		t.Error()
	}
	if writer.String() != "10.00 0:1.50 1:2.50 2:3.50\n2.10 0:2.50 1:3.50 2:1.50\n-4.00 0:1.50 2:2.50\n" {
		t.Fail()
	}
}

func TestWriteLibSVMFile(t *testing.T) {
	section := Section{[]Row{
		Row{Empty: false, Schema: []int{0, 1, 2}, Features: []float64{1.5, 2.5, 3.5}, Label: 10},
		Row{Empty: false, Schema: []int{0, 1, 2}, Features: []float64{2.5, 3.5, 1.5}, Label: 2.1},
		Row{Empty: false, Schema: []int{0, 2}, Features: []float64{1.5, 2.5}, Label: -4}}}

	options := writeOptions{Precision: 2, Append: false}

	err := writeLibSVMFile("test.svm", &section, &options)
	if err != nil {
		t.Error()
	}
}

func float64SlicesEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func intSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestReadDenseCSV(t *testing.T) {
	section, err := readCSV("test_dense.csv", &readOptions{})
	if err != nil {
		t.Error()
	}

	if len(section.Rows) != 6 {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[0].Features, []float64{1, 2, 3}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[1].Features, []float64{1, 0, 2}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[2].Features, []float64{0.5, 0, 1}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[3].Features, []float64{2, 0.25, -3}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[4].Features, []float64{1.5, 2, 1}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[5].Features, []float64{-1, 0.5, 0.75}) {
		t.Fail()
	}
	for i, row := range section.Rows {
		if !intSlicesEqual(row.Schema, []int{1, 2, 3}) {
			fmt.Println(i)
			fmt.Println(row.Schema)
			t.Fail()
		}
	}
}

func TestReadSparseCSVWithNA(t *testing.T) {
	section, err := readCSV("test_sparse_NA.csv", &readOptions{})
	if err != nil {
		t.Error()
	}

	if len(section.Rows) != 6 {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[0].Features, []float64{1, 2}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[1].Features, []float64{1, 2}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[2].Features, []float64{0.5, 0, 1}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[3].Features, []float64{2, 0.25, -3}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[4].Features, []float64{2, 1}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[5].Features, []float64{-1, 0.75}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[0].Schema, []int{1, 2}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[1].Schema, []int{1, 3}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[2].Schema, []int{1, 2, 3}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[3].Schema, []int{1, 2, 3}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[4].Schema, []int{2, 3}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[5].Schema, []int{1, 3}) {
		t.Fail()
	}
}
func TestReadSparseCSVWithBlank(t *testing.T) {
	section, err := readCSV("test_sparse_blank.csv", &readOptions{})
	if err != nil {
		t.Error()
	}

	if len(section.Rows) != 6 {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[0].Features, []float64{1, 2}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[1].Features, []float64{1, 2}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[2].Features, []float64{0.5, 0, 1}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[3].Features, []float64{2, 0.25, -3}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[4].Features, []float64{2, 1}) {
		t.Fail()
	}
	if !float64SlicesEqual(section.Rows[5].Features, []float64{-1, 0.75}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[0].Schema, []int{1, 2}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[1].Schema, []int{1, 3}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[2].Schema, []int{1, 2, 3}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[3].Schema, []int{1, 2, 3}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[4].Schema, []int{2, 3}) {
		t.Fail()
	}
	if !intSlicesEqual(section.Rows[5].Schema, []int{1, 3}) {
		t.Fail()
	}
}
