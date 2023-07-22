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

func TestAliasDeleteCommand(t *testing.T) {
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
			"four-blocks - by ip",
			args{
				[]string{"192.168.100.51"},
				"",
				"testdata/four-blocks.txt",
				"testdata/delete/four-blocks_one-alias_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"four-blocks - by ip - in block name",
			args{
				[]string{"192.168.100.51", "-b", "pet-prj2"},
				"",
				"testdata/four-blocks.txt",
				"testdata/delete/four-blocks_one-alias_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"four-blocks - by ip - in block id",
			args{
				[]string{"192.168.100.51", "-b", "4"},
				"",
				"testdata/four-blocks.txt",
				"testdata/delete/four-blocks_one-alias_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"four-blocks - by alias",
			args{
				[]string{"users.example.com"},
				"",
				"testdata/four-blocks.txt",
				"testdata/delete/four-blocks_one-alias_result.txt",
				"",
				"",
			},
			true,
		},
		{
			"four-blocks - by alias from multialias line",
			args{
				[]string{"awards.example.com"},
				"",
				"testdata/four-blocks.txt",
				"testdata/delete/four-blocks_one-alias-multiline_result.txt",
				"",
				"",
			},
			true,
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
				tstp := testtools.RunHelperProcess("TestAliasDeleteCommand", tt.name)
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

			cmd := NewCmdAliasDelete()
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
