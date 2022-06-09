package redis

// See this article for another example:
// https://blog.logrocket.com/how-to-use-redis-as-a-database-with-go-redis/

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Database struct {
	*redis.Client
}

func NewClient(connectionString string) (*Database, error) {
	opt, err := redis.ParseURL(connectionString)
	return &Database{redis.NewClient(opt)}, err
}

func GetKey(rdb *Database, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	return val, err
}

func SetKey(rdb *Database, key, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}
