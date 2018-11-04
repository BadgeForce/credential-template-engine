package state

import (
	"time"
	"crypto/md5"
	"fmt"
	"encoding/json"
	"github.com/BadgeForce/badgeforce-chain-node/core/common"
)

type CredentialTemplate struct {
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	Version   string `json:"version"`
	Owner     string `json:"owner"`
	Data      string `json:"data"`
	CheckSum  string `json:"check_sum"`
}

func (c *CredentialTemplate) VerifyChecksum() bool {
	if fmt.Sprintf("%x", md5.Sum([]byte(c.Data))) != c.CheckSum {
		return false
	}
	return true
}

func (c *CredentialTemplate) UpdateChecksum() {
	c.CheckSum = fmt.Sprintf("%x", md5.Sum([]byte(c.Data)))
}

func (c *CredentialTemplate) AsBytes() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("error invalid credential temlate (%s)", err)
	}

	return b, nil
}

func (c *CredentialTemplate) StateAddress() string {
	owner := common.NewPart(c.Owner, 0, 30)
	name := common.NewPart(c.Name, 0, 30)
	version := common.NewPart(c.Version, 0, 4)

	addressParts := []*common.AddressPart{owner, name, version}
	address, _ := common.NewAddress(NameSpaceMngr.NameSpaces[0]).AddParts(addressParts...).Build()
	return address
}

func NewCredentialTemplate(name, owner, version, data string) *CredentialTemplate {
	checkSum := md5.Sum([]byte(data))
	return &CredentialTemplate{
		Name: name,
		Version: version,
		CreatedAt: time.Now().Unix(),
		Owner:     owner,
		Data:      data,
		CheckSum:  fmt.Sprintf("%x", checkSum),
	}
}
