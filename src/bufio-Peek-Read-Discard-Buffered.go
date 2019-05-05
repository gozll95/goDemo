/*
// bufio 包实现了带缓存的 I/O 操作

------------------------------------------------------------

type Reader struct { ... }

// NewReaderSize 将 rd 封装成一个带缓存的 bufio.Reader 对象，
// 缓存大小由 size 指定（如果小于 16 则会被设置为 16）。
// 如果 rd 的基类型就是有足够缓存的 bufio.Reader 类型，则直接将
// rd 转换为基类型返回。
func NewReaderSize(rd io.Reader, size int) *Reader

// NewReader 相当于 NewReaderSize(rd, 4096)
func NewReader(rd io.Reader) *Reader

// bufio.Reader 实现了如下接口：
// io.Reader
// io.WriterTo
// io.ByteScanner
// io.RuneScanner

// Peek 返回缓存的一个切片，该切片引用缓存中前 n 个字节的数据，
// 该操作不会将数据读出，只是引用，引用的数据在下一次读取操作之
// 前是有效的。如果切片长度小于 n，则返回一个错误信息说明原因。
// 如果 n 大于缓存的总大小，则返回 ErrBufferFull。
func (b *Reader) Peek(n int) ([]byte, error)

// Read 从 b 中读出数据到 p 中，返回读出的字节数和遇到的错误。
// 如果缓存不为空，则只能读出缓存中的数据，不会从底层 io.Reader
// 中提取数据，如果缓存为空，则：
// 1、len(p) >= 缓存大小，则跳过缓存，直接从底层 io.Reader 中读
// 出到 p 中。
// 2、len(p) < 缓存大小，则先将数据从底层 io.Reader 中读取到缓存
// 中，再从缓存读取到 p 中。
func (b *Reader) Read(p []byte) (n int, err error)

// Buffered 返回缓存中未读取的数据的长度。
func (b *Reader) Buffered() int

// Discard 跳过后续的 n 个字节的数据，返回跳过的字节数。
// 如果结果小于 n，将返回错误信息。
// 如果 n 小于缓存中的数据长度，则不会从底层提取数据。
func (b *Reader) Discard(n int) (discarded int, err error)

// ReadSlice 在 b 中查找 delim 并返回 delim 及其之前的所有数据。
// 该操作会读出数据，返回的切片是已读出的数据的引用，切片中的数据
// 在下一次读取操作之前是有效的。
//
// 如果找到 delim，则返回查找结果，err 返回 nil。
// 如果未找到 delim，则：
// 1、缓存不满，则将缓存填满后再次查找。
// 2、缓存是满的，则返回整个缓存，err 返回 ErrBufferFull。
//
// 如果未找到 delim 且遇到错误（通常是 io.EOF），则返回缓存中的所
// 有数据和遇到的错误。
//
// 因为返回的数据有可能被下一次的读写操作修改，所以大多数操作应该
// 使用 ReadBytes 或 ReadString，它们返回的是数据的拷贝。
func (b *Reader) ReadSlice(delim byte) (line []byte, err error)

// ReadLine 是一个低水平的行读取原语，大多数情况下，应该使用
// ReadBytes('\n') 或 ReadString('\n')，或者使用一个 Scanner。
//
// ReadLine 通过调用 ReadSlice 方法实现，返回的也是缓存的切片。用于
// 读取一行数据，不包括行尾标记（\n 或 \r\n）。
//
// 只要能读出数据，err 就为 nil。如果没有数据可读，则 isPrefix 返回
// false，err 返回 io.EOF。
//
// 如果找到行尾标记，则返回查找结果，isPrefix 返回 false。
// 如果未找到行尾标记，则：
// 1、缓存不满，则将缓存填满后再次查找。
// 2、缓存是满的，则返回整个缓存，isPrefix 返回 true。
//
// 整个数据尾部“有一个换行标记”和“没有换行标记”的读取结果是一样。
//
// 如果 ReadLine 读取到换行标记，则调用 UnreadByte 撤销的是换行标记，
// 而不是返回的数据。
func (b *Reader) ReadLine() (line []byte, isPrefix bool, err error)

// ReadBytes 功能同 ReadSlice，只不过返回的是缓存的拷贝。
func (b *Reader) ReadBytes(delim byte) (line []byte, err error)

// ReadString 功能同 ReadBytes，只不过返回的是字符串。
func (b *Reader) ReadString(delim byte) (line string, err error)

// Reset 将 b 的底层 Reader 重新指定为 r，同时丢弃缓存中的所有数据，复位
// 所有标记和错误信息。 bufio.Reader。
func (b *Reader) Reset(r io.Reader)
*/

package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {

	sr := strings.NewReader("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890") //36bytes

	buf := bufio.NewReaderSize(sr, 0) //小于16会被置于16
	b := make([]byte, 10)

	fmt.Println(buf.Buffered()) //0
	s, _ := buf.Peek(7)
	s[0], s[1], s[2] = 'a', 'b', 'c'
	fmt.Printf("%d   %q\n", buf.Buffered(), s) //16   "abcDE"

	buf.Discard(1)

	for n, err := 0, error(nil); err == nil; {
		n, err = buf.Read(b)
		fmt.Printf("%d   %q   %v\n", buf.Buffered(), b[:n], err)
	}
}

// 5   "bcDEFGHIJK"   <nil>
// 0   "LMNOP"   <nil>
// 6   "QRSTUVWXYZ"   <nil>
// 0   "123456"   <nil>
// 0   "7890"   <nil>
// 0   ""   EOF
