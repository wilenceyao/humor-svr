package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-api/config"
	emq_client "github.com/wilenceyao/humor-api/pkg/emq-client"
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
	err = emq_client.InitEmqClient(config.Config.Mqtt.Username, config.Config.Mqtt.Password,
		config.Config.Mqtt.Ip, config.Config.Mqtt.Port)
	if err != nil {
		log.Error().Msgf("InitEmqClient err: %+v", err)
		return err
	}
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.Logger
	DefaultServer = &Server{
		router: gin.Default(),
		impl:   &ApiImpl{},
	}

	DefaultServer.addApi()
	return DefaultServer.router.Run(fmt.Sprintf(":%d", config.Config.Server.Port))
}

func (s *Server) addApi() {
	s.router.POST("getDevices", s.impl.GetDevices)
	s.router.POST("sendTts", s.impl.SendTts)
}
