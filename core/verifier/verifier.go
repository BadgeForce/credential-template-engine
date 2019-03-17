package verifier

import (
	"fmt"
	"github.com/BadgeForce/sawtooth-utils"
	"github.com/BadgeForce/sawtooth-utils/protos/templates_pb"
	"github.com/golang/protobuf/proto"
)

// VerifyTemplate ...
func VerifyTemplate(txtSignerPub string, template *templates_pb.Template) error {
	b, err := proto.Marshal(template.GetData())
	if err != nil {
		return fmt.Errorf("error: could not marshal proto (%s)", err)
	}
	expectedHash := template.GetVerification().GetProofOfIntegrityHash()
	if hash, ok := utils.VerifyPOIHash(b, expectedHash); !ok {
		return fmt.Errorf("error: proof of integrity hash invalid got (%s) want (%s)", hash, expectedHash)
	}

	issuerPub := template.GetData().GetIssuerPub()
	sig := template.GetVerification().GetSignature()

	if ok := utils.VerifySig(sig, []byte(txtSignerPub), b, false); !ok {
		return fmt.Errorf("error: transaction signer must also be owner who signs template (%s)", txtSignerPub)
	} else if txtSignerPub != issuerPub {
		return fmt.Errorf("error: transaction signer public key must match template issuer got (%s) want (%s)", issuerPub, txtSignerPub)
	}

	return nil
}

// HasValidOwnership using the first 30 bytes of a public key this func will
// verify that the pub key is the first 30 bytes of each address indicating ownership.
// If validation for one address fails, the entire validation process is will fail, array of address
// that failed validation is returned along with a bool
func HasValidOwnership(issuerPub string, addresses ...string) ([]string, bool) {
	invalid := make([]string, 0)
	prefix := issuerPub[0:30]
	for _, address := range addresses {
		if address[6:30] != prefix {
			invalid = append(invalid, address)
		}
	}

	return invalid, len(invalid) == 0
}
