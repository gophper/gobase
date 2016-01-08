package gobase

import (
	// "fmt"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	CacheTable = NewCache()
	CacheTable.Set("setKey1", "setValue1", 3)
	CacheTable.Set("setKey2", "setValue2", 0)
	CacheValue1, _ := CacheTable.Get("setKey1")
	CacheValue2, _ := CacheTable.Get("setKey2")
	if CacheValue1 != "setValue1" {
		t.Error("CacheTable Set Error")
	}
	if CacheValue2 != "setValue2" {
		t.Error("CacheTable Set Error")
	}

	time.Sleep(5 * time.Second)
	CacheValue1, _ = CacheTable.Get("setKey1")
	CacheValue2, _ = CacheTable.Get("setKey2")
	if CacheValue1 == "setValue1" {
		t.Error("5 * time.Second CacheTable Set Error")
	}
	if CacheValue2 != "setValue2" {
		t.Error("5 * time.Second CacheTable Set Error")
	}

}
