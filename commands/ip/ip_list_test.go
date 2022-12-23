package ip

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestIpListCommand(t *testing.T) {
	type args struct {
		args       []string
		inputFile  string
		outputFile string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"empty",
			args{
				[]string{},
				"testdata/empty.txt",
				"testdata/list/default_empty.txt",
			},
			true,
		},
		{
			"one line",
			args{
				[]string{},
				"testdata/one-ip.txt",
				"testdata/list/default_one-ip.txt",
			},
			true,
		},
		{
			"two blocks",
			args{
				[]string{},
				"testdata/two-sys-blocks.txt",
				"testdata/list/default_two-sys-blocks.txt",
			},
			true,
		},
		// json set of tests
		{
			"empty json",
			args{
				[]string{"-o", "json"},
				"testdata/empty.txt",
				"testdata/list/json_empty.txt",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			fs := afero.NewMemMapFs()
			fn := hosts.EtcHosts.Path()
			f, err := fs.Create(fn)
			if err != nil {
				t.Errorf("Can't create %v", fn)
				t.FailNow()
			}
			data, err := os.ReadFile(tt.args.inputFile)
			if err != nil {
				t.Errorf("Can't read %v", tt.args.inputFile)
				t.FailNow()
			}
			sdata := string(data)
			f.WriteString(sdata)
			f.Close()

			expectData, err := os.ReadFile(tt.args.outputFile)
			if err != nil {
				t.Errorf("Can't read %v", tt.args.outputFile)
				t.FailNow()
			}
			expectOut := string(expectData)

			ctx := common.WithCustomFilesystem(context.Background(), fs)
			out := &strings.Builder{}

			cmd := NewCmdIpList()
			cmd.SetArgs(tt.args.args)
			cmd.SetOutput(out)
			cmd.SetContext(ctx)

			// act
			c, err := cmd.ExecuteC()

			// assert
			if tt.want {
				assert.NoError(t, err, "command should succeed")
			} else {
				assert.Error(t, err, "command should fail")
			}

			assert.Same(t, cmd, c)

			s := out.String()
			assert.Equal(t, expectOut, s)

			// Sample output:
			// Add GRP and SYS columns to output
			// GRP  SYS  IP          ALIAS
			// [+]   +   127.0.0.1   localhost
			//           127.0.0.2   router
			//           127.0.0.3   printer
			// [+]   +   ::1         ip6-localhost
			//       +   ::1         ip6-loopback
			//       +   e00::0      ip6-localnet
			//       +   e00::0      ip6-mcastprefix
			//       +   e00::0      ip6-allnodes
			//       +   e00::0      ip6-allrouters
			// [+]       10.0.0.101  hhost1
			//           10.0.0.102  hhost2
		})
	}
}
