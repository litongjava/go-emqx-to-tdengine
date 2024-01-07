package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"go-emqx-to-tdengine/eqmx"
	"go-emqx-to-tdengine/tdengine"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v \n", err)
}

func main() {
	configFile := flag.String("c", "config.toml", "Configuration file")
	flag.Parse()

	if err := loadConfig(*configFile); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	broker := viper.GetString("emqxBroker")
	port := viper.GetInt("emqxPort")
	username := viper.GetString("emqxUsername")
	password := viper.GetString("emqxPassword")
	topic := viper.GetString("emqxTopic")
	tdDriverName := viper.GetString("tdDriverName")
	tdDataSourceName := viper.GetString("tdDataSourceName")
	createSql := viper.GetString("createSql")
	insertSql := viper.GetString("insertSql")

	sqlDb := tdengine.ConnectTdEngine(tdDriverName, tdDataSourceName)
	defer sqlDb.Close()
	tdengine.CreateTable(sqlDb, createSql)

	sqlTemplate := getInsertSqlTemplate(insertSql)
	messagePubHandler := BuildMessageHandler(sqlTemplate, sqlDb, tdDriverName, tdDataSourceName)
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log.Println("Connected to mq")
		sub(client, topic)
	}
	eqmx.ConnectEqmx(broker, port, username, password, messagePubHandler, connectHandler, connectLostHandler)
	//publish(client, topic)
	// 阻止程序退出
	select {}
}

func getInsertSqlTemplate(insertSql string) string {
	sqlTemplate := ""
	if insertSql == "" {
		log.Println("insert sql is empty")
	} else {
		//读取createSql的值
		bytes, err := ioutil.ReadFile(insertSql)
		if err != nil {
			log.Fatalln("An error occurred while reading the file:", err)
		}
		//create_sql_test.go
		sqlTemplate = string(bytes)
		sqlTemplate = strings.ReplaceAll(sqlTemplate, "\n", " ") // Remove newlines from SQL template
		sqlTemplate = strings.ReplaceAll(sqlTemplate, "\t", " ") // Replace tabs with spaces if any
		sqlTemplate = strings.ReplaceAll(sqlTemplate, "\r", " ") // Replace \r with spaces if any
		sqlTemplate = strings.ReplaceAll(sqlTemplate, "  ", " ")
		sqlTemplate = strings.ReplaceAll(sqlTemplate, "  ", " ")
		sqlTemplate = strings.ReplaceAll(sqlTemplate, "  ", " ")
		sqlTemplate = strings.TrimSpace(sqlTemplate)
	}
	return sqlTemplate
}

//func publish(client mqtt.Client, topic string) {
//	num := 10
//	for i := 0; i < num; i++ {
//		text := fmt.Sprintf("Message %d", i)
//		token := client.Publish(topic, 0, false, text)
//		token.Wait()
//		time.Sleep(time.Second)
//	}
//}

func loadConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")
	return viper.ReadInConfig()
}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic: %s \n", topic)
}

func reconnect(tdDriverName string, tdDataSourceName string) (*sql.DB, error) {
	// 尝试重新建立连接
	db, err := sql.Open(tdDriverName, tdDataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func BuildMessageHandler(sqlTemplate string, sqlDb *sql.DB, tdDriverName string, tdDataSourceName string) mqtt.MessageHandler {
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		var payload map[string]interface{}
		err := json.Unmarshal(msg.Payload(), &payload)
		if err != nil {
			log.Println("Error in parsing JSON:", err)
		}

		// Replace placeholders in the SQL template
		if sqlTemplate != "" {
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
			// 使用示例
			_, err = sqlDb.Exec(sqlString)
			if err != nil {
				log.Println("error", sqlString)
				maxRetries := 10
				for i := 0; i < maxRetries; i++ {
					newDb, err := reconnect(tdDriverName, tdDataSourceName)
					if err != nil {
						log.Println("error in reconnect")
					} else {
						sqlDb = newDb // 更新外部的 sqlDb 变量
						_, err = sqlDb.Exec(sqlString)
						if err == nil {
							break // 如果执行成功则跳出循环
						}
						log.Println("error in reconnect success and exec sql")
					}
				}
			}
		}
	}
	return messagePubHandler
}
