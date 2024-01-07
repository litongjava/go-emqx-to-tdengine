package tdengine

import (
	"database/sql"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"testing"
)

func TestCrateSql(t *testing.T) {
	viper.SetConfigFile("../config.toml")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("error in read file", err)
	}
	tdDriverName := viper.GetString("tdDriverName")
	tdDataSourceName := viper.GetString("tdDataSourceName")
	//createSql := viper.GetString("createSql")
	//createSql := "../create.sql"
	createSql := ""
	if createSql == "" {
		log.Println("skip create sql")
		return
	}
	log.Println(tdDriverName, tdDataSourceName)
	taos, err := sql.Open(tdDriverName, tdDataSourceName)
	defer taos.Close()
	if err != nil {
		log.Fatalln("error in connect to td-engine:", err)
	}
	//读取createSql的值
	bytes, err := ioutil.ReadFile(createSql)
	if err != nil {
		log.Fatalln("An error occurred while reading the file:", err)
	}

	sql := string(bytes)
	_, err = taos.Exec(sql)
	if err != nil {
		log.Fatalln("error in create table", err)
	}
}
