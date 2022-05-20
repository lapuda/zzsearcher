package util

import (
	"fmt"
	"github.com/jeanphorn/log4go"
	"io/ioutil"
	"net/http"
)

var Logger = log4go.NewDefaultLogger(log4go.FINE)

func HttpGetContent(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		message := fmt.Sprintf(
			"Request Body Exception! code is : %d", resp.StatusCode)
		Logger.Error(message)
		panic(message)
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(all)
}
