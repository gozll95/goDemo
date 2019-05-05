package client

import (
	"crypto/sha1"
	"crypto/sha256"
	"io"
)

var (
	Helpers *_Helper
)

type _Helper struct{}

func (_ *_Helper) MakeSha1(data []byte) []byte {
	hash := sha1.New()
	hash.Write(data)

	return hash.Sum(nil)
}

func (_ *_Helper) MakeSha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)

	return hash.Sum(nil)
}

func (_ *_Helper) MakeSha256Reader(reader io.ReadSeeker) []byte {
	// reset reader offset after hash
	start, _ := reader.Seek(0, 1)
	defer reader.Seek(start, 0)

	hash := sha256.New()
	io.Copy(hash, reader)

	return hash.Sum(nil)
}
