package main

import (
	"ia-tarefa-arff/csv_to_arff"
	"ia-tarefa-arff/models"
)

func main() {
	parsedLinesLimit := int64(5)
	filePath := "input/casos_obitos_doencas_preexistentes.csv"
	columns := []models.Column{
		*models.NewColumn("idade", nil),
		*models.NewColumn("cs_sexo", []string{"masculino", "feminino"}),
		*models.NewColumn("obito", []string{"1", "0"}),
		*models.NewColumn("asma", []string{"1", "0"}),
		*models.NewColumn("cardiopatia", []string{"1", "0"}),
		*models.NewColumn("diabetes", []string{"1", "0"}),
		*models.NewColumn("doenca_hematologica", []string{"1", "0"}),
		*models.NewColumn("doenca_hepatica", []string{"1", "0"}),
		*models.NewColumn("doenca_neurologica", []string{"1", "0"}),
		*models.NewColumn("doenca_renal", []string{"1", "0"}),
		*models.NewColumn("imunodepressao", []string{"1", "0"}),
		*models.NewColumn("obesidade", []string{"1", "0"}),
		*models.NewColumn("outros_fatores_de_risco", []string{"1", "0"}),
		*models.NewColumn("pneumopatia", []string{"1", "0"}),
		*models.NewColumn("puerpera", []string{"1", "0"}),
		*models.NewColumn("sindrome_de_down", []string{"1", "0"}),
	}
	csvSeparator := ";"

	csvToArff := csv_to_arff.NewCsvToArff(filePath, columns, csvSeparator, parsedLinesLimit)

	csvToArff.Parse()
}
