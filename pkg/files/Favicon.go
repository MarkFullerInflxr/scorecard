package files

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

type Favicon struct {
}

func NewFiles() *Favicon {
	f := Favicon{}
	return &f
}

func (f *Favicon) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileName := c.Param("file")
		b, err := os.ReadFile("./templates/public" + string(os.PathSeparator) + strings.ReplaceAll(fileName, string(os.PathSeparator), ""))
		if err != nil {
			return
		}

		c.Header("Content-Type", "image/png")
		c.Writer.Write(b)
	}
}
