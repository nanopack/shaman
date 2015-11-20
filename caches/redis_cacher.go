package caches

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

func (rc redisCacher) GetRecord(key string) (string, error) {
	conn := rc.connection.Get()
	defer conn.Close()
	record, err := redis.String(conn.Do("GET", key))
	if rc.expires > 0 {
		// refresh the expires
		_, err = conn.Do("EXPIRE", key, rc.expires)
	}
	return record, err
}

func (rc redisCacher) SetRecord(key string, value string) error {
	conn := rc.connection.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if rc.expires > 0 {
		// set the expires
		_, err = conn.Do("EXPIRE", key, rc.expires)
	}
	return err
}

func (rc redisCacher) ReviseRecord(key string, value string) error {
	conn := rc.connection.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	if rc.expires > 0 {
		// set the expires
		_, err = conn.Do("EXPIRE", key, rc.expires)
	}
	return err
}

func (rc redisCacher) DeleteRecord(key string) error {
	conn := rc.connection.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}
