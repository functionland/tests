import logo from './logo.svg';
import './App.css';
import { useEffect, useState } from 'react';
import * as wn from 'webnative'
import { createInMemoryIPFS } from './helpers/ipfs'
import useIpfsFactory from './hooks/use-ipfs-factory';
import { click } from '@testing-library/user-event/dist/click';
import { init } from './helpers/wnfs'

const App = () => {
  // const { ipfs, ipfsInitError } = useIpfsFactory()

  const [cid, setCid] = useState(false)
  const [state, setState] = useState(null)
  const [path, setPath] = useState("public/Docs")
  const [contentText, setContentText] = useState("")

  const [ipfs, setIpfs] = useState(false)
  useEffect(() => {
    const createIpfs = async () => {
      // const i = await createInMemoryIPFS()
      // setIpfs(i)
      const state = await init(wn)

      if (state)
        setState(state)
    }
    createIpfs()
  }, [])

  // useEffect(() => {
  //   if (state)
  //     testFS()
  // }, [state])




  const ls = async () => {
    if (!state)
      return

    let efs = state.fs

    let p = wn.path.directory(...path.split("/"))
    let l = await efs.ls(p)
    console.log(`Ls result on ${path}`, l)
  }

  const writeFile = async () => {
    if (!state)
      return

    let efs = state.fs

    let p = wn.path.file(...path.split("/"))

    console.log("Writing file at", path)
    await efs.add(p, contentText)

    console.log("Added file", p, contentText)
    return
  }

  const readFile = async () => {
    if (!state)
      return

    let p = wn.path.file(...path.split("/"))

    console.log("Reading file at", path)
    let rfc = await state.fs.cat(p)

    console.log("File read", rfc)
  }


  const publish = async () => {
    console.log("Publishing rootDir")
    let cid = await state.fs.publish()

    console.log("rootDir published, cid:", cid.toV0().toString())
  }

const testIPFS = async () => {
  console.log("There")
  let root = await ipfs.dag.get(cid)
  console.log({ root })
}

return (
  <div className="App">
    <header className="App-header">
      <img src={logo} className="App-logo" alt="logo" />
      <p>
        Edit <code>src/App.js</code> and save to reload.
      </p>
      <input type="text" value={path} placeholder="Path" onChange={evt => setPath(evt.target.value)}></input>
      <input type="text" placeholder="File Content" onChange={evt => setContentText(evt.target.value)}></input>
      <button onClick={ls}>LS</button>
      <button onClick={readFile}>Read File</button>
      <button onClick={writeFile}>Write File</button>
      <button onClick={publish}>Publish</button>
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
