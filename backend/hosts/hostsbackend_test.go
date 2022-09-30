package hosts

import (
	"reflect"
	"testing"

	"github.com/0xcfff/dnssync/model"
	"github.com/spf13/afero"
)

func Test_hostsBackend_ReadState(t *testing.T) {
	type fields struct {
		etcHostsPath string
		fs           afero.Fs
	}
	tests := []struct {
		name    string
		fields  fields
		want    model.BackendState
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := &hostsBackend{
				etcHostsPath: tt.fields.etcHostsPath,
				fs:           tt.fields.fs,
			}
			got, err := backend.ReadState()
			if (err != nil) != tt.wantErr {
				t.Errorf("hostsBackend.ReadState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hostsBackend.ReadState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hostsBackend_UpdateState(t *testing.T) {
	type fields struct {
		etcHostsPath string
		fs           afero.Fs
	}
	type args struct {
		changeSet model.BackendStateChangeSet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.BackendState
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := &hostsBackend{
				etcHostsPath: tt.fields.etcHostsPath,
				fs:           tt.fields.fs,
			}
			got, err := backend.UpdateState(tt.args.changeSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("hostsBackend.UpdateState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hostsBackend.UpdateState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBackend(t *testing.T) {
	tests := []struct {
		name string
		want model.Backend
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultBackend(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultBackend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBackend(t *testing.T) {
	type args struct {
		hostsFilePath *string
		fs            afero.Fs
	}
	tests := []struct {
		name string
		args args
		want model.Backend
	}{
		{"test1", args{nil, afero.NewOsFs()}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := DefaultBackend()
			backend.ReadState()
			// if got := NewBackend(tt.args.hostsFilePath, tt.args.fs); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewBackend() = %v, want %v", got, tt.want)
			// }
		})
	}
}
