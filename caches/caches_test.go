package caches_test

import (
	"github.com/golang/mock/gomock"
	"github.com/jcelliott/lumber"
	"github.com/nanopack/shaman/caches"
	"github.com/nanopack/shaman/caches/mock_caches"
	"github.com/nanopack/shaman/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.Log = lumber.NewConsoleLogger(lumber.ERROR)
	if testing.Verbose() {
		config.Log = lumber.NewConsoleLogger(lumber.DEBUG)
	}
	os.Exit(m.Run())
}

func initializeCaches(t *testing.T) (*mock_caches.MockCacher, *mock_caches.MockCacher) {
	caches.InitCache()
	ctrl1 := gomock.NewController(t)
	defer ctrl1.Finish()
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	l1 := mock_caches.NewMockCacher(ctrl1)
	l2 := mock_caches.NewMockCacher(ctrl2)
	caches.L1 = l1
	caches.L2 = l2
	go caches.StartCache()
	return l1, l2
}

func TestFindRecordL1(t *testing.T) {
	l1, _ := initializeCaches(t)
	gomock.InOrder(
		l1.EXPECT().GetRecord("1-key").Return("found", nil),
	)
	findReturn := make(chan caches.FindReturn)
	findOp := caches.FindOp{Key: "1-key", Resp: findReturn}
	caches.FindOps <- findOp
	findRet := <-findReturn
	err := findRet.Err
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	record := findRet.Value
	if record != "found" {
		t.Errorf("bad result from L1: %s", record)
	}
}

func TestFindRecordL2(t *testing.T) {
	l1, l2 := initializeCaches(t)
	gomock.InOrder(
		l1.EXPECT().GetRecord("1-key").Return("", nil),
		l2.EXPECT().GetRecord("1-key").Return("found", nil),
		l1.EXPECT().SetRecord("1-key", "found").Return(nil),
	)
	findReturn := make(chan caches.FindReturn)
	findOp := caches.FindOp{Key: "1-key", Resp: findReturn}
	caches.FindOps <- findOp
	findRet := <-findReturn
	err := findRet.Err
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	record := findRet.Value
	if record != "found" {
		t.Errorf("bad result from L1: %s", record)
	}
}

func TestAddRecord(t *testing.T) {
	l1, l2 := initializeCaches(t)
	gomock.InOrder(
		l2.EXPECT().SetRecord("1-key", "found").Return(nil),
		l1.EXPECT().SetRecord("1-key", "found").Return(nil),
	)
	resp := make(chan error)
	addOp := caches.AddOp{Key: "1-key", Value: "found", Resp: resp}
	caches.AddOps <- addOp
	err := <-resp
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestUpdateRecord(t *testing.T) {
	l1, l2 := initializeCaches(t)
	gomock.InOrder(
		l2.EXPECT().ReviseRecord("1-key", "found").Return(nil),
		l1.EXPECT().ReviseRecord("1-key", "found").Return(nil),
	)
	resp := make(chan error)
	updateOp := caches.UpdateOp{Key: "1-key", Value: "found", Resp: resp}
	caches.UpdateOps <- updateOp
	err := <-resp
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestRemoveRecord(t *testing.T) {
	l1, l2 := initializeCaches(t)
	gomock.InOrder(
		l2.EXPECT().DeleteRecord("1-key").Return(nil),
		l1.EXPECT().DeleteRecord("1-key").Return(nil),
	)
	resp := make(chan error)
	removeOp := caches.RemoveOp{Key: "1-key", Resp: resp}
	caches.RemoveOps <- removeOp
	err := <-resp
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestListRecords(t *testing.T) {
	_, l2 := initializeCaches(t)
	gomock.InOrder(
		l2.EXPECT().ListRecords().Return([]string{"found"}, nil),
	)
	resp := make(chan caches.ListReturn)
	listOp := caches.ListOp{Resp: resp}
	caches.ListOps <- listOp
	listReturn := <-resp
	err := listReturn.Err
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(listReturn.Values) != 1 && listReturn.Values[0] != "found" {
		t.Errorf("Bad return: %s", listReturn.Values)
	}
}
