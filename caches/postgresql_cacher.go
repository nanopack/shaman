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
	"time"
)

type postgresqlCacher struct {
	expires int
	db      *sql.DB
}

// Initializer for the postgres cacher
func NewPostgresCacher(connection string, expires int) (*postgresqlCacher, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	pc := postgresqlCacher{
		expires: expires,
		db:      db,
	}
	createTable(pc)
	return &pc, nil
}

// Create the database table to store entries
func createTable(pc postgresqlCacher) error {
	stmt, err := pc.db.Prepare("CREATE TABLE IF NOT EXISTS dns_entries ( key varchar(128) UNIQUE, value text, expires bigint)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	stmt.Query()
	return nil
}

// Retrieve record and check to make sure it isn't expired, update expires if needed.
func (self postgresqlCacher) GetRecord(key string) (string, error) {
	stmt, err := self.db.Prepare("SELECT value, expires FROM dns_entries WHERE key = $1")
	defer stmt.Close()
	var value string
	var expires int64
	if err != nil {
		return "", err
	}
	err = stmt.QueryRow(key).Scan(&value, &expires)
	if err != nil {
		return "", err
	}

	if self.expires > 0 {
		now := time.Now().Unix()
		if expires < now {
			// expired
			self.DeleteRecord(key)
			return "", nil
		}
		newExpires := now + int64(self.expires)
		stmt2, err := self.db.Prepare("UPDATE dns_entries SET expires=$1 WHERE key = $2")
		defer stmt2.Close()
		if err != nil {
			return value, err
		}
		stmt2.Query(newExpires, key)
	}
	return value, nil
}

// Insert new record in the database, update expires if needed.
func (self postgresqlCacher) SetRecord(key string, value string) error {
	stmt, err := self.db.Prepare("INSERT INTO dns_entries (key, value) VALUES ($1, $2)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Query(key, value)
	if err != nil {
		return err
	}
	if self.expires > 0 {
		now := time.Now().Unix()
		newExpires := now + int64(self.expires)
		stmt2, err := self.db.Prepare("UPDATE dns_entries SET expires=$1 WHERE key = $2")
		defer stmt2.Close()
		if err != nil {
			return err
		}
		_, err = stmt2.Query(newExpires, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update existing record, update expires if needed.
func (self postgresqlCacher) ReviseRecord(key string, value string) error {
	stmt, err := self.db.Prepare("UPDATE dns_entries SET value=$1 WHERE key = $2")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Query(value, key)
	if err != nil {
		return err
	}
	if self.expires > 0 {
		now := time.Now().Unix()
		newExpires := now + int64(self.expires)
		stmt2, err := self.db.Prepare("UPDATE dns_entries SET expires=$1 WHERE key = $2")
		defer stmt2.Close()
		if err != nil {
			return err
		}
		_, err = stmt2.Query(newExpires, key)
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove record from database.
func (self postgresqlCacher) DeleteRecord(key string) error {
	stmt, err := self.db.Prepare("DELETE FROM dns_entries WHERE key = $1")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Query(key)
	if err != nil {
		return err
	}
	return nil
}
