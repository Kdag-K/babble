package hashgraph

import (
	"testing"

	"github.com/Kdag-K/kdag/src/crypto/keys"
	"github.com/Kdag-K/kdag/src/peers"
)

func createTestBlock() *Block {
	block := NewBlock(0, 1,
		[]byte("framehash"),
		[]*peers.Peer{},
		[][]byte{
			[]byte("abc"),
			[]byte("def"),
			[]byte("ghi"),
		},
		[]InternalTransaction{
			NewInternalTransaction(PEER_ADD, *peers.NewPeer("peer1", "paris", "peer1")),
			NewInternalTransaction(PEER_REMOVE, *peers.NewPeer("peer2", "london", "peer2")),
		},
		0,
	)

	receipts := []InternalTransactionReceipt{}
	for _, itx := range block.InternalTransactions() {
		receipts = append(receipts, itx.AsAccepted())
	}
	block.Body.InternalTransactionReceipts = receipts

	return block
}

func TestSignBlock(t *testing.T) {
	privateKey, _ := keys.GenerateECDSAKey()

	block := createTestBlock()

	sig, err := block.Sign(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	res, err := block.Verify(sig)
	if err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}
	if !res {
		t.Fatal("Verify returned false")
	}
}

func TestAppendSignature(t *testing.T) {
	privateKey, _ := keys.GenerateECDSAKey()

	block := createTestBlock()

	sig, err := block.Sign(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	err = block.SetSignature(sig)
	if err != nil {
		t.Fatal(err)
	}

	blockSignature, err := block.GetSignature(keys.PublicKeyHex(&privateKey.PublicKey))
	if err != nil {
		t.Fatal(err)
	}

	res, err := block.Verify(blockSignature)
	if err != nil {
		t.Fatalf("Error verifying signature: %v", err)
	}
	if !res {
		t.Fatal("Verify returned false")
	}
}

func TestNewBlockFromFrame(t *testing.T) {

	frameTimestamp := int64(123456789)

	transactions := [][]byte{
		[]byte("transaction1"),
		[]byte("transaction2"),
		[]byte("transaction3"),
		[]byte("transaction4"),
		[]byte("transaction5"),
		[]byte("transaction6"),
		[]byte("transaction7"),
		[]byte("transaction8"),
		[]byte("transaction9"),
	}

	internalTransactions := []InternalTransaction{
		NewInternalTransaction(PEER_ADD, *peers.NewPeer("peer1000.pub", "peer1000.addr", "peer1000")),
		NewInternalTransaction(PEER_ADD, *peers.NewPeer("peer1001.pub", "peer1001.addr", "peer1001")),
		NewInternalTransaction(PEER_ADD, *peers.NewPeer("peer1002.pub", "peer1002.addr", "peer1002")),
	}

	frame := &Frame{
		Round: 56,
		Peers: []*peers.Peer{
			peers.NewPeer("peer1.pub", "peer1.addr", "peer1"),
			peers.NewPeer("peer2.pub", "peer2.addr", "peer2"),
			peers.NewPeer("peer3.pub", "peer3.addr", "peer3"),
		},
		Roots: nil,
		Events: []*FrameEvent{
			{
				Core: &Event{
					Body: EventBody{
						Transactions:         transactions[0:3],
						InternalTransactions: internalTransactions[:1],
					},
				},
			},
			{
				Core: &Event{
					Body: EventBody{
						Transactions:         transactions[3:6],
						InternalTransactions: internalTransactions[1:2],
					},
				},
			},
			{
				Core: &Event{
					Body: EventBody{
						Transactions:         transactions[6:],
						InternalTransactions: internalTransactions[2:],
					},
				},
			},
		},
		Timestamp: frameTimestamp,
	}

	block, err := NewBlockFromFrame(10, frame)
	if err != nil {
		t.Fatal(err)
	}

}
