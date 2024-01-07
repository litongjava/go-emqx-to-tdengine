package tdengine

import (
	"database/sql"
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"testing"
)

func TestInsertSql(t *testing.T) {
	viper.SetConfigFile("../config.toml")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("error in read file", err)
	}
	tdDriverName := viper.GetString("tdDriverName")
	tdDataSourceName := viper.GetString("tdDataSourceName")
	//createSql := viper.GetString("createSql")
	createSql := "../insert.sql"
	//createSql := ""
	if createSql == "" {
		log.Println("skip insert sqlString")
		return
	}
	//读取createSql的值
	bytes, err := ioutil.ReadFile(createSql)
	if err != nil {
		log.Fatalln("An error occurred while reading the file:", err)
	}
	//create_sql_test.go
	sqlTemplate := string(bytes)
	payloadString := "{\"temperature\":76.71,\"humidity\":43.62,\"volume\":139.53,\"PM10\":152.52,\"pm25\":138.79,\"SO2\":49.08,\"NO2\":17.5,\"CO\":25.38,\"area\":2,\"ts\":1704572289024,\"id\":\"mock_client_91\"}"
	//_, err = taos.Exec(sqlString)
	//if err != nil {

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(payloadString), &payload)
	if err != nil {
		log.Fatalln("Error in parsing JSON:", err)
	}

	// Replace placeholders in the SQL template
	sqlString := sqlTemplate
	for key, value := range payload {
		placeholder := "${payload." + key + "}"
		var valueStr string
		switch v := value.(type) {
		case float64:
			valueStr = strconv.FormatFloat(v, 'f', -1, 64)
		case string:
			valueStr = v
		case int, int64:
			valueStr = strconv.FormatInt(v.(int64), 10) // Convert to int64 for uniformity
		}
		sqlString = strings.Replace(sqlString, placeholder, valueStr, -1)
	}

	println(sqlString)

	log.Println(tdDriverName, tdDataSourceName)
	taos, err := sql.Open(tdDriverName, tdDataSourceName)
	defer taos.Close()
	if err != nil {
		log.Fatalln("An error occurred while connect to td-engine:", err)
	}
	_, err = taos.Exec(sqlString)
	if err != nil {
		log.Fatalln("An error occurred while insert data", err)
	}
}
