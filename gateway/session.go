package main

import (
	//"encoding/json"
	//"errors"
	. "game-im/lib/mio"
	"io"
	"runtime/debug"
	//"strconv"
	"game-im/config"
	"sync"
	"time"
)

type Session struct {
	Sessionid               string
	conn                    *Connection
	InputPackBufferChannel  chan *Gpack
	OutputPackBufferChannel chan *Gpack
	UpTime                  int64
	closed                  bool
	wait                    sync.WaitGroup
}

const (
	CHANNEL_BUFF_SIZE = 10
)

func NewSession(conn *Connection) (s *Session) {
	s = &Session{
		conn: conn,
		InputPackBufferChannel:  make(chan *Gpack, CHANNEL_BUFF_SIZE),
		OutputPackBufferChannel: make(chan *Gpack, CHANNEL_BUFF_SIZE),
		UpTime:                  time.Now().Unix(),
	}

	return
}

func (session *Session) GetConn() (conn *Connection) {
	return session.conn
}

func (session *Session) Close() {
	defer func() {
		if err := recover(); err != nil {
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))

		}

	}()

	if !session.closed {
		close(session.InputPackBufferChannel)
		close(session.OutputPackBufferChannel)
		session.closed = true
	}

	session.conn.Close()
	server.delSessionMap(session.Sessionid)
}

func (session *Session) HandleInputGpackBuffer() {
	var err error
	var gpack *Gpack
	defer func() {
		if err := recover(); err != nil {
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))
			session.conn.WriteGpack(GetErrGpack("error", gpack.Sid))
		}
		noticeLogger.Println("HandleInputGpackBuffer close...")
		session.Close()
		session.wait.Done()

	}()

	for { //读入 分发
		if session.closed {
			return
		}
		gpack, err = session.conn.ReadGpack()
		//infoLogger.Println("read gpack : " + string(gpack.ToBytes()))
		if err != nil {
			if err == io.EOF {
				noticeLogger.Println("read eof, close...")
				return
			}
			errorLogger.Println("handle conn: read gpack failed: %s", err)
			session.conn.WriteGpack(GetErrGpack("read failed", ""))
			return
		}

		gpack.SetFlag(session.Sessionid)
		session.UpTime = time.Now().Unix()

		if gpack.Cmd == "ping" {
			ret := NewGpack("pong", gpack.Sid, map[string]string{"ec": "0"})
			session.OutputPackBufferChannel <- ret
		} else {
			host := session.getLogicServer(gpack.GetFlag())
			session.sendToLogic(host, gpack)
		}
	}
}

func (session *Session) HandleOuputGpackBuffer() {
	var err error
	var gpack *Gpack
	defer func() {
		if err := recover(); err != nil {
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))
			session.conn.WriteGpack(GetErrGpack("error", gpack.Sid))
		}
		noticeLogger.Println("HandleOuputGpackBuffer close...")
		session.Close()
		session.wait.Done()

	}()
	for {
		if session.closed {
			return
		}
		gpack = <-session.OutputPackBufferChannel
		if gpack == nil {
			return
		}
		err = session.conn.WriteGpack(gpack)

		//infoLogger.Println("write gpack : " + string(gpack.ToBytes()))
		if err != nil {
			errorLogger.Println("write gpack error : ", err)
			session.Close()
		}
	}
}

func (session *Session) HandleMsgGpackBuffer() {

}

func (session *Session) getLogicServer(gameid string) (host string) {
	host = config.REDIS_LOGIC_1
	return
}

func (session *Session) sendToLogic(host string, gpqak *Gpack) {
	redis, ok := logicQueueRedis[host]
	if !ok {
		redis = NewRedisClient(host)
		logicQueueRedis[host] = redis
	}
	redis.Rpush(config.KEY_LOGIC_IN, string(gpqak.ToBytes()))
}
