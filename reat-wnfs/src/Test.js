import React, { useEffect, useState } from 'react'
import wasmInit, {PublicDirectory, PrivateDirectory, Namefilter, PrivateForest} from 'fx-wnfs'
// import wasmInit, {PublicDirectory} from 'wnfs'
import { MemoryBlockStore } from './store'




class Rng {
  /** Returns random bytes of specified length */
  randomBytes(count){
    const array = new Uint8Array(count);

    window.crypto.getRandomValues(array);
    return array;
  }
}

const Test = () => {

  const [init, setInit] = useState(false)

  const testPublic = async () => {
    if (!init) {
      await wasmInit()
      setInit(true)
    }
    

    const time = new Date()
    const dir = new PublicDirectory(time);
    const store = new MemoryBlockStore();

    // Create a /pictures/cats directory.
    var { rootDir } = await dir.mkdir(["pictures", "cats"], time, store);

    // Get a sample CIDv1.
    const cid = Uint8Array.from([
      1, 112, 18, 32, 195, 196, 115, 62, 200, 175, 253, 6, 207, 158, 159, 245, 15,
      252, 107, 205, 46, 200, 90, 97, 112, 0, 75, 183, 9, 102, 156, 49, 222, 148,
      57, 26,
    ]);

    // Add a file to /pictures/cats.
    var { rootDir } = await rootDir.write(
      ["pictures", "cats", "tabby.png"],
      cid,
      time,
      store
    );

    // Create and add a file to /pictures/dogs directory.
    var { rootDir } = await rootDir.write(
      ["pictures", "dogs", "billie.jpeg"],
      cid,
      time,
      store
    );

    // Delete /pictures/cats directory.
    var { rootDir } = await rootDir.rm(["pictures", "cats"], store);

    // // List all files in /pictures directory.
    var { result } = await rootDir.ls(["pictures"], store);

    console.log("Files in /pictures directory:", result);

    var img = await rootDir.read(["pictures", "dogs", "billie.jpeg"])

    console.log({ img })
  }

  const testPrivate = async () => {
    if (!init) {
      await wasmInit()
      setInit(true)
    }

    const time = new Date()
    const rng = new Rng()
    const nf = new Namefilter()
    const dir = new PrivateDirectory(nf, time, rng);
    const store = new MemoryBlockStore();
    const pf = new PrivateForest()

    const searchLatest = true

    console.log({dir})

    // Create a /pictures/cats directory.
    var { rootDir } = await dir.mkdir(["pictures", "cats"], searchLatest, time, pf, store, rng);

    console.log({rootDir})

    // Get a sample CIDv1.
    const cid = Uint8Array.from([
      1, 112, 18, 32, 195, 196, 115, 62, 200, 175, 253, 6, 207, 158, 159, 245, 15,
      252, 107, 205, 46, 200, 90, 97, 112, 0, 75, 183, 9, 102, 156, 49, 222, 148,
      57, 26,
    ]);

    // Add a file to /pictures/cats.
    var { rootDir } = await rootDir.write(
      ["pictures", "cats", "tabby.png"],
      searchLatest,
      cid,
      time,
      pf,
      store,
      rng
    );

    // Create and add a file to /pictures/dogs directory.
    var { rootDir } = await rootDir.write(
      ["pictures", "dogs", "billie.jpeg"],
      searchLatest,
      cid,
      time,
      pf,
      store,
      rng
    );

    // Delete /pictures/cats directory.
    var { rootDir } = await rootDir.rm(["pictures", "cats"], searchLatest, pf, store, rng);

    // // List all files in /pictures directory.
    var { result } = await rootDir.ls(["pictures"], searchLatest, pf, store);

    console.log("Files in /pictures directory:", result);

    var img = await rootDir.read(["pictures", "dogs", "billie.jpeg"], searchLatest, pf, store)

    console.log({ img })

    console.log({dir})
  }

  useEffect(() => {
    // setPublicDir(pd)
    // if (!init)
      testPublic()
      // testPrivate()
      // console.log({wasm})
  }, [])
  return (
    <div>Test</div>
  )
}

export default Test