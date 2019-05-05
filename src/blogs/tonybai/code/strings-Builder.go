package main 

import(
	"bytes"
	"strings"
)

type BuilderByByteBuffer struct{
	b bytes.Buffer
}

func(b *BuilderByByteBuffer)WriteString(s string)error{
	_,err:=b.b.WriteString(s)
	return err
}