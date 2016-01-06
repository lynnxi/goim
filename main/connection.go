package main

import (
	"bufio"
	"game-im/lib/mio"
	//"game-im/lib/stdlog"
	"encoding/json"
	"game-im/config"
	"github.com/vmihailenco/msgpack"
	"net"
	"time"
)

type Connection struct {
	net.Conn
	reader        *mio.Reader
	writer        *bufio.Writer
	stats         map[string]int
	lastStatsTime int64
}

func NewConnection(conn net.Conn) (c *Connection) {
	//这里net.conn为什么不是引用？
	c = &Connection{
		Conn: conn,
	}
	c.reader = mio.NewReader(c.Conn)
	c.writer = bufio.NewWriter(c.Conn)
	c.stats = make(map[string]int)

	return
}

func (c *Connection) ReadGpack() (gpack *Gpack, err error) {
	var line []byte
	if line, err = c.reader.ReadLine(); err != nil {
		return
	}

	gpack = new(Gpack)
	msgpack.Unmarshal(line, gpack)

	now := time.Now().Unix()
	if now == c.lastStatsTime {
		c.stats["read_qpack"]++
	} else {
		c.stats["read_qpack"] = 0
		c.lastStatsTime = now
	}

	return
}

func (c *Connection) WriteGpack(gpack *Gpack) (err error) {
	data, err = msgpack.Marshal(gpack)
	_, err = c.writer.Write(data)
	_, err = c.writer.Write([]byte(config.CRLF))
	c.writer.Flush()

	now := time.Now().Unix()
	if now == c.lastStatsTime {
		c.stats["write_qpack"]++
	} else {
		c.stats["write_qpack"] = 0
		c.lastStatsTime = now
	}
	return
}
