package caches

// TODO:
//  - add logging
//  - test

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type postgresqlCacher struct {
	expires int
	db      *sql.DB
}

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

func createTable(pc postgresqlCacher) error {
	stmt, err := pc.db.Prepare("CREATE TABLE IF NOT EXISTS dns_entries ( key varchar(128) UNIQUE, value text, expires bigint)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	stmt.Query()
	return nil
}

func (pc postgresqlCacher) GetRecord(key string) (string, error) {
	stmt, err := pc.db.Prepare("SELECT value, expires FROM dns_entries WHERE key = $1")
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

	if pc.expires > 0 {
		now := time.Now().Unix()
		if expires < now {
			// expired
			return "", nil
		}
		newExpires := now + int64(pc.expires)
		stmt2, err := pc.db.Prepare("UPDATE dns_entries SET expires=$1 WHERE key = $2")
		defer stmt2.Close()
		if err != nil {
			return value, err
		}
		stmt2.Query(newExpires, key)
	}
	return value, nil
}

func (pc postgresqlCacher) SetRecord(key string, value string) error {
	stmt, err := pc.db.Prepare("INSERT INTO dns_entries (key, value) VALUES ($1, $2)")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Query(key, value)
	if err != nil {
		return err
	}
	if pc.expires > 0 {
		now := time.Now().Unix()
		newExpires := now + int64(pc.expires)
		stmt2, err := pc.db.Prepare("UPDATE dns_entries SET expires=$1 WHERE key = $2")
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

func (pc postgresqlCacher) ReviseRecord(key string, value string) error {
	stmt, err := pc.db.Prepare("UPDATE dns_entries SET value=$1 WHERE key = $2")
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Query(value, key)
	if err != nil {
		return err
	}
	if pc.expires > 0 {
		now := time.Now().Unix()
		newExpires := now + int64(pc.expires)
		stmt2, err := pc.db.Prepare("UPDATE dns_entries SET expires=$1 WHERE key = $2")
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

func (pc postgresqlCacher) DeleteRecord(key string) error {
	stmt, err := pc.db.Prepare("DELETE FROM dns_entries WHERE key = $1")
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
