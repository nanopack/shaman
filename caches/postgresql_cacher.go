package caches

// This stores entries in a Postgresql database.
// Postgresql doesn't handle expiring data automatically, this will
// need to verify data hasn't expired yet.

// TODO:
//  - add logging
//  - test
//  - add routine for removing old data

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/nanopack/shaman/config"
	"time"
)

type postgresqlCacher struct {
	expires int
	db      *sql.DB
}

// Initializer for the postgres cacher
func NewPostgresCacher(connection string, expires int) (*postgresqlCacher, error) {
	config.Log.Info("creating postgresql cacher")
	db, err := sql.Open("postgres", connection)
	if err != nil {
		config.Log.Error("error: %s", err)
		return nil, err
	}
	pc := postgresqlCacher{
		expires: expires,
		db:      db,
	}
	return &pc, nil
}

func (self postgresqlCacher) InitializeDatabase() error {
	rows, err := self.db.Query("CREATE TABLE IF NOT EXISTS dns_entries ( key varchar(128) UNIQUE, value text, expires bigint)")
	if err != nil {
		config.Log.Error("error: %s", err)
		return err
	}
	defer rows.Close()
	return nil
}

func (self postgresqlCacher) ClearDatabase() error {
	rows, err := self.db.Query("DELETE FROM dns_entries")
	if err != nil {
		config.Log.Error("error: %s", err)
		return err
	}
	defer rows.Close()
	return nil
}

// Retrieve record and check to make sure it isn't expired, update expires if needed.
func (self postgresqlCacher) GetRecord(key string) (string, error) {
	var value string
	var expires int64
	var err error
	var rows *sql.Rows
	err = self.db.QueryRow("SELECT value, expires FROM dns_entries WHERE key = $1", key).Scan(&value, &expires)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		} else {
			config.Log.Error("error: %s, %s", err, value)
			return value, err
		}
	}

	if self.expires > 0 {
		now := time.Now().Unix()
		if expires < now {
			// expired
			self.DeleteRecord(key)
			return "", nil
		}
		newExpires := now + int64(self.expires)
		rows, err = self.db.Query("UPDATE dns_entries SET expires=$1 WHERE key = $2", newExpires, key)
		if err != nil {
			config.Log.Error("error: %s", err)
			return value, err
		}
		defer rows.Close()
	}
	return value, nil
}

// Insert new record in the database, update expires if needed.
func (self postgresqlCacher) SetRecord(key string, value string) error {
	now := time.Now().Unix()
	expires := now + int64(self.expires)
	rows, err := self.db.Query("INSERT INTO dns_entries (key, value, expires) VALUES ($1, $2, $3)", key, value, expires)
	if err != nil {
		config.Log.Error("error: %s", err)
		return err
	}
	defer rows.Close()
	return nil
}

// Update existing record, update expires if needed.
func (self postgresqlCacher) ReviseRecord(key string, value string) error {
	now := time.Now().Unix()
	expires := now + int64(self.expires)
	rows, err := self.db.Query("UPDATE dns_entries SET value=$1, expires=$2 WHERE key = $3", value, expires, key)
	if err != nil {
		config.Log.Error("error: %s", err)
		return err
	}
	defer rows.Close()
	return nil
}

// Remove record from database.
func (self postgresqlCacher) DeleteRecord(key string) error {
	rows, err := self.db.Query("DELETE FROM dns_entries WHERE key = $1", key)
	if err != nil {
		config.Log.Error("error: %s", err)
		return err
	}
	defer rows.Close()
	return nil
}

func (self postgresqlCacher) ListRecords() ([]string, error) {
	entries := make([]string, 0)
	now := time.Now().Unix()
	var value string
	var expires int64
	rows, err := self.db.Query("SELECT value, expires FROM dns_entries")
	if err != nil {
		config.Log.Error("error: %s", err)
		return entries, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&value, &expires)
		if err != nil {
			config.Log.Error("Error: %s", err)
		}
		if self.expires > 0 {
			if expires > now {
				entries = append(entries, value)
			}
		} else {
			entries = append(entries, value)
		}

	}
	return entries, nil
}
