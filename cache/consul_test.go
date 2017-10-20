package cache_test

import (
	"testing"

	"github.com/nanopack/shaman/cache"
	"github.com/nanopack/shaman/config"
	shaman "github.com/nanopack/shaman/core/common"
)

// test consul cache init
func TestConsulInitialize(t *testing.T) {
	config.L2Connect = "consul://127.0.0.1:8500"
	err := cache.Initialize()
	cache.Initialize()
	if err != nil {
		t.Errorf("Failed to initalize consul cacher - %v", err)
	}
}

// test consul cache addRecord
func TestConsulAddRecord(t *testing.T) {
	consulReset()
	err := cache.AddRecord(&nanopack)
	if err != nil {
		t.Errorf("Failed to add record to consul cacher - %v", err)
	}
}

// test consul cache getRecord
func TestConsulGetRecord(t *testing.T) {
	consulReset()
	cache.AddRecord(&nanopack)
	_, err := cache.GetRecord("nanobox.io")
	_, err2 := cache.GetRecord("nanopack.io")
	if err == nil || err2 != nil {
		t.Errorf("Failed to get record from consul cacher - %v%v", err, err2)
	}
}

// test consul cache updateRecord
func TestConsulUpdateRecord(t *testing.T) {
	consulReset()
	err := cache.UpdateRecord("nanobox.io", &nanopack)
	err2 := cache.UpdateRecord("nanopack.io", &nanopack)
	if err != nil || err2 != nil {
		t.Errorf("Failed to update record in consul cacher - %v%v", err, err2)
	}
}

// test consul cache deleteRecord
func TestConsulDeleteRecord(t *testing.T) {
	consulReset()
	err := cache.DeleteRecord("nanobox.io")
	cache.AddRecord(&nanopack)
	err2 := cache.DeleteRecord("nanopack.io")
	if err != nil || err2 != nil {
		t.Errorf("Failed to delete record from consul cacher - %v%v", err, err2)
	}
}

// test consul cache resetRecords
func TestConsulResetRecords(t *testing.T) {
	consulReset()
	err := cache.ResetRecords(&nanoBoth)
	if err != nil {
		t.Errorf("Failed to reset records in consul cacher - %v", err)
	}
}

// test consul cache listRecords
func TestConsulListRecords(t *testing.T) {
	consulReset()
	_, err := cache.ListRecords()
	cache.ResetRecords(&nanoBoth)
	_, err2 := cache.ListRecords()
	if err != nil || err2 != nil {
		t.Errorf("Failed to list records in consul cacher - %v%v", err, err2)
	}
}

func consulReset() {
	config.L2Connect = "consul://127.0.0.1:8500"
	cache.Initialize()
	blank := make([]shaman.Resource, 0, 0)
	cache.ResetRecords(&blank)
}
