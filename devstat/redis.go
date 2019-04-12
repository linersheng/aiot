package devstat

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

const (
	SUBSCRIBE_KEYSPACE = 0x01
	SUBSCRIBE_KEYEVENT = 0x02
)

type Redis struct {
	Addr      string
	Passwd    string
	Subscribe int
	pool      *redis.Pool
}

func (r *Redis) runSubscribe() {
	c := r.pool.Get()
	defer c.Close()

	if (r.Subscribe & SUBSCRIBE_KEYSPACE) != 0 {
		c.Send("PSUBSCRIBE", "__keyspace@0__:*")
	}

	if (r.Subscribe & SUBSCRIBE_KEYEVENT) != 0 {
		c.Send("PSUBSCRIBE", "__keyevent@0__:*")
	}

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

func (r *Redis) InitRedis() {
	r.pool = &redis.Pool{
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.Addr)
			if err != nil {
				return nil, err
			}
			if r.Passwd != "" {
				if _, err := c.Do("AUTH", r.Passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}

	go r.runSubscribe()
}

func (r *Redis) TryGetKeyInDB(id string) (string) {
	c := r.pool.Get()
	defer c.Close()

	arr, err := redis.Values(c.Do("SCAN", 0, "MATCH", id+"*", "COUNT", 1))
	if nil != err {
		fmt.Println("err:", err)
		return ""
	}
	fmt.Printf("%v\n", arr)

	iter, err := redis.Int(arr[0], nil)
	if nil != err || iter == 0{
		fmt.Println("err:", err)
		return ""
	}

	fmt.Println("cursor",iter)

	vals, err := redis.Strings(arr[1], nil)
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}
	fmt.Println("vals:", vals)
	return vals[0]
}

func (r *Redis) SetExpired(key string, timeout int) error {
	c := r.pool.Get()
	defer c.Close()
	fmt.Println("SetExpired(), key:", key)
	if _, err := c.Do("EXPIRE", key, timeout); err != nil {
		return err
	}
	return nil
}

func (r *Redis) AddKeyWithTimeout(key string, timeout int) error {
	c := r.pool.Get()
	defer c.Close()

	fmt.Println("AddKeyWithTimeout(), key:", key)
	if _, err := c.Do("SET", key, 0, "EX", timeout); err != nil {
		return err
	}
	return nil
}
