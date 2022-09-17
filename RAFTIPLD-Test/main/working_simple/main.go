package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/raft"
	"github.com/libp2p/go-libp2p"
	libp2praft "github.com/libp2p/go-libp2p-raft"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"

	"github.com/functionland/go-fula/event"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/storage/memstore"
	"github.com/multiformats/go-multicodec"
)

var lp = cidlink.LinkPrototype{
	Prefix: cid.Prefix{
		Version:  1,
		Codec:    uint64(multicodec.DagCbor),
		MhType:   uint64(multicodec.Sha2_256),
		MhLength: -1,
	},
}

func main() {
	fmt.Println("start")

	// This example shows how to use go-libp2p-raft to create a cluster
	// which agrees on a State. In order to do it, it defines a state,
	// creates three Raft nodes and launches them. We call a function which
	// lets the cluster leader repeteadly update the state. At the
	// end of the execution we verify that all members have agreed on the
	// same state.
	//
	// Some error handling has been excluded for simplicity

	// error handling ommitted
	newPeer := func(listenPort int) host.Host {
		h, _ := libp2p.New(
			libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		)
		return h
	}

	// Create peers and make sure they know about each others.
	peer1 := newPeer(9997)
	peer2 := newPeer(9998)
	peer3 := newPeer(9999)
	defer peer1.Close()
	defer peer2.Close()
	defer peer3.Close()
	peer1.Peerstore().AddAddrs(peer2.ID(), peer2.Addrs(), peerstore.PermanentAddrTTL)
	peer1.Peerstore().AddAddrs(peer3.ID(), peer3.Addrs(), peerstore.PermanentAddrTTL)
	peer2.Peerstore().AddAddrs(peer1.ID(), peer1.Addrs(), peerstore.PermanentAddrTTL)
	peer2.Peerstore().AddAddrs(peer3.ID(), peer3.Addrs(), peerstore.PermanentAddrTTL)
	peer3.Peerstore().AddAddrs(peer1.ID(), peer1.Addrs(), peerstore.PermanentAddrTTL)
	peer3.Peerstore().AddAddrs(peer2.ID(), peer2.Addrs(), peerstore.PermanentAddrTTL)

	// Create the consensus instances and initialize them with a state.
	// Note that state is just used for local initialization, and that,
	// only states submitted via CommitState() alters the state of the
	// cluster.
	first := &event.Event{
		Version:   event.Version0,
		Peer:      "0",
		Signature: []byte("sig1"),
	}
	type raftState struct {
		Value int
	}
	consensus1 := libp2praft.NewConsensus(&raftState{3})
	consensus2 := libp2praft.NewConsensus(&raftState{3})
	consensus3 := libp2praft.NewConsensus(&raftState{3})

	// Create LibP2P transports Raft
	transport1, err := libp2praft.NewLibp2pTransport(peer1, time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}
	transport2, err := libp2praft.NewLibp2pTransport(peer2, time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}
	transport3, err := libp2praft.NewLibp2pTransport(peer3, time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer transport1.Close()
	defer transport2.Close()
	defer transport3.Close()

	// Create Raft servers configuration for bootstrapping the cluster
	// Note that both IDs and Address are set to the Peer ID.
	servers := make([]raft.Server, 0)
	for _, h := range []host.Host{peer1, peer2, peer3} {
		servers = append(servers, raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(h.ID().Pretty()),
			Address:  raft.ServerAddress(h.ID().Pretty()),
		})
	}
	serversCfg := raft.Configuration{Servers: servers}

	// Create Raft Configs. The Local ID is the PeerOID
	config1 := raft.DefaultConfig()
	config1.LogOutput = io.Discard
	config1.Logger = nil
	config1.LocalID = raft.ServerID(peer1.ID().Pretty())

	config2 := raft.DefaultConfig()
	config2.LogOutput = io.Discard
	config2.Logger = nil
	config2.LocalID = raft.ServerID(peer2.ID().Pretty())

	config3 := raft.DefaultConfig()
	config3.LogOutput = io.Discard
	config3.Logger = nil
	config3.LocalID = raft.ServerID(peer3.ID().Pretty())

	// Create snapshotStores. Use FileSnapshotStore in production.
	snapshots1 := raft.NewInmemSnapshotStore()
	snapshots2 := raft.NewInmemSnapshotStore()
	snapshots3 := raft.NewInmemSnapshotStore()

	// Create the InmemStores for use as log store and stable store.
	logStore1 := raft.NewInmemStore()
	logStore2 := raft.NewInmemStore()
	logStore3 := raft.NewInmemStore()

	// Bootsrap the stores with the serverConfigs
	raft.BootstrapCluster(config1, logStore1, logStore1, snapshots1, transport1, serversCfg.Clone())
	raft.BootstrapCluster(config2, logStore2, logStore2, snapshots2, transport2, serversCfg.Clone())
	raft.BootstrapCluster(config3, logStore3, logStore3, snapshots3, transport3, serversCfg.Clone())

	// Create Raft objects. Our consensus provides an implementation of
	// Raft.FSM
	raft1, err := raft.NewRaft(config1, consensus1.FSM(), logStore1, logStore1, snapshots1, transport1)
	if err != nil {
		log.Fatal(err)
	}
	raft2, err := raft.NewRaft(config2, consensus2.FSM(), logStore2, logStore2, snapshots2, transport2)
	if err != nil {
		log.Fatal(err)
	}
	raft3, err := raft.NewRaft(config3, consensus3.FSM(), logStore3, logStore3, snapshots3, transport3)
	if err != nil {
		log.Fatal(err)
	}

	// Create the actors using the Raft nodes
	actor1 := libp2praft.NewActor(raft1)
	actor2 := libp2praft.NewActor(raft2)
	actor3 := libp2praft.NewActor(raft3)

	// Set the actors so that we can CommitState() and GetCurrentState()
	consensus1.SetActor(actor1)
	consensus2.SetActor(actor2)
	consensus3.SetActor(actor3)

	ls := cidlink.DefaultLinkSystem()
	store := &memstore.Store{}
	ls.SetReadStorage(store)
	ls.SetWriteStorage(store)

	ctx := context.Background()
	// This function updates the cluster state commiting 1000 updates.
	last := first
	updateState := func(c *libp2praft.Consensus) {
		nUpdates := 0
		for {
			if nUpdates >= 3 {
				break
			}

			newState := &event.Event{
				Version:   event.Version0,
				Peer:      "2",
				Signature: []byte("sig"),
			}

			n := bindnode.Wrap(last, event.Prototypes.Event.Type())
			if l, err := ls.Store(ipld.LinkContext{Ctx: ctx}, lp, n); err != nil {
				fmt.Println(err)
			} else {
				newState.Previous = l
				fmt.Printf("Link to last event: %s\n", l.String())
				last = newState
			}

			//newState := &raftState{nUpdates * 2}

			// CommitState() blocks until the state has been
			// agreed upon by everyone
			/*cidstring, err := fmt.Printf("link: %s", last)
			if err != nil {
				fmt.Println(err)
				continue
			}*/

			a := &raftState{nUpdates * 2}

			agreedState, err := c.CommitState(a)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if agreedState == nil {
				fmt.Println("agreedState is nil: commited on a non-leader?")
				continue
			} else {
				fmt.Println("agreedState is commited")
			}

			nUpdates++

			agreedRaftState := agreedState.(*raftState)
			if nUpdates%200 == 0 {
				stringagreedRaftState, err := json.Marshal(agreedRaftState)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("Performed %d updates. Current state value: %s\n", nUpdates, stringagreedRaftState)
			}
		}
	}

	// Provide some time for leader election
	time.Sleep(5 * time.Second)

	// Run the 1000 updates on the leader
	// Barrier() will wait until updates have been applied
	if actor1.IsLeader() {
		updateState(consensus1)
	} else if actor2.IsLeader() {
		updateState(consensus2)
	} else if actor3.IsLeader() {
		updateState(consensus3)
	}

	// Wait for updates to arrive.
	time.Sleep(5 * time.Second)

	// Shutdown raft and wait for it to complete
	// (ignoring errors)
	raft1.Shutdown().Error()
	raft2.Shutdown().Error()
	raft3.Shutdown().Error()

	// Final states
	finalState1, err := consensus1.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}
	finalState2, err := consensus2.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}
	finalState3, err := consensus3.GetCurrentState()
	if err != nil {
		fmt.Println(err)
		return
	}
	finalRaftState1 := finalState1.(*raftState)
	finalRaftState2 := finalState2.(*raftState)
	finalRaftState3 := finalState3.(*raftState)

	stringOutput1, err1 := json.Marshal(finalRaftState1)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	stringOutput2, err2 := json.Marshal(finalRaftState2)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	stringOutput3, err3 := json.Marshal(finalRaftState3)
	if err3 != nil {
		fmt.Println(err3)
		return
	}

	fmt.Printf("Raft1 final state: %s\n", string(stringOutput1))
	fmt.Printf("Raft2 final state: %s\n", string(stringOutput2))
	fmt.Printf("Raft3 final state: %s\n", string(stringOutput3))
	// Output:
	// Performed 200 updates. Current state value: 398
	// Performed 400 updates. Current state value: 798
	// Performed 600 updates. Current state value: 1198
	// Performed 800 updates. Current state value: 1598
	// Performed 1000 updates. Current state value: 1998
	// Raft1 final state: 1998
	// Raft2 final state: 1998
	// Raft3 final state: 1998
}
