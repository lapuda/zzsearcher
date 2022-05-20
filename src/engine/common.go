package engine

import (
	"fmt"
	"github.com/jeanphorn/log4go"
	"zzsearcher/src/engine/types"
)

var Logger = log4go.NewDefaultLogger(log4go.FINE)

func worker(r types.Request, needCache bool) (types.ParserResult, error) {
	Logger.Info(fmt.Sprintf("Fetch url:%s", r.Url))
	body, title, err := r.FetchFunc(r.Url, needCache)
	if err != nil {
		return types.ParserResult{}, err
	}
	parseResult := r.ParseFunc(body)
	Logger.Info(fmt.Sprintf("Parse url %s  [%s] \n", r.Url, string(title)))
	return parseResult, nil
}
