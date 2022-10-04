package main

import (
	"context"

	wnfs "github.com/functionland/wnfs-go"
	base "github.com/functionland/wnfs-go/base"
	mockblocks "github.com/functionland/wnfs-go/mockblocks"
	private "github.com/functionland/wnfs-go/private"
	ratchet "github.com/functionland/wnfs-go/private/ratchet"
	public "github.com/functionland/wnfs-go/public"
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

func newMemTestStore(ctx context.Context) public.Store {
	return public.NewStore(ctx, mockblocks.NewOfflineMemBlockservice())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := newMemTestStore(ctx)
	rs := ratchet.NewMemStore(ctx)
	fsys, err := wnfs.NewEmptyFS(ctx, store.Blockservice(), rs, testRootKey)
	if err != nil {
		println("oh")
	}

	pathStr := "public/foo/hello.txt"
	fileContents := []byte("hello!")
	f := base.NewMemfileBytes("hello.txt", fileContents)

	err = fsys.Write(pathStr, f)
	_, err = fsys.Commit()

}
