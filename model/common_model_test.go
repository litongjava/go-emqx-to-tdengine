package model

import (
	"github.com/spf13/viper"
	"log"
	"testing"
)

func TestLoadFile(t *testing.T) {
	viper.SetConfigFile("../config.toml")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	emqxBrokder := viper.GetString("emqxBroker")
	log.Println("emqxBrokder:", emqxBrokder)
	tdengineHost := viper.GetString("tdengineHost")
	if tdengineHost != "" {
		log.Println("connect to tdengine")
	}
}
