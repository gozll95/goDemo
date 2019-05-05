package log

import (
	"go-kit-gin/service"
)

//LoggingMiddleware 日志中间件
func LoggingMiddleware() service.SvcMiddleware {
	return func(next service.AppService) service.AppService {
		return logmw{next}
	}
}

type logmw struct {
	service.AppService
}
