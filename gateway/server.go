package main

import (
	"encoding/json"
	"game-im/config"
	//"game-im/lib/stdlog"
	"errors"
	. "game-im/lib/mio"
	"net"
	"runtime/debug"
	"time"
)

type Server struct {
	host         string
	port         string
	local_host   string
	local_port   string
	gwQueueRedis *RedisClient
	SessionMap   map[string]*Session
}

func newServer(host string, port string, local_host string, local_port string) (server *Server) {
	server = &Server{
		host:         host,
		port:         port,
		local_host:   local_host,
		local_port:   local_port,
		gwQueueRedis: NewRedisClient(local_host + ":" + local_port),
		SessionMap:   map[string]*Session{},
	}
	return
}

func (server *Server) Start() {
	go server.checkSessionMap()
	go server.PopGpack()
	server.Listen()

}

/**
 *
 */
func (server *Server) Listen() error {
	listener, err := net.Listen("tcp", server.host+":"+server.port)

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()

		infoLogger.Println("accept conn..." + conn.RemoteAddr().String() + conn.RemoteAddr().Network())
		if err != nil {
			return err
		}
		session := NewSession(NewConnection(conn))
		go server.handleConnection(session)
	}

	return nil
}

func (server *Server) handleConnection(session *Session) {
	defer func() {
		if err := recover(); err != nil { //子协程崩溃貌似catch不住
			panicLogger.Println("panic : ", err)
			panicLogger.Println("stack : ", string(debug.Stack()))

			session.conn.WriteGpack(GetErrGpack("error", ""))
		}

		infoLogger.Println("conn closed...")
		session.Close() //这里去redis删除了session， 但服务器重启或被kill后没能处理session

	}()

	//身份认证
	reqPack, err := session.conn.ReadGpack()
	if err != nil {
		errorLogger.Println("handle conn: read gpack failed: %s", err)
		return
	}

	sessionid, err := GpackAuthHandler(reqPack)
	if err != nil {
		errorLogger.Println("handle conn: auth failed: %s", err)
		resPack := GetErrGpack("auth failed", reqPack.Sid)
		session.conn.WriteGpack(resPack)
		return
	}
	noticeLogger.Println("handle conn: auth success: %s", reqPack.GetFlag())
	resPack := NewRetGpack(reqPack.Sid, map[string]string{"ec": "0"})
	session.conn.WriteGpack(resPack)

	/**
	* @TODO 踢下线
	**/
	server.setSessionMap(sessionid, session)

	// //取群lv
	// key := "gmsg:1:version"
	// v, err := msgRedis.Get(key)

	// if err != nil {
	// 	errorLogger.Println("msg version get err ", err)
	// 	return
	// }
	// var lv = 0
	// if v == nil {
	// 	errorLogger.Println("msg version get nil ", key)
	// } else {
	// 	lv, _ = strconv.Atoi(string(v.([]byte)))
	// }

	// session.conn.WriteGpack(NewReqGpack(reqPack.Gameid, reqPack.Appid, "gmsg_psh", map[string]interface{}{"gid": "1", "lv": lv}))

	session.wait.Add(2)
	go session.HandleInputGpackBuffer() //处理输入
	go session.HandleOuputGpackBuffer() //处理输出

	session.wait.Wait()
}

func (server *Server) setSessionMap(sessionid string, session *Session) {
	session.Sessionid = sessionid
	server.SessionMap[sessionid] = session

	key := "session_map"
	_server, _ := json.Marshal(map[string]interface{}{"host": server.host,
		"port":       server.port,
		"local_host": server.local_host,
		"local_port": server.local_port})
	SessionRedis.Hset(key, sessionid, string(_server))

}

func (server *Server) delSessionMap(sessionid string) {
	delete(server.SessionMap, sessionid)
	key := "session_map"
	SessionRedis.Hdel(key, sessionid)
}

func (server *Server) getSessionMap() (m map[string]*Session) {
	return server.SessionMap
}

func (server *Server) checkSessionMap() {
	for {
		sessionMap := server.getSessionMap()

		_session := make(map[string]int64)
		_stats := make(map[string]interface{})
		now := time.Now().Unix()
		for sessionid, session := range sessionMap {
			if now-session.UpTime > 5*60 { //超时检测
				session.Close()
			}
			_session[sessionid] = session.UpTime
			_stats[sessionid] = session.GetConn().GetStats()
		}

		out, _ := json.Marshal(map[string]interface{}{"session": _session, "stats": _stats})

		infoLogger.Println("session map : ", string(out))
		time.Sleep(time.Second)
	}
}

func (server *Server) PopGpack() {
	var gpack *Gpack
	for {
		v, err := server.gwQueueRedis.Blpop(config.KEY_GW_IN, 0)
		if err != nil {
			errorLogger.Println("BLPOP ERR : ", err)
			continue
		}
		gpack = nil
		if err = json.Unmarshal(v.([]byte), &gpack); err != nil {
			errorLogger.Println("RedisPool: json decoding failed: %s", err)
			errorLogger.Println(string(v.([]byte)))
			continue
		}
		/**
		 * @todo 根据包找到session 发送到ouput channel
		 **/
		session, ok := server.SessionMap[gpack.GetFlag()]
		if ok {
			defer func() {
				if err := recover(); err != nil { //捕获此时channel被关闭的异常
					panicLogger.Println("panic : ", err)
					panicLogger.Println("stack : ", string(debug.Stack()))

					infoLogger.Println("conn closed...")
					session.Close() //这里去redis删除了session， 但服务器重启或被kill后没能处理session
				}

			}()
			session.OutputPackBufferChannel <- gpack
		} else {
			errorLogger.Println("gpack not found session %s", gpack.GetFlag())
		}

	}
}

func GpackAuthHandler(req *Gpack) (sessionid string, err error) {
	if req.Cmd != "auth" {
		err = errors.New("error")
		return
	}
	appid := req.Body.(map[string]interface{})["appid"].(string)
	token := req.Body.(map[string]interface{})["token"].(string)

	key := "app:" + appid + ":baseinfo"
	v, err := GameInfoRedisSlave.Get(key)
	if err != nil {
		errorLogger.Println("app base info get error ", err)
		return
	}
	if v == nil {
		errorLogger.Println("app base info get nil ", key)
		err = errors.New("auth: pp base info get nil " + key)
		return
	}

	var gameInfo map[string]interface{}
	if err = json.Unmarshal(v.([]byte), &gameInfo); err != nil {
		errorLogger.Println("game info: data json decoding failed: %s", err)
	}

	if err != nil {
		errorLogger.Println("req data get error ", err)
		return
	}

	if token == gameInfo["app_secret"].(string) {
		sessionid = appid
	} else {
		errorLogger.Println("auth: token invalid ")
		err = errors.New("auth: token invalid ")
	}

	return
}
