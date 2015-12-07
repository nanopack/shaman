package caches

// General layered caching for the shaman dns server
// l1 is a short-term quick response lookup for entries
// l2 is a long-term storage for entries
// l1 and l2 can be configured to use different backend caches or databases

// TODO:
//  - implement caching backends
//  - add logging
//  - test

import (
	"fmt"
	"github.com/nanopack/shaman/config"
	"net/url"
)

type Cacher interface {
	GetRecord(string) (string, error)
	SetRecord(string, string) error
	ReviseRecord(string, string) error
	DeleteRecord(string) error
}

type cacheEntry struct {
	expires int64
	value   string
}

var (
	l1 Cacher
	l2 Cacher
)

// Determine the backend cache to initialize based off of the connection string
// Pass the connection string and TTL into the backend constructor
func initializeCacher(connection string, expires int) (Cacher, error) {
	u, err := url.Parse(connection)
	if err != nil {

	}
	var cacher Cacher
	switch u.Scheme {
	case "redis":
		cacher, err = NewRedisCacher(connection, expires)
	case "postgres":
		cacher, err = NewPostgresCacher(connection, expires)
	case "scribble":
		cacher, err = NewScribbleCacher(connection, expires)
	default:
		cacher, err = NewMapCacher(connection, expires)
	}
	if err != nil {
		return nil, err
	}
	return cacher, nil
}

// Create l1 and l2 from the config
func Init() error {
	l1, err := initializeCacher(config.L1Connect, config.L1Expires)
	_ = l1
	if err != nil {
		return err
	}
	l2, err := initializeCacher(config.L2Connect, config.L2Expires)
	_ = l2
	if err != nil {
		return err
	}
	return nil
}

// Create a lookup key based off of the domain and type of record
func Key(domain string, rtype uint16) string {
	return fmt.Sprintf("%d-%s", rtype, domain)
}

// Add record into the caches. First insert into the long term,
// then try the short term.
func AddRecord(key string, value string) error {
	if l2 != nil {
		err := l2.SetRecord(key, value)
		if err != nil {
			return nil
		}
	}
	if l1 != nil {
		err := l1.SetRecord(key, value)
		if err != nil {
			return nil
		}
	}
	return nil
}

// Remove record from the long term storage, then remove from short term.
func RemoveRecord(key string) error {
	if l2 != nil {
		err := l2.DeleteRecord(key)
		if err != nil {
			return err
		}
	}
	if l1 != nil {
		err := l1.DeleteRecord(key)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update the long term storage, then update the short term storage.
func UpdateRecord(key string, value string) error {
	if l2 != nil {
		err := l2.ReviseRecord(key, value)
		if err != nil {
			return err
		}
	}
	if l1 != nil {
		err := l1.ReviseRecord(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// Look for the record in the short term, if it isn't there, check the
// long term, and put it in the short term.
func FindRecord(key string) (string, error) {
	var record string
	if l1 != nil {
		record, err := l1.GetRecord(key)
		if err != nil {
			return record, err
		}
	}
	if record != "" {
		return record, nil
	}
	if l2 != nil {
		record, err := l2.GetRecord(key)
		if record != "" {
			l1.SetRecord(key, record)
			return record, err
		}
	}
	return "", nil
}
