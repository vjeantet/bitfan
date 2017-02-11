package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vjeantet/bitfan/processors/doc"
)

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
func docsByKind(kind string) map[string]*doc.Processor {
	data := map[string]*doc.Processor{}
	for name, proc := range plugins[kind] {
		if name == "when" {
			continue
		}
		data[name] = proc().Doc()
	}
	return data
}

func docsByKindByName(kind string, name string) (*doc.Processor, error) {
	if _, ok := plugins[kind][name]; !ok {
		return nil, fmt.Errorf("not found")
	}
	return plugins[kind][name]().Doc(), nil
}
