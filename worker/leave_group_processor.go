package main

import (
	"encoding/json"
	"game-im/config"
	. "game-im/lib/mio"
	//"game-im/lib/stdlog"
	"game-im/lib/mcrypt"
	"runtime/debug"
	"strconv"
)

func LeaveGroup() {
	defer func() {
		errorLogger.Println("leave group shut down...")
	}()

	for _, host := range config.ImqRedisConfig {
		wait.Add(1)
		go leaveGroup(host)
	}

	wait.Wait()
}

func leaveGroup(host string) {
	
}
