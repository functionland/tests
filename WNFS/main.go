package main

import (
	"context"
	"fmt"

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
type fataler interface {
	Name() string
	Helper()
	Fatal(args ...interface{})
}

func newMemTestStore(ctx context.Context) public.Store {
	return public.NewStore(ctx, mockblocks.NewOfflineMemBlockservice())
}
func newFileTestStore(ctx context.Context) (st public.Store, cleanup func()) {
	bserv, cleanup, err := mockblocks.NewOfflineFileBlockservice("test")
	if err != nil {
		println(err.Error())
	}

	store := public.NewStore(ctx, bserv)
	return store, cleanup
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rs := ratchet.NewMemStore(ctx)
	store, cleanup := newFileTestStore(ctx)
	defer cleanup()
	fsys, err := wnfs.NewEmptyFS(ctx, store.Blockservice(), rs, testRootKey)
	if err != nil {
		println(err)
	}

	pathStr := "private/foo/hello.txt"
	fileContents := []byte("hello!")
	f := base.NewMemfileBytes("hello.txt", fileContents)

	err1 := fsys.Write(pathStr, f)
	if err1 != nil {
		fmt.Printf("Erro happened %s", err1.Error())
	}

	gotFileContents, err := fsys.Cat(pathStr)
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("%s\n", string(gotFileContents))

	println("End")

}
