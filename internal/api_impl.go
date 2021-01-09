package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	agentapi "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humor-api/common"
	"github.com/wilenceyao/humor-api/svr/rest"
	"github.com/wilenceyao/humor-svr/config"
	"github.com/wilenceyao/humors"
	"net/http"
)

type ApiImpl struct {
	Adaptor *humors.HumorAdaptor
}

func (a *ApiImpl) Weather(c *gin.Context) {
	var req rest.WeatherRequest
	res := &rest.WeatherResponse{
		Response: &common.BaseResponse{
		},
	}
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		res.Response.Code = common.ErrorCode_INVALID_PARAMETERS
		c.JSON(http.StatusOK, res)
		return
	}
	agentReq := &agentapi.WeatherRequest{
		Request: req.Request,
	}
	agentRes := &agentapi.WeatherResponse{
		Response: &common.BaseResponse{},
	}
	err = a.Adaptor.Call(config.Config.AgentClientID, int32(agentapi.Action_WEATHER), agentReq, agentRes)
	if err != nil {
		log.Error().Msgf("call agent err: %v", err)
		res.Response.Code = common.ErrorCode_INTERNAL_ERROR
	} else {
		res.Response = agentRes.Response
	}
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
	agentReq := &agentapi.TtsRequest{
		Request: req.Request,
		Text:    req.Text,
	}
	agentRes := &agentapi.TtsResponse{
		Response: &common.BaseResponse{},
	}
	err = a.Adaptor.Call(config.Config.AgentClientID, int32(agentapi.Action_TTS), agentReq, agentRes)
	if err != nil {
		log.Error().Msgf("call agent err: %v", err)
		res.Response.Code = common.ErrorCode_INTERNAL_ERROR
	} else {
		res.Response = agentRes.Response
	}
	c.JSON(http.StatusOK, res)
}
