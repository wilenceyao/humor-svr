package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-api/config"
	emq_client "github.com/wilenceyao/humor-api/pkg/emq-client"
	"github.com/wilenceyao/humor-api/pkg/mqttrpc"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	stdlog "log"
)

var DefaultServer *Server
var LogFileWriter io.Writer

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
	rpcConfig := &mqttrpc.MqttRpcConfig{
		IP:       config.Config.MqttServer.IP,
		Port:     config.Config.MqttServer.Port,
		Username: config.Config.MqttServer.Username,
		Password: config.Config.MqttServer.Password,
		RpcTopic: mqttrpc.HUMOR_API_RPC_TOPIC,
	}
	err = mqttrpc.InitMqttClient(rpcConfig)
	if err != nil {
		log.Error().Msgf("InitMqttClient err: %+v", err)
		return err
	}
	gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = LogFileWriter
	DefaultServer.router = gin.Default()
	DefaultServer.impl = &ApiImpl{
	}
	DefaultServer.addApi()
	return DefaultServer.router.Run(fmt.Sprintf(":%d", config.Config.Server.Port))
}

func (s *Server) addApi() {
	s.router.POST("getDevices", s.impl.GetDevices)
	s.router.POST("sendTts", s.impl.SendTts)
}

func setupLog() {
	LogFileWriter = &lumberjack.Logger{
		Filename:   config.Config.LogFile,
		MaxSize:    100,  // megabytes
		MaxBackups: 10,   // 最多50个日志文件，因而只保留49个旧日志备份
		MaxAge:     10,   //days
		Compress:   true, // disabled by default
	}
	stdlog.SetOutput(LogFileWriter)
	log.Logger = log.With().Caller().Logger().Output(LogFileWriter)
}
