package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

const (
	KEY_PREFIX = "auth_cookie_"
)

var (
	blog_redis      *redis.Client
	blog_redis_once sync.Once
)

func createRedisClient(address, passwd string, db int) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: passwd,
		DB:       db,
	})
	if err := cli.Ping().Err(); err != nil {
		panic(fmt.Errorf("connect to redis %d failed %v", db, err))
	} else {
		fmt.Printf("connect to redis %d\n", db) //能ping成功才说明连接成功
	}
	return cli
}

func GetRedisClient() *redis.Client {
	blog_redis_once.Do(func() {
		if blog_redis == nil {
			blog_redis = createRedisClient("127.0.0.1:6379", "", 0)
		}
	})
	return blog_redis
}

// 把<cookie_value, uid>写入redis
func SetCookieAuth(cookieValue, uid string) {
	client := GetRedisClient()
	if err := client.Set(KEY_PREFIX+cookieValue, uid, time.Hour*24*30).Err(); err != nil { //30天之后过期
		fmt.Printf("write  pair(%s, %s) to redis failed: %s\n", cookieValue, uid, err)
	} else {
		// fmt.Printf("write  pair(%s, %s) to redis\n", cookieValue, uid)
	}
}

// 根据cookie_value获取uid
func GetCookieAuth(cookieValue string) (uid string) {
	client := GetRedisClient()
	var err error
	if uid, err = client.Get(KEY_PREFIX + cookieValue).Result(); err != nil {
		if err != redis.Nil {
			fmt.Printf("get auth info %s failed: %s\n", cookieValue, err)
		}
	} else {
		// fmt.Printf("get uid %s by auth key %s\n", uid, cookieValue)
	}
	return
}
