package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	log = logrus.New()
)

func init() {
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logrus.SetOutput(os.Stdout)
	//设置最低loglevel
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	log.WithFields(logrus.Fields{
		"filename": "124.txt",
	}).Info("打开文件失败")

	entry := logrus.WithFields(logrus.Fields{"request_id": 1, "user_ip": "192.168.1.1"})
	entry.Info("something happened on that request")
	entry.Warn("something not great happened")

	loga:=logrus.New()
	loga.WithFields(logrus.Fields{"xxxx": 1, "xxx": "192.168.1.1"}).WithFields(logrus.Fields{"request_id": 1, "user_ip": "192.168.1.2"}).Info("xx")

	// 	entry := logger.LogEntry
	// f := logger.fields
	// f["timedate"] = time.Now()
	// entry.WithFields(log.Fields(f))

	a:=make(map[string]interface{})
	a["class"]="class"
	loga.WithFields(logrus.Fields(a)).Info("yy")


}

/*
flower@:~/workspace/learngo/src/myGoNotes$ go run logrus1.go
INFO[0000] A walrus appears                              animal=walrus
*/

/*
Fields

有时候我们需要固定的fields,不需要向每行都重复写,只需要生成一 logrus.Entry


INFO[0000] something happened on that request            request_id=1 user_ip=192.168.1.1
WARN[0000] something not great happened                  request_id=1 user_ip=192.168.1.1
*/

/*
Entry

logrus.WithFields会自动返回一个 *Entry，Entry里面的有些变量会被自动加上

time:entry被创建时的时间戳
msg:在调用.Info()等方法时被添加
level
Thread safety

默认的logger在并发写的时候是被mutex保护的，比如当同时调用hook和写log时mutex就会被请求，有另外一种情况，文件是以appending mode打开的，
此时的并发操作就是安全的，可以用logger.SetNoLock()来关闭它


logrus的hook可以发给外面的邮件之类的
*/
