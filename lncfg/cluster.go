package lncfg

import (
	"context"
	"fmt"
	"os"

	"github.com/brronsuite/broln/cluster"
)

const (
	// DefaultEtcdElectionPrefix is used as election prefix if none is provided
	// through the config.
	DefaultEtcdElectionPrefix = "/leader/"
)

// Cluster holds configuration for clustered broln.
type Cluster struct {
	EnableLeaderElection bool `long:"enable-leader-election" description:"Enables leader election if set."`

	LeaderElector string `long:"leader-elector" choice:"etcd" description:"Leader elector to use. Valid values: \"etcd\"."`

	EtcdElectionPrefix string `long:"etcd-election-prefix" description:"Election key prefix when using etcd leader elector. Defaults to \"/leader/\"."`

	ID string `long:"id" description:"Identifier for this node inside the cluster (used in leader election). Defaults to the hostname."`
}

// DefaultCluster creates and returns a new default DB config.
func DefaultCluster() *Cluster {
	hostname, _ := os.Hostname()
	return &Cluster{
		LeaderElector:      cluster.EtcdLeaderElector,
		EtcdElectionPrefix: DefaultEtcdElectionPrefix,
		ID:                 hostname,
	}
}

// MakeLeaderElector is a helper method to construct the concrete leader elector
// based on the current configuration.
func (c *Cluster) MakeLeaderElector(electionCtx context.Context, db *DB) (
	cluster.LeaderElector, error) {

	if c.LeaderElector == cluster.EtcdLeaderElector {
		return cluster.MakeLeaderElector(
			electionCtx, c.LeaderElector, c.ID,
			c.EtcdElectionPrefix, db.Etcd,
		)
	}

	return nil, fmt.Errorf("unsupported leader elector")
}

// Validate validates the Cluster config.
func (c *Cluster) Validate() error {
	if !c.EnableLeaderElection {
		return nil
	}

	switch c.LeaderElector {

	case cluster.EtcdLeaderElector:
		if c.EtcdElectionPrefix == "" {
			return fmt.Errorf("etcd-election-prefix must be set")
		}
		return nil

	default:
		return fmt.Errorf("unknown leader elector, valid values are: "+
			"\"%v\"", cluster.EtcdLeaderElector)
	}
}

// Compile-time constraint to ensure Workers implements the Validator interface.
var _ Validator = (*Cluster)(nil)
