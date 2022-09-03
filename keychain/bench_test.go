package keychain

import (
	"testing"

	"github.com/brronsuite/brond/bronec"
	"github.com/stretchr/testify/require"
)

func BenchmarkDerivePrivKey(t *testing.B) {
	cleanUp, wallet, err := createTestbronwallet(
		CoinTypeBrocoin,
	)
	if err != nil {
		t.Fatalf("unable to create wallet: %v", err)
	}

	keyRing := NewbronwalletKeyRing(wallet, CoinTypeBrocoin)

	defer cleanUp()

	var (
		privKey *bronec.PrivateKey
	)

	keyDesc := KeyDescriptor{
		KeyLocator: KeyLocator{
			Family: KeyFamilyMultiSig,
			Index:  1,
		},
	}

	t.ReportAllocs()
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		privKey, err = keyRing.DerivePrivKey(keyDesc)
	}
	require.NoError(t, err)
	require.NotNil(t, privKey)
}
