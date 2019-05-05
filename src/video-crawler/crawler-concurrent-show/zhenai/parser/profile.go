package parser

import (
	"regexp"
	"strconv"
	"video-crawler/crawler-concurrent-show/engine"
	"video-crawler/crawler-concurrent-show/model"
)

var nameRe = regexp.MustCompile(`<td><span class="label">年龄：</span>([\d]+)岁</td>`)
var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>([\d]+)岁</td>`)
var heightRe = regexp.MustCompile(`<td><span class="label">身高：</span>([\d]+)CM</td>`)
var weightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">([\d]+)</span></td>`)

var incomeRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var genderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var hukouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var educationRe = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var hoseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var carRe = regexp.MustCompile(` <td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)
var xingzuo = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)

var guessRe = regexp.MustCompile(`<a class="exp-user-name"[^>]*href="(http://album.zhenai.com/u/[\d]+)">([^<]+)</a>`)

var idUrlre = regexp.MustCompile(`http://album.zhenai.com/u/([\d]+)`)

func parseProfile(contents []byte, name string, url string) engine.ParseResult {
	profile := model.Profile{}
	profile.Name = name
	profile.Age = extractInt(contents, ageRe)
	profile.Height = extractInt(contents, heightRe)
	profile.Weight = extractInt(contents, weightRe)
	profile.Income = extractString(contents, incomeRe)
	profile.Marriage = extractString(contents, marriageRe)
	profile.Gender = extractString(contents, genderRe)
	profile.Hukou = extractString(contents, hukouRe)
	profile.House = extractString(contents, hoseRe)
	profile.Car = extractString(contents, carRe)
	profile.Xingzuo = extractString(contents, xingzuo)
	id := extractString([]byte(url), idUrlre)
	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "zhenai",
				Id:      id,
				Payload: profile,
			},
		},
	}
	matches := guessRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		url := string(m[1])
		name := string(m[2])
		result.Requests = append(result.Requests, engine.Request{
			Url: url,
			// ParserFunc: ProfileParser(name),
			ParserFunc: NewProfileParser(name),
		})
	}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	}
	return ""
}

func extractInt(contents []byte, re *regexp.Regexp) int {
	i, err := strconv.Atoi(extractString(contents, re))
	if err == nil {
		return i
	}
	return 0
}

type ProfileParser struct {
	userName string
}

func NewProfileParser(name string) *ProfileParser {
	return &ProfileParser{
		userName: name,
	}
}

func (p ProfileParser) Parse(contents []byte, url string) ParseResult {
	return parseProfile(contents, url, p.userName)
}

func (p ProfileParser) Serialize() (name string, args interface{}) {
	return "ProfileParser", p.userName
}

// func ProfileParser(name string) engine.ParserFunc {
// 	return func(c []byte, url string) engine.ParseResult {
// 		return parseProfile(c, name, url)
// 	}
// }
