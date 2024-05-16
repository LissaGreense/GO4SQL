package engine

// Rows - Contain rows that store values, alternative to Table, some operations are easier
type Rows struct {
	rows []map[string]ValueInterface
}

// MapTableToRows - transform Table struct into Rows
func MapTableToRows(table *Table) Rows {
	rows := make([]map[string]ValueInterface, 0)

	numberOfRows := len(table.Columns[0].Values)

	for rowIndex := 0; rowIndex < numberOfRows; rowIndex++ {
		row := make(map[string]ValueInterface)
		for _, column := range table.Columns {
			row[column.Name] = column.Values[rowIndex]
		}
		rows = append(rows, row)
	}
	return Rows{rows: rows}
}
