package example

import (
	"github.com/aenjoy/aemkvDB"
	"github.com/aenjoy/aemkvDB/backends/json"
)

func jsonDB() {
	config := aemkvDB.ConfigStandDb{Data: "./data.json", Mode: 0} //if file not exist, use mode 0; otherwise use mode 1.
	// if data is json string, use mode 2
	// config := aemkvDB.ConfigStandDb{Data: "{}",Mode: 2}
	db := json.New(config)
	defer db.Close()
	db.Set("key", "value")
	db.Set("key2", "value2")
	print(db.Get("key"), db.Get("key2")) //value value2
}
