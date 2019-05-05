/*
// Scanner 提供了一个方便的接口来读取数据，例如遍历多行文本中的行。Scan 方法会通过
// 一个“匹配函数”读取数据中符合要求的部分，跳过不符合要求的部分。“匹配函数”由调
// 用者指定。本包中提供的匹配函数有“行匹配函数”、“字节匹配函数”、“字符匹配函数”
// 和“单词匹配函数”，用户也可以自定义“匹配函数”。默认的“匹配函数”为“行匹配函
// 数”，用于获取数据中的一行内容（不包括行尾标记）
//
// Scanner 使用了缓存，所以匹配部分的长度不能超出缓存的容量。默认缓存容量为 4096 -
// bufio.MaxScanTokenSize，用户可以通过 Buffer 方法指定自定义缓存及其最大容量。
//
// Scan 在遇到下面的情况时会终止扫描并返回 false（扫描一旦终止，将无法再继续）：
// 1、遇到 io.EOF
// 2、遇到读写错误
// 3、“匹配部分”的长度超过了缓存的长度
//
// 如果需要对错误进行更多的控制，或“匹配部分”超出缓存容量，或需要连续扫描，则应该
// 使用 bufio.Reader
type Scanner struct { ... }

// NewScanner 创建一个 Scanner 来扫描 r，默认匹配函数为 ScanLines。
func NewScanner(r io.Reader) *Scanner

// Buffer 用于设置自定义缓存及其可扩展范围，如果 max 小于 len(buf)，则 buf 的尺寸将
// 固定不可调。Buffer 必须在第一次 Scan 之前设置，否则会引发 panic。
// 默认情况下，Scanner 会使用一个 4096 - bufio.MaxScanTokenSize 大小的内部缓存。
func (s *Scanner) Buffer(buf []byte, max int)

// Split 用于设置“匹配函数”，这个函数必须在调用 Scan 前执行。
func (s *Scanner) Split(split SplitFunc)

// SplitFunc 用来定义“匹配函数”，data 是缓存中的数据。atEOF 标记数据是否读完。
// advance 返回 data 中已处理的数据的长度。token 返回找到的“匹配部分”，“匹配
// 部分”可以是缓存的切片，也可以是自己新建的数据（比如 bufio.errorRune）。“匹
// 配部分”将在 Scan 之后通过 Bytes 和 Text 反馈给用户。err 返回错误信息。
//
// 如果在 data 中无法找到一个完整的“匹配部分”则应返回 (0, nil, nil)，以便告诉
// Scanner 向缓存中填充更多数据，然后再次扫描（Scan 会自动重新扫描）。如果缓存已
// 经达到最大容量还没有找到，则 Scan 会终止并返回 false。
// 如果 data 为空，则“匹配函数”将不会被调用，意思是在“匹配函数”中不必考虑
// data 为空的情况。
//
// 如果 err != nil，扫描将终止，如果 err == ErrFinalToken，则 Scan 将返回 true，
// 表示扫描正常结束，如果 err 是其它错误，则 Scan 将返回 false，表示扫描出错。错误
// 信息可以在 Scan 之后通过 Err 方法获取。
//
// SplitFunc 的作用很简单，从 data 中找出你感兴趣的数据，然后返回，同时返回已经处理
// 的数据的长度。
type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)

// Scan 开始一次扫描过程，如果匹配成功，可以通过 Bytes() 或 Text() 方法取出结果，
// 如果遇到错误，则终止扫描，并返回 false。
func (s *Scanner) Scan() bool

// Bytes 将最后一次扫描出的“匹配部分”作为一个切片引用返回，下一次的 Scan 操作会覆
// 盖本次引用的内容。
func (s *Scanner) Bytes() []byte

// Text 将最后一次扫描出的“匹配部分”作为字符串返回（返回副本）。
func (s *Scanner) Text() string

// Err 返回扫描过程中遇到的非 EOF 错误，供用户调用，以便获取错误信息。
func (s *Scanner) Err() error

// ScanBytes 是一个“匹配函数”用来找出 data 中的单个字节并返回。
func ScanBytes(data []byte, atEOF bool) (advance int, token []byte, err error)

// ScanRunes 是一个“匹配函数”，用来找出 data 中单个 UTF8 字符的编码。如果 UTF8 编
// 码错误，则 token 会返回 "\xef\xbf\xbd"（即：U+FFFD），但只消耗 data 中的一个字节。
// 这使得调用者无法区分“真正的U+FFFD字符”和“解码错误的返回值”。
func ScanRunes(data []byte, atEOF bool) (advance int, token []byte, err error)

// ScanLines 是一个“匹配函数”，用来找出 data 中的单行数据并返回（包括空行）。
// 行尾标记可以是 \n 或 \r\n（返回值不包含行尾标记）
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error)

// ScanWords 是一个“匹配函数”，用来找出 data 中以空白字符分隔的单词。
// 空白字符由 unicode.IsSpace 定义。
func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error)
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 逗号分隔的字符串，最后一项为空
	const input = "1,2,3,4,"
	scanner := bufio.NewScanner(strings.NewReader(input))
	// 定义匹配函数（查找逗号分隔的字符串）
	onComma := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == ',' {
				return i + 1, data[:i], nil
			}
		}
		if atEOF {
			// 告诉 Scanner 扫描结束。
			return 0, data, bufio.ErrFinalToken
		} else {
			// 告诉 Scanner 没找到匹配项，让 Scan 填充缓存后再次扫描。
			return 0, nil, nil
		}
	}
	// 指定匹配函数
	scanner.Split(onComma)
	// 开始扫描
	for scanner.Scan() {
		fmt.Printf("%q ", scanner.Text())
	}
	// 检查是否因为遇到错误而结束
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
}

/*
// 示例：带检查扫描
func main() {
	const input = "1234 5678 1234567901234567890 90"
	scanner := bufio.NewScanner(strings.NewReader(input))
	// 自定义匹配函数
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// 获取一个单词
		advance, token, err = bufio.ScanWords(data, atEOF)
		// 判断其能否转换为整数，如果不能则返回错误
		if err == nil && token != nil {
			_, err = strconv.ParseInt(string(token), 10, 32)
		}
		// 这里包含了 return 0, nil, nil 的情况
		return
	}
	// 设置匹配函数
	scanner.Split(split)
	// 开始扫描
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}
}
*/

//默认的“匹配函数”为“行匹配函数”，用于获取数据中的一行内容（不包括行尾标记）
