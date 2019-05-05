# book-datafile

1."前景":
假设需要创建一个文件存放数据,同一时刻可能会有多个goroutine分别对该文件进行写操作和读操作。每一次写操作都应该向该文件
写入若干字节的数据,这若干个字节的数据应该作为一个独立的数据块存在。这就意味着,写操作之间不能彼此干扰，数据库之间也不能
出现穿插和混淆的情况。另一个方面,每一次读操作都从这个文件读取一个独立的、完整的数据库。它们读取的数据块不能重复,且需要
按顺序读取。例如,第一个读操作读取了数据库1,那么第二个读操作读取数据库2,而第三个读操作读取数据库3,以此类推。对于这些
读操作是否可以并发进行,这里并不作要求。即使它们并发进行,程序也应该分辨出它们的先后顺序。

为了避免一些额外工作量,我规定每个数据库的长度都相同,该长度在读写操作进行前给定。若写操作实际欲写入数据的长度超过了该值,
则超过部分会被截掉。


2."接口":
这个接口需要有的方法:
Read()
Write()
GetReadIndex()获取当前读到哪里了
GetWriteIndex()获取当前写到哪里了

一般我们不知道一个具体的类型 统统定义[]byte 
比如这里的"数据块"

// Data 代表数据的类型。
type Data []byte

// DataFile 代表数据文件的接口类型。
type DataFile interface {
	// Read 会读取一个数据块。
	Read() (rsn int64, d Data, err error)
	// Write 会写入一个数据块。
	Write(d Data) (wsn int64, err error)
	// RSN 会获取最后读取的数据块的序列号。
	RSN() int64
	// WSN 会获取最后写入的数据块的序列号。
	WSN() int64
	// DataLen 会获取数据块的长度。
	DataLen() uint32
	// Close 会关闭数据文件。
	Close() error
}

分析:

读操作:
lock( it is read locker)
get read index
update read index
unlock(it is read locker)

rwlock( it is a rwlocker)
do read action
rw unlock( it is a rwlocker)


写操作:
lock( it is writer locker)
get write index 
update write index 
unlock( it is writer locker)

rwlock
do write action

所以:
// myDataFile 代表数据文件的实现类型。
type myDataFile struct {
	f       *os.File     // 文件。
	fmutex  sync.RWMutex // 被用于文件的读写锁。
	woffset int64        // 写操作需要用到的偏移量。
	roffset int64        // 读操作需要用到的偏移量。
	wmutex  sync.Mutex   // 写操作需要用到的互斥锁。
	rmutex  sync.Mutex   // 读操作需要用到的互斥锁。
	dataLen uint32       // 数据块长度。
}