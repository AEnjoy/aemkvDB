package example

import "github.com/aenjoy/aemkvDB"
import "github.com/aenjoy/aemkvDB/backends/buntdb"

func buntDb() {
	// using buntdb backend (store in memory)
	config := aemkvDB.ConfigMemDb{}
	config.Database = ":memory:"
	db := buntdb.New(config)
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
