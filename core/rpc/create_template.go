package rpc

import (
	"github.com/BadgeForce/sawtooth-utils/protos/templates_pb"

	"github.com/BadgeForce/credential-template-engine/core/state"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
)

// CreateTemplateHandler RPC handler to handle creating templates
type CreateTemplateHandler struct {
	method string
}

// Handle ...
func (handler *CreateTemplateHandler) Handle(request *processor_pb2.TpProcessRequest, context *processor.Context, reqData interface{}) error {
	create := reqData.(*templates_pb.Create)

	return state.NewTemplateState(context).Save(create.GetParams())
}

// Method ...
func (handler *CreateTemplateHandler) Method() string {
	return handler.method
}

// CreateHandle ...
var CreateHandle = &CreateTemplateHandler{templates_pb.Method_CREATE.String()}
