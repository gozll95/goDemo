
type T struct {
	a int
}

// method type
func (t T) Get() int       { return t.a }
func (t *T) Set(a int) int { t.a = a; return t.a }

var t T
t.Get()
t.Set(1)

// method expression 类型名调用方法
var t T
T.Get(t)
(*T).Set(&t, 1)

f1 := (*T).Set //函数类型：func (t *T, int)int
f2 := T.Get //函数类型：func(t T)int
f1(&t, 3)
fmt.Println(f2(t))


f3 := (&t).Set //函数类型：func(int)int
f3(4)
f4 := t.Get//函数类型：func()int   
fmt.Println(f4())
