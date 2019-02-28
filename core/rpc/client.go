package rpc

import (
	template_pb "github.com/BadgeForce/credential-template-engine/core/template_pb"
	"github.com/BadgeForce/sawtooth-utils"
	"github.com/golang/protobuf/proto"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
)

// Client ...
var Client *utils.RPCClient

var delegateCB = func(request *processor_pb2.TpProcessRequest) (string, interface{}, error) {
	var rpcRequest template_pb.RPCRequest
	err := proto.Unmarshal(request.GetPayload(), &rpcRequest)
	if err != nil {
		return "", nil, &processor.InvalidTransactionError{Msg: "unable to unmarshal RPC request"}
	}

	switch method := rpcRequest.Method.(type) {
	case *template_pb.RPCRequest_Create:
		return template_pb.Method_CREATE.String(), method.Create, nil
	case *template_pb.RPCRequest_Delete:
		return template_pb.Method_DELETE.String(), method.Delete, nil
	default:
		return "", nil, &processor.InvalidTransactionError{Msg: "invalid RPC method"}
	}
}

func init() {
	handlers := []utils.MethodHandler{CreateHandle, DeleteHandle}
	Client = utils.NewRPCClient(handlers, delegateCB)
}
