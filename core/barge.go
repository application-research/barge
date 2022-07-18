package core

import (
	"context"
	"fmt"
	"github.com/application-research/estuary/pinner/types"
	"github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-filestore"
	"github.com/labstack/gommon/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/rand"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var BargeAddCmd = &cli.Command{
	Name:        "add",
	Description: `'barge add <file>' is a command to add a file'`,
	Usage:       "barge add <file>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "progress",
		},
	},
	Action: func(cctx *cli.Context) error {
		r, err := OpenRepo()
		if err != nil {
			return err
		}

		progress := cctx.Bool("progress")

		var paths []string
		// TODO: this expansion could be done in parallel to speed things up on large directories
		for _, f := range cctx.Args().Slice() {
			matches, err := filepath.Glob(f)
			if err != nil {
				return err
			}

			for _, m := range matches {
				// TODO: reuse these stats...
				st, err := os.Stat(m)
				if err != nil {
					return err
				}

				if st.IsDir() {
					sub, err := expandDirectory(m)
					if err != nil {
						return err
					}
					// expand!
					paths = append(paths, sub...)
				} else {
					paths = append(paths, m)
				}
			}
		}

		progcb := func(int64) {}
		incrTotal := func(int64) {}
		finish := func() {}

		if progress {
			bar := pb.New64(0)

			bar.Set(pb.Bytes, true)
			bar.SetTemplate(pb.Full)
			bar.Start()

			progcb = func(amt int64) {
				bar.Add64(amt)
			}

			var total int64
			var totlk sync.Mutex

			incrTotal = func(amt int64) {
				totlk.Lock()
				total += amt
				bar.SetTotal(total)
				totlk.Unlock()
			}

			finish = func() {
				bar.Finish()
			}

		}

		type addJob struct {
			Path  string
			Found []File
			Stat  os.FileInfo
		}

		type updateJob struct {
			Path  string
			Found []File
			Stat  os.FileInfo
			Cid   cid.Cid
		}

		tocheck := make(chan string, 1)
		tobuffer := make(chan *addJob, 128)
		toadd := make(chan *addJob)
		toupdate := make(chan updateJob, 128)

		go func() {
			defer close(tocheck)
			for _, f := range paths {
				tocheck <- f
			}
		}()

		go func() {
			defer close(tobuffer)
			for p := range tocheck {
				st, err := os.Stat(p)
				if err != nil {
					fmt.Println(err)
					return
				}

				incrTotal(st.Size())

				var found []File
				if err := r.DB.Find(&found, "path = ?", p).Error; err != nil {
					fmt.Println(err)
					return
				}

				if len(found) > 0 {
					existing := found[0]

					// have it already... check if its changed
					if st.ModTime().Equal(existing.Mtime) {
						// mtime the same, assume its the same file...
						continue
					}
				}

				tobuffer <- &addJob{
					Path:  p,
					Found: found,
					Stat:  st,
				}
			}
		}()

		go func() {
			defer close(toadd)
			var next *addJob
			var buffer []*addJob
			var out chan *addJob
			var inputDone bool

			for {
				select {
				case aj, ok := <-tobuffer:
					if !ok {
						inputDone = true
						if next == nil && len(buffer) == 0 {
							return
						}
						continue
					}
					if out == nil {
						next = aj
						out = toadd
					} else {
						buffer = append(buffer, aj)
					}
				case out <- next:
					if len(buffer) > 0 {
						next = buffer[0]
						buffer = buffer[1:]
					} else {
						out = nil
						next = nil
						if inputDone {
							return
						}
					}
				}
			}
		}()

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				for aj := range toadd {
					nd, _, err := filestoreAdd(r.Filestore, aj.Path, progcb)
					if err != nil {
						fmt.Println(err)
						return
					}

					toupdate <- updateJob{
						Path:  aj.Path,
						Found: aj.Found,
						Cid:   nd.Cid(),
						Stat:  aj.Stat,
					}
				}
			}()
		}

		go func() {
			wg.Wait()
			close(toupdate)
		}()

		var batchCreates []*File
		for uj := range toupdate {
			if len(uj.Found) > 0 {
				existing := uj.Found[0]
				if existing.Cid != uj.Cid.String() {
					if err := r.DB.Model(File{}).Where("id = ?", existing.ID).UpdateColumns(map[string]interface{}{
						"cid":   uj.Cid.String(),
						"mtime": uj.Stat.ModTime(),
					}).Error; err != nil {
						return err
					}
				}

				continue
			}

			abs, err := filepath.Abs(uj.Path)
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(r.Dir, abs)
			if err != nil {
				return err
			}

			batchCreates = append(batchCreates, &File{
				Path:  rel,
				Cid:   uj.Cid.String(),
				Mtime: uj.Stat.ModTime(),
			})

			if len(batchCreates) > 200 {
				if err := r.DB.CreateInBatches(batchCreates, 100).Error; err != nil {
					return err
				}
				batchCreates = nil
			}
		}

		if err := r.DB.CreateInBatches(batchCreates, 100).Error; err != nil {
			return err
		}

		finish()

		return nil
	},
}

var BargeStatusCmd = &cli.Command{
	Name:        "status",
	Description: `'barge status' is a command to check the status of the file'`,
	Usage:       "barge status",
	Action: func(cctx *cli.Context) error {
		r, err := OpenRepo()
		if err != nil {
			return err
		}

		var allfiles []File
		if err := r.DB.Order("path asc").Find(&allfiles).Error; err != nil {
			return err
		}

		fmt.Println("Changes not yet staged:")

		var unpinned []File
		for _, f := range allfiles {
			ch, reason, err := maybeChanged(f)
			if err != nil {
				return err
			}

			var pins []Pin
			if err := r.DB.Find(&pins, "file = ?", f.ID).Error; err != nil {
				return err
			}

			if !ch {
				if len(pins) > 0 {
					pin := pins[0]

					if pin.Status == types.PinningStatusPinned {
						// unchanged and pinned, no need to print anything
						continue
					}
				}

				unpinned = append(unpinned, f)
				continue
			}

			fmt.Printf("\t%s: %s\n", reason, f.Path)
		}

		if len(unpinned) > 0 {
			fmt.Println()
			fmt.Println("Unpinned files:")
			for _, f := range unpinned {
				fmt.Printf("\t%s\n", f.Path)
			}
		}

		return nil
	},
}

var BargeSyncCmd = &cli.Command{
	Name:        "sync",
	Description: `'barge sync' is a command to synchronize the state of the objects in this barge instance'`,
	Usage:       "barge sync",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "recover",
		},
		&cli.Int64Flag{
			Name: "new-pin-limit",
		},
	},
	Action: func(cctx *cli.Context) error {
		ctx := cctx.Context
		r, err := OpenRepo()
		if err != nil {
			return err
		}

		c, err := LoadClient(cctx)
		if err != nil {
			return err
		}

		coluuid := r.Cfg.GetString("collection.uuid")
		if coluuid == "" {
			return fmt.Errorf("barge repo does not have a collection set")
		}

		/*
			var files []File
			if err := r.DB.Find(&files).Error; err != nil {
				return err
			}
		*/

		var filespins []FileWithPin
		if err := r.DB.Model(File{}).Joins("left join pins on pins.file = files.id AND pins.cid = files.cid").Select("files.id as file_id, pins.id as pin_id, path, status, request_id, files.cid as cid").Scan(&filespins).Error; err != nil {
			return err
		}

		pc, err := setupBitswap(ctx, r.Filestore)
		if err != nil {
			return err
		}

		h := pc.host

		var addrs []string
		for _, a := range h.Addrs() {
			addrs = append(addrs, fmt.Sprintf("%s/p2p/%s", a, h.ID()))
		}

		var pinComplete []FileWithPin
		var needsNewPin []FileWithPin
		var inProgress []*Pin
		var checkProgress []FileWithPin
		for _, f := range filespins {
			if f.PinID == 0 {
				needsNewPin = append(needsNewPin, f)
				continue
			}

			if f.Status == types.PinningStatusPinned {
				// TODO: add flag to allow a forced rechecking
				continue
			}

			checkProgress = append(checkProgress, f)
		}

		batchSize := 500
		fmt.Printf("need to check progress of %d pins\n", len(checkProgress))
		for i := 0; i < len(checkProgress); i += batchSize {
			log.Printf("getting pin statuses: %d / %d\n", i, len(checkProgress))
			end := i + batchSize
			if end > len(checkProgress) {
				end = len(checkProgress)
			}

			var reqids []string
			for _, p := range checkProgress[i:end] {
				reqids = append(reqids, p.RequestID)
			}

			resp, err := c.PinStatuses(ctx, reqids)
			if err != nil {
				return fmt.Errorf("failed to recheck pin statuses: %w", err)
			}

			for _, fp := range checkProgress[i:end] {
				st, ok := resp[fp.RequestID]
				if !ok {
					return fmt.Errorf("did not get status back for requestid %s", fp.RequestID)
				}

				switch st.Status {
				case types.PinningStatusPinned:
					pinComplete = append(pinComplete, fp)
					if err := r.DB.Model(Pin{}).Where("id = ?", fp.PinID).UpdateColumn("status", st.Status).Error; err != nil {
						return err
					}
				case types.PinningStatusFailed:
					needsNewPin = append(needsNewPin, fp)
					if err := r.DB.Delete(Pin{ID: fp.PinID}).Error; err != nil {
						return err
					}
				default:
					// pin is technically in progress? do nothing for now
					inProgress = append(inProgress, &Pin{
						ID:        fp.PinID,
						File:      fp.FileID,
						Status:    fp.Status,
						RequestID: fp.RequestID,
					})
				}
			}
		}

		if cctx.Bool("recover") {
			fmt.Println("recovery requested, searching for pins on estuary not tracked locally...")
			for i, nnp := range needsNewPin {
				fmt.Printf("                                \r")
				fmt.Printf("[%d / %d]\r", i, len(needsNewPin))
				// TODO: can batch this
				st, err := c.PinStatusByCid(ctx, []string{nnp.Cid})
				if err != nil {
					fmt.Println("failed to get pin status: ", err)
					continue
				}

				pin, ok := st[nnp.Cid]
				if !ok {
					continue
				}

				if pin.Status == types.PinningStatusFailed {
					// dont bother recording
					continue
				}

				if err := r.DB.Create(&Pin{
					File:      nnp.FileID,
					Cid:       nnp.Cid,
					RequestID: pin.RequestID,
					Status:    pin.Status,
				}).Error; err != nil {
					return err
				}
			}

			return nil
		}

		fmt.Printf("need to make %d new pins\n", len(needsNewPin))
		if lim := cctx.Int64("new-pin-limit"); lim > 0 {
			if int64(len(needsNewPin)) > lim {
				needsNewPin = needsNewPin[:lim]
				fmt.Printf("only making %d for now...\n", lim)
			}
		}

		var dplk sync.Mutex
		var donePins int
		var wg sync.WaitGroup
		newpins := make([]*Pin, len(needsNewPin))
		errs := make([]error, len(needsNewPin))
		sema := make(chan struct{}, 20)
		var delegates []string
		for i := range needsNewPin {
			wg.Add(1)
			go func(ix int) {
				defer wg.Done()

				f := needsNewPin[ix]

				fcid, err := cid.Decode(f.Cid)
				if err != nil {
					errs[ix] = err
					return
				}

				sema <- struct{}{}
				defer func() {
					<-sema
				}()

				resp, err := c.PinAdd(ctx, fcid, filepath.Base(f.Path), addrs, map[string]interface{}{
					"coluuid": coluuid,
					"colpath": "/" + f.Path,
				})
				if err != nil {
					errs[ix] = err
					return
				}

				dplk.Lock()
				delegates = append(delegates, resp.Delegates...)
				donePins++
				fmt.Printf("                                                 \r")
				fmt.Printf("creating new pins %d/%d", donePins, len(needsNewPin))
				dplk.Unlock()

				p := &Pin{
					File:      f.FileID,
					Cid:       fcid.String(),
					RequestID: resp.RequestID,
					Status:    resp.Status,
				}

				newpins[ix] = p
			}(i)
		}
		wg.Wait()

		if err := connectToDelegates(ctx, h, delegates); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "failed to connect to deletegates for new pin: %s\n", err)
			if err != nil {
				return err
			}
		}

		var tocreate []*Pin
		for _, p := range newpins {
			if p != nil {
				tocreate = append(tocreate, p)
				inProgress = append(inProgress, p)
			}
		}

		if len(tocreate) > 0 {
			if err := r.DB.CreateInBatches(tocreate, 100).Error; err != nil {
				return err
			}
		}

		for _, err := range errs {
			if err != nil {
				return err
			}
		}

		fmt.Println()
		fmt.Println("transferring data...")

		complete := make(map[string]bool)
		failed := make(map[string]bool)
		for range time.Tick(time.Second * 2) {

		loopstart:
			var tocheck []string
			for _, p := range inProgress {
				if complete[p.RequestID] || failed[p.RequestID] {
					continue
				}

				tocheck = append(tocheck, p.RequestID)

				if len(tocheck) >= 300 {
					break
				}
			}

			// if we have a lot of pins still to check, start randomly selecting some to look at
			if len(inProgress)-(len(complete)+len(failed)) > batchSize*2 {
				for i := 0; i < 200; i++ {
					p := inProgress[rand.Intn(len(inProgress))]
					if complete[p.RequestID] || failed[p.RequestID] {
						continue
					}

					tocheck = append(tocheck, p.RequestID)
				}
			}

			statuses, err := c.PinStatuses(ctx, tocheck)
			if err != nil {
				return fmt.Errorf("failed to check pin statuses: %w", err)
			}

			var newdone int
			for _, req := range tocheck {
				status, ok := statuses[req]
				if !ok {
					fmt.Printf("didnt get expected pin status back in request: %s\n", req)
					continue
				}

				switch status.Status {
				case types.PinningStatusPinned:
					newdone++
					complete[req] = true
					if err := r.DB.Model(Pin{}).Where("request_id = ?", req).UpdateColumn("status", types.PinningStatusPinned).Error; err != nil {
						return err
					}
				case types.PinningStatusFailed:
					newdone++
					failed[req] = true
					if err := r.DB.Model(Pin{}).Where("request_id = ?", req).Delete(Pin{}).Error; err != nil {
						return err
					}
				default:
				}

				if err := connectToDelegates(ctx, h, status.Delegates); err != nil {
					fmt.Println("failed to connect to pin delegates: ", err)
				}
			}

			st := pc.bwc.GetBandwidthForProtocol("/ipfs/bitswap/1.2.0")
			fmt.Printf("pinned: %d, pinning: %d, failed: %d, xfer rate: %s/s (connections: %d)\n", len(complete), len(inProgress)-(len(complete)+len(failed)), len(failed), humanize.Bytes(uint64(st.RateOut)), len(h.Network().Conns()))

			if len(failed)+len(complete) >= len(inProgress) {
				break
			}

			// dont wait if we get a high enough proportion of new info
			if newdone > 100 {
				goto loopstart
			}
		}

		return nil

	},
}

var BargeCheckCmd = &cli.Command{
	Name:        "check",
	Description: `'barge check' to check the state of the object'`,
	Usage:       "barge check <cid>",
	Action: func(cctx *cli.Context) error {
		r, err := OpenRepo()
		if err != nil {
			return err
		}

		for _, path := range cctx.Args().Slice() {
			var file File
			if err := r.DB.First(&file, "path = ?", path).Error; err != nil {
				return err
			}

			fcid, err := cid.Decode(file.Cid)
			if err != nil {
				return err
			}

			ctx := context.TODO()
			lres := filestore.Verify(ctx, r.Filestore, fcid)
			fmt.Println(lres.Status.String())
			fmt.Println(lres.ErrorMsg)
		}

		return nil
	},
}

var BargeShareCmd = &cli.Command{
	Name:        "share",
	Description: `'barge check' to share objects'`,
	Usage:       "barge share <cid>",
	Action: func(cctx *cli.Context) error {
		r, err := OpenRepo()
		if err != nil {
			return err
		}

		pc, err := setupBitswap(cctx.Context, r.Filestore)
		if err != nil {
			return err
		}

		h := pc.host

		for _, a := range h.Addrs() {
			fmt.Printf("%s/p2p/%s\n", a, h.ID())
		}

		select {}

		//return nil
	},
}

func maybeChanged(f File) (bool, string, error) {
	st, err := os.Stat(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, "deleted", nil
		}
		return false, "", err
	}

	if f.Mtime.Equal(st.ModTime()) {
		return false, "", nil
	}

	return true, "modified", nil
}

func expandDirectory(dir string) ([]string, error) {
	dirents, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, ent := range dirents {
		if strings.HasPrefix(ent.Name(), ".") {
			continue
		}

		if ent.IsDir() {
			sub, err := expandDirectory(filepath.Join(dir, ent.Name()))
			if err != nil {
				return nil, err
			}

			for _, s := range sub {
				out = append(out, s)
			}
		} else {
			out = append(out, filepath.Join(dir, ent.Name()))
		}
	}

	return out, nil
}
