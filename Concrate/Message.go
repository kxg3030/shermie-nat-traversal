package Concrate

import (
	"encoding/json"
)

const (
	ActionConnect = 0x1
	ActionData    = 0x2
)

// 消息体
type Message struct {
	Action    uint8  `json:"action,omitempty"`
	UniqueId  int32  `json:"unique_id,omitempty"`
	HeaderLen uint32 `json:"header_len,omitempty"`
	Data      []byte `json:"data,omitempty"`
	To        int32  `json:"to,omitempty"`
}

func NewMessage() *Message {
	return &Message{}
}

func (i *Message) Json() string {
	jsonBody, _ := json.Marshal(i)
	return string(jsonBody)
}

func (i *Message) DataString() string {
	dataByte := i.Data
	return string(dataByte)
}
