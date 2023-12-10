package spreadsheet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"text/template"
)

type Scorecard struct {
	tmpl      map[string]*template.Template
	data      map[string]*Table
	tableConn map[string][]*websocket.Conn
	homeConn  []*websocket.Conn
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
				notify := func(row int, col int, tNot bool) {
					if tNot {
						s.notifyTable(tableName)
					} else {
						s.notifyCells(tableName, row, col, nil) //nil ->  need to notify self
					}
				}
				s.data[tableName].UpdateCell(event.locators[0], event.locators[1], v.(string), notify)
			}
			if event.event == "headerCellValue" {
				s.data[tableName].changeHeaderVal(event.locators[0], v.(string))
				s.notifyHeaderCell(tableName, event.locators[0])
			}
			if event.event == "addRow" {
				s.data[tableName].addRowAt(event.locators[0])
				s.notifyTable(tableName)
			}
			if event.event == "addColumn" {
				s.data[tableName].addColumnAt(event.locators[0])
				s.notifyTable(tableName)
			}
		}
	}
}

func (s *Scorecard) notifyCells(tableName string, row int, col int, self *websocket.Conn) {
	for _, conn := range s.tableConn[tableName] {
		if conn == self {
			continue // dont echo messages to yourself, that way we can keep focus on the <input>
		}
		conn.WriteMessage(websocket.TextMessage, s.BuildCell(tableName, row, col))
	}
}

func (s *Scorecard) notifyHeaderCell(tableName string, col int) {
	for _, conn := range s.tableConn[tableName] {
		conn.WriteMessage(websocket.TextMessage, s.BuildHeaderCell(tableName, col))
	}
}

func (s *Scorecard) notifyTable(tableName string) {
	returnTable := s.BuildTable(tableName)
	for _, conn := range s.tableConn[tableName] {
		err := conn.WriteMessage(websocket.TextMessage, returnTable)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Scorecard) notifyList() {
	returnList := s.BuildTableList()
	for _, conn := range s.homeConn {
		err := conn.WriteMessage(websocket.TextMessage, returnList)
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

func CreateCellData(rowIndex int, colIndex int, tableName string, cellVal CellData) CellArgs {
	cd := CellArgs{rowIndex, colIndex, tableName, string(cellVal.Value), cellVal.Type}
	return cd
}

func CreateHeaderData(colIndex int, tableName string, headerVal HeaderData) HeaderArgs {
	cd := HeaderArgs{colIndex, tableName, headerVal.Value}
	return cd
}

func parseTemplateFiles(names ...string) (*template.Template, error) {
	return template.ParseFiles(names...)
}

func NewScorecard() *Scorecard {
	tmpl := make(map[string]*template.Template)

	fm := template.FuncMap{"CreateCellData": CreateCellData, "CreateHeaderData": CreateHeaderData}

	parseTemplates := func(templateName string, files ...string) {
		t := template.New(templateName)
		t.Funcs(fm)
		tmpl[templateName] = t
		parsedFiles, err := tmpl[templateName].ParseFiles(files...)
		if err != nil {
			return
		}
		tmpl[templateName] = parsedFiles
	}

	parseTemplates("table/index", "./templates/views/table/table.html", "./templates/views/table/cell.html", "./templates/views/table/headercell.html", "./templates/views/table/index.html")
	parseTemplates("table/table", "./templates/views/table/table.html", "./templates/views/table/cell.html", "./templates/views/table/headercell.html")
	parseTemplates("list/index", "./templates/views/list/tablelist.html", "./templates/views/list/index.html")
	parseTemplates("list/tablelist", "./templates/views/list/tablelist.html")

	tmpl["table/cell"] = template.Must(template.ParseFiles("./templates/views/table/cell.html"))
	tmpl["table/headercell"] = template.Must(template.ParseFiles("./templates/views/table/headercell.html"))

	s := Scorecard{}
	s.tmpl = tmpl
	s.data = make(map[string]*Table)
	s.data["default_table"] = NewTable("default_table")

	s.tableConn = make(map[string][]*websocket.Conn)

	return &s
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
	if strings.HasPrefix(s, "ac") { //
		return TableUpdateEvent{event: "addColumn", locators: []int{-1}}
	}
	return TableUpdateEvent{}
}
