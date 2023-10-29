package syntax

import (
	"bytes"
	"strings"
	"testing"

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
			doc, _ := parse(r)
			w := &strings.Builder{}

			// act
			err := format(w, doc, FmtDefault)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, content, w.String())
		})
	}
}
