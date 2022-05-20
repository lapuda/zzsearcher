package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jeanphorn/log4go"
	"net/http"
	"zzsearcher/src/engine/types"
)

var Logger = log4go.NewDefaultLogger(log4go.FINE)

type RemoteStore struct {
}

func (r *RemoteStore) Save(m interface{}, collection string) {
	if data, error := json.Marshal(m); error != nil {
		fmt.Printf("Table:%s save data is:%v \n", collection, m)
	} else {
		switch collection {
		case "manga":
			r.Post("/comic/create", data)
		case "chapter":
			r.Post("/chapter/create", data)
		}
	}

}

func (r *RemoteStore) Post(api string, data []byte) {
	Logger.Info(fmt.Sprintf("Request API:%s data:%s", api, string(data)))
	rsp, error := http.Post("http://127.0.0.1:8081"+api, "application/json", bytes.NewBuffer(data))
	if error != nil {
		println(error)
		return
	}
	var rdata []byte = make([]byte, 1024)
	rsp.Body.Read(rdata)
	Logger.Info(string(rdata))
}

func NewRemoteStore() *RemoteStore {
	return &RemoteStore{}
}

func ItemRemote(ctx context.Context) chan map[types.ParserType]interface{} {
	store = NewRemoteStore()
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
						store.Save(item, "chapter")
					default:
						store.Save(item, "default")
					}
				}
			}
		}
	}()
	return out
}
