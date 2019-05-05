GO语言提供的字典类型并不是并发安全的,因此需要使用一些同步方法对它进行扩展。这看起来好像并不困难,貌似只要用读写锁把读操作和写操作保护起来就可以了。
确实,读写锁是我们首先想到的同步工具。不过,使用锁进行并发访问控制太重了。

#1.先要确定并发安全的字典类型的行为。这显然需要一个接口类型。

```
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
```

#2.下面开始编写接口的实现类型,这里使用结构体类型。
```
// myConcurrentMap 代表ConcurrentMap接口的实现类型。
type myConcurrentMap struct {
	concurrency int
	segments    []Segment
	total       uint64
}
```

concurrency字段表示并发量,同时也代表了segments字段的长度。在这个并发安全字典的实现类型中，一个Segment类型值代表
一个散列段。***每个散列段都提供对其包含的键-元素对的读写操作***。在这里的读写操作需要由互斥锁保证其并发安全性.有多少个散列
段就有多少个互斥锁加以保护。这样的加锁方式称为"***分段锁***",是一种非常流行的并发控制实现。分段锁可以在适当降低互斥锁的开销
的同时保护共享资源。***在同一时刻,同一个散列段中的键-元素对只能有一个goroutine进行读写,但是不同散列段中的键-元素对是可以并发访问的,并且是安全的,若concurrency字段值为16,就可以有16个goroutine同时访问同一个并发安全字典,只要它们访问的散列段是不同的****这就是分段锁的意义和优势所在。

键-值元素对总数的增加只影响各个散列段的容量,而不影响它们的数量。散列段数量的固定可能会使得键-元素对分布不均，但是这从总体上来看并不是什么大问题。因为它可以通过良好的***段定位算法***和***设置足够多的并发量***来缓解，而且我还会在散列段中做键-元素对的***负载均衡***。

最后,total字段用于实时反映当前字典中键-元素对的实际数量,目的是让对字典容量的获取更加直接、简单和快速。uint64的类型让我可以对它实施原子操作。


#3.用于创建并初始化一个并发安全字典实例的函数是这样的
// NewConcurrentMap 会创建一个ConcurrentMap类型的实例。
// 参数pairRedistributor可以为nil。
func NewConcurrentMap(
	concurrency int,
	pairRedistributor PairRedistributor) (ConcurrentMap, error) {
	if concurrency <= 0 {
		return nil, newIllegalParameterError("concurrency is too small")
	}
	if concurrency > MAX_CONCURRENCY {
		return nil, newIllegalParameterError("concurrency is too large")
	}
	cmap := &myConcurrentMap{}
	cmap.concurrency = concurrency
	cmap.segments = make([]Segment, concurrency)
	for i := 0; i < concurrency; i++ {
		cmap.segments[i] =
			newSegment(DEFAULT_BUCKET_NUMBER, pairRedistributor)
	}
	return cmap, nil
}

DEFAULT_BUCKET_NUMBER代表一个散列段中默认包含的散列桶的数量。

#4.关于Put方法的声明
func (cmap *myConcurrentMap) Put(key string, element interface{}) (bool, error) {
	p, err := newPair(key, element)
	if err != nil {
		return false, err
	}
	s := cmap.findSegment(p.Hash())
	ok, err := s.Put(p)
	if ok {
		atomic.AddUint64(&cmap.total, 1)
	}
	return ok, err
}

在该方法中,先将两个参数值封装成了一个表示键-元素对的Pair类型值。
Pair类型实际上是一个接口

#5.关于Pair接口
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

#6.Pair的实现类型
// pair 代表键-元素对的类型。
type pair struct { -->这里实体类型是小写,接口是大写
	key string
	// hash 代表键的哈希值。
	hash    uint64
	element unsafe.Pointer
	next    unsafe.Pointer
}

注意,element和next字段都是unsafe.Pointer类型的。后者的实例可以代表一个可寻址的值的指针值。
***对于unsafe.Pointer类型的值是可以实施原子操作的。***
// newPair 会创建一个Pair类型的实例。
func newPair(key string, element interface{}) (Pair, error) {
	p := &pair{
		key:  key,
		hash: hash(key),
	}
	if element == nil {
		return nil, newIllegalParameterError("element is nil")
	}
	p.element = unsafe.Pointer(&element)
	return p, nil
}


请注意newPair函数中调用的函数hash,其功能是生成给定字符串的散列值。hash函数的优劣会影响到键-元素对是否能够均匀地分布到多个散列段以及散列桶中。分布越均匀,并发安全字典的读写操作耗时也就越稳定,也就意味着整体性能会更好。同时,散列值计算在读写操作耗时中占比也比较大。所以,在并发安全字典进行性能调优的时候,你应该优先考虑对hash函数的优化。

#7.findSegment //段定位算法
// findSegment 会根据给定参数寻找并返回对应散列段。
func (cmap *myConcurrentMap) findSegment(keyHash uint64) Segment {
	if cmap.concurrency == 1 {
		return cmap.segments[0]
	}
	var keyHash32 uint32
	if keyHash > math.MaxUint32 {
		keyHash32 = uint32(keyHash >> 32)
	} else {
		keyHash32 = uint32(keyHash)
	}
	return cmap.segments[int(keyHash32>>16)%(cmap.concurrency-1)]
}

可以看到,该算法的核心思想就是使用高位的几个字节来决定散列段的索引。这样可以让键-元素对在segments中分布得更广、更均匀一些。

#8.散列段
一旦定位到了散列段,就可以调用该散列段的Put方法放入当前的键-元素对实例。只要散列段的Put方法返回的第一个结果值是true,就需要用原子操作对myConcurrentMap的total字段+1,这表示添加了一个新的键-元素对实例。***可以说,并发安全字典的Put方法的核心主要是靠相应散列段的Put方法实现的***。Get方法、Delete方法其实也都是这样的。

func (cmap *myConcurrentMap) Get(key string) interface{} {
	keyHash := hash(key)
	s := cmap.findSegment(keyHash)
	pair := s.GetWithHash(key, keyHash)
	if pair == nil {
		return nil
	}
	return pair.Element()
}



func (cmap *myConcurrentMap) Delete(key string) bool {
	s := cmap.findSegment(hash(key))
	if s.Delete(key) {
		atomic.AddUint64(&cmap.total, ^uint64(0))
		return true
	}
	return false
}

***这实际上是把复杂度留给了散列段,也理应如此。因为一个散列段就相当于一个并发安全的字典,只不过我又在上面封装了一层,以求把互斥锁的开销分摊并降低。在散列段中,我使用互斥锁对键-元素对的读写操作进行全面保护。***

散列段的接口声明是这样的

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

它与ConcurrentMap接口的声明很相似。其中的GetWithHash方法,纯粹是为了在某些情况下避免重复计算键的散列值而声明。

Segment接口的实现类型是segment的基本结构如下:

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

***pairRedistributor***字段用于存储使用者通过NewConcurrentMap函数传入的键-值对的***再分布器***,用于把散列段中的所有键-值元素对均匀地分布到所有散列桶中。

// NewSegment 会创建一个Segment类型的实例。
func newSegment(
	bucketNumber int, pairRedistributor PairRedistributor) Segment {
	if bucketNumber <= 0 {
		bucketNumber = DEFAULT_BUCKET_NUMBER
	}
	if pairRedistributor == nil {
		pairRedistributor =
			newDefaultPairRedistributor(
				DEFAULT_BUCKET_LOAD_FACTOR, bucketNumber)
	}
	buckets := make([]Bucket, bucketNumber)
	for i := 0; i < bucketNumber; i++ {
		buckets[i] = newBucket()
	}
	return &segment{
		buckets:           buckets,
		bucketsLen:        bucketNumber,
		pairRedistributor: pairRedistributor,
	}
}


#9.再分布器

####很重要:再分布
其中DEFAULT_BUCKET_NUMBER常量前面已经解释过了,这里着重看第二条if语句。当使用者未传入有效的键-元素对再分布器时,就使用一个默认的实现。这个实现需要两个参数,一个是散列桶因子,一个是当前散列段中的散列桶数量。我会用当前散列段中的键-元素对的总数和散列桶数量计算出一个平均值,这个平均值表示在均衡的情况下每个散列桶应该包含多少个键-元素对。不会直接使用这个平均值,而是用它计算出一个单个散列桶可包含的键-元素对的数量的上限,即阈值。这个阈值对触发键-元素对的再分布操作非常有用。这里用到了散列桶装载因子。我用DEFAULT_BUCKET_LOAD_FACTOR常量表示其默认值。平均值乘以装载因子就可以得到阈值。这就是newDefaultPairRedistributor函数对其新建的再分布器的初始化。

我会在散列段的一些方法中用到pairRedistributor字段。所以这里先展示一下PairRedistributor接口的声明:

// PairRedistributor 代表针对键-元素对的再分布器。
// 用于当散列段内的键-元素对分布不均时进行重新分布。
type PairRedistributor interface {
	//  UpdateThreshold 会根据键-元素对总数和散列桶总数计算并更新阈值。
	UpdateThreshold(pairTotal uint64, bucketNumber int)
	// CheckBucketStatus 用于检查散列桶的状态。
	CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	// Redistribe 用于实施键-元素对的再分布。
	Redistribe(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}

下面开始说散列段的几个方法

#10.散列段的方法
//涉及到再分布
func (s *segment) Put(p Pair) (bool, error) {
	s.lock.Lock()
	b := s.buckets[int(p.Hash()%uint64(s.bucketsLen))]
	ok, err := b.Put(p, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, 1)
		s.redistribute(newTotal, b.Size())
	}
	s.lock.Unlock()
	return ok, err
}

func (s *segment) GetWithHash(key string, keyHash uint64) Pair {
	s.lock.Lock()
	b := s.buckets[int(keyHash%uint64(s.bucketsLen))]
	s.lock.Unlock()
	return b.Get(key)
}

func (s *segment) Delete(key string) bool {
	s.lock.Lock()
	b := s.buckets[int(hash(key)%uint64(s.bucketsLen))]
	ok := b.Delete(key, nil)
	if ok {
		newTotal := atomic.AddUint64(&s.pairTotal, ^uint64(0))
		s.redistribute(newTotal, b.Size())
	}
	s.lock.Unlock()
	return ok
}

func (s *segment) Size() uint64 {
	return atomic.LoadUint64(&s.pairTotal)
}

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


#11.散列桶
字典通常是根据键散列值存取键-元素对的,并且同一个键在一个字典中只能存有一份,后存入的键-元素对会替代先存入的相同键的键-元素对。你应该知道,不同字符串的散列值有可能是相同的,这取决于使用的散列函数。这种现象称之为***"散列值碰撞"***。那么碰撞发生之后应该怎么解决呢?

***你可以想象有一个桶,桶中装有且只装有键散列值相同的键-元素对。这些键-元素对之前由单向链相连。只要获取到桶中的第一个键-元素对，就可以顺藤摸瓜的查出桶中所有的键-元素对，因此通常只要记录下前者即可。这样的桶通常称为"散列桶"。***如此一来,在查找一个键-元素对的时候就需要进行***两次对比***。
- 第一次是对比键散列值，从而找到对应的散列桶;这在前文所述的散列段实现中已有所体现。
- 第二次是通过单链表遍历桶中的所有的键-元素对,逐一比较键本身。由于第一次对比已经极大的缩小了查找范围。因此有效减少了时间复杂度为O(n)的第二次对比的实际耗时。这也是散列桶的价值所在。

现在扩展一下,一个散列段包含若干散列桶。***但是,我会为每个散列桶指定一个散列值的集合,而不是单一的散列值。***只要一个键-元素对的键散列值在某个集合中，就会被放入对应的那个散列桶。这也是前面的代码s.buckets[int(p.Hash()%uint64(s.bucketsLen))]所遵循的规则。它可以大大减少散列桶的数量,而不失散列桶的本质。

实际上,很多使用不同编程语言实现的并发安全字典都基于上述几点。

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

其中的方法Put、Delete和Clear都接受一个sync.Locker类型的参数。这也就意味着对这些方法的调用需要由锁保护:使用者要么传入一个锁,要么自行加锁。在segment类的相关方法中,我使用了自行加锁的方法。

***为什么散列桶的Get方法和GetFirstPair方法不用加锁?***这是因为其中使用了一些小技巧,在无锁的情况下消除了散列桶的***读操作之间***,以及***读操作与写操作***之间的竞态条件。

？？？这里始终不懂

实现Bucket接口的是Bucket类型

// bucket 代表并发安全的散列桶的类型。
type bucket struct {
	// firstValue 存储的是键-元素对列表的表头。
	firstValue atomic.Value
	size       uint64
}

// 占位符。
// 由于原子值不能存储nil，所以当散列桶空时用此符占位。
var placeholder Pair = &pair{}

// newBucket 会创建一个Bucket类型的实例。
func newBucket() Bucket {
	b := &bucket{}
	b.firstValue.Store(placeholder)
	return b
}

#散列桶的方法
func (b *bucket) GetFirstPair() Pair {
	if v := b.firstValue.Load(); v == nil {
		return nil
	} else if p, ok := v.(Pair); !ok || p == placeholder {
		return nil
	} else {
		return p
	}
}

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


注意倒数第5行至倒数第2行的代码,它们处理的就是当前桶中未包含参数p的键的情况,是无锁化的关键。即使在Put方法执行期间Get方法被调用了,也不会产生竞态条件。

***对由firstValue字段表示的键-元素对单链表表头的变更总是原子的。键-元素对添加操作只会把参数p的值链向原有的表头,并把它变成新的表头,而新表头后面的原有单链表中的每个键-元素对,以及它们的链接关系都会原封不动。如此一来,键-元素对获取操作无论何时都可以原子获取到一个表头,并可以并发安全地向表尾遍历。***

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


# 总结
到此为止,我展示并说明了一个并发安全的字典的实现方法。
其中有4层封装,从下至上:
- 封装键-元素对的Pair接口:意义-检查键值和元素值的有效性,并先行计算键的散列值以备后用。同时,在需要对键-元素对的元素值(由element字段表示)或链接(由next字段表示)进行替换时进行原子操作。这比互斥锁要快得多。
- 封装Pair的单链表的Bucket接口:意义-存储键散列值在同一范围内的键-元素对,并使用单链表和原子值消除读操作之间以及读操作与写操作之间的竞态条件。这一方案是无锁化的,大大提高了操作的性能。不过,写操作之间的竞态条件只能用互斥锁来消除了。因为它们都可能建立新的单链表并替换旧链表,如果任由它们并发进行就有可能出现混乱。
- 封装Bucket切片的Segment接口的实现-意义:为了让单个散列桶集合能够自动伸缩,并承担字典内部的负载均衡。正因此,这一层上的读写操作都需要加锁。对于读操作来说,仅需要再依据键散列值定位散列桶这一步上加锁。当成功向散列段添加或者删除一个键-元素对的时候,就会触发对其中所有键-元素对的负载均衡。当然,真正执行负载均衡还需要满足一系列条件。负载均衡的频率不能过低也不能过高,否则就会影响性能。负载均衡的先决条件判定和执行是由字典使用者传入的,或由默认的键-元素对再分布器负责。
- 封装Segment切片的ConcurrentMap接口-意义:会根据使用者的需要初始化若干个散列段。散列段的数量会影响并发安全字典在当前应用场景下的整体性能,所以需要跟进实际情况去设定。本层完全把对并发安全的保证和键-元素对的负载均衡下放到了第3层,而只负责根据键-散列值找到对应的散列段并下达操作指令。这样就可以分摊同步方法的使用以及负载均衡的执行带来的开销,消除了重量级的全局锁,大幅提高了性能。






