package main

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/myapp/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
}


//for signal
l := &lumberjack.Logger{}
log.SetOutput(l)
c := make(chan os.Signal, 1)
signal.Notify(c, syscall.SIGHUP)

go func() {
    for {
        <-c
        l.Rotate()
    }
}()

//这个包用于log的组件，不是单独出来用的