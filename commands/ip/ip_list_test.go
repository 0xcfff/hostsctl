package ip

import (
	"context"
	"strings"
	"testing"

	"github.com/0xcfff/hostsctl/commands/common"
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

		ctx := common.WithCustomFilesystem(context.Background(), fs)
		out := &strings.Builder{}

		cmd := NewCmdIpList()
		cmd.SetArgs(make([]string, 0))
		cmd.SetOutput(out)

		expectOut := "GRP  SYS  IP  ALIAS\n"

		// TODO: Refactor tests

		c, err := cmd.ExecuteContextC(ctx)

		assert.NoError(t, err, "command should succeed")
		assert.Same(t, cmd, c)

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

		s := out.String()
		assert.Equal(t, expectOut, s)
	})
	t.Run("list one line", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		f, _ := fs.Create(hosts.EtcHosts.Path())
		f.WriteString("127.0.0.1   localhost")
		f.Close()

		ctx := common.WithCustomFilesystem(context.Background(), fs)
		out := &strings.Builder{}

		cmd := NewCmdIpList()
		cmd.SetArgs(make([]string, 0))
		cmd.SetOutput(out)

		expectOut :=
			`GRP  SYS  IP         ALIAS
[+]  +    127.0.0.1  localhost
`

		c, err := cmd.ExecuteContextC(ctx)

		assert.NoError(t, err, "command should succeed")
		assert.Same(t, cmd, c)

		s := out.String()
		assert.Equal(t, expectOut, s)
	})

	t.Run("list empty json", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		f, _ := fs.Create(hosts.EtcHosts.Path())
		f.WriteString("")
		f.Close()

		ctx := common.WithCustomFilesystem(context.Background(), fs)
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
