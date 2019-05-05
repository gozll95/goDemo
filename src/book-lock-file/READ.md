假设需要创建一个文件存放数据,同一时刻可能会有多个goroutine分别对该文件进行写操作和读操作。每一次写操作都应该向该文件
写入若干字节的数据,这若干个字节的数据应该作为一个独立的数据块存在。这就意味着,写操作之间不能彼此干扰，数据库之间也不能
出现穿插和混淆的情况。另一个方面,每一次读操作都从这个文件读取一个独立的、完整的数据库。它们读取的数据块不能重复,且需要
按顺序读取。例如,第一个读操作读取了数据库1,那么第二个读操作读取数据库2,而第三个读操作读取数据库3,以此类推。对于这些
读操作是否可以并发进行,这里并不作要求。即使它们并发进行,程序也应该分辨出它们的先后顺序。

为了避免一些额外工作量,我规定每个数据库的长度都相同,该长度在读写操作进行前给定。若写操作实际欲写入数据的长度超过了该值,
则超过部分会被截掉。

为了实现上述需求，我创建了一个接口类型:

// 用于表示数据文件的接口类型
type DataFile interface{
    // 读取一个数据块
    Read()(rsn int64,d Data,err error)
    // 写入一个数据块
    Write(d Data)(wsn int64,err error)
    // 获取最后读取的数据块的序列号
    RSN() int64
    // 获取最后写入的数据块的序列号
    WSN() int64
    // 获取数据块的长度
    DataLen()uint32
    // 关闭数据文件
    Close() error
}

其中,类型Data被声明为一个[]byte的别名类型

//用于表示数据的类型
type Data []byte

下面来编写DataFile接口的实现类型,将其命名为myDataFile,它的基本结构如下:
//用于表示数据文件的实现类型
type myDataFile struct{
    f *os.File //文件
    fmutex sync.RWMutex //用于文件的读写锁
    woffset int64 //写操作需要用到的偏移量
    roffset int64 //读操作需要用到的偏移量
    wmutex sync.Mutex //写操作需要用到的互斥量
    rmutex sync.Mutex //读操作需要用到的互斥量
    dataLen uint32 //数据块长度
}

// 新建一个数据文件的实例
func NewDataFile(path string,dataLen uint32)(DataFile,error){
    f,err:=os.Create(path)
    if err!=nil{
        return nil,err
    }
    if dataLen==0{
        return nil,errors.New("Invalid data length!")
    }
    df:=&myDataFile{f:f,dataLen:dataLen}
    return df,nil
}

## 先来看*myDataFile类型的Read方法
该方法应该按照如下步骤实现:
- 获取并更新读偏移量
- 根据读偏移量从文件中读取一块数据
- 把该数据块封装成一个Data类型值并将其作为结果值返回

其中,步骤1在执行的时候应该由互斥锁rmutex保护起来,因为我要求多个读操作不能读取同一个数据块,并且它们应该按照顺序读取文件中的数据库。
而步骤2页会用读写锁fmutex加以保护。

下面是这个Read方法的第一个版本:
func(df *myDataFile)Read()(rsn int64,d Data,err error){
    //读取并更新偏移量
    var offset int64
    df.rmutex.Lock()
    offset=df.roffset
    df.roffset+=int64(df.dataLen)
    df.rmutex.Unlock()

    //读取一个数据块
    rsn=offset/int64(df.dataLen)
    df.fmutex.RLock()
    defer df.fmutex.RUnlock()
    bytes:=make([]byte,df.dataLen)
    _,err=df.f.ReadAt(bytes,offset)
    if err!=nil{
        return
    }
    d=bytes
    return
}

在读取一个数据库的时候，我适时地进行了fmutex字段的读锁定和读解锁,这可以保证在这里读取到的是完整的数据库。不过,这个完整的数据块却并不一定是正确的。为什么这么说呢?
***important***请想象这样的场景,在这个程序中,有3个goroutine并发的执行某个*myDataFile类型值的Read方法,并有2个goroutine并发的执行该值的Write方法。通过前3个goroutine的运行,数据文件中的数据块被依次读取出来。但是,由于进行写操作的goroutine比进行读操作的goroutine少。因此过不了多久,读偏移量roffset的值就会等于甚至大于偏移量woffset的值。也就是说,读操作很快就会没有数据可读了。这种情况会使上面的df.f.ReadAt方法返回的第二个结果值为io.EOF。io.EOF是一个变量,代表无更多数据可读的状态.(EOF实为End of File的缩写)。该变量虽然是error类型的,但我们不应该把它视为错误的代表,而应该看成是一种边界情况。

在这个版本的Read方法中并没有对这种***边界情况***作出正确的处理,该方法在遇到这种情况时会直接把错误返回给调用方。调用方会得到读取
出错的数据块的序列号,但却无法再次尝试读取这个数据块。由于其他正在或后续进行的Read方法会继续增加读偏移量roffset的值。因此当该调用方再次调用这个Read方法的时候,只能读取到在此数据块后面的其他数据块。注意,执行Read方法时遇到这种情况的次数越多,被漏读的数据块也就会越多。为了解决这个问题,我编写了Read方法的第二个版本:

func(df *myDataFile)Read()(rsn int64,d Data,err error){
    // 读取并更新偏移量
    // 省略若干diamante

    // 读取一个数据块
    rns=offset/int64(df.dataLen)
    bytes:=make([]byte,df.dataLen)
    for{
        df.fmutex.RLock()
        _,err:=df.f.ReadAt(bytes,offset)
        if err!=nil{
            if err==io.EOF{
                df.fmutex.RUnlock()
                continue
            }
            df.fmutex.RUlock()
            return
        }
        d=bytes
        df.fmutex.RUnlock()
        return
    }
}

***第二个版本的Read方法***使用for语句是为了达到这样一个目的:在其中的df.f.ReadAt方法返回io.EOF的时候,继续尝试获取同一个数据块,直到获取成功为止。注意,如果在该for代码块执行期间一直让读写锁fmutex处于读锁定状态,那么针对它的写锁定操作将永远不会成功,且相应
的goroutine也会一直阻塞。所以,我不得不在该循环中的每条return语句和continue语句都加入一个针对fmutex的读解锁操作。并在每次迭代开始时都会对fmutex进行一次读锁定。显然,这样的代码看起来有些丑陋。冗余的代码会使代码的维护成本和出错概率大大增加。并且,当for代码中
的代码引发运行时恐慌时,并不会及时对读写锁fmutex进行读解锁。因为要处理一种边界情况,而去掉了第一版中的defer df.fmutex.RUnlock()语句。这种做法利弊参半。

## 下面来考虑*myDataFile的Write方法
与Read方法相比,Write方法的实现会简单一些,因为后者不会涉及到边界情况。
步骤:
- 获取并更新写偏移量
- 向文件写入一个数据块

func(df *myDataFile)Write(d Data)(wsn int64,err error){
    // 读取并更新偏移量
    var offset int64
    df.wmutex.Lock()
    offset=df.woffset
    df.woffset+=int64(df.dataLen)
    df.wmutex.Unlock()

    //写入一个数据块
    wsn=offset/int64(df.dataLen)
    var bytes []byte
    if len(d)>int(df.dataLen){
        bytes=d[0:df.dataLen]
    }else{
        bytes=d
    }
    df.fmutex.Lock()
    defer df.fmutex.Unlock()
    _,err=df.f.Write(bytes)
    return
}

有了编写前面两个方法的经验,很容易编写出*myDataFile类型的RSN方法和WSN方法
func(df *myDataFile)RSN()int64{
    df.rmutex.Lock()
    defer df.rmutex.Unlock()
    return df.roffset/int64(df.dataLen) 
}

func(df *myDataFile)WNS()int64{
    df.wmutex.Lock()
    defer df.wmutex.Unlock()
    return df.woffset/int64(df.dataLen)
}

编写上面这个完整实例的主要目的是:展示互斥锁和读写锁在实际场景中的应用。由于还没有讲到Go语言提供的其他同步工具,因此相关方法中所有需要
同步的地方都是用锁来实现的。实际上,其中的一些问题用锁来解决是不足够或不合适的;我会在后面逐步对它们进行改进。

