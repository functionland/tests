
// import {fs} from "memfs"
import * as Ipfs from "ipfs-core"
// import tempDir from "ipfs-utils/src/temp-dir.js"
import { createRepo } from "ipfs-repo"
import { MemoryDatastore } from "datastore-core/memory"
import { MemoryBlockstore } from "blockstore-core/memory"

// Codecs
import * as dagPB from "@ipld/dag-pb"
import * as dagCBOR from "@ipld/dag-cbor"
import * as raw from "multiformats/codecs/raw"

const CODECS = {
  [dagPB.code]: dagPB,
  [dagPB.name]: dagPB,
  [dagCBOR.code]: dagCBOR,
  [dagCBOR.name]: dagCBOR,
  [raw.code]: raw,
  [raw.name]: raw,
}

export async function createInMemoryIPFS() {
//   const dir = fs.mkdirSync()
//   fs.mkdirSync("/ipfs")

  const memoryDs = new MemoryDatastore()
  const memoryBs = new MemoryBlockstore()

  const ipfs = await Ipfs.create({
    offline: true,
    silent: true,
    preload: {
      enabled: false,
    },
    config: {
    //   Addresses: {
    //     Swarm: ["/ip4/0.0.0.0/tcp/4002"],
    //     API: "/ip4/127.0.0.1/tcp/5002",
    //     Gateway: "/ip4/127.0.0.1/tcp/9090"
    //   },
      Discovery: {
        MDNS: {
          Enabled: false
        },
        webRTCStar: {
          Enabled: false
        }
      },
      Pubsub: {
        Enabled: false
      }
    },
    libp2p: {
      peerDiscovery: [],
      connectionManager: { autoDial: false }
    },
    repo: createRepo(
      "/ipfs",
      codeOrName => Promise.resolve(CODECS[codeOrName]), {
      root: memoryDs,
      blocks: memoryBs,
      keys: memoryDs,
      datastore: memoryDs,
      pins: memoryDs
    }, {
      repoLock: {
        lock: async () => ({ close: async () => { return } }),
        locked: async () => false
      },
      autoMigrate: false,
    }
    )
  })

  return ipfs
}