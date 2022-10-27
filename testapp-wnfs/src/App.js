import logo from './logo.svg';
import './App.css';
import {useEffect} from 'react';
import { MemoryBlockStore } from 'ipfs-car/blockstore/memory' 
import init, { PrivateDirectory, PrivateForest, Namefilter } from "fx-wnfs";
import * as nacl from 'tweetnacl';

const cid = Uint8Array.from([
  1, 112, 18, 32, 195, 196, 115, 62, 200, 175, 253, 6, 207, 158, 159, 245, 15,
  252, 107, 205, 46, 200, 90, 97, 112, 0, 75, 183, 9, 102, 156, 49, 222, 148,
  57, 26,
]);
const time = new Date();
const rng= nacl;
function App() {
    //

    const store = new MemoryBlockStore();
    store.putBlock = store.put;
    useEffect(() => {
       init().then(() => {
    
       
         const initialHamt = new PrivateForest();
         const dir = new PrivateDirectory(new Namefilter(), new Date(), rng);
         const fetchData = async () => {
        //START
        console.log(dir);
            var { rootDir, hamt } = await dir.mkdir(
              ["pictures", "cats"],
              true,
              new Date(),
              initialHamt,
              store,
              rng
            );
            var { rootDir, hamt } = await rootDir.write(
              ["pictures", "cats", "tabby.png"],
              cid,
              time,
              store
            );
            var { result } = await rootDir.ls(["pictures"], store);
            
            console.log(result);
        //END

          }
         fetchData();
          /*dir.mkdir(["pictures", "cats"], new Date(), store).then((res)=>{
             console.log(res);
             // Create a sample CIDv1.
            const cid = Uint8Array.from([
              1, 112, 18, 32, 195, 196, 115, 62, 200, 175, 253, 6, 207, 158, 159, 245, 15,
              252, 107, 205, 46, 200, 90, 97, 112, 0, 75, 183, 9, 102, 156, 49, 222, 148,
              57, 26,
            ]);

            // Add a file to /pictures/cats.
            res.rootDir.write(
              ["pictures", "cats", "tabby.png"],
              cid,
              new Date(),
              store
            ).then((res2)=>{
                res2.rootDir.ls(["pictures"], store).then((res3)=>{
                    console.log(res3)
                });
            });
         });
         */
       });
 }, []);
    
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
