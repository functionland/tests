import logo from './logo.svg';
import './App.css';
import {useEffect} from 'react';
import { MemoryBlockStore } from 'ipfs-car/blockstore/memory';
import { MemoryBlockstore as Blockstore } from 'blockstore-core/memory';
import { MemoryDatastore } from 'datastore-core/memory';
import { BlockstoreDatastoreAdapter } from 'blockstore-datastore-adapter';
import init, { PrivateDirectory, PrivateForest, Namefilter, PublicDirectory } from "fx-wnfs";
import { CID } from 'multiformats/cid';
import * as json from 'multiformats/codecs/json';
import { sha256 } from 'multiformats/hashes/sha2';

const CID2 = require('cids');
const multihashing = require('multihashing-async')
var sjcl = require('randombytes');

const cidi = Uint8Array.from([
  1, 112, 18, 32, 195, 196, 115, 62, 200, 175, 253, 6, 207, 158, 159, 245, 15,
  252, 107, 205, 46, 200, 90, 97, 112, 0, 75, 183, 9, 102, 156, 49, 222, 148,
  57, 26,
]);

const time = new Date();
const rng = {randomBytes: sjcl};
function App() {
    //

    const store = new MemoryBlockStore();
    const store2 = new Blockstore();
    const store3 = new BlockstoreDatastoreAdapter(new MemoryDatastore());
    const store4 = {data:{}};

    store.putBlock = (cid, data) => {
      console.log({cid: cid, data: data});
      return store.put(cid, data);
    }
    store.getBlock = store.get;

    store2.putBlock = (cid, data) => {
      console.log({cid: cid, data: data});
      return store2.put(cid, data);
    }
    store2.getBlock = (cid) => {
      console.log({cid: cid,});
      return store2.get(cid);
    }

    store3.putBlock = (cid, data=null, options=null) => {
      console.log("put", {cid: cid, data: data, options:options});
      return store3.put(cid, data);
    }
    store3.getBlock = (cid, options=null) => {
      console.log("get",{cid: cid, options:options});
      return store3.get(cid);
    }

    store4.putBlock = async (data, codec, options=null) => {
      let hash = await sha256.digest(data)
      let cid = CID.create(1, json.code, hash)

      const hash2 = await multihashing(data, 'sha2-256')
      let cid2 = new CID2(1, 'dag-cbor', hash2);
      store4.data[cid2] = data;
      let validateCid = CID2.validateCID(cid2);
      console.log("dir",{original: {codec: codec, data: data, options:options}, cid: cid2, cidstring: cid2.toString(), validateCid:validateCid});
      return cid2.toString();
    }
    store4.getBlock = (cid, options=null) => {
      console.log("get",{cid: cid, options:options});
      return store4.data[cid];
    }

    
    
    useEffect(() => {
      const fetchData = async () => {
          //START

          //PUBLIC TEST
          console.log("START OF PUBLIC TEST");
          const dirp = new PublicDirectory(new Date());
          console.log("dirp", dirp);
          var { rootDir } = await dirp.mkdir(["pictures", "cats"], new Date(), store4);
          console.log("rootDir", rootDir);
          var { rootDir } = await rootDir.write(
            ["pictures", "cats", "tabby.png"],
            cidi,
            new Date(),
            store4
          );
          var { result } = await rootDir.ls(["pictures", "cats"], store4);
          console.log("Files in /pictures directory:", result);
          //END
        //END OF PUBLIC TEST

        //PRIVATE TEST
        console.log("START OF PRIVATE TEST");
        var hamt = await new PrivateForest();
        console.log("initialHamt", hamt);
        const dir = await new PrivateDirectory(await new Namefilter(), new Date(), rng);
        var { rootDir, hamt } = await dir.mkdir(
          ["pictures", "cats"],
          true,
          new Date(),
          hamt,
          store4,
          rng
        );
        //END OF PRIVATE TEST
      }
       init().then(() => {
        
        fetchData();
       });
 }, [store4]);
    
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
