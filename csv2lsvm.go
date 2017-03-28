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
var CSV_READ_ERROR = errors.New("failure to read CSV row")

// Row represents a line of numerical data from a CSV or libSVM file, mapping a series of
// features to a label
type Row struct {
	Empty    bool
	Schema   []int
	Features []float64
	Label    float64
}

// ToString outputs a libSVM representation of a Row
func (row *Row) ToString(precision int) string {
	// TODO: implement variable precision
	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatFloat(row.Label, 'f', precision, 64))
	for i, feature := range row.Features {
		buffer.WriteRune(' ')
		buffer.WriteString(strconv.Itoa(row.Schema[i]))
		buffer.WriteRune(':')
		buffer.WriteString(strconv.FormatFloat(feature, 'f', precision, 64))
	}
	buffer.WriteRune('\n')
	return buffer.String()
}

type readOptions struct {
	StartRow     int // not implemented
	NumberOfRows int
	Columns      []int // not implemented
}

type writeOptions struct {
	Precision int
	Append    bool // not implemented
}

// Section is an array of Rows representing the contens of a file or a section
// of a file
type Section struct {
	Rows []Row
}

// WriteLibSVM sends a libSVM-formatted representation of a Section to a
// buffered Writer
func (section *Section) WriteLibSVM(buffer *bufio.Writer, precision int) error {
	var s string
	for _, row := range section.Rows {
		s = row.ToString(precision)
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
	if err != nil {
		return &row, err
	}

	colNum := 0
	lineLength := len(line)
	for i, rn := range line {

		if rn == ',' || i == lineLength-1 {
			if len(buffer) != 0 {
				row.Empty = false
				value, err = strconv.ParseFloat(string(buffer), 64)
				if err == nil {
					if colNum != 0 {
						row.Schema = append(row.Schema, colNum)
						row.Features = append(row.Features, value)
					} else {
						row.Label = value
					}
				} else if colNum == 0 {
					return &row, UNREADABLE_LABEL_ERROR
				}
				buffer = buffer[:0]
			}
		} else if rn != ' ' && rn != '\t' && rn != '\r' { // exclude whitespace
			buffer = append(buffer, byte(rn))
		}

		if rn == ',' {
			colNum++
		}

	}
	return &row, nil
}

func readCSV(filePath string, options *readOptions) (*Section, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return &Section{}, err
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
			return &Section{}, errors.New(fmt.Sprintf("failure to read CSV row %d", rowCount))
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
		return err
	}
	defer f.Close()

	buffer := bufio.NewWriter(f)
	err = content.WriteLibSVM(buffer, options.Precision)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var precision = flag.Int("p", 4, "decimal precision")
	var output = flag.String("o", "out.svm",
		"output file; if not provided, data is written to out.svm")
	flag.Parse()
	input := flag.Arg(0)
	fmt.Println(input)
	fmt.Println(*output)

	fmt.Println(time.Now())
	section, err := readCSV(input, &readOptions{})
	fmt.Println(time.Now())
	if err != nil {
		fmt.Println("failure to read CSV")
		fmt.Println(err)
		os.Exit(1)
	}
	err = writeLibSVMFile(*output, section, &writeOptions{Precision: *precision})
	fmt.Println(time.Now())
	if err != nil {
		fmt.Println("failure to write libSVM file")
		fmt.Println(err)
		os.Exit(1)
	}
}
