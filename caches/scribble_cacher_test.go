package caches_test

import (
	"github.com/nanopack/shaman/caches"
	"testing"
	"time"
)

func initializeScribbleCacher(t *testing.T, expires int) caches.Cacher {
	cacher, err := caches.NewScribbleCacher("scribble://localhost/tmp/shaman-test", expires)
	cacher.ClearDatabase()
	if err != nil {
		t.Errorf("Error from initializeScribbleCacher in ScribbleCacher: %s", err)
	}
	return cacher
}

func scribbleSet(t *testing.T, scribbleCacher caches.Cacher, key string, value string) {
	err := scribbleCacher.SetRecord(key, value)
	if err != nil {
		t.Errorf("Error from SetRecord in ScribbleCacher: %s", err)
	}
}

func scribbleGet(t *testing.T, scribbleCacher caches.Cacher, key string, checkValue string) {
	value, err := scribbleCacher.GetRecord(key)
	if err != nil {
		t.Errorf("Error from GetRecord in ScribbleCacher: %s", err)
	}
	if value != checkValue {
		t.Errorf("Unexpected result from ScribbleCacher: %s", value)
	}
}

func scribbleRevise(t *testing.T, scribbleCacher caches.Cacher, key string, value string) {
	err := scribbleCacher.ReviseRecord(key, value)
	if err != nil {
		t.Errorf("Error from SetRecord in ScribbleCacher: %s", err)
	}
}

func scribbleDelete(t *testing.T, scribbleCacher caches.Cacher, key string) {
	err := scribbleCacher.DeleteRecord("1-key")
	if err != nil {
		t.Errorf("Error from DeleteRecord in ScribbleCacher: %s", err)
	}
}

func scribbleList(t *testing.T, scribbleCacher caches.Cacher, key string, checkValues []string) {
	values, err := scribbleCacher.ListRecords()
	if err != nil {
		t.Errorf("Error from ListRecord in ScribbleCacher: %s", err)
	}
	if len(values) != len(checkValues) {
		t.Errorf("Unexpected length from ListRecord in ScribbleCacher: %d", len(values))
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
			t.Errorf("Unexpected values from ListRecord in ScribbleCacher: %s", values)
		}
	}
}

func TestScribbleSet(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleSet(t, scribbleCacher, "1-key", "found")
}

func TestScribbleGet(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleGet(t, scribbleCacher, "1-key", "")
}

func TestScribbleGetAfterSet(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	time.Sleep(2 * time.Second)
	scribbleGet(t, scribbleCacher, "1-key", "found")
}

func TestScribbleGetAfterSetWithExpiresNoSleep(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 1)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	scribbleGet(t, scribbleCacher, "1-key", "found")
}

func TestScribbleGetAfterSetWithExpires(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 1)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	time.Sleep(2 * time.Second)
	scribbleGet(t, scribbleCacher, "1-key", "")
}

func TestScribbleRevise(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleRevise(t, scribbleCacher, "1-key", "found")
}

func TestScribbleReviseAfterSet(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	scribbleRevise(t, scribbleCacher, "1-key", "found")
}

func TestScribbleGetAfterReviseAfterSet(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	scribbleRevise(t, scribbleCacher, "1-key", "found too")
	scribbleGet(t, scribbleCacher, "1-key", "found too")
}

func TestScribbleGetAfterReviseAfterSetWithExpires(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 1)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	scribbleRevise(t, scribbleCacher, "1-key", "found too")
	time.Sleep(2 * time.Second)
	scribbleGet(t, scribbleCacher, "1-key", "")
}

func TestScribbleDelete(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	scribbleDelete(t, scribbleCacher, "1-key")
	scribbleGet(t, scribbleCacher, "1-key", "")
}

func TestScribbleDeleteToo(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	scribbleSet(t, scribbleCacher, "2-key", "found too")
	scribbleDelete(t, scribbleCacher, "1-key")
	scribbleGet(t, scribbleCacher, "1-key", "")
	scribbleGet(t, scribbleCacher, "2-key", "found too")
}

func TestScribbleList(t *testing.T) {
	scribbleCacher := initializeScribbleCacher(t, 0)
	scribbleSet(t, scribbleCacher, "1-key", "found")
	scribbleSet(t, scribbleCacher, "2-key", "found too")
	scribbleList(t, scribbleCacher, "2-key", []string{"found", "found too"})
}
