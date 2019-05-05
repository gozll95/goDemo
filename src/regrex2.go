package main

import (
	"fmt"
	"regexp"
)

const text = `
                    "status":"active",
                    "updated_at":"2018-07-19T11:04:22.401+08:00",
                    "created_at":"2018-07-19T11:04:22.401+08:00",
					"created_at":"2018-07-19T11:04:22.401+08:00",
                    "freezed_status":"",
                    "freezed_at":"0001-01-01T00:00:00Z",
                    "blocked_status":"",
                    "blocked_at":"0001-01-01T00:00:00Z",
                    "is_restricted":false
`

const sample = `aabproductsaaa`

func main() {
	re := regexp.MustCompile("([0-9]{4}-[0-9]{2}-[0-9]{2})T(\\d+:\\d+):\\d+.\\d+(\\+08:00)")
	match := re.FindAllString(text, -1)
	fmt.Println(match)

	match2 := re.FindAllStringSubmatch(text, -1)
	fmt.Println(match2)

	for _, m := range match2 {
		fmt.Println(m)
	}

	fmt.Printf("%s", re.ReplaceAllString(text, "$1 $2"))

}
