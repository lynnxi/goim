package main

import (
	"fmt"
	"game-im/lib/stdlog"
	"os"
	"runtime"
	"time"
)

var (
	infoLogger   = stdlog.Log("info")
	noticeLogger = stdlog.Log("notice")
	errorLogger  = stdlog.Log("error")
	panicLogger  = stdlog.Log("panic")
)

func main() {

	outLogInfo, _ := os.OpenFile("/home/deploy/log/worker_info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	infoLogger.SetOutput(outLogInfo)

	outLogNotice, _ := os.OpenFile("/home/deploy/log/worker_notice.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	noticeLogger.SetOutput(outLogNotice)

	outLogError, _ := os.OpenFile("/home/deploy/log/worker_error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	errorLogger.SetOutput(outLogError)

	outLogPanic, _ := os.OpenFile("/home/deploy/log/worker_panic.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	panicLogger.SetOutput(outLogPanic)

	// （可选）设置函数前缀
	stdlog.SetPrefix(func() string {
		t := time.Now()
		return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d] ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	})

	runtime.GOMAXPROCS(4)

	infoLogger.Println("start...")

	go PopGmsg()      //pop 群消息
	go LeaveGroup()   //退群
	go DisbandGroup() //解散群

	for {
		time.Sleep(time.Second)
	}

	infoLogger.Println("stop...")
}
