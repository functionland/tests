package main

import (
	"context"
	"fmt"

	wnfs "github.com/functionland/wnfs-go"
	base "github.com/functionland/wnfs-go/base"
	mockblocks "github.com/functionland/wnfs-go/mockblocks"
	private "github.com/functionland/wnfs-go/private"
	ratchet "github.com/functionland/wnfs-go/private/ratchet"
	cmp "github.com/google/go-cmp/cmp"
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

func newFileTestStorePrivate(ctx context.Context) (st private.Store, cleanup func()) {
	bserv, cleanup, err := mockblocks.NewOfflineFileBlockservice("test")
	if err != nil {
		println(err.Error())
	}
	rs := ratchet.NewMemStore(ctx)
	store, err := private.NewStore(ctx, bserv, rs)
	if err != nil {
		println(err)
	}
	return store, cleanup
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rs := ratchet.NewMemStore(ctx)
	store, cleanup := newFileTestStorePrivate(ctx)
	defer cleanup()
	fsys, err := wnfs.NewEmptyFS(ctx, store.Blockservice(), rs, testRootKey)
	if err != nil {
		println(err)
	}
	fmt.Printf("wnfs root CID: %s\n", fsys.Cid())

	pathStr := "private/foo/hello.txt"
	fileContents := []byte("hello!")
	f := base.NewMemfileBytes("hello.txt", fileContents)

	err1 := fsys.Write(pathStr, f)
	if err1 != nil {
		fmt.Printf("Erro happened %s\n", err1.Error())
	}
	fmt.Printf("wnfs new root CID: %s\n", fsys.Cid())

	gotFileContents, err := fsys.Cat(pathStr)
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("%s\n", string(gotFileContents))

	if diff := cmp.Diff(fileContents, gotFileContents); diff != "" {
		fmt.Printf("\nresult mismatch. (-want +got):\n%s", diff)
	}

	ls, err := fsys.Ls("private/foo")
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("\nFolder structure: %s", ls)

	println("\nEnd")

}
