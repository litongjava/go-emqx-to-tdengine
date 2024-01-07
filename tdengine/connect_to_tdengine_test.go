package tdengine

import (
	"database/sql"
	"fmt"
	_ "github.com/taosdata/driver-go/v3/taosSql"
	"log"
	"testing"
)

func TestConnectToTdEngine(t *testing.T) {
	var taosUri = "root:taosdata@tcp(192.168.3.9:6030)/"
	taos, err := sql.Open("taosSql", taosUri)
	if err != nil {
		fmt.Println("failed to connect TDengine, err:", err)
		return
	}
	log.Println(taos)
	defer taos.Close()
	rows, err := taos.Query("show databases")

	if err != nil {
		log.Fatalln("Faild to execute query")
	}
	defer rows.Close()
	for rows.Next() {
		var r struct {
			name string
		}
		err = rows.Scan(&r.name)
		if err != nil {
			log.Fatalln("scan error:\n", err)
			return
		}
		log.Println(r.name)
	}
}
