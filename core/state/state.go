package state

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/BadgeForce/credential-template-engine/core/template_pb"
	utils "github.com/BadgeForce/sawtooth-utils"
	"github.com/rberg2/sawtooth-go-sdk/logging"
	"github.com/rberg2/sawtooth-go-sdk/processor"
)

//TODO: move CredentialTemplatePrefix to configuration

// CredentialTemplatePrefix ...
const CredentialTemplatePrefix = "credential:templates"

var (
	logger        = logging.Get()
	NameSpaceMngr *utils.NamespaceMngr
)

type State struct {
	instance *utils.State
}

type TransactionReceipt struct {
	Date            int64  `json:"date"`
	TemplateName    string `json:"template_name"`
	TemplateVersion string `json:"template_version"`
	StateAddress    string `json:"state_address"`
	Method          string `json:"method"`
}

func (s *State) NewTemplateSavedReceipt(name, version, address string) (*TransactionReceipt, []byte, error) {
	receipt := &TransactionReceipt{
		Date:            time.Now().Unix(),
		TemplateName:    name,
		TemplateVersion: version,
		StateAddress:    address,
		Method:          template_pb.Method_CREATE.String(),
	}

	b, err := json.Marshal(receipt)
	return receipt, b, err
}
func (s *State) NewTemplateDeleteReceipt(name, version, address string) (*TransactionReceipt, []byte, error) {
	receipt := &TransactionReceipt{
		Date:            time.Now().Unix(),
		TemplateName:    name,
		TemplateVersion: version,
		StateAddress:    address,
		Method:          template_pb.Method_DELETE.String(),
	}

	b, err := json.Marshal(receipt)
	return receipt, b, err
}

func (s *State) Context() *processor.Context {
	return s.instance.Context
}

func (s *State) Save(template *CredentialTemplate) error {
	address := template.StateAddress()
	_, receiptBytes, err := s.NewTemplateSavedReceipt(template.Name, template.Version, address)
	if err != nil {
		logger.Warnf("unable to generate transaction receipt for template saved (%s)", err)
	}

	b, err := template.AsBytes()
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("(%s)", err)}
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

func (s *State) Delete(addresses ...string) error {
	_, err := s.Context().DeleteState(addresses)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("unable to delete credential temlate (%s)", err)}
	}

	for _, address := range addresses {
		_, receiptBytes, err := s.NewTemplateDeleteReceipt("", "", address)
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

func (s *State) GetTemplates(address ...string) ([]CredentialTemplate, error) {
	state, err := s.Context().GetState(address)
	if err != nil {
		return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not get state (%s)", err)}
	}

	templates := make([]CredentialTemplate, 0)
	for _, value := range state {
		var template CredentialTemplate
		err := json.Unmarshal(value, &template)
		if err != nil {
			return nil, &processor.InvalidTransactionError{Msg: fmt.Sprintf("could not unmarshal json data (%s)", err)}
		}
		templates = append(templates, template)
	}

	return templates, nil
}

func NewTemplateState(context *processor.Context) *State {
	return &State{utils.NewStateInstance(context)}
}

func init() {
	NameSpaceMngr = utils.NewNamespaceMngr().RegisterNamespaces(CredentialTemplatePrefix)
}
