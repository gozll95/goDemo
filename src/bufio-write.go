/*
type Writer struct { ... }

// NewWriterSize 将 wr 封装成一个带缓存的 bufio.Writer 对象，
// 缓存大小由 size 指定（如果小于 4096 则会被设置为 4096）。
// 如果 wr 的基类型就是有足够缓存的 bufio.Writer 类型，则直接将
// wr 转换为基类型返回。
func NewWriterSize(wr io.Writer, size int) *Writer

// NewWriter 相当于 NewWriterSize(wr, 4096)
func NewWriter(wr io.Writer) *Writer

// bufio.Writer 实现了如下接口：
// io.Writer
// io.ReaderFrom
// io.ByteWriter

// WriteString 功能同 Write，只不过写入的是字符串
func (b *Writer) WriteString(s string) (int, error)

// WriteRune 向 b 写入 r 的 UTF-8 编码，返回 r 的编码长度。
func (b *Writer) WriteRune(r rune) (size int, err error)

// Flush 将缓存中的数据提交到底层的 io.Writer 中
func (b *Writer) Flush() error

// Available 返回缓存中未使用的空间的长度
func (b *Writer) Available() int

// Buffered 返回缓存中未提交的数据的长度
func (b *Writer) Buffered() int

// Reset 将 b 的底层 Writer 重新指定为 w，同时丢弃缓存中的所有数据，复位
// 所有标记和错误信息。相当于创建了一个新的 bufio.Writer。
func (b *Writer) Reset(w io.Writer)
*/

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	buf := bufio.NewWriterSize(os.Stdout, 0)
	fmt.Println(buf.Available(), buf.Buffered())

	buf.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	fmt.Println(buf.Available(), buf.Buffered()) // 4070 26

	//缓存后统一输出，避免中断频繁刷新，影响速度
	buf.Flush() // ABCDEFGHIJKLMNOPQRSTUVWXYZ
}
