package ip

import (
	"context"
	"strings"
	"testing"

	"github.com/0xcfff/hostsctl/hosts"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdIpList(t *testing.T) {
	t.Run("list empty", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		f, _ := fs.Create(hosts.EtcHosts.Path())
		f.WriteString("")
		f.Close()

		ctx := context.WithValue(context.Background(), ctxCustomFs, fs)
		out := &strings.Builder{}

		cmd := NewCmdIpList()
		cmd.SetArgs(make([]string, 0))
		cmd.SetOutput(out)

		c, err := cmd.ExecuteContextC(ctx)

		assert.NoError(t, err, "command should succeed")
		assert.Same(t, cmd, c)

		assert.Fail(t, "Implement correct result check and command support of custom FS")

	})
	t.Run("list empty json", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		f, _ := fs.Create(hosts.EtcHosts.Path())
		f.WriteString("")
		f.Close()

		ctx := context.WithValue(context.Background(), ctxCustomFs, fs)
		out := &strings.Builder{}

		cmd := NewCmdIpList()
		cmd.SetArgs([]string{"-o", "json"})
		cmd.SetOutput(out)

		c, err := cmd.ExecuteContextC(ctx)

		assert.NoError(t, err, "command should succeed")
		assert.Same(t, cmd, c)

		assert.Fail(t, "Implement correct result check and command support of custom FS")

	})
}
