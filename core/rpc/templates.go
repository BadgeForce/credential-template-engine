package rpc

import (
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/propsproject/pending-props/core/proto/pending_props_pb"
	"github.com/BadgeForce/credential-template-engine/core/state"
)

var createHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *pending_props_pb.RPCRequest) error {
	// validate ownership
	// validate template data
	return state.NewState(context).SaveTemplate(nil)
}

var CREATE = &MethodHandler{createHandle, pending_props_pb.Method_ISSUE.String()}
