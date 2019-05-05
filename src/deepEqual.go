package main

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func test() {
	got := strings.Split("a:b:c", ":")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		fmt.Println("!=")
	} else {
		fmt.Println("==")
	}
}

func validateEqualArgs(expected, actual interface{}) error {
	if isFunction(expected) || isFunction(actual) {
		return errors.New("cannot take func type as argument")
	}
	return nil
}

func isFunction(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}

func ObjectsAreEqual(expected, actual interface{}) bool {

	if expected == nil || actual == nil {
		return expected == actual
	}
	if exp, ok := expected.([]byte); ok {
		act, ok := actual.([]byte)
		if !ok {
			return false
		} else if exp == nil || act == nil {
			return exp == nil && act == nil
		}
		return bytes.Equal(exp, act)
	}
	return reflect.DeepEqual(expected, actual)

}

func myEqual(expected, actual interface{}) (isSame bool) {
	if err := validateEqualArgs(expected, actual); err != nil {
		panic(err)
	}
	return ObjectsAreEqual(expected, actual)
}

type A struct {
	AA int
	BB string
}

// func testCase1() {
// 	got := strings.Split("a:b:c", ":")
// 	want := []string{"a", "b", "c"}
// 	isSame, err := myEqual(got, want)
// 	fmt.Println(isSame, err)

// }

func main() {
	fmt.Println(myEqual([]int{1, 2, 3}, []int{1, 2, 3}))        // "true"
	fmt.Println(myEqual([]string{"foo"}, []string{"bar"}))      // "false"
	fmt.Println(myEqual([]string(nil), []string{}))             // "true"
	fmt.Println(myEqual(map[string]int(nil), map[string]int{})) // "true"

	a := A{
		AA: 1,
		BB: "1",
	}

	b := A{
		AA: 1,
		BB: "1",
	}
	fmt.Println(myEqual(a, b)) // "true"
}
