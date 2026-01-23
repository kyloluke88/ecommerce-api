package bootstrap

import (
	"api/pkg/config"
	"api/pkg/redis"
	"fmt"
	"strconv"
)

// SetupRedis 初始化 Redis
func SetupRedis() {

	// todo 这里为什么Get[int]获取不到int
	dbNum, _ := strconv.Atoi(config.Get[string]("redis.database"))

	// 建立 Redis 连接
	redis.ConnectRedis(
		fmt.Sprintf("%v:%v", config.Get[string]("redis.host"), config.Get[string]("redis.port")),
		config.Get[string]("redis.username"),
		config.Get[string]("redis.password"),
		dbNum,
	)
}
