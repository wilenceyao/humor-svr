package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
)

var Config *GlobalConfig

type GlobalConfig struct {
	Mqtt MqttConfig
	// 本地管控服务
	Server  ServerConfig
	LogFile string
}

type ServerConfig struct {
	Port int
}

type QCloudConfig struct {
	SecretId  string
	SecretKey string
	TtsConfig QCloudTtsConfig
}

type MqttConfig struct {
	Ip       string
	Port     uint
	Username string
	Password string
}

type QCloudTtsConfig struct {
	VoiceType int64
}

func InitConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Error().Msgf("read config file err: %+v", err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = &GlobalConfig{}
	err = decoder.Decode(Config)
	if err != nil {
		log.Error().Msgf("decode config file err: %+v", err)
		return err
	}
	log.Info().Msgf("config init finished")
	return nil
}
