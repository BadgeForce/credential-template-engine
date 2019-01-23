package rpc

import (
	"fmt"

	"github.com/BadgeForce/credential-template-engine/core/template_pb"

	"github.com/BadgeForce/credential-template-engine/core/state"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
)

type CreateTemplateHandler struct {
	method string
}

func (handler *CreateTemplateHandler) Handle(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq interface{}) error {
	return handler.createTemplate(request, context, rpcReq.(*template_pb.RPCRequest))
}

func (handler *CreateTemplateHandler) Method() string {
	return handler.method
}

func (handler *CreateTemplateHandler) createTemplate(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *template_pb.RPCRequest) error {
	data, err := NewPayloadDecoder(rpcReq).UnmarshalCreate()
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal template data from rpc request (%s)", err)}
	}

	return state.NewTemplateState(context).Save(state.NewCredentialTemplate(data.Name, data.Owner, data.Version, data.Data))
}

var CreateHandle = &CreateTemplateHandler{template_pb.Method_CREATE.String()}
