package caches_test

import (
	"testing"
	"time"

	"github.com/nanopack/shaman/caches"
)

func initializeMapCacher(expires int) caches.Cacher {
	cacher, _ := caches.NewMapCacher("", expires)
	cacher.ClearDatabase()
	return cacher
}

func mapSet(t *testing.T, mapCacher caches.Cacher, key string, value string) {
	err := mapCacher.SetRecord(key, value)
	if err != nil {
		t.Errorf("Error from SetRecord in MapCacher: %s", err)
	}
}

func mapGet(t *testing.T, mapCacher caches.Cacher, key string, checkValue string) {
	value, err := mapCacher.GetRecord(key)
	if err != nil {
		t.Errorf("Error from GetRecord in MapCacher: %s", err)
	}
	if value != checkValue {
		t.Errorf("Unexpected result from MapCacher: %s", value)
	}
}

func mapRevise(t *testing.T, mapCacher caches.Cacher, key string, value string) {
	err := mapCacher.ReviseRecord(key, value)
	if err != nil {
		t.Errorf("Error from SetRecord in MapCacher: %s", err)
	}
}

func mapDelete(t *testing.T, mapCacher caches.Cacher, key string) {
	err := mapCacher.DeleteRecord("1-key")
	if err != nil {
		t.Errorf("Error from DeleteRecord in MapCacher: %s", err)
	}
}

func mapList(t *testing.T, mapCacher caches.Cacher, key string, checkValues []string) {
	values, err := mapCacher.ListRecords()
	if err != nil {
		t.Errorf("Error from ListRecord in MapCacher: %s", err)
	}
	if len(values) != len(checkValues) {
		t.Errorf("Unexpected length from ListRecord in MapCacher: %d", len(values))
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
			t.Errorf("Unexpected values from ListRecord in MapCacher: %s", values)
		}
	}
}

func TestMapSet(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapSet(t, mapCacher, "1-key", "found")
}

func TestMapGet(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapGet(t, mapCacher, "1-key", "")
}

func TestMapGetAfterSet(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapSet(t, mapCacher, "1-key", "found")
	mapGet(t, mapCacher, "1-key", "found")
}

func TestMapGetAfterSetWithExpiresNoSleep(t *testing.T) {
	mapCacher := initializeMapCacher(1)
	mapSet(t, mapCacher, "1-key", "found")
	mapGet(t, mapCacher, "1-key", "found")
}

func TestMapGetAfterSetWithExpires(t *testing.T) {
	mapCacher := initializeMapCacher(1)
	mapSet(t, mapCacher, "1-key", "found")
	time.Sleep(2 * time.Second)
	mapGet(t, mapCacher, "1-key", "")
}

func TestMapRevise(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapRevise(t, mapCacher, "1-key", "found")
}

func TestMapReviseAfterSet(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapSet(t, mapCacher, "1-key", "found")
	mapRevise(t, mapCacher, "1-key", "found")
}

func TestMapGetAfterReviseAfterSet(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapSet(t, mapCacher, "1-key", "found")
	mapRevise(t, mapCacher, "1-key", "found too")
	mapGet(t, mapCacher, "1-key", "found too")
}

func TestMapGetAfterReviseAfterSetWithExpires(t *testing.T) {
	mapCacher := initializeMapCacher(1)
	mapSet(t, mapCacher, "1-key", "found")
	mapRevise(t, mapCacher, "1-key", "found too")
	time.Sleep(2 * time.Second)
	mapGet(t, mapCacher, "1-key", "")
}

func TestMapDelete(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapSet(t, mapCacher, "1-key", "found")
	mapDelete(t, mapCacher, "1-key")
	mapGet(t, mapCacher, "1-key", "")
}

func TestMapDeleteToo(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapSet(t, mapCacher, "1-key", "found")
	mapSet(t, mapCacher, "2-key", "found too")
	mapDelete(t, mapCacher, "1-key")
	mapGet(t, mapCacher, "1-key", "")
	mapGet(t, mapCacher, "2-key", "found too")
}

func TestMapList(t *testing.T) {
	mapCacher := initializeMapCacher(0)
	mapSet(t, mapCacher, "1-key", "found")
	mapSet(t, mapCacher, "2-key", "found too")
	mapList(t, mapCacher, "2-key", []string{"found", "found too"})
}
