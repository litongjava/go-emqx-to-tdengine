package model

// AppConfig 结构体保存 MQTT 配置信息
type AppConfig struct {
	Broker   string
	Port     int
	Username string
	Password string
	Topic    string
}
