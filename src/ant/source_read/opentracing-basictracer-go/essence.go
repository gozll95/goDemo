import (
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
)

//************************************Tracer********************************
/*
 * Tracer
 */
type Tracer interface {
	opentracing.Tracer
	Options() Options
}

/*
 * Options
 */
type Options struct {
	ShouldSample               func(traceID uint64) bool
	TrimUnsampledSpans         bool
	Recorder                   SpanRecorder
	NewSpanEventListener       func() func(SpanEvent)
	DropAllLogs                bool
	MaxLogsPerSpan             int
	DebugAssertSingleGoroutine bool
	DebugAssertUseAfterFinish  bool
	EnableSpanPool             bool
}

/*
 * tracer实例化
 */
type tracerImpl struct {
	options            Options
	textPropagator     *textMapPropagator
	binaryPropagator   *binaryPropagator
	accessorPropagator *accessorPropagator
}

/*
 * func NewWithOptions(opts Options) opentracing.Tracer
 * func New(recorder SpanRecorder) opentracing.Tracer
 */

/*
 * func (t *tracerImpl) StartSpan(
 *	operationName string,
 *	opts ...opentracing.StartSpanOption,
 * ) opentracing.Span
 */
func (t *tracerImpl) StartSpan(
	operationName string,
	opts ...opentracing.StartSpanOption,
) opentracing.Span {
	sso := opentracing.StartSpanOptions{}
	// 作用于这个👍👍👍
	for _, o := range opts {
		o.Apply(&sso)
	}
	return t.StartSpanWithOptions(operationName, sso)
}

/*
 * sp := t.getSpan() // 这里会通过是否从pool拿的标志位决定是否从sync.Pool中复用spanImpl对象
 */

/*
 * Inject
 */
 func (t *tracerImpl) Inject(sc opentracing.SpanContext, format interface{}, carrier interface{}) error {
	switch format {
	case opentracing.TextMap, opentracing.HTTPHeaders:
		return t.textPropagator.Inject(sc, carrier)
	case opentracing.Binary:
		return t.binaryPropagator.Inject(sc, carrier)
	}
	if _, ok := format.(delegatorType); ok {
		return t.accessorPropagator.Inject(sc, carrier)
	}
	return opentracing.ErrUnsupportedFormat
}

//*************************************************************************

//************************************ Span ********************************
/*
 * interface: Span
 */
type Span interface {
	opentracing.Span
	Operation() string
	Start() time.Time
}

/*
 * Span的实例化
 */
type spanImpl struct {
	tracer     *tracerImpl
	event      func(SpanEvent)
	sync.Mutex // protects the fields below
	raw        RawSpan
	// The number of logs dropped because of MaxLogsPerSpan.
	numDroppedLogs int
}

/*
 * Span的内容
 */
type RawSpan struct {
	Context SpanContext
	ParentSpanID uint64
	Operation string
	Start    time.Time
	Duration time.Duration
	Tags opentracing.Tags
	Logs []opentracing.LogRecord
}

/*
 * Span Context,用于边界传输
 */
type SpanContext struct {
	TraceID uint64
	SpanID uint64
	Sampled bool
	Baggage map[string]string // initialized on first use
}

/*
 * spanPool
 */
var spanPool = &sync.Pool{New: func() interface{} {
	return &spanImpl{}
}}


/*
 * Finish
 * 恰
 */
func (s *spanImpl) FinishWithOptions(opts opentracing.FinishOptions)
	// 其实是在这里记录span的 👍👍👍
	s.tracer.options.Recorder.RecordSpan(s.raw)


/*
 * spanImpl有很多onCreate/onLog/onFinish等方法
 */


//************************************ Propagator ********************************
type textMapPropagator struct {
	tracer *tracerImpl
}
type binaryPropagator struct {
	tracer *tracerImpl
}


/*
 * Inject
 */
 func (p *textMapPropagator) Inject(
	spanContext opentracing.SpanContext,
	opaqueCarrier interface{},
) error 
sc, ok := spanContext.(SpanContext)
carrier, ok := opaqueCarrier.(opentracing.TextMapWriter)
carrier.Set(fieldNameTraceID, strconv.FormatUint(sc.TraceID, 16))


/*
 * Extract
 */
 func (p *textMapPropagator) Extract(
	opaqueCarrier interface{},
) (opentracing.SpanContext, error)
carrier, ok := opaqueCarrier.(opentracing.TextMapReader)
err = carrier.ForeachKey(func(k, v string) error{
	// 对k/v做操作
})


/*
 * 网络大终端写
 */
 func (p *binaryPropagator) Inject(
	spanContext opentracing.SpanContext,
	opaqueCarrier interface{},
) error {
	sc, ok := spanContext.(SpanContext)
	if !ok {
		return opentracing.ErrInvalidSpanContext
	}
	carrier, ok := opaqueCarrier.(io.Writer)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}

	state := wire.TracerState{}
	state.TraceId = sc.TraceID
	state.SpanId = sc.SpanID
	state.Sampled = sc.Sampled
	state.BaggageItems = sc.Baggage

	b, err := proto.Marshal(&state)
	if err != nil {
		return err
	}

	// Write the length of the marshalled binary to the writer.
	length := uint32(len(b))
	if err := binary.Write(carrier, binary.BigEndian, &length); err != nil {
		return err
	}

	_, err = carrier.Write(b)
	return err
}


/*
 * 网络大终端读
 */
 func (p *binaryPropagator) Extract(
	opaqueCarrier interface{},
) (opentracing.SpanContext, error) {
	carrier, ok := opaqueCarrier.(io.Reader)
	if !ok {
		return nil, opentracing.ErrInvalidCarrier
	}

	// Read the length of marshalled binary. io.ReadAll isn't that performant
	// since it keeps resizing the underlying buffer as it encounters more bytes
	// to read. By reading the length, we can allocate a fixed sized buf and read
	// the exact amount of bytes into it.
	var length uint32
	if err := binary.Read(carrier, binary.BigEndian, &length); err != nil {
		return nil, opentracing.ErrSpanContextCorrupted
	}
	buf := make([]byte, length)
	if n, err := carrier.Read(buf); err != nil {
		if n > 0 {
			return nil, opentracing.ErrSpanContextCorrupted
		}
		return nil, opentracing.ErrSpanContextNotFound
	}

	ctx := wire.TracerState{}
	if err := proto.Unmarshal(buf, &ctx); err != nil {
		return nil, opentracing.ErrSpanContextCorrupted
	}

	return SpanContext{
		TraceID: ctx.TraceId,
		SpanID:  ctx.SpanId,
		Sampled: ctx.Sampled,
		Baggage: ctx.BaggageItems,
	}, nil
}


/*
 * DelegatingCarrier:是一个灵活的Carrier,可以被(存储trace metadata并且知道如何序列化它们)
 */
type DelegatingCarrier interface {
	SetState(traceID, spanID uint64, sampled bool)
	State() (traceID, spanID uint64, sampled bool)
	SetBaggageItem(key, value string)
	GetBaggage(func(key, value string))
}


/*
 * 以DelegatingCarrier为载体的
 */
type accessorPropagator struct {
	tracer *tracerImpl
}

func (p *accessorPropagator) Inject(
	spanContext opentracing.SpanContext,
	carrier interface{},
) error {
	dc, ok := carrier.(DelegatingCarrier)
	if !ok || dc == nil {
		return opentracing.ErrInvalidCarrier
	}
	sc, ok := spanContext.(SpanContext)
	if !ok {
		return opentracing.ErrInvalidSpanContext
	}
	dc.SetState(sc.TraceID, sc.SpanID, sc.Sampled)
	for k, v := range sc.Baggage {
		dc.SetBaggageItem(k, v)
	}
	return nil
}

func (p *accessorPropagator) Extract(
	carrier interface{},
) (opentracing.SpanContext, error) {
	dc, ok := carrier.(DelegatingCarrier)
	if !ok || dc == nil {
		return nil, opentracing.ErrInvalidCarrier
	}

	traceID, spanID, sampled := dc.State()
	sc := SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
		Sampled: sampled,
		Baggage: nil,
	}
	dc.GetBaggage(func(k, v string) {
		if sc.Baggage == nil {
			sc.Baggage = map[string]string{}
		}
		sc.Baggage[k] = v
	})

	return sc, nil
}



//*************************************************************************


//************************************精华 ********************************
a衍生了b,b包含a


// NetTraceIntegrator can be passed into a basictracer as NewSpanEventListener
// and causes all traces to be registered with the net/trace endpoint.
var NetTraceIntegrator = func() func(basictracer.SpanEvent) {
	var tr trace.Trace
	return func(e basictracer.SpanEvent) {
		switch t := e.(type) {
		case basictracer.EventCreate:
			tr = trace.New("tracing", t.OperationName)
			tr.SetMaxEvents(1000)
		case basictracer.EventFinish:
			tr.Finish()
		case basictracer.EventTag:
			tr.LazyPrintf("%s:%v", t.Key, t.Value)
		case basictracer.EventLogFields:
			var buf bytes.Buffer
			for i, f := range t.Fields {
				if i > 0 {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%s:%v", f.Key(), f.Value())
			}

			tr.LazyPrintf("%s", buf.String())
		case basictracer.EventLog:
			if t.Payload != nil {
				tr.LazyPrintf("%s (payload %v)", t.Event, t.Payload)
			} else {
				tr.LazyPrintf("%s", t.Event)
			}
		}
	}
}

//*************************************************************************
