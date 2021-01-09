package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-api/api/common"
	"github.com/wilenceyao/humor-api/api/rest"
	emq_client "github.com/wilenceyao/humor-api/pkg/emq-client"
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
	log.Info().Msgf("before getde")
	deviceId := req.DeviceID
	if deviceId == "" {
		deviceId, err = a.getDefaultDevice()
		if err != nil {
			res.Response.Code = common.ErrorCode_UNSUPPORTED_OPERATION
			res.Response.Msg = "no device available"
			c.JSON(http.StatusOK, res)
			return
		}
	}
	// TODO

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
