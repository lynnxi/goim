package main

import (
	"flag"
	"fmt"
	. "game-im/lib/mio"
	"game-im/lib/stdlog"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

var (
	infoLogger              = stdlog.Log("info")
	LogRedis   *RedisClient = NewRedisClient("")
)

func main() {

	outLogInfo, _ := os.OpenFile("/home/deploy/log/web_info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	infoLogger.SetOutput(outLogInfo)

	// （可选）设置函数前缀
	stdlog.SetPrefix(func() string {
		t := time.Now()
		return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d] ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	})

	host := flag.String("h", "0.0.0.0", "server host")
	port := flag.String("p", "8080", "server port")

	flag.Parse()

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		uid := r.FormValue("uid")

		key_queue := "uid:" + uid + ":log"
		key_session := "uid:" + uid + ":session"
		LogRedis.SetEx(key_session, "1", 3600*24)
		defer func() {
			if err := recover(); err != nil {
				infoLogger.Println("panic : ", err)
				infoLogger.Println("stack : ", string(debug.Stack()))
			}
			LogRedis.Del(key_session)
		}()

		for {
			v, err := LogRedis.Blpop(key_queue, 10)
			if err != nil {
				infoLogger.Println("BLPOP ERR : ", err)
				continue
			}
			if v == nil {
				continue
			}
			out := string(v.([]byte))
			infoLogger.Println(out)
			fmt.Fprintf(w, out)
			if f, ok := w.(http.Flusher); ok {
				infoLogger.Println("flush")
				f.Flush()
			} else {
				infoLogger.Println("Damn, no flush")
			}
		}

	})
	infoLogger.Println(http.ListenAndServe(*host+":"+*port, nil))
}
