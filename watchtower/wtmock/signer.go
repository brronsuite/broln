package wtmock

import (
	"sync"

	"github.com/brronsuite/broln/input"
	"github.com/brronsuite/broln/keychain"
	"github.com/brronsuite/brond/bronec"
	"github.com/brronsuite/brond/txscript"
	"github.com/brronsuite/brond/wire"
)

// MockSigner is an input.Signer that allows one to add arbitrary private keys
// and sign messages by passing the assigned keychain.KeyLocator.
type MockSigner struct {
	mu sync.Mutex

	index uint32
	keys  map[keychain.KeyLocator]*bronec.PrivateKey
}

// NewMockSigner returns a fresh MockSigner.
func NewMockSigner() *MockSigner {
	return &MockSigner{
		keys: make(map[keychain.KeyLocator]*bronec.PrivateKey),
	}
}

// SignOutputRaw signs an input on the passed transaction using the input index
// in the sign descriptor. The returned signature is the raw DER-encoded
// signature without the signhash flag.
func (s *MockSigner) SignOutputRaw(tx *wire.MsgTx,
	signDesc *input.SignDescriptor) (input.Signature, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	witnessScript := signDesc.WitnessScript
	amt := signDesc.Output.Value

	privKey, ok := s.keys[signDesc.KeyDesc.KeyLocator]
	if !ok {
		panic("cannot sign w/ unknown key")
	}

	sig, err := txscript.RawTxInWitnessSignature(
		tx, signDesc.SigHashes, signDesc.InputIndex, amt,
		witnessScript, signDesc.HashType, privKey,
	)
	if err != nil {
		return nil, err
	}

	return bronec.ParseDERSignature(sig[:len(sig)-1], bronec.S256())
}

// ComputeInputScript is not implemented.
func (s *MockSigner) ComputeInputScript(tx *wire.MsgTx,
	signDesc *input.SignDescriptor) (*input.Script, error) {
	panic("not implemented")
}

// AddPrivKey records the passed privKey in the MockSigner's registry of keys it
// can sign with in the future. A unique key locator is returned, allowing the
// caller to sign with this key when presented via an input.SignDescriptor.
func (s *MockSigner) AddPrivKey(privKey *bronec.PrivateKey) keychain.KeyLocator {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyLoc := keychain.KeyLocator{
		Index: s.index,
	}
	s.index++

	s.keys[keyLoc] = privKey

	return keyLoc
}

// Compile-time constraint ensuring the MockSigner implements the input.Signer
// interface.
var _ input.Signer = (*MockSigner)(nil)
