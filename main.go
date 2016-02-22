package main

import (
	"flag"
	"runtime"
	"runtime/debug"
)

var configPath string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(300)
	flag.StringVar(&configPath, "conf", "./config.json", "path of the file config.json")
}

/*

SELECT * FROM `information_schema`.`TABLES` Where TABLE_SCHEMA='Customers';
SELECT * FROM `information_schema`.`VIEWS` Where TABLE_SCHEMA='Customers';

SELECT * FROM `information_schema`.`COLUMNS` Where TABLE_SCHEMA='Customers' AND TABLE_NAME='Customer';

SELECT * FROM `information_schema`.`KEY_COLUMN_USAGE` Where TABLE_SCHEMA='Customers' AND TABLE_NAME='Customer' AND CONSTRAINT_NAME='PRIMARY';

SELECT C.TABLE_NAME, C.COLUMN_NAME, C.IS_NULLABLE, C.DATA_TYPE, C.CHARACTER_MAXIMUM_LENGTH, C.NUMERIC_PRECISION, C.NUMERIC_SCALE, C.COLUMN_TYPE, C.COLUMN_KEY
FROM information_schema.COLUMNS C WHERE C.TABLE_SCHEMA = 'Messages' ORDER BY TABLE_NAME, C.ORDINAL_POSITION;

*/

// Start the code generator
func main() {

}
