package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/functionland/go-fula/pool"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func main() {

	const poolName = "my-pool"
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Elevate log level to show internal communications.
	if err := logging.SetLogLevel("*", "error"); err != nil {
		panic(err)
	}

	// Use a deterministic random generator to generate deterministic
	// output for the example.
	rng := rand.New(rand.NewSource(42))

	// Instantiate the first node in the pool
	pid1, _, err := crypto.GenerateECDSAKeyPair(rng)
	if err != nil {
		panic(err)
	}
	ma := "/ip4/192.168.0.2/tcp/64658"
	h1, err := libp2p.New(libp2p.Identity(pid1), libp2p.ListenAddrStrings(ma))
	if err != nil {
		panic(err)
	}
	n1, err := pool.New(pool.WithPoolName(poolName), pool.WithHost(h1))
	if err != nil {
		panic(err)
	}
	if err := n1.Start(ctx); err != nil {
		panic(err)
	}
	defer n1.Shutdown(ctx)
	fmt.Printf("Instantiated node in pool %s with ID: %s\nHost multiaddr: %s\n", poolName, h1.ID().String(), ma)

	// Generate a sample DAG and store it on node 1 (n1) in the pool
	leaf := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
		na.AssembleEntry("this").AssignBool(true)
	})
	leafLink, err := n1.Store(ctx, leaf)
	if err != nil {
		panic(err)
	}
	root := fluent.MustBuildMap(basicnode.Prototype.Map, 2, func(na fluent.MapAssembler) {
		na.AssembleEntry("that").AssignInt(42)
		na.AssembleEntry("leafLink").AssignLink(leafLink)
	})
	rootLink, err := n1.Store(ctx, root)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s stored IPLD data with links:\n    root: %s\n    leaf:%s\n", h1.ID(), rootLink, leafLink)

	for {

	}

}
