package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rtgo "github.com/go-redis/redis/v8"
)

const (
	REDIS_PREFIX = "myapp"
	RFC3339Day   = "2006-01-02"
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

//
// Example of redis data serries. More info:
//  - https://redis.io/commands/zadd/
//  - https://redis.io/commands/zrange/
//

type Domain struct {
	Req2xx int64 `json:"2xx"`
	Req3xx int64 `json:"3xx"`
	Req4xx int64 `json:"4xx"`
	Req5xx int64 `json:"5xx"`
}

// AggregateDomainStats sums all the stats values in the `domainStats` list and returns
// one `Domain` struct with sum stats values.
func AggregateDomainStats(domainStats []Domain) (d Domain) {
	for _, ds := range domainStats {
		d.Req2xx += ds.Req2xx
		d.Req3xx += ds.Req3xx
		d.Req4xx += ds.Req4xx
		d.Req5xx += ds.Req5xx
	}
	return
}

// AddDomainSeries adds a `Domain` value to the `domain` which is a sorted set in Redis.
func AddDomainSeries(rdb *RedisDatabase, domain string, value Domain, timestamp time.Time) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.ZAdd(ctx, fmt.Sprintf("%s-domain-%s", REDIS_PREFIX, domain), &rtgo.Z{
		Score:  float64(timestamp.UTC().UnixMilli()),
		Member: p,
	}).Err()
}

// GetDomainSeries returns aggregated `Domain` statistics in the sorted set at
// `domain` key with a score between `start` and `stop` time.
// This function is the counterpart of the `AddDomainSeries`.
func GetDomainSeries(rdb *RedisDatabase, domain string, start, stop time.Time) (Domain, error) {
	results, err := rdb.ZRangeByScore(ctx, fmt.Sprintf("%s-domain-%s", REDIS_PREFIX, domain), &rtgo.ZRangeBy{
		Min: fmt.Sprint(start.UTC().UnixMilli()),
		Max: fmt.Sprint(stop.UTC().UnixMilli()),
	}).Result()

	if err != nil {
		return Domain{}, err
	}

	domainStats := []Domain{}
	for _, r := range results {
		var tmpDomain Domain
		err := json.Unmarshal([]byte(r), &tmpDomain)
		if err != nil {
			return Domain{}, err
		}
		domainStats = append(domainStats, tmpDomain)
	}
	return AggregateDomainStats(domainStats), nil
}

// GetAllDomainSeries does the same thing as `GetDomainSeries` except it returns
// `Domain` stats for all active pages (i.e. they have `config.json` file) we have in
// our Redis store.
func GetAllDomainSeries(rdb *RedisDatabase, start, stop time.Time) (map[string]Domain, error) {
	store := map[string]Domain{}
	domains := GetDomains()
	for _, domain := range domains {
		d, err := GetDomainSeries(rdb, domain, start, stop)
		if err != nil {
			return map[string]Domain{}, err
		}
		store[domain] = d
	}
	return store, nil
}

// GetDomains is a mocked implementation of a database or some other source of domains.
func GetDomains() []string {
	return []string{"foo.com", "foobar.com", "bar.com", "barfoo.com"}
}

// enumerateTimeRange is a helper function that creates `[]time.Time` with
// start and end dates.
// Let's say the start time is "2020-01-01 12:34" and end time is "2020-01-05 14:56" the
// output will be:
//		[]time.Time{2020-01-01, 2020-01-02, 2020-01-03, 2020-01-04}
//
func enumerateTimeRange(start, end time.Time) (enumeratedTime []time.Time) {
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

	days := 0
	for {
		newDate := startDate.AddDate(0, 0, days)
		enumeratedTime = append(enumeratedTime, newDate)
		days++

		if newDate.Equal(endDate) {
			break
		}
	}
	return
}
