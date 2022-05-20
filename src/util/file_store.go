package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func FileStore(file string, content []byte) {
	err := ioutil.WriteFile(file, content, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

func FileReStore(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return data
}

func MD5(source string) string {
	hash := md5.New()
	hash.Write([]byte(source))
	return hex.EncodeToString(hash.Sum(nil))
}
