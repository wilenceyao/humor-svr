package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/wilenceyao/humor-api/api/mqtt"
	"github.com/wilenceyao/humor-api/api/rest"
	emq_client "github.com/wilenceyao/humor-api/pkg/emq-client"
	"net/http"
)

type ApiImpl struct {
}

func (a *ApiImpl) GetDevices(c *gin.Context) {
	res, _ := emq_client.DefaultEmqClient.GetClients()
	c.JSON(http.StatusOK, res)
}

func (a *ApiImpl) SendTts(c *gin.Context) {
	var req rest.SendTtsRequest
	res := &rest.SendTtsResponse{}
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Code = rest.INVALID_PARAMETERS
		c.JSON(http.StatusOK, res)
		return
	}
	if req.TraceID == "" || req.DeviceID == "" || req.Text == "" {
		res.Code = rest.INVALID_PARAMETERS
		c.JSON(http.StatusOK, res)
		return
	}
	payload := &mqtt.TtsPayload{
		Text: req.Text,
	}
	payloadBtArr, _ := proto.Marshal(payload)
	mqttMsg := &mqtt.Message{
		TraceID: req.TraceID,
		Action:  mqtt.Action_TTS,
		Payload: payloadBtArr,
	}
	msgBtArr, _ := proto.Marshal(mqttMsg)
	mqttReq := &emq_client.MqttPublishRequest{
		Clientid: "humor-api",
		Payload:  string(msgBtArr),
		Topic:    fmt.Sprintf("device/%s", req.DeviceID),
	}
	mqttRes, err := emq_client.DefaultEmqClient.MqttPublish(mqttReq)
	if err != nil {
		res.Code = rest.INTERNAL_ERROR
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}
	if mqttRes.Code != 0 {
		res.Code = rest.EXTERNAL_ERROR
	}
	c.JSON(http.StatusOK, res)
}
