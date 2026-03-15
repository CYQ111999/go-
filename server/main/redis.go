package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// 定义一个全局的pool
var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) (*redis.Pool, error) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
	return pool, nil
}
