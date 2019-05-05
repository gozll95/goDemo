package pipeline

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

var startTime time.Time

func Init() {
	startTime = time.Now()
}

func ArraySource(a ...int) <-chan int {
	//out := make(chan int)
	// 优化,增加buffer
	out := make(chan int, 1024)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

func InMemSort(in <-chan int) <-chan int {
	//out := make(chan int)
	// 优化,增加buffer
	out := make(chan int, 1024)
	go func() {
		// Read into memory
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read done:", time.Now().Sub(startTime))

		// Sort
		sort.Ints(a)
		fmt.Println("InMemSort done:", time.Now().Sub(startTime))

		// Output
		for _, v := range a {
			out <- v
		}
		close(out)
		fmt.Println("Merge done:", time.Now().Sub(startTime))
	}()
	return out
}

func Merge(int1, int2 <-chan int) <-chan int {
	//out := make(chan int)
	// 优化,增加buffer
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-int1
		v2, ok2 := <-int2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-int1
			} else {
				out <- v2
				v2, ok2 = <-int2
			}
		}
		close(out)
	}()
	return out
}

// func ReaderSource(reader io.Reader) <-chan int {
// 	out := make(chan int)
// 	go func() {
// 		buffer := make([]byte, 8)
// 		for {
// 			n, err := reader.Read(buffer)
// 			// 第一个n表示读了多少个字节,err表示读了EOF之类的
// 			if n > 0 {
// 				v := int(binary.BigEndian.Uint64(buffer))
// 				out <- v
// 			}
// 			if err != nil {
// 				break
// 			}
// 		}
// 		close(out)
// 	}()
// 	return out
// }

func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(
			buffer, uint64(v))
		writer.Write(buffer)
	}
}

func RandomSource(count int) <-chan int {
	//out := make(chan int)
	// 优化,增加buffer
	out := make(chan int, 1024)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

// 多路的两两归并
func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}

	m := len(inputs) / 2
	// merge inputs[0..m] and inputs[m..end]
	return Merge(
		MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}

// Reader的分块
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	//out := make(chan int)
	// 优化,增加buffer
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		bytesRead := 0
		for {
			n, err := reader.Read(buffer)
			bytesRead += n
			// 第一个n表示读了多少个字节,err表示读了EOF之类的
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil ||
				(chunkSize != -1 && bytesRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}
