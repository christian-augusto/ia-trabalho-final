package models

type Column struct {
	Name          string
	AllowedValues []string
	IndexOnCsv    int
}

func NewColumn(name string, allowedValues []string) *Column {
	return &Column{
		Name:          name,
		AllowedValues: allowedValues,
		IndexOnCsv:    -1,
	}
}
