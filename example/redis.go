package example

import (
	"context"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/redis"
)

func redisDB() {
	config := aemkvDB.ConfigMemDb{}
	config.Addr = []string{"127.0.0.1:6379"}
	config.Ctx = context.Background()
	config.Password = "123456"
	config.Database = "0"
	db := redis.New(config)
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
