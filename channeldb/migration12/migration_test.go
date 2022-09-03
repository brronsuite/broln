package migration12_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/brronsuite/broln/channeldb/migration12"
	"github.com/brronsuite/broln/channeldb/migtest"
	"github.com/brronsuite/broln/kvdb"
	"github.com/brronsuite/broln/lntypes"
)

var (
	// invoiceBucket is the name of the bucket within the database that
	// stores all data related to invoices no matter their final state.
	// Within the invoice bucket, each invoice is keyed by its invoice ID
	// which is a monotonically increasing uint32.
	invoiceBucket = []byte("invoices")

	preimage = lntypes.Preimage{
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
	}

	hash = preimage.Hash()

	beforeInvoice0Htlcs = []byte{
		0x0b, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72,
		0x6c, 0x64, 0x09, 0x62, 0x79, 0x65, 0x20, 0x77, 0x6f, 0x72,
		0x6c, 0x64, 0x06, 0x70, 0x61, 0x79, 0x72, 0x65, 0x71, 0x00,
		0x00, 0x00, 0x20, 0x00, 0x00, 0x4e, 0x94, 0x91, 0x4f, 0x00,
		0x00, 0x0f, 0x01, 0x00, 0x00, 0x00, 0x0e, 0x77, 0xc4, 0xd3,
		0xd5, 0x00, 0x00, 0x00, 0x00, 0xfe, 0x20, 0x0f, 0x01, 0x00,
		0x00, 0x00, 0x0e, 0x77, 0xd5, 0xc8, 0x1c, 0x00, 0x00, 0x00,
		0x00, 0xfe, 0x20, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x03, 0xe8, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xa4,
	}

	afterInvoice0Htlcs = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb8, 0x00, 0x0b,
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c,
		0x64, 0x01, 0x06, 0x70, 0x61, 0x79, 0x72, 0x65, 0x71, 0x02,
		0x0f, 0x01, 0x00, 0x00, 0x00, 0x0e, 0x77, 0xc4, 0xd3, 0xd5,
		0x00, 0x00, 0x00, 0x00, 0xfe, 0x20, 0x03, 0x0f, 0x01, 0x00,
		0x00, 0x00, 0x0e, 0x77, 0xd5, 0xc8, 0x1c, 0x00, 0x00, 0x00,
		0x00, 0xfe, 0x20, 0x04, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x05, 0x05, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x06, 0x06, 0x20, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x07, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xe8, 0x08, 0x04, 0x00,
		0x00, 0x00, 0x20, 0x09, 0x08, 0x00, 0x00, 0x4e, 0x94, 0x91,
		0x4f, 0x00, 0x00, 0x0a, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, 0x00, 0x0c,
		0x01, 0x03, 0x0d, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0xa4,
	}

	testHtlc = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x41,
		0x01, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08,
		0x03, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
		0x05, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x64,
		0x07, 0x04, 0x00, 0x00, 0x00, 0x58, 0x09, 0x08, 0x00, 0x13,
		0xbc, 0xbf, 0x72, 0x4e, 0x1e, 0x00, 0x0b, 0x08, 0x00, 0x17,
		0xaf, 0x4c, 0x22, 0xc4, 0x24, 0x00, 0x0d, 0x04, 0x00, 0x00,
		0x23, 0x1d, 0x0f, 0x01, 0x02,
	}

	beforeInvoice1Htlc = append([]byte{
		0x0b, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72,
		0x6c, 0x64, 0x09, 0x62, 0x79, 0x65, 0x20, 0x77, 0x6f, 0x72,
		0x6c, 0x64, 0x06, 0x70, 0x61, 0x79, 0x72, 0x65, 0x71, 0x00,
		0x00, 0x00, 0x20, 0x00, 0x00, 0x4e, 0x94, 0x91, 0x4f, 0x00,
		0x00, 0x0f, 0x01, 0x00, 0x00, 0x00, 0x0e, 0x77, 0xc4, 0xd3,
		0xd5, 0x00, 0x00, 0x00, 0x00, 0xfe, 0x20, 0x0f, 0x01, 0x00,
		0x00, 0x00, 0x0e, 0x77, 0xd5, 0xc8, 0x1c, 0x00, 0x00, 0x00,
		0x00, 0xfe, 0x20, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x03, 0xe8, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xa4,
	}, testHtlc...)

	afterInvoice1Htlc = append([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb8, 0x00, 0x0b,
		0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c,
		0x64, 0x01, 0x06, 0x70, 0x61, 0x79, 0x72, 0x65, 0x71, 0x02,
		0x0f, 0x01, 0x00, 0x00, 0x00, 0x0e, 0x77, 0xc4, 0xd3, 0xd5,
		0x00, 0x00, 0x00, 0x00, 0xfe, 0x20, 0x03, 0x0f, 0x01, 0x00,
		0x00, 0x00, 0x0e, 0x77, 0xd5, 0xc8, 0x1c, 0x00, 0x00, 0x00,
		0x00, 0xfe, 0x20, 0x04, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x05, 0x05, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x06, 0x06, 0x20, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42,
		0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x07, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xe8, 0x08, 0x04, 0x00,
		0x00, 0x00, 0x20, 0x09, 0x08, 0x00, 0x00, 0x4e, 0x94, 0x91,
		0x4f, 0x00, 0x00, 0x0a, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, 0x00, 0x0c,
		0x01, 0x03, 0x0d, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0xa4,
	}, testHtlc...)
)

type migrationTest struct {
	name            string
	beforeMigration func(kvdb.RwTx) error
	afterMigration  func(kvdb.RwTx) error
}

var migrationTests = []migrationTest{
	{
		name:            "no invoices",
		beforeMigration: func(kvdb.RwTx) error { return nil },
		afterMigration:  func(kvdb.RwTx) error { return nil },
	},
	{
		name:            "zero htlcs",
		beforeMigration: genBeforeMigration(beforeInvoice0Htlcs),
		afterMigration:  genAfterMigration(afterInvoice0Htlcs),
	},
	{
		name:            "one htlc",
		beforeMigration: genBeforeMigration(beforeInvoice1Htlc),
		afterMigration:  genAfterMigration(afterInvoice1Htlc),
	},
}

// genBeforeMigration creates a closure that inserts an invoice serialized under
// the old format under the test payment hash.
func genBeforeMigration(beforeBytes []byte) func(kvdb.RwTx) error {
	return func(tx kvdb.RwTx) error {
		invoices, err := tx.CreateTopLevelBucket(
			invoiceBucket,
		)
		if err != nil {
			return err
		}

		return invoices.Put(hash[:], beforeBytes)
	}
}

// genAfterMigration creates a closure that verifies the tlv invoice migration
// succeeded, but comparing the resulting encoding of the invoice to the
// expected serialization. In addition, the decoded invoice is compared against
// the expected invoice for equality.
func genAfterMigration(afterBytes []byte) func(kvdb.RwTx) error {
	return func(tx kvdb.RwTx) error {
		invoices := tx.ReadWriteBucket(invoiceBucket)
		if invoices == nil {
			return fmt.Errorf("invoice bucket not found")
		}

		// Fetch the new invoice bytes and check that they match our
		// expected serialization.
		invoiceBytes := invoices.Get(hash[:])
		if !bytes.Equal(invoiceBytes, afterBytes) {
			return fmt.Errorf("invoice bytes mismatch, "+
				"want: %x, got: %x",
				invoiceBytes, afterBytes)
		}

		return nil
	}
}

// TestTLVInvoiceMigration executes a suite of migration tests for moving
// invoices to use TLV for their bodies. In the process, feature bits and
// payment addresses are added to the invoice while the receipt field is
// dropped. We test a few different invoices with a varying number of HTLCs, as
// well as the case where there are no invoices present.
//
// NOTE: The test vectors each include a receipt that is not present on the
// final struct, but verifies that the field is properly removed.
func TestTLVInvoiceMigration(t *testing.T) {
	for _, test := range migrationTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			migtest.ApplyMigration(
				t,
				test.beforeMigration,
				test.afterMigration,
				migration12.MigrateInvoiceTLV,
				false,
			)
		})
	}
}
