package core

import (
	"context"
	"fmt"
	"github.com/application-research/estuary/util"
	"github.com/ipfs/go-bitswap"
	bsnet "github.com/ipfs/go-bitswap/network"
	"github.com/ipfs/go-cidutil"
	chunker "github.com/ipfs/go-ipfs-chunker"
	metri "github.com/ipfs/go-metrics-interface"
	"github.com/ipfs/go-unixfs/importer/balanced"
	ihelper "github.com/ipfs/go-unixfs/importer/helpers"
	unixfsio "github.com/ipfs/go-unixfs/io"
	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/metrics"
	rhelp "github.com/libp2p/go-libp2p-routing-helpers"
	mh "github.com/multiformats/go-multihash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	"github.com/ipfs/go-filestore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbor "github.com/ipfs/go-ipld-cbor"
	format "github.com/ipfs/go-ipld-format"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/application-research/estuary/pinner/types"
	dagsplit "github.com/application-research/estuary/util/dagsplit"
)

var PlumbCmd = &cli.Command{
	Name:        "plumb",
	Hidden:      true,
	Description: "low level plumbing commands",
	Usage:       "plumb <command> [<args>]",
	Subcommands: []*cli.Command{
		PlumbPutFileCmd,
		PlumbPutCarCmd,
		PlumbSplitAddFileCmd,
		PlumbPutDirCmd,
	},
}

var PlumbPutFileCmd = &cli.Command{
	Name:  "put-file",
	Usage: "put-file <file> [<name>]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "specify alternate name for file to be added with",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "specify password to encrypt the file with a password",
		},
	},
	Action: func(cctx *cli.Context) error {
		if !cctx.Args().Present() {
			return fmt.Errorf("must specify filename to upload")
		}

		c, err := LoadClient(cctx)
		if err != nil {
			return err
		}

		f := cctx.Args().First()
		fname := filepath.Base(f)
		if oname := cctx.String("name"); oname != "" {
			fname = oname
		}

		resp, err := c.AddFile(f, fname)
		if err != nil {
			return err
		}

		fmt.Println(resp.Cid)
		return nil
	},
}

func PlumbAddFile(ctx *cli.Context, fpath string, fname string) (*util.ContentAddResponse, error) {
	c, err := LoadClient(ctx)
	if err != nil {
		return nil, err
	}

	return c.AddFile(fpath, fname)
}

func PlumbAddCar(ctx *cli.Context, fpath string, fname string) (*util.ContentAddResponse, error) {
	c, err := LoadClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.AddCar(fpath, fname)

}

var PlumbPutDirCmd = &cli.Command{
	Name: "put-dir",
	Action: func(cctx *cli.Context) error {
		ctx := cctx.Context
		client, err := LoadClient(cctx)
		if err != nil {
			return err
		}

		ds := dsync.MutexWrap(datastore.NewMapDatastore())
		fsm := filestore.NewFileManager(ds, "/")
		bs := blockstore.NewBlockstore(ds)

		fsm.AllowFiles = true
		fstore := filestore.NewFilestore(bs, fsm)
		dserv := merkledag.NewDAGService(blockservice.New(fstore, nil))
		fname := cctx.Args().First()

		dnd, err := addDirectory(ctx, fstore, dserv, fname)

		if err != nil {
			return err
		}

		fmt.Println("imported directory: ", dnd.Cid())
		return doAddPin(ctx, fstore, client, dnd.Cid(), fname)
	},
}

var PlumbPutCarCmd = &cli.Command{
	Name: "put-car",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "specify alternate name for file to be added with",
		},
	},
	Action: func(cctx *cli.Context) error {
		if !cctx.Args().Present() {
			return fmt.Errorf("must specify car file to upload")
		}

		c, err := LoadClient(cctx)
		if err != nil {
			return err
		}

		c.DoProgress = true

		f := cctx.Args().First()
		fname := filepath.Base(f)
		if oname := cctx.String("name"); oname != "" {
			fname = oname
		}

		resp, err := c.AddCar(f, fname)
		if err != nil {
			return err
		}

		fmt.Println(resp.Cid)
		return nil
	},
}

var PlumbSplitAddFileCmd = &cli.Command{
	Name: "split-add",
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:  "chunk",
			Value: uint64(abi.PaddedPieceSize(16 << 30).Unpadded()),
		},
		&cli.BoolFlag{
			Name: "no-pin-only-split",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := cctx.Context
		client, err := LoadClient(cctx)
		if err != nil {
			return err
		}

		ds := dsync.MutexWrap(datastore.NewMapDatastore())
		fsm := filestore.NewFileManager(ds, "/")

		bs := blockstore.NewBlockstore(ds)

		fsm.AllowFiles = true
		fstore := filestore.NewFilestore(bs, fsm)
		cst := cbor.NewCborStore(fstore)

		fname := cctx.Args().First()

		progcb := func(int64) {}
		nd, _, err := filestoreAdd(fstore, fname, progcb)
		if err != nil {
			return err
		}

		fmt.Println("imported file: ", nd.Cid())

		dserv := merkledag.NewDAGService(blockservice.New(fstore, nil))
		builder := dagsplit.NewBuilder(dserv, cctx.Uint64("chunk"), 0)

		if err := builder.Pack(ctx, nd.Cid()); err != nil {
			return err
		}

		for i, box := range builder.Boxes() {
			cc, err := cst.Put(ctx, box)
			if err != nil {
				return err
			}

			tsize := 0
			/* old way, maybe wrong?
			if err := merkledag.Walk(ctx, dserv.GetLinks, cc, func(c cid.Cid) bool {
				size, err := fstore.GetSize(c)
				if err != nil {
					panic(err)
				}

				tsize += size
				return true
			}); err != nil {
				return err
			}
			*/
			cset := cid.NewSet()
			if err := merkledag.Walk(ctx, func(ctx context.Context, c cid.Cid) ([]*ipld.Link, error) {
				node, err := dserv.Get(ctx, c)
				if err != nil {
					return nil, err
				}

				tsize += len(node.RawData())

				return node.Links(), nil
			}, cc, cset.Visit); err != nil {
				return err
			}
			fmt.Printf("%d: %s %d\n", i, cc, tsize)
		}

		if cctx.Bool("no-pin-only-split") {
			return nil
		}
		/*
			if err := builder.Add(cctx.Context, nd.Cid()); err != nil {
				return err
			}
		*/

		pc, err := setupBitswap(ctx, fstore)
		if err != nil {
			return err
		}

		h := pc.host

		var addrs []string
		for _, a := range h.Addrs() {
			addrs = append(addrs, fmt.Sprintf("%s/p2p/%s", a, h.ID()))
		}
		fmt.Println("addresses: ", addrs)

		basename := filepath.Base(fname)

		var pins []string
		var cids []cid.Cid
		for i, box := range builder.Boxes() {
			cc, err := cst.Put(ctx, box)
			if err != nil {
				return err
			}

			cids = append(cids, cc)
			fmt.Println("box: ", i, cc)

			st, err := client.PinAdd(ctx, cc, fmt.Sprintf("%s-%d", basename, i), addrs, nil)
			if err != nil {
				return xerrors.Errorf("failed to pin box %d to estuary: %w", i, err)
			}

			if err := connectToDelegates(ctx, h, st.Delegates); err != nil {
				fmt.Println("failed to connect to pin delegates: ", err)
			}

			pins = append(pins, st.RequestID)
		}

		for range time.Tick(time.Second * 2) {
			var pinning, queued, pinned, failed int
			for _, p := range pins {
				status, err := client.PinStatus(ctx, p)
				if err != nil {
					fmt.Println("error getting pin status: ", err)
					continue
				}

				switch status.Status {
				case types.PinningStatusPinned:
					pinned++
				case types.PinningStatusFailed:
					failed++
				case types.PinningStatusPinning:
					pinning++
				case types.PinningStatusQueued:
					queued++
				}

				if err := connectToDelegates(ctx, h, status.Delegates); err != nil {
					fmt.Println("failed to connect to pin delegates: ", err)
				}
			}

			fmt.Printf("pinned: %d, pinning: %d, queued: %d, failed: %d (num conns: %d)\n", pinned, pinning, queued, failed, len(h.Network().Conns()))
			if failed+pinned >= len(pins) {
				break
			}
		}

		fmt.Println("finished pinning: ", nd.Cid())

		return nil
	},
}

func setupBitswap(ctx context.Context, bstore blockstore.Blockstore) (*PinClient, error) {
	cmgr, err := connmgr.NewConnManager(2000, 3000)
	if err != nil {
		return nil, err
	}

	bwc := metrics.NewBandwidthCounter()
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(cmgr),
		//libp2p.Identity(peerkey),
		libp2p.BandwidthReporter(bwc),
		libp2p.DefaultTransports,
	)
	if err != nil {
		return nil, err
	}

	bsNetFromIpfsHost := bsnet.NewFromIpfsHost(h, rhelp.Null{})
	bsctx := metri.CtxScope(ctx, "barge.exch")

	bswap := bitswap.New(bsctx, bsNetFromIpfsHost, bstore,
		bitswap.EngineBlockstoreWorkerCount(600),
		bitswap.TaskWorkerCount(600),
		bitswap.MaxOutstandingBytesPerPeer(10<<20),
	)

	return &PinClient{
		host:    h,
		bitswap: bswap,
		bwc:     bwc,
	}, nil
}

func addDirectory(ctx context.Context, fstore *filestore.Filestore, dserv ipld.DAGService, dir string) (*merkledag.ProtoNode, error) {

	//	create the root directory
	dirNode := unixfsio.NewDirectory(dserv)
	prefix, err := merkledag.PrefixForCidVersion(1)
	prefix.MhType = mh.SHA2_256
	dirNode.SetCidBuilder(cidutil.InlineBuilder{
		Builder: prefix,
	})

	dirents, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	progCb := func(int64) {}

	for _, d := range dirents {
		name := filepath.Join(dir, d.Name())
		if d.IsDir() {

			dirn, err := addDirectory(ctx, fstore, dserv, name)
			if err != nil {
				return nil, err
			}
			if err := dirNode.AddChild(ctx, name, dirn); err != nil {
				return nil, err
			}

			fmt.Printf("imported directory: %s | %s \n", d.Name(), dirn.Cid())
		} else {
			node, _, err := filestoreAdd(fstore, name, progCb)
			if err != nil {
				return nil, err
			}

			if err := dirNode.AddChild(ctx, name, node); err != nil {
				return nil, err
			}
		}
	}
	node, err := dirNode.GetNode()
	stats, err := node.Stat()

	fmt.Println("links", stats.NumLinks)
	return node.(*merkledag.ProtoNode), nil
}

type FilestoreFile struct {
	*os.File
	absPath string
	st      os.FileInfo
	cb      func(int64)
}

func (ff *FilestoreFile) AbsPath() string {
	return ff.absPath
}

func (ff *FilestoreFile) Size() (int64, error) {
	finfo, err := ff.File.Stat()
	if err != nil {
		return 0, err
	}

	return finfo.Size(), nil
}

func (ff *FilestoreFile) Stat() os.FileInfo {
	return ff.st
}

func (ff *FilestoreFile) Read(b []byte) (int, error) {
	n, err := ff.File.Read(b)
	ff.cb(int64(n))
	return n, err
}

func newFF(fpath string, cb func(int64)) (*FilestoreFile, error) {
	fi, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	absp, err := filepath.Abs(fpath)
	if err != nil {
		return nil, err
	}

	st, err := fi.Stat()
	if err != nil {
		return nil, err
	}

	return &FilestoreFile{
		File:    fi,
		absPath: absp,
		st:      st,
		cb:      cb,
	}, nil
}

func importFile(dserv ipld.DAGService, fi io.Reader) (ipld.Node, error) {
	prefix, err := merkledag.PrefixForCidVersion(1)

	if err != nil {
		return nil, err
	}
	prefix.MhType = mh.SHA2_256

	spl := chunker.NewSizeSplitter(fi, 1024*1024)
	dbp := ihelper.DagBuilderParams{
		Maxlinks:  1024,
		RawLeaves: true,

		CidBuilder: cidutil.InlineBuilder{
			Builder: prefix,
			Limit:   32,
		},

		Dagserv: dserv,
		NoCopy:  true,
	}

	db, err := dbp.New(spl)
	if err != nil {
		return nil, err
	}

	nd, err := balanced.Layout(db)
	if err != nil {
		return nil, err
	}
	//fmt.Println("imported file ----: ", nd.Cid())

	return nd, err
}
func filestoreAdd(fstore *filestore.Filestore, fpath string, progcb func(int64)) (format.Node, uint64, error) {
	ff, err := newFF(fpath, progcb)
	if err != nil {
		return nil, 0, err
	}
	defer func(ff *FilestoreFile) {
		err := ff.Close()
		if err != nil {

		}
	}(ff)

	dserv := merkledag.NewDAGService(blockservice.New(fstore, nil))
	nd, err := importFile(dserv, ff)

	if err != nil {
		return nil, 0, err
	}

	size, err := nd.Size()
	if err != nil {
		return nil, 0, err
	}
	fmt.Printf("imported file: %s | %s \n", fpath, nd.Cid())
	return nd, size, nil
}

func connectToDelegates(ctx context.Context, h host.Host, delegates []string) error {
	peers := make(map[peer.ID][]multiaddr.Multiaddr)
	for _, d := range delegates {
		ai, err := peer.AddrInfoFromString(d)
		if err != nil {
			return err
		}

		peers[ai.ID] = append(peers[ai.ID], ai.Addrs...)
	}

	for p, addrs := range peers {
		h.Peerstore().AddAddrs(p, addrs, time.Hour)

		if h.Network().Connectedness(p) != network.Connected {
			if err := h.Connect(ctx, peer.AddrInfo{
				ID: p,
			}); err != nil {
				return err
			}

			h.ConnManager().Protect(p, "pinning")
		}
	}

	return nil
}

func doAddPin(ctx context.Context, bstore blockstore.Blockstore, client *EstClient, root cid.Cid, fname string) error {
	pc, err := setupBitswap(ctx, bstore)
	if err != nil {
		return err
	}

	h := pc.host

	var addrs []string
	for _, a := range h.Addrs() {
		addrs = append(addrs, fmt.Sprintf("%s/p2p/%s", a, h.ID()))
	}
	fmt.Println("addresses: ", addrs)

	basename := filepath.Base(fname)

	st, err := client.PinAdd(ctx, root, basename, addrs, nil)
	if err != nil {
		return xerrors.Errorf("failed to pin %s to estuary: %w", root, err)
	}

	fmt.Println("Delegates: ", st.Delegates)
	if err := connectToDelegates(ctx, h, st.Delegates); err != nil {
		fmt.Println("failed to connect to pin delegates: ", err)
	}

	pins := []string{st.RequestID}
	for range time.Tick(time.Second * 2) {
		var pinning, queued, pinned, failed int
		for _, p := range pins {
			status, err := client.PinStatus(ctx, p)
			if err != nil {
				fmt.Println("error getting pin status: ", err)
				continue
			}

			switch status.Status {
			case types.PinningStatusPinned:
				pinned++
			case types.PinningStatusFailed:
				failed++
			case types.PinningStatusPinning:
				pinning++
			case types.PinningStatusQueued:
				queued++
			}

			if err := connectToDelegates(ctx, h, status.Delegates); err != nil {
				fmt.Println("failed to connect to pin delegates: ", err)
			}
		}

		st := pc.bwc.GetBandwidthForProtocol("/ipfs/bitswap/1.2.0")
		fmt.Printf("pinned: %d, pinning: %d, queued: %d, failed: %d, xfer rate: %s/s (num conns: %d)\n", pinned, pinning, queued, failed, humanize.Bytes(uint64(st.RateOut)), len(h.Network().Conns()))
		if failed+pinned >= len(pins) {
			break
		}
	}

	fmt.Println("finished pinning: ", root)

	return nil

}
