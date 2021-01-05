package humor_util

import "fmt"

func FormatRpcReqTopic(deviceID string) string {
	return fmt.Sprintf("rpc/req/%s", deviceID)
}

func FormatRpcResTopic(deviceID string) string {
	return fmt.Sprintf("rpc/res/%s", deviceID)
}
