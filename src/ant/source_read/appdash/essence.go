import (
	"sync"
	"time"
)
//************************************ Store ********************************
/*
 * interface:存储和检索spans
 * Store: 满足了Collector的interface
 */
type Store interface {
	Collector

	// Trace gets a trace (a tree of spans) given its trace ID. If no
	// such trace exists, ErrTraceNotFound is returned.
	Trace(ID) (*Trace, error)
}

type Collector interface {
	Collect(SpanID, ...Annotation) error
}

/*
 * 列出traces的参数
 */
type TracesOpts struct {
	Timespan Timespan
	TraceIDs []ID
}

/*
 * interface:Query,输入query参数,输出对应的Traces
 */
type Queryer interface {
	// Traces returns an implementation-defined list of traces according to the options.
	Traces(opts TracesOpts) ([]*Trace, error)
}

/*
 * AggregatedResult:一组traces的被关心的metrics
 */
type AggregatedResult struct {
	RootSpanName string
	Average, Min, Max, StdDev time.Duration
	Samples int64
	Slowest []ID
}
/*
 * interface: Aggregatr是一种类型的store,可以汇总trace data
 */
type Aggregator interface {
	Aggregate(start, end time.Duration) ([]*AggregatedResult, error)
}

/*
 * Memory store
 */
 type MemoryStore struct {
	trace map[ID]*Trace        // trace ID -> trace tree
	span  map[ID]map[ID]*Trace // trace ID -> span ID -> trace (sub)tree

	sync.Mutex // protects trace

	log bool
}




/*
 * interface: PersistentStore,可以持久化到r/w中
 */
 type PersistentStore interface {
	Write(io.Writer) error
	ReadFrom(io.Reader) (int64, error)
	Store
}

/*
 * Method:
 *		  // 定时持久化到文件中
 *        func PersistEvery(s PersistentStore, interval time.Duration, file string) error 
 */

 /*
  * interface: DeleteStore
  */
type DeleteStore interface {
	Store
	Delete(...ID) error
}


/*
 * Store的实现一: 基于时间过期
 */

type RecentStore struct {
	//
	MinEvictAge time.Duration
	DeleteStore
	// 是否打印debug log
	Debug bool
	// trace ID 第一次的 unix time
	created map[ID]int64
	//
	lastEvicted time.Time
	mu          sync.Mutex
}

/*
 * Store的实现二: 环式
 */
 type LimitStore struct {
	 // 最多保持多少trace
	Max int 
	DeleteStore
	mu            sync.Mutex
	// traces的set,用于快速觉得trace是否在ring里
	traces        map[ID]struct{} 
	// 环
	ring          []int64  
	// ring下一次插的位置
	nextInsertIdx int    
}

/*
 * Collect(id SpanID, anns ...Annotation) error
 *        traces里没有-->插到ring里,将ring里被替换的trace删除
 *        traces里有-->直接调用 store的collect
 */


//*************************************************************************






//************************************ Span ********************************
const (
	SpanIDDelimiter = "/"
)


/*
 * SpanID
 */
type SpanID struct {
	Trace  ID
	Span   ID
	Parent ID
}

/*
 * Method:是否是root Span
 */


/*
 * Span
 */
type Span struct {
	// ID probabilistically uniquely identifies this span.
	ID SpanID

	Annotations
}

/*
 * Method:Span<-->Wire_Span
 */


/*
 * Annotations
 */
 type Annotations []Annotation

 type Annotation struct {
	Key   string
	Value []byte
}

/*
 * Method:Annotations<-->Wire_Anno
 *        get key
 */


//*************************************************************************



//************************************Trace********************************
// Trace 是一颗树
type Trace struct {
	Span          // Root span
	Sub  []*Trace // Children
}

/*
 * Method: String()string // string化输出
 *         FindSpan(spanID ID) *Trace // 根据spanId找到子树
 *         treeString(w io.Writer, depth int)
 *         TimespanEvent() (TimespanEvent, error) // tracer里的annotations->event->findTraceTimes(events)
 * 		   findTraceTimes(events []Event) (start, end time.Time, _ bool) // 得到trace的时间跨度
 */


/*
 * tracesByIDSpan
 */
type tracesByIDSpan []*Trace

func (t tracesByIDSpan) Len() int           { return len(t) }
func (t tracesByIDSpan) Less(i, j int) bool { return t[i].Span.ID.Span < t[j].Span.ID.Span }
func (t tracesByIDSpan) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

//*************************************************************************






type ID uint64

// Compile-time "implements" check.
var _ interface {
	Store
	Queryer
} = (*MemoryStore)(nil)



//************************************Event***********************************

// 有很多注册的Event
var registeredEvents = map[string]Event{} // event schema -> event type

func init() {
	RegisterEvent(SpanNameEvent{})
	RegisterEvent(logEvent{})
	RegisterEvent(msgEvent{})
	RegisterEvent(timespanEvent{})
	RegisterEvent(Timespan{})
}

/*
 * Event
 */
type Event interface {
	Schema() string
}

/*
 * 以下是几种实现了Event的类型
 */
type TimespanEvent interface {
	Event
	Start() time.Time
	End() time.Time
}

type timespanEvent struct {
	S, E time.Time
}

func (timespanEvent) Schema() string      { return "timespan" }
func (ev timespanEvent) Start() time.Time { return ev.S }
func (ev timespanEvent) End() time.Time   { return ev.E }

//*************************************************************************

Annotation->schema->Event



//************************************ Collector ********************************
/*
 * interface: Collector
 */
type Collector interface {
	Collect(SpanID, ...Annotation) error
}

/*
 * Collector接口的封装一: ChunkedCollector,用于缓存一些collections,然后整体发出去
 * 使用场景: 一些对时间敏感的collector
 */
type ChunkedCollector struct {
	Collector
	// Flush的最小时间间隔
	MinInterval time.Duration
	// Flush操作的timeout
	FlushTimeout time.Duration
	// 最大的queue size
	MaxQueueSize uint64
	Log *log.Logger
	// Flush操作开始时执行的操作
	OnFlush func(queueSize int)
	lastErr error
	started, stopped bool
	stopChan         chan struct{}
	queueSizeBytes  uint64
	pendingBySpanID map[SpanID]Annotations
	mu sync.Mutex
}

/*
 * Start()
 */
 func (cc *ChunkedCollector) start() {
	cc.stopChan = make(chan struct{})
	cc.started = true
	go func() {
		for {
			t := time.After(cc.MinInterval)
			select {
			case <-t:
				if err := cc.Flush(); err != nil {
					cc.mu.Lock()
					cc.lastErr = err
					cc.mu.Unlock()
				}
			case <-cc.stopChan:
				return // stop
			}
		}
	}()
}

/*
 * Flush将 pendingBySpanID 一次性推到collector中
 */


/*
 * Collector接口的封装二,👍👍👍,用于protobuf
 */
 type RemoteCollector struct {
	addr string
	dial func() (net.Conn, error)
	mu    sync.Mutex      // guards pconn
	pconn pio.WriteCloser // delimited-protobuf remote connection
	Log   *log.Logger
	logMu sync.Mutex
	Debug bool 	// Debug is whether to log debug messages.
}

/*
 * 这里定义了dial方法
 */
func NewTLSRemoteCollector(addr string, tlsConfig *tls.Config) *RemoteCollector {
	return &RemoteCollector{
		addr: addr,
		dial: func() (net.Conn, error) {
			return tls.Dial("tcp", addr, tlsConfig)
		},
	}
}


/*
 * 如果pconn已经有了就不用重新连了,👍👍👍
 * 重连的话会重置rc.pconn->nil->dail()
 * 如果连接有问题会重连一次
 */
 func (rc *RemoteCollector) collectAndRetry(p *wire.CollectPacket) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.pconn != nil {
		if err := rc.collect(p); err == nil {
			return nil
		}
		if rc.Debug {
			rc.log().Printf("Reconnecting to send %v", spanIDFromWire(p.Spanid))
		}
	}
	if err := rc.connect(); err != nil {
		return err
	}
	return rc.collect(p)
}


/*
 * Collector接口的封装三,👍👍👍
 * 用于从client端接收序列化的,再存储到collector中
 */
 type CollectorServer struct {
	c Collector
	l net.Listener
	Log   *log.Logger
	logMu sync.Mutex
	Debug bool
	Trace bool
}

//*************************************************************************




//************************************Tracer********************************
/*
 * Options
 */
type Options struct {
	// 采样的方法 👍: 使用方法是个非常灵活的实现
	ShouldSample func(traceID uint64) bool
	Verbose bool
	Logger *log.Logger
}
//*************************************************************************


//************************************ Recorder ********************************
/*
 * Recorder结构体
 */
type Recorder struct {
	collector appdash.Collector
	logOnce   sync.Once
	verbose   bool
	Log       *log.Logger
}

//*************************************************************************

//************************************ Event ********************************
/*
 * Event interface
 * 事件,比如: sql event
 */
type Event interface {
	Schema() string
}

type ImportantEvent interface {
	Important() []string
}

/*
 * Marshaler interface: 将event->Anno
 */
type EventMarshaler interface {
	MarshalEvent() (Annotations, error)
}


/*
 * EventUnmarshaler interface: 将Anno->event
 */
type EventUnmarshaler interface {
	UnmarshalEvent(Annotations) (Event, error)
}

//*************************************************************************

//************************************高端玩法********************************

/*
 * reflect的玩法
 * *[]interface的玩法
 */
func UnmarshalEvents(anns Annotations, events *[]Event) error {
	schemas := anns.schemas()
	for _, schema := range schemas {
		ev := registeredEvents[schema]
		if ev == nil {
			continue
		}
		evv := reflect.New(reflect.TypeOf(ev))
		if err := UnmarshalEvent(anns, evv.Interface().(Event)); err != nil {
			return err
		}
		*events = append(*events, evv.Elem().Interface().(Event))
	}
	return nil
}

/*
 * append的玩法
 */
func (id SpanID) Format(s string, args ...interface{}) string {
	args = append([]interface{}{id.String()}, args...)
	return fmt.Sprintf(s, args...)
}
//*************************************************************************

wire<->Span

Annotation <--> wire_Anno




ioutil.TempFile("", "appdash")

ls.nextInsertIdx = (ls.nextInsertIdx + 1) % ls.Max // increment & wrap


flattenValue("", reflect.ValueOf(e), func(k, v string) {
	as = append(as, Annotation{Key: k, Value: []byte(v)})
})


func flattenValue(prefix string, v reflect.Value, f func(k, v string)) {
	switch o := v.Interface().(type) {
	case time.Time:
		f(prefix, o.Format(time.RFC3339Nano))
		return
	case time.Duration:
		ms := float64(o.Nanoseconds()) / float64(time.Millisecond)
		f(prefix, strconv.FormatFloat(ms, 'f', -1, 64))
		return
	case fmt.Stringer:
		f(prefix, o.String())
		return
	}

	switch v.Kind() {
	case reflect.Ptr:
		flattenValue(prefix, v.Elem(), f)
	case reflect.Bool:
		f(prefix, strconv.FormatBool(v.Bool()))
	case reflect.Float32, reflect.Float64:
		f(prefix, strconv.FormatFloat(v.Float(), 'f', -1, 64))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f(prefix, strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f(prefix, strconv.FormatUint(v.Uint(), 10))
	case reflect.String:
		f(prefix, v.String())
	case reflect.Struct:
		for i, name := range fieldNames(v) {
			flattenValue(nest(prefix, name), v.Field(i), f)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			// small bit of cuteness here: use flattenValue on the key first,
			// then on the value
			flattenValue("", key, func(_, k string) {
				flattenValue(nest(prefix, k), v.MapIndex(key), f)
			})
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			flattenValue(nest(prefix, strconv.Itoa(i)), v.Index(i), f)
		}
	default:
		f(prefix, fmt.Sprintf("%+v", v.Interface()))
	}
}



//************************************蛮好的地方********************************
/*
 * json tag
 */

type SendCmd struct {
	CollectorAddr  string `short:"c" long:"collector" description:"collector listen address" default:":7701"`
	CollectorProto string `short:"p" long:"proto" description:"collector protocol (tcp or tls)" default:"tcp"`
	ServerName     string `short:"s" long:"server-name" description:"server name (required for TLS)"`
	Debug          bool   `short:"d" long:"debug" description:"debug log"`
}
//*************************************************************************