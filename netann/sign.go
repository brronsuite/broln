package netann

import (
	"fmt"

	"github.com/brronsuite/broln/input"
	"github.com/brronsuite/broln/keychain"
	"github.com/brronsuite/broln/lnwallet"
	"github.com/brronsuite/broln/lnwire"
)

// SignAnnouncement signs any type of gossip message that is announced on the
// network.
func SignAnnouncement(signer lnwallet.MessageSigner, keyLoc keychain.KeyLocator,
	msg lnwire.Message) (input.Signature, error) {

	var (
		data []byte
		err  error
	)

	switch m := msg.(type) {
	case *lnwire.ChannelAnnouncement:
		data, err = m.DataToSign()
	case *lnwire.ChannelUpdate:
		data, err = m.DataToSign()
	case *lnwire.NodeAnnouncement:
		data, err = m.DataToSign()
	default:
		return nil, fmt.Errorf("can't sign %T message", m)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to get data to sign: %v", err)
	}

	return signer.SignMessage(keyLoc, data, true)
}
