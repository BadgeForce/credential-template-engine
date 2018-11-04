package rpc

import (
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/BadgeForce/credential-template-engine/core/proto"
	"fmt"
	"github.com/BadgeForce/credential-template-engine/core/state"
)

type UpdateTemplateHandler struct {
	method string
}

func (handler *UpdateTemplateHandler) Handle(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq interface{}) error {
	return handler.updateTemplate(request, context, rpcReq.(*credential_template_engine_pb.RPCRequest))
}

func (handler *UpdateTemplateHandler) Method() string {
	return handler.method
}

func (handler *UpdateTemplateHandler) updateTemplate(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *credential_template_engine_pb.RPCRequest) error {
	template, err := NewPayloadDecoder(rpcReq).UnmarshalUpdate()
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal template data from rpc request (%s)", err)}
	}

	return state.NewTemplateState(context).Save(template)
}

var UpdateHandle = &UpdateTemplateHandler{ credential_template_engine_pb.Method_UPDATE.String()}