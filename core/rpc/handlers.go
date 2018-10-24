package rpc

import (
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/BadgeForce/credential-template-engine/core/proto"
)

type MethodHandler struct {
	Handle func(*processor_pb2.TpProcessRequest, *processor.Context, *credential_template_engine_pb.RPCRequest) error
	Method string
}

type RPCClient struct {
	MethodHandlers map[string]*MethodHandler
}

func (r *RPCClient) registerMethod(handler *MethodHandler) *RPCClient {
	r.MethodHandlers[handler.Method] = handler
	return r
}

func (r *RPCClient) DelegateMethod(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	var rpcRequest credential_template_engine_pb.RPCRequest
	err := proto.Unmarshal(request.GetPayload(), &rpcRequest)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: "malformed payload data"}
	}

	return r.delegate(request, context, rpcRequest)
}

func (r *RPCClient) delegate(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcRequest credential_template_engine_pb.RPCRequest) error {
	method := rpcRequest.GetMethod().String()

	if methodHandler, exists := r.MethodHandlers[method]; exists {
		return methodHandler.Handle(request, context, &rpcRequest)
	}

	return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not determine RPC method: %v", method)}
}

func NewClient() *RPCClient {
	client := &RPCClient{make(map[string]*MethodHandler)}
	return client.registerMethod(CREATE)
}
