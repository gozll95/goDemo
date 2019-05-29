
import (
	"math"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go/log"
)

// *********************************关于 Tracer*******************************************
/*
 * 有一个全局的私有的 Tracer: globalTracer
 */
var (
	globalTracer = registeredTracer{NoopTracer{}, false}
)

/*
 * 注册Tracer(tracer+是否被注册)
 */
type registeredTracer struct {
	tracer       Tracer
	isRegistered bool
}

/*
 * Tracer Interface
 */
type Tracer interface {
	// 创建一个span
	StartSpan(operationName string, opts ...StartSpanOption) Span
	// 将sm以format形式携带到carrier
	Inject(sm SpanContext, format interface{}, carrier interface{}) error
	//
	Extract(format interface{}, carrier interface{}) (SpanContext, error)
}

/*
 * StartSpanOption的若干个实例可以被传递到Tracer.StartSpan中
 * 这个是个很赞的方法👍👍👍
 */
type StartSpanOption interface {
	Apply(*StartSpanOptions)
}

/*
 * StartSpanOption
 */
type StartSpanOptions struct {
	// 若干个,Span之间的关系
	References []SpanReference
	StartTime  time.Time
	Tags       map[string]interface{}
}

// *********************************关于 Span 之间的关系*******************************************

/*
 * Span之间的关系
 */
type SpanReference struct {
	Type              SpanReferenceType
	ReferencedContext SpanContext
}

/*
 * 两种span之间的关系
 */
type SpanReferenceType int

const (
	ChildOfRef SpanReferenceType = iota
	FollowsFromRef
)

func ChildOf(sc SpanContext) SpanReference
func FollowsFrom(sc SpanContext) SpanReference

/*
 * 满足StartSpanOption接口
 * 将自己的reference添加到已有的reference
 */
func (r SpanReference) Apply(o *StartSpanOptions) {
	if r.ReferencedContext != nil {
		o.References = append(o.References, r)
	}
}

/*
 * 以下都是满足了StartSpanOption接口
 */
type StartTime time.Time

func (t StartTime) Apply(o *StartSpanOptions)

type Tags map[string]interface{}

func (t Tags) Apply(o *StartSpanOptions)

// ****************************************************************************

/*
 * Span上下文,包含{trace_id,span_id,sampled}必须被传递给后代span
 */
type SpanContext interface {
	// 对上下文中的每个baggae item遍历执行handler
	ForeachBaggageItem(handler func(k, v string) bool)
}

// ****************************************************************************

// *********************************关于 Span*******************************************
/*
 * span
 */
type Span interface {
	Finish()
	FinishWithOptions(opts FinishOptions)
	Context() SpanContext
	SetOperationName(operationName string) Span
	SetTag(key string, value interface{}) Span
	LogFields(fields ...log.Field)
	LogKV(alternatingKeyValues ...interface{})
	SetBaggageItem(restrictedKey, value string) Span
	BaggageItem(restrictedKey string) string
	Tracer() Tracer
}

type FinishOptions struct {
	FinishTime time.Time
	LogRecords []LogRecord
}

type LogRecord struct {
	Timestamp time.Time
	Fields    []log.Field
}

// ****************************************************************************

// *********************************关于 Context*******************************************

// ****************************************************************************

// *********************************关于 Propagation *******************************************
type BuiltinFormat byte

const (
	// carrier 为 io.Writer/io.Reader
	Binary BuiltinFormat = iota
	// carrier 为 TextMapWriter/TextMapReader
	TextMap
	HTTPHeaders
)

type TextMapWriter interface {
	Set(key, val string)
}

/* 
 * 高度抽象的方法 👍👍👍
 */ 
type TextMapReader interface {
	// The "foreach" callback pattern reduces unnecessary copying in some cases
	// and also allows implementations to hold locks while the map is read.
	ForeachKey(handler func(key, val string) error) error
}

/*
 * 具体实现
 */
type TextMapCarrier map[string]string

func (c TextMapCarrier) ForeachKey(handler func(key, val string) error) error
func (c TextMapCarrier) Set(key, val string)

type HTTPHeadersCarrier http.Header

func (c HTTPHeadersCarrier) Set(key, val string)
func (c HTTPHeadersCarrier) ForeachKey(handler func(key, val string) error) error

// ********************************* Log *******************************************
/*
 * Log
 */
type Field struct {
	key string
	// 形式
	fieldType    fieldType
	numericVal   int64
	stringVal    string
	interfaceVal interface{}
}

const (
	stringType fieldType = iota
	boolType
	intType
	int32Type
	uint32Type
	int64Type
	uint64Type
	float32Type
	float64Type
	errorType
	objectType
	lazyLoggerType
	noopType
)

/*
 * 允许用户定义,晚一些,👍👍👍
 */
type LazyLogger func(fv Encoder)

func Lazy(ll LazyLogger) Field {
	return Field{
		fieldType:    lazyLoggerType,
		interfaceVal: ll,
	}
}

/*
 * 不被Tracer记录的
 */
func Noop() Field

/*
 * Encoder接入Filed的内容
 */
type Encoder interface {
	EmitString(key, value string)
	EmitBool(key string, value bool)
	EmitInt(key string, value int)
	EmitInt32(key string, value int32)
	EmitInt64(key string, value int64)
	EmitUint32(key string, value uint32)
	EmitUint64(key string, value uint64)
	EmitFloat32(key string, value float32)
	EmitFloat64(key string, value float64)
	EmitObject(key string, value interface{})
	EmitLazyLogger(value LazyLogger)
}

/*
 * 根据Field实例的类型去执行对应的方法,👍👍👍
 */
func (lf Field) Marshal(visitor Encoder) {
	switch lf.fieldType {
	case stringType:
		visitor.EmitString(lf.key, lf.stringVal)
	case boolType:
		visitor.EmitBool(lf.key, lf.numericVal != 0)
	case intType:
		visitor.EmitInt(lf.key, int(lf.numericVal))
	case int32Type:
		visitor.EmitInt32(lf.key, int32(lf.numericVal))
	case int64Type:
		visitor.EmitInt64(lf.key, int64(lf.numericVal))
	case uint32Type:
		visitor.EmitUint32(lf.key, uint32(lf.numericVal))
	case uint64Type:
		visitor.EmitUint64(lf.key, uint64(lf.numericVal))
	case float32Type:
		visitor.EmitFloat32(lf.key, math.Float32frombits(uint32(lf.numericVal)))
	case float64Type:
		visitor.EmitFloat64(lf.key, math.Float64frombits(uint64(lf.numericVal)))
	case errorType:
		if err, ok := lf.interfaceVal.(error); ok {
			visitor.EmitString(lf.key, err.Error())
		} else {
			visitor.EmitString(lf.key, "<nil>")
		}
	case objectType:
		visitor.EmitObject(lf.key, lf.interfaceVal)
	case lazyLoggerType:
		visitor.EmitLazyLogger(lf.interfaceVal.(LazyLogger))
	case noopType:
		// intentionally left blank
	}
}

// ****************************************************************************

// *************************************** mock span ***********************************

/*
 * Span
 */
type MockSpan struct {
	sync.RWMutex

	ParentID int

	OperationName string
	StartTime     time.Time
	FinishTime    time.Time

	// All of the below are protected by the embedded RWMutex.
	SpanContext MockSpanContext
	tags        map[string]interface{}
	logs        []MockLogRecord
	tracer      *MockTracer
}

/*
 * MockSpanContext,要传递下去
 */
type MockSpanContext struct {
	TraceID int
	SpanID  int
	Sampled bool
	Baggage map[string]string
}

/*
 * 根据Field实例的类型去执行对应的方法,👍👍👍
 */
func (s *MockSpan) logFieldsWithTimestamp(ts time.Time, fields ...log.Field) {
	lr := MockLogRecord{
		Timestamp: ts,
		Fields:    make([]MockKeyValue, len(fields)),
	}
	// 对每个MockKeyValue做处理
	for i, f := range fields {
		outField := &(lr.Fields[i])
		f.Marshal(outField)
	}

	s.Lock()
	defer s.Unlock()
	s.logs = append(s.logs, lr)
}

type MockLogRecord struct {
	Timestamp time.Time
	Fields    []MockKeyValue
}

/*
 * MockKeyValue实现了log.Encoder的interface,👍👍👍
 */
type MockKeyValue struct {
	Key string

	// 擅长使用kind来标识不同的类型
	ValueKind   reflect.Kind
	ValueString string
}

// ****************************************************************************

// *************************************** mock tracer ***********************************
/*
 * Tracer
 */
type MockTracer struct {
	sync.RWMutex
	finishedSpans []*MockSpan
	injectors     map[interface{}]Injector
	extractors    map[interface{}]Extractor
}

/*
 * 高度抽象
 */
type Injector interface {
	Inject(ctx MockSpanContext, carrier interface{}) error
}

type Extractor interface {
	Extract(carrier interface{}) (MockSpanContext, error)
}

func (t *MockTracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	spanContext, ok := sm.(MockSpanContext)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}
	injector, ok := t.injectors[format]
	if !ok {
		return opentracing.ErrUnsupportedFormat
	}
	return injector.Inject(spanContext, carrier)
}

// ****************************************************************************

/*
 * 👍👍👍
 */
func (m *MockKeyValue) EmitLazyLogger(value log.LazyLogger) {
	var meta MockKeyValue
	value(&meta)
	m.Key = meta.Key
	m.ValueKind = meta.ValueKind
	m.ValueString = meta.ValueString
}

// 大量使用 初始化a, o.Apply(a)修改了a,

// 这里所有的handler func 都给上层很大的灵活性 👍👍👍