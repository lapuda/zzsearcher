package store

import (
	"context"
	"zzsearcher/src/engine/types"
)

var store DataStore
var dbName = "manga007"
var dbUri = "mongodb://root:123456@mongodb.269r.com:30001,mongodb.269r.com:30002/?replicaSet=wy_repl2&authSource=admin"


func ItemSaver(ctx context.Context) chan map[types.ParserType]interface{}  {
	store = NewMongoDBStore(ctx,dbName, dbUri)
	out := make(chan map[types.ParserType]interface{})
	go func() {
		item := <-out
		for key,item := range item {
			switch key {
			case types.Manga:
				store.Save(item, "manga")
			case types.Chapter:
				store.Save(item, "charpter")
			}
		}
	}()
	return out
}
