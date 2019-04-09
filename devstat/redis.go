package main

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

func (r *Redis) CheckInDB(id string) bool {
	c := r.pool.Get()
	defer c.Close()

	reply, err := c.Do("SCAN", 0, "MATCH", id, "COUNT", 1)
	if nil != err {
		return false
	}

	arr, err := redis.MultiBulk(reply, nil)
	if nil != err {
		return false
	}

	hint, err := redis.Int(arr[0], nil)
	if nil != err || hint == 0 {
		return false
	}

	return true
}

func (r *Redis) SetExpired(id string, timeout int) error {
	c := r.pool.Get()
	defer c.Close()

	if _, err := c.Do("EXPIRE", id, timeout); err != nil {
		return err
	}

	return nil
}
