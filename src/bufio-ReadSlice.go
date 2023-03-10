package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	// 尾部有换行标记
	buf := bufio.NewReaderSize(strings.NewReader("ABCDEFG\n"), 0)

	for line, err := []byte{0}, error(nil); len(line) > 0 && err == nil; {
		line, err = buf.ReadSlice('\n')
		fmt.Printf("%q   %v\n", line, err)
	}
	// "ABCDEFG\n"   <nil>
	// ""   EOF

	fmt.Println("----------")

	// 尾部没有换行标记
	buf = bufio.NewReaderSize(strings.NewReader("ABCDEFG"), 0)

	for line, err := []byte{0}, error(nil); len(line) > 0 && err == nil; {
		line, err = buf.ReadSlice('\n')
		fmt.Printf("%q   %v\n", line, err)
	}
	// "ABCDEFG"   EOF
}
