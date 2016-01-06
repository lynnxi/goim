package mio

import (
	"fmt"
	"game-im/lib/redigo/redis"
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

func (redisClient *RedisClient) SetBit(key string, offset int, value int) (err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	_, err = conn.Do("SETBIT", key, offset, value)

	return
}

func (redisClient *RedisClient) GetBit(key string, offset int) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	value, err = conn.Do("GETBIT", key, offset)

	return
}

func (redisClient *RedisClient) Del(key string) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	value, err = conn.Do("DEL", key)

	return
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

func (redisClient *RedisClient) SetEx(key string, value string, ttl int) (err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	_, err = conn.Do("SETEX", key, ttl, value)

	return
}

func (redisClient *RedisClient) Test(key string) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)

	value, err = conn.Do("HGETALL", key)

	return
}

func (redisClient *RedisClient) HgetAll(key string) (value map[string]interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)

	v, err := conn.Do("HGETALL", key)

	if err != nil {
		return
	}

	if v == nil || len(v.([]interface{})) == 0 {
		return
	}

	var k string
	value = map[string]interface{}{}
	for i, b := range v.([]interface{}) {
		if i%2 == 0 {
			k = string(b.([]byte))
		} else {
			value[k] = string(b.([]byte))
		}
	}

	return
}

func (redisClient *RedisClient) Hget(key string, field string) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)

	value, err = conn.Do("HGET", key, field)

	return
}

func (redisClient *RedisClient) Hset(key string, field string, value string) (ret interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)

	ret, err = conn.Do("HSET", key, field, value)

	return
}

func (redisClient *RedisClient) Hdel(key string, field string) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)

	value, err = conn.Do("HDEL", key, field)

	return
}

func (redisClient *RedisClient) Lpop(key string) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	value, err = conn.Do("LPOP", key)

	return
}

func (redisClient *RedisClient) Blpop(key string, timeout int) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	value, err = conn.Do("BLPOP", key, timeout)

	if err == nil && value != nil {

		value = value.([]interface{})[1] //blop多条批量回复 第一条是key  第二条是value，我们只用一个key，所以只要第二条

	}

	return
}

func (redisClient *RedisClient) Brpop(key string, timeout int) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	value, err = conn.Do("BRPOP", key, timeout)

	if err == nil && value != nil {
		value = value.([]interface{})[1]
	}
	return
}

func (redisClient *RedisClient) Rpush(key string, value string) (err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	_, err = conn.Do("RPUSH", key, value)

	return
}

func (redisClient *RedisClient) Lpush(key string, value string) (err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	_, err = conn.Do("LPUSH", key, value)

	return
}

func (redisClient *RedisClient) Zrange(key string, start int, end int, withscores bool) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	if withscores {
		value, err = conn.Do("ZRANGE", key, start, end, "WITHSCORES")
	} else {
		value, err = conn.Do("ZRANGE", key, start, end)
	}

	return
}

func (redisClient *RedisClient) ZrangeByScore(key string, min int, max int, withscores bool) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)

	var m interface{}
	if max == -1 {
		m = "+inf"
	} else {
		m = max
	}

	if withscores {
		value, err = conn.Do("ZRANGEBYSCORE", key, min, m, "WITHSCORES")
	} else {
		value, err = conn.Do("ZRANGEBYSCORE", key, min, m)
	}

	return
}

func (redisClient *RedisClient) Incr(key string) (value interface{}, err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	value, err = conn.Do("INCR", key)

	return
}

func (redisClient *RedisClient) Zadd(key string, member string, score int) (err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	_, err = conn.Do("ZADD", key, score, member)

	return
}

func (redisClient *RedisClient) ZremRangeByScore(key string, min int, max int) (err error) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)
	_, err = conn.Do("ZREMRANGEBYSCORE", key, min, max)

	return
}

func (redisClient *RedisClient) Monitor(f func(key string, event string)) {
	conn := redisClient.GetConn()
	defer redisClient.CloseConn(conn)

	config_key := "notify-keyspace-events"
	notification_config := "xE"
	event := "__keyspace@0__:expire_callback"
	conn.Send("CONFIG", "SET", config_key, notification_config)
	conn.Send("SUBSCRIBE", event)
	conn.Flush()
	for {
		reply, err := conn.Receive()
		if err != nil {
			fmt.Println("monitor subscribe err...", err)
			return
		}
		// process pushed message

		if data, ok := reply.([]string); ok {
			if len(data) != 3 || data[0] != "message" {
			} else {
				key := data[1]
				event := data[2]
				f(key, event)

			}
		} else {
		}
	}
}
