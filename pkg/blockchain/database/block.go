package database

import "github.com/sphierex/blockchain/pkg/blockchain/merkle"

// =============================================================================

// BlockData represents what can be serialized to disk and over the network.
type BlockData struct {
	Hash   string      `json:"hash"`
	Header BlockHeader `json:"block"`
	Trans  []BlockTx   `json:"trans"`
}

// =============================================================================

// BlockHeader represents common information required for each block.
type BlockHeader struct {
	Number        uint64    `json:"number"`
	PrevBlockHash string    `json:"prev_block_hash"`
	TimeStamp     uint64    `json:"timestamp"`
	BeneficiaryID AccountID `json:"beneficiary"`
	Difficulty    uint16    `json:"difficulty"`
	MiningReward  uint64    `json:"mining_reward"`
	StateRoot     string    `json:"state_root"`
	TransRoot     string    `json:"trans_root"`
	Nonce         uint64    `json:"nonce"`
}

// Block represents a group of transactions batched together.
type Block struct {
	Header     BlockHeader
	MerkleTree *merkle.Tree[BlockTx]
}
