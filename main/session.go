package main

import (
	"game-im/config"
)

type Session struct {
	conn                 *Connection
	ReqPackBufferChannel chan *Gpack
	ResPackBufferChannel chan *Gpack
	Appid                string
}

func NewSession(conn *Connection) (s *Session) {
	s = &Session{
		conn:                 conn,
		ReqPackBufferChannel: make(chan *Gpack, config.IMJ_CHANNEL_BUFF_SIZE),
		ResPackBufferChannel: make(chan *Gpack, config.IMJ_CHANNEL_BUFF_SIZE),
	}

	return
}

func (session *Session) Close() {
	session.conn.Close()
}

func (session *Session) HandleGpackBuffer() {
	var reqPack, resPack *Gpack
	for {
		reqPack = <-session.ReqPackBufferChannel
		session.conn.WriteGpack(reqPack)
		resPack = <-session.ResPackBufferChannel
		if resPack.Appid != "nil" {
		}
	}
}
