package caches

// Caching implementation using scribble as the backend

// TODO:
//  - add logging
//  - test

import (
	scribble "github.com/nanobox-io/golang-scribble"
	"net/url"
	"time"
)

type scribbleCacher struct {
	expires    int
	scribbleDb *scribble.Driver
}

// Initialize a new scribble cacher
func NewScribbleCacher(connection string, expires int) (*scribbleCacher, error) {
	u, err := url.Parse(connection)
	if err != nil {
		return nil, err
	}
	dir := u.Path
	db, err := scribble.New(dir, nil)
	if err != nil {
		return nil, err
	}
	sC := scribbleCacher{
		expires:    expires,
		scribbleDb: db,
	}
	return &sC, nil
}

// Retrieve a record from the scribble database, update the expires if
func (self scribbleCacher) GetRecord(key string) (string, error) {
	ce := cacheEntry{}
	if err := self.scribbleDb.Read("records", key, ce); err != nil {
		return "", err
	}
	if self.expires > 0 {
		now := time.Now().Unix()
		if ce.expires < now {
			// expired
			self.DeleteRecord(key)
			return "", nil
		}
		newExpires := now + int64(self.expires)
		ce.expires = newExpires
		if err := self.scribbleDb.Write("records", key, ce); err != nil {
			return ce.value, nil
		}
	}

	return ce.value, nil
}

// Set record in scribble database
func (self scribbleCacher) SetRecord(key string, value string) error {
	var expires int64
	if self.expires > 0 {
		expires = time.Now().Unix() + int64(self.expires)
	}
	ce := cacheEntry{
		expires: expires,
		value:   value,
	}
	return self.scribbleDb.Write("records", key, ce)
}

// Update record in scribble database
func (self scribbleCacher) ReviseRecord(key string, value string) error {
	return self.SetRecord(key, value)
}

// Remove record from scribble database
func (self scribbleCacher) DeleteRecord(key string) error {
	return self.scribbleDb.Delete("records", key)
}
