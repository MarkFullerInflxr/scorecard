package spreadsheet

import (
	"bytes"
	"fmt"
	utils "influxer/scorecard/utilities"
)

func (s *Scorecard) BuildCell(tableName string, row int, col int) []byte {
	table := s.data[tableName]

	if table == nil {
		s.data[tableName] = NewTable(tableName)
		table = s.data[tableName]
	}

	var tplBuffer bytes.Buffer
	err := s.tmpl["table/cell"].ExecuteTemplate(&tplBuffer, "table/cell", CellArgs{
		RowIndex:  row,
		ColIndex:  col,
		TableName: tableName,
		CellVal:   string(table.Matrix[row][col].Value),
		CellType:  table.Matrix[row][col].Type,
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return tplBuffer.Bytes()

}

func (s *Scorecard) BuildIndex(tableName string) []byte {
	table := s.data[tableName]

	if table == nil {
		s.data[tableName] = NewTable(tableName)
		table = s.data[tableName]
		s.notifyList()
	}

	var tplBuffer bytes.Buffer
	err := s.tmpl["table/index"].ExecuteTemplate(&tplBuffer, "table/index", table)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return tplBuffer.Bytes()

}

func (s *Scorecard) BuildListIndex() []byte {
	tables := utils.MapMap(s.data, func(k string, v *Table) *Table {
		return v
	})

	var tplBuffer bytes.Buffer
	err := s.tmpl["list/index"].ExecuteTemplate(&tplBuffer, "list/index", tables)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return tplBuffer.Bytes()
}

func (s *Scorecard) BuildTableList() []byte {
	tables := utils.MapMap(s.data, func(k string, v *Table) *Table {
		return v
	})

	var tplBuffer bytes.Buffer
	err := s.tmpl["list/tablelist"].ExecuteTemplate(&tplBuffer, "list/tablelist", tables)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return tplBuffer.Bytes()
}
func (s *Scorecard) BuildTable(tableName string) []byte {
	table := s.data[tableName]

	if table == nil {
		s.data[tableName] = NewTable(tableName)
		table = s.data[tableName]
	}

	var tplBuffer bytes.Buffer
	err := s.tmpl["table/table"].ExecuteTemplate(&tplBuffer, "table/table", table)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return tplBuffer.Bytes()

}

func (s *Scorecard) BuildHeaderCell(tableName string, row int) []byte {
	var tplBuffer bytes.Buffer
	s.data[tableName].changeHeaderVal(row, s.data[tableName].Headers[row].Value)
	err := s.tmpl["table/headercell"].ExecuteTemplate(&tplBuffer, "table/headercell",
		HeaderArgs{
			row,
			tableName,
			string(s.data[tableName].Headers[row].Value)})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tplBuffer.Bytes()
}
