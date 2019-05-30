//************************************ 精华 ********************************
"github.com/spf13/viper"

strings.NewReplacer("-", "_", ".", "_")

//*************************************************************************



//************************************ memoryStore ********************************
/*
 * Options	
 */
 const limit = "memory.max-traces"


/*
 * Memory Options
 */

type Options struct {
	Configuration config.Configuration
}

/*
 * Config
 */
 type Configuration struct {
	MaxTraces int `yaml:"max-traces"`
}

/*
 * 初始化from Viper
 */
 func (opt *Options) InitFromViper(v *viper.Viper) {
	opt.Configuration.MaxTraces = v.GetInt(limit)
}


/*
 * 初始化 from CLI
 */
 func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.Int(limit, opt.Configuration.MaxTraces, "The maximum amount of traces to store in memory")
}


/*
 * Store实体
 */
 type Store struct {
	sync.RWMutex
	// 结合index维护一个环形队列
	ids        []*model.TraceID
	// 	m.traces[span.TraceID].Spans = append(m.traces[span.TraceID].Spans, span)
	traces     map[model.TraceID]*model.Trace
	// m.services[span.Process.ServiceName] = struct{}{}
	services   map[string]struct{}
	// m.operations[span.Process.ServiceName][span.OperationName] = struct{}{}
	operations map[string]map[string]struct{}
	deduper    adjuster.Adjuster
	config     config.Configuration
	// 目前保存的游标,这里维护的是一个环形队列
	index      int
}

/*
 * WriteSpan(span *model.Span)error
 */
- m.operations[span.Process.ServiceName][span.OperationName] = struct{}{}
- m.services[span.Process.ServiceName] = struct{}{}
- ids + index + m.config.MaxTraces 维护一个环形队列,插入新的,删除旧的
- m.traces[span.TraceID].Spans = append(m.traces[span.TraceID].Spans, span)



/*
 * Methods:
 *			查询svcs
 * 			查询svc下的operations name
 * 		   	根据query查询traces/trace_ids
 * 
 */
span.Tags
span.Process.Tags
span.Logs


/*
 * Factory
 */
 type Factory struct {
	options        Options
	metricsFactory metrics.Factory
	logger         *zap.Logger
	store          *Store
}

/*
 * 通过工厂方法初始化store
 */
 func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger
	f.store = WithConfiguration(f.options.Configuration)
	logger.Info("Memory storage initialized", zap.Any("configuration", f.store.config))
	return nil
}

// 实现了
spanstore.Reader
spanstore.Writer
dependencystore.Reader



// 疑问: deduper    adjuster.Adjuster


//*************************************************************************

// 思考如果自己做的话,是怎么做,然后再看开源的实现


//************************************ Config 这是 Tracer的配置 ********************************
config.FromEnv


// Configuration configures and creates Jaeger Tracer
type Configuration struct {
	// JAEGER_SERVICE_NAME
	ServiceName string `yaml:"serviceName"`
	// JAEGER_DISABLED
	Disabled bool `yaml:"disabled"`
	// AEGER_RPC_METRICS
	RPCMetrics bool `yaml:"rpc_metrics"`
	// JAEGER_TAGS
	Tags []opentracing.Tag `yaml:"tags"`
	Sampler             *SamplerConfig             `yaml:"sampler"`
	Reporter            *ReporterConfig            `yaml:"reporter"`
	Headers             *jaeger.HeadersConfig      `yaml:"headers"`
	BaggageRestrictions *BaggageRestrictionsConfig `yaml:"baggage_restrictions"`
	Throttler           *ThrottlerConfig           `yaml:"throttler"`
}

// 采样配置
type SamplerConfig struct {
	// sampler的类型: const, probabilistic, rateLimiting,remote
	// 同样可以被JAEGER_SAMPLER_TYPE定义
	Type string `yaml:"type"`

	// Param is a value passed to the sampler.
	// Valid values for Param field are:
	// - for "const" sampler, 0 or 1 for always false/true respectively
	// - for "probabilistic" sampler, a probability between 0 and 1
	// - for "rateLimiting" sampler, the number of spans per second
	// - for "remote" sampler, param is the same as for "probabilistic"
	//   and indicates the initial sampling rate before the actual one
	//   is received from the mothership.
	// Can be set by exporting an environment variable named JAEGER_SAMPLER_PARAM
	Param float64 `yaml:"param"`

	// jaegrer-agent的HTTP采样server
	// Can be set by exporting an environment variable named JAEGER_SAMPLER_MANAGER_HOST_PORT
	SamplingServerURL string `yaml:"samplingServerURL"`

	// sampler可以追踪的最大operation数目
	// MaxOperations is the maximum number of operations that the sampler
	// will keep track of. If an operation is not tracked, a default probabilistic
	// sampler will be used rather than the per operation specific sampler.
	// Can be set by exporting an environment variable named JAEGER_SAMPLER_MAX_OPERATIONS
	MaxOperations int `yaml:"maxOperations"`

	// 这个字段sampler从jaeger-agent获取采样策略的频率
	// SamplingRefreshInterval controls how often the remotely controlled sampler will poll
	// jaeger-agent for the appropriate sampling strategy.
	// Can be set by exporting an environment variable named JAEGER_SAMPLER_REFRESH_INTERVAL
	SamplingRefreshInterval time.Duration `yaml:"samplingRefreshInterval"`
}

// reporter配置
type ReporterConfig struct {
	// reporter的缓存span的大小,这个缓存队列会被一个后台goroutine抽取
	// Can be set by exporting an environment variable named JAEGER_REPORTER_MAX_QUEUE_SIZE
	QueueSize int `yaml:"queueSize"`

	// 缓存队列多久被强制flushed
	// Can be set by exporting an environment variable named JAEGER_REPORTER_FLUSH_INTERVAL
	BufferFlushInterval time.Duration

	// 如果被置为true,则允许LoggingReporter可以与main reporter并行运行,并且记录所有已经提交的spans
	// Can be set by exporting an environment variable named JAEGER_REPORTER_LOG_SPANS
	LogSpans bool `yaml:"logSpans"`

	// reporter发送span给jaeger-agent的地址
	// Can be set by exporting an environment variable named JAEGER_AGENT_HOST / JAEGER_AGENT_PORT
	LocalAgentHostPort string `yaml:"localAgentHostPort"`

	// reporter发送span给jaeger-collector
	// Can be set by exporting an environment variable named JAEGER_ENDPOINT
	CollectorEndpoint string `yaml:"collectorEndpoint"`

	// http basic 认证 user
	// Can be set by exporting an environment variable named JAEGER_USER
	User string `yaml:"user"`

	// http basic 认证 passwd
	// jaeger-collector. Can be set by exporting an environment variable named JAEGER_PASSWORD
	Password string `yaml:"password"`
}


// Baggage配置
type BaggageRestrictionsConfig struct {
	// 暂时忽略吧
	DenyBaggageOnInitializationFailure bool `yaml:"denyBaggageOnInitializationFailure"`
	// jaeger-agent的baggage restrictions server的地址
	HostPort string `yaml:"hostPort"`
	// 控制多久baggage restriction manager将从 jaeger-agent 那里poll 最新的 baggage 策略
	RefreshInterval time.Duration `yaml:"refreshInterval"`
}

// Throttler配置
type ThrottlerConfig struct {
	// jaeger-agent的受信server
	HostPort string `yaml:"hostPort"`
	// 多久从jaeger-agent去poll节流策略
	RefreshInterval time.Duration `yaml:"refreshInterval"`
	// 是否一看到一个operation就从jaeger-agent去同步credits
	SynchronousInitialization bool `yaml:"synchronousInitialization"`
}



/*
 * 从Cfg里利用Option生成Tracer,一个io.Closer用于在shutdown之前flush buffers
 */
 func (c Configuration) NewTracer(options ...Option) (opentracing.Tracer, io.Closer, error)
- 设置ServerName
- opts := applyOptions(options...)
- metrics
- 如果 c.RPCMetrics => 将rpc observers 加入 c.observers
- 一些字段的判断,if nil -> init default
	- sampler // 如果opt.sample为nil则生成默认sampler
	- reporter // 如果opt.reporter为nil则生成默认reporter
	- 生成 []jaeger.TracerOption
- 	tracer, closer := jaeger.NewTracer(
	c.ServiceName,
	sampler,
	reporter,
	tracerOptions...,
)


/*
 * 根据 config 生成sampler
 * TODO: 🍠🍠 这里需要注意有 串行化 options
 */
 func (sc *SamplerConfig) NewSampler(
	serviceName string,
	metrics *jaeger.Metrics,
) (jaeger.Sampler, error)


/*
 * 根据 config 生成 reporter
 * 这里先生成了sender
 */

func (rc *ReporterConfig) NewReporter(
	serviceName string,
	metrics *jaeger.Metrics,
	logger jaeger.Logger,
) (jaeger.Reporter, error) {
	sender, err := rc.newTransport()
	if err != nil {
		return nil, err
	}
	// 这里再次借助了option
	reporter := jaeger.NewRemoteReporter(
		sender,
		jaeger.ReporterOptions.QueueSize(rc.QueueSize),
		jaeger.ReporterOptions.BufferFlushInterval(rc.BufferFlushInterval),
		jaeger.ReporterOptions.Logger(logger),
		jaeger.ReporterOptions.Metrics(metrics))
	// 链式👍👍👍
	if rc.LogSpans && logger != nil {
		logger.Infof("Initializing logging reporter\n")
		reporter = jaeger.NewCompositeReporter(jaeger.NewLoggingReporter(logger), reporter)
	}
	return reporter, err
}

/*
 * 生成transport/sender
 */

func (rc *ReporterConfig) newTransport() (jaeger.Transport, error) {
	switch {
	// 认证的http
	case rc.CollectorEndpoint != "" && rc.User != "" && rc.Password != "":
		return transport.NewHTTPTransport(rc.CollectorEndpoint, transport.HTTPBatchSize(1),
			transport.HTTPBasicAuth(rc.User, rc.Password)), nil
	// 非认证的http
	case rc.CollectorEndpoint != "":
		return transport.NewHTTPTransport(rc.CollectorEndpoint, transport.HTTPBatchSize(1)), nil
	// udp
	default:
		return jaeger.NewUDPTransport(rc.LocalAgentHostPort, 0)
	}
}


//*************************************************************************


//************************************ Demo ********************************
metrics
log
observer

//*************************************************************************


//************************************ Option ********************************
// 这里充分展示了interface的特性
type Options struct {
	// 发出span的数据
	metrics             metrics.Factory
	// 用于log reporter errors 和 span 的 logspans
	logger              jaeger.Logger
	reporter            jaeger.Reporter
	sampler             jaeger.Sampler
	// 被Tracer注册去收集新span的提醒
	contribObservers    []jaeger.ContribObserver
	// 被Tracer注册去收集新span的提醒
	observers           []jaeger.Observer
	// 是否生成128bit的trace ID
	gen128Bit           bool
	// 是否按照zipkin的方式client和server共享一个span
	zipkinSharedRPCSpan bool
	// tag value的最大值
	maxTagValueLength   int
	tags                []opentracing.Tag
	// injector注册
	injectors           map[interface{}]jaeger.Injector
	// extractor注册
	extractors          map[interface{}]jaeger.Extractor
}


// Option是一种方法: 在client端设置一些选择 👍👍👍 
// 作用于 Options
type Option func(c *Options)  // 

// demo
func Metrics(factory metrics.Factory) Option {
	return func(c *Options) {
		c.metrics = factory
	}
}

// 很好的一种应用方式
func applyOptions(options ...Option) Options {
	opts := Options{
		injectors:  make(map[interface{}]jaeger.Injector),
		extractors: make(map[interface{}]jaeger.Extractor),
	}
	for _, option := range options {
		option(&opts)
	}
	if opts.metrics == nil {
		opts.metrics = metrics.NullFactory
	}
	if opts.logger == nil {
		opts.logger = jaeger.NullLogger
	}
	return opts
}


//************************************************************************* 


//************************************ Sampler ********************************

// interface: 决定一个新的trace是否应该被采样
type Sampler interface {
	// 是否应该被采样,返回采样的参数
	IsSampled(id TraceID, operation string) (sampled bool, tags []Tag)
	// stop 所有的后台 goroutine
	Close()
}

/*
 * 
 */
type ConstSampler struct {
	Decision bool
	tags     []Tag
}

/*
 * 这里????不是很懂
 */
// TraceID represents unique 128bit identifier of a trace
type TraceID struct {
	High, Low uint64
}


//*************************************************************************** 


//************************************ Reporter ********************************
// 当一个span已经完成的时候, tracer 调用 Reporter 将 span 发送到 collector 中
type Reporter interface {
	Report(span *Span)
	Close()
}


//*************************************************************************


//************************************ Transport ********************************
// Transport指定了跨进程send span的方法
// 实体不需要线程安全
// remote reporter只需要在同一个goroutine里调用Transport方法即可
type Transport interface {
	// 将span->wire,并将其添加到sender的内部buffer,如果buffer超了就调用Flush
	// 如果出现了error那么意味这是一个失败的span,就会被report to metrics
	Append(span *Span) (int, error)
	// Flush提交internal buffer去remote server,如果失败,就report to metrics
	Flush() (int, error)
	io.Closer
}


/*
 * HTTP Transport
 * 实现了interface: Transport 将 span 转发到一个 http server
 */
 type HTTPTransport struct {
	url             string
	client          *http.Client
	batchSize       int
	// thrift形式的span
	spans           []*j.Span
	// thrift形式的process
	process         *j.Process
	httpCredentials *HTTPBasicAuthCredentials
}

/*
 * HTTP basic auth
 */
type HTTPBasicAuthCredentials struct {
	username string
	password string
}

/*
 * Option作用于上述model
 * 很好的方法,👍👍👍
 */
type HTTPOption func(c *HTTPTransport)

// demo
func HTTPTimeout(duration time.Duration) HTTPOption {
	return func(c *HTTPTransport) { c.client.Timeout = duration }
}

/*
 * 一个生成实例的方法,options作用于自身
 */
func NewHTTPTransport(url string, options ...HTTPOption) *HTTPTransport {
	c := &HTTPTransport{
		url:       url,
		client:    &http.Client{Timeout: defaultHTTPTimeout},
		batchSize: 100,
		spans:     []*j.Span{},
	}

	for _, option := range options {
		option(c)
	}
	return c
}

/*
 * 实现 Transport 接口
 */
func (c *HTTPTransport) Append(span *jaeger.Span) (int, error) {
	// 🍠🍠🍠
	if c.process == nil {
		c.process = jaeger.BuildJaegerProcessThrift(span)
	}
	// 生成thfift type的span
	jSpan := jaeger.BuildJaegerThrift(span)
	c.spans = append(c.spans, jSpan)
	// oversize就Flush
	if len(c.spans) >= c.batchSize {
		// thrift send 
		return c.Flush()
	}
	return 0, nil
}

/*
 * 后面还是UDP Transport 需要看下
 */

//*************************************************************************



//************************************ 令牌桶 ********************************
详见jaeger/vendor/github.com/uber/jaeger-client-go/utils/rate_limiter.go
//*************************************************************************