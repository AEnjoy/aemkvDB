package example

import (
	"context"
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/etcd"
)

func etcdDB() {
	config := aemkvDB.ConfigMemDb{}
	config.Addr = []string{"127.0.0.1:2379"}
	config.Ctx = context.Background()
	db := etcd.New(config)
	defer db.Close()
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
