package spreadsheet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

type Scorecard struct {
	tmpl map[string]*template.Template
	data map[string]*Table
	conn map[string][]*websocket.Conn
}

func (s *Scorecard) handleEvent(meta map[string]interface{}, conn *websocket.Conn) {
	headersInterface, ok := meta["HEADERS"].(map[string]interface{})
	if !ok {
		fmt.Println("malformed headers")
		return
	}

	currentUrlI := headersInterface["HX-Current-URL"].(string)

	tableName := getTableName(currentUrlI)

	for k, v := range meta { //c0-1,d
		if k != "HEADERS" {
			event := parseWhatToken(k)

			if event.event == "bodyCellValue" {
				s.data[tableName].changeCellVal(event.locators[0], event.locators[1], v.(string))
				s.notifyCells(tableName, event.locators[0], event.locators[1], conn)
			}
			if event.event == "headerCellValue" {
				s.data[tableName].changeHeaderVal(event.locators[0], v.(string))
				s.notifyHeaderCell(tableName, event.locators[0])
			}
			if event.event == "addRow" {
				s.data[tableName].addRowAt(event.locators[0])
				s.notifyTable(tableName)
			}

		}
	}
}

func (s *Scorecard) notifyCells(tableName string, row int, col int, self *websocket.Conn) {
	for _, conn := range s.conn[tableName] {
		if conn == self {
			continue // dont echo messages to yourself, that way we can keep focus on the <input>
		}
		conn.WriteMessage(websocket.TextMessage, s.BuildCell(tableName, row, col))
	}
}

func (s *Scorecard) notifyHeaderCell(tableName string, col int) {
	for _, conn := range s.conn[tableName] {
		conn.WriteMessage(websocket.TextMessage, s.BuildHeaderCell(tableName, col))
	}
}

func (s *Scorecard) notifyTable(tableName string) {
	returnTable := s.BuildTable(tableName)
	for _, conn := range s.conn[tableName] {
		err := conn.WriteMessage(websocket.TextMessage, returnTable)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getTableName(s string) string {
	lastIndex := strings.LastIndex(s, "/")

	if lastIndex != -1 {
		substringAfterLastSlash := s[lastIndex+1:]
		return substringAfterLastSlash
	} else {
		fmt.Println("Failed to parse headers")
		return ""
	}
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

	s.conn = make(map[string][]*websocket.Conn)

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

// add a row between two indices
// if its negative count back from the end
// it its positive just shift
func (s *Table) addRowAt(i int) {
	s.mut.Lock()
	defer s.mut.Unlock()

	if i > 0 {
		newRow := make([]CellData, len(s.Matrix[0]))
		s.Matrix = append(s.Matrix[:i+1], nil)
		copy(s.Matrix[i+1:], s.Matrix[i:])
		s.Matrix[i] = newRow
	} else {
		/*
			-3|0: [],
			-2|1: []
			-1|2: nil
		*/
		width := len(s.Matrix[0]) //1
		height := len(s.Matrix)
		newRow := make([]CellData, width)

		bottomEnd := [][]CellData{}
		copy(bottomEnd, s.Matrix[height+i:])

		s.Matrix = append(s.Matrix[:height+i+1], nil) // s[0:2] -> only the first row
		height = len(s.Matrix)                        //3
		copy(s.Matrix[height+i+1:], bottomEnd)        //s[1:] , s[2-1-1:]
		s.Matrix[height+i] = newRow
	}
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

	if strings.HasPrefix(s, "ar") { //
		return TableUpdateEvent{event: "addRow", locators: []int{-1}}
	}
	return TableUpdateEvent{}
}
