package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func RedisGet(rdb *redis.Client, key string) ([]byte, error) {
	serialized, err := rdb.Get(Ctx, key).Bytes()
	return serialized, err
}

func PopFromQueue(rdb *redis.Client) (interface{}, error) {
	var response interface{}

	jsonValue, err := rdb.RPop(Ctx, KeyName).Result()
	if err != nil {
		return response, fmt.Errorf("could not pop from Redis queue: %v", err)
	}

	err = json.Unmarshal([]byte(jsonValue), &response)
	if err != nil {
		return response, fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	return response, nil
}
