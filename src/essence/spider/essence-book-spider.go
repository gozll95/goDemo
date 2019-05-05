# book-spider

1."并发结果值最好用chan来传递":
ErrorChan()<-chan error


2."调度器或者其他类似manager的方法":
- Start(xxx)
- Stop(xxx)
- Status() Status
- Summary() xxx

3."数据缓冲池":
- 请求缓冲池:传输请求类型
- 响应缓冲池:传输响应类型
- 条目缓冲池:传输条目类型
- 错误缓冲池:传输错误类型

每个缓冲池需要2个参数:缓冲池中单一缓冲器的容量+缓冲池包含的缓冲器的最大数量。

// DataArgs 代表数据相关的参数容器的类型。
type DataArgs struct {
	// ReqBufferCap 代表请求缓冲器的容量。
	ReqBufferCap uint32 `json:"req_buffer_cap"`
	// ReqMaxBufferNumber 代表请求缓冲器的最大数量。
	ReqMaxBufferNumber uint32 `json:"req_max_buffer_number"`
	// RespBufferCap 代表响应缓冲器的容量。
	RespBufferCap uint32 `json:"resp_buffer_cap"`
	// RespMaxBufferNumber 代表响应缓冲器的最大数量。
	RespMaxBufferNumber uint32 `json:"resp_max_buffer_number"`
	// ItemBufferCap 代表条目缓冲器的容量。
	ItemBufferCap uint32 `json:"item_buffer_cap"`
	// ItemMaxBufferNumber 代表条目缓冲器的最大数量。
	ItemMaxBufferNumber uint32 `json:"item_max_buffer_number"`
	// ErrorBufferCap 代表错误缓冲器的容量。
	ErrorBufferCap uint32 `json:"error_buffer_cap"`
	// ErrorMaxBufferNumber 代表错误缓冲器的最大数量。
	ErrorMaxBufferNumber uint32 `json:"error_max_buffer_number"`
}

缓冲池和缓冲器是一对程序实体,缓冲器是缓冲池的底层支持。缓冲器是缓冲池的底层支持,缓冲池是缓冲器的再封装。缓冲池利用它持有的缓冲器实现***数据存取***的功能,并可以根据情况***自动地增减它持有的缓冲器的数量***。

// Pool 代表数据缓冲池的接口类型。
type Pool interface {
	// BufferCap 用于获取池中缓冲器的统一容量。
	BufferCap() uint32
	// MaxBufferNumber 用于获取池中缓冲器的最大数量。
	MaxBufferNumber() uint32
	// BufferNumber 用于获取池中缓冲器的数量。
	BufferNumber() uint32
	// Total 用于获取缓冲池中数据的总数。
	Total() uint64
	// Put 用于向缓冲池放入数据。
	// 注意！本方法应该是阻塞的。
	// 若缓冲池已关闭则会直接返回非nil的错误值。
	Put(datum interface{}) error
	// Get 用于从缓冲池获取数据。
	// 注意！本方法应该是阻塞的。
	// 若缓冲池已关闭则会直接返回非nil的错误值。
	Get() (datum interface{}, err error)
	// Close 用于关闭缓冲池。
	// 若缓冲池之前已关闭则返回false，否则返回true。
	Close() bool
	// Closed 用于判断缓冲池是否已关闭。
	Closed() bool
}

要点:缓冲池"动态伸缩"缓冲器

"缓冲器":
// Buffer 代表FIFO的缓冲器的接口类型。
type Buffer interface {
	// Cap 用于获取本缓冲器的容量。
	Cap() uint32
	// Len 用于获取本缓冲器中的数据数量。
	Len() uint32
	// Put 用于向缓冲器放入数据。
	// 注意！本方法应该是非阻塞的。
	// 若缓冲器已关闭则会直接返回非nil的错误值。
	Put(datum interface{}) (bool, error)
	// Get 用于从缓冲器获取器。
	// 注意！本方法应该是非阻塞的。
	// 若缓冲器已关闭则会直接返回非nil的错误值。
	Get() (interface{}, error)
	// Close 用于关闭缓冲器。
	// 若缓冲器之前已关闭则返回false，否则返回true。
	Close() bool
	// Closed 用于判断缓冲器是否已关闭。
	Closed() bool
}








4."数据"
检查数据的有效性

// Args 代表参数容器的接口类型。
type Args interface {
	// Check 用于自检参数的有效性。
	// 若结果值为nil，则说明未发现问题，否则就意味着自检未通过。
	Check() error
}

## 状态[要包含中间状态]


## 5.channel注意
- 向一个已关闭的通道发送值
- 关闭一个已关闭的通道
上述两个操作都会发恐慌


5."多重读取器"
可以多次提供基于同一底层数据(可以是[]byte类型的)io.ReadCloser类型值的类型

// MultipleReader 代表多重读取器的接口。
type MultipleReader interface {
	// Reader 用于获取一个可关闭读取器的实例。
	// 后者会持有本多重读取器中的数据。
	Reader() io.ReadCloser
}


// myMultipleReader 代表多重读取器的实现类型。
type myMultipleReader struct {
	data []byte
}

// NewMultipleReader 用于新建并返回一个多重读取器的实例。
func NewMultipleReader(reader io.Reader) (MultipleReader, error) {
	var data []byte
	var err error
	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("multiple reader: couldn't create a new one: %s", err)
		}
	} else {
		data = []byte{}
	}
	return &myMultipleReader{
		data: data,
	}, nil
}

func (rr *myMultipleReader) Reader() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(rr.data))
}

6.tips:
"ioutil.NopCloser":通常用这个函数包装无需关闭的读取器,这就是NopCloser的含义。

7."调度器":
	ctx context.Context
	// cancelFunc 代表取消函数，用于停止调度器。
	cancelFunc context.CancelFunc
	// status 代表状态。
	status Status
	// statusLock 代表专用于状态的读写锁。
	statusLock sync.RWMutex
	// summary 代表摘要信息。
	summary SchedSummary

综上所述,Init方法就做了4件事:
- 检查调度器状态
- 检查参数
- 初始化内部字段
- 注册组件实例

Init
	- check status and set 
	- defer check err and set status
	- check args valid
	- register modulesArgs

Start
5."启动并开启goroutine":
xx.Init()
go func() {
    errChan := sched.ErrorChan()
    for {
        err, ok := <-errChan
        if !ok {
            break
        }
        t.Errorf("An error occurs when running scheduler: %s", err)
    }
}()


6."缓冲器":
// myBuffer 代表缓冲器接口的实现类型。
type myBuffer struct {
	// ch 代表存放数据的通道。
	ch chan interface{}
	// closed 代表缓冲器的关闭状态：0-未关闭；1-已关闭。
	closed uint32
	// closingLock 代表为了消除因关闭缓冲器而产生的竞态条件的读写锁。
	closingLock sync.RWMutex
}

### Put方法的实现
func (buf *myBuffer) Put(datum interface{}) (ok bool, err error) {
	buf.closingLock.RLock()
	defer buf.closingLock.RUnlock()
	if buf.Closed() {
		return false, ErrClosedBuffer
	}
	select {
	case buf.ch <- datum:
		ok = true
	default:
		ok = false
	}
	return
}

select语句主要是为了让Put方法永远不会阻塞在发送操作上,在default分支中把结果变量ok的值设置为false,加之这时的结果变量err必为ni,就可以告知调用方放入数据的操作未成功,且原因并不是缓冲器已关闭,而是缓冲器已满。


###Close方法的实现
再说Close方法,在关闭通道之前,先要避免重复操作。因为重复关闭一个通道也会引发运行时恐慌。***避免措施就是先检查closed字段的值。当然,必须使用原子操作***。

func (buf *myBuffer) Close() bool {
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		buf.closingLock.Lock()
		close(buf.ch)
		buf.closingLock.Unlock()
		return true
	} 
	return false
}

"根据判断结果来执行后续操作->推荐使用CAS"

###Closed方法的实现
在Closed方法中***读取closed字段的值***,也一定要使用***原子操作***

func (buf *myBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}

#######重点:千万不要假设读取共享资源就是并发安全的,除非资源本身做出了这种保证。
"原子地读取值"

"双层通道":
### pay attation to:
注意:bufCh字段的类型是chan Buffer,一个元素类型为Buffer的通道类型。这与缓冲器同样是通道类型的ch字段联合起来看,就是一个***双层通道***的设计。***在放入或获取数据时,我会先从bufCh拿到一个缓冲器,再向该缓冲器放入数据或从该缓冲器获取数据,然后再把它发送回bufCh***。这样的设计有如下几点好处:
- bufCh中的每个缓冲器一次只会被一个goroutine中的程序(以下简称并发程序)拿到。并且,在放回bufCh之前,它对其他并发程序都是不可见的。一个缓冲器每次只会被并发程序放入或取走一个数据。即使同一个程序连续调用多次Put方法或Get方法,也会这样。缓冲器不至于一下被填满或取空。
- 更进一步看,bufCh是FIFO的。当把先前拿出的缓冲器归还给bufCh时,该缓冲器总会被放在队尾。也就是说,池中缓冲器的操作频率可以降到最低,这也有利于池中数据的均匀分布。
- 在从bufCh拿到缓冲器后,我可以判断是否需要缩减缓冲器的数量。如果需要并且该缓冲器已空,就可以直接把它关掉,并且不还给bufCh。另一方面,如果在放入数据时发现所有缓冲器都已满并且在一段时间内就没有空位,就可以新建一个缓冲器并放入bufCh。总之,这让缓冲池***自动伸缩功能***的实现变得简单了。
- 最后也是最重要的是,bufCh本身提供了对并发安全的保障。


"缓冲池"也很值得看:
// myPool 代表数据缓冲池接口的实现类型。
type myPool struct {
	// bufferCap 代表缓冲器的统一容量。
	bufferCap uint32
	// maxBufferNumber 代表缓冲器的最大数量。
	maxBufferNumber uint32
	// bufferNumber 代表缓冲器的实际数量。
	bufferNumber uint32
	// total 代表池中数据的总数。
	total uint64
	// bufCh 代表存放缓冲器的通道。
	bufCh chan Buffer
	// closed 代表缓冲池的关闭状态：0-未关闭；1-已关闭。
	closed uint32
	// lock 代表保护内部共享资源的读写锁。
	rwlock sync.RWMutex
}

#### Put方法
Put方法有两个主要的功能:
- 向缓冲池中放入数据
- 当发现所有的缓冲器都已满一段时间后,新建一个缓冲器并将其放入缓冲池。当然,如果当前缓冲池持有的缓冲器已达最大数量,就不能这么做了。所以,这里我们首先需要建立一个***发现和触发追加缓冲器操作的机制***。我规定当对池中所有缓冲器的操作的失败次数都达到5次时,就追加一个缓冲器入池。


func (pool *myPool) Put(datum interface{}) (err error) {
	if pool.Closed() {
		return ErrClosedBufferPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 5
	var ok bool
	for buf := range pool.bufCh {
		ok, err = pool.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

实际上,放入操作的核心逻辑在myPool类型的putData方法中。Put方法本身做的主要是不断的取出池中的缓冲器,并持有一个统一的***"已满"***计数。请注意count和maxCount变量的初始值。

#### PutData方法

func (pool *myPool) putData(
	buf Buffer, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
	...省略代码
}

##### 第一段
putData为了及时响应缓冲池的关闭,需要在一开始就检***查缓冲池的状态***。并且在方法执行结束前还要检查一次,以便***及时释放资源***。

if pool.Closed() {
	return false, ErrClosedBufferPool
}
defer func() {
	pool.rwlock.RLock()
	if pool.Closed() {
		atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
		err = ErrClosedBufferPool
	} else {
		pool.bufCh <- buf
	}
	pool.rwlock.RUnlock()
}()


##### 第二段 
执行向拿到的缓冲器放入数据的操作,并在必要时增加***已满***计数:

	ok, err = buf.Put(datum)
	if ok {
		atomic.AddUint64(&pool.total, 1)
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已满而未放入数据就递增计数。
	(*count)++


请注意那两条return语句以及最后的(*count)++。在试图向缓冲器放入数据后,我们需要立即判断操作结果。如果ok的值是true,就说明放入成功,此时就可以在递增total字段的值后立即返回。如果err的值不为nil,就是说缓冲器已关闭,这时就不需要再执行后面的语句了。除了这两种情况,我们就需要递增count的值。因为这时说明缓冲器已满。

这里的count值递增操作与第三段代码息息相关,这涉及对追加缓冲器的操作的触发。
	// 如果尝试向缓冲器放入数据的失败次数达到阈值，
	// 并且池中缓冲器的数量未达到最大值，
	// 那么就尝试创建一个新的缓冲器，先放入数据再把它放入池。
	if *count >= maxCount &&
		pool.BufferNumber() < pool.MaxBufferNumber() {
		pool.rwlock.Lock()
		if pool.BufferNumber() < pool.MaxBufferNumber() {
			if pool.Closed() {
				pool.rwlock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(pool.bufferCap)
			newBuf.Put(datum)
			pool.bufCh <- newBuf
			atomic.AddUint32(&pool.bufferNumber, 1)
			atomic.AddUint64(&pool.total, 1)
			ok = true
		}
		pool.rwlock.Unlock()
		*count = 0
	}
	return

在这段代码中,我用到了***双检锁***。如果第一次条件判断通过,就会立即再做一次条件判断。不过这之前,我会先锁定rwlock的写锁。这有两个作用:第一,防止向已关闭的缓冲池追加缓冲器。第二,防止缓冲器的数量超过最大值。在确保这两种情况不会发生后,我就会把一个已放入那个数据的缓冲器追加到缓冲池中。


#### Get方法
Get方法的总体流程与Put方法基本一致:

func (pool *myPool) Get() (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPool
	}
	var count uint32
	maxCount := pool.BufferNumber() * 10
	for buf := range pool.bufCh {
		datum, err = pool.getData(buf, &count, maxCount)
		if datum != nil || err != nil {
			break
		}
	}
	return
}

我把"已空"计数的上线maxCount设为缓冲器数量的10倍。也就是说,若在遍历所有缓冲器10次之后仍无法获取到数据。Get方法就会从缓冲池中去掉一个空的缓冲器。

#### getData方法
getData方法声明如下:

// getData 用于从给定的缓冲器获取数据，并在必要时把缓冲器归还给池。
func (pool *myPool) getData(
	buf Buffer, count *uint32, maxCount uint32) (datum interface{}, err error) {
	if pool.Closed() {
		return nil, ErrClosedBufferPool
	}
	defer func() {
		// 如果尝试从缓冲器获取数据的失败次数达到阈值，
		// 同时当前缓冲器已空且池中缓冲器的数量大于1，
		// 那么就直接关掉当前缓冲器，并不归还给池。
		if *count >= maxCount &&
			buf.Len() == 0 &&
			pool.BufferNumber() > 1 {
			buf.Close()
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			*count = 0
			return
		}
		pool.rwlock.RLock()
		if pool.Closed() {
			atomic.AddUint32(&pool.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPool
		} else {
			pool.bufCh <- buf
		}
		pool.rwlock.RUnlock()
	}()
	datum, err = buf.Get()
	if datum != nil {
		atomic.AddUint64(&pool.total, ^uint64(0))
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已空未取出数据就递增计数。
	(*count)++
	return
}

#### Close方法
func (pool *myPool) Close() bool {
	if !atomic.CompareAndSwapUint32(&pool.closed, 0, 1) {
		return false
	}
	pool.rwlock.Lock()
	defer pool.rwlock.Unlock()
	close(pool.bufCh)
	for buf := range pool.bufCh {
		buf.Close()
	}
	return true
}

#### Closed方法
func (pool *myPool) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}




7."扩展":
了解container/List/Ring类型

8."双检锁":
"判断-锁-判断"
这里的count值递增操作与第三段代码息息相关,这涉及对追加缓冲器的操作的触发。
	// 如果尝试向缓冲器放入数据的失败次数达到阈值，
	// 并且池中缓冲器的数量未达到最大值，
	// 那么就尝试创建一个新的缓冲器，先放入数据再把它放入池。
	if *count >= maxCount &&
		pool.BufferNumber() < pool.MaxBufferNumber() {
		pool.rwlock.Lock()
		if pool.BufferNumber() < pool.MaxBufferNumber() {
			if pool.Closed() {
				pool.rwlock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(pool.bufferCap)
			newBuf.Put(datum)
			pool.bufCh <- newBuf
			atomic.AddUint32(&pool.bufferNumber, 1)
			atomic.AddUint64(&pool.total, 1)
			ok = true
		}
		pool.rwlock.Unlock()
		*count = 0
	}
	return

在这段代码中,我用到了***双检锁***。如果第一次条件判断通过,就会立即再做一次条件判断。不过这之前,我会先锁定rwlock的写锁。这有两个作用:第一,防止向已关闭的缓冲池追加缓冲器。第二,防止缓冲器的数量超过最大值。在确保这两种情况不会发生后,我就会把一个已放入那个数据的缓冲器追加到缓冲池中。



