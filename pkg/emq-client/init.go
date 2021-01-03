package emq_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

var (
	AppID     string
	AppSecret string
	Ip        string
	Port      int
)
var DefaultEmqClient *EmqClient

type EmqClient struct {
	httpClient *http.Client
	appId      string
	appSecret  string
	ip         string
	port       uint
}

func InitEmqClient(id, secret, serverIp string, serverPort uint) error {
	DefaultEmqClient = &EmqClient{
		httpClient: &http.Client{},
		appId:      id,
		appSecret:  secret,
		ip:         serverIp,
		port:       serverPort,
	}
	return nil
}

func (c *EmqClient) get(uri string, resObj interface{}) error {
	url := fmt.Sprintf("http://%s:%d%s", c.ip, c.port, uri)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Msgf("uri: %s, build request err: %+v", uri, err)
		return err
	}
	req.SetBasicAuth(c.appId, c.appSecret)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Msgf("uri: %s, request to broker err: %+v", uri, err)
		return err
	}
	btArr, err := ioutil.ReadAll(resp.Body)
	log.Debug().Msgf("[EMQ API GET] url: %s, res: %s", url, string(btArr))
	return json.Unmarshal(btArr, resObj)
}

func (c *EmqClient) post(uri string, reqObj, resObj interface{}) error {
	url := fmt.Sprintf("http://%s:%d%s", c.ip, c.port, uri)
	reqBtArr, _ := json.Marshal(reqObj)
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBtArr))
	if err != nil {
		log.Error().Msgf("uri: %s, build request err: %+v", uri, err)
		return err
	}
	req.SetBasicAuth(c.appId, c.appSecret)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().Msgf("uri: %s, request to broker err: %+v", uri, err)
		return err
	}
	btArr, err := ioutil.ReadAll(resp.Body)
	log.Debug().Msgf("[EMQ API POST] uri: %s, req: %s, res: %s",
		url, string(reqBtArr), string(btArr))
	return json.Unmarshal(btArr, resObj)
}
