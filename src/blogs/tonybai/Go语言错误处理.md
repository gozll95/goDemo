# 破解

## 1、checkError style

对于一些在error handle时可以选择goroutine exit（注意：如果仅存main goroutine一个goroutine，调用runtime.Goexit会导致program以crash形式退出）或os.Exit的情形，我们可以选择类似常见的checkError方式简化错误处理，例如：
```
func checkError(err error) {
    if err != nil {
        fmt.Println("Error is ", err)
        os.Exit(-1)
    }
}

func foo() {
    err := doStuff1()
    checkError(err)

    err = doStuff2()
    checkError(err)

    err = doStuff3()
    checkError(err)
}
```
这种方式有些类似于C中用宏(macro)简化错误处理过程代码，只是由于Go不支持宏，使得这种方式的应用范围有限。

## 2、聚合error handle functions

有些时候，我们会遇到这样的情况：
```
err := doStuff1()
if err != nil {
    //handle A
    //handle B
    ... ...
}

err = doStuff2()
if err != nil {
    //handle A
    //handle B
    ... ...
}

err = doStuff3()
if err != nil {
    //handle A
    //handle B
    ... ...
}
```

在每个错误处理过程，处理过程相似，都是handle A、handle B等，我们可以通过Go提供的defer + 闭包的方式，将handle A、handle B…聚合到一个defer匿名helper function中去：
```
func handleA() {
    fmt.Println("handle A")
}
func handleB() {
    fmt.Println("handle B")
}

func foo() {
    var err error
    defer func() {
        if err != nil {
            handleA()
            handleB()
        }
    }()

    err = doStuff1()
    if err != nil {
        return
    }

    err = doStuff2()
    if err != nil {
        return
    }

    err = doStuff3()
    if err != nil {
        return
    }
}
```

## 3.将doStuff和error处理绑定
```
    b := bufio.NewWriter(fd)
    b.Write(p0[a:b])
    b.Write(p1[c:d])
    b.Write(p2[e:f])
    // and so on
    if b.Flush() != nil {
            return b.Flush()
        }
    }
```

```
type Writer struct {
    err error
    buf []byte
    n   int
    wr  io.Writer
}


func (b *Writer) Write(p []byte) (nn int, err error) {
	for len(p) > b.Available() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, b.err = b.wr.Write(p)
		} else {
			n = copy(b.buf[b.n:], p)
			b.n += n
			b.Flush()
		}
		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.n:], p)
	b.n += n
	nn += n
	return nn, nil
}

```

我们可以看到,错误处理被绑定在Writer.Write的内部了,Writer定义中有一个err作为一个错误状态值,与Writer的实例绑定在了一起,并且在每次Write入口判断是否为!=nil。一旦!=nil,Write其实什么都没做就return了。

