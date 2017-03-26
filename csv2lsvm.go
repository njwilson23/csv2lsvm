package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

var UNREADABLE_LABEL_ERROR = errors.New("unreadable label")

// Row represents a line of numerical data from a CSV or libSVM file, mapping a series of
// features to a label
type Row struct {
	Empty    bool
	Schema   []int
	Features []float64
	Label    float64
}

// ToString outputs a libSVM representation of a Row
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

type readOptions struct {
	StartRow     int
	NumberOfRows int
	Columns      []int
}

type writeOptions struct {
	Append bool
}

// Section is an array of Rows representing the contens of a file or a section
// of a file
type Section struct {
	Rows []Row
}

// WriteLibSVM sends a libSVM-formatted representation of a Section to a
// buffered Writer
func (section *Section) WriteLibSVM(buffer *bufio.Writer) error {
	var s string
	for _, row := range section.Rows {
		s = row.ToString()
		buffer.WriteString(s)
	}
	buffer.Flush()
	return nil
}

func readCSVRow(readBuffer *bufio.Reader) (*Row, error) {
	var err error
	var value float64
	row := Row{Empty: true}
	buffer := []byte{}

	line, err := readBuffer.ReadString('\n')
	if err == io.EOF && len(buffer) == 0 {
		return &row, io.EOF
	} else if err != nil {
		panic(err)
	}

	colNum := 0
	for i, r := range line {
		if r == ',' || i == len(line)-1 {
			if len(buffer) != 0 {
				row.Empty = false
				value, err = strconv.ParseFloat(string(buffer), 64)
				if err == nil {
					if colNum == 0 {
						row.Label = value
					} else {
						row.Schema = append(row.Schema, colNum)
						row.Features = append(row.Features, value)
					}
				} else if colNum == 0 {
					return &row, UNREADABLE_LABEL_ERROR
				}
				colNum++
				buffer = buffer[:0]
			}
		} else if r != ' ' && r != '\t' && r != '\r' { // exclude whitespace
			buffer = append(buffer, byte(r))
		}
	}
	return &row, nil
}

func readCSV(filePath string, options *readOptions) (*Section, error) {
	f, err := os.Open(filePath)
	if err != nil {
		panic("failure to open file for reading")
	}
	buffer := bufio.NewReader(f)

	// count the number of rows
	b := make([]byte, 1)
	var n int
	nRows := 1
	for {
		n, err = buffer.Read(b)
		if n == 0 {
			return &Section{}, io.EOF
		}
		if b[0] == ',' {
			nRows++
		} else if b[0] == '\n' {
			break
		}
	}

	// for each row, create a Row assuming the first column are labels
	var row *Row
	rows := []Row{}
	rowCount := 0
	for {
		row, err = readCSVRow(buffer)
		if err == io.EOF || err == UNREADABLE_LABEL_ERROR {
			break
		} else if row.Empty == true {
			break
		} else if err != nil {
			panic("failure to read CSV row")
		}
		rows = append(rows, *row)
		rowCount++
		if rowCount == options.NumberOfRows {
			break
		}
	}
	return &Section{rows}, nil
}

func writeLibSVMFile(filePath string, content *Section, options *writeOptions) error {
	f, err := os.Create(filePath)
	if err != nil {
		panic("failure to create file for writing")
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
	var output = flag.String("o", "out.svm",
		"output file; if not provided, data is written to out.svm")
	flag.Parse()
	input := flag.Arg(0)
	fmt.Println(input)
	fmt.Println(*output)

	fmt.Println(time.Now())
	section, err := readCSV(input, &readOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now())
	err = writeLibSVMFile(*output, section, &writeOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Now())
}
