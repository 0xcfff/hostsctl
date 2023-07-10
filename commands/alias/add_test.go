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

func TestAliasAddCommand(t *testing.T) {
	type args struct {
		args       []string
		stdin      string
		inputFile  string
		outputFile string
		stdout     string
		errorText  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"empty - args",
			args{
				[]string{"127.0.0.1", "my.domain.test"},
				"",
				"testdata/empty.txt",
				"testdata/add/empty_one-alias_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"one line - args",
			args{
				[]string{"127.0.0.1", "my.domain.test"},
				"",
				"testdata/one-ip.txt",
				"testdata/add/one-ip_one-alias_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"one line - stdin",
			args{
				[]string{},
				"127.0.0.1 my.domain.test",
				"testdata/one-ip.txt",
				"testdata/add/one-ip_one-alias_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"one line - args + comment",
			args{
				[]string{"127.0.0.1", "my.domain.test", "-c", "My custom service domain"},
				"127.0.0.1 my.domain.test",
				"testdata/one-ip.txt",
				"testdata/add/one-ip_one-alias-and-comment_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"one line - stdin + comment",
			args{
				[]string{},
				"127.0.0.1 my.domain.test # My custom service domain",
				"testdata/one-ip.txt",
				"testdata/add/one-ip_one-alias-and-comment_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"one line - args + missing block",
			args{
				[]string{"127.0.0.1", "my.domain.test", "-b", "local-k8s"},
				"",
				"testdata/one-ip.txt",
				"",
				"",
				"aliases block 'local-k8s' was not found",
			},
			false,
		},
		// {
		// 	"default two blocks",
		// 	args{
		// 		[]string{},
		// 		"testdata/two-sys-blocks.txt",
		// 		"testdata/list/default_two-sys-blocks.txt",
		// 	},
		// 	true,
		// },
		// // short
		// {
		// 	"plain empty",
		// 	args{
		// 		[]string{"-o", "short"},
		// 		"testdata/empty.txt",
		// 		"testdata/list/short_empty.txt",
		// 	},
		// 	true,
		// },
		// {
		// 	"plain one line",
		// 	args{
		// 		[]string{"-o", "short"},
		// 		"testdata/one-ip.txt",
		// 		"testdata/list/short_one-ip.txt",
		// 	},
		// 	true,
		// },
		// {
		// 	"plain two blocks",
		// 	args{
		// 		[]string{"-o", "short"},
		// 		"testdata/two-sys-blocks.txt",
		// 		"testdata/list/short_two-sys-blocks.txt",
		// 	},
		// 	true,
		// },
		// // wide
		// {
		// 	"wide empty",
		// 	args{
		// 		[]string{"-o", "wide"},
		// 		"testdata/empty.txt",
		// 		"testdata/list/wide_empty.txt",
		// 	},
		// 	true,
		// },
		// {
		// 	"wide one line",
		// 	args{
		// 		[]string{"-o", "wide"},
		// 		"testdata/one-ip.txt",
		// 		"testdata/list/wide_one-ip.txt",
		// 	},
		// 	true,
		// },
		// {
		// 	"wide two blocks",
		// 	args{
		// 		[]string{"-o", "wide"},
		// 		"testdata/two-sys-blocks.txt",
		// 		"testdata/list/wide_two-sys-blocks.txt",
		// 	},
		// 	true,
		// },
		// // plain
		// {
		// 	"plain empty",
		// 	args{
		// 		[]string{"-o", "plain"},
		// 		"testdata/empty.txt",
		// 		"testdata/list/plain_empty.txt",
		// 	},
		// 	true,
		// },
		// {
		// 	"plain one line",
		// 	args{
		// 		[]string{"-o", "plain"},
		// 		"testdata/one-ip.txt",
		// 		"testdata/list/plain_one-ip.txt",
		// 	},
		// 	true,
		// },
		// {
		// 	"plain two blocks",
		// 	args{
		// 		[]string{"-o", "plain"},
		// 		"testdata/two-sys-blocks.txt",
		// 		"testdata/list/plain_two-sys-blocks.txt",
		// 	},
		// 	true,
		// },
		// // json set of tests
		// {
		// 	"json empty",
		// 	args{
		// 		[]string{"-o", "json"},
		// 		"testdata/empty.txt",
		// 		"testdata/list/json_empty.txt",
		// 	},
		// 	true,
		// },
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
				tstp := testtools.RunHelperProcess("TestAliasAddCommand", tt.name)
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
			expectRes := string(expectData)
			expectOut := tt.args.stdout

			ctx := common.WithCustomFilesystem(context.Background(), fs)
			in := strings.NewReader(tt.args.stdin)
			out := &strings.Builder{}

			cmd := NewCmdAliasAdd()
			cmd.SilenceErrors = true
			cmd.SetArgs(tt.args.args)
			cmd.SetIn(in)
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

			fr, _ := afero.ReadFile(fs, fn)
			assert.Equal(t, expectRes, string(fr))
		})
	}
}
