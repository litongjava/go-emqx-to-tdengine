package tdengine

import (
	"database/sql"
	_ "github.com/taosdata/driver-go/v3/taosSql"
	"io/ioutil"
	"log"
)

func ConnectTdEngine(tdDriverName string, tdDataSourceName string) *sql.DB {
	log.Println(tdDriverName, tdDataSourceName)
	taos, err := sql.Open(tdDriverName, tdDataSourceName)
	if err != nil {
		log.Fatalln("An error occurred while connect to td-engine:", err)
	}
	return taos
}

func CreateTable(sqlDb *sql.DB, createSql string) {
	if createSql == "" {
		log.Println("skip create sql")
		return
	}
	//读取createSql的值
	bytes, err := ioutil.ReadFile(createSql)
	if err != nil {
		log.Fatalln("An error occurred while reading the file:", err)
	}

	sql := string(bytes)
	_, err = sqlDb.Exec(sql)
	if err != nil {
		log.Fatalln("error in create table", err)
	} else {
		log.Println("create sql exec finished")
	}
}
