package cache_test

import (
	"testing"

	"github.com/nanopack/shaman/cache"
	"github.com/nanopack/shaman/config"
	shaman "github.com/nanopack/shaman/core/common"
)

// test postgres cache init
func TestPostgresInitialize(t *testing.T) {
	config.L2Connect = "postgres://postgres@127.0.0.1?sslmode=disable" // default
	err := cache.Initialize()
	config.L2Connect = "postgresql://postgres@127.0.0.1:9999?sslmode=disable" // unable to init?
	err2 := cache.Initialize()
	if err != nil || err2 != nil {
		t.Errorf("Failed to initalize postgres cacher - %v%v", err, err2)
	}
}

// test postgres cache addRecord
func TestPostgresAddRecord(t *testing.T) {
	postgresReset()
	err := cache.AddRecord(&nanopack)
	if err != nil {
		t.Errorf("Failed to add record to postgres cacher - %v", err)
	}

	err = cache.AddRecord(&nanopack)
	if err != nil {
		t.Errorf("Failed to add record to postgres cacher - %v", err)
	}
}

// test postgres cache getRecord
func TestPostgresGetRecord(t *testing.T) {
	postgresReset()
	cache.AddRecord(&nanopack)
	_, err := cache.GetRecord("nanobox.io.")
	_, err2 := cache.GetRecord("nanopack.io")
	if err == nil || err2 != nil {
		t.Errorf("Failed to get record from postgres cacher - %v%v", err, err2)
	}
}

// test postgres cache updateRecord
func TestPostgresUpdateRecord(t *testing.T) {
	postgresReset()
	err := cache.UpdateRecord("nanobox.io", &nanopack)
	err2 := cache.UpdateRecord("nanopack.io", &nanopack)
	if err != nil || err2 != nil {
		t.Errorf("Failed to update record in postgres cacher - %v%v", err, err2)
	}
}

// test postgres cache deleteRecord
func TestPostgresDeleteRecord(t *testing.T) {
	postgresReset()
	err := cache.DeleteRecord("nanobox.io")
	cache.AddRecord(&nanopack)
	err2 := cache.DeleteRecord("nanopack.io")
	if err != nil || err2 != nil {
		t.Errorf("Failed to delete record from postgres cacher - %v%v", err, err2)
	}
}

// test postgres cache resetRecords
func TestPostgresResetRecords(t *testing.T) {
	postgresReset()
	err := cache.ResetRecords(&nanoBoth)
	if err != nil {
		t.Errorf("Failed to reset records in postgres cacher - %v", err)
	}
}

// test postgres cache listRecords
func TestPostgresListRecords(t *testing.T) {
	postgresReset()
	_, err := cache.ListRecords()
	cache.ResetRecords(&nanoBoth)
	_, err2 := cache.ListRecords()
	if err != nil || err2 != nil {
		t.Errorf("Failed to list records in postgres cacher - %v%v", err, err2)
	}
}

func postgresReset() {
	config.L2Connect = "postgres://postgres@127.0.0.1?sslmode=disable"
	cache.Initialize()
	blank := make([]shaman.Resource, 0, 0)
	cache.ResetRecords(&blank)
}
