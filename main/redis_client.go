package main

import (
	"game-im/lib/redigo/redis"
	//"game-im/lib/stdlog"
	"time"
)

type RedisClient struct {
	//conn *Connection
	host        string
	pool        *redis.Pool
	maxPoolSize int

	poolIn  chan redis.Conn
	poolOut chan redis.Conn
}

func NewRedisClient(config string) (redisClient *RedisClient) {
	redisClient = &RedisClient{host: config}
	// c, err := net.DialTimeout("tcp", config, time.Millisecond*1000)
	// if err != nil {
	// 	return
	// }
	// proxy.conn = NewConnection(c)
	redisClient.maxPoolSize = 2000
	redisClient.poolOut = make(chan redis.Conn, redisClient.maxPoolSize)
	redisClient.poolIn = make(chan redis.Conn, redisClient.maxPoolSize)

	redisClient.pool = &redis.Pool{
		MaxIdle:     2000,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", ""); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, err
		},
		// TestOnBorrow: func(c redis.Conn, t time.Time) error {
		// 	_, err := c.Do("PING")
		// 	return err
		// },
	}
	// go redisClient.producer()
	// go redisClient.consumer()
	return
}

func (redisClient *RedisClient) consumer() {
	var conn redis.Conn
	for {
		conn = <-redisClient.poolIn
		conn.Close()
	}
}

func (redisClient *RedisClient) producer() {
	var conn redis.Conn
	for {
		conn = redisClient.pool.Get()
		redisClient.poolOut <- conn
	}
}

func (redisClient *RedisClient) GetConn() redis.Conn {
	return redisClient.pool.Get()
	//return <-redisClient.poolOut
}

func (redisClient *RedisClient) CloseConn(conn redis.Conn) {
	conn.Close()
	//redisClient.poolIn <- conn
}

func (redisClient *RedisClient) Get(key string) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	// conn.Do("set", "b", "2")
	// _r, err := conn.Do("get", "b")
	// fmt.Println(string(_r.([]byte)[0]))
	//stdlog.Println("proxy cmd " + proxy.host + " : " + cmd.String())
	value, err = conn.Do("GET", key)

	return
}

func (redisClient *RedisClient) Set(key string, value string) (err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	_, err = conn.Do("SET", key, value)

	return
}

func (redisClient *RedisClient) Lpop(key string) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	value, err = conn.Do("LPOP", key)

	return
}

func (redisClient *RedisClient) Rpush(key string, value string) (err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	_, err = conn.Do("RPUSH", key, value)

	return
}
