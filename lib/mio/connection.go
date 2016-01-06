package mio

import (
	"bufio"
	"encoding/json"
	"game-im/config"
	//"game-im/lib/vmihailenco/msgpack"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Connection struct {
	net.Conn
	reader        *Reader
	writer        *bufio.Writer
	stats         map[string]int
	lastWriteTime int64
	lastReadTime  int64
}

func NewConnection(conn net.Conn) (c *Connection) {
	//这里net.conn为什么不是引用？
	c = &Connection{
		Conn: conn,
	}
	c.reader = NewReader(c.Conn)
	c.writer = bufio.NewWriter(c.Conn)
	c.stats = make(map[string]int)

	return
}

func (c *Connection) GetStats() map[string]int {
	return c.stats
}

func (c *Connection) flush() {
	c.writer.Flush()
}

func (c *Connection) ReadGpack() (gpack *Gpack, err error) {
	var line []byte
	if line, err = c.reader.ReadLine(); err != nil {
		return
	}

	sz, err := strconv.Atoi(string(line))
	if err != nil {
		return
	}
	if sz < 0 {
		return
	}
	var buf = make([]byte, sz)
	var p = buf
	for {
		var n int
		n, err = c.reader.Read(p)
		if err != nil {
			return
		}
		if n < len(p) {
			p = p[n:]
		} else {
			break
		}
	}
	gpack = new(Gpack)
	err = json.Unmarshal(buf, gpack)

	fmt.Println("read gpack : " + string(gpack.ToBytes()))

	now := time.Now().Unix()
	if now == c.lastReadTime {
		c.stats["read_qpack"]++
	} else {
		c.stats["read_qpack"] = 1
		c.lastReadTime = now
	}

	return
}

func (c *Connection) WriteGpack(gpack *Gpack) (err error) {
	data, err := json.Marshal(gpack)
	_, err = c.writer.Write([]byte(strconv.Itoa(len(data))))
	_, err = c.writer.Write([]byte(config.CRLF))
	_, err = c.writer.Write(data)

	err = c.writer.Flush()
	if err != nil {
		return
	}

	fmt.Println("write gpack : " + string(gpack.ToBytes()))

	now := time.Now().Unix()
	if now == c.lastWriteTime {
		c.stats["write_qpack"]++
	} else {
		c.stats["write_qpack"] = 1
		c.lastWriteTime = now
	}
	return
}
