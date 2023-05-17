package dom

import (
	"strings"
	"testing"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		content := ""
		syndoc, _ := syntax.Read(strings.NewReader(content))
		doc := parse(syndoc)

		assert.Equal(t, syndoc, doc.originalDocument)
		assert.Equal(t, 0, doc.BlocksCount())
	})
	t.Run("ip only", func(t *testing.T) {
		content := "127.0.0.1 localhost"
		syndoc, _ := syntax.Read(strings.NewReader(content))
		doc := parse(syndoc)

		assert.Equal(t, syndoc, doc.originalDocument)
		assert.Equal(t, 1, doc.BlocksCount())
		assert.Equal(t, IPList, doc.Blocks()[0].Type())
	})
	t.Run("named ip block only", func(t *testing.T) {
		content := `# system ips
127.0.0.1 localhost`
		syndoc, _ := syntax.Read(strings.NewReader(content))
		doc := parse(syndoc)

		assert.Equal(t, syndoc, doc.originalDocument)
		assert.Equal(t, 1, doc.BlocksCount())
		assert.Equal(t, IPList, doc.Blocks()[0].Type())
		b0 := doc.Blocks()[0].(*IPAliasesBlock)
		assert.Equal(t, 1, len(b0.origHeader))
		assert.Equal(t, "system ips", b0.Note())
		assert.Equal(t, 1, len(b0.Entries()))
		assert.Equal(t, 1, b0.Id())
		assert.Equal(t, "", b0.Name())
		assert.Equal(t, "system ips", b0.Note())
	})
	t.Run("properly named ip block only", func(t *testing.T) {
		content := `# [101] proj-01 - system ips
127.0.0.1 localhost`
		syndoc, _ := syntax.Read(strings.NewReader(content))
		doc := parse(syndoc)

		assert.Equal(t, syndoc, doc.originalDocument)
		assert.Equal(t, 1, doc.BlocksCount())
		assert.Equal(t, IPList, doc.Blocks()[0].Type())
		b0 := doc.Blocks()[0].(*IPAliasesBlock)
		assert.Equal(t, 1, len(b0.origHeader))
		assert.Equal(t, "[101] proj-01 - system ips", b0.origHeader[0].CommentText())
		assert.Equal(t, 1, len(b0.Entries()))
		assert.Equal(t, 101, b0.Id())
		assert.Equal(t, "proj-01", b0.Name())
		assert.Equal(t, "system ips", b0.Note())
	})
}
