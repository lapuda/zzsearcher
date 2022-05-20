package types

type ParserType int

type MODE int8

const (
	SERVER = iota
	TEST
)

const (
	Seed = iota
	MangaList
	Manga
	Chapter
)

type Request struct {
	Url       string
	ParseFunc func([]byte) ParserResult
	FetchFunc func(string, bool) ([]byte, []byte, error)
	Seq       int
}

type ParserResult struct {
	Requests   []Request
	ParserType ParserType
	Items      []interface{}
}
