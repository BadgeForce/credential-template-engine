package rpc

import (
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/BadgeForce/credential-template-engine/core/proto"
	"fmt"
	"github.com/BadgeForce/credential-template-engine/core/state"
)

type DeleteTemplatesHandler struct {
	method string
}

func (handler *DeleteTemplatesHandler) Handle(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq interface{}) error {
	return handler.createTemplate(request, context, rpcReq.(*credential_template_engine_pb.RPCRequest))
}

func (handler *DeleteTemplatesHandler) Method() string {
	return handler.method
}

func (handler *DeleteTemplatesHandler) createTemplate(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *credential_template_engine_pb.RPCRequest) error {
	payload, err := NewPayloadDecoder(rpcReq).UnmarshalDelete()
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal template data from rpc request (%s)", err)}
	}
	return state.NewTemplateState(context).Delete(payload.Addresses...)
}

var DeleteHandle = &DeleteTemplatesHandler{credential_template_engine_pb.Method_DELETE.String()}
