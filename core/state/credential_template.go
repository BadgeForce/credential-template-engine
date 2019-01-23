package state

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	utils "github.com/BadgeForce/sawtooth-utils"
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
	return TemplateStateAddress(c.Name, c.Owner, c.Version)
}

func TemplateStateAddress(owner, name, version string) string {
	o := utils.NewPart(owner, 0, 30)
	n := utils.NewPart(name, 0, 30)
	v := utils.NewPart(version, 0, 4)

	addressParts := []*utils.AddressPart{o, n, v}
	address, _ := utils.NewAddress(NameSpaceMngr.NameSpaces[0]).AddParts(addressParts...).Build()
	return address
}

func NewCredentialTemplate(name, owner, version, data string) *CredentialTemplate {
	checkSum := md5.Sum([]byte(data))
	return &CredentialTemplate{
		Name:      name,
		Version:   version,
		CreatedAt: time.Now().Unix(),
		Owner:     owner,
		Data:      data,
		CheckSum:  fmt.Sprintf("%x", checkSum),
	}
}
