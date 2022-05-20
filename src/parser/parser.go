package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jeanphorn/log4go"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zzsearcher/src/engine/types"
	"zzsearcher/src/fatcher"
	"zzsearcher/src/model"
	"zzsearcher/src/util"
)

var Logger = log4go.NewDefaultLogger(log4go.FINE)

type ParseRule struct {
	ContentRule string `json:"contentRule"`
	SeedUrlRule string `json:"seedUrlRule"`
	ParseFunc   string `json:"parse_func"`
	FetchFunc   string `json:"fetch_func"`
}

// 章节名称  string
// 更新时间  string
// 作者     string
// 类别     string
// 状态     string
// 描述     string
// 封面     string

type MangaRule struct {
	*ParseRule
	Name       string `json:"name"`
	UpdateDate string `json:"update_date"`
	Author     string `json:"author"`
	Type       string `json:"type"`
	MangaId    string `json:"manga_id"`
	Status     string `json:"status"`
	Describe   string `json:"describe"`
	Cover      string `json:"cover"`
}

type ChapterRule struct {
	*ParseRule
	Name      string `json:"name"`
	MangaId   string `json:"manga_id"`
	ChapterId string `json:"chapter_id"`
	ImageBase string `json:"image_base"`
}

type SeedRule struct {
	BaseUrl           string            `json:"base_url"`
	SeedListRule      ParseRule         `json:"seed_list_rule"`
	MangaListRule     ParseRule         `json:"manga_list_rule"`
	Manga             MangaRule         `json:"manga_rule"`
	ChapterRule       ChapterRule       `json:"chapter_rule"`
	ChapterListRule   ParseRule         `json:"chapter_list_rule"`
	Base64ChapterRule Base64ChapterRule `json:"base_64_chapter_rule"`
}

type Base64ChapterRule struct {
	*ParseRule
	Name          string `json:"name"`
	MangaId       string `json:"manga_id"`
	ChapterId     string `json:"chapter_id"`
	ImageBase     string `json:"image_base"`
	Base64Content string `json:"base64_content"`
}

var parserRule SeedRule

func init() {
	f, err := os.Open("parse_rule.json")
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(f).Decode(&parserRule)
	if err != nil {
		panic(err)
	}
}

func NilParser([]byte) types.ParserResult {
	return types.ParserResult{}
}

func Parser(contents []byte, rule ParseRule, baseUrl string) types.ParserResult {
	Logger.Info("parse start")
	// 第一步拿到区域内容
	secondContent := contents
	if rule.ContentRule != "" {
		Logger.Info("content rule is not empty,go to parse content start....")
		re := regexp.MustCompile(rule.ContentRule)
		find := re.FindAllSubmatch(contents, -1)
		Logger.Info("Parse is complete，but match nothing!")
		if len(find) == 0 {
			Logger.Info("Parse is complete，but match nothing!")
			return types.ParserResult{}
		}
		secondContent = find[0][0]
	}
	// 第二步拿到这个页面的种子
	a := regexp.MustCompile(rule.SeedUrlRule)
	findSecond := a.FindAllSubmatch(secondContent, -1)
	if len(findSecond) == 0 {
		println("没有匹配到任何的种子")
	}
	result := types.ParserResult{}
	for index, submatch := range findSecond {
		result.Items = append(result.Items, string(submatch[2]))
		result.Requests = append(
			result.Requests,
			types.Request{
				baseUrl + string(submatch[1]),
				ParserFunc(rule.ParseFunc),
				fatcher.FetcherFunc(rule.FetchFunc),
				index + 1,
			},
		)
	}
	Logger.Info("Parse complete!")
	return result
}

func SeedListParser(contents []byte) types.ParserResult {
	Logger.Info("SeedListParser Parse start....")
	return Parser(contents, parserRule.SeedListRule, parserRule.BaseUrl)
}

func MangaListParser(contents []byte) types.ParserResult {
	Logger.Info("MangaListParser Parse start....")
	return Parser(contents, parserRule.MangaListRule, parserRule.BaseUrl)
}

func MangaParser(contents []byte) types.ParserResult {
	Logger.Info("MangaParser Parse start....")
	rule := parserRule.Manga
	model := model.Manga{}
	model.Name = util.ContentMatch(rule.Name, contents)
	model.Author = util.ContentMatch(rule.Author, contents)
	model.Status = "更新中"
	model.Type = util.ContentMatch(rule.Type, contents)
	model.Cover = util.ContentMatch(rule.Cover, contents)
	model.UpdateDate = util.ContentMatch(rule.UpdateDate, contents)
	model.Describe = util.ContentMatch(rule.Describe, contents)
	model.CreateTime = time.Now().UnixNano() / 1000
	model.MangaID, _ = strconv.Atoi(util.ContentMatch(rule.MangaId, contents))

	secondContent := contents
	// 第一步拿到区域内容
	if rule.ContentRule != "" {
		re := regexp.MustCompile(rule.ContentRule)
		find := re.FindAllSubmatch(contents, -1)
		if len(find) < 0 {
			return types.ParserResult{}
		}
		secondContent = find[0][0]
	}

	// 第二步拿到
	a := regexp.MustCompile(rule.SeedUrlRule)
	findSecond := a.FindAllSubmatch(secondContent, -1)
	result := types.ParserResult{}
	result.ParserType = types.Manga
	result.Items = append(result.Items, model)

	for index, submatch := range findSecond {
		result.Requests = append(
			result.Requests,
			types.Request{
				parserRule.BaseUrl + string(submatch[1]),
				ParserFunc(rule.ParseFunc),
				fatcher.FetcherFunc(rule.FetchFunc),
				index + 1,
			},
		)
	}
	return result
}

func ChapterParser(contents []byte) types.ParserResult {
	Logger.Info("ChapterParser Parse start....")
	rule := parserRule.ChapterRule
	model := model.Chapter{}

	secondContent := contents
	// 第一步拿到区域内容
	if rule.ContentRule != "" {
		re := regexp.MustCompile(rule.ContentRule)
		find := re.FindAllSubmatch(contents, -1)
		if len(find) < 0 {
			return types.ParserResult{}
		}
		secondContent = find[0][0]
	}
	// 第二步拿到
	a := regexp.MustCompile(rule.SeedUrlRule)
	findSecond := a.FindAllSubmatch(secondContent, -1)
	model.MangaID, _ = strconv.Atoi(util.ContentMatch(rule.MangaId, contents))
	model.ChapterName = util.ContentMatch(rule.Name, contents)
	model.ChapterID, _ = strconv.Atoi(util.ContentMatch(rule.ChapterId, contents))
	model.UpdateTime = time.Now().UnixNano() / 1000

	result := types.ParserResult{}

	var images []string
	for _, submatch := range findSecond {
		imgURL := strings.Replace(string(submatch[1]), "\\", "", -1)
		images = append(images, rule.ImageBase+imgURL)
	}
	if len(images) > 0 {
		model.Images = strings.Join(images, ",")
	}
	result.ParserType = types.Chapter
	result.Items = append(result.Items, model)
	return result
}

func Base64ChapterParser(contents []byte) types.ParserResult {
	Logger.Info("Base64ChapterParser Parse start....")
	rule := parserRule.Base64ChapterRule
	model := model.Chapter{}
	secondContent := contents
	// 第一步拿到区域内容
	if rule.ContentRule != "" {
		re := regexp.MustCompile(rule.ContentRule)
		find := re.FindAllSubmatch(contents, -1)
		if len(find) < 0 {
			return types.ParserResult{}
		}
		secondContent = find[0][0]
	}
	// 第二步拿到
	model.MangaID, _ = strconv.Atoi(util.ContentMatch(rule.MangaId, secondContent))
	model.ChapterName = util.ContentMatch(rule.Name, secondContent)
	model.ChapterID, _ = strconv.Atoi(util.ContentMatch(rule.ChapterId, secondContent))
	imageBase := util.ContentSecondMatch(rule.Base64Content, secondContent)

	decoded, err := base64.StdEncoding.DecodeString(imageBase)
	if err != nil {
		fmt.Println("decode error:", err)
		panic(err)
	}

	var jsonImages []string
	err = json.Unmarshal(decoded, &jsonImages)
	if err != nil {
		println(err)
	}

	var images []string
	for _, str := range jsonImages {
		uri := strings.Split(str, "|")
		if len(uri) >= 2 {
			images = append(images, strings.TrimSpace(uri[1]))
		} else {
			println(model.ChapterID, model.ChapterName, model.MangaID)
		}
	}
	if len(images) > 0 {
		model.Images = strings.Join(images, ",")
	}

	model.UpdateTime = time.Now().UnixNano() / 1000

	result := types.ParserResult{}
	result.ParserType = types.Chapter
	result.Items = append(result.Items, model)
	return result
}

func ChapterListParser(contents []byte) types.ParserResult {
	Logger.Info("ChapterListParser Parse start....")
	return Parser(contents, parserRule.ChapterListRule, parserRule.BaseUrl)
}

func ZeroParser(contents []byte) types.ParserResult {
	Logger.Info("ZeroParser Parse start....")
	return types.ParserResult{}
}

func ParserFunc(funcName string) func(contents []byte) types.ParserResult {
	switch funcName {
	case "seed_list_parser":
		return SeedListParser
	case "manga_list_parser":
		return MangaListParser
	case "manga_parser":
		return MangaParser
	case "chapter_list_parser":
		return ChapterListParser
	case "chapter_parser":
		return ChapterParser
	case "base64_chapter_parser":
		return Base64ChapterParser
	case "":
		return ZeroParser
	default:
		panic("func " + funcName + " not support")
	}
}
