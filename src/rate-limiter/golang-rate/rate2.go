//https://www.jianshu.com/p/4ce68a31a71d

接口介绍:

type Limiter
type Limiter struct {
    // contains filtered or unexported fields
}
Limter限制时间的发生频率，采用令牌池的算法实现。这个池子一开始容量为b，装满b个令牌，然后每秒往里面填充r个令牌。
由于令牌池中最多有b个令牌，所以一次最多只能允许b个事件发生，一个事件花费掉一个令牌。

Limter提供三中主要的函数 Allow, Reserve, and Wait. 大部分时候使用Wait。

func NewLimiter
func NewLimiter(r Limit, b int) *Limiter
NewLimiter 返回一个新的Limiter。

func (*Limiter) [Allow]
func (lim *Limiter) Allow() bool
Allow 是函数 AllowN(time.Now(), 1)的简化函数。

func (*Limiter) AllowN
func (lim *Limiter) AllowN(now time.Time, n int) bool
AllowN标识在时间now的时候，n个事件是否可以同时发生(也意思就是now的时候是否可以从令牌池中取n个令牌)。如果你需要在事件超出频率的时候丢弃或跳过事件，就使用AllowN,否则使用Reserve或Wait.

func (*Limiter) Reserve
func (lim *Limiter) Reserve() *Reservation
Reserve是ReserveN(time.Now(), 1).的简化形式。

func (*Limiter) ReserveN
func (lim *Limiter) ReserveN(now time.Time, n int) *Reservation
ReserveN 返回对象Reservation ，标识调用者需要等多久才能等到n个事件发生(意思就是等多久令牌池中至少含有n个令牌)。

如果ReserveN 传入的n大于令牌池的容量b，那么返回false.
使用样例如下：

r := lim.ReserveN(time.Now(), 1)
if !r.OK() {
  // Not allowed to act! Did you remember to set lim.burst to be > 0 ?我只要1个事件发生仍然返回false，是不是b设置为了0？
  return
}
time.Sleep(r.Delay())
Act()
如果希望根据频率限制等待和降低事件发生的速度而不丢掉事件，就使用这个方法。
我认为这里要表达的意思就是如果事件发生的频率是可以由调用者控制的话，可以用ReserveN 来控制事件发生的速度而不丢掉事件。如果要使用context的截止日期或cancel方法的话，使用WaitN。

func (*Limiter) Wait
func (lim *Limiter) Wait(ctx context.Context) (err error)
Wait是WaitN(ctx, 1)的简化形式。

func (*Limiter) WaitN
func (lim *Limiter) WaitN(ctx context.Context, n int) (err error)
WaitN 阻塞当前直到lim允许n个事件的发生。

如果n超过了令牌池的容量大小则报错。
如果Context被取消了则报错。
如果lim的等待时间超过了Context的超时时间则报错。

