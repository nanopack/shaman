package shaman_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jcelliott/lumber"

	"github.com/nanopack/shaman/config"
	"github.com/nanopack/shaman/core"
	sham "github.com/nanopack/shaman/core/common"
)

var (
	nanopack  = sham.Resource{Domain: "nanopack.io.", Records: []sham.Record{{Address: "127.0.0.1"}}}
	nanopack2 = sham.Resource{Domain: "nanopack.io.", Records: []sham.Record{{Address: "127.0.0.3"}}}
	nanobox   = sham.Resource{Domain: "nanobox.io.", Records: []sham.Record{{Address: "127.0.0.2"}}}
	nanoBoth  = []sham.Resource{nanopack, nanobox}
)

func TestMain(m *testing.M) {
	shamanClear()
	// manually configure
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt("FATAL"))

	// run tests
	rtn := m.Run()

	os.Exit(rtn)
}

func TestAddRecord(t *testing.T) {
	shamanClear()
	err := shaman.AddRecord(&nanopack)
	err = shaman.AddRecord(&nanopack)
	err2 := shaman.AddRecord(&nanopack2)
	if err != nil || err2 != nil {
		t.Errorf("Failed to add record - %v%v", err, err2)
	}
}

func TestGetRecord(t *testing.T) {
	shamanClear()
	_, err := shaman.GetRecord("nanopack.io")
	shaman.AddRecord(&nanopack)
	_, err2 := shaman.GetRecord("nanopack.io")
	if err == nil || err2 != nil {
		// t.Errorf("Failed to get record - %v%v", err, "hi")
		t.Errorf("Failed to get record - %v%v", err, err2)
	}
}

func TestUpdateRecord(t *testing.T) {
	shamanClear()
	err := shaman.UpdateRecord("nanopack.io", &nanopack)
	err2 := shaman.UpdateRecord("nanobox.io", &nanopack)
	if err != nil || err2 != nil {
		t.Errorf("Failed to update record - %v%v", err, err2)
	}
}

func TestDeleteRecord(t *testing.T) {
	shamanClear()
	err := shaman.DeleteRecord("nanobox.io")
	shaman.AddRecord(&nanopack)
	err2 := shaman.DeleteRecord("nanopack.io")
	if err != nil || err2 != nil {
		t.Errorf("Failed to delete record - %v%v", err, err2)
	}
}

func TestResetRecords(t *testing.T) {
	shamanClear()
	err := shaman.ResetRecords(&nanoBoth)
	err2 := shaman.ResetRecords(&nanoBoth, true)
	if err != nil || err2 != nil {
		t.Errorf("Failed to reset records - %v%v", err, err2)
	}
}

func TestListDomains(t *testing.T) {
	shamanClear()
	domains := shaman.ListDomains()
	if fmt.Sprint(domains) != "[]" {
		t.Errorf("Failed to list domains - %+q", domains)
	}
	shaman.ResetRecords(&nanoBoth)
	domains = shaman.ListDomains()
	if len(domains) != 2 {
		t.Errorf("Failed to list domains - %+q", domains)
	}
}

func TestListRecords(t *testing.T) {
	shamanClear()
	resources := shaman.ListRecords()
	if fmt.Sprint(resources) != "[]" {
		t.Errorf("Failed to list records - %+q", resources)
	}
	shaman.ResetRecords(&nanoBoth)
	resources = shaman.ListRecords()
	if len(resources) == 2 && (resources[0].Domain != "nanopack.io." && resources[0].Domain != "nanobox.io.") {
		t.Errorf("Failed to list records - %+q", resources)
	}
}

func TestExists(t *testing.T) {
	shamanClear()
	if shaman.Exists("nanopack.io") {
		t.Errorf("Failed to list records")
	}
	shaman.AddRecord(&nanopack)
	if !shaman.Exists("nanopack.io") {
		t.Errorf("Failed to list records")
	}
}

func shamanClear() {
	shaman.Answers = make(map[string]sham.Resource, 0)
}
