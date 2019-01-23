package processor

import (
	"sync"

	"github.com/BadgeForce/credential-template-engine/core/rpc"
	"github.com/BadgeForce/credential-template-engine/core/state"
	utils "github.com/BadgeForce/sawtooth-utils"
	"github.com/rberg2/sawtooth-go-sdk/processor"
)

//TODO: move FamilyName to configuration

// FamilyName processor family name
const FamilyName = "credential-templates"

//TODO: move NameSpaces to configuration

// Namespaces transaction processor namespaces
var Namespaces = state.NameSpaceMngr.NameSpaces

//TODO: move FamilyVersions to configuration

// FamilyVersions transaction processor versions
var FamilyVersions = []string{"1.0"}

var once sync.Once
var transactionHandler *utils.TransactionHandler
var transactionProcessor *utils.TransactionProcessor

// TransactionProcessor instantiates transaction processor once and return it
func TransactionProcessor(validator string) *processor.TransactionProcessor {
	once.Do(func() {
		transactionHandler = utils.NewTransactionHandler(rpc.NewClient(), FamilyName, FamilyVersions, Namespaces)
		transactionProcessor = utils.NewTransactionProcessor(validator, transactionHandler)
	})

	return transactionProcessor
}
