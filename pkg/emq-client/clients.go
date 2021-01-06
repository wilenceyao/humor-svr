package emq_client

// GetClients 返回集群下所有信息
func (c *EmqClient) GetClients() (*GetClientsResponse, error) {
	res := &GetClientsResponse{}
	err := c.get("/proto/v4/clients", res)
	return res, err
}
