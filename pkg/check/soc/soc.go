package soc

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ethersphere/bee/pkg/cac"
	"github.com/ethersphere/bee/pkg/crypto"
	"github.com/ethersphere/bee/pkg/soc"
	"github.com/ethersphere/beekeeper/pkg/bee"
	"github.com/ethersphere/beekeeper/pkg/check"
	"github.com/prometheus/client_golang/prometheus/push"
)

// compile check whether Check implements interface
var _ check.Check = (*Check2)(nil)

// TODO: rename to Check
// Check instance
type Check2 struct{}

// NewCheck returns new check
func NewCheck() check.Check {
	return &Check2{}
}

// Options represents check options
type Options struct {
	NodeGroup string
	Seed      int64
}

func (c *Check2) Run(ctx context.Context, cluster *bee.Cluster, o interface{}) (err error) {
	return
}

// Check sends a SOC chunk and retrieves with the address.
func Check(c *bee.Cluster, o Options, pusher *push.Pusher, pushMetrics bool) error {

	payload := []byte("Hello Swarm :)")

	ng := c.NodeGroup(o.NodeGroup)
	sortedNodes := ng.NodesSorted()

	privKey, err := crypto.GenerateSecp256k1Key()
	if err != nil {
		return err
	}
	signer := crypto.NewDefaultSigner(privKey)

	ch, err := cac.New(payload)
	if err != nil {
		return err
	}

	idBytes, err := randomID()
	if err != nil {
		return err
	}
	sch, err := soc.New(idBytes, ch).Sign(signer)
	if err != nil {
		return err
	}

	chunkData := sch.Data()
	signatureBytes := chunkData[soc.IdSize : soc.IdSize+soc.SignatureSize]

	publicKey, err := signer.PublicKey()
	if err != nil {
		return err
	}

	ownerBytes, err := crypto.NewEthereumAddress(*publicKey)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	nodeName := sortedNodes[0]
	node := ng.NodeClient(nodeName)

	owner := hex.EncodeToString(ownerBytes)
	id := hex.EncodeToString(idBytes)
	sig := hex.EncodeToString(signatureBytes)

	fmt.Printf("soc: submitting soc chunk %s to node %s\n", sch.Address().String(), nodeName)
	fmt.Printf("soc: owner %s\n", owner)
	fmt.Printf("soc: id %s\n", id)
	fmt.Printf("soc: sig %s\n", sig)

	ref, err := node.UploadSOC(ctx, owner, id, sig, ch.Data())
	if err != nil {
		return err
	}

	retrieved, err := node.DownloadChunk(ctx, ref, "")
	if err != nil {
		return err
	}

	if !bytes.Equal(retrieved, chunkData) {
		return errors.New("soc: retrieved chunk data does NOT match soc chunk")
	}

	return nil
}

func randomID() ([]byte, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}