package caches_test

import (
	"testing"
	"time"

	"github.com/nanopack/shaman/caches"
)

func initializeRedisCacher(t *testing.T, expires int) caches.Cacher {
	cacher, err := caches.NewRedisCacher("redis://localhost:6379", expires)
	cacher.ClearDatabase()
	if err != nil {
		t.Errorf("Error from initializeRedisCacher in RedisCacher: %s", err)
	}
	return cacher
}

func redisSet(t *testing.T, redisCacher caches.Cacher, key string, value string) {
	err := redisCacher.SetRecord(key, value)
	if err != nil {
		t.Errorf("Error from SetRecord in RedisCacher: %s", err)
	}
}

func redisGet(t *testing.T, redisCacher caches.Cacher, key string, checkValue string) {
	value, err := redisCacher.GetRecord(key)
	if err != nil {
		t.Errorf("Error from GetRecord in RedisCacher: %s", err)
	}
	if value != checkValue {
		t.Errorf("Unexpected result from RedisCacher: %s", value)
	}
}

func redisRevise(t *testing.T, redisCacher caches.Cacher, key string, value string) {
	err := redisCacher.ReviseRecord(key, value)
	if err != nil {
		t.Errorf("Error from SetRecord in RedisCacher: %s", err)
	}
}

func redisDelete(t *testing.T, redisCacher caches.Cacher, key string) {
	err := redisCacher.DeleteRecord("1-key")
	if err != nil {
		t.Errorf("Error from DeleteRecord in RedisCacher: %s", err)
	}
}

func redisList(t *testing.T, redisCacher caches.Cacher, key string, checkValues []string) {
	values, err := redisCacher.ListRecords()
	if err != nil {
		t.Errorf("Error from ListRecord in RedisCacher: %s", err)
	}
	if len(values) != len(checkValues) {
		t.Errorf("Unexpected length from ListRecord in RedisCacher: %d", len(values))
	}
	for value := range values {
		found := false
		for checkValue := range checkValues {
			if checkValue == value {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected values from ListRecord in RedisCacher: %s", values)
		}
	}
}

func TestRedisSet(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisSet(t, redisCacher, "1-key", "found")
}

func TestRedisGet(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisGet(t, redisCacher, "1-key", "")
}

func TestRedisGetAfterSet(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisSet(t, redisCacher, "1-key", "found")
	redisGet(t, redisCacher, "1-key", "found")
}

func TestRedisGetAfterSetWithExpiresNoSleep(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 1)
	redisSet(t, redisCacher, "1-key", "found")
	redisGet(t, redisCacher, "1-key", "found")
}

func TestRedisGetAfterSetWithExpires(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 1)
	redisSet(t, redisCacher, "1-key", "found")
	time.Sleep(2 * time.Second)
	redisGet(t, redisCacher, "1-key", "")
}

func TestRedisRevise(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisRevise(t, redisCacher, "1-key", "found")
}

func TestRedisReviseAfterSet(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisSet(t, redisCacher, "1-key", "found")
	redisRevise(t, redisCacher, "1-key", "found")
}

func TestRedisGetAfterReviseAfterSet(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisSet(t, redisCacher, "1-key", "found")
	redisRevise(t, redisCacher, "1-key", "found too")
	redisGet(t, redisCacher, "1-key", "found too")
}

func TestRedisGetAfterReviseAfterSetWithExpires(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 1)
	redisSet(t, redisCacher, "1-key", "found")
	redisRevise(t, redisCacher, "1-key", "found too")
	time.Sleep(2 * time.Second)
	redisGet(t, redisCacher, "1-key", "")
}

func TestRedisDelete(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisSet(t, redisCacher, "1-key", "found")
	redisDelete(t, redisCacher, "1-key")
	redisGet(t, redisCacher, "1-key", "")
}

func TestRedisDeleteToo(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisSet(t, redisCacher, "1-key", "found")
	redisSet(t, redisCacher, "2-key", "found too")
	redisDelete(t, redisCacher, "1-key")
	redisGet(t, redisCacher, "1-key", "")
	redisGet(t, redisCacher, "2-key", "found too")
}

func TestRedisList(t *testing.T) {
	redisCacher := initializeRedisCacher(t, 0)
	redisSet(t, redisCacher, "1-key", "found")
	redisSet(t, redisCacher, "2-key", "found too")
	redisList(t, redisCacher, "2-key", []string{"found", "found too"})
}
