package core

import (
	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/ipfs/go-filestore"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestOpenRepo(t *testing.T) {
	type args struct {
		cctx *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *Repo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OpenRepo(tt.args.cctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenRepo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Close(t *testing.T) {
	type fields struct {
		DB        *gorm.DB
		Filestore *filestore.Filestore
		Dir       string
		leveldb   *leveldb.Datastore
		Cfg       *viper.Viper
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				DB:        tt.fields.DB,
				Filestore: tt.fields.Filestore,
				Dir:       tt.fields.Dir,
				leveldb:   tt.fields.leveldb,
				Cfg:       tt.fields.Cfg,
			}
			if err := r.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findRepo(t *testing.T) {
	type args struct {
		cctx *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findRepo(tt.args.cctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("findRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findRepo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
