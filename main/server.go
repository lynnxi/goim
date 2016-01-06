package main

import (
	"game-im/config"
	//"game-im/lib/stdlog"
	"io"
	"net"
	//"time"
)

type Server struct {
	SessionMap         map[string]*Session
	DazaRedisSlave     *RedisClient
	DazaRedisMaster    *RedisClient
	GameInfoRedisSlave *RedisClient
}

func newServer() (server *Server) {
	server = &Server{SessionMap: map[string]*Session{},
		DazaRedisSlave:     NewRedisClient(config.REDIS_DAZA_SLAVE),
		DazaRedisMaster:    NewRedisClient(config.REDIS_DAZA_MASTER),
		GameInfoRedisSlave: NewRedisClient(config.REDIS_GAMEINFO_SLAVE)}
	server.SessionMap["test"] = NewSession(NewConnection(nil))
	return
}

func (server *Server) Listen(host string) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		infoLogger.Println("accept conn...")
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
		infoLogger.Println("conn closed...")
		session.Close()
		if session.Appid != "" {
			delete(server.SessionMap, session.Appid)
		}
	}()

	var reqPack, resPack *Gpack
	var err error

	reqPack, err = session.conn.ReadGpack()
	if err != nil {
		errorLogger.Println("handle conn: read goack failed: %s", err)
		return
	}

	resPack, err = GpackAuthHandler(reqPack, server)
	if err != nil {
		errorLogger.Println("handle conn: auth failed: %s", err)
		return
	}

	session.Appid = reqPack.Appid
	server.SessionMap[reqPack.Appid] = session
	session.conn.WriteGpack(resPack)

	go session.HandleGpackBuffer()
	for {
		reqPack, err = session.conn.ReadGpack()
		if err != nil {
			if err == io.EOF {
				return
			}
			errorLogger.Println("handle conn: json decoding failed: %s", err)
			session.conn.WriteGpack(GetErrGpack())
			continue
		}

		if reqPack.Action == "ret" {
			if len(session.ResPackBufferChannel) < config.IMJ_CHANNEL_BUFF_SIZE {
				session.ResPackBufferChannel <- reqPack
			}

		} else {
			resPack, err = Action2Handler[reqPack.Action](reqPack, server)
			session.conn.WriteGpack(resPack)
		}
	}

}
