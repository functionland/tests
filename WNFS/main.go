package main

import (
	"context"
	"io/fs"

	base "github.com/functionland/wnfs-go/base"
	mockblocks "github.com/functionland/wnfs-go/mockblocks"
	private "github.com/functionland/wnfs-go/private"
	ratchet "github.com/functionland/wnfs-go/private/ratchet"
	public "github.com/functionland/wnfs-go/public"
	blockservice "github.com/ipfs/go-blockservice"
	cid "github.com/ipfs/go-cid"
)

type fileSystem struct {
	store public.Store
	ctx   context.Context
	root  *rootTree
}

type (
	Node         = base.Node
	HistoryEntry = base.HistoryEntry
	PrivateName  = private.Name
	Key          = private.Key
)

type rootHeader struct {
	Info     *public.Info
	Previous *cid.Cid
	Metadata *cid.Cid
	Pretty   *cid.Cid
	Public   *cid.Cid
	Private  *cid.Cid
}

type rootTree struct {
	store   public.Store
	pstore  private.Store
	id      cid.Cid
	tx      cid.Cid // transaction start CID
	rootKey Key

	h *rootHeader

	// Pretty   *base.BareTree
	metadata *public.LDFile
	Public   *public.Tree
	Private  *private.Root
}

type PosixFS interface {
	// directories (trees)
	Ls(pathStr string) ([]fs.DirEntry, error)
	Mkdir(pathStr string) error

	// files
	Write(pathStr string, f fs.File) error
	Cat(pathStr string) ([]byte, error)
	Open(pathStr string) (fs.File, error)

	// general
	// Mv(from, to string) error
	Cp(pathStr, srcPathStr string, src fs.FS) error
	Rm(pathStr string) error
}

type PrivateFS interface {
	RootKey() private.Key
	PrivateName() (PrivateName, error)
}

type CommitResult struct {
	Root        cid.Cid
	PrivateName *PrivateName
	PrivateKey  *Key
}

type WNFS interface {
	fs.FS
	fs.ReadDirFile // wnfs root is a directory file
	PosixFS
	PrivateFS

	Cid() cid.Cid
	History(ctx context.Context, pathStr string, generations int) ([]HistoryEntry, error)
	Commit() (CommitResult, error)
}

const (
	// PreviousLinkName is the string for a historical backpointer in wnfs
	PreviousLinkName = "previous"
	// FileHierarchyNamePrivate is the root of encrypted files on WNFS
	FileHierarchyNamePrivate = "private"
	// FileHierarchyNamePublic is the root of public files on WNFS
	FileHierarchyNamePublic = "public"
	// FileHierarchyNamePretty is a link to a read-only branch at the root of a WNFS
	FileHierarchyNamePretty = "p"
)

var testRootKey Key = [32]byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2,
}

func newMemTestStore(ctx context.Context) public.Store {
	return public.NewStore(ctx, mockblocks.NewOfflineMemBlockservice())
}

func newEmptyRootTree(store public.Store, rs ratchet.Store, rootKey Key) (root *rootTree, err error) {
	root = &rootTree{
		store:   store,
		rootKey: rootKey,

		h: &rootHeader{
			Info: public.NewInfo(base.NTDir),
		},
		Public: public.NewEmptyTree(store, FileHierarchyNamePublic),
		// Pretty: &base.BareTree{},
	}

	root.pstore, err = private.NewStore(context.TODO(), store.Blockservice(), rs)
	if err != nil {
		return nil, err
	}

	privateRoot, err := private.NewEmptyRoot(store.Context(), root.pstore, FileHierarchyNamePrivate, rootKey)
	if err != nil {
		return nil, err
	}
	root.Private = privateRoot
	return root, nil
}

func NewEmptyFS(ctx context.Context, bserv blockservice.BlockService, rs ratchet.Store, rootKey Key) (WNFS, error) {
	store := public.NewStore(ctx, bserv)
	fs := &fileSystem{
		ctx:   ctx,
		store: store,
	}

	root, err := newEmptyRootTree(store, rs, rootKey)
	if err != nil {
		return nil, err
	}

	fs.root = root

	// put all root tree to establish base hashes for all top level directories in
	// the file hierarchy
	if _, err := root.Public.Put(); err != nil {
		return nil, err
	}
	if _, err := root.Private.Put(); err != nil {
		return nil, err
	}

	return fs, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := newMemTestStore(ctx)
	rs := ratchet.NewMemStore(ctx)
	fsys, err := NewEmptyFS(ctx, store.Blockservice(), rs, testRootKey)
	if err != nil {
		println("oh")
	}

	pathStr := "public/foo/hello.txt"
	fileContents := []byte("hello!")
	f := base.NewMemfileBytes("hello.txt", fileContents)

	err = fsys.Write(pathStr, f)
	_, err = fsys.Commit()

}
