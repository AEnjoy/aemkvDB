package example

import (
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/file"
)

func fileDB() {
	config := aemkvDB.ConfigStandDb{Data: "a.db"}
	db := file.New(config)
	defer db.Close()
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
