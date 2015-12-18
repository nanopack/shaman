package caches

// Caching implementation using scribble as the backend

// TODO:
//  - add logging
//  - test

import (
	"encoding/json"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/nanopack/shaman/config"
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

func (self scribbleCacher) InitializeDatabase() error {
	return nil
}

func (self scribbleCacher) ClearDatabase() error {
	self.scribbleDb.Delete("records", "")
	return nil
}

// Retrieve a record from the scribble database, update the expires if
func (self scribbleCacher) GetRecord(key string) (string, error) {
	ce := CacheEntry{}
	if err := self.scribbleDb.Read("records", key, &ce); err != nil {
		config.Log.Error("Error: %s", err)
		return "", nil
	}
	if self.expires > 0 {
		now := time.Now().Unix()
		if ce.Expires < now {
			// expired
			self.DeleteRecord(key)
			return "", nil
		}
		newExpires := now + int64(self.expires)
		ce.Expires = newExpires
		if err := self.scribbleDb.Write("records", key, ce); err != nil {
			return ce.Value, nil
		}
	}
	return ce.Value, nil
}

// Set record in scribble database
func (self scribbleCacher) SetRecord(key string, value string) error {
	var expires int64
	if self.expires > 0 {
		expires = time.Now().Unix() + int64(self.expires)
	}
	ce := CacheEntry{
		Expires: expires,
		Value:   value,
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

func (self scribbleCacher) ListRecords() ([]string, error) {
	entries := make([]string, 0)
	now := time.Now().Unix()
	values, err := self.scribbleDb.ReadAll("records")
	if err != nil {
		return entries, err
	}
	for i := range values {
		var ce CacheEntry
		json.Unmarshal([]byte(values[i]), &ce)
		if self.expires != 0 {
			if ce.Expires > now {
				entries = append(entries, ce.Value)
			}
		} else {
			entries = append(entries, ce.Value)
		}
	}
	return entries, nil
}
