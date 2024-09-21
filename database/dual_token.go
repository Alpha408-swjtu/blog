package database

import (
	"blog/util"
	"time"

	"github.com/go-redis/redis"
)

/**
什么时候使用Redis?
1. 高并发，低延时。redis比mysql快一到两个数量级。
2. redis可靠性没mysql高，万一redis挂了对业务影响不大，好修复。
3. redis通常存储string型value，此时它相对于mysql的性能优势更明显。
*/

const (
	TOKEN_PREFIX = "dual_token_"
	TOKEN_EXPIRE = 7 * 24 * time.Hour //一次登录7天有效
)

// 把<refreshToken, authToken>写入redis
func SetToken(refreshToken, authToken string) {
	client := GetRedisClient()
	if err := client.Set(TOKEN_PREFIX+refreshToken, authToken, TOKEN_EXPIRE).Err(); err != nil { //7天之后就拿不到authToken了
		util.LogRus.Errorf("write token pair(%s, %s) to redis failed: %s", refreshToken, authToken, err)
	}
}

// 根据refreshToken获取authToken
func GetToken(refreshToken string) (authToken string) {
	client := GetRedisClient()
	var err error
	if authToken, err = client.Get(TOKEN_PREFIX + refreshToken).Result(); err != nil {
		if err != redis.Nil {
			util.LogRus.Errorf("get auth token of refresh token %s failed: %s", refreshToken, err)
		}
	}
	return
}
