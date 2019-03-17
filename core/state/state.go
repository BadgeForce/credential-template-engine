package state

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/BadgeForce/credential-template-engine/core/verifier"
	"github.com/BadgeForce/sawtooth-utils"
	"github.com/BadgeForce/sawtooth-utils/protos/templates_pb"
	"github.com/rberg2/sawtooth-go-sdk/logging"
	"github.com/rberg2/sawtooth-go-sdk/processor"
)

//TODO: move CredentialTemplatePrefix to configuration

// CredentialTemplatePrefix ...
const CredentialTemplatePrefix = "credential:templates"

var (
	logger = logging.Get()

	// NameSpaceMngr ...
	NameSpaceMngr *utils.NamespaceMngr
)

// State ...
type State struct {
	instance *utils.State
}

// Context ...
func (s *State) Context() *processor.Context {
	return s.instance.Context
}

// GetTxtRecpt returns a transaction receipt with correct data
func (s *State) GetTxtRecpt(rpcMethod templates_pb.Method, stateAddress string, template *templates_pb.Template) (*templates_pb.Receipt, []byte, error) {
	var recpt templates_pb.Receipt
	recpt.Date = time.Now().Unix()
	recpt.StateAddress = stateAddress
	recpt.RpcMethod = rpcMethod
	recpt.Template = template

	b, err := proto.Marshal(&recpt)
	return &recpt, b, err
}

// Save saves a template template to state
func (s *State) Save(template *templates_pb.Template) error {
	address := TemplateStateAddress(
		template.GetData().GetIssuerPub(),
		template.GetData().GetName(),
		template.GetData().GetVersion(),
	)

	_, receiptBytes, err := s.GetTxtRecpt(templates_pb.Method_CREATE, address, template)
	if err != nil {
		logger.Warnf("unable to generate transaction receipt for template saved (%s)", err)
	}
	b, err := proto.Marshal(template)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("unable to marshal template proto (%s)", err)}
	}

	_, err = s.Context().SetState(map[string][]byte{address: b})
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("unable to save credential temlate (%s)", err)}
	}

	err = s.Context().AddReceiptData(receiptBytes)
	if err != nil {
		logger.Warnf("unable to add transaction receipt for template saved (%s)", err)
	}

	return nil
}

// Delete delete stored at each specified address in state
func (s *State) Delete(issuerPub string, addresses ...string) error {
	_, err := s.Context().DeleteState(addresses)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("unable to delete credential temlate (%s)", err)}
	}

	for _, address := range addresses {
		_, receiptBytes, err := s.GetTxtRecpt(templates_pb.Method_DELETE, address, nil)
		if err != nil {
			logger.Warnf("unable to generate transaction receipt for template saved (%s)", err)
		}

		err = s.Context().AddReceiptData(receiptBytes)
		if err != nil {
			logger.Warnf("unable to add transaction receipt for template saved (%s)", err)
		}

	}

	return nil
}

// GetTemplates get some templates stored at each specified address from state
func (s *State) GetTemplates(issuerPub string, address ...string) ([]*templates_pb.Template, error) {
	if addrs, ok := verifier.HasValidOwnership(issuerPub, address...); !ok {
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state invalid ownership of templates (%s)", addrs)}
	}

	state, err := s.Context().GetState(address)
	if err != nil {
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
	}

	templates := make([]*templates_pb.Template, 0)
	for _, value := range state {
		var template templates_pb.Template
		err := proto.Unmarshal(value, &template)
		if err != nil {
			return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// TemplateStateAddress ...
func TemplateStateAddress(issuerPub, name string, version *templates_pb.Version) string {
	vrsn := fmt.Sprintf("%x.%x.%x", version.GetMajor(), version.GetMinor(), version.GetPatch())
	o := utils.NewPart(issuerPub, 0, 30)
	n := utils.NewPart(name, 0, 30)
	v := utils.NewPart(vrsn, 0, 4)

	addressParts := []*utils.AddressPart{o, n, v}
	address, _ := utils.NewAddress(NameSpaceMngr.NameSpaces[0]).AddParts(addressParts...).Build()
	return address
}

// NewTemplateState ...
func NewTemplateState(context *processor.Context) *State {
	return &State{utils.NewStateInstance(context)}
}

func init() {
	NameSpaceMngr = utils.NewNamespaceMngr().RegisterNamespaces(CredentialTemplatePrefix)
}
