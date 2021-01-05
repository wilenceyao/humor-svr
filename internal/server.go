package internal

import (
	"fmt"
	"github.com/eclipse/paho.golang/paho"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-api/config"
	emq_client "github.com/wilenceyao/humor-api/pkg/emq-client"
	"github.com/wilenceyao/humor-api/pkg/util"
	"gopkg.in/natefinch/lumberjack.v2"
)

var DefaultServer *Server

type Server struct {
	router *gin.Engine
	impl   *ApiImpl
}

func RunServer() error {
	err := config.InitConfig("config.json")
	if err != nil {
		log.Error().Msgf("InitConfig err: %+v", err)
		return err
	}
	setupLog()
	err = emq_client.InitEmqClient(config.Config.MqttAdmin.Username, config.Config.MqttAdmin.Password,
		config.Config.MqttAdmin.IP, config.Config.MqttAdmin.Port)
	if err != nil {
		log.Error().Msgf("InitEmqClient err: %+v", err)
		return err
	}
	DefaultServer = &Server{}
	rpcConfig := &util.MqttRpcConfig{
		ClientID:    config.Config.ClientID,
		IP:          config.Config.MqttServer.IP,
		Port:        config.Config.MqttServer.Port,
		Username:    config.Config.MqttServer.Username,
		Password:    config.Config.MqttServer.Password,
		RecvHandler: DefaultServer.rpcDispatcher,
	}
	mqttRpcReqHandler, err := util.NewMqttRpcHandler(rpcConfig)
	if err != nil {
		log.Error().Msgf("InitMqttRpc err: %+v", err)
		return err
	}
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Logger
	DefaultServer.router = gin.Default()
	DefaultServer.impl = &ApiImpl{
		MqttRpcReqHandler: mqttRpcReqHandler,
	}
	DefaultServer.addApi()
	return DefaultServer.router.Run(fmt.Sprintf(":%d", config.Config.Server.Port))
}
func (s *Server) rpcDispatcher(m *paho.Publish) {
}

func (s *Server) addApi() {
	s.router.POST("getDevices", s.impl.GetDevices)
	s.router.POST("sendTts", s.impl.SendTts)
}

func setupLog() {
	h := &lumberjack.Logger{
		Filename:   config.Config.LogFile,
		MaxSize:    100,  // megabytes
		MaxBackups: 10,   // 最多50个日志文件，因而只保留49个旧日志备份
		MaxAge:     10,   //days
		Compress:   true, // disabled by default
	}
	log.Logger = log.With().Caller().Logger().Output(h)
}
