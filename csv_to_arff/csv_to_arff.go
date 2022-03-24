package csv_to_arff

import (
	"bufio"
	"fmt"
	"ia-tarefa-arff/models"
	"ia-tarefa-arff/utils"
	"log"
	"os"
	"strings"
)

type csvToArff struct {
	filePath         string
	columns          []models.Column
	csvSeparator     string
	parsedLinesLimit int64
}

func NewCsvToArff(filePath string, columns []models.Column, csvSeparator string, parsedLinesLimit int64) *csvToArff {
	return &csvToArff{
		filePath:         filePath,
		columns:          columns,
		csvSeparator:     csvSeparator,
		parsedLinesLimit: parsedLinesLimit,
	}
}

func (cta *csvToArff) Parse() {
	file, err := os.Open(cta.filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	parsedLines := int64(0)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), cta.csvSeparator)

		if lineNumber == 1 {
			for columnIndex, value := range line {
				currentColumnIndex := indexOfColumnName(cta.columns, value)

				if currentColumnIndex != -1 {
					cta.columns[currentColumnIndex].IndexOnCsv = columnIndex
				}
			}
		} else {
			values := make([]string, len(cta.columns))
			err = nil

			for i, column := range cta.columns {
				values[i], err = processColumn(column, line)

				if err != nil {
					break
				}
			}

			if err == nil {
				// write in arff
				fmt.Println(lineNumber)
				fmt.Println(values)
				fmt.Println(line)
				fmt.Println("--------")
				parsedLines++
			}
		}

		if parsedLines == cta.parsedLinesLimit {
			break
		}

		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func processColumn(column models.Column, line []string) (string, error) {
	value := line[column.IndexOnCsv]

	value = strings.ToLower(value)
	value = strings.TrimSpace(value)
	value = utils.ReplaceSpecialCharOfString(value)

	switch column.Name {
	case "asma":
		value = convertBoolean(value)
	case "cardiopatia":
		value = convertBoolean(value)
	case "diabetes":
		value = convertBoolean(value)
	case "doenca_hematologica":
		value = convertBoolean(value)
	case "doenca_hepatica":
		value = convertBoolean(value)
	case "doenca_neurologica":
		value = convertBoolean(value)
	case "doenca_renal":
		value = convertBoolean(value)
	case "imunodepressao":
		value = convertBoolean(value)
	case "obesidade":
		value = convertBoolean(value)
	case "outros_fatores_de_risco":
		value = convertBoolean(value)
	case "pneumopatia":
		value = convertBoolean(value)
	case "puerpera":
		value = convertBoolean(value)
	case "sindrome_de_down":
		value = convertBoolean(value)
	}

	if column.AllowedValues != nil {
		if indexOfString(column.AllowedValues, value) == -1 {
			return "", fmt.Errorf("%v is not allowed to column %v", value, column.Name)
		}
	}

	return value, nil
}

func convertBoolean(value string) string {
	if value == "sim" {
		return "1"
	} else if value == "nao" {
		return "0"
	}

	return value
}

func indexOfString(array []string, target string) int {
	for i, value := range array {
		if value == target {
			return i
		}
	}

	return -1
}

func indexOfColumnName(array []models.Column, target string) int {
	for i, value := range array {
		if value.Name == target {
			return i
		}
	}

	return -1
}
