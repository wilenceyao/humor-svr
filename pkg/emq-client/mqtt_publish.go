package emq_client

func (c *EmqClient) MqttPublish(req *MqttPublishRequest) (*MqttPublishResponse, error) {
	res := &MqttPublishResponse{}
	err := c.post("/api/v4/mqtt/publish", req, res)
	return res, err
}
