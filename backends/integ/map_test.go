package integ

import (
	"errors"
	"testing"
)

// test success

func TestMapDBOpen(t *testing.T) {
	db := &MapDB{opened: false}

	err := db.open()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if !db.opened {
		t.Errorf("Expected db.opened to be true, but got false")
	}

	if len(db.Db) != 0 {
		t.Errorf("Expected db.Db to be empty, but got length: %d", len(db.Db))
	}

	err = db.open()

	if err == nil || err.Error() != "db is opened" {
		t.Errorf("Expected error: 'db is opened', but got: %v", err)
	}
}

func TestMapDBSet(t *testing.T) {
	db := &MapDB{}
	err := db.open()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	err = db.Set("key", "value")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if db.Db["key"] != "value" {
		t.Errorf("Expected db.Db['key'] to be 'value', but got: %v", db.Db["key"])
	}

	err = db.Set("key", "new value")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if db.Db["key"] != "new value" {
		t.Errorf("Expected db.Db['key'] to be 'new value', but got: %v", db.Db["key"])
	}

	err = db.Close()
	if err != nil {
		return
	}

	err = db.Set("key", "value")

	if err == nil || err.Error() != "db is not opened" {
		t.Errorf("Expected error: 'db is not opened', but got: %v", err)
	}
}

func TestMapDBGet(t *testing.T) {
	db := &MapDB{}
	err := db.open()
	if err != nil {
		return
	}
	db.Db["key"] = "value"

	value := db.Get("key")

	if value == nil || value != "value" {
		t.Errorf("Expected db.Get('key') to be 'value', but got: %v", value)
	}

	err = db.Close()
	if err != nil {
		return
	}

	value = db.Get("key")

	if value == nil {
		t.Errorf("Expected db.Get('key') to return error: 'db is not opened', but got: %v", value)
	}
}

func TestMapDBDelete(t *testing.T) {
	db := &MapDB{}
	err := db.open()
	if err != nil {
		return
	}
	db.Db["key"] = "value"

	err = db.Delete("key")

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if db.Db["key"] != nil {
		t.Errorf("Expected db.Db['key'] to be nil, but got: %v", db.Db["key"])
	}

	db.opened = false

	err = db.Delete("key")

	if err == nil || err.Error() != "db is not opened" {
		t.Errorf("Expected error: 'db is not opened', but got: %v", err)
	}
}

func TestMapDBClose(t *testing.T) {
	db := &MapDB{}
	err := db.open()
	if err != nil {
		return
	}
	err = db.Close()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if db.opened {
		t.Errorf("Expected db.opened to be false, but got true")
	}

	if len(db.Db) != 0 {
		t.Errorf("Expected db.Db to be empty, but got length: %d", len(db.Db))
	}

}

func TestMapDBIsOpened(t *testing.T) {
	db := &MapDB{}
	err := db.open()
	if err != nil {
		return
	}
	if !db.IsOpened() {
		t.Errorf("Expected db.IsOpened() to be true, but got false")
	}

	db.opened = false

	if db.IsOpened() {
		t.Errorf("Expected db.IsOpened() to be false, but got true")
	}
}

func TestMapSyncDBOpen(t *testing.T) {
	db := &MapSyncDB{}
	err := db.open()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	if !db.IsOpened() {
		t.Errorf("Expected db to be opened, but got false")
	}
}

func TestMapSyncDBOpenAlreadyOpened(t *testing.T) {
	db := &MapSyncDB{opened: true}
	err := db.open()
	expectedErr := errors.New("db is already opened")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error \"%v\", but got \"%v\"", expectedErr, err)
	}
	if !db.IsOpened() {
		t.Errorf("Expected db to be opened, but got false")
	}
}

func TestMapSyncDBSet(t *testing.T) {
	db := &MapSyncDB{}
	err := db.open()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	err = db.Set("key", "value")
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	value := db.Get("key")
	if value == nil || value != "value" {
		t.Errorf("Expected value to be \"value\", but got %v", value)
	}
}

func TestMapSyncDBSetNotOpened(t *testing.T) {
	db := &MapSyncDB{}
	err := db.Set("key", "value")
	expectedErr := errors.New("db is not opened")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error \"%v\", but got \"%v\"", expectedErr, err)
	}
}

func TestMapSyncDBGet(t *testing.T) {
	db := &MapSyncDB{}
	err := db.open()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	db.Set("key", "value")
	value := db.Get("key")
	if value == nil || value != "value" {
		t.Errorf("Expected value to be \"value\", but got %v", value)
	}
}

func TestMapSyncDBGetNotOpened(t *testing.T) {
	db := &MapSyncDB{}
	value := db.Get("key")
	if value != nil {
		t.Errorf("Expected value to be nil, but got %v", value)
	}
}

func TestMapSyncDBDelete(t *testing.T) {
	db := &MapSyncDB{}
	err := db.open()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	db.Set("key", "value")
	err = db.Delete("key")
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	value := db.Get("key")
	if value != nil {
		t.Errorf("Expected value to be nil, but got %v", value)
	}
}

func TestMapSyncDBDeleteNotOpened(t *testing.T) {
	db := &MapSyncDB{}
	err := db.Delete("key")
	expectedErr := errors.New("db is not opened")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error \"%v\", but got \"%v\"", expectedErr, err)
	}
}

func TestMapSyncDBClose(t *testing.T) {
	db := &MapSyncDB{}
	err := db.open()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	err = db.Close()
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	if db.IsOpened() {
		t.Errorf("Expected db to be closed, but got true")
	}
}

func TestMapSyncDBCloseNotOpened(t *testing.T) {
	db := &MapSyncDB{}
	err := db.Close()
	expectedErr := errors.New("db is not opened")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Errorf("Expected error \"%v\", but got \"%v\"", expectedErr, err)
	}
}
