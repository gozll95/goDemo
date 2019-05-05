1."请求参数"
type RequestArgs struct {
	// AcceptedDomains 代表可以接受的URL的主域名的列表。
	// URL主域名不在列表中的请求都会被忽略，
	AcceptedDomains []string `json:"accepted_primary_domains"`
	// maxDepth 代表了需要被爬取的最大深度。
	// 实际深度大于此值的请求都会被忽略。
	MaxDepth uint32 `json:"max_depth"`
}

2."数据参数"
type DataArgs struct {
	// ReqBufferCap 代表请求缓冲器的容量。
	ReqBufferCap uint32 `json:"req_buffer_cap"`
	// ReqMaxBufferNumber 代表请求缓冲器的最大数量。
	ReqMaxBufferNumber uint32 `json:"req_max_buffer_number"`
	// RespBufferCap 代表响应缓冲器的容量。
	RespBufferCap uint32 `json:"resp_buffer_cap"`
	// RespMaxBufferNumber 代表响应缓冲器的最大数量。
	RespMaxBufferNumber uint32 `json:"resp_max_buffer_number"`
	// ItemBufferCap 代表条目缓冲器的容量。
	ItemBufferCap uint32 `json:"item_buffer_cap"`
	// ItemMaxBufferNumber 代表条目缓冲器的最大数量。
	ItemMaxBufferNumber uint32 `json:"item_max_buffer_number"`
	// ErrorBufferCap 代表错误缓冲器的容量。
	ErrorBufferCap uint32 `json:"error_buffer_cap"`
	// ErrorMaxBufferNumber 代表错误缓冲器的最大数量。
	ErrorMaxBufferNumber uint32 `json:"error_max_buffer_number"`
}


3."组件参数"
type ModuleArgs struct {
	// Downloaders 代表下载器列表。
	Downloaders []module.Downloader
	// Analyzers 代表分析器列表。
	Analyzers []module.Analyzer
	// Pipelines 代表条目处理管道管道列表。
	Pipelines []module.Pipeline
}


goquery