package caches_test

import (
	"testing"
	"time"

	"github.com/nanopack/shaman/caches"
)

func initializePostgresqlCacher(t *testing.T, expires int) caches.Cacher {
	cacher, err := caches.NewPostgresCacher("postgres://postgres@localhost/travis_ci_test?sslmode=disable", expires)
	cacher.InitializeDatabase()
	cacher.ClearDatabase()
	if err != nil {
		t.Errorf("Error from initializePostgresqlCacher in PostgresqlCacher: %s", err)
	}
	return cacher
}

func postgresqlSet(t *testing.T, postgresqlCacher caches.Cacher, key string, value string) {
	err := postgresqlCacher.SetRecord(key, value)
	if err != nil {
		t.Errorf("Error from SetRecord in PostgresqlCacher: %s", err)
	}
}

func postgresqlGet(t *testing.T, postgresqlCacher caches.Cacher, key string, checkValue string) {
	value, err := postgresqlCacher.GetRecord(key)
	if err != nil {
		t.Errorf("Error from GetRecord in PostgresqlCacher: %s", err)
	}
	if value != checkValue {
		t.Errorf("Unexpected result from PostgresqlCacher: %s", value)
	}
}

func postgresqlRevise(t *testing.T, postgresqlCacher caches.Cacher, key string, value string) {
	err := postgresqlCacher.ReviseRecord(key, value)
	if err != nil {
		t.Errorf("Error from SetRecord in PostgresqlCacher: %s", err)
	}
}

func postgresqlDelete(t *testing.T, postgresqlCacher caches.Cacher, key string) {
	err := postgresqlCacher.DeleteRecord("1-key")
	if err != nil {
		t.Errorf("Error from DeleteRecord in PostgresqlCacher: %s", err)
	}
}

func postgresqlList(t *testing.T, postgresqlCacher caches.Cacher, key string, checkValues []string) {
	values, err := postgresqlCacher.ListRecords()
	if err != nil {
		t.Errorf("Error from ListRecord in PostgresqlCacher: %s", err)
	}
	if len(values) != len(checkValues) {
		t.Errorf("Unexpected length from ListRecord in PostgresqlCacher: %d", len(values))
	}
	for value := range values {
		found := false
		for checkValue := range checkValues {
			if checkValue == value {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected values from ListRecord in PostgresqlCacher: %s", values)
		}
	}
}

func TestPostgresqlSet(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
}

func TestPostgresqlGet(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlGet(t, postgresqlCacher, "1-key", "")
}

func TestPostgresqlGetAfterSet(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	time.Sleep(2 * time.Second)
	postgresqlGet(t, postgresqlCacher, "1-key", "found")
}

func TestPostgresqlGetAfterSetWithExpiresNoSleep(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 1)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	postgresqlGet(t, postgresqlCacher, "1-key", "found")
}

func TestPostgresqlGetAfterSetWithExpires(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 1)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	time.Sleep(2 * time.Second)
	postgresqlGet(t, postgresqlCacher, "1-key", "")
}

func TestPostgresqlRevise(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlRevise(t, postgresqlCacher, "1-key", "found")
}

func TestPostgresqlReviseAfterSet(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	postgresqlRevise(t, postgresqlCacher, "1-key", "found")
}

func TestPostgresqlGetAfterReviseAfterSet(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	postgresqlRevise(t, postgresqlCacher, "1-key", "found too")
	postgresqlGet(t, postgresqlCacher, "1-key", "found too")
}

func TestPostgresqlGetAfterReviseAfterSetWithExpires(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 1)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	postgresqlRevise(t, postgresqlCacher, "1-key", "found too")
	time.Sleep(2 * time.Second)
	postgresqlGet(t, postgresqlCacher, "1-key", "")
}

func TestPostgresqlDelete(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	postgresqlDelete(t, postgresqlCacher, "1-key")
	postgresqlGet(t, postgresqlCacher, "1-key", "")
}

func TestPostgresqlDeleteToo(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	postgresqlSet(t, postgresqlCacher, "2-key", "found too")
	postgresqlDelete(t, postgresqlCacher, "1-key")
	postgresqlGet(t, postgresqlCacher, "1-key", "")
	postgresqlGet(t, postgresqlCacher, "2-key", "found too")
}

func TestPostgresqlList(t *testing.T) {
	postgresqlCacher := initializePostgresqlCacher(t, 0)
	postgresqlSet(t, postgresqlCacher, "1-key", "found")
	postgresqlSet(t, postgresqlCacher, "2-key", "found too")
	postgresqlList(t, postgresqlCacher, "2-key", []string{"found", "found too"})
}
