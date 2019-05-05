Logrus
其是一个结构化后的日志包，API完全兼容标准包logger。docker(moby)源码中使用logrus来做日志记录

获取

项目地址: sirupsen/logrus

目前项目地址为sirupsen/logrus，如果你是Sirupsen/logrus，请转到小写的名称项目中

go get github.com/sirupsen/logrus
输出内容格式化

内建日志格式化有：

logrus.TextFormatter
当你登录TTY时，其输出内容会是彩色标注；如果不在TTY，则需要将ForceColors设置为true

logrus.JSONFormatter
以JSON格式为输出

六个日志等级

log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
// 随后会触发os.Exit(1)
log.Fatal("Bye.")
// 随后会触发panic()
log.Panic("I'm bailing.")
基本实例

代码下载： logrus_study.go

package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	// 以JSON格式为输出，代替默认的ASCII格式
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// 以Stdout为输出，代替默认的stderr
	logrus.SetOutput(os.Stdout)
	// 设置日志等级
	logrus.SetLevel(logrus.WarnLevel)
}
func main() {
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	logrus.WithFields(logrus.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")
}
运行：

go run logrus_study.go

<<'COMMENT'
{"level":"warning","msg":"The group's number increased tremendously!","number":122,"omg":true,"time":"2017-09-18T17:53:13+08:00"}
{"level":"fatal","msg":"The ice breaks!","number":100,"omg":true,"time":"2017-09-18T17:53:13+08:00"}
COMMENT
Logger

如果多个地方使用logging，可以创建一个logrus实例Logger

代码下载： logrus_logger.go

package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()


func main() {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
	file.Close()
}
执行后会在GOPATH路径下穿件logrus.log文件，内容如下：

time="2018-09-18T18:09:32+08:00" level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
Fields

如果有固定Fields，可以创建一个logrus.Entry

requestLogger := log.WithFields(log.Fields{"request_id": request_id, "user_ip": user_ip})
requestLogger.Info("something happened on that request") # will log request_id and user_ip
requestLogger.Warn("something not great happened")
从函数WithFields中可以看出其会返回*Entry，Entry中会包含一些变量

// WithFields函数
func WithFields(fields Fields) *Entry {
	return std.WithFields(fields)
}

// Entry结构体
type Entry struct {
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// When formatter is called in entry.log(), an Buffer may be set to entry
	Buffer *bytes.Buffer
}
Hooks

可以与外面的控件联合，例如

使用github.com/multiplay/go-slack与slack/bearchat一些企业团队协作平台/软件联合使用
使用https://github.com/zbindenren/logrus_mail可以发送email，例如以下实例
logrus_mail.go

安装go包：

go get github.com/zbindenren/logrus_mail
代码下载： logrus_email.go

package main

import (
	"time"

	"github.com/logrus_mail"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	hook, err := logrus_mail.NewMailAuthHook(
		"logrus_email",
		"smtp.gmail.com",
		587,
		"chenjian158978@gmail.com",
		"271802559@qq.com",
		"chenjian158978@gmail.com",
		"xxxxxxx",
	)
	if err == nil {
		logger.Hooks.Add(hook)
	}
	//生成*Entry
	var filename = "123.txt"
	contextLogger := logger.WithFields(logrus.Fields{
		"file":    filename,
		"content": "GG",
	})
	//设置时间戳和message
	contextLogger.Time = time.Now()
	contextLogger.Message = "这是一个hook发来的邮件"
	//只能发送Error,Fatal,Panic级别的log
	contextLogger.Level = logrus.ErrorLevel

	//使用Fire发送,包含时间戳，message
	hook.Fire(contextLogger)
}

邮件截图:


Thread safety

默认的Logger在并发写入的时候由mutex保护，其在调用hooks和写入logs时被启用。如果你认为此锁是没有必要的，可以添加logger.SetNoLock()来让锁失效。