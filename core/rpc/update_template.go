package rpc

import (
	"fmt"

	"github.com/BadgeForce/credential-template-engine/core/state"
	"github.com/BadgeForce/credential-template-engine/core/template_pb"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
)

type UpdateTemplateHandler struct {
	method string
}

func (handler *UpdateTemplateHandler) Handle(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq interface{}) error {
	return handler.updateTemplate(request, context, rpcReq.(*template_pb.RPCRequest))
}

func (handler *UpdateTemplateHandler) Method() string {
	return handler.method
}

func (handler *UpdateTemplateHandler) updateTemplate(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *template_pb.RPCRequest) error {
	template, err := NewPayloadDecoder(rpcReq).UnmarshalUpdate()
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal template data from rpc request (%s)", err)}
	}

	return state.NewTemplateState(context).Save(template)
}

var UpdateHandle = &UpdateTemplateHandler{template_pb.Method_UPDATE.String()}
