package main

import (
	"github.com/go-redis/redis/v8"
	"strconv"
)

func IncreaseVisitorCount(rdb *redis.Client) (int, error) {
	valStr, err := RedisGet(rdb, "visitor-count", "0")
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, err
	}

	err = RedisSet(rdb, "visitor-count", strconv.Itoa(int(val+1)))
	if err != nil {
		return int(val), err
	}

	return int(val), nil
}

func RedisGet(rdb *redis.Client, key string, defaultValue string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return defaultValue, nil
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func RedisSet(rdb *redis.Client, key string, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}
