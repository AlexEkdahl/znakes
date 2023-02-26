package network

import (
	"github.com/AlexEkdahl/snakes/pkg/network/protobuf"
	"google.golang.org/protobuf/proto"
)

type Messenger interface {
	EncodeMessage(proto.Message) ([]byte, error)
	DecodeMessage([]byte) (*protobuf.Message, error)
}

type protobufHandler struct {
	Message interface{}
}

func (mh *protobufHandler) EncodeMessage(pbMessage proto.Message) ([]byte, error) {
	return proto.Marshal(pbMessage)
}

// Decode decodes a protobuf message into a message
func (mh *protobufHandler) DecodeMessage(data []byte) (*protobuf.Message, error) {
	var msg protobuf.Message
	if err := proto.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
