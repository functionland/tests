# React Native Pool Test

This repository is a place for testing new features of `go-fula/pools` inside a react-native application.

It uses `go-fula/mobile` package as a libp2p host that provides pool features inside a react-native application. The `FulaMobile` Java module is compiled from `go-fula/mobile` using `gomobile/gobind` tool.

## Usage
to install the app on an android device
```bash
npm install
npm run android
```

to start the libp2p host (for testing graphsync exchange)
```bash
go run blox.go
```
after running the server, it will show you the info you need to enter in the android app.
- PeerId: The peer ID for this libp2p host
- MultiAddr: The multiaddr for this app. You may need to change the `ma` variable in the `blox.go` depending on your machine's ip address in the LAN.
- root and leaf CIDs: these are the addresses for two sample IPLD nodes that are getting created and stored in the libp2p host. you can use either of them for testing.

put the info in the corresponding inputs in the application and press FETCH to see the IPLD node's data is getting transfered using graphsync.
