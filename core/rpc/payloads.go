package rpc

import (
	"encoding/json"
	"github.com/BadgeForce/credential-template-engine/core/proto"
	"github.com/BadgeForce/credential-template-engine/core/state"
)

// PayloadDecoder helper struct that decodes JSON payloads from RPC request
type PayloadDecoder struct {
	req *credential_template_engine_pb.RPCRequest
}

// DeletePayload expected payload for delete template requests to handler
type DeletePayload struct {
	Addresses []string `json:"addresses"`
}

// Bytes returns Params from RPCReq serialized to []byte
func (p *PayloadDecoder) Bytes() []byte {
	return []byte(p.req.Params)
}

// UnmarshalCreate un-marshals RPC Payloads for request with method CREATE in expected JSON format
func (p *PayloadDecoder) UnmarshalCreate() (*state.CredentialTemplate, error) {
	var template state.CredentialTemplate
	err := json.Unmarshal(p.Bytes(), &template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// UnmarshalCreate un-marshals RPC Payloads for request with method UPDATE in expected JSON format
func (p *PayloadDecoder) UnmarshalUpdate() (*state.CredentialTemplate, error) {
	var template state.CredentialTemplate
	err := json.Unmarshal(p.Bytes(), &template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}

// UnmarshalCreate un-marshals RPC Payloads for request with method DELETE in expected JSON format
func (p *PayloadDecoder) UnmarshalDelete() (*DeletePayload, error) {
	var payload DeletePayload
	err := json.Unmarshal(p.Bytes(), &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func NewPayloadDecoder(req *credential_template_engine_pb.RPCRequest) *PayloadDecoder {
	return &PayloadDecoder{req}
}