package app

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lukasbischof/luk4s.dev/app/forum"
	"strconv"
)

func IncreaseVisitorCount(rdb *redis.Client, ctx context.Context) (int, error) {
	valStr, err := RedisGet(rdb, ctx, "visitor-count", "0")
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, err
	}

	if err := RedisSet(rdb, ctx, "visitor-count", strconv.Itoa(int(val+1))); err != nil {
		return int(val), err
	}

	return int(val), nil
}

func SaveForumEntry(rdb *redis.Client, ctx context.Context, forumEntry *forum.Entry) error {
	json, err := forumEntry.ToJson()
	if err != nil {
		return err
	}

	if err := rdb.SAdd(ctx, "forum", json).Err(); err != nil {
		return err
	}

	fmt.Printf("SADD %s\n", json)

	return nil
}

func GetForumEntries(rdb *redis.Client, ctx context.Context) ([]*forum.Entry, error) {
	result, err := rdb.SMembers(ctx, "forum").Result()
	if err != nil {
		return []*forum.Entry{}, err
	}

	entries := make([]*forum.Entry, len(result))
	for i, json := range result {
		entry, err := forum.FromJson([]byte(json))
		if err != nil {
			return []*forum.Entry{}, err
		}

		entries[i] = entry
	}

	return entries, nil
}

func RedisGet(rdb *redis.Client, ctx context.Context, key string, defaultValue string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return defaultValue, nil
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func RedisSet(rdb *redis.Client, ctx context.Context, key string, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}
