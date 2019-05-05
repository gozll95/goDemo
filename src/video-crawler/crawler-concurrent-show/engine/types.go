package engine

type Request struct {
	Url string
	//ParserFunc ParserFunc
	ParserFunc Parser
}

// 加上Parser interface
type Parser interface {
	Parse(contents []byte, url string) ParseResult
	Serialize() (name string, args interface{})
}

// type SerializerParser struct {
// 	Name string
// 	Args interface{}
// }

// {"ParseCityList",nil},{"ProfileParser",userName} // func(contents []byte, url string) 这个不用写

type ParseResult struct {
	Requests []Request
	Items    []Item
}

type Item struct {
	Url     string
	Id      string
	Type    string
	Payload interface{}
}

// func NilParser([]byte) ParseResult {
// 	return ParseResult{}
// }

// 添加
type NilParser struct{}

func (NilParser) Parse(_ []byte, _ string) ParseResult {
	return ParseResult{}
}
func (NilParser) Serialize() (name string, args interface{}) {
	return "NilParser", nil
}

// 工厂函数
type ParserFunc func(contents []byte, url string) ParseResult

type FuncParser struct {
	parser ParserFunc
	name   string
}

func NewFuncParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		name:   name,
	}
}

func (f FuncParser) Parse(contents []byte, url string) ParseResult {
	return f.parser(contents, url)

}
func (f FuncParser) Serialize() (name string, args interface{}) {
	return f.name, nil
}
