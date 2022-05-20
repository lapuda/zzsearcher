package util

import (
	"regexp"
	"strings"
)

// 内容匹配
func ContentMatch(reg string, data []byte) string {
	re := regexp.MustCompile(reg)
	find := re.FindSubmatch(data)
	if len(find) == 0 {
		return ""
	}
	content := string(find[0])
	if len(find) > 1 {
		content = strings.Replace(string(find[1]), "\n", "", -1)
	}
	return content
}

// 内容匹配
func ContentSecondMatch(reg string, data []byte) string {
	re := regexp.MustCompile(reg)
	findSecond := re.FindSubmatch(data)
	if len(findSecond) > 1 {
		return string(findSecond[1])
	}
	return ""
}
