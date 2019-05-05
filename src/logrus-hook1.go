import (
    "github.com/lestrrat/go-file-rotatelogs"
    "github.com/rifflock/lfshook"
    log "github.com/sirupsen/logrus"
    "time"
    "os"
    "github.com/pkg/errors"
    "path"
    "time"
)
// config logrus log to local filesystem, with file rotation
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
    baseLogPaht := path.Join(logPath, logFileName)
    writer, err := rotatelogs.New(
        baseLogPaht+".%Y%m%d%H%M",
        rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
        rotatelogs.WithMaxAge(maxAge), // 文件最大保存时间
        rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
    )
    if err != nil {
        log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
    }
    lfHook := lfshook.NewHook(lfshook.WriterMap{
        log.DebugLevel: writer, // 为不同级别设置不同的输出目的
        log.InfoLevel:  writer,
        log.WarnLevel:  writer,
        log.ErrorLevel: writer,
        log.FatalLevel: writer,
        log.PanicLevel: writer,
    })
    log.AddHook(lfHook)
}

func main(){
    d, _ := time.ParseDuration("-8h")
    ConfigLocalFilesystemLogger("./","test.log",d)
}