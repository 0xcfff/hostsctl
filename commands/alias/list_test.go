package alias

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/0xcfff/hostsctl/commands/common"
	"github.com/0xcfff/hostsctl/hosts"
	"github.com/0xcfff/hostsctl/testtools"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestAliasListCommand(t *testing.T) {
	type args struct {
		args       []string
		inputFile  string
		outputFile string
		errorText  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// default
		{
			"default empty",
			args{
				[]string{},
				"testdata/empty.txt",
				"testdata/list/default_empty.txt",
				"",
			},
			true,
		},
		{
			"default one line",
			args{
				[]string{},
				"testdata/one-ip.txt",
				"testdata/list/default_one-ip.txt",
				"",
			},
			true,
		},
		{
			"default two blocks",
			args{
				[]string{},
				"testdata/two-sys-blocks.txt",
				"testdata/list/default_two-sys-blocks.txt",
				"",
			},
			true,
		},
		// short
		{
			"plain empty",
			args{
				[]string{"-o", "short"},
				"testdata/empty.txt",
				"testdata/list/short_empty.txt",
				"",
			},
			true,
		},
		{
			"plain one line",
			args{
				[]string{"-o", "short"},
				"testdata/one-ip.txt",
				"testdata/list/short_one-ip.txt",
				"",
			},
			true,
		},
		{
			"plain two blocks",
			args{
				[]string{"-o", "short"},
				"testdata/two-sys-blocks.txt",
				"testdata/list/short_two-sys-blocks.txt",
				"",
			},
			true,
		},
		// wide
		{
			"wide empty",
			args{
				[]string{"-o", "wide"},
				"testdata/empty.txt",
				"testdata/list/wide_empty.txt",
				"",
			},
			true,
		},
		{
			"wide one line",
			args{
				[]string{"-o", "wide"},
				"testdata/one-ip.txt",
				"testdata/list/wide_one-ip.txt",
				"",
			},
			true,
		},
		{
			"wide two blocks",
			args{
				[]string{"-o", "wide"},
				"testdata/two-sys-blocks.txt",
				"testdata/list/wide_two-sys-blocks.txt",
				"",
			},
			true,
		},
		// plain
		{
			"plain empty",
			args{
				[]string{"-o", "plain"},
				"testdata/empty.txt",
				"testdata/list/plain_empty.txt",
				"",
			},
			true,
		},
		{
			"plain one line",
			args{
				[]string{"-o", "plain"},
				"testdata/one-ip.txt",
				"testdata/list/plain_one-ip.txt",
				"",
			},
			true,
		},
		{
			"plain two blocks",
			args{
				[]string{"-o", "plain"},
				"testdata/two-sys-blocks.txt",
				"testdata/list/plain_two-sys-blocks.txt",
				"",
			},
			true,
		},
		// json set of tests
		{
			"json empty",
			args{
				[]string{"-o", "json"},
				"testdata/empty.txt",
				"testdata/list/json_empty.txt",
				"",
			},
			true,
		},
		{
			"json one line",
			args{
				[]string{"-o", "json"},
				"testdata/one-ip.txt",
				"testdata/list/json_one-ip.txt",
				"",
			},
			true,
		},
		// yaml set of tests
		{
			"yaml empty",
			args{
				[]string{"-o", "yaml"},
				"testdata/empty.txt",
				"testdata/list/yaml_empty.txt",
				"",
			},
			true,
		},
		{
			"yaml one line",
			args{
				[]string{"-o", "yaml"},
				"testdata/one-ip.txt",
				"testdata/list/yaml_one-ip.txt",
				"",
			},
			true,
		},
		// arrangement cases
		{
			"arange - two blocks - raw",
			args{
				[]string{"-a", "raw"},
				"testdata/two-mixed-blocks.txt",
				"testdata/list/arange_two-mixed-blocks_raw.txt",
				"",
			},
			true,
		},
		{
			"arange - two blocks - group",
			args{
				[]string{"-a", "group"},
				"testdata/two-mixed-blocks.txt",
				"testdata/list/arange_two-mixed-blocks_group.txt",
				"",
			},
			true,
		},
		{
			"arange - two blocks - ungroup",
			args{
				[]string{"-a", "ungroup"},
				"testdata/two-mixed-blocks.txt",
				"testdata/list/arange_two-mixed-blocks_ungroup.txt",
				"",
			},
			true,
		},

		// wrong arguments check
		{
			"error - wrong format",
			args{
				[]string{"-o", "not_supported"},
				"testdata/empty.txt",
				"",
				"Error: value not_supported is not support; not supported output format",
			},
			false,
		},
		{
			"error - wrong grouping",
			args{
				[]string{"-a", "not_supported"},
				"testdata/empty.txt",
				"",
				"Error: value not_supported is not support; wrong argument value",
			},
			false,
		},
	}

	for _, tt := range tests {
		inHelperProcess := os.Getenv("GO_TEST_HELPER_PROCESS") == "1"
		if inHelperProcess {
			testName := os.Getenv("GO_TEST_TEST_NAME")
			if !strings.EqualFold(testName, tt.name) {
				continue
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if !tt.want && !inHelperProcess {
				tstp := testtools.RunHelperProcess("TestAliasListCommand", tt.name)
				out, _ := tstp.CombinedOutput()
				fmt.Println(string(out))
				assert.NotEqual(t, 0, tstp.ProcessState.ExitCode())
				assert.Contains(t, string(out), tt.args.errorText)
				return
			}

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

			expectData := bytes.NewBufferString("").Bytes()
			if tt.args.outputFile != "" {
				expectData, err = os.ReadFile(tt.args.outputFile)
				if err != nil {
					t.Errorf("Can't read %v", tt.args.outputFile)
					t.FailNow()
				}
			}
			expectOut := string(expectData)

			ctx := common.WithCustomFilesystem(context.Background(), fs)
			out := &strings.Builder{}

			cmd := NewCmdAliasList()
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
		})
	}
}
