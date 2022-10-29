# React-WNFS
This is a repository to test features of wnfs (the wasm version published as fx-wnfs) inside a react app.

## Usage
install dependencies
```bash
npm install
```

start dev server
```
npm start
```

The `App.js` file is for testing the `webnative` package, it is written in typescript and needs to init the webnative for permissions and proofs.

The `Test.js` file is for testing the `fx-wnfs` package, it is the wasm version.