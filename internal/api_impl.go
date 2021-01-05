package internal

import (
	"fmt"
	"github.com/eclipse/paho.golang/paho"
	"github.com/eclipse/paho.golang/paho/extensions/rpc"
	"github.com/gin-gonic/gin"
	"github.com/wilenceyao/humor-api/api/mqtt"
	"github.com/wilenceyao/humor-api/api/rest"
	emq_client "github.com/wilenceyao/humor-api/pkg/emq-client"
	"github.com/wilenceyao/humor-api/pkg/util"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type ApiImpl struct {
	MqttRpcReqHandler *rpc.Handler
}

func (a *ApiImpl) GetDevices(c *gin.Context) {
	res, _ := emq_client.DefaultEmqClient.GetClients()
	c.JSON(http.StatusOK, res)
}

func (a *ApiImpl) SendTts(c *gin.Context) {
	var req rest.SendTtsRequest
	res := &rest.SendTtsResponse{}
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		res.Code = rest.INVALID_PARAMETERS
		c.JSON(http.StatusOK, res)
		return
	}
	if req.TraceID == "" || req.Text == "" {
		res.Code = rest.INVALID_PARAMETERS
		c.JSON(http.StatusOK, res)
		return
	}
	deviceId := req.DeviceID
	if deviceId == "" {
		deviceId, err = a.getDefaultDevice()
		if err != nil {
			res.Code = rest.UNSUPPORTED_OPERATION
			res.Msg = "no device"
			c.JSON(http.StatusOK, res)
			return
		}
	}
	payload := &mqtt.TtsRequest{
		Text: req.Text,
	}
	payloadBtArr, _ := proto.Marshal(payload)
	mqttReqMsg := &mqtt.Message{
		TraceID: req.TraceID,
		Action:  mqtt.Action_TTS,
		Payload: payloadBtArr,
	}
	reqMsgBtArr, _ := proto.Marshal(mqttReqMsg)
	rpcRes, err := a.MqttRpcReqHandler.Request(&paho.Publish{
		Topic:   util.FormatRpcTopic(deviceId),
		Payload: reqMsgBtArr,
	})
	if err != nil {
		res.Code = rest.INTERNAL_ERROR
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	mqttResMsg := &mqtt.Message{}
	_ = proto.Unmarshal(rpcRes.Payload, mqttResMsg)
	resPayload := &mqtt.TtsReply{}
	_ = proto.Unmarshal(mqttResMsg.Payload, resPayload)
	if resPayload.Reply.Code != mqtt.ErrorCode_SUCCESS {
		res.Code = rest.EXTERNAL_ERROR
	}
	c.JSON(http.StatusOK, res)
}

func (a *ApiImpl) getDefaultDevice() (string, error) {
	getClientsRes, err := emq_client.DefaultEmqClient.GetClients()
	if err != nil {
		return "", err
	}
	if len(getClientsRes.Data) == 0 {
		return "", fmt.Errorf("no devices")
	}
	return getClientsRes.Data[0].Clientid, nil
}
