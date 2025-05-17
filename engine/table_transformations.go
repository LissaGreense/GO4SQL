package engine

import (
	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
	"hash/adler32"
)

// TableTransformer defines a function that takes a Table as input and then applies a transformation
type TableTransformer func(*Table) (*Table, error)

// --- Basic SELECT ---
func (engine *DbEngine) withVanillaSelect(cmd *ast.SelectCommand) TableTransformer {
	return func(tbl *Table) (*Table, error) {
		return engine.selectFromProvidedTable(cmd, tbl)
	}
}

// --- WHERE only ---
func (engine *DbEngine) withWhere(cmd *ast.SelectCommand) TableTransformer {
	return func(tbl *Table) (*Table, error) {
		if len(tbl.Columns) == 0 || len(tbl.Columns[0].Values) == 0 {
			return engine.selectFromProvidedTable(cmd, &Table{Columns: []*Column{}})
		}
		filtered, err := engine.getFilteredTable(tbl, cmd.WhereCommand, false, cmd.Name.GetToken().Literal)
		if err != nil {
			return nil, err
		}
		return engine.selectFromProvidedTable(cmd, filtered)
	}
}

// --- ORDER BY only ---
func (engine *DbEngine) withOrderBy(cmd *ast.SelectCommand) TableTransformer {
	return func(tbl *Table) (*Table, error) {
		emptyTable := getCopyOfTableWithoutRows(tbl)
		sorted, err := engine.getSortedTable(cmd.OrderByCommand, tbl, emptyTable, cmd.Name.GetToken().Literal)
		if err != nil {
			return nil, err
		}
		return engine.selectFromProvidedTable(cmd, sorted)
	}
}

// --- OFFSET + LIMIT ---
func (engine *DbEngine) withOffsetLimit(cmd *ast.SelectCommand) TableTransformer {
	return func(tbl *Table) (*Table, error) {
		var offset = 0
		var limitRaw = -1

		if cmd.HasLimitCommand() {
			limitRaw = cmd.LimitCommand.Count
		}
		if cmd.HasOffsetCommand() {
			offset = cmd.OffsetCommand.Count
		}

		for _, column := range tbl.Columns {
			var limit int

			if limitRaw == -1 || limitRaw+offset > len(column.Values) {
				limit = len(column.Values)
			} else {
				limit = limitRaw + offset
			}

			if offset > len(column.Values) || limit == 0 {
				column.Values = make([]ValueInterface, 0)
			} else {
				column.Values = column.Values[offset:limit]
			}
		}

		return tbl, nil
	}
}

// --- DISTINCT ---
func (engine *DbEngine) withDistinct() TableTransformer {
	return func(tbl *Table) (*Table, error) {
		distinctTable := getCopyOfTableWithoutRows(tbl)

		rowsCount := len(tbl.Columns[0].Values)

		checksumSet := map[uint32]struct{}{}

		for iRow := 0; iRow < rowsCount; iRow++ {

			mergedColumnValues := ""
			for iColumn := range tbl.Columns {
				fieldValue := tbl.Columns[iColumn].Values[iRow].ToString()
				if tbl.Columns[iColumn].Type.Literal == token.TEXT {
					fieldValue = "'" + fieldValue + "'"
				}
				mergedColumnValues += fieldValue
			}
			checksum := adler32.Checksum([]byte(mergedColumnValues))

			_, exist := checksumSet[checksum]
			if !exist {
				checksumSet[checksum] = struct{}{}
				for i, column := range distinctTable.Columns {
					column.Values = append(column.Values, tbl.Columns[i].Values[iRow])
				}
			}
		}

		return distinctTable, nil
	}
}
