package main

import (
	"fmt"
	"game-im/lib/stdlog"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

var (
	infoLogger   = stdlog.Log("info")
	noticeLogger = stdlog.Log("notice")
	errorLogger  = stdlog.Log("error")
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("211.152.99.46:6060", nil))
	}()

	outLogInfo, _ := os.OpenFile("/home/deploy/log/info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	infoLogger.SetOutput(outLogInfo)

	outLogNotice, _ := os.OpenFile("/home/deploy/log/notice.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	noticeLogger.SetOutput(outLogNotice)

	outLogError, _ := os.OpenFile("/home/deploy/log/error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	errorLogger.SetOutput(outLogError)

	// （可选）设置函数前缀
	stdlog.SetPrefix(func() string {
		t := time.Now()
		return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d] ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	})

	runtime.GOMAXPROCS(4)

	infoLogger.Println("start...")
	server := newServer()
	go server.PopImj()
	go server.PopImj()
	go server.PopImj()
	go server.PopImj()
	go server.PopImj()
	go server.PopImj()
	go server.PopImj()
	go server.PopImj()
	server.Listen("211.152.99.46:6320")
	infoLogger.Println("stop...")
}
