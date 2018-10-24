package rpc

import (
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/BadgeForce/credential-template-engine/core/state"
	"github.com/BadgeForce/credential-template-engine/core/proto"
	"encoding/json"
	"fmt"
)

var createHandle = func(request *processor_pb2.TpProcessRequest, context *processor.Context, rpcReq *credential_template_engine_pb.RPCRequest) error {
	// validate ownership
	// validate template data
	var templateData state.CredentialTemplate
	err := json.Unmarshal([]byte(rpcReq.Params), &templateData)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal template data from rpc request (%s)", err)}
	}

	template := state.NewCredentialTemplate(templateData.Name, templateData.Owner, templateData.Version, templateData.Data)
	return state.NewState(context).SaveTemplate(template)
}

var CREATE = &MethodHandler{createHandle, credential_template_engine_pb.Method_CREATE.String()}
