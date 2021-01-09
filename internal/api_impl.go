package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	humoragent "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humor-api/common"
	"github.com/wilenceyao/humor-api/svr/rest"
	emq_client "github.com/wilenceyao/humor-svr/pkg/emq-client"
	"github.com/wilenceyao/humors"
	"net/http"
)

type ApiImpl struct {
	Adaptor *humors.HumorAdaptor
}

func (a *ApiImpl) GetDevices(c *gin.Context) {
	res, _ := emq_client.DefaultEmqClient.GetClients()
	c.JSON(http.StatusOK, res)
}

func (a *ApiImpl) SendTts(c *gin.Context) {
	var req rest.TtsRequest
	res := &rest.TtsResponse{
		Response: &common.BaseResponse{
		},
	}
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		res.Response.Code = common.ErrorCode_INVALID_PARAMETERS
		c.JSON(http.StatusOK, res)
		return
	}
	if req.Request.RequestID == "" || req.Text == "" {
		res.Response.Code = common.ErrorCode_INVALID_PARAMETERS
		c.JSON(http.StatusOK, res)
		return
	}
	clientID := req.ClientID
	if clientID == "" {
		clientID, err = a.getDefaultClient()
		if err != nil {
			res.Response.Code = common.ErrorCode_UNSUPPORTED_OPERATION
			res.Response.Msg = "no device available"
			c.JSON(http.StatusOK, res)
			return
		}
	}
	agentReq := &humoragent.TtsRequest{
		Request: req.Request,
		Text:    req.Text,
	}
	agentRes := &humoragent.TtsResponse{
		Response: &common.BaseResponse{},
	}
	err = a.Adaptor.Call(clientID, int32(humoragent.Action_TTS), agentReq, agentRes)
	if err != nil {
		log.Error().Msgf("call agent err: %v", err)
		res.Response.Code = common.ErrorCode_INTERNAL_ERROR
	} else {
		res.Response.Code = agentRes.Response.Code
		res.Response.Msg = agentRes.Response.Msg
	}
	c.JSON(http.StatusOK, res)
}

func (a *ApiImpl) getDefaultClient() (string, error) {
	getClientsRes, err := emq_client.DefaultEmqClient.GetClients()
	if err != nil {
		return "", err
	}
	if len(getClientsRes.Data) == 0 {
		return "", fmt.Errorf("no devices")
	}
	return getClientsRes.Data[0].Clientid, nil
}
