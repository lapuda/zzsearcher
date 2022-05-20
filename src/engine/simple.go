package engine

import (
	"fmt"
	"zzsearcher/src/engine/types"
	"zzsearcher/src/store"
)

type SimpleEngine struct {
	Store store.DataStore
	mode  types.MODE
}

func (s *SimpleEngine) Run(seeds ...types.Request) {
	var requests []types.Request = seeds
	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]
		parseResult, err := worker(r, false)
		if err != nil {
			Logger.Info(err)
			continue
		}
		requests = append(requests, parseResult.Requests...)
		if len(parseResult.Items) > 0 {
			for _, item := range parseResult.Items {
				if s.mode == types.TEST {
					fmt.Printf("item : %v", item)
					break
				}
				switch parseResult.ParserType {
				case types.Manga:
					s.Store.Save(item, "manga")
				case types.Chapter:
					s.Store.Save(item, "charpter")
				}
			}
		}
		fmt.Printf("length :%d,url: %s \n", len(requests), r.Url)
	}
}
