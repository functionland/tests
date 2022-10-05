package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	wnfs "github.com/functionland/wnfs-go"
	base "github.com/functionland/wnfs-go/base"
	private "github.com/functionland/wnfs-go/private"
	ratchet "github.com/functionland/wnfs-go/private/ratchet"
	cmp "github.com/google/go-cmp/cmp"

	// we do not need the whole ipfslite below
	ipfslite "github.com/hsanjuan/ipfs-lite"
	bitswap "github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-bitswap/network"
	blockservice "github.com/ipfs/go-blockservice"
	datastore "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	provider "github.com/ipfs/go-ipfs-provider"
	ipld "github.com/ipfs/go-ipld-format"
	merkledag "github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var testRootKey Key = [32]byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2,
}

type (
	Node         = base.Node
	HistoryEntry = base.HistoryEntry
	PrivateName  = private.Name
	Key          = private.Key
)

// Config wraps configuration options for the Peer.
type Config struct {
	// The DAGService will not announce or retrieve blocks from the network
	Offline bool
	// ReprovideInterval sets how often to reprovide records to the DHT
	ReprovideInterval time.Duration
}

// Peer is an IPFS-Lite peer. It provides a DAG service that can fetch and put
// blocks from/to the IPFS network.
type Peer struct {
	ctx context.Context

	cfg *Config

	host  host.Host
	dht   routing.Routing
	store datastore.Batching

	ipld.DAGService // become a DAG service
	exch            exchange.Interface
	bstore          blockstore.Blockstore
	bserv           blockservice.BlockService
	reprovider      provider.System
}

var secret = "2cc2c79ea52c9cc85dfd3061961dd8c4230cce0b09f182a0822c1536bf1d5f21"

// New creates an IPFS-Lite Peer. It uses the given datastore, libp2p Host and
// Routing (usuall the DHT). The Host and the Routing may be nil if
// config.Offline is set to true, as they are not used in that case. Peer
// implements the ipld.DAGService interface.
func New(
	ctx context.Context,
	store datastore.Batching,
	host host.Host,
	dht routing.Routing,
	cfg *Config,
) (*Peer, error) {

	if cfg == nil {
		cfg = &Config{}
	}

	p := &Peer{
		ctx:   ctx,
		cfg:   cfg,
		host:  host,
		dht:   dht,
		store: store,
	}

	return p, nil
}

func setupPeers(ctx context.Context) (p1 *Peer, closer func()) {

	ds1 := ipfslite.NewInMemoryDatastore()

	priv1, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		println(err)
	}

	psk, err := hex.DecodeString(secret)
	if err != nil {
		println(err)
	}

	listen, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")
	h1, dht1, err := ipfslite.SetupLibp2p(
		ctx,
		priv1,
		psk,
		[]multiaddr.Multiaddr{listen},
		nil,
	)
	if err != nil {
		println(err)
	}

	closer = func() {
		for _, cl := range []io.Closer{dht1, h1} {
			err := cl.Close()
			if err != nil {
				println(err)
			}
		}
	}
	p1, err = New(ctx, ds1, h1, dht1, nil)
	if err != nil {
		closer()
		println(err)
	}

	return
}

func newFileTestStorePrivate(ctx context.Context, p *Peer) (st private.Store) {

	//func (p *Peer) setupBlockstore()
	bs := blockstore.NewBlockstore(p.store)
	bs = blockstore.NewIdStore(bs)
	cachedbs, err1 := blockstore.CachedBlockstore(p.ctx, bs, blockstore.DefaultCacheOpts())
	if err1 != nil {
		println(err1.Error())
	}
	p.bstore = cachedbs

	p.DAGService = merkledag.NewDAGService(p.bserv)

	//func (p *Peer) setupBlockService()
	if p.cfg.Offline {
		p.bserv = blockservice.New(p.bstore, offline.Exchange(p.bstore))
	} else {
		bswapnet := network.NewFromIpfsHost(p.host, p.dht)
		bswap := bitswap.New(p.ctx, bswapnet, p.bstore, bitswap.ProvideEnabled(true))
		p.bserv = blockservice.New(p.bstore, bswap)
		p.exch = bswap
	}

	/*bserv, cleanup, err := mockblocks.NewOfflineFileBlockservice("test", bswap)
	if err != nil {
		println(err.Error())
	}*/
	rs := ratchet.NewMemStore(ctx)
	store, err := private.NewStore(ctx, p.bserv, rs)
	if err != nil {
		println(err)
	}

	return store
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rs := ratchet.NewMemStore(ctx)
	p1, cleanup := setupPeers(ctx)
	store := newFileTestStorePrivate(ctx, p1)
	//cleanup()
	fsys, err := wnfs.NewEmptyFS(ctx, store.Blockservice(), rs, testRootKey)
	if err != nil {
		println(err)
		cleanup()
	}
	fmt.Printf("wnfs root CID: %s\n", fsys.Cid())

	ls_i, err_i := fsys.Ls("private")
	if err_i != nil {
		println(err_i.Error())
	}
	fmt.Printf("\ninitial Folder structure: %s\n", ls_i)

	pathStr := "private/foo/hello.txt"
	fileContents := []byte("hello!")
	f := base.NewMemfileBytes("hello.txt", fileContents)

	err1 := fsys.Write(pathStr, f)
	if err1 != nil {
		fmt.Printf("Erro happened %s\n", err1.Error())
	}
	fsys.Commit()
	fmt.Printf("wnfs new root CID: %s\n", fsys.Cid())

	gotFileContents, err := fsys.Cat(pathStr)
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("file content: %s\n", string(gotFileContents))

	if diff := cmp.Diff(fileContents, gotFileContents); diff != "" {
		fmt.Printf("\nresult mismatch. (-want +got):\n%s", diff)
	}

	ls, err := fsys.Ls("private/foo")
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("\nFolder structure: %s", ls)

	println("\nEnd")

}
