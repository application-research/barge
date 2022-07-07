package core

import (
	"context"
	"github.com/application-research/estuary/pinner/types"
	"github.com/application-research/estuary/util"
	"github.com/ipfs/go-cid"
	"github.com/urfave/cli/v2"
	"reflect"
	"testing"
)

func TestEstClient_AddCar(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		fpath string
		name  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *util.ContentAddResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.AddCar(tt.args.fpath, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddCar() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_AddFile(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		fpath    string
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *util.ContentAddResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.AddFile(tt.args.fpath, tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_CollectionsCreate(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx  context.Context
		name string
		desc string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.CollectionsCreate(tt.args.ctx, tt.args.name, tt.args.desc)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectionsCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectionsCreate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_CollectionsList(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Collection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.CollectionsList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectionsList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectionsList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_CollectionsListDir(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx     context.Context
		coluuid string
		path    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []collectionListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.CollectionsListDir(tt.args.ctx, tt.args.coluuid, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CollectionsListDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectionsListDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_PinAdd(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx     context.Context
		root    cid.Cid
		name    string
		origins []string
		meta    map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.IpfsPinStatusResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.PinAdd(tt.args.ctx, tt.args.root, tt.args.name, tt.args.origins, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("PinAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PinAdd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_PinStatus(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx   context.Context
		reqid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.IpfsPinStatusResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.PinStatus(tt.args.ctx, tt.args.reqid)
			if (err != nil) != tt.wantErr {
				t.Errorf("PinStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PinStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_PinStatusByCid(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx  context.Context
		cids []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.IpfsPinStatusResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.PinStatusByCid(tt.args.ctx, tt.args.cids)
			if (err != nil) != tt.wantErr {
				t.Errorf("PinStatusByCid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PinStatusByCid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_PinStatuses(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx    context.Context
		reqids []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]*types.IpfsPinStatusResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.PinStatuses(tt.args.ctx, tt.args.reqids)
			if (err != nil) != tt.wantErr {
				t.Errorf("PinStatuses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PinStatuses() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_Viewer(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *util.ViewerResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.Viewer(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Viewer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Viewer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_doRequest(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx    context.Context
		method string
		path   string
		body   interface{}
		resp   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.doRequest(tt.args.ctx, tt.args.method, tt.args.path, tt.args.body, tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("doRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEstClient_doRequestRetries(t *testing.T) {
	type fields struct {
		Host       string
		Shuttle    string
		Tok        string
		DoProgress bool
		LogTimings bool
	}
	type args struct {
		ctx     context.Context
		method  string
		path    string
		body    interface{}
		resp    interface{}
		retries int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EstClient{
				Host:       tt.fields.Host,
				Shuttle:    tt.fields.Shuttle,
				Tok:        tt.fields.Tok,
				DoProgress: tt.fields.DoProgress,
				LogTimings: tt.fields.LogTimings,
			}
			got, err := c.doRequestRetries(tt.args.ctx, tt.args.method, tt.args.path, tt.args.body, tt.args.resp, tt.args.retries)
			if (err != nil) != tt.wantErr {
				t.Errorf("doRequestRetries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("doRequestRetries() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadClient(t *testing.T) {
	type args struct {
		cctx *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *EstClient
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadClient(tt.args.cctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpStatusError_Error(t *testing.T) {
	type fields struct {
		Status     string
		StatusCode int
		Extra      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hse := httpStatusError{
				Status:     tt.fields.Status,
				StatusCode: tt.fields.StatusCode,
				Extra:      tt.fields.Extra,
			}
			if got := hse.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shouldRetry(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRetry(tt.args.err); got != tt.want {
				t.Errorf("shouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}
