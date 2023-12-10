package spreadsheet

import (
	utils "influxer/scorecard/utilities"
	"strings"
	"sync"
)

type Table struct {
	mut       sync.Mutex
	TableName string
	Headers   []HeaderData
	Matrix    [][]CellData
}

func NewTable(name string) *Table {
	mat := make([][]CellData, 2)
	mat[0] = make([]CellData, 3)
	mat[1] = make([]CellData, 3)

	for i, _ := range mat[0] {
		mat[0][i].Value = []byte("")
		mat[1][i].Value = []byte("")
	}

	t := Table{
		mut:       sync.Mutex{},
		TableName: name,
		Headers:   []HeaderData{{}, {}, {}},
		Matrix:    mat,
	}
	return &t
}

type CellData struct {
	Value []byte
	Type  string
}

type HeaderData struct {
	Value string
}

type HeaderArgs struct {
	ColIndex  int
	TableName string
	CellVal   any
}

type CellArgs struct {
	RowIndex  int
	ColIndex  int
	TableName string
	CellVal   any
	CellType  string
}

func (s *Table) changeCellVal(row int, col int, val string) {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.Matrix[row][col].Value = []byte(val)
}

func (s *Table) changeHeaderVal(col int, val string) {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.Headers[col].Value = val
}

func (s *Table) addRowBottom() {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.Matrix = append(s.Matrix, make([]CellData, len(s.Matrix[0])))
}

// add a row between two indices
// if its negative count back from the end
// it its positive just shift
func (s *Table) addRowAt(newRowIndex int) {
	s.mut.Lock()
	defer s.mut.Unlock()

	if newRowIndex < 0 {
		newRowIndex = len(s.Matrix)
	}

	s.Matrix = utils.InsertAt(s.Matrix, make([]CellData, len(s.Matrix[0])), newRowIndex)
}

func (s *Table) addColumnAt(newColIndex int) {
	s.mut.Lock()
	defer s.mut.Unlock()

	if newColIndex < 0 {
		newColIndex = len(s.Matrix[0])
	}

	for i, _ := range s.Matrix {
		newCell := CellData{
			Value: []byte(""),
			Type:  "string",
		}
		s.Matrix[i] = utils.InsertAt(s.Matrix[i], newCell, newColIndex)
	}

	// add header col
	s.Headers = utils.InsertAt(s.Headers, HeaderData{}, len(s.Headers))
}

func unsignedIndex(index int, length int) int {
	if index < 0 {
		return length + index + 1
	} else {
		return index
	}
}

func (s *Table) UpdateCell(row int, col int, val string, notify func(int, int, bool)) {
	altered := [][]int{}
	tableNotify := false

	if strings.ContainsRune(val, 0x09) || strings.ContainsRune(val, 0x0A) || strings.ContainsRune(val, 0x20) { // is a mutticell copy paste
		tmp, shouldNotify := s.handleMulticellChangeViaPaste(row, col, val)
		tableNotify = shouldNotify
		altered = append(altered, tmp...)
	} else {
		s.changeCellVal(row, col, val)
		altered = append(altered, []int{row, col})
	}

	if tableNotify {
		notify(0, 0, true)
	} else {
		for _, v := range altered {
			notify(v[0], v[1], false)
		}
	}
}

func (s *Table) handleMulticellChangeViaPaste(row int, col int, val string) ([][]int, bool) {
	var triggerTableNotify bool
	updated := [][]int{}
	rows := utils.Split(val, 0x0A, 0x20)

	for rowindex, r := range rows {
		cols := strings.Split(r, string(0x09))
		for colindex, c := range cols {
			if row+rowindex >= len(s.Matrix) {
				s.addRowAt(-1)
				triggerTableNotify = true
			}
			if col+colindex >= len(s.Matrix[0]) {
				s.addColumnAt(-1)
				triggerTableNotify = true
			}
			s.Matrix[row+rowindex][col+colindex] = CellData{Value: []byte(c), Type: "string"}
			updated = append(updated, []int{rowindex + row, colindex + col})
		}
	}
	return updated, triggerTableNotify
}
