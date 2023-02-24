package network

import (
	"github.com/AlexEkdahl/snakes/network/protobuf"
	"github.com/golang/protobuf/proto"
)

type Messenger interface {
	EncodeMessage() ([]byte, error)
	DecodeMessage([]byte) (*protobuf.Message, error)
}

type protobufHandler struct {
	Message interface{}
}

func (mh *protobufHandler) EncodeMessage() ([]byte, error) {
	var pbMessage proto.Message

	switch msg := mh.Message.(type) {
	case *protobuf.MoveMessage:
		pbMessage = &protobuf.Message{
			Type: &protobuf.Message_Move{Move: msg},
		}
	case *protobuf.DisconnectMessage:
		pbMessage = &protobuf.Message{
			Type: &protobuf.Message_Disconnect{Disconnect: msg},
		}
	default:
		joinMsg := &protobuf.Message{
			Type: &protobuf.Message_Join{
				Join: &protobuf.JoinMessage{},
			},
		}
		return proto.Marshal(joinMsg)
	}

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
