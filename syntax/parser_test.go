package syntax

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		content := []byte("")
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 0, len(doc.Elements()))
	})

	t.Run("newline only file", func(t *testing.T) {
		content := []byte("\n")
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 2, len(doc.Elements()))
		assert.Equal(t, Empty, doc.Elements()[0].Type())
		assert.Equal(t, Empty, doc.Elements()[1].Type())
	})

	t.Run("newline and whitespaces file", func(t *testing.T) {
		content := []byte("  \n   ")
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 2, len(doc.Elements()))

		assert.Equal(t, Empty, doc.Elements()[0].Type())
		el0 := doc.Elements()[0].(*EmptyLine)
		assert.Equal(t, 1, el0.originalLineIndex)
		assert.Equal(t, "  ", *el0.originalLineText)

		assert.Equal(t, Empty, doc.Elements()[1].Type())
		el1 := doc.Elements()[1].(*EmptyLine)
		assert.Equal(t, 2, el1.originalLineIndex)
		assert.Equal(t, "   ", *el1.originalLineText)
	})

	t.Run("ipv4 only file", func(t *testing.T) {
		content := []byte("127.0.0.1 localhost")
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 1, len(doc.Elements()))
		assert.Equal(t, IPMapping, doc.Elements()[0].Type())
		el0 := doc.Elements()[0].(*IPMappingLine)
		assert.Equal(t, "127.0.0.1", el0.IPAddress())
		assert.Equal(t, 1, len(el0.DomainNames()))
		assert.Equal(t, "localhost", el0.DomainNames()[0])
	})

	t.Run("ipv6 only file", func(t *testing.T) {
		content := []byte("fe00::0 ip6-localnet")
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 1, len(doc.Elements()))
		assert.Equal(t, IPMapping, doc.Elements()[0].Type())
		el0 := doc.Elements()[0].(*IPMappingLine)
		assert.Equal(t, "fe00::0", el0.IPAddress())
		assert.Equal(t, 1, len(el0.DomainNames()))
		assert.Equal(t, "ip6-localnet", el0.DomainNames()[0])
	})

	t.Run("ipv6 only with 2 fqdn file", func(t *testing.T) {
		content := []byte("::1     ip6-localhost ip6-loopback")
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 1, len(doc.Elements()))
		assert.Equal(t, IPMapping, doc.Elements()[0].Type())
		el0 := doc.Elements()[0].(*IPMappingLine)
		assert.Equal(t, "::1", el0.IPAddress())
		assert.Equal(t, 2, len(el0.DomainNames()))
		assert.Equal(t, "ip6-localhost", el0.DomainNames()[0])
		assert.Equal(t, "ip6-loopback", el0.DomainNames()[1])
	})

	t.Run("incorrect ipv4 file", func(t *testing.T) {
		content := []byte("12a.0.0.1 localhost")
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 1, len(doc.Elements()))
		assert.Equal(t, Unknown, doc.Elements()[0].Type())
		el0 := doc.Elements()[0].(*UnrecognizedLine)
		assert.Equal(t, "12a.0.0.1 localhost", *el0.originalLineText)
	})

	t.Run("incorrect ipv6 file", func(t *testing.T) {
		content := []byte(":t:1     ip6-localhost")
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 1, len(doc.Elements()))
		assert.Equal(t, Unknown, doc.Elements()[0].Type())
		el0 := doc.Elements()[0].(*UnrecognizedLine)
		assert.Equal(t, ":t:1     ip6-localhost", *el0.originalLineText)
	})

	t.Run("all elements", func(t *testing.T) {
		content := []byte(`# ipv4 mappings  
 127.0.0.1    localhost  
 
  # ipv6 mappings
  ::1     ip6-localhost ip6-loopback
  fe00::0 ip6-localnet
`)
		reader := bytes.NewReader(content)

		doc, err := Parse(reader)

		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, 7, len(doc.Elements()))

		assert.Equal(t, Comment, doc.Elements()[0].Type())
		el0 := doc.Elements()[0].(*CommentLine)
		assert.Equal(t, 1, el0.originalLineIndex)
		assert.Equal(t, "# ipv4 mappings  ", *el0.originalLineText)
		assert.Equal(t, "ipv4 mappings", el0.CommentText())

		assert.Equal(t, IPMapping, doc.Elements()[1].Type())
		el1 := doc.Elements()[1].(*IPMappingLine)
		assert.Equal(t, 2, el1.originalLineIndex)
		assert.Equal(t, " 127.0.0.1    localhost  ", *el1.originalLineText)
		assert.Equal(t, "127.0.0.1", el1.IPAddress())
		assert.Equal(t, 1, len(el1.DomainNames()))
		assert.Equal(t, "localhost", el1.DomainNames()[0])

		assert.Equal(t, Empty, doc.Elements()[2].Type())
		el2 := doc.Elements()[2].(*EmptyLine)
		assert.Equal(t, 3, el2.originalLineIndex)
		assert.Equal(t, " ", *el2.originalLineText)
	})

}
