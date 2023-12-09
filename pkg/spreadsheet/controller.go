package spreadsheet

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Spreadsheet interface {
	GetSheet() gin.HandlerFunc
	//UpdateCell() gin.HandlerFunc
	//AddRow() gin.HandlerFunc
	HandleLiveComm() gin.HandlerFunc
	ListSheets() gin.HandlerFunc
}

func (s *Scorecard) GetSheet() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("tableName")

		page := s.BuildIndex(tableName)

		c.Header("Content-Type", "text/html")
		write, err := c.Writer.Write(page)
		if err != nil || write < 1 {
			fmt.Println("Failed to write template")
			return
		}
	}
}

func (s *Scorecard) ListSheets() gin.HandlerFunc {
	return func(c *gin.Context) {
		page := s.BuildListIndex()

		c.Header("Content-Type", "text/html")
		write, err := c.Writer.Write(page)
		if err != nil || write < 1 {
			fmt.Println("Failed to write template")
			return
		}
	}
}

//func (s *Scorecard) UpdateCell() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		tableName := c.Param("tableName")
//
//		buf := new(bytes.Buffer)
//		_, err := buf.ReadFrom(c.Request.Body)
//		if err != nil {
//			return
//		}
//
//		updatedDataTags := strings.Split(buf.String(), "=") //c0-1=a //h0=a
//		updatedValue := updatedDataTags[1]
//		whatGotUpdated := parseWhatToken(updatedDataTags[0]) //h1 | c0-1
//
//		var tplBuffer bytes.Buffer
//		if whatGotUpdated.event == "headerCellValue" {
//			s.data[tableName].changeHeaderVal(whatGotUpdated.locators[0], updatedValue)
//			err = s.tmpl["table/headercell"].ExecuteTemplate(&tplBuffer, "table/headercell",
//				HeaderArgs{
//					whatGotUpdated.locators[0],
//					tableName,
//					string(s.data[tableName].Headers[whatGotUpdated.locators[0]].Value)})
//		}
//		if whatGotUpdated.event == "bodyCellValue" {
//			s.data[tableName].changeCellVal(whatGotUpdated.locators[0], whatGotUpdated.locators[1], updatedValue)
//			err = s.tmpl["table/cell"].ExecuteTemplate(&tplBuffer, "table/cell",
//				CellArgs{whatGotUpdated.locators[0],
//					whatGotUpdated.locators[1],
//					tableName,
//					string(s.data[tableName].Matrix[whatGotUpdated.locators[0]][whatGotUpdated.locators[1]].Value),
//					s.data[tableName].Matrix[whatGotUpdated.locators[0]][whatGotUpdated.locators[1]].Type})
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//		}
//
//		c.Header("Content-Type", "text/html")
//		write, err := c.Writer.Write(tplBuffer.Bytes())
//		if err != nil || write < 1 {
//			fmt.Println(err)
//			fmt.Println("Failed to write template")
//			return
//		}
//	}
//}

//func (s *Scorecard) AddRow() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		tableName := c.Param("tableName")
//
//		s.data[tableName].addRowBottom()
//
//		var tplBuffer bytes.Buffer
//		err := s.tmpl["table/table"].ExecuteTemplate(&tplBuffer, "table/table", s.data[tableName])
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//
//		c.Header("Content-Type", "text/html")
//		write, err := c.Writer.Write(tplBuffer.Bytes())
//		if err != nil || write < 1 {
//			fmt.Println("Failed to write template")
//			return
//		}
//	}
//}
