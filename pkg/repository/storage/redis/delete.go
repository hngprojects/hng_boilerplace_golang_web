package redis

import "github.com/go-redis/redis/v8"

func RedisDelete(rdb *redis.Client, key string) (int64, error) {
	deleted, err := rdb.Del(Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
