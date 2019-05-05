package main

import (
	"bufio"
	"fmt"
	"os"
	"pipelinedemo/pipeline"
)

func main() {
	/* v1
	p := pipeline.ArraySource(3, 2, 6, 7, 4)
	for v := range p {
		fmt.Println(v)
	}
	*/

	/*v2
		p := pipeline.InMemSort(
		pipeline.ArraySource(3, 2, 6, 7, 4))
	for v := range p {
		fmt.Println(v)
	}
	/*
	2
	3
	4
	6
	7
	*/

	/*v3
	p := pipeline.Merge(
		pipeline.InMemSort(
			pipeline.ArraySource(
				3, 2, 6, 7, 4)),
		pipeline.InMemSort(
			pipeline.ArraySource(
				7, 4, 0, 3, 2, 8, 13, 8)))

	for v := range p {
		fmt.Println(v)
	}
	*/
	/*
		0
		2
		2
		3
		3
		4
		4
		6
		7
		7
		8
		8
		13
	*/

	// const filename = "small.in"
	// const n = 64

	// const filename = "large.in"
	// const n = 100000000

	const filename = "large.in"
	const n = 100000000

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	p := pipeline.RandomSource(n)

	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, p)
	writer.Flush() // 这里需要看bufio的知识点

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	p = pipeline.ReaderSource(bufio.NewReader(file), -1)
	for v := range p {
		fmt.Println(v)
	}

}

func mergeDemo() {
	p := pipeline.Merge(
		pipeline.InMemSort(
			pipeline.ArraySource(
				3, 2, 6, 7, 4)),
		pipeline.InMemSort(
			pipeline.ArraySource(
				7, 4, 0, 3, 2, 8, 13, 8)))

	for v := range p {
		fmt.Println(v)
	}
}
