package util

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/paho"
	"github.com/eclipse/paho.golang/paho/extensions/rpc"
	"github.com/rs/zerolog/log"
	"net"
)

type MqttRpcConfig struct {
	ClientID    string
	IP          string
	Port        uint
	Username    string
	Password    string
	RecvHandler paho.MessageHandler
}

func FormatRpcTopic(deviceID string) string {
	return fmt.Sprintf("rpc/%s", deviceID)
}

func NewMqttRpcHandler(rpcCfg *MqttRpcConfig) (*rpc.Handler, error) {
	log.Info().Msg("NewMqttRpcHandler start")
	server := fmt.Sprintf("%s:%d", rpcCfg.IP, rpcCfg.Port)
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Error().Msgf("connect to mqtt server %s err: %+v", server, err)
		return nil, err
	}
	c := paho.NewClient(paho.ClientConfig{
		Router: paho.NewSingleHandlerRouter(nil),
		Conn:   conn,
	})

	cp := &paho.Connect{
		KeepAlive:    30,
		CleanStart:   true,
		ClientID:     rpcCfg.ClientID,
		Username:     rpcCfg.Username,
		Password:     []byte(rpcCfg.Password),
		UsernameFlag: true,
		PasswordFlag: true,
	}
	ca, err := c.Connect(context.Background(), cp)
	if err != nil {
		log.Error().Msgf("connect mqtt server err: %+v", err)
		return nil, err
	}
	if ca.ReasonCode != 0 {
		err = fmt.Errorf("connect mqtt server failed: %d", ca.ReasonCode)
		log.Err(err)
		return nil, err
	}
	c.Router = paho.NewSingleHandlerRouter(rpcCfg.RecvHandler)
	rpcTopic := FormatRpcTopic(rpcCfg.ClientID)
	_, err = c.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			rpcTopic: {QoS: 0},
		},
	})
	if err != nil {
		log.Error().Msgf("subscribe topic %s err: %+v", rpcTopic, err)
		return nil, err
	}
	rpcHandler, err := rpc.NewHandler(c)
	if err != nil {
		log.Error().Msgf("build rpc request handler err: %+v", err)
		return nil, err
	}
	log.Info().Msg("NewMqttRpcHandler end")
	return rpcHandler, nil
}
