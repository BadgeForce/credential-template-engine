package state

import (
	"github.com/rberg2/sawtooth-go-sdk/processor"
	"time"
	"encoding/json"
	"github.com/rberg2/sawtooth-go-sdk/logging"
	"fmt"
)

var logger = logging.Get()

type State struct {
	context *processor.Context
}

type TemplateSavedReceipt struct {
	Date            int64  `json:"date"`
	TemplateName    string `json:"template_name"`
	TemplateVersion string `json:"template_version"`
	StateAddress    string `json:"state_address"`
}

func NewTemplateSavedReceipt(name, version, address string) (*TemplateSavedReceipt, []byte, error) {
	receipt := &TemplateSavedReceipt{
		Date: time.Now().Unix(),
		TemplateName: name,
		TemplateVersion: version,
		StateAddress: address,
	}

	b, err := json.Marshal(receipt)
	return receipt, b, err
}

func NewState(context *processor.Context) *State {
	return &State{context: context}
}

func (s *State) SaveTemplate(template *CredentialTemplate) error {
	address := MakeIdentifierAddress(TemplatePrefix, template.Owner, template.Name)
	_, receiptBytes, err := NewTemplateSavedReceipt(template.Name, template.Version, address)
	if err != nil {
		logger.Warnf("unable to generate transaction receipt for template saved (%s)", err)
	}

	b, err := json.Marshal(template)
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("invalid credential temlate (%s)", err)}
	}

	_, err = s.context.SetState(map[string][]byte{
		address: b,
	})
	if err != nil {
		return &processor.InvalidTransactionError{Msg: fmt.Sprintf("unable to save credential temlate (%s)", err)}
	}

	err = s.context.AddReceiptData(receiptBytes)
	if err != nil {
		logger.Warnf("unable to add transaction receipt for template saved (%s)", err)
	}

	return nil
}
