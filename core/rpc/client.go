package rpc

import (
	"github.com/BadgeForce/sawtooth-utils"
	"github.com/BadgeForce/sawtooth-utils/protos/templates_pb"
	"github.com/golang/protobuf/proto"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
)

// Client ...
var Client *utils.RPCClient

var delegateCB = func(request *processor_pb2.TpProcessRequest) (string, interface{}, error) {
	var rpcRequest templates_pb.RPCRequest
	err := proto.Unmarshal(request.GetPayload(), &rpcRequest)
	if err != nil {
		return "", nil, &processor.InvalidTransactionError{Msg: "unable to unmarshal RPC request"}
	}

	switch method := rpcRequest.Method.(type) {
	case *templates_pb.RPCRequest_Create:
		return templates_pb.Method_CREATE.String(), method.Create, nil
	case *templates_pb.RPCRequest_Delete:
		return templates_pb.Method_DELETE.String(), method.Delete, nil
	default:
		return "", nil, &processor.InvalidTransactionError{Msg: "invalid RPC method"}
	}
}

func init() {
	handlers := []utils.MethodHandler{CreateHandle, DeleteHandle}
	Client = utils.NewRPCClient(handlers, delegateCB)
}
