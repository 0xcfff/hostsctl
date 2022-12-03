package iptools

import "testing"

func TestIsIPv4(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid ipv4", args{"127.0.0.1"}, true},
		{"fail as ipv6", args{"684D:1111:222:3333:4444:5555:6:77"}, false},
		{"random text", args{"mytext"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv4(tt.args.value); got != tt.want {
				t.Errorf("IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"fail as ipv4", args{"127.0.0.1"}, false},
		{"valid ipv6", args{"684D:1111:222:3333:4444:5555:6:77"}, true},
		{"valid short ipv6", args{"fe00::0"}, true},
		{"random text", args{"mytext"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv6(tt.args.value); got != tt.want {
				t.Errorf("IsIPv6() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIP(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid ipv4", args{"127.0.0.1"}, true},
		{"valid ipv6", args{"684D:1111:222:3333:4444:5555:6:77"}, true},
		{"valid short ipv6", args{"fe00::0"}, true},
		{"random text", args{"mytext"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIP(tt.args.value); got != tt.want {
				t.Errorf("IsIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
