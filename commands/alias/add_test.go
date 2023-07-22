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
		// add to specific block
		{
			"add - 3rd block + no block specified",
			args{
				[]string{"192.168.100.100", "local-service"},
				"",
				"testdata/four-blocks.txt",
				"testdata/add/four-blocks_one-alias_no-block_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"add - 3rd block + by id",
			args{
				[]string{"192.168.100.100", "local-service", "-b", "3"},
				"",
				"testdata/four-blocks.txt",
				"testdata/add/four-blocks_one-alias_3rd-block_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"add - 3rd block + by name",
			args{
				[]string{"192.168.100.100", "local-service", "-b", "pet-prj1"},
				"",
				"testdata/four-blocks.txt",
				"testdata/add/four-blocks_one-alias_3rd-block_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"add - 5th block + by id",
			args{
				[]string{"192.168.100.100", "local-service", "-b", "5", "--force"},
				"",
				"testdata/four-blocks.txt",
				"testdata/add/four-blocks_one-alias-by-id_5th-block_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"add - 5th block + by name",
			args{
				[]string{"192.168.100.100", "local-service", "-b", "prj-pet009", "--force"},
				"",
				"testdata/four-blocks.txt",
				"testdata/add/four-blocks_one-alias-by-name_5th-block_result.txt",
				"",
				"",
			},
			true,
		},

		// errors cases
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
