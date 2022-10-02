package hosts

import (
	"strings"
	"testing"

	"github.com/0xcfff/dnspipe/model"
	"github.com/stretchr/testify/assert"
)

func TestHostsFile_AppendSource(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		orig := ``
		expected := `# @sync kubectl://microk8s
`
		source := model.NewSourceConfig(map[string]string{
			"source": "kubectl://microk8s",
		})

		// act
		f := NewHostsFile([]byte(orig))
		err := f.AppendSource(source)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, string(f.content))
	})
	t.Run("empty many props", func(t *testing.T) {
		orig := ``
		expected := `# @sync kubectl://microk8s
# @props map=nodeip objects=ingress,endpoints
`
		source := model.NewSourceConfig(map[string]string{
			"source":  "kubectl://microk8s",
			"objects": "ingress,endpoints",
			"map":     "nodeip",
		})

		// act
		f := NewHostsFile([]byte(orig))
		err := f.AppendSource(source)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, string(f.content))
	})

	t.Run("content exists", func(t *testing.T) {
		orig := `127.0.0.1	localhost
127.0.1.1	box01 box02`
		expected := `127.0.0.1	localhost
127.0.1.1	box01 box02

# @sync kubectl://microk8s
# @props map=nodeip objects=ingress,endpoints
`
		source := model.NewSourceConfig(map[string]string{
			"source":  "kubectl://microk8s",
			"objects": "ingress,endpoints",
			"map":     "nodeip",
		})

		// act
		f := NewHostsFile([]byte(orig))
		err := f.AppendSource(source)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, string(f.content))
	})
}

func Test_writeSourceConfig(t *testing.T) {
	type args struct {
		b      *strings.Builder
		source model.SourceConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeSourceConfig(tt.args.b, tt.args.source); (err != nil) != tt.wantErr {
				t.Errorf("writeSourceConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_osDefaultNewLine(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := osDefaultNewLine(); got != tt.want {
				t.Errorf("osDefaultNewLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
