package spreadsheet

import (
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

type Scorecard struct {
	tmpl map[string]*template.Template
	data map[string]*Table
	comm *websocket.Conn
}
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

func CreateCellData(rowIndex int, colIndex int, tableName string, cellVal CellData) CellArgs {
	cd := CellArgs{rowIndex, colIndex, tableName, string(cellVal.Value), cellVal.Type}
	return cd
}

func CreateHeaderData(colIndex int, tableName string, headerVal HeaderData) HeaderArgs {
	cd := HeaderArgs{colIndex, tableName, headerVal.Value}
	return cd
}

func NewScorecard() *Scorecard {
	tmpl := make(map[string]*template.Template)

	fm := template.FuncMap{"CreateCellData": CreateCellData, "CreateHeaderData": CreateHeaderData}

	indexTemp := template.New("index")
	indexTemp.Funcs(fm)
	indexFiles, err := indexTemp.ParseFiles("./templates/table.html", "./templates/cell.html", "./templates/headercell.html", "./templates/index.html")
	if err != nil {
		return nil
	}
	tmpl["index.html"] = indexFiles

	tableTemp := template.New("table")
	tableTemp.Funcs(fm)
	tableFiles, err := tableTemp.ParseFiles("./templates/table.html", "./templates/cell.html", "./templates/headercell.html")
	if err != nil {
		return nil
	}
	tmpl["table.html"] = tableFiles

	tmpl["cell.html"] = template.Must(template.ParseFiles("./templates/cell.html"))
	tmpl["headercell.html"] = template.Must(template.ParseFiles("./templates/headercell.html"))

	s := Scorecard{}
	s.tmpl = tmpl
	s.data = make(map[string]*Table)
	s.data["Table_1"] = NewTable("Table_1")

	return &s
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

type TableUpdateEvent struct {
	event    string
	locators []int
}

func parseWhatToken(s string) TableUpdateEvent {
	if strings.HasPrefix(s, "h") { // h0
		// header val
		loc := strings.Replace(s, "h", "", 1)
		iloc, _ := strconv.Atoi(loc)
		return TableUpdateEvent{
			event:    "headerCellValue",
			locators: []int{iloc},
		}

	}
	if strings.HasPrefix(s, "c") { //c0-1
		// cell val
		loc := strings.Replace(s, "c", "", 1)
		locs := strings.Split(loc, "-")
		ilocx, _ := strconv.Atoi(locs[0])
		ilocy, _ := strconv.Atoi(locs[1])
		return TableUpdateEvent{
			event:    "bodyCellValue",
			locators: []int{ilocx, ilocy},
		}
	}
	return TableUpdateEvent{}
}
