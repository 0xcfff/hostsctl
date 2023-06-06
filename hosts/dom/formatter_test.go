package dom

import (
	"bytes"
	"strings"
	"testing"

	"github.com/0xcfff/hostsctl/hosts/syntax"
	"github.com/stretchr/testify/assert"
)

func Test_format_sameAsOriginal(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"empty file", ""},
		{"newline only file", "\n"},
		{"newline and whitespaces file", "  \n   "},
		{"ipv4 only file", "127.0.0.1 localhost"},
		{"ipv4 only with a comment file", "127.0.0.1 localhost # my own IP"},
		{"ipv6 only file", "fe00::0 ip6-localnet"},
		{"ipv6 only with 2 fqdn file", "::1     ip6-localhost ip6-loopback"},
		{"incorrect ipv6 file", ":t:1     ip6-localhost"},
		{"comments only", "# line 1\n# line2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			content := tt.content
			r := bytes.NewReader([]byte(content))
			sdoc, _ := syntax.Read(r)
			doc := parse(sdoc)
			w := &strings.Builder{}

			// act
			fsdoc := constructSyntax(doc)
			syntax.Write(w, fsdoc, syntax.FmtDefault)

			// assert
			assert.NotSame(t, sdoc, fsdoc)       // different document instances
			assert.Equal(t, content, w.String()) // but same content
		})
	}
}
