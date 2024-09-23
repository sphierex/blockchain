package genesis

import (
	"encoding/json"
	"os"
	"time"
)

// Genesis represents the genesis file.
type Genesis struct {
	Date          time.Time         `json:"date"`
	ChainID       uint16            `json:"chain_id"`
	TransPerBlock uint16            `json:"trans_per_block"`
	Difficulty    uint16            `json:"difficulty"`
	MiningReward  uint64            `json:"mining_reward"`
	GasPrice      uint64            `json:"gas_price"`
	Balances      map[string]uint64 `json:"balances"`
}

// Load opens and consumes the genesis file.
func Load() (Genesis, error) {
	path := "zblock/genesis.json"
	content, err := os.ReadFile(path)
	if err != nil {
		return Genesis{}, err
	}

	var genesis Genesis
	err = json.Unmarshal(content, &genesis)
	if err != nil {
		return Genesis{}, err
	}

	return genesis, nil
}
