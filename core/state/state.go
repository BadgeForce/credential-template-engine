package state

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/BadgeForce/credential-template-engine/core/template_pb"
	utils "github.com/BadgeForce/sawtooth-utils"
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
func (s *State) GetTxtRecpt(rpcMethod template_pb.Method, stateAddress string, template *template_pb.Template) (*template_pb.Receipt, []byte, error) {
	var recpt template_pb.Receipt
	recpt.Date = time.Now().Unix()
	recpt.StateAddress = stateAddress
	recpt.RpcMethod = rpcMethod
	recpt.Template = template

	b, err := proto.Marshal(&recpt)
	return &recpt, b, err
}

// VerifyTemplate ...
func VerifyTemplate(txtSignerPub string, template *template_pb.Template) (bool, error) {
	b, err := proto.Marshal(template.GetData())
	if err != nil {
		return false, fmt.Errorf("error: could not marshal proto (%s)", err)
	}

	expectedHash := template.GetVerification().GetProofOfIntegrityHash()
	if hash, ok := utils.VerifyPOIHash(b, expectedHash); !ok {
		return false, fmt.Errorf("error: proof of integrity hash invalid got (%s) want (%s)", hash, expectedHash)
	}

	issuerPub := template.GetData().GetIssuerPub()
	sig := template.GetVerification().GetSignature()

	if ok := utils.VerifySig(sig, []byte(txtSignerPub), b, false); !ok {
		return false, fmt.Errorf("error: transaction signer must also be owner who signs template (%s)", txtSignerPub)
	} else if txtSignerPub != issuerPub {
		return false, fmt.Errorf("error: transaction signer public key must match template issuer got (%s) want (%s)", issuerPub, txtSignerPub)
	}

	return true, nil
}

// Save saves a template template to state
func (s *State) Save(template *template_pb.Template) error {
	address := TemplateStateAddress(
		template.GetData().GetIssuerPub(),
		template.GetData().GetName(),
		template.GetData().GetVersion(),
	)

	_, receiptBytes, err := s.GetTxtRecpt(template_pb.Method_CREATE, address, template)
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
		_, receiptBytes, err := s.GetTxtRecpt(template_pb.Method_DELETE, address, nil)
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
func (s *State) GetTemplates(issuerPub string, address ...string) ([]*template_pb.Template, error) {

	if addrs, ok := HasValidOwnership(issuerPub, address...); !ok {
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state invalid ownership of templates (%s)", addrs)}
	}

	state, err := s.Context().GetState(address)
	if err != nil {
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
	}

	templates := make([]*template_pb.Template, 0)
	for _, value := range state {
		var template template_pb.Template
		err := proto.Unmarshal(value, &template)
		if err != nil {
			return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal proto data (%s)", err)}
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// HasValidOwnership using the first 30 bytes of a public key this func will
// verify that the pub key is the first 30 bytes of each address indicating ownership.
// If validation for one address fails, the entire validation process is will fail, array of address
// that failed validation is returned along with a bool
func HasValidOwnership(issuerPub string, addresses ...string) ([]string, bool) {
	invalid := make([]string, 0)
	prefix := issuerPub[0:30]
	for _, address := range addresses {
		if address[0:30] != prefix {
			invalid = append(invalid, address)
		}
	}

	return invalid, len(invalid) == 0
}

// TemplateStateAddress ...
func TemplateStateAddress(issuerPub, name string, version *template_pb.Version) string {
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
