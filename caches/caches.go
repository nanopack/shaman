package caches

// General layered caching for the shaman dns server
// L1 is a short-term quick response lookup for entries
// L2 is a long-term storage for entries
// L1 and L2 can be configured to use different backend caches or databases

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
	InitializeDatabase() error
	GetRecord(string) (string, error)
	SetRecord(string, string) error
	ReviseRecord(string, string) error
	DeleteRecord(string) error
	ListRecords() ([]string, error)
}

type cacheEntry struct {
	expires int64
	value   string
}

type FindReturn struct {
	Err   error
	Value string
}

type ListReturn struct {
	Err    error
	Values []string
}

type AddOp struct {
	Key   string
	Value string
	Resp  chan error
}

type UpdateOp struct {
	Key   string
	Value string
	Resp  chan error
}

type RemoveOp struct {
	Key  string
	Resp chan error
}

type FindOp struct {
	Key  string
	Resp chan FindReturn
}

type ListOp struct {
	Resp chan ListReturn
}

var (
	AddOps    = make(chan AddOp)
	RemoveOps = make(chan RemoveOp)
	UpdateOps = make(chan UpdateOp)
	FindOps   = make(chan FindOp)
	ListOps   = make(chan ListOp)
	L1        Cacher
	L2        Cacher
)

func StartCache() error {
	initCache()
	for {
		select {
		case addOp := <-AddOps:
			addOp.Resp <- addRecord(addOp.Key, addOp.Value)
		case removeOp := <-RemoveOps:
			removeOp.Resp <- removeRecord(removeOp.Key)
		case updateOp := <-UpdateOps:
			updateOp.Resp <- updateRecord(updateOp.Key, updateOp.Value)
		case findOp := <-FindOps:
			value, err := findRecord(findOp.Key)
			findOp.Resp <- FindReturn{Err: err, Value: value}
		case listOp := <-ListOps:
			values, err := listRecords()
			listOp.Resp <- ListReturn{Err: err, Values: values}
		}
	}
	return nil
}

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
	err = cacher.InitializeDatabase()
	if err != nil {
		return nil, err
	}
	return cacher, nil
}

// Create L1 and L2 from the config
func initCache() {
	config.Log.Info("Initializing caches")
	var err error
	L1, err = initializeCacher(config.L1Connect, config.L1Expires)
	if err != nil {
		config.Log.Error("Error with L1: %s", err)
	}
	L2, err = initializeCacher(config.L2Connect, config.L2Expires)
	if err != nil {
		config.Log.Error("Error with L2: %s", err)
	}
}

// Create a lookup key based off of the domain and type of record
func Key(domain string, rtype uint16) string {
	return fmt.Sprintf("%d-%s", rtype, domain)
}

// Add record into the caches. First insert into the long term,
// then try the short term.
func addRecord(key string, value string) error {
	config.Log.Info("Adding key: %s, value: %s", key, value)
	if L2 != nil {
		err := L2.SetRecord(key, value)
		if err != nil {
			config.Log.Error("Error adding to L2: %s", err)
			return err
		}
	}
	if L1 != nil {
		err := L1.SetRecord(key, value)
		if err != nil {
			config.Log.Error("Error adding to L1: %s", err)
			return err
		}
	}
	return nil
}

// Remove record from the long term storage, then remove from short term.
func removeRecord(key string) error {
	config.Log.Info("Removing key: %s", key)
	if L2 != nil {
		err := L2.DeleteRecord(key)
		if err != nil {
			config.Log.Error("Error removing from L2: %s", err)
			return err
		}
	}
	if L1 != nil {
		err := L1.DeleteRecord(key)
		if err != nil {
			config.Log.Error("Error removing from L1: %s", err)
			return err
		}
	}
	return nil
}

// Update the long term storage, then update the short term storage.
func updateRecord(key string, value string) error {
	config.Log.Info("Updating key: %s, value: %s", key, value)
	if L2 != nil {
		err := L2.ReviseRecord(key, value)
		if err != nil {
			config.Log.Error("Error updating L2: %s", err)
			return err
		}
	}
	if L1 != nil {
		err := L1.ReviseRecord(key, value)
		if err != nil {
			config.Log.Error("Error updating L1: %s", err)
			return err
		}
	}
	return nil
}

// Look for the record in the short term, if it isn't there, check the
// long term, and put it in the short term.
func findRecord(key string) (string, error) {
	config.Log.Info("Finding key: %s", key)
	var record string
	var err error
	if L1 != nil {
		record, err = L1.GetRecord(key)
		if err != nil {
			config.Log.Error("Error finding L1: %s", err)
			return record, err
		}
	}
	if record != "" {
		return record, nil
	}
	if L2 != nil {
		record, err = L2.GetRecord(key)
		if err != nil {
			config.Log.Error("Error finding L2: %s", err)
		}
		if record != "" {
			L1.SetRecord(key, record)
			return record, err
		}
	}
	return "", nil
}

func listRecords() ([]string, error) {
	if L2 != nil {
		return L2.ListRecords()
	}
	if L1 != nil {
		return L1.ListRecords()
	}
	return []string{}, nil
}
