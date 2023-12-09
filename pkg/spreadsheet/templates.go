package spreadsheet

import (
	"bytes"
	"fmt"
)

func (s *Scorecard) BuildCell(tableName string, row int, col int) []byte {
	table := s.data[tableName]

	if table == nil {
		s.data[tableName] = NewTable(tableName)
		table = s.data[tableName]
	}

	var tplBuffer bytes.Buffer
	err := s.tmpl["cell.html"].ExecuteTemplate(&tplBuffer, "cell", CellArgs{
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
	}

	var tplBuffer bytes.Buffer
	err := s.tmpl["index.html"].ExecuteTemplate(&tplBuffer, "index", table)
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
	err := s.tmpl["table.html"].ExecuteTemplate(&tplBuffer, "table", table)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return tplBuffer.Bytes()

}

func (s *Scorecard) BuildHeaderCell(tableName string, row int) []byte {
	var tplBuffer bytes.Buffer
	s.data[tableName].changeHeaderVal(row, s.data[tableName].Headers[row].Value)
	err := s.tmpl["headercell.html"].ExecuteTemplate(&tplBuffer, "headercell",
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
