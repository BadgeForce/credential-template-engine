package core

import (
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"github.com/rberg2/sawtooth-go-sdk/protobuf/processor_pb2"
	"github.com/BadgeForce/credential-template-engine/core/rpc"
	"github.com/BadgeForce/credential-template-engine/core/state"
)

// TransactionHandler ...
type TransactionHandler struct {
	FName     string   `json:"familyName"`
	FVersions []string `json:"familyVersions"`
	NSpace    []string `json:"nameSpace"`
	RPCClient *rpc.RPCClient
}

// FamilyName ...
const FamilyName = "bf-credential-templates" // move to configuration

// FamilyVersions ...
var FamilyVersions = []string{"1.0"} // move to configuration

// FamilyName ...
func (t *TransactionHandler) FamilyName() string {
	return t.FName
}

// FamilyVersions ...
func (t *TransactionHandler) FamilyVersions() []string {
	return t.FVersions
}

// Namespaces ...
func (t *TransactionHandler) Namespaces() []string {
	return t.NSpace
}

// Apply ...
func (t *TransactionHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	return t.RPCClient.DelegateMethod(request, context)
}

// NewTransactionHandler returns a new transaction handler
func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{
		FName:     FamilyName,
		FVersions: FamilyVersions,
		NSpace:    []string{state.TemplatePrefix},
		RPCClient: rpc.NewClient(),
	}
}
