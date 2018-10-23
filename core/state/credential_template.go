package state

import (
	"time"
	"crypto/md5"
	"fmt"
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
