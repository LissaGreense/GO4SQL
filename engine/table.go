package engine

import "github.com/LissaGreense/GO4SQL/token"

type Table struct {
	Columns []*Column
}

func (table *Table) isEqual(secondTable *Table) bool {
	if len(table.Columns) != len(secondTable.Columns) {
		return false
	}

	for i := 0; i < len(table.Columns); i++ {
		if table.Columns[i].Name != secondTable.Columns[i].Name {
			return false
		}
		if table.Columns[i].Type.Literal != secondTable.Columns[i].Type.Literal {
			return false
		}
		if table.Columns[i].Type.Type != secondTable.Columns[i].Type.Type {
			return false
		}
		if len(table.Columns[i].Values) != len(secondTable.Columns[i].Values) {
			return false
		}
		for j := 0; j < len(table.Columns[i].Values); j++ {
			if table.Columns[i].Values[j].ToString() != secondTable.Columns[i].Values[j].ToString() {
				return false
			}
		}
	}

	return true
}

func (table *Table) ToString() string {
	columWidths := getColumWidths(table.Columns)
	bar := getBar(columWidths)
	result := bar + "\n"

	result += "|"
	for i := 0; i < len(table.Columns); i++ {
		result += " "
		for j := 0; j < columWidths[i]-len(table.Columns[i].Name); j++ {
			result += " "
		}
		result += table.Columns[i].Name
		result += " |"
	}
	result += "\n" + bar + "\n"

	rowsCount := len(table.Columns[0].Values)

	for iRow := 0; iRow < rowsCount; iRow++ {
		result += "|"

		for iColumn := 0; iColumn < len(table.Columns); iColumn++ {
			result += " "

			printedValue := table.Columns[iColumn].Values[iRow].ToString()
			if table.Columns[iColumn].Type.Literal == token.TEXT {
				printedValue = "'" + printedValue + "'"
			}
			for i := 0; i < columWidths[iColumn]-len(printedValue); i++ {
				result += " "
			}

			result += printedValue + " |"
		}

		result += "\n"
	}

	return result + bar
}

func (table *Table) appendRow(providedTable []*Column, rowIndex int) {
	for columnIndex := range table.Columns {
		table.Columns[columnIndex].Values = append(table.Columns[columnIndex].Values, providedTable[columnIndex].Values[rowIndex])
	}
}

func getBar(columWidths []int) string {
	bar := "+"

	for i := 0; i < len(columWidths); i++ {
		bar += "-"
		for j := 0; j < columWidths[i]; j++ {
			bar += "-"
		}
		bar += "-+"
	}

	return bar
}

func getColumWidths(columns []*Column) []int {
	widths := make([]int, 0)

	for iColumn := 0; iColumn < len(columns); iColumn++ {
		maxLength := len(columns[iColumn].Name)
		for iRow := 0; iRow < len(columns[iColumn].Values); iRow++ {
			valueLength := len(columns[iColumn].Values[iRow].ToString())
			if columns[iColumn].Type.Literal == token.TEXT {
				valueLength += 2 // double "'"
			}
			if valueLength > maxLength {
				maxLength = valueLength
			}
		}
		widths = append(widths, maxLength)
	}

	return widths
}
