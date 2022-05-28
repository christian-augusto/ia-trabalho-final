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

	// arffFile, err = os.Create("output/casos_obitos_doencas_preexistentes.arff")
	arffFile, err = os.Create("output/" + cta.customName() + ".arff")

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

		if lineNumber%1000 == 0 {
			fmt.Print(".")
		}

		if lineNumber == 1 {
			for columnIndex, value := range line {
				currentColumnIndex := cta.indexOfColumnName(cta.columns, value)

				if currentColumnIndex != -1 {
					cta.columns[currentColumnIndex].IndexOnCsv = columnIndex
				}
			}
		} else {
			values := make([]string, len(cta.columns))
			err = nil

			for i, column := range cta.columns {
				values[i], err = cta.processColumn(column, line)

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

	fmt.Print("\n")

	fmt.Printf("Total lines: %v\n", lineNumber)
	fmt.Printf("Lines parsed: %v\n", parsedLines)

	if err := scanner.Err(); err != nil {
		return err
	}

	return err
}

func (cta *csvToArff) processColumn(column models.Column, line []string) (string, error) {
	var err error

	value := line[column.IndexOnCsv]

	value = strings.ToLower(value)
	value = strings.TrimSpace(value)
	value = utils.ReplaceSpecialCharOfString(value)

	switch column.Name {
	case "idade":
		value, err = cta.convertAge(value)
	case "cs_sexo":
		value = cta.convertGender(value)
	case "asma":
		value = cta.convertBoolean(value)
	case "cardiopatia":
		value = cta.convertBoolean(value)
	case "diabetes":
		value = cta.convertBoolean(value)
	case "doenca_hematologica":
		value = cta.convertBoolean(value)
	case "doenca_hepatica":
		value = cta.convertBoolean(value)
	case "doenca_neurologica":
		value = cta.convertBoolean(value)
	case "doenca_renal":
		value = cta.convertBoolean(value)
	case "imunodepressao":
		value = cta.convertBoolean(value)
	case "obesidade":
		value = cta.convertBoolean(value)
	case "outros_fatores_de_risco":
		value = cta.convertBoolean(value)
	case "pneumopatia":
		value = cta.convertBoolean(value)
	case "puerpera":
		value = cta.convertBoolean(value)
	case "sindrome_de_down":
		value = cta.convertBoolean(value)
	}

	if err != nil {
		return "", err
	}

	if column.AllowedValues != nil {
		if cta.indexOfString(column.AllowedValues, value) == -1 {
			return "", fmt.Errorf("%v is not allowed to column %v", value, column.Name)
		}
	}

	return value, err
}

func (cta *csvToArff) customName() string {
	var name string

	if len(cta.columns) > 0 {
		name = cta.columns[0].Name
	}

	for i := 1; i < len(cta.columns); i++ {
		name += "-" + cta.columns[i].Name
	}

	return name
}

func (cta *csvToArff) convertAge(value string) (string, error) {
	var age int
	var err error

	age, err = strconv.Atoi(value)

	if err != nil {
		return "", err
	}

	if age < 13 {
		return "crianca", err
	} else if age < 20 {
		return "adolescente", err
	} else if age < 36 {
		return "jovem", err
	} else if age < 60 {
		return "adulto", err
	} else {
		return "idoso", err
	}
}

func (cta *csvToArff) convertBoolean(value string) string {
	if value == "sim" {
		return "1"
	} else if value == "nao" {
		return "0"
	}

	return value
}

func (cta *csvToArff) convertGender(value string) string {
	if value == "masculino" {
		return "m"
	} else if value == "feminino" {
		return "f"
	}

	return value
}

func (cta *csvToArff) indexOfString(array []string, target string) int {
	for i, value := range array {
		if value == target {
			return i
		}
	}

	return -1
}

func (cta *csvToArff) indexOfColumnName(array []models.Column, target string) int {
	for i, value := range array {
		if value.Name == target {
			return i
		}
	}

	return -1
}
