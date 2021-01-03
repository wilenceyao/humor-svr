package emq_client

func (c *EmqClient) MqttPublish(req *MqttPublishRequest) (*MqttPublishResponse, error) {
	res := &MqttPublishResponse{}
	err := c.post("/api/v4/mqtt/publish", res, req)
	return res, err
}
