package rpc

import (
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/golang/protobuf/proto"
	"github.com/BadgeForce/credential-template-engine/core/proto"
	"github.com/BadgeForce/badgeforce-chain-node/core/common"
)

var delegateCB = func(request *processor_pb2.TpProcessRequest) (string, interface{}, error) {
	var rpcRequest credential_template_engine_pb.RPCRequest
	err := proto.Unmarshal(request.GetPayload(), &rpcRequest)
	if err != nil {
		return "", nil, &processor.InvalidTransactionError{Msg: "malformed payload data"}
	}

	return rpcRequest.GetMethod().String(), &rpcRequest, nil
}

func NewClient() *common.RPCClient {
	handlers := []common.MethodHandler{CreateHandle, UpdateHandle, DeleteHandle}
	return common.NewRPCClient(handlers, delegateCB)
}
