package engine

import (
	"hash/adler32"

	"github.com/LissaGreense/GO4SQL/ast"
	"github.com/LissaGreense/GO4SQL/token"
)

// TableTransformer defines a function that takes a Table as input and then applies a transformation
type TableTransformer func(*Table) (*Table, error)

// SelectProcessor handles the step-by-step processing of a SELECT query using a builder pattern.
type SelectProcessor struct {
	engine       *DbEngine
	cmd          *ast.SelectCommand
	transformers []TableTransformer
}

// NewSelectProcessor creates a new SelectProcessor.
func NewSelectProcessor(engine *DbEngine, cmd *ast.SelectCommand) *SelectProcessor {
	return &SelectProcessor{
		engine:       engine,
		cmd:          cmd,
		transformers: []TableTransformer{}, // Initialize as empty slice
	}
}

// WithVanillaSelectClause adds the vanilla select (projection) transformation.
func (sp *SelectProcessor) WithVanillaSelectClause() *SelectProcessor {
	sp.transformers = append(sp.transformers, sp.getVanillaSelectTransformer())
	return sp
}

// WithWhereClause adds the WHERE clause transformation.
func (sp *SelectProcessor) WithWhereClause() *SelectProcessor {
	sp.transformers = append(sp.transformers, sp.getWhereTransformer())
	return sp
}

// WithOrderByClause adds the ORDER BY clause transformation.
func (sp *SelectProcessor) WithOrderByClause() *SelectProcessor {
	sp.transformers = append(sp.transformers, sp.getOrderByTransformer())
	return sp
}

// WithOffsetLimitClause adds the OFFSET and LIMIT clause transformation.
func (sp *SelectProcessor) WithOffsetLimitClause() *SelectProcessor {
	sp.transformers = append(sp.transformers, sp.getOffsetLimitTransformer())
	return sp
}

// WithDistinctClause adds the DISTINCT clause transformation.
func (sp *SelectProcessor) WithDistinctClause() *SelectProcessor {
	sp.transformers = append(sp.transformers, sp.getDistinctTransformer())
	return sp
}

// Process applies the configured pipeline of transformations to the initialTable.
func (sp *SelectProcessor) Process(initialTable *Table) (*Table, error) {
	table := initialTable
	var err error

	for _, transform := range sp.transformers {
		table, err = transform(table)
		if err != nil {
			return nil, err
		}
	}
	return table, nil
}

// --- Private Transformer Getters ---

func (sp *SelectProcessor) getVanillaSelectTransformer() TableTransformer {
	return func(tbl *Table) (*Table, error) {
		return sp.engine.selectFromProvidedTable(sp.cmd, tbl)
	}
}

func (sp *SelectProcessor) getWhereTransformer() TableTransformer {
	return func(tbl *Table) (*Table, error) {
		if len(tbl.Columns) == 0 || (len(tbl.Columns) > 0 && len(tbl.Columns[0].Values) == 0) {
			return sp.engine.selectFromProvidedTable(sp.cmd, &Table{Columns: []*Column{}})
		}
		filtered, err := sp.engine.getFilteredTable(tbl, sp.cmd.WhereCommand, false, sp.cmd.Name.GetToken().Literal)
		if err != nil {
			return nil, err
		}
		return sp.engine.selectFromProvidedTable(sp.cmd, filtered)
	}
}

func (sp *SelectProcessor) getOrderByTransformer() TableTransformer {
	return func(tbl *Table) (*Table, error) {
		emptyTable := getCopyOfTableWithoutRows(tbl)
		sorted, err := sp.engine.getSortedTable(sp.cmd.OrderByCommand, tbl, emptyTable, sp.cmd.Name.GetToken().Literal)
		if err != nil {
			return nil, err
		}
		return sp.engine.selectFromProvidedTable(sp.cmd, sorted)
	}
}

func (sp *SelectProcessor) getOffsetLimitTransformer() TableTransformer {
	return func(tbl *Table) (*Table, error) {
		var offset = 0
		var limitRaw = -1

		if sp.cmd.HasLimitCommand() {
			limitRaw = sp.cmd.LimitCommand.Count
		}
		if sp.cmd.HasOffsetCommand() {
			offset = sp.cmd.OffsetCommand.Count
		}

		if len(tbl.Columns) == 0 {
			return tbl, nil
		}

		for _, column := range tbl.Columns {
			var limit int

			if limitRaw == -1 || limitRaw+offset > len(column.Values) {
				limit = len(column.Values)
			} else {
				limit = limitRaw + offset
			}

			if offset >= len(column.Values) {
				column.Values = make([]ValueInterface, 0)
			} else if offset < len(column.Values) && limit > offset {
				column.Values = column.Values[offset:limit]
			} else {
				column.Values = make([]ValueInterface, 0)
			}
		}
		return tbl, nil
	}
}

func (sp *SelectProcessor) getDistinctTransformer() TableTransformer {
	return func(tbl *Table) (*Table, error) {
		if len(tbl.Columns) == 0 || len(tbl.Columns[0].Values) == 0 {
			return tbl, nil
		}

		distinctTable := getCopyOfTableWithoutRows(tbl)
		rowsCount := len(tbl.Columns[0].Values)
		checksumSet := make(map[uint32]struct{})

		for iRow := range rowsCount {
			mergedColumnValues := ""
			for iColumn := range tbl.Columns {
				if iRow < len(tbl.Columns[iColumn].Values) {
					fieldValue := tbl.Columns[iColumn].Values[iRow].ToString()
					if tbl.Columns[iColumn].Type.Literal == token.TEXT {
						fieldValue = "'" + fieldValue + "'"
					}
					mergedColumnValues += fieldValue
				} else {
					mergedColumnValues += "<missing_val>"
				}
			}
			checksum := adler32.Checksum([]byte(mergedColumnValues))

			if _, exist := checksumSet[checksum]; !exist {
				checksumSet[checksum] = struct{}{}
				for i, column := range distinctTable.Columns {
					if iRow < len(tbl.Columns[i].Values) {
						column.Values = append(column.Values, tbl.Columns[i].Values[iRow])
					}
				}
			}
		}
		return distinctTable, nil
	}
}
