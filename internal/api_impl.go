package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	emq_client "github.com/wilenceyao/humor-api/pkg/emq-client"
	"github.com/wilenceyao/humor-api/pkg/mqttrpc"
	"github.com/wilenceyao/humor-api/proto/common"
	"github.com/wilenceyao/humor-api/proto/mqtt"
	"github.com/wilenceyao/humor-api/proto/rest"
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
	log.Info().Msgf("before getde")
	deviceId := req.DeviceID
	if deviceId == "" {
		deviceId, err = a.getDefaultDevice()
		if err != nil {
			res.Code = rest.UNSUPPORTED_OPERATION
			res.Msg = "no device available"
			c.JSON(http.StatusOK, res)
			return
		}
	}
	mqttReq := &mqtt.TtsRequest{
		Text: req.Text,
	}
	mqttRes := &mqtt.TtsReply{}
	err = mqttrpc.CallMqttRpc(mqttrpc.FormatAgentRpcTopic(deviceId), req.TraceID, mqtt.Action_TTS, mqttReq, mqttRes)
	if mqttRes.Reply.Code != common.ErrorCode_SUCCESS {
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
