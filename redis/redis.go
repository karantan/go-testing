package redis

import (
	"context"
	"encoding/json"
	"time"

	rtgo "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisDatabase struct {
	*rtgo.Client
}

// NewRedis is an instance of Redis client for this environment
func NewRedis(connectionString string) (*RedisDatabase, error) {
	c, err := rtgo.ParseURL(connectionString)
	return &RedisDatabase{rtgo.NewClient(c)}, err
}

// GetKey Redis `GET key` command.
func GetKey(rdb *RedisDatabase, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	return val, err
}

// SetKey Redis `SET key value [expiration]` command.
func SetKey(rdb *RedisDatabase, key, value string, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}

// GetStruct populates `dest` with whatever it is in `key`. Make sure you pass `dest`
// as a pointer!
func GetStruct(rdb *RedisDatabase, key string, dest interface{}) error {
	p, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(p), dest)
}

// SetStruct does the Redis `SET key value [expiration]` command.
func SetStruct(rdb *RedisDatabase, key string, value interface{}, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = rdb.Set(ctx, key, p, expiration).Result()
	return err
}
