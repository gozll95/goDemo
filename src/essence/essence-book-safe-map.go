# book-safe-map


1."锁进行并发访问控制太重了"
// ConcurrentMap 代表并发安全的字典的接口。
type ConcurrentMap interface {
	// Concurrency 会返回并发量。
	Concurrency() int
	// Put 会推送一个键-元素对。
	// 注意！参数element的值不能为nil。
	// 第一个返回值表示是否新增了键-元素对。
	// 若键已存在，新元素值会替换旧的元素值。
	Put(key string, element interface{}) (bool, error)
	// Get 会获取与指定键关联的那个元素。
	// 若返回nil，则说明指定的键不存在。
	Get(key string) interface{}
	// Delete 会删除指定的键-元素对。
	// 若结果值为true则说明键已存在且已删除，否则说明键不存在。
	Delete(key string) bool
	// Len 会返回当前字典中键-元素对的数量。
	Len() uint64
}


2."负载均衡-散列段-分段锁":
同时并发访问多段-每段加锁

cmap.segments[i] = newSegment(DEFAULT_BUCKET_NUMBER, pairRedistributor)

散列段-散列桶


3."单向链表":
// linkedPair 代表单向链接的键-元素对的接口。
type linkedPair interface {
	// Next 用于获得下一个键-元素对。
	// 若返回值为nil，则说明当前已在单链表的末尾。
	Next() Pair
	// SetNext 用于设置下一个键-元素对。
	// 这样就可以形成一个键-元素对的单链表。
	SetNext(nextPair Pair) error
}


// Pair 代表并发安全的键-元素对的接口。
type Pair interface {
	// linkedPair 代表单链键-元素对接口。
	linkedPair
	// Key 会返回键的值。
	Key() string
	// Hash 会返回键的哈希值。
	Hash() uint64
	// Element 会返回元素的值。
	Element() interface{}
	// Set 会设置元素的值。
	SetElement(element interface{}) error
	// Copy 会生成一个当前键-元素对的副本并返回。
	Copy() Pair
	// String 会返回当前键-元素对的字符串表示形式。
	String() string
}

Pair接口首先嵌入了linkedPair接口,后者是***包级私有***的。这主要是为了保护一些需要接口化的方法,使之不被包外代码访问。
实现linkedPair接口,可以让多个键-元素对形成一个***单链表***。

之所以有Hash方法,原因是:一个键-元素对值的键不可改变。因此,其键的散列值也是永远不变的。因此,在创建键-元素对的时候,先计算出这个散列值并存储起来以后备用。这样可以节省一些后续计算,提高效率。

// pair 代表键-元素对的类型。
type pair struct { -->这里实体类型是小写,接口是大写
	key string
	// hash 代表键的哈希值。
	hash    uint64
	element unsafe.Pointer
	next    unsafe.Pointer
}

对于unsafe.Pointer类型的值是可以实施原子操作的

4."xx":

Put(key string,value interface{})
	- key+value
		- pair(interface)[Hash()]
			-->find segment
				-->segment.Put(pair)
					-->atomic.AddUint64(xxx)


Get(key string)interface{}
	- key->hash->findSetment(keyHash)->segment GetWithHash(key,keyHash)-->pair.Element()

Delete(key string)
	- findSegment(hash(key)):segment--->segment Delete(key)-->atomic.AddUint64(..-1)


"这实际上是把复杂度留给了散列段,也理应如此。因为一个散列段就相当于一个并发安全的字典,只不过我又在上面封装了一层,以求把互斥锁的开销分摊并降低。在散列段中,我使用互斥锁对键-元素对的读写操作进行全面保护"

5."散列段":

// Segment 代表并发安全的散列段的接口。
type Segment interface {
	// Put 会根据参数放入一个键-元素对。
	// 第一个返回值表示是否新增了键-元素对。
	Put(p Pair) (bool, error)
	// Get 会根据给定参数返回对应的键-元素对。
	// 该方法会根据给定的键计算哈希值。
	Get(key string) Pair
	// GetWithHash 会根据给定参数返回对应的键-元素对。
	// 注意！参数keyHash应该是基于参数key计算得出哈希值。
	GetWithHash(key string, keyHash uint64) Pair
	// Delete 会删除指定键的键-元素对。
	// 若返回值为true则说明已删除，否则说明未找到该键。
	Delete(key string) bool
	// Size 用于获取当前段的尺寸（其中包含的散列桶的数量）。
	Size() uint64
}

// segment 代表并发安全的散列段的类型。
type segment struct {
	// buckets 代表散列桶切片。
	buckets []Bucket
	// bucketsLen 代表散列桶切片的长度。
	bucketsLen int
	// pairTotal 代表键-元素对总数。
	pairTotal uint64
	// pairRedistributor 代表键-元素对的再分布器。
	pairRedistributor PairRedistributor
	lock              sync.Mutex
}


type PairRedistributor interface {
	//  UpdateThreshold 会根据键-元素对总数和散列桶总数计算并更新阈值。
	UpdateThreshold(pairTotal uint64, bucketNumber int)
	// CheckBucketStatus 用于检查散列桶的状态。
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	// Redistribe 用于实施键-元素对的再分布。
	Redistribe(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}

散列段的方法:
Put(p Pair):
	- segment Lock()
		- b:=s.buckets[xxx] // 通过hash得到对应的散列桶
		- bucket Put(p,nil) //散列桶 Put
		- atomic.AddUint64 //更新 pairTotal 
		- segment redistribute(pairTotal,bucket size) //根据total pair和 散列桶size 再分布
	- segment Unlock()

GetWithHash(key string, keyHash uint64):
	- segment Lock()
	- b:=s.buckeets[xx] //通过hash得到对应的散列桶
	- segment Unlock()
	- bucket Get(key)

Delete ... 同上,类似

"这种get的操作优先使用atomic.Load"
Size()uint64:
	- atomic.LoadUint64



// redistribute 会检查给定参数并设置相应的阈值和计数，
// 并在必要时重新分配所有散列桶中的所有键-元素对。
// 注意！必须在互斥锁的保护下调用本方法！
func (s *segment) redistribute(pairTotal uint64, bucketSize uint64) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = newPairRedistributorError(pErr.Error())
			} else {
				err = newPairRedistributorError(fmt.Sprintf("%s", p))
			}
		}
	}()
	s.pairRedistributor.UpdateThreshold(pairTotal, s.bucketsLen)
	bucketStatus := s.pairRedistributor.CheckBucketStatus(pairTotal, bucketSize)
	newBuckets, changed := s.pairRedistributor.Redistribe(bucketStatus, s.buckets)
	if changed {
		s.buckets = newBuckets
		s.bucketsLen = len(s.buckets)
	}
	return nil
}


6."散列桶":-"存储散列值相同的something-扩展-为每个散列桶指定一个散列值的集合"
字典通常是根据键散列值存取键-元素对的,并且同一个键在一个字典中只能存有一份,后存入的键-元素对会替代先存入的相同键的键-元素对。你应该知道,不同字符串的散列值有可能是相同的,这取决于使用的散列函数。这种现象称之为***"散列值碰撞"***。那么碰撞发生之后应该怎么解决呢?

***你可以想象有一个桶,桶中装有且只装有键散列值相同的键-元素对。这些键-元素对之前由单向链相连。只要获取到桶中的第一个键-元素对，就可以顺藤摸瓜的查出桶中所有的键-元素对，因此通常只要记录下前者即可。这样的桶通常称为"散列桶"。***如此一来,在查找一个键-元素对的时候就需要进行***两次对比***。
- 第一次是对比键散列值，从而找到对应的散列桶;这在前文所述的散列段实现中已有所体现。
- 第二次是通过单链表遍历桶中的所有的键-元素对,逐一比较键本身。由于第一次对比已经极大的缩小了查找范围。因此有效减少了时间复杂度为O(n)的第二次对比的实际耗时。这也是散列桶的价值所在。

现在扩展一下,一个散列段包含若干散列桶。***但是,我会为每个散列桶指定一个散列值的集合,而不是单一的散列值。***只要一个键-元素对的键散列值在某个集合中，就会被放入对应的那个散列桶。这也是前面的代码s.buckets[int(p.Hash()%uint64(s.bucketsLen))]所遵循的规则。它可以大大减少散列桶的数量,而不失散列桶的本质。


// Bucket 代表并发安全的散列桶的接口。
type Bucket interface {
	// Put 会放入一个键-元素对。
	// 第一个返回值表示是否新增了键-元素对。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Put(p Pair, lock sync.Locker) (bool, error)
	// Get 会获取指定键的键-元素对。
	Get(key string) Pair
	// GetFirstPair 会返回第一个键-元素对。
	GetFirstPair() Pair
	// Delete 会删除指定的键-元素对。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Delete(key string, lock sync.Locker) bool
	// Clear 会清空当前散列桶。
	// 若在调用此方法前已经锁定lock，则不要把lock传入！否则必须传入对应的lock！
	Clear(lock sync.Locker)
	// Size 会返回当前散列桶的尺寸。
	Size() uint64
	// String 会返回当前散列桶的字符串表示形式。
	String() string
}


其实这里就算调用Bucket的Put方法之前,我们lock了segment的Lock,但是本着设计原则:
- 要假想任何场景,所以要设计并发安全,不能仅仅供给某一场景使用


"很厉害":
### Put方法
Put方法会先调用GetFirstPair方法,以获取桶中的第一个键-元素对,如果后者的返回结果为nil,那么前者就直接把参数值作为第一个键-元素对存入firstValue,否则就利用键-元素对实例的Next方法遍历所有的键-元素对,并判断该键是否存在。若已经存在,则直接替换与该键对应的元素值。***(这里替换是原子的)***。若不存在,则调用参数值的SetNext方法,把当前的第一个键-元素对指定为参数值的单链目标,然后把参数值存入firstValue。

func (b *bucket) Put(p Pair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newIllegalParameterError("pair is nil")
	}
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}
	var target Pair
	key := p.Key()
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			break
		}
	}
	if target != nil {
		target.SetElement(p.Element())
		return false, nil
	}
	p.SetNext(firstPair)
	b.firstValue.Store(p)
	atomic.AddUint64(&b.size, 1)
	return true, nil
}



func (b *bucket) Delete(key string, lock sync.Locker) bool {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return false
	}
	var prevPairs []Pair
	var target Pair
	var breakpoint Pair
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			breakpoint = v.Next()
			break
		}
		prevPairs = append(prevPairs, v)
	}
	if target == nil {
		return false
	}
	newFirstPair := breakpoint
	for i := len(prevPairs) - 1; i >= 0; i-- {
		pairCopy := prevPairs[i].Copy()
		pairCopy.SetNext(newFirstPair)
		newFirstPair = pairCopy
	}
	if newFirstPair != nil {
		b.firstValue.Store(newFirstPair)
	} else {
		b.firstValue.Store(placeholder)
	}
	atomic.AddUint64(&b.size, ^uint64(0))
	return true
}


func (b *bucket) Get(key string) Pair {
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return nil
	}
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			return v
		}
	}
	return nil
}


func (b *bucket) Clear(lock sync.Locker) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	atomic.StoreUint64(&b.size, 0)
	b.firstValue.Store(placeholder)
}



# 2018.3.19

if newNumber > currentNumber {
	for i := uint64(0); i < currentNumber; i++ {
		buckets[i].Clear(nil)
	}
	for j := newNumber - currentNumber; j > 0; j-- {
		buckets = append(buckets, newBucket())
	}
} else {
	buckets = make([]Bucket, newNumber)
	for i := uint64(0); i < newNumber; i++ {
		buckets[i] = newBucket()
	}
}

1.error

// IllegalParameterError 代表非法的参数的错误类型。
type IllegalParameterError struct {
	msg string
}

// newIllegalParameterError 会创建一个IllegalParameterError类型的实例。
func newIllegalParameterError(errMsg string) IllegalParameterError {
	return IllegalParameterError{
		msg: fmt.Sprintf("concurrent map: illegal parameter: %s", errMsg),
	}
}

func (ipe IllegalParameterError) Error() string {
	return ipe.msg
}