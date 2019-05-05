package config

const (
	// Service ports
	ItemSaverPort = 1234
	// ElasticSearch
	ElasticIndex = "dating_profile"
	// RPC Endpoints
	ItemSaverRpc    = "ItemSaverService.Save"
	CrawlServiceRpc = "CrawlService.Process"

	// Parser names
	ParseCity     = "ParseCity"
	ParseCityList = "ParseCityList"
	ParseProfile  = "ParseProfile"
	NilParser     = "NilParser"

	WorkerPort0 = 1235

	// Rate limiting
	Qps = 20
)
