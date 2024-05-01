package etcd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test success

func TestEtcdDB(t *testing.T) {
	// 创建EtcdDB实例
	db := &EtcdDB{}

	assert.False(t, db.IsOpened(), "db should not be opened initially")

	t.Log("Testing open method")
	addresses := []string{"192.168.53.53:2379"}
	err := db.open(addresses)
	assert.NoError(t, err, "open should not return error")
	assert.True(t, db.IsOpened(), "db should be opened after open method")

	t.Log("Testing openWithCTX method")
	ctx := context.Background()
	err = db.openWithCTX(addresses, ctx)
	assert.EqualError(t, err, "db already opened", "openWithCTX should return error when db is already opened")

	t.Log("Testing Set method")
	key := "test_key"
	value := "test_value"
	err = db.Set(key, value)
	assert.NoError(t, err, "Set should not return error")
	assert.Equal(t, value, db.Get(key), "Get should return the value set by Set method")

	t.Log("Testing Get method")
	assert.Equal(t, "", db.Get("non_existent_key"), "Get should return empty string for non-existent key")

	t.Log("Testing Delete method")
	err = db.Delete(key)
	assert.NoError(t, err, "Delete should not return error")
	assert.Equal(t, "", db.Get(key), "Get should return empty string after key is deleted")

	t.Log("Testing Close method")
	err = db.Close()
	assert.NoError(t, err, "Close should not return error")
	assert.False(t, db.IsOpened(), "db should not be opened after Close method")
}

func TestEtcdDB_ErrorConditions(t *testing.T) {
	// when db is not opened
	db := &EtcdDB{}

	assert.False(t, db.IsOpened(), "db should not be opened initially")

	err := db.Set("test_key", "test_value")
	assert.EqualError(t, err, "db is not opened", "Set should return error when db is not opened")

	assert.Equal(t, "", db.Get("test_key"), "Get should return empty string when db is not opened")

	err = db.Delete("test_key")
	assert.EqualError(t, err, "db is not opened", "Delete should return error when db is not opened")

	err = db.Close()
	assert.EqualError(t, err, "db is not opened or has been closed.", "Close should return error when db is not opened")

	err = db.Close()
	assert.EqualError(t, err, "db is not opened or has been closed.", "Close should return error when db is already closed")
}

func TestEtcdDB_ClientMethods(t *testing.T) {
	// 创建EtcdDB实例
	db := &EtcdDB{}

	// 测试open方法
	addresses := []string{"localhost:2379"}
	err := db.open(addresses)
	assert.NoError(t, err, "open should not return error")

	// 测试Set方法
	key := "test_key"
	value := "test_value"
	err = db.Set(key, value)
	assert.NoError(t, err, "Set should not return error")

	// 测试Get方法
	get, err := db.kv.Get(db.ctx, key)
	assert.NoError(t, err, "kv.Get should not return error")
	assert.Equal(t, value, string(get.Kvs[0].Value), "kv.Get should return the value set by Set method")

	// 测试Delete方法
	err = db.Delete(key)
	assert.NoError(t, err, "Delete should not return error")

	// 测试Close方法
	err = db.Close()
	assert.NoError(t, err, "Close should not return error")
}
