package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
)

// test success

var HostInfo = "192.168.53.53:6379"
var passWord = "123456"

func TestDbOpen(t *testing.T) {
	db := &Db{}
	err := db.open(HostInfo, passWord, 0)
	if err != nil {
		t.Errorf("open() returned error: %v", err)
	}

	if !db.IsOpened() {
		t.Errorf("open() did not set opened to true")
	}

	if db.Db == nil {
		t.Errorf("open() did not create redis client")
	}

	err = db.open("localhost:6379", "", 0)
	if err == nil || err.Error() != "db is already opened" {
		t.Errorf("open() did not return error when db is already opened")
	}
}

func TestDbOpenWithCTX(t *testing.T) {
	db := &Db{}
	ctx := context.Background()
	err := db.openWithCTX(HostInfo, passWord, 0, ctx)
	if err != nil {
		t.Errorf("openWithCTX() returned error: %v", err)
	}

	if !db.IsOpened() {
		t.Errorf("openWithCTX() did not set opened to true")
	}

	if db.Db == nil {
		t.Errorf("openWithCTX() did not create redis client")
	}

	if db.ctx != ctx {
		t.Errorf("openWithCTX() did not set context")
	}

	err = db.openWithCTX(HostInfo, passWord, 0, ctx)
	if err == nil || err.Error() != "db is already opened" {
		t.Errorf("openWithCTX() did not return error when db is already opened")
	}
}

func TestDbCheckLink(t *testing.T) {
	db := &Db{
		Db: redis.NewClient(&redis.Options{
			Addr:     HostInfo,
			Password: passWord,
			DB:       0,
		}),
		opened: true,
		ctx:    context.Background(),
	}

	if !db.checkLink() {
		t.Errorf("checkLink() returned false with a valid connection")
	}

	db.Db.Close()
	db.opened = false

	if db.checkLink() {
		t.Errorf("checkLink() returned true with an invalid connection")
	}
}

func TestDbIncr(t *testing.T) {
	db := &Db{
		Db: redis.NewClient(&redis.Options{
			Addr:     HostInfo,
			Password: passWord,
			DB:       0,
		}),
		opened: true,
		ctx:    context.Background(),
	}
	db.Set("key", "0")
	err := db.Incr("key")
	if err != nil {
		t.Errorf("Incr() returned error: %v", err)
	}

	v, err := db.Db.Get(db.ctx, "key").Result()
	if err != nil || v != "1" {
		t.Errorf("Incr() did not increment key correctly")
	}

	db.Db.Close()

	db.opened = false

	err = db.Incr("key")
	if err == nil || err.Error() != "Db is not opened" {
		t.Errorf("Incr() did not return error when db is not opened")
	}

	err = db.Incr("key")
	if err == nil {
		t.Error("Incr() did not return error when connection to Redis failed", err)
	}
}

func TestDbDecr(t *testing.T) {
	db := &Db{
		Db: redis.NewClient(&redis.Options{
			Addr:     HostInfo,
			Password: passWord,
			DB:       0,
		}),
		opened: true,
		ctx:    context.Background(),
	}
	db.Set("key", "0")
	err := db.Decr("key")
	if err != nil {
		t.Errorf("Decr() returned error: %v", err)
	}

	v, err := db.Db.Get(db.ctx, "key").Result()
	if err != nil || v != "-1" {
		t.Errorf("Decr() did not decrement key correctly")
	}

	db.Db.Close()
	db.opened = false

	err = db.Decr("key")
	if err == nil || err.Error() != "Db is not opened" {
		t.Errorf("Decr() did not return error when db is not opened")
	}

	err = db.Decr("key")
	if err == nil {
		t.Errorf("Decr() did not return error when connection to Redis failed")
	}
}

func TestDbSet(t *testing.T) {
	db := &Db{
		Db: redis.NewClient(&redis.Options{
			Addr:     HostInfo,
			Password: passWord,
			DB:       0,
		}),
		opened: true,
		ctx:    context.Background(),
	}

	err := db.Set("key", "value")
	if err != nil {
		t.Errorf("Set() returned error: %v", err)
	}

	v, err := db.Db.Get(db.ctx, "key").Result()
	if err != nil || v != "value" {
		t.Errorf("Set() did not set key correctly")
	}

	db.Db.Close()
	db.opened = false

	err = db.Set("key", "value")
	if err == nil || err.Error() != "Db is not opened" {
		t.Errorf("Set() did not return error when db is not opened")
	}

	err = db.Set("key", "value")
	if err == nil {
		t.Errorf("Set() did not return error when connection to Redis failed")
	}
}

func TestDbGet(t *testing.T) {
	db := &Db{
		Db: redis.NewClient(&redis.Options{
			Addr:     HostInfo,
			Password: passWord,
			DB:       0,
		}),
		opened: true,
		ctx:    context.Background(),
	}

	db.Db.Set(db.ctx, "key", "value", 0)

	v := db.Get("key")
	if v != "value" {
		t.Errorf("Get() did not retrieve key correctly")
	}

	db.Db.Close()
	db.opened = false

	v = db.Get("key")
	if v != "" {
		t.Errorf("Get() did not return empty string when db is not opened")
	}

	v = db.Get("key")
	if v != "" {
		t.Errorf("Get() did not return empty string when connection to Redis failed")
	}
}

func TestDbDelete(t *testing.T) {
	db := &Db{
		Db: redis.NewClient(&redis.Options{
			Addr:     HostInfo,
			Password: passWord,
			DB:       0,
		}),
		opened: true,
		ctx:    context.Background(),
	}

	db.Db.Set(db.ctx, "key", "value", 0)

	err := db.Delete("key")
	if err != nil {
		t.Errorf("Delete() returned error: %v", err)
	}

	v, err := db.Db.Get(db.ctx, "key").Result()
	if err == nil || v != "" {
		t.Errorf("Delete() did not delete key correctly")
	}

	err = db.Db.Close()
	if err != nil {
		return
	}
	db.opened = false

	err = db.Delete("key")
	if err == nil || err.Error() != "db is not opened" {
		t.Errorf("Delete() did not return error when db is not opened")
	}

	err = db.Delete("key")
	if err == nil {
		t.Errorf("Delete() did not return error when connection to Redis failed")
	}
}

func TestDbClose(t *testing.T) {
	db := &Db{
		Db: redis.NewClient(&redis.Options{
			Addr:     HostInfo,
			Password: passWord,
			DB:       0,
		}),
		opened: true,
		ctx:    context.Background(),
	}

	err := db.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	if db.opened {
		t.Errorf("Close() did not set opened to false")
	}

	err = db.Close()
	if err == nil || err.Error() != "db is not opened" {
		t.Errorf("Close() did not return error when db is not opened")
	}
}
