package buntdb

import (
	"testing"
)

// test success

func TestNewMemDB(t *testing.T) {
	db := &MemDB{}
	err := db.newMemDB()

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	if !db.IsOpened() {
		t.Errorf("Expected db to be opened, but got false")
	}
}

func TestSet(t *testing.T) {
	db := &MemDB{}
	db.open(":memory")

	err := db.Set("key", "value")

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	value := db.Get("key")
	if value != "value" {
		t.Errorf("Expected value 'value', but got %v", value)
	}
}

func TestSetWhenNotOpened(t *testing.T) {
	db := &MemDB{}

	err := db.Set("key", "value")

	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	} else if err.Error() != "Db is not opened" {
		t.Errorf("Expected error message 'Db is not opened', but got %v", err.Error())
	}
}

func TestGet(t *testing.T) {
	db := &MemDB{}
	db.open(":memory")
	db.Set("key", "value")

	value := db.Get("key")
	if value != "value" {
		t.Errorf("Expected value 'value', but got %v", value)
	}
}

func TestGetWhenNotOpened(t *testing.T) {
	db := &MemDB{}

	value := db.Get("key")
	if value != "" {
		t.Errorf("Expected empty value, but got %v", value)
	}
}

func TestDelete(t *testing.T) {
	db := &MemDB{}
	db.open(":memory")
	db.Set("key", "value")
	db.Delete("key")

	value := db.Get("key")
	if value != "" {
		t.Errorf("Expected empty value, but got %v", value)
	}
}

func TestDeleteWhenNotOpened(t *testing.T) {
	db := &MemDB{}

	err := db.Delete("key")

	if err == nil {
		t.Errorf("Expected non-nil error, but got nil")
	} else if err.Error() != "db is not opened" {
		t.Errorf("Expected error message 'db is not opened', but got %v", err.Error())
	}
}

func TestClose(t *testing.T) {
	db := &MemDB{}
	db.open(":memory")

	err := db.Close()

	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	if db.IsOpened() {
		t.Errorf("Expected db to be closed, but got true")
	}
}

func TestIsOpened(t *testing.T) {
	db := &MemDB{}
	db.open(":memory")

	if !db.IsOpened() {
		t.Errorf("Expected db to be opened, but got false")
	}
}
