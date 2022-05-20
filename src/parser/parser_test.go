package parser

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	const listSize = 470
	requestTestData := map[int]string{
		10:  "http://www.zhenai.com/zhenghun/anyang",
		100: "http://www.zhenai.com/zhenghun/gansu",
		20:  "http://www.zhenai.com/zhenghun/baoji",
		60:  "http://www.zhenai.com/zhenghun/chuxiong",
		90:  "http://www.zhenai.com/zhenghun/fengxian2",
		210: "http://www.zhenai.com/zhenghun/lanzhou",
		355: "http://www.zhenai.com/zhenghun/tongliao",
	}
	dataFile, err := os.Open("parse_data.html")
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(dataFile)
	result := ListParser(content)
	// 测试解析的长度
	if len(result.Requests) != listSize {
		t.Errorf("result size should have %d"+
			",but get %d", listSize, len(result.Requests))
	}
	// random test
	for index, url := range requestTestData {
		if index >= len(result.Requests) {
			t.Errorf("Request index(%d) is overwritten, index scope should be 0-%d",
				index, len(result.Requests)-1)
		}
		actualRequest := result.Requests[index]
		if url != actualRequest.Url {
			t.Errorf("Request expect  is %s"+
				",but actual %s", url, actualRequest.Url)
		}
	}

}
