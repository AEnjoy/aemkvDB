package file

import (
	"testing"
)

// test success

func TestFileDB(t *testing.T) {
	db := &FileDB{}
	err := db.open("test.db")
	if err != nil {
		t.Errorf("Failed to open database: Because: %s", err.Error())
	}
	if !db.IsOpened() {
		t.Errorf("Failed to open database. Db isn't opened.")
		return
	}

	err = db.Set("key1", "value1")
	if err != nil {
		t.Errorf("Failed to set key-value pair: %s", err.Error())
	}

	value := db.Get("key1")
	if value != "value1" {
		t.Errorf("Incorrect value retrieved. Expected: value1, Got: %s", value)
	}

	err = db.Delete("key1")
	if err != nil {
		t.Errorf("Failed to delete key-value pair: %s", err.Error())
	}

	value = db.Get("key1")
	if value != "" {
		t.Errorf("Value should be empty after deletion, but got: %s", value)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Failed to close database: %s", err.Error())
	}
}
