package mqttrpc

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/paho"
	"github.com/eclipse/paho.golang/paho/extensions/rpc"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-api/proto/common"
	"github.com/wilenceyao/humor-api/proto/mqtt"
	"google.golang.org/protobuf/proto"
	"net"
	"reflect"
	"time"
)

var HUMOR_API_RPC_TOPIC string = "$share/humor-proto/rpc"
var TRACEID_KEY = "traceID"

var DefaultMqttRpcHandler *rpc.Handler
var DefaultMqttClient *paho.Client
var MqttRpcFuncs map[mqtt.Action]*MqttRpcFun

type MqttRpcConfig struct {
	ClientID string
	IP       string
	Port     uint
	Username string
	Password string
	RpcTopic string
}

type MqttRpcFun struct {
	ReqType reflect.Type
	ResType reflect.Type
	Fun     func(ctx context.Context, req interface{}, res interface{})
}

func FormatAgentRpcTopic(deviceID string) string {
	return fmt.Sprintf("rpc/%s", deviceID)
}

func CallMqttRpc(topic, traceID string, action mqtt.Action, serviceReq proto.Message,
	serviceRes proto.Message) error {
	log.Info().Msgf("CallMqttRpc request: %s %s %v %+v", traceID, topic, action, serviceReq)
	payloadBtArr, _ := proto.Marshal(serviceReq)
	mqttReq := &mqtt.RequestMessage{
		TraceID: traceID,
		Action:  action,
		Payload: payloadBtArr,
	}
	mqttReqBtArr, _ := proto.Marshal(mqttReq)
	rpcRes, err := DefaultMqttRpcHandler.Request(&paho.Publish{
		Topic:   topic,
		Payload: mqttReqBtArr,
	})
	if err != nil {
		return err
	}
	mqttRes := &mqtt.ReplyMessage{}
	_ = proto.Unmarshal(rpcRes.Payload, mqttRes)
	err = proto.Unmarshal(mqttRes.Payload, serviceRes)
	log.Info().Msgf("CallMqttRpc response: %s %s %v %+v", traceID, topic, action, serviceRes)
	return err
}

func RegisterMqttRpcFun(action mqtt.Action, req interface{}, res interface{},
	fun func(ctx context.Context, req interface{}, res interface{})) {
	log.Info().Msgf("RegisterMqttRpcFun, action: %v, fun: %+v", action, fun)
	f := &MqttRpcFun{
		ReqType: reflect.TypeOf(req),
		ResType: reflect.TypeOf(res),
		Fun:     fun,
	}
	MqttRpcFuncs[action] = f
}

func dispatch(reqMsg *mqtt.RequestMessage) *mqtt.ReplyMessage {
	resMsg := &mqtt.ReplyMessage{
		TraceID: reqMsg.TraceID,
		Action:  reqMsg.Action,
	}
	handler, ok := MqttRpcFuncs[reqMsg.Action]
	if !ok {
		resMsg.Code = common.ErrorCode_UNSUPPORTED_OPERATION
		resMsg.Msg = "func not found"
		return resMsg
	}
	ctx := context.WithValue(context.Background(), TRACEID_KEY, reqMsg.TraceID)
	req := reflect.New(handler.ReqType).Interface()
	res := reflect.New(handler.ResType).Interface()
	_ = proto.Unmarshal(reqMsg.Payload, req.(proto.Message))
	handler.Fun(ctx, req, res)
	resBtArr, _ := proto.Marshal(res.(proto.Message))
	resMsg.Payload = resBtArr
	return resMsg
}

func InitMqttClient(rpcCfg *MqttRpcConfig) error {
	log.Info().Msg("InitMqttClient start")
	server := fmt.Sprintf("%s:%d", rpcCfg.IP, rpcCfg.Port)
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Error().Msgf("connect to mqtt server %s err: %+v", server, err)
		return err
	}
	DefaultMqttClient = paho.NewClient(paho.ClientConfig{
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
	ca, err := DefaultMqttClient.Connect(context.Background(), cp)
	if err != nil {
		log.Error().Msgf("connect mqtt server err: %+v", err)
		return err
	}
	if ca.ReasonCode != 0 {
		err = fmt.Errorf("connect mqtt server failed: %d", ca.ReasonCode)
		log.Err(err)
		return err
	}
	DefaultMqttClient.Router = paho.NewSingleHandlerRouter(RpcMessageHandler)
	_, err = DefaultMqttClient.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			rpcCfg.RpcTopic: {QoS: 0},
		},
	})
	if err != nil {
		log.Error().Msgf("subscribe topic %s err: %+v", rpcCfg.RpcTopic, err)
		return err
	}
	DefaultMqttRpcHandler, err = rpc.NewHandler(DefaultMqttClient)
	if err != nil {
		log.Error().Msgf("build rpc request handler err: %+v", err)
		return err
	}
	MqttRpcFuncs = make(map[mqtt.Action]*MqttRpcFun)
	log.Info().Msg("InitMqttClient end")
	return nil
}

func RpcMessageHandler(m *paho.Publish) {
	if m.Properties != nil && m.Properties.CorrelationData != nil && m.Properties.ResponseTopic != "" {
		startT := time.Now()
		req, err := decode(m.Payload)
		if err != nil {
			log.Error().Msgf("decode msg err: %+v", err)
		}
		res := dispatch(req)
		endT := time.Now()
		log.Info().Msgf("[MQTT RPC] action: %s, elapsed: %v", req.Action.String(), endT.Sub(startT))
		btArr, err := encode(res)
		if err != nil {
			log.Error().Msgf("encode msg err: %+v", err)
		}
		_, err = DefaultMqttClient.Publish(context.Background(), &paho.Publish{
			Properties: &paho.PublishProperties{
				CorrelationData: m.Properties.CorrelationData,
			},
			Topic:   m.Properties.ResponseTopic,
			Payload: btArr,
		})
		if err != nil {
			log.Error().Msgf("publish msg err: %+v", err)
		}
	}
}
