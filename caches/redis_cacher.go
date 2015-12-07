package caches

// This stores entries in a Redis cache.
// Redis handles the expires, this only needs to refresh the expire.

// TODO:
//  - add logging
//  - test

import (
	"github.com/garyburd/redigo/redis"
	"net/url"
	"time"
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
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
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
	password, _ := u.User.Password()
	if err != nil {
		return nil, err
	}
	rc := redisCacher{
		expires:    expires,
		connection: newPool(u.Host, password),
	}
	return &rc, nil
}

// Retrieve record from Redis, update its expire time
func (self redisCacher) GetRecord(key string) (string, error) {
	conn := self.connection.Get()
	defer conn.Close()
	record, err := redis.String(conn.Do("GET", key))
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
