package internal

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-svr/config"
	emq_client "github.com/wilenceyao/humor-svr/pkg/emq-client"
	"github.com/wilenceyao/humors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	stdlog "log"
)

var DefaultServer *Server
var LogFileWriter io.Writer

type Server struct {
	router   *gin.Engine
	impl     *ApiImpl
	humorSys *humors.Humors
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
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", config.Config.MqttServer.IP, config.Config.MqttServer.Port))
	opts.SetUsername(config.Config.MqttServer.Username)
	opts.SetPassword(config.Config.MqttServer.Password)
	opts.SetClientID(config.Config.MqttServer.ClientID)
	h, err := humors.NewHumors(opts)
	if err != nil {
		log.Error().Msgf("humors init err: %v", err)
		return err
	}
	h.InitAdaptor(2000)
	DefaultServer = &Server{
		humorSys: h,
	}
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = stdlog.Writer()
	DefaultServer.router = gin.Default()
	DefaultServer.impl = &ApiImpl{
		Adaptor: h.Adaptor,
	}
	DefaultServer.addApi()
	return DefaultServer.router.Run(fmt.Sprintf(":%d", config.Config.Server.Port))
}

func (s *Server) addApi() {
	s.router.POST("tts", s.impl.SendTts)
	s.router.POST("weather", s.impl.Weather)
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
