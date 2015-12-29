package caches

// This stores entries in a Redis cache.
// Redis handles the expires, this only needs to refresh the expire.

// TODO:
//  - add logging
//  - test

import (
	"net/url"
	"time"

	"github.com/garyburd/redigo/redis"
)

type redisCacher struct {
	expires    int
	connection *redis.Pool
}

// Redis connection pool
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Initialize a Redis cacher
func NewRedisCacher(connection string, expires int) (*redisCacher, error) {
	u, err := url.Parse(connection)
	var password string
	user := u.User
	if user != nil {
		password, _ = user.Password()
	}
	if err != nil {
		return nil, err
	}
	rc := redisCacher{
		expires:    expires,
		connection: newPool(u.Host, password),
	}
	return &rc, nil
}

func (self redisCacher) InitializeDatabase() error {
	return nil
}

func (self redisCacher) ClearDatabase() error {
	conn := self.connection.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	return err
}

// Retrieve record from Redis, update its expire time
func (self redisCacher) GetRecord(key string) (string, error) {
	conn := self.connection.Get()
	defer conn.Close()
	record, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return "", nil
	}
	if self.expires > 0 {
		// refresh the expires
		_, err = conn.Do("EXPIRE", key, self.expires)
	}
	return record, err
}

// Insert record in the cache, reset the expire time
func (self redisCacher) SetRecord(key string, value string) error {
	conn := self.connection.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if self.expires > 0 {
		// set the expires
		_, err = conn.Do("EXPIRE", key, self.expires)
	}
	return err
}

// Update entry
func (self redisCacher) ReviseRecord(key string, value string) error {
	return self.SetRecord(key, value)
}

// Remove entry
func (self redisCacher) DeleteRecord(key string) error {
	conn := self.connection.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

func (self redisCacher) ListRecords() ([]string, error) {
	entries := make([]string, 0)
	iter := 0
	conn := self.connection.Get()
	defer conn.Close()
	for {
		arr, err := redis.MultiBulk(conn.Do("SCAN", iter))
		if err != nil {
			return entries, err
		}
		iter, _ = redis.Int(arr[0], nil)
		keys, _ := redis.Strings(arr[1], nil)
		for key := range keys {
			record, err := redis.String(conn.Do("GET", keys[key]))
			if err != nil {
				if err != redis.ErrNil {
					return entries, err
				}
			} else {
				entries = append(entries, record)
			}
		}

		if iter == 0 {
			break
		}
	}
	return entries, nil
}
