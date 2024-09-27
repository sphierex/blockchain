package signature

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// ZeroHash represents a hash code of zeros.
const ZeroHash string = "0x0000000000000000000000000000000000000000000000000000000000000000"

// sophiaID is an arbitrary number for signing messages. This will make it
// clear that the signature comes from the Ardan blockchain.
// Ethereum and Bitcoin do this as well, but they use the value of 27.
const sophiaID = 44

// ============================================================================

// Hash returns a unique string for the value.
func Hash(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ZeroHash
	}

	hash := sha256.Sum256(data)
	return "0x" + hex.EncodeToString(hash[:])
}

// Sign uses the specified private key to sign the transaction.
func Sign(value any, privateKey *ecdsa.PrivateKey) (v, r, s *big.Int, err error) {

	// Prepare the transaction for signing.
	data, err := stamp(value)
	if err != nil {
		return nil, nil, nil, err
	}

	// Sign the hash with private key to produce a signature.
	sig, err := crypto.Sign(data, privateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Convert the 65 byte signature into the [R|S|V] format.
	v, r, s = toSignatureValues(sig)

	return v, r, s, nil
}

// VerifySignature verifies the signature confirms to our standards and
// is associated with the data claimed to be signed.
func VerifySignature(value any, v, r, s *big.Int) error {
	// Check the recovery id is either 0 or 1.
	uintV := v.Uint64() - sophiaID
	if uintV != 0 && uintV != 1 {
		return errors.New("invalid recovery id")
	}

	// Check the signature values are valid.
	if !crypto.ValidateSignatureValues(byte(uintV), r, s, false) {
		return errors.New("invalid signature values")
	}

	// Prepare the transaction for recovery and validation.
	tran, err := stamp(value)
	if err != nil {
		return err
	}

	// Convert the [R|S|V] format into the original 65 bytes.
	sig := ToSignatureBytes(v, r, s)

	// Capture the uncompressed public key associated with this signature.
	sigPublicKey, err := crypto.Ecrecover(tran, sig)

	// Check that the given public key created the signature over the data.
	rs := sig[:crypto.RecoveryIDOffset]
	if !crypto.VerifySignature(sigPublicKey, tran, rs) {
		return errors.New("invalid signature")
	}

	return nil
}

// FromAddress extracts the address for the account that signed the transaction.
func FromAddress(value any, v, r, s *big.Int) (string, error) {

	// Prepare the transaction for public key extraction.
	data, err := stamp(value)
	if err != nil {
		return "", err
	}

	// Convert the [R|S|V] format into the original 65 bytes.
	sig := ToSignatureBytes(v, r, s)

	// Validate the signature since there can be conversion issues
	// between [R|S|V] to []bytes, Leading 0's are truncated by big package.
	// publicKey, err := crypto.SigToPub(data, sig)
	var sigPublicKey []byte
	{
		sigPublicKey, err := crypto.Ecrecover(data, sig)
		if err != nil {
			return "", err
		}

		rs := sig[:crypto.RecoveryIDOffset]
		if !crypto.VerifySignature(sigPublicKey, data, rs) {
			return "", errors.New("invalid signature")
		}
	}

	// Capture the public key associated with this signature.
	// x, y := elliptic.Unmarshal(crypto.S256(), sigPublicKey)
	x, y := crypto.S256().Unmarshal(sigPublicKey)
	publicKey := ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}

	// Extract the account address from the public key.
	return crypto.PubkeyToAddress(publicKey).String(), nil
}

// String returns the signature as a string.
// as SignatureString.
func String(v, r, s *big.Int) string {
	return "0x" + hex.EncodeToString(ToSignatureBytesWithSophiaID(v, r, s))
}

// ============================================================================

// stamp returns a hash of 32 bytes that represents this transaction with
// the Sophia stamp embedded into the final hash.
func stamp(value any) ([]byte, error) {

	// Marshal the data.
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// Hash the transaction data into a 32 bytes array. This will provide
	// a data length consistency with all transactions.
	txHash := crypto.Keccak256Hash(data)

	// Convert the stamp into a slice of bytes. This stamp is
	// used so signatures we produce when signing transactions
	// are always unique to the Sophia blockchain.
	stamp := []byte("\x19Sophia Signed Message:\n32")

	// Hash the stamp and txHash together in a final 32 byte array
	// the represents the transaction data.
	tran := crypto.Keccak256Hash(stamp, txHash.Bytes())

	return tran.Bytes(), nil
}

// toSignatureValues converts the signature into the r, s, v values.
func toSignatureValues(sig []byte) (v, r, s *big.Int) {
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + sophiaID})

	return v, r, s
}

// ToSignatureBytes converts the r, s, v values into a slice of bytes
// with the removal of the sophiaID.
func ToSignatureBytes(v, r, s *big.Int) []byte {
	sig := make([]byte, crypto.SignatureLength)

	rBytes := r.Bytes()
	if len(rBytes) == 31 {
		copy(sig[1:], rBytes)
	} else {
		copy(sig, rBytes)
	}

	sBytes := s.Bytes()
	if len(sBytes) == 31 {
		copy(sig[33:], sBytes)
	} else {
		copy(sig[32:], sBytes)
	}

	sig[64] = byte(v.Uint64() - sophiaID)

	return sig
}

// ToSignatureBytesWithSophiaID converts the r, s, v values into a slice of bytes
// keeping the Sophia id.
func ToSignatureBytesWithSophiaID(v, r, s *big.Int) []byte {
	sig := ToSignatureBytes(v, r, s)
	sig[64] = byte(v.Uint64())

	return sig
}
