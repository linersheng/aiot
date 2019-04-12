package main

import (
	"github.com/linersheng/aiot/devstat"
	//"time"
)

const (
	REDIS_ADDR = "192.168.2.58:6379"
	REDIS_PWD  = "linersheng"
)

/*
var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", REDIS_ADDR)
		},
	}
}

func subscribeByConn() {
	c := pool.Get()
	defer c.Close()
	c.Do("auth", REDIS_PWD)
	psc := redis.PubSubConn{Conn: c}
	psc.PSubscribe("__keyspace@0__:*", "__keyevent@0__:*")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println("error:", v)
			return
		}
	}
}
func subscribe() {
	c := pool.Get()
	defer c.Close()
	c.Do("auth", REDIS_PWD)
	c.Send("PSUBSCRIBE", "__keyspace@0__:*")
	c.Send("PSUBSCRIBE", "__keyevent@0__:*")
	c.Flush()
	for {
		reply, err := c.Receive()
		if err != nil {
			fmt.Println("recive err:", err)
			return
		}
		vacs, ok := reply.([]interface{})
		if ok {
			for _, r := range vacs {
				if v, ok := r.([]byte); ok {
					fmt.Println("[]byte value:", string(v))
				} else if v, ok := r.(int64); ok {
					fmt.Println("int value:", v)
				} else {
					fmt.Printf("convert to string fail.(%v)\n", r)
				}
				//fmt.Println( reflect.TypeOf(r))
			}
		} else {
			fmt.Println("convert to [] interface{} fail.")
		}
	}
}
*/

func main() {
	//c := pool.Get()
	//defer c.Close()
	//defer pool.Close()
	//redis.Subscription{}
	//go subscribeByConn()
	//go subscribe()
	//c.Do("auth", REDIS_PWD)
	//c.Do("set", "lines", 123)

	// r := &Redis{
	// 	Addr:      REDIS_ADDR,
	// 	Passwd:    REDIS_PWD,
	// 	Subscribe: 0x03,
	// }

	// r.InitRedis()

	// for {
	// 	//v, _ := redis.Int(c.Do("get", "lines"))
	// 	//fmt.Println("read lines:", v)
	// 	time.Sleep(time.Second * 10)
	// }

	c := devstat.Config{
		RedisAddr:   REDIS_ADDR,
		RedisPasswd: REDIS_PWD,
		HeartTime:   60,
		Port:        4455,
	}
	devstat.SetConfig(c)
	devstat.Run()
}
