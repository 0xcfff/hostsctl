package model

import (
	"reflect"
	"testing"
)

func Test_sourceConfig_Property(t *testing.T) {

	emptyProps := make(map[string]string)

	sourceOnlyProps := make(map[string]string)
	sourceOnlyProps["source"] = "http"

	type fields struct {
		properties map[string]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue string
		wantOk    bool
	}{
		// TODO: Add test cases.
		{name: "source missing", fields: fields{properties: emptyProps}, args: args{name: "source"}, wantValue: "", wantOk: false},
		{name: "source exists", fields: fields{properties: sourceOnlyProps}, args: args{name: "source"}, wantValue: "http", wantOk: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sourceConfig{
				properties: tt.fields.properties,
			}
			gotValue, gotOk := s.Property(tt.args.name)
			if gotValue != tt.wantValue {
				t.Errorf("sourceConfig.Property() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("sourceConfig.Property() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_sourceConfig_Properties(t *testing.T) {
	type fields struct {
		properties map[string]string
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sourceConfig{
				properties: tt.fields.properties,
			}
			if got := s.Properties(tt.args.names); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sourceConfig.Properties() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sourceConfig_VerifySchema(t *testing.T) {
	type fields struct {
		properties map[string]string
	}
	type args struct {
		schema *ConfigSchema
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
			s := &sourceConfig{
				properties: tt.fields.properties,
			}
			if err := s.VerifySchema(tt.args.schema); (err != nil) != tt.wantErr {
				t.Errorf("sourceConfig.VerifySchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sourceConfig_ConfigHash(t *testing.T) {
	type fields struct {
		properties map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sourceConfig{
				properties: tt.fields.properties,
			}
			if got := s.ConfigHash(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sourceConfig.ConfigHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsourceConfig(t *testing.T) {
	type args struct {
		properties map[string]string
	}
	tests := []struct {
		name string
		args args
		want SourceConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSourceConfig(tt.args.properties); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewsourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
