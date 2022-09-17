import logo from './logo.svg';
import './App.css';
import * as UCAN from "@ipld/dag-ucan";
import { identity } from "multiformats/hashes/identity";
import { encode as dagencode, decode as dagdecode } from '@ipld/dag-cbor';

const lib = require("ucans")

const { validate, build, encode, EdKeypair, RsaKeypair, parse, capability } = lib;
  
function App() {
	/*const jwt = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCIsInVjdiI6IjAuOC4xIn0.eyJhdWQiOiJkaWQ6a2V5Ono2TWtmZkRaQ2tDVFdyZWc4ODY4ZkcxRkdGb2djSmo1WDZQWTkzcFBjV0RuOWJvYiIsImF0dCI6W3sid2l0aCI6InduZnM6Ly9ib3Jpcy5maXNzaW9uLm5hbWUvcHVibGljL3Bob3Rvcy8iLCJjYW4iOiJjcnVkL0RFTEVURSJ9LHsid2l0aCI6InduZnM6Ly9ib3Jpcy5maXNzaW9uLm5hbWUvcHJpdmF0ZS84NE1aN2Fxd0tuN3NOaU1Hc1NiYXhzRWE2RVBuUUxvS1liWEJ5eE5CckNFciIsImNhbiI6InduZnMvQVBQRU5EIn0seyJ3aXRoIjoibWFpbHRvOmJvcmlzQGZpc3Npb24uY29kZXMiLCJjYW4iOiJtc2cvU0VORCJ9XSwiZXhwIjoxNjUwNTAwODQ5LCJpc3MiOiJkaWQ6a2V5Ono2TWtrODliQzNKclZxS2llNzFZRWNjNU0xU01WeHVDZ054NnpMWjhTWUpzeEFMaSIsInByZiI6W119.OqM4_glZJq8GRvg7k8U3OJNgjJ_N8ORM5cOKA0O84lE9Ttzy9YJQe6e4QkOhS0uIkzIvxCdWB0DWsFhTc1rtBA";
	
	const ucan = UCAN.parse(jwt);
	const did = ucan.issuer.did();
	const ucanjwt = UCAN.decode(UCAN.encode(ucan));*/
	/////////
	const alice = EdKeypair.fromSecretKey(
	  "U+bzp2GaFQHso587iSFWPSeCzbSfn/CbNHEz7ilKRZ1UQMmMS7qq4UhTzKn3X9Nj/4xgrwa+UqhMOeo4Ki8JUw=="
	);
	const bob = EdKeypair.fromSecretKey(
	  "G4+QCX1b3a45IzQsQd4gFMMe0UB1UOx9bCsh8uOiKLER69eAvVXvc8P2yc4Iig42Bv7JD2zJxhyFALyTKBHipg=="
	);
	
	const ucan2 = UCAN.issue({
	  issuer: alice,
	  audience: bob,
	  capabilities: [
		{
		  can: "fs/read",
		  with: `storage://${alice.did()}/public/photos/`,
		},
		{
		  can: "pin/add",
		  with: alice.did(),
		},
	  ],
	});
	console.log(alice);
	ucan2.then(
		(res)=>{
			console.log(res.issuer.did());
			console.log(UCAN.format(res));
			UCAN.link(res, {hasher: identity}).then((res2)=>{
				console.log(res2.toV1());
				console.log(res2.toV1().toString());
				console.log(dagdecode(UCAN.encode(res)));
			});
		}
	);
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
