package csv_to_arff

import (
	"bufio"
	"fmt"
	"ia-tarefa-arff/models"
	"ia-tarefa-arff/utils"
	"os"
	"strconv"
	"strings"
)

type csvToArff struct {
	csvPath          string
	columns          []models.Column
	csvSeparator     string
	parsedLinesLimit int
}

func NewCsvToArff(csvPath string, columns []models.Column, csvSeparator string, parsedLinesLimit int) *csvToArff {
	return &csvToArff{
		csvPath:          csvPath,
		columns:          columns,
		csvSeparator:     csvSeparator,
		parsedLinesLimit: parsedLinesLimit,
	}
}

func (cta *csvToArff) Parse() error {
	var err error
	var csvFile *os.File
	var arffFile *os.File

	csvFile, err = os.Open(cta.csvPath)

	if err != nil {
		return err
	}

	defer csvFile.Close()

	arffFile, err = os.Create("output/casos_obitos_doencas_preexistentes.arff")

	if err != nil {
		return err
	}

	defer arffFile.Close()

	_, err = arffFile.WriteString("@RELATION casos_obitos_doencas_preexistentes\n")

	if err != nil {
		return fmt.Errorf("Error in write relation")
	}

	for _, column := range cta.columns {
		_, err = arffFile.WriteString(
			fmt.Sprintf("@ATTRIBUTE %v {%v}\n", column.Name, strings.Join(column.AllowedValues, ",")),
		)

		if err != nil {
			return fmt.Errorf("Error in write attribute %v: %v", column.Name, err)
		}
	}

	_, err = arffFile.WriteString("@DATA\n")

	if err != nil {
		return fmt.Errorf("Error in write @DATA")
	}

	scanner := bufio.NewScanner(csvFile)
	lineNumber := 1
	parsedLines := 0

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
				parsedLines++

				_, err = arffFile.WriteString(strings.Join(values, ",") + "\n")

				if err != nil {
					return fmt.Errorf("Error in write line %v: %v", lineNumber, err)
				}
			}
		}

		if parsedLines == cta.parsedLinesLimit {
			break
		}

		lineNumber++
	}

	fmt.Printf("Total lines: %v\n", lineNumber)
	fmt.Printf("Lines parsed: %v\n", parsedLines)

	if err := scanner.Err(); err != nil {
		return err
	}

	return err
}

func processColumn(column models.Column, line []string) (string, error) {
	var err error

	value := line[column.IndexOnCsv]

	value = strings.ToLower(value)
	value = strings.TrimSpace(value)
	value = utils.ReplaceSpecialCharOfString(value)

	switch column.Name {
	case "idade":
		value, err = convertAge(value)
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

	if err != nil {
		return "", err
	}

	if column.AllowedValues != nil {
		if indexOfString(column.AllowedValues, value) == -1 {
			return "", fmt.Errorf("%v is not allowed to column %v", value, column.Name)
		}
	}

	return value, err
}

func convertBoolean(value string) string {
	if value == "sim" {
		return "1"
	} else if value == "nao" {
		return "0"
	}

	return value
}

func convertAge(value string) (string, error) {
	var age int
	var err error

	age, err = strconv.Atoi(value)

	if err != nil {
		return "", err
	}

	if age < 13 {
		return "child", err
	} else if age < 20 {
		return "teen", err
	} else if age < 36 {
		return "young", err
	} else if age < 60 {
		return "adult", err
	} else {
		return "elderly", err
	}
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
