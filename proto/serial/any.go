package serial

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func MarshalAny(message proto.Message) (*any.Any, error) {
	return ptypes.MarshalAny(message)
}

func UnmarshalAny(any *any.Any) (proto.Message, error) {
	protoMessage, err := ptypes.Empty(any)
	if err != nil {
		return nil, err
	}
	if err := ptypes.UnmarshalAny(any, protoMessage); err != nil {
		return nil, err
	}
	return protoMessage, nil
}
