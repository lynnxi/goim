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
	infoLogger      = stdlog.Log("info")
	noticeLogger    = stdlog.Log("notice")
	errorLogger     = stdlog.Log("error")
	panicLogger     = stdlog.Log("panic")
	logicQueueRedis = map[string]*RedisClient{
		config.REDIS_LOGIC_1: NewRedisClient(config.REDIS_LOGIC_1),
	}

	server *Server
)

func main() {

	outLogInfo, _ := os.OpenFile(config.PATH_LOG_GW+"/server_info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	infoLogger.SetOutput(outLogInfo)

	outLogNotice, _ := os.OpenFile(config.PATH_LOG_GW+"/server_notice.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	noticeLogger.SetOutput(outLogNotice)

	outLogError, _ := os.OpenFile(config.PATH_LOG_GW+"/server_error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	errorLogger.SetOutput(outLogError)

	outLogPanic, _ := os.OpenFile(config.PATH_LOG_GW+"/server_panic.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	panicLogger.SetOutput(outLogPanic)

	// （可选）设置函数前缀
	stdlog.SetPrefix(func() string {
		t := time.Now()
		return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d] ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	})

	runtime.GOMAXPROCS(4)

	host := flag.String("h", "0.0.0.0", "server host")
	port := flag.String("p", "2602", "server port")
	local_host := flag.String("lh", "127.0.0.1", "server host")
	local_port := flag.String("lp", "1603", "server port")

	flag.Parse()

	infoLogger.Println("start tcp server..." + *host + " " + *port + " " + *local_host + " " + *local_port)
	server = newServer(*host, *port, *local_host, *local_port)
	server.Start()
}
