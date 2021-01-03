package emq_client

type Meta struct {
	Count   int64 `json:"count"`
	Hasnext bool  `json:"hasnext"`
	Limit   int64 `json:"limit"`
	Page    int64 `json:"page"`
}

type Client struct {
	AwaitingRel      int64  `json:"awaiting_rel"`
	CleanStart       bool   `json:"clean_start"`
	Clientid         string `json:"clientid"`
	Connected        bool   `json:"connected"`
	ConnectedAt      string `json:"connected_at"`
	CreatedAt        string `json:"created_at"`
	ExpiryInterval   int64  `json:"expiry_interval"`
	HeapSize         int64  `json:"heap_size"`
	Inflight         int64  `json:"inflight"`
	IPAddress        string `json:"ip_address"`
	IsBridge         bool   `json:"is_bridge"`
	Keepalive        int64  `json:"keepalive"`
	MailboxLen       int64  `json:"mailbox_len"`
	MaxAwaitingRel   int64  `json:"max_awaiting_rel"`
	MaxInflight      int64  `json:"max_inflight"`
	MaxMqueue        int64  `json:"max_mqueue"`
	MaxSubscriptions int64  `json:"max_subscriptions"`
	Mountpoint       string `json:"mountpoint"`
	MqueueDropped    int64  `json:"mqueue_dropped"`
	MqueueLen        int64  `json:"mqueue_len"`
	Node             string `json:"node"`
	Port             int64  `json:"port"`
	ProtoName        string `json:"proto_name"`
	ProtoVer         int64  `json:"proto_ver"`
	RecvCnt          int64  `json:"recv_cnt"`
	RecvMsg          int64  `json:"recv_msg"`
	RecvOct          int64  `json:"recv_oct"`
	RecvPkt          int64  `json:"recv_pkt"`
	Reductions       int64  `json:"reductions"`
	SendCnt          int64  `json:"send_cnt"`
	SendMsg          int64  `json:"send_msg"`
	SendOct          int64  `json:"send_oct"`
	SendPkt          int64  `json:"send_pkt"`
	SubscriptionsCnt int64  `json:"subscriptions_cnt"`
	Username         string `json:"username"`
	Zone             string `json:"zone"`
}

type BaseResponse struct {
	Code int64 `json:"code"`
}
type GetClientsResponse struct {
	BaseResponse
	Data []Client `json:"data"`
	Meta Meta     `json:"meta"`
}

type MqttPublishRequest struct {
	Clientid string `json:"clientid"`
	Payload  string `json:"payload"`
	Qos      int64  `json:"qos"`
	Retain   bool   `json:"retain"`
	Topic    string `json:"topic"`
}

type MqttPublishResponse struct {
	BaseResponse
}
