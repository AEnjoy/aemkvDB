package sql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// test success

func TestSqlDB_Set(t *testing.T) {
	db := &SqlDB{}
	err := db.open(1, "test.db")
	assert.NoError(t, err)
	assert.True(t, db.IsOpened())

	err = db.Set("key1", "value1")
	assert.NoError(t, err)

	err = db.Set("key2", "value2")
	assert.NoError(t, err)

	// Add more test cases as needed
}

func TestSqlDB_Get(t *testing.T) {
	db := &SqlDB{}
	err := db.open(1, "test.db")
	assert.NoError(t, err)
	assert.True(t, db.IsOpened())

	err = db.Set("key1", "value1")
	assert.NoError(t, err)

	value := db.Get("key1")
	assert.Equal(t, "value1", value)

	value = db.Get("key:?")
	assert.Equal(t, "", value)

	// Add more test cases as needed
}

func TestSqlDB_Delete(t *testing.T) {
	db := &SqlDB{}
	err := db.open(1, "test.db")
	assert.NoError(t, err)
	assert.True(t, db.IsOpened())

	err = db.Set("key1", "value1")
	assert.NoError(t, err)

	err = db.Delete("key1")
	assert.NoError(t, err)

	value := db.Get("key1")
	assert.Equal(t, "", value)

	// Add more test cases as needed
}

func TestSqlDB_Close(t *testing.T) {
	db := &SqlDB{}
	err := db.open(1, "test.db")
	assert.NoError(t, err)
	assert.True(t, db.IsOpened())

	err = db.Close()
	assert.NoError(t, err)
	assert.False(t, db.IsOpened())

	// Add more test cases as needed
}
