package caches

// Simple cache that stores data in a simple go map.
// map doesn't automatically evict expired data, this will need to
// check to ensure data isn't already expired.

// TODO:
//  - add logging
//  - test
//  - add routine for removing old data

import "time"

type mapCacher struct {
	expires int
	db      map[string]cacheEntry
}

// Map cacher initializer
func NewMapCacher(connection string, expires int) (*mapCacher, error) {
	mc := mapCacher{expires: expires, db: make(map[string]cacheEntry)}
	return &mc, nil
}

// Get record from the map cacher and make sure it hasn't expired yet
func (self mapCacher) GetRecord(key string) (string, error) {
	ce := cacheEntry{}
	ce, ok := self.db[key]
	if !ok {
		return "", nil
	}
	if self.expires > 0 {
		if time.Now().Unix() > ce.expires {
			// expired
			self.DeleteRecord(key)
			return "", nil
		}
		ce.expires = time.Now().Unix() + int64(self.expires)
		self.db[key] = ce
	}
	return ce.value, nil
}

// Insert/update entry in the map cacher
func (self mapCacher) SetRecord(key, val string) error {
	ce := cacheEntry{expires: time.Now().Unix() + int64(self.expires), value: val}
	self.db[key] = ce
	return nil
}

// Update entry in the map cacher
func (self mapCacher) ReviseRecord(key, val string) error {
	return self.SetRecord(key, val)
}

// remove entry from the map cacher
func (self mapCacher) DeleteRecord(key string) error {
	delete(self.db, key)
	return nil
}
