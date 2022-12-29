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
		b0 := doc.Blocks()[0].(*IPListBlock)
		assert.Equal(t, 1, len(b0.HeaderCommentLines()))
		assert.Equal(t, "system ips", b0.HeaderCommentLines()[0])
		assert.Equal(t, 1, len(b0.BodyElements()))
	})
}
