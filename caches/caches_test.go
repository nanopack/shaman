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

func TestGetRecordL1(t *testing.T) {
	caches.InitCache()
	ctrl1 := gomock.NewController(t)
	defer ctrl1.Finish()
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	L1 := mock_caches.NewMockCacher(ctrl1)
	L2 := mock_caches.NewMockCacher(ctrl2)
	caches.L1 = L1
	caches.L2 = L2
	go caches.StartCache()
	gomock.InOrder(
		L1.EXPECT().GetRecord("1-key").Return("found", nil),
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

func TestGetRecordL2(t *testing.T) {
	caches.InitCache()
	ctrl1 := gomock.NewController(t)
	defer ctrl1.Finish()
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	L1 := mock_caches.NewMockCacher(ctrl1)
	L2 := mock_caches.NewMockCacher(ctrl2)
	caches.L1 = L1
	caches.L2 = L2
	go caches.StartCache()
	gomock.InOrder(
		L1.EXPECT().GetRecord("1-key").Return("", nil),
		L2.EXPECT().GetRecord("1-key").Return("found", nil),
		L1.EXPECT().SetRecord("1-key", "found").Return(nil),
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
	caches.InitCache()
	ctrl1 := gomock.NewController(t)
	defer ctrl1.Finish()
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	L1 := mock_caches.NewMockCacher(ctrl1)
	L2 := mock_caches.NewMockCacher(ctrl2)
	caches.L1 = L1
	caches.L2 = L2
	go caches.StartCache()
	gomock.InOrder(
		L2.EXPECT().SetRecord("1-key", "found").Return(nil),
		L1.EXPECT().SetRecord("1-key", "found").Return(nil),
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
	caches.InitCache()
	ctrl1 := gomock.NewController(t)
	defer ctrl1.Finish()
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	L1 := mock_caches.NewMockCacher(ctrl1)
	L2 := mock_caches.NewMockCacher(ctrl2)
	caches.L1 = L1
	caches.L2 = L2
	go caches.StartCache()
	gomock.InOrder(
		L2.EXPECT().ReviseRecord("1-key", "found").Return(nil),
		L1.EXPECT().ReviseRecord("1-key", "found").Return(nil),
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
	caches.InitCache()
	ctrl1 := gomock.NewController(t)
	defer ctrl1.Finish()
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	L1 := mock_caches.NewMockCacher(ctrl1)
	L2 := mock_caches.NewMockCacher(ctrl2)
	caches.L1 = L1
	caches.L2 = L2
	go caches.StartCache()
	gomock.InOrder(
		L2.EXPECT().DeleteRecord("1-key").Return(nil),
		L1.EXPECT().DeleteRecord("1-key").Return(nil),
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
	caches.InitCache()
	ctrl1 := gomock.NewController(t)
	defer ctrl1.Finish()
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	L1 := mock_caches.NewMockCacher(ctrl1)
	L2 := mock_caches.NewMockCacher(ctrl2)
	caches.L1 = L1
	caches.L2 = L2
	go caches.StartCache()
	gomock.InOrder(
		L2.EXPECT().ListRecords().Return([]string{"found"}, nil),
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
