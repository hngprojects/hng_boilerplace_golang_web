package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func RedisSet(rdb *redis.Client, key string, value interface{}) error {
	serialized, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(Ctx, key, serialized, 24*time.Hour).Err()
}

func PushToQueue(rdb *redis.Client, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("could not marshal struct: %v", err)
	}

	err = rdb.LPush(Ctx, KeyName, jsonValue).Err()
	if err != nil {
		fmt.Println("could not push to Redis queue: ", err)
	}

	return nil
}
