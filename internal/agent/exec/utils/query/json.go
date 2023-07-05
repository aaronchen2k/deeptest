package queryUtils

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/antchfx/jsonquery"
	"strings"
)

func JsonQuery(content string, expression string) (result string) {
	doc, err := jsonquery.Parse(strings.NewReader(content))
	if err != nil {
		result = consts.ContentErr
		return
	}
	elem, err := jsonquery.Query(doc, expression)

	if err != nil || elem == nil {
		result = consts.ExtractorErr
		return
	}

	obj := elem.Value()

	switch obj.(type) {
	case string:
		result = obj.(string)
	default:
		bytes, err := json.Marshal(obj)
		if err != nil {
			result = err.Error()
		} else {
			result = string(bytes)
		}
	}

	return
}
