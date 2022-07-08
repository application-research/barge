package core

import (
	"fmt"
	flatfs "github.com/ipfs/go-ds-flatfs"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/ipfs/go-filestore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type Repo struct {
	DB        *gorm.DB
	Filestore *filestore.Filestore
	Dir       string

	leveldb *leveldb.Datastore
	Cfg     *viper.Viper
}

func (r *Repo) Close() error {

	if err := r.leveldb.Close(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, "failed to close leveldb: ", err)
		if err != nil {
			return err
		}
	}

	return nil
}

func findRepo() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	wd = filepath.Clean(wd)

	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	for wd != "/" {
		if wd == home {
			return "", fmt.Errorf("barge directory not found, make sure to run `barge init` first")
		}

		dir := filepath.Join(wd, ".barge")
		st, err := os.Stat(dir)
		if err != nil {
			if !os.IsNotExist(err) {
				return "", err
			}

			wd = filepath.Dir(wd)
			continue
		}

		if !st.IsDir() {
			return "", fmt.Errorf("found .barge, it wasnt a file")
		}

		return dir, nil
	}

	return "", fmt.Errorf("barge directory not found, make sure to run `barge init` first")
}

func OpenRepo() (*Repo, error) {
	dir, err := findRepo()
	if err != nil {
		return nil, err
	}

	reporoot := filepath.Dir(dir)

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := v.WriteConfigAs("config"); err != nil {
				return nil, err
			}
		} else {
			fmt.Printf("read err: %#v\n", err)
			return nil, err
		}
	}

	dbdir := dir
	cfgdbdir, ok := v.Get("database.directory").(string)
	if ok && cfgdbdir != "" {
		dbdir = cfgdbdir
	}

	db, err := gorm.Open(sqlite.Open(filepath.Join(dbdir, "barge.db")), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	errOnMigrateFile := db.AutoMigrate(&File{})
	if errOnMigrateFile != nil {
		return nil, errOnMigrateFile
	}
	errOnMigratePin := db.AutoMigrate(&Pin{})
	if errOnMigratePin != nil {
		return nil, errOnMigratePin
	}

	lds, err := leveldb.NewDatastore(filepath.Join(dbdir, "leveldb"), &leveldb.Options{
		NoSync: true,
	})
	if err != nil {
		return nil, err
	}

	fsmgr := filestore.NewFileManager(lds, reporoot)
	fsmgr.AllowFiles = true

	ffs, err := flatfs.CreateOrOpen(filepath.Join(dbdir, "flatfs"), flatfs.IPFS_DEF_SHARD, false)
	if err != nil {
		return nil, err
	}

	fbs := blockstore.NewBlockstoreNoPrefix(ffs)
	fstore := filestore.NewFilestore(fbs, fsmgr)

	return &Repo{
		DB:        db,
		Filestore: fstore,
		Dir:       reporoot,
		Cfg:       v,
	}, nil
}
