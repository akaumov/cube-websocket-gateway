package js

import (
	"encoding/json"
)

type MessageType int

const (
	TEXT   MessageType = 0
	BINARY MessageType = 1
)

type OnReceiveMessageParams struct {
	InputTime int64       `json:"inputTime"`
	UserId    *string     `json:"userId"`
	DeviceId  *string     `json:"deviceId"`
	Type      MessageType `json:"type"`
	Body      []byte      `json:"body"`
}

type CloseDeviceConnectionsParams struct {
	UserId   string `json:"userId"`
	DeviceId string `json:"deviceId"`
	Reason   string `json:"reason"`
}

type CloseUserConnectionsParams struct {
	UserId string `json:"userId"`
	Reason string `json:"reason"`
}

type SendMessageParams struct {
	UserId   *string     `json:"userId"`
	DeviceId *string     `json:"deviceId"`
	Type     MessageType `json:"type"`
	Body     []byte      `json:"body"`
}

type RoutingPacket struct {
	Endpoint string          `json:"endpoint"`
	Payload  json.RawMessage `json:"payload"`
}
