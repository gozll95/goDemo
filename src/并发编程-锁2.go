package main

import (
	"errors"
	"io"
	"os"
	"sync"
)

//数据文件的接口类型
type DataFile interface {
	//读取一个数据块
	Read() (rsn int64, d Data, err error)
	//写入一个数据块
	Write(d Data) (wsn int64, err error)
	//获取最后读取的数据块的序列号
	Rsn() int64
	//获取最后写入的数据块的序列号
	Wsn() int64
	//获取数据块的长度
	DataLen() unit32
}

//数据的类型
type Data []byte

//数据文件的实现类型
type myDataFile struct {
	f       *os.File     //文件
	fmutex  sync.RWMutex //被用于文件的读写锁
	woffset int64        //写操作需要用到的偏移量
	roffset int64        //读操作需要用到的偏移量
	wmutex  sync.Mutex   //写操作需要用到的互斥锁
	rmutex  sync.Mutex   //读操作需要用到的互斥锁
	dataLen unit32       //数据块长度
}

func NewDataFile(path string, dataLen unit32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, errors.New("Invalid data length!")
	}
	df := &myDataFile{f: f, dataLen: dataLen}
	return df, nil
}

//编写*myDataFile类型的Read方法
/*
1.获取并更新偏移量
2.根据读偏移量从文件中读取一块数据
3.把该数据块封装成一个Data类型值并将其作为结果值返回

其中,第一个步骤在被执行的时候应该由互斥锁rmutex保护起来,因为,我们要求多个读操作
不能读同一个数据块，并且它们应该按顺序的读取文件中的数据块。
而第二个步骤，我们也会用读写锁fmutex加以保护。
*/

/*
func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	//读取并更新读偏移量
	var offset int64
	df.rmutex.Lock()
	offset = df.roffset
	df.roffset += int64(df.dataLen)
	df.rmutex.Unlock()

	//读取一个数据块
	rsn = offset / int64(df.dataLen)
	df.fmutex.RLock()
	defer df.fmutex.RUnlock()
	bytes := make([]byte, df.dataLen)
	_, err = df.f.ReadAt(bytes, offset)
	if err != nil {
		return
	}
	d = bytes
	return
}
*/

//但是要考虑到边界情况
func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	//读取并更新读偏移量
	var offset int64
	df.rmutex.Lock()
	offset = df.roffset
	df.roffset += int64(df.dataLen)
	df.rmutex.Unlock()

	//读取一个数据块
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)

	//一直读到没有io.EOF为止
	for {
		df.fmutex.RLock()
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				df.fmutex.Unlock()
				continue
			}
			df.fmutex.RUnlock()
			return
		}
		d = bytes
		df.fmutex.RUnlock()
		return
	}

}

//Write方法
/*
1.获取更新写偏移量
2.向文件写入一个数据库
*/

func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	//读取并更新写偏移量
	var offset int64
	df.wmutex.Lock()
	offset = df.woffset
	df.woffset += int64(df.dataLen)
	df.wmutex.Unlock()

	//写入一个数据块
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	df.fmutex.Lock()
	defer df.fmutex.Unlock()
	_, err = df.f.Write(bytes)
	return
}

//Rsn
func (df *myDataFile) Rsn() int64 {
	df.rmutex.Lock()
	defer df.rmutex.Unlock()

	return df.roffset / int64(df.dataLen)
}

func (df *myDataFile) Wsn() int64 {
	df.wmutex.Lock()
	defer df.wmutex.Unlock()

	return df.woffset / int64(df.dataLen)
}
