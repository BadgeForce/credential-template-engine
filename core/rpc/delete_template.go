package rpc

import (
	"fmt"
	"github.com/BadgeForce/credential-template-engine/core/verifier"

	"github.com/BadgeForce/credential-template-engine/core/state"
	template_pb "github.com/BadgeForce/credential-template-engine/core/template_pb"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
)

// DeleteTemplatesHandler RPC handler for deleteing templates
type DeleteTemplatesHandler struct {
	method string
}

// Handle ...
func (handler *DeleteTemplatesHandler) Handle(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcData interface{}) error {
	delete := rpcData.(*template_pb.Delete)
	issuerPub := request.GetHeader().GetSignerPublicKey()
	if addresses, valid := verifier.HasValidOwnership(issuerPub, delete.GetAddresses()...); !valid {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("unable to delete credential template(s) invalid ownership (%s)", addresses)}
	}

	return state.NewTemplateState(context).Delete(issuerPub, delete.GetAddresses()...)
}

// Method ...
func (handler *DeleteTemplatesHandler) Method() string {
	return handler.method
}

// DeleteHandle ...
var DeleteHandle = &DeleteTemplatesHandler{template_pb.Method_DELETE.String()}
