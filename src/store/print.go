package store

import (
	"context"
	"encoding/json"
	"fmt"
	"zzsearcher/src/engine/types"
)

type PrintStore struct {
}

func (receiver *PrintStore) Save(m interface{}, collection string) {
	if data, error := json.Marshal(m); error != nil {
		fmt.Printf("Table:%s save data is:%v \n", collection, m)
	} else {
		fmt.Printf("Table:%s save data is:%v \n", collection, string(data))
	}

}

func NewPrintStore() *PrintStore {
	return &PrintStore{}
}

func ItemPrinter(ctx context.Context) chan map[types.ParserType]interface{} {
	store = NewPrintStore()
	out := make(chan map[types.ParserType]interface{})
	go func() {
		for {
			select {
			case items := <-out:
				for key, item := range items {
					switch key {
					case types.Manga:
						store.Save(item, "manga")
					case types.Chapter:
						store.Save(item, "charpter")
					default:
						store.Save(item, "default")
					}
				}
			}
		}
	}()
	return out
}
