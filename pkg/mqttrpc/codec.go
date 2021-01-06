package mqttrpc

import (
	"github.com/wilenceyao/humor-api/proto/mqtt"
	"google.golang.org/protobuf/proto"
)

func decode(payload []byte) (*mqtt.RequestMessage, error) {
	req := &mqtt.RequestMessage{}
	err := proto.Unmarshal(payload, req)
	return req, err
}

func encode(message *mqtt.ReplyMessage) ([]byte, error) {
	return proto.Marshal(message)
}
