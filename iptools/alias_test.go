package iptools

import "testing"

func TestIsSystemAlias(t *testing.T) {
	type args struct {
		ip    string
		alias string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"localhost", args{"127.0.0.1", "localhost"}, true},
		{"localhost umixed registry", args{"127.0.0.1", "LocalHost"}, true},
		{"127.0.0.01 localhost", args{"127.0.0.01", "localhost"}, false}, // not supported by go
		{"127.0.0.1 to custom", args{"127.0.0.1", "custom"}, false},
		{"public to custom", args{"55.03.04.99", "custom"}, false},
		{"::1 to ipv6-localhost", args{"::1", "ip6-localhost"}, true},
		{"::01 to ipv6-localhost", args{"::01", "ip6-localhost"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSystemAlias(tt.args.ip, tt.args.alias); got != tt.want {
				t.Errorf("IsSystemAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}
