package database

import (
	"sync"

	"github.com/sphierex/blockchain/pkg/blockchain/genesis"
)

// Database manages data related to accounts who have transacted on the blockchain.
type Database struct {
	mu      sync.RWMutex
	genesis genesis.Genesis
	// latest Block
	accounts map[AccountID]Account
}
