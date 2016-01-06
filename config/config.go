package config

import ()

const (
	CR           = '\r'
	LF           = '\n'
	CRLF         = "\r\n"
	KEY_IMQ      = "_GAME_IMJ_"
	KEY_IMQ_TEST = "_GAME_IMJ_TEST_"
	KEY_IMQ_V2   = "_GAME_IMJ_V2_"
	SALT         = 1602

	KEY_LOGIC_IN  = "_GAME_LOGIC_IN_"
	KEY_LOGIC_OUT = "_GAME_LOGIC_OUT_"
	KEY_GW_IN     = "_GAME_GW_IN_"
	KEY_GW_OUT    = "_GAME_GW_OUT_"

	PATH_LOG_GW    = "/home/deploy/log"
	PATH_LOG_LOGIC = "/home/deploy/log"
)

const (
	CHANNEL_BUFF_SIZE = 10
)

var (
	REDIS_LOCAL_QUEUE = "0.0.0.0:1603"
	ImqRedisConfig    = [...]string{REDIS_IMQ_1, REDIS_IMQ_2}
)
