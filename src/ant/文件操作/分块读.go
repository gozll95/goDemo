import (
	"encoding/binary"
	"io"
)

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
