package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Row struct {
	RowNum   int
	Schema   []int
	Features []float64
	Label    float64
}

func (row *Row) ToString() string {
	// TODO: implement variable precision
	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatFloat(row.Label, 'f', 2, 64))
	for i, feature := range row.Features {
		buffer.WriteRune(' ')
		buffer.WriteString(strconv.Itoa(row.Schema[i]))
		buffer.WriteRune(':')
		buffer.WriteString(strconv.FormatFloat(feature, 'f', 2, 64))
	}
	buffer.WriteRune('\n')
	return buffer.String()
}

type ReadOptions struct {
	StartRow     int
	NumberOfRows int
	Columns      []int
}

type WriteOptions struct {
	Append bool
}

type Section struct {
	Rows []Row
}

func (section *Section) WriteLibSVM(buffer *bufio.Writer) error {
	var s string
	for _, row := range section.Rows {
		s = row.ToString()
		buffer.WriteString(s)
	}
	buffer.Flush()
	return nil
}

func readCSVRow(f *os.File) (*Row, error) {
	var n int
	var err error
	var value float64
	row := Row{RowNum: -1}
	b := make([]byte, 1)
	buffer := []byte{}
	colNum := 0
	for {
		n, err = f.Read(b)
		if err != nil {
			if err == io.EOF {
				return &row, nil
			}
			return &row, err
		}
		if n == 0 {
			return &row, nil
		}
		if b[0] == 44 || b[0] == 10 { // comma or newline
			if len(buffer) != 0 {
				row.RowNum = 0
				value, err = strconv.ParseFloat(string(buffer), 64)
				if err != nil {
					// TODO: handle this, as it will happen whenever the CSV contains a value
					// unlike a float
					panic(fmt.Sprintf("failure parsing float: '%s'", string(buffer)))
				}
				if colNum == 0 {
					row.Label = value
				} else {
					row.Schema = append(row.Schema, colNum)
					row.Features = append(row.Features, value)
				}
				colNum++
				buffer = buffer[:0]
			}
		} else if b[0] != 32 && b[0] != 116 { // exclude whitespace
			buffer = append(buffer, b...)
		}

		if b[0] == 10 { // newline
			return &row, nil
		}
	}
}

func readCSV(filePath string, options *ReadOptions) (*Section, error) {
	f, err := os.Open(filePath)
	if err != nil {
		panic("failure to open file for reading")
	}

	// count the number of rows
	b := make([]byte, 1)
	var n int
	nRows := 1
	for {
		n, err = f.Read(b)
		if n == 0 {
			return &Section{}, errors.New("end of file encountered unexpectedly")
		}
		if b[0] == 44 {
			nRows++
		} else if b[0] == 10 {
			break
		}
	}

	// for each row, create a Row assuming the first column are labels
	var row *Row
	rows := []Row{}
	rowCount := 0
	for {
		row, err = readCSVRow(f)
		if err != nil {
			panic("failure to read CSV row")
		}
		if row.RowNum == -1 {
			break
		}
		rows = append(rows, *row)
		rowCount++
		if rowCount == options.NumberOfRows {
			break
		}
	}
	return &Section{rows}, nil
}

func writeLibSVMFile(filePath string, content *Section, options *WriteOptions) error {
	f, err := os.Create(filePath)
	if err != nil {
		panic("failure to create file for writing")
		//return err
	}
	defer f.Close()

	buffer := bufio.NewWriter(f)
	err = content.WriteLibSVM(buffer)
	if err != nil {
		return err
	}
	return nil
}

func main() {
}
