package utils

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func GetRedisClusterDetails(addr, password string) map[string]string {

	redisClusterOptions := redis.ClusterOptions{
		Addrs: []string{fmt.Sprintf("%s:6379", addr)},
	}

	if password != "" {
		redisClusterOptions.Password = password
	}

	currentShards := map[string]string{}
	ctx := context.Background()
	rdb := redis.NewClusterClient(&redisClusterOptions)
	shards, _ := rdb.ClusterShards(ctx).Result()
	for _, shard := range shards {
		for _, node := range shard.Nodes {
			currentShards[node.ID] = node.IP
		}
	}
	return currentShards
}
