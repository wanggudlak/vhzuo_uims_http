package gredis_test

import (
	redis2 "github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"uims/boot"
	"uims/pkg/gredis"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

func TestPing(t *testing.T) {
	pong, err := gredis.Def().Ping().Result()
	assert.Nil(t, err)
	t.Logf("Pong %s", pong)
}

func TestGetSet(t *testing.T) {
	// test exists key
	err := gredis.Def().Set("test", "test", 0*time.Second).Err()
	assert.Nil(t, err)
	ret, err := gredis.Def().Get("test").Result()
	assert.Equal(t, "test", ret)
	assert.Nil(t, err)

	err = gredis.Def().Del("test").Err()
	assert.Nil(t, err)

	// test not exists key
	ret2, err := gredis.Def().Get("test_no_exists").Result()
	assert.Equal(t, "", ret2)
	assert.NotNil(t, err)
	assert.IsType(t, redis2.Nil, err)
}

func TestExists(t *testing.T) {
	i, err := gredis.Def().Exists("test").Result()
	j := gredis.Def().Exists("test").Val()
	t.Log(j)
	assert.Nil(t, err)
	t.Log(i)

	err = gredis.Def().Set("test", "test", 5*time.Second).Err()
	assert.Nil(t, err)
	i, err = gredis.Def().Exists("test").Result()
	assert.Nil(t, err)
	t.Log(i)
}

func TestH(t *testing.T) {
	err := gredis.Def().HSet("testh", "a", true).Err()
	assert.Nil(t, err)
	all := gredis.Def().HGetAll("testh").Val()
	assert.NotNil(t, all)
	if i, ok := all["a"]; ok {
		assert.Equal(t, "1", i)
	} else {
		t.Error("应该存在的")
		t.FailNow()
	}
}
