package example

import (
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/integ"
)

func mapDB() {
	config := aemkvDB.ConfigStandDb{Mode: 0}
	db := integ.New(config).(integ.MapDB)
	defer db.Close()
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}

func mapSyncDB() {
	config := aemkvDB.ConfigStandDb{Mode: 1}
	db := integ.New(config).(*integ.MapSyncDB)
	defer db.Close()
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
