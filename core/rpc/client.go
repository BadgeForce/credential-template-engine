package rpc

import (
	"github.com/BadgeForce/credential-template-engine/core/template_pb"
	utils "github.com/BadgeForce/sawtooth-utils"
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
		return "", nil, &processor.InvalidTransactionError{Msg: "malformed payload data"}
	}

	return rpcRequest.GetMethod().String(), &rpcRequest, nil
}

func init() {
	handlers := []utils.MethodHandler{CreateHandle, UpdateHandle, DeleteHandle}
	Client = utils.NewRPCClient(handlers, delegateCB)
}
