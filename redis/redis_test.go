package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	assert := assert.New(t)
	mr, err := miniredis.Run()
	assert.NoError(err)

	client, err := NewRedis(fmt.Sprintf("redis://%s/0", mr.Addr()))
	assert.NoError(err)
	assert.IsType(&RedisDatabase{}, client)
}

func TestSetGetKey(t *testing.T) {
	assert := assert.New(t)
	mr, err := miniredis.Run()
	assert.NoError(err)
	client, err := NewRedis(fmt.Sprintf("redis://%s/0", mr.Addr()))
	assert.NoError(err)

	err = SetKey(client, "foo", "bar", 0)
	assert.NoError(err)

	got, err := GetKey(client, "foo")
	want := "bar"
	assert.NoError(err)
	assert.Equal(want, got)
}

func TestSetGetStruct(t *testing.T) {
	assert := assert.New(t)
	mr, err := miniredis.Run()
	assert.NoError(err)
	client, err := NewRedis(fmt.Sprintf("redis://%s/0", mr.Addr()))
	assert.NoError(err)

	type Foo struct {
		Key   string
		Value string
	}

	err = SetStruct(client, "foo", Foo{"admin", "secret"}, 0)
	assert.NoError(err)

	var got Foo
	err = GetStruct(client, "foo", &got)
	assert.NoError(err)
	assert.Equal(Foo{"admin", "secret"}, got)
}

func TestAddDomainSerries(t *testing.T) {
	t.Run("add domain serries and get the data back", func(t *testing.T) {
		assert := assert.New(t)
		mr, err := miniredis.Run()
		assert.Nil(err)
		client, err := NewRedis(fmt.Sprintf("redis://%s/0", mr.Addr()))
		assert.Nil(err)

		err = AddDomainSeries(
			client,
			"foo.com",
			Domain{1, 0, 0, 0},
			time.Date(2009, time.November, 10, 2, 1, 2, 3, time.UTC),
		)
		err = AddDomainSeries(
			client,
			"foo.com",
			Domain{12, 0, 0, 0},
			time.Date(2009, time.November, 11, 2, 1, 2, 3, time.UTC),
		)
		err = AddDomainSeries(
			client,
			"foo.com",
			Domain{13, 0, 0, 0},
			time.Date(2009, time.November, 12, 2, 1, 2, 3, time.UTC),
		)
		err = AddDomainSeries(
			client,
			"foo.com",
			Domain{14, 0, 0, 0},
			time.Date(2009, time.November, 13, 2, 1, 2, 3, time.UTC),
		)
		err = AddDomainSeries(
			client,
			"foo.com",
			Domain{15, 0, 0, 0},
			time.Date(2009, time.November, 14, 2, 1, 2, 3, time.UTC),
		)
		assert.NoError(err)

		got, err := GetDomainSeries(
			client,
			"foo.com",
			time.Date(2009, time.November, 10, 2, 1, 2, 3, time.UTC),
			time.Date(2009, time.November, 14, 2, 1, 2, 3, time.UTC),
		)
		assert.Nil(err)
		assert.Equal(Domain{55, 0, 0, 0}, got)

		got, err = GetDomainSeries(
			client,
			"foo.com",
			time.Date(2009, time.November, 12, 2, 1, 2, 3, time.UTC),
			time.Date(2009, time.November, 14, 2, 1, 2, 3, time.UTC),
		)
		assert.Nil(err)
		assert.Equal(Domain{42, 0, 0, 0}, got)
	})
}

func TestAggregateDomainStats(t *testing.T) {
	assert := assert.New(t)
	want := Domain{5, 6, 7, 8}
	got := AggregateDomainStats(
		[]Domain{
			{1, 1, 1, 1},
			{1, 2, 1, 1},
			{1, 1, 3, 1},
			{1, 1, 1, 4},
			{1, 1, 1, 1},
		},
	)
	assert.Equal(want, got)
}

func TestGetAllDomainSeries(t *testing.T) {
	assert := assert.New(t)
	mr, err := miniredis.Run()
	assert.Nil(err)
	client, err := NewRedis(fmt.Sprintf("redis://%s/0", mr.Addr()))
	assert.Nil(err)

	AddDomainSeries(
		client,
		"foo.com",
		Domain{1, 0, 0, 0},
		time.Date(2009, time.November, 10, 2, 1, 2, 3, time.UTC),
	)
	AddDomainSeries(
		client,
		"foo.com",
		Domain{12, 0, 0, 0},
		time.Date(2009, time.November, 11, 2, 1, 2, 3, time.UTC),
	)
	AddDomainSeries(
		client,
		"bar.com",
		Domain{1, 1, 0, 1},
		time.Date(2009, time.November, 12, 2, 1, 2, 3, time.UTC),
	)
	AddDomainSeries(
		client,
		"bar.com",
		Domain{1, 1, 1, 0},
		time.Date(2009, time.November, 13, 2, 1, 2, 3, time.UTC),
	)
	AddDomainSeries(
		client,
		"bar.com",
		Domain{1, 0, 1, 1},
		time.Date(2009, time.November, 14, 2, 1, 2, 3, time.UTC),
	)

	want := map[string]Domain{}
	want["foo.com"] = Domain{13, 0, 0, 0}
	want["bar.com"] = Domain{3, 2, 2, 2}
	want["barfoo.com"] = Domain{0, 0, 0, 0}
	want["foobar.com"] = Domain{0, 0, 0, 0}
	got, _ := GetAllDomainSeries(
		client,
		time.Date(2009, time.November, 10, 2, 1, 2, 3, time.UTC),
		time.Date(2009, time.November, 14, 2, 1, 2, 3, time.UTC),
	)

	assert.Equal(want, got)
}

func Test_enumerateTimeRange(t *testing.T) {
	assert := assert.New(t)
	start := time.Date(2009, time.November, 10, 2, 1, 2, 3, time.UTC)
	end := time.Date(2009, time.November, 15, 23, 2, 3, 4, time.UTC)

	want := []time.Time{
		time.Date(2009, time.November, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.November, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.November, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.November, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.November, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2009, time.November, 15, 0, 0, 0, 0, time.UTC),
	}
	got := enumerateTimeRange(start, end)
	assert.Equal(want, got)
}
