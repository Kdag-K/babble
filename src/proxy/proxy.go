package proxy

import (
	"github.com/Kdag-K/kdag/src/hashgraph"
	"github.com/Kdag-K/kdag/src/node/state"
)

// AppGateway defines the interface which is used by Kdag to communicate with
// the App
type AppGateway interface {
	SubmitCh() chan []byte
	CommitBlock(block hashgraph.Block) (CommitResponse, error)
	GetSnapshot(blockIndex int) ([]byte, error)
	Restore(snapshot []byte) error
	OnStateChanged(state.State) error
}
