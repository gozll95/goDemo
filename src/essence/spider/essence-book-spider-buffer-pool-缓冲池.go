# book-spider-buffer-pool

1."设计思路":

缓冲池下是缓冲器,可以自动伸缩缓冲器

双层channel

2."buffer的数据结构":

buffer用于存放数据,不阻塞

这里 注意点在于 closed close 


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

type myBuffer struct {
	// ch 代表存放数据的通道。
	ch chan interface{}
	// closed 代表缓冲器的关闭状态：0-未关闭；1-已关闭。
	closed uint32
	// closingLock 代表为了消除因关闭缓冲器而产生的竞态条件的读写锁。
	closingLock sync.RWMutex
}

- "buffer methods":

1) "Cap()"
return uint32(cap(buf.ch))

2) "Len()uint32"
return uint32(len(buf.ch))

3)"Put(datum interface{}) (ok bool, err error)"
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

4)"Get() (interface{}, error)" 
func (buf *myBuffer) Get() (interface{}, error) {
	select {
	case datum, ok := <-buf.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
		return nil, nil
	}
}

5)"func (buf *myBuffer) Close() bool"

func (buf *myBuffer) Close() bool {
	if atomic.CompareAndSwapUint32(&buf.closed, 0, 1) {
		buf.closingLock.Lock()
		close(buf.ch)
		buf.closingLock.Unlock()
		return true
	}
	return false
}

6)"Closed()bool"
func (buf *myBuffer) Closed() bool {
	if atomic.LoadUint32(&buf.closed) == 0 {
		return false
	}
	return true
}


3."注意事项":
当获取一个可变的全局变量的时候 使用aotmic.Load 
当根据结果操作的时候优先选择CAS 
学会atomic操作和lock并用


3."Pool":
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

func(pool *myPool)BufferCap()uint32{
	return pool.bufferCap 
}

func (pool *myPool) MaxBufferNumber() uint32 {
	return pool.maxBufferNumber
}

func(pool *myPool)BufferNumber()uint32{
	return atomic.LoadUint64(&pool.total)
}


// putData 用于向给定的缓冲器放入数据，并在必要时把缓冲器归还给池。
func (pool *myPool) putData(
	buf Buffer, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
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
}

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

func (pool *myPool) Closed() bool {
	if atomic.LoadUint32(&pool.closed) == 1 {
		return true
	}
	return false
}


6."测试"

	testingFunc := func(datum interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			ok, err := buf.Put(datum)
			if err != nil {
				t.Fatalf("An error occurs when putting a datum to the buffer: %s (datum: %d)",
					err, datum)
			}
			if !ok && buf.Len() < buf.Cap() {
				t.Fatalf("Couldn't put datum to the buffer! (datum: %d)",
					datum)
			}
		}
	}

		t.Run("Put in parallel(1)", func(t *testing.T) {
		for _, datum := range data[:size/2] {
			t.Run(fmt.Sprintf("Datum=%d", datum), testingFunc(datum, t))
		}
	})


	c:=func(a,b) func(a){

	}

	c(a,b)


t.Parallel()? ? ?


	t.Parallel() 是并行的 

	写了t.Parallel()的函数都会并行执行


	testingFunc := func(datum interface{}, t *testing.T) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			err := pool.Put(datum)
			if err != nil {
				t.Fatalf("An error occurs when putting a datum to the buffer pool: %s (datum: %d)",
					err, datum)
			}
			atomic.AddUint32(&count, 1)
			currentCount := atomic.LoadUint32(&count)
			if uint64(currentCount) > pool.Total() {
				t.Fatalf("Inconsistent data total: %d > %d (old > new)",
					currentCount, pool.Total())
			}
		}
	}
	t.Run("Put in parallel(1)", func(t *testing.T) {
		for _, datum := range data[:dataLen/2] {
			t.Run(fmt.Sprintf("Datum=%d", datum), testingFunc(datum, t))
		}
	})

	https://blog.csdn.net/wdy_yx/article/details/78210526