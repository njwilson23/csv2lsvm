package main

import (
	"bufio"
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
	s = row.ToString()
	if s != "10.00 0:1.50 1:2.50 2:3.50\n" {
		t.Fail()
	}

	row = Row{Empty: false, Schema: []int{1, 2}, Features: []float64{2.5, 3.5}, Label: -10}
	s = row.ToString()
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
	err := section.WriteLibSVM(buffer)
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

	options := writeOptions{false}

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

func TestReadCSV(t *testing.T) {
	section, err := readCSV("test.csv", &readOptions{})
	if err != nil {
		t.Error()
	}
	for _, row := range section.Rows {
		fmt.Print(row.Label)
		fmt.Print(" ")
		fmt.Println(row.Features)
	}
}
