package cache_test

import (
	"os"
	"testing"

	"github.com/nanopack/shaman/cache"
	"github.com/nanopack/shaman/config"
)

// test scribble cache init
func TestScribbleInitialize(t *testing.T) {
	config.L2Connect = "/tmp/shamanCache" // default
	err := cache.Initialize()
	config.L2Connect = "!@#$%^&*()" // unparse-able
	err2 := cache.Initialize()
	config.L2Connect = "scribble:///roots/file" // unable to init? (test no sudo)
	err3 := cache.Initialize()
	config.L2Connect = "scribble:///" // defaulting to "/var/db"
	err4 := cache.Initialize()
	if err != nil || err2 == nil || err3 == nil || err4 != nil {
		t.Errorf("Failed to initalize scribble cacher - %v%v%v%v", err, err2, err3, err4)
	}
}

// test scribble cache addRecord
func TestScribbleAddRecord(t *testing.T) {
	scribbleReset()
	err := cache.AddRecord(&nanopack)
	if err != nil {
		t.Errorf("Failed to add record to scribble cacher - %v", err)
	}
}

// test scribble cache getRecord
func TestScribbleGetRecord(t *testing.T) {
	scribbleReset()
	cache.AddRecord(&nanopack)
	_, err := cache.GetRecord("nanobox.io")
	_, err2 := cache.GetRecord("nanopack.io")
	if err == nil || err2 != nil {
		t.Errorf("Failed to get record from scribble cacher - %v%v", err, err2)
	}
}

// test scribble cache updateRecord
func TestScribbleUpdateRecord(t *testing.T) {
	scribbleReset()
	err := cache.UpdateRecord("nanobox.io", &nanopack)
	err2 := cache.UpdateRecord("nanopack.io", &nanopack)
	if err != nil || err2 != nil {
		t.Errorf("Failed to update record in scribble cacher - %v%v", err, err2)
	}
}

// test scribble cache deleteRecord
func TestScribbleDeleteRecord(t *testing.T) {
	scribbleReset()
	err := cache.DeleteRecord("nanobox.io")
	cache.AddRecord(&nanopack)
	err2 := cache.DeleteRecord("nanopack.io")
	if err != nil || err2 != nil {
		t.Errorf("Failed to delete record from scribble cacher - %v%v", err, err2)
	}
}

// test scribble cache resetRecords
func TestScribbleResetRecords(t *testing.T) {
	scribbleReset()
	err := cache.ResetRecords(&nanoBoth)
	if err != nil {
		t.Errorf("Failed to reset records in scribble cacher - %v", err)
	}
}

// test scribble cache listRecords
func TestScribbleListRecords(t *testing.T) {
	scribbleReset()
	_, err := cache.ListRecords()
	cache.ResetRecords(&nanoBoth)
	_, err2 := cache.ListRecords()
	if err != nil || err2 != nil {
		t.Errorf("Failed to list records in scribble cacher - %v%v", err, err2)
	}
}

func scribbleReset() {
	os.RemoveAll("/tmp/shamanCache")
	config.L2Connect = "scribble:///tmp/shamanCache"
	cache.Initialize()
}
