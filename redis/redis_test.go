package redis

import (
	"fmt"
	"testing"

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
