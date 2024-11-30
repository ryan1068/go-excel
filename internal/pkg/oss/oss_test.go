package oss

import (
	"cst/internal/pkg/config"
	"flag"
	aliyunOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
		want *oss
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_oss_DeleteObject(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	type args struct {
		objectKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &oss{
				cfg: tt.fields.cfg,
			}
			if err := o.DeleteObject(tt.args.objectKey); (err != nil) != tt.wantErr {
				t.Errorf("DeleteObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_oss_DownloadFile(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	type args struct {
		objectKey string
		filePath  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &oss{
				cfg: tt.fields.cfg,
			}
			got, err := o.DownloadFile(tt.args.objectKey, tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DownloadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var flagConfig = flag.String("config", "../../../configs/dev.yml", "path to the config file")

func Test_oss_UploadFile(t *testing.T) {

	flag.Parse()
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		panic(err)
	}

	type fields struct {
		cfg *config.Config
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			fields: struct{ cfg *config.Config }{cfg: cfg},
			args:   struct{ filePath string }{filePath: "timg.jpg"},
			want:   "https://test-bq.oss-cn-hangzhou.aliyuncs.com/mirco/20240509/timg.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &oss{
				cfg: tt.fields.cfg,
			}
			got, err := o.UploadFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UploadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_oss_bucket(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		want    *aliyunOss.Bucket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &oss{
				cfg: tt.fields.cfg,
			}
			got, err := o.bucket()
			if (err != nil) != tt.wantErr {
				t.Errorf("bucket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bucket() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_oss_client(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		want    *aliyunOss.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &oss{
				cfg: tt.fields.cfg,
			}
			got, err := o.client()
			if (err != nil) != tt.wantErr {
				t.Errorf("client() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client() got = %v, want %v", got, tt.want)
			}
		})
	}
}
