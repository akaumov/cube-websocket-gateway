package js

type MessageType int

const (
	TEXT   MessageType = 0
	BINARY MessageType = 1
)

type MessageParams struct {
	InputTime int64       `json:"inputTime"`
	UserId    *string     `json:"userId"`
	DeviceId  *string     `json:"deviceId"`
	Type      MessageType `json:"type"`
	Body      []byte      `json:"body"`
}
