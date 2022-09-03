package postgres

import "github.com/brronsuite/bronwallet/walletdb"

type Fixture interface {
	DB() walletdb.DB
	Dump() (map[string]interface{}, error)
}
