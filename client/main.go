package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"
	//"game-im/config"
	"game-im/lib/stdlog"
	//"io/ioutil"
	//"net/http"
	//"game-im/lib/vmihailenco/msgpack"
	//"game-im/lib/mcrypt"
	. "game-im/lib/mio"
	"io"
	"net"
	"runtime"
	"time"
)

var (
	ReqPackBufferChannel   chan *Gpack = make(chan *Gpack)
	RetPackBufferChannel   chan *Gpack = make(chan *Gpack)
	OuputPackBufferChannel chan *Gpack = make(chan *Gpack)
)
var Action2Handler = map[string]func(req *Gpack) (res *Gpack, err error){
	"msg":      GpackMsgHandler,
	"msg_psh":  GpackMsgPshHandler,
	"msg_ack":  GpackMsgAckHandler,
	"gmsg_psh": GpackGmsgPshHandler,
	"gmsg_ack": GpackGmsgAckHandler,
}

var (
	gameid string
	appid  string
)

func main() {
	
}
