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
	t.Run("unrecognized lines", func(t *testing.T) {
		content := `unrecognized line 1
# comment
unrecognized line 2`
		syndoc, _ := syntax.Read(strings.NewReader(content))
		doc := parse(syndoc)

		assert.Equal(t, syndoc, doc.originalDocument)
		assert.Equal(t, 3, doc.BlocksCount())
		assert.Equal(t, Unknown, doc.Blocks()[0].Type())
		b0 := doc.Blocks()[0].(*UnrecognizedBlock)
		assert.Equal(t, 1, len(b0.BodyElements()))
		assert.Equal(t, Comments, doc.Blocks()[1].Type())
		b1 := doc.Blocks()[1].(*CommentsBlock)
		assert.Equal(t, 1, len(b1.origComments))
		assert.Equal(t, Unknown, doc.Blocks()[2].Type())
		b2 := doc.Blocks()[2].(*UnrecognizedBlock)
		assert.Equal(t, 1, len(b2.BodyElements()))
	})
	t.Run("empty ip block", func(t *testing.T) {
		content := `# custom ips
# <<placeholder>>`
		syndoc, _ := syntax.Read(strings.NewReader(content))
		doc := parse(syndoc)

		assert.Equal(t, syndoc, doc.originalDocument)
		assert.Equal(t, 1, doc.BlocksCount())
		assert.Equal(t, IPList, doc.Blocks()[0].Type())
		b0 := doc.Blocks()[0].(*IPAliasesBlock)
		assert.Equal(t, 1, len(b0.origHeader))
		assert.Equal(t, 1, b0.Id())
		assert.Equal(t, "", b0.Name())
		assert.Equal(t, "custom ips", b0.Note())
		assert.Equal(t, 0, len(b0.Entries()))
	})
}
