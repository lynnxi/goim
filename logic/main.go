package main

import (
	"flag"
	"fmt"
	"game-im/config"
	. "game-im/lib/mio"
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

var (
	gwQueueRedis = map[string]*RedisClient{
		config.REDIS_GW_1: NewRedisClient(config.REDIS_GW_1),
	}

	logic *Logic
	Moa   *MoaClient = NewMoaClient()
)

func main() {

	outLogInfo, _ := os.OpenFile(config.PATH_LOG_LOGIC+"/logic_info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	infoLogger.SetOutput(outLogInfo)

	outLogNotice, _ := os.OpenFile(config.PATH_LOG_LOGIC+"/logic_notice.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	noticeLogger.SetOutput(outLogNotice)

	outLogError, _ := os.OpenFile(config.PATH_LOG_LOGIC+"/logic_error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	errorLogger.SetOutput(outLogError)

	outLogPanic, _ := os.OpenFile(config.PATH_LOG_LOGIC+"/logic_panic.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	panicLogger.SetOutput(outLogPanic)

	// （可选）设置函数前缀
	stdlog.SetPrefix(func() string {
		t := time.Now()
		return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d] ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	})

	runtime.GOMAXPROCS(4)

	local_host := flag.String("lh", "127.0.0.1", "server host")
	local_port := flag.String("lp", "2603", "server port")

	flag.Parse()

	infoLogger.Println("start...", *local_host, " ", *local_port)
	logic = newLogic(*local_host, *local_port)
	logic.Start()

	infoLogger.Println("stop...")
}
