package state

import (
	"strings"
	"encoding/hex"
	"crypto/sha512"
)

const (
	TEMPLATE_PREFIX = "credential:templates"
)

var (
	TemplatePrefix = ComputePrefix(TEMPLATE_PREFIX)
)

// ComputePrefix returns namespace prefix of 6 bytes
func ComputePrefix(prefix string) string {
	return Hexdigest(prefix)[:6]
}

// Hexdigest
func Hexdigest(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

// MakeAddress . . .
func MakeAddress(namespacePrefix, namespaceSuffix string) string {
	return namespacePrefix + Hexdigest(namespaceSuffix)[:64]
}

// MakeIdentifierAddress . . .
func MakeIdentifierAddress(prefix, owner, postfix string) string {
	return prefix + Hexdigest(owner)[:32] + Hexdigest(postfix)[:32]
}

