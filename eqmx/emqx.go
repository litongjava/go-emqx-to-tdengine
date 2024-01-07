package eqmx

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func ConnectEqmx(broker string, port int, username string, password string,
	messagePubHandler mqtt.MessageHandler, connectHandler mqtt.OnConnectHandler, connectLostHandler mqtt.ConnectionLostHandler) mqtt.Client {
	opts := mqtt.NewClientOptions()

	log.Println("broker port", broker, port)
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")

	log.Println("username password", username, password)
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
