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

	if err := rdb.HSetNX(ctx, "forum", forumEntry.Id, json).Err(); err != nil {
		return err
	}

	fmt.Printf("HSETNX forum %s %s\n", forumEntry.Id, json)

	return nil
}

func GetForumEntries(rdb *redis.Client, ctx context.Context) ([]*forum.Entry, error) {
	result, err := rdb.HGetAll(ctx, "forum").Result()
	if err != nil {
		return []*forum.Entry{}, err
	}

	entries := make([]*forum.Entry, len(result))
	i := 0
	for id, json := range result {
		entry, err := forum.FromJson([]byte(json))
		if err != nil {
			return []*forum.Entry{}, err
		}

		entry.Id = id
		entries[i] = entry
		i++
	}

	return entries, nil
}

func DeleteForumEntry(rdb *redis.Client, ctx context.Context, id string) error {
	if err := rdb.HDel(ctx, "forum", id).Err(); err != nil {
		return err
	}

	fmt.Printf("HDEL forum %s\n", id)

	return nil
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
