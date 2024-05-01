package json

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

// test fail

func TestJsonDBOpen(t *testing.T) {
	db := &JsonDB{}
	err := db.open(0, "testdata/test.json")
	assert.NoError(t, err)
	assert.True(t, db.IsOpened())
	assert.Equal(t, "testdata/test.json", db.path)

	err = db.open(0, "testdata/test.json")
	assert.Error(t, err)
	assert.Equal(t, "Db is opened", err.Error())

	db.Close()

	err = db.open(1, "testdata/new.json")
	assert.NoError(t, err)
	assert.True(t, db.IsOpened())
	assert.Equal(t, "testdata/new.json", db.path)

	db.Close()

	err = db.open(2, `{"key":"value"}`)
	assert.NoError(t, err)
	assert.True(t, db.IsOpened())
}

func TestJsonDBSetStr(t *testing.T) {
	db := &JsonDB{}
	err := db.open(1, "testdata/new.json")
	assert.NoError(t, err)

	err = db.SetStr("key1", "value1")
	assert.NoError(t, err)

	err = db.SetStr("key2", "value2")
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBSetInt(t *testing.T) {
	db := &JsonDB{}
	err := db.open(1, "testdata/new.json")
	assert.NoError(t, err)

	err = db.SetInt("key1", strconv.Itoa(100))
	assert.NoError(t, err)

	err = db.SetInt("key2", strconv.Itoa(200))
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBSetFloat64(t *testing.T) {
	db := &JsonDB{}
	err := db.open(1, "testdata/new.json")
	assert.NoError(t, err)

	err = db.SetFloat64("key1", 1.1)
	assert.NoError(t, err)

	err = db.SetFloat64("key2", 2.2)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBSetBool(t *testing.T) {
	db := &JsonDB{}
	err := db.open(1, "testdata/new.json")
	assert.NoError(t, err)

	err = db.SetBool("key1", true)
	assert.NoError(t, err)

	err = db.SetBool("key2", false)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBGetStr(t *testing.T) {
	db := &JsonDB{}
	err := db.open(0, "testdata/new.json")
	assert.NoError(t, err)
	err = db.SetStr("key", "value")
	assert.NoError(t, err)

	//db.Close()
	//err = db.open(0, "testdata/test.json")
	t.Log(db.GetStr("key"))
	assert.Equal(t, "value", db.GetStr("key"))

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBGetInt(t *testing.T) {
	db := &JsonDB{}
	err := db.open(0, "testdata/test.json")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), db.GetInt("num"))

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBGetBool(t *testing.T) {
	db := &JsonDB{}
	err := db.open(0, "testdata/test.json")
	assert.NoError(t, err)
	assert.Equal(t, true, db.GetBool("bool"))

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBGetFloat64(t *testing.T) {
	db := &JsonDB{}
	err := db.open(0, "testdata/test.json")
	assert.NoError(t, err)
	assert.Equal(t, 3.3, db.GetFloat64("float"))

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBSaveDB(t *testing.T) {
	db := &JsonDB{}
	err := db.open(1, "testdata/new.json")
	assert.NoError(t, err)

	err = db.SetStr("key1", "value1")
	assert.NoError(t, err)

	flag := db.saveDB("testdata/new.json")
	if !flag {
		t.Error("save db error")
	}

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBSet(t *testing.T) {
	db := &JsonDB{}
	err := db.open(1, "testdata/new.json")
	assert.NoError(t, err)

	err = db.Set("key1", 100)
	assert.NoError(t, err)

	err = db.Set("key2", "value2")
	assert.NoError(t, err)

	err = db.Set("key3", 3.3)
	assert.NoError(t, err)

	err = db.Set("key4", true)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBDelete(t *testing.T) {
	db := &JsonDB{}
	err := db.open(0, "testdata/test.json")
	assert.NoError(t, err)

	err = db.Delete("key")
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBGet(t *testing.T) {
	db := &JsonDB{}
	err := db.open(0, "testdata/test.json")
	assert.NoError(t, err)

	val := db.Get("num")
	assert.Equal(t, "10", val)

	val = db.Get("key")
	assert.Equal(t, "value", val)

	err = db.Close()
	assert.NoError(t, err)
}

func TestJsonDBClose(t *testing.T) {
	db := &JsonDB{}
	err := db.open(1, "testdata/new.json")
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)

	err = db.Close()
	assert.Equal(t, "db is not opened or has been closed", err.Error())
}
