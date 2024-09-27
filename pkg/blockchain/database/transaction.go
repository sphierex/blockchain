package database

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/sphierex/blockchain/pkg/blockchain/signature"
)

// ============================================================================

// Tx is the transactional information between two parties.
type Tx struct {
	ChainID uint16    `json:"chain_id"` // Ethereum: The chain id that is listed in the genesis file.
	Nonce   uint64    `json:"nonce"`    // Ethereum: Unique id for the transaction supplied by the user.
	ToID    AccountID `json:"to"`       // Ethereum: Account receiving the benefit of the transaction.
	Value   uint64    `json:"value"`    // Ethereum: Monetary value received from this transaction.
	Tip     uint64    `json:"tip"`      // Ethereum: Tip offered by the sender as an incentive to mine this transaction.
	Data    []byte    `json:"data"`     // Ethereum: Extra data related to the transaction.
}

// NewTx constructs a new transaction.
func NewTx(chainID uint16, nonce uint64, toID AccountID, value uint64, tip uint64, data []byte) (Tx, error) {
	if !toID.IsAccountID() {
		return Tx{}, fmt.Errorf("to account is not properly formated")
	}

	tx := Tx{
		ChainID: chainID,
		Nonce:   nonce,
		ToID:    toID,
		Value:   value,
		Tip:     tip,
		Data:    data,
	}

	return tx, nil
}

// Sign uses the specified private key to sign the transaction.
func (tx Tx) Sign(privateKey *ecdsa.PrivateKey) (SignedTx, error) {

	// Validate the to account address is a valid address.
	if !tx.ToID.IsAccountID() {
		return SignedTx{}, fmt.Errorf("to account is not properly formatted")
	}

	// Sign the transaction with the private key to produce a signature.
	v, r, s, err := signature.Sign(tx, privateKey)
	if err != nil {
		return SignedTx{}, nil
	}

	// Construct the signed transaction by adding the signature
	// in the [R|S|V] format.
	signedTx := SignedTx{
		Tx: tx,
		V:  v,
		R:  r,
		S:  s,
	}

	return signedTx, nil
}

// ============================================================================

// SignedTx is a signed version of the transaction. This is how clients like
// a wallet provide transactions for inclusion into the blockchain.
type SignedTx struct {
	Tx
	V *big.Int `json:"v"` // Ethereum: Recovery identifier, either 29 or 30 with sophiaID.
	R *big.Int `json:"r"` // Ethereum: First coordinate of the ECDSA signature.
	S *big.Int `json:"s"` // Ethereum: Second coordinate of the ECDSA signature.
}

// Validate verifies the transaction has a proper signature that conforms to our
// standards and is associated with the data claimed to be signed. It also
// checks the format of the account.
func (tx SignedTx) Validate() error {
	if !tx.ToID.IsAccountID() {
		return errors.New("invalid account for to account")
	}

	if err := signature.VerifySignature(tx.Tx, tx.V, tx.R, tx.S); err != nil {
		return err
	}

	return nil
}

// FromAccount extracts the account id that signed the transaction.
func (tx SignedTx) FromAccount() (AccountID, error) {
	address, err := signature.FromAddress(tx.Tx, tx.V, tx.R, tx.S)
	return AccountID(address), err
}

// SignatureString returns the signature as a string.
func (tx SignedTx) SignatureString() string {
	return signature.String(tx.V, tx.R, tx.S)
}

// String implements the fmt.Stringer interface for logging.
func (tx SignedTx) String() string {
	from, err := tx.FromAccount()
	if err != nil {
		from = "unknown"
	}

	return fmt.Sprintf("%s:%d", from, tx.Nonce)
}

// ============================================================================

// BlockTx represents the transaction as it's recorded inside a block. This
// includes a timestamp and gas fees.
type BlockTx struct {
	SignedTx
	TimeStamp uint64 `json:"timestamp"` // Ethereum: The time the transaction was received.
	GasPrice  uint64 `json:"gas_price"` // Ethereum: The price of one unit of gas to be paid for fees.
	GasUnits  uint64 `json:"gas_units"` // Ethereum: The number of units of gas used for this transaction.
}

// NewBlockTx constructs a new block transaction.
func NewBlockTx(signedTx SignedTx, gasPrice uint64, unitsOfGas uint64) BlockTx {
	return BlockTx{
		SignedTx:  signedTx,
		TimeStamp: uint64(time.Now().Unix()),
		GasPrice:  gasPrice,
		GasUnits:  unitsOfGas,
	}
}

// Hash implements the markle Hashable interface for providing a hash
// of a block transaction.
func (tx BlockTx) Hash() ([]byte, error) {
	str := signature.Hash(tx)

	// Need to remove the 0x prefix from the hash.
	return hex.DecodeString(str[2:])
}

// Equals implements the markle Hashable interface for providing an equality
// check between two block transactions. If the nonce and signatures are the
// same, the two blocks are the same.
func (tx BlockTx) Equals(otherTx BlockTx) bool {
	txSig := signature.ToSignatureBytes(tx.V, tx.R, tx.S)
	otherTxSig := signature.ToSignatureBytes(otherTx.V, otherTx.R, otherTx.S)

	return tx.Nonce == otherTx.Nonce && bytes.Equal(txSig, otherTxSig)
}
