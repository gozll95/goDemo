//sync.Pool适用于无状态的对象的复用,而不适用于连接池之类的。

//在fmt包中有一个很好的使用池的例子,它维护一个动态大小的临时输出缓冲区。

package main

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func timeNow() time.Time {
	return time.Unix(1136214245, 0)
}

func Log(w io.Writer, key, val string) {
	// 获取临时对象，没有的话会自动创建
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(timeNow().UTC().Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(val)
	w.Write(b.Bytes())
	// 将临时对象放回到 Pool 中
	bufPool.Put(b)
}

func main() {
	Log(os.Stdout, "path", "/search?q=flowers")
}

//2006-01-02T15:04:05Z path=/search?q=flowers
