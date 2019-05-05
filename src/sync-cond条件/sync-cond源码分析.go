sync.Cond的介绍和源码观察
Cond用于在并发环境下routine的等待和通知

结构体定义
type Cond struct {
    noCopy noCopy //不允许复制,一个结构体,有一个Lock()方法,嵌入别的结构体中,表示不允许复制
    L Locker    //锁
    notify  notifyList  //通知列表,调用Wait()方法的routine会被放入list中,每次唤醒,从这里取出
    checker copyChecker //复制检查,检查cond实例是否被复制
}

‘构造’方法
// NewCond returns a new Cond with Locker l.
//通过一个Locker实例初始化,传参数的时候必须是引用或指针,比如&sync.Mutex{}或new(sync.Mutex)
//不然会报异常:cannot use lock (type sync.Mutex) as type sync.Locker in argument to sync.NewCond:
//sync.Mutex does not implement sync.Locker (Lock method has pointer receiver)
func NewCond(l Locker) *Cond {
    return &Cond{L: l}
}

常用方法 
Wait
//调用此方法会将此routine加入通知列表,并等待获取通知,调用此方法必须先Lock,不然方法里会调用Unlock(),报错.
func (c *Cond) Wait() {
    c.checker.check()   //检查是否被复制
    t := runtime_notifyListAdd(&c.notify) //加入通知列表
    c.L.Unlock() // 释放锁
    runtime_notifyListWait(&c.notify, t) //等待通知
    c.L.Lock() //被通知了,获取锁,继续运行
}

- Signal

//唤醒在Wait的routine中的一个
    func (c *Cond) Signal() {
    c.checker.check() //检查是否被复制
    runtime_notifyListNotifyOne(&c.notify) //通知等待列表中的一个
}

- Broadcast

//唤醒所有等待的
func (c *Cond) Broadcast() {
    c.checker.check()
    runtime_notifyListNotifyAll(&c.notify)
}