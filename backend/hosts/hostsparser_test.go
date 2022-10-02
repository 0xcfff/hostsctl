package hosts

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HostsParser_DSLCanParseSamples(t *testing.T) {
	basePath := "testdata/parsesamples"
	testCases, _ := os.ReadDir(basePath)

	for _, tc := range testCases {
		if !tc.IsDir() {
			tclocal := tc
			t.Run(fmt.Sprintf("parsing: %s", tc.Name()), func(t *testing.T) {
				f, err := os.Open(path.Join(basePath, tclocal.Name()))
				if err != nil {
					t.Errorf("Error readin file %s %s", tclocal.Name(), err)
				}

				defer f.Close()

				// data, err := HostsParser.Parse(tclocal.Name(), f)
				data, err := ParseHostsFileWithSources(f, Strict)

				if err != nil {
					t.Errorf("Error parsing '%s', %s ", tclocal.Name(), err)
				}

				if data == nil {
					t.Errorf("Error parsing '%s' gave empty result", tclocal.Name())
				}

			})
		}
	}
}

func Test_ParseHostsFile_ParsesDataCorrectly(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		// arrange
		rawContent := ``

		// act
		r := strings.NewReader(rawContent)
		data, err := ParseHostsFileWithSources(r, Strict)

		// assert
		assert.NoError(t, err, "Error parsing content")
		assert.NotNil(t, data)
		assert.Nil(t, data.IPRecords)
		assert.Nil(t, data.SyncBlocks)
	})
	t.Run("IPv4Only", func(t *testing.T) {
		// arrange
		rawContent := `
		127.0.0.1	localhost
		127.0.1.1	box01 box02 # local dev`

		expected := &HostsFileContent{
			IPRecords: []*IPRecord{
				{
					Pos:     Position{Line: 2},
					IP:      "127.0.0.1",
					Aliases: []string{"localhost"},
				},
				{
					Pos:     Position{Line: 3},
					IP:      "127.0.1.1",
					Aliases: []string{"box01", "box02"},
					Notes:   "local dev",
				},
			},
		}

		// act
		r := strings.NewReader(rawContent)
		actual, err := ParseHostsFileWithSources(r, Strict)

		// hack
		expected.ContentHash = actual.ContentHash

		// assert
		assert.NoError(t, err, "parsing should be successful")
		assert.NotNil(t, actual, "parse should return if not error")
		assert.Equal(t, expected, actual)
		assert.NotEmpty(t, actual.ContentHash)
	})
	t.Run("IPv6Only", func(t *testing.T) {
		// arrange
		rawContent := `
		::1     ip6-localhost ip6-loopback
		fe00::0 ip6-localnet
		2001:0db8:85a3:0000:0000:8a2e:0370:7334 ipv6.domain.com`

		expected := &HostsFileContent{
			IPRecords: []*IPRecord{
				{
					Pos:     Position{Line: 2},
					IP:      "::1",
					Aliases: []string{"ip6-localhost", "ip6-loopback"},
				},
				{
					Pos:     Position{Line: 3},
					IP:      "fe00::0",
					Aliases: []string{"ip6-localnet"},
				},
				{
					Pos:     Position{Line: 4},
					IP:      "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
					Aliases: []string{"ipv6.domain.com"},
				},
			},
		}

		// act
		r := strings.NewReader(rawContent)
		actual, err := ParseHostsFileWithSources(r, Strict)

		// hack
		expected.ContentHash = actual.ContentHash

		// assert
		assert.NoError(t, err, "parsing should be successful")
		assert.NotNil(t, actual, "parse should return if not error")
		assert.Equal(t, expected, actual)
		assert.NotEmpty(t, actual.ContentHash)
	})
	t.Run("sync minimal", func(t *testing.T) {
		// arrange
		rawContent := `# @sync http://192.168.64.1:9053/dns?src=ingress,lb`

		expected := &HostsFileContent{
			SyncBlocks: []*SyncBlock{
				{
					Pos: Position{Line: 1},
					InlineProps: []*InlineProperty{
						{
							Pos:   Position{Line: 1},
							Name:  "source",
							Value: "http://192.168.64.1:9053/dns?src=ingress,lb",
						},
					},
					PosEndHeader: Position{Line: 1},
				},
			},
		}

		// act
		r := strings.NewReader(rawContent)
		actual, err := ParseHostsFileWithSources(r, Strict)

		// hack
		expected.ContentHash = actual.ContentHash

		// assert
		assert.NoError(t, err, "parsing should be successful")
		assert.NotNil(t, actual, "parse should return if not error")
		assert.Equal(t, expected, actual)
		assert.NotEmpty(t, actual.ContentHash)
	})
	t.Run("sync onelineprop", func(t *testing.T) {
		// arrange
		rawContent := `# @sync source=http://192.168.100.2/dns-records?domain=.sync-ops.com&domain=.sync-svc.com, unsafe=true, interval=auto`

		expected := &HostsFileContent{
			SyncBlocks: []*SyncBlock{
				{
					Pos: Position{Line: 1},
					InlineProps: []*InlineProperty{
						{
							Pos:   Position{Line: 1},
							Name:  "source",
							Value: "http://192.168.100.2/dns-records?domain=.sync-ops.com&domain=.sync-svc.com",
						},
						{
							Pos:   Position{Line: 1},
							Name:  "unsafe",
							Value: "true",
						},
						{
							Pos:   Position{Line: 1},
							Name:  "interval",
							Value: "auto",
						},
					},
					PosEndHeader: Position{Line: 1},
				},
			},
		}

		// act
		r := strings.NewReader(rawContent)
		actual, err := ParseHostsFileWithSources(r, Strict)

		// hack
		expected.ContentHash = actual.ContentHash

		// assert
		assert.NoError(t, err, "parsing should be successful")
		assert.NotNil(t, actual, "parse should return if not error")
		assert.Equal(t, expected, actual)
		assert.NotEmpty(t, actual.ContentHash)
	})
	t.Run("sync multiline", func(t *testing.T) {
		// arrange
		rawContent := `
		# @sync kubectl
		# @context microk8s
		# @interval auto`

		expected := &HostsFileContent{
			SyncBlocks: []*SyncBlock{
				{
					Pos: Position{Line: 2},
					InlineProps: []*InlineProperty{
						{
							Pos:   Position{Line: 2},
							Name:  "source",
							Value: "kubectl",
						},
						{
							Pos:   Position{Line: 3},
							Name:  "context",
							Value: "microk8s",
						},
						{
							Pos:   Position{Line: 4},
							Name:  "interval",
							Value: "auto",
						},
					},
					PosEndHeader: Position{Line: 4},
				},
			},
		}

		// act
		r := strings.NewReader(rawContent)
		actual, err := ParseHostsFileWithSources(r, Strict)

		// hack
		expected.ContentHash = actual.ContentHash

		// assert
		assert.NoError(t, err, "parsing should be successful")
		assert.NotNil(t, actual, "parse should return if not error")
		assert.Equal(t, expected, actual)
		assert.NotEmpty(t, actual.ContentHash)
	})
	t.Run("sync mixed", func(t *testing.T) {
		// arrange
		rawContent := `# @sync kubectl://microk8s
					   # @props unsafe=true, interval=auto
					   # @log verbose
					   # @kubeargs -v -d -c`

		expected := &HostsFileContent{
			SyncBlocks: []*SyncBlock{
				{
					Pos: Position{Line: 1},
					InlineProps: []*InlineProperty{
						{
							Pos:   Position{Line: 1},
							Name:  "source",
							Value: "kubectl://microk8s",
						},
						{
							Pos:   Position{Line: 2},
							Name:  "unsafe",
							Value: "true",
						},
						{
							Pos:   Position{Line: 2},
							Name:  "interval",
							Value: "auto",
						},
						{
							Pos:   Position{Line: 3},
							Name:  "log",
							Value: "verbose",
						},
						{
							Pos:   Position{Line: 4},
							Name:  "kubeargs",
							Value: "-v -d -c",
						},
					},
					PosEndHeader: Position{Line: 4},
				},
			},
		}

		// act
		r := strings.NewReader(rawContent)
		actual, err := ParseHostsFileWithSources(r, Strict)

		// hack
		expected.ContentHash = actual.ContentHash

		// assert
		assert.NoError(t, err, "parsing should be successful")
		assert.NotNil(t, actual, "parse should return if not error")
		assert.Equal(t, expected, actual)
		assert.NotEmpty(t, actual.ContentHash)
	})
	t.Run("sync full", func(t *testing.T) {
		// arrange
		rawContent := `127.0.0.1	localhost
					   127.0.1.1	box01 box02
					   ::1     		ip6-localhost ip6-loopback
					   fe00::0 		ip6-localnet

					   # synchronization blocks go below
					   # ================================
					   # @sync http://192.168.64.1/dns-records?domain=.ops.com
					   # @props unsafe=true, interval=auto

					   # @sync source=kubectl://microk8s
					   # @props unsafe=false, interval=5s
					   # @begin_sync
					   192.168.64.4 kibana.localops.com
					   192.168.64.4 grafana.localops.com
					   # @end_sync

					   127.0.3.1	box03 box04
					   `

		expected := &HostsFileContent{
			IPRecords: []*IPRecord{
				{
					Pos:     Position{Line: 1},
					IP:      "127.0.0.1",
					Aliases: []string{"localhost"},
				},
				{
					Pos:     Position{Line: 2},
					IP:      "127.0.1.1",
					Aliases: []string{"box01", "box02"},
				},
				{
					Pos:     Position{Line: 3},
					IP:      "::1",
					Aliases: []string{"ip6-localhost", "ip6-loopback"},
				},
				{
					Pos:     Position{Line: 4},
					IP:      "fe00::0",
					Aliases: []string{"ip6-localnet"},
				},
				{
					Pos:     Position{Line: 18},
					IP:      "127.0.3.1",
					Aliases: []string{"box03", "box04"},
				},
			},
			SyncBlocks: []*SyncBlock{
				{
					Pos: Position{Line: 8},
					InlineProps: []*InlineProperty{
						{
							Pos:   Position{Line: 8},
							Name:  "source",
							Value: "http://192.168.64.1/dns-records?domain=.ops.com",
						},
						{
							Pos:   Position{Line: 9},
							Name:  "unsafe",
							Value: "true",
						},
						{
							Pos:   Position{Line: 9},
							Name:  "interval",
							Value: "auto",
						},
					},
					PosEndHeader: Position{Line: 9},
				},
				{
					Pos: Position{Line: 11},
					InlineProps: []*InlineProperty{
						{
							Pos:   Position{Line: 11},
							Name:  "source",
							Value: "kubectl://microk8s",
						},
						{
							Pos:   Position{Line: 12},
							Name:  "unsafe",
							Value: "false",
						},
						{
							Pos:   Position{Line: 12},
							Name:  "interval",
							Value: "5s",
						},
					},
					PosEndHeader: Position{Line: 12},
					Data: &SyncDataBlock{
						Pos: Position{Line: 13},
						IPRecords: []*IPRecord{
							{
								Pos:     Position{Line: 14},
								IP:      "192.168.64.4",
								Aliases: []string{"kibana.localops.com"},
							},
							{
								Pos:     Position{Line: 15},
								IP:      "192.168.64.4",
								Aliases: []string{"grafana.localops.com"},
							},
						},
						PosEndData: Position{Line: 16},
					},
				},
			},
		}

		// act
		r := strings.NewReader(rawContent)
		actual, err := ParseHostsFileWithSources(r, Strict)

		// hack
		expected.ContentHash = actual.ContentHash

		// assert
		assert.NoError(t, err, "parsing should be successful")
		assert.NotNil(t, actual, "parse should return if not error")
		assert.Equal(t, expected, actual)
		assert.NotEmpty(t, actual.ContentHash)
	})
}

func Test_HostsPaster_DLSDebug(t *testing.T) {
	t.Run("debug DSL", func(t *testing.T) {
		// var reader = strings.NewReader("# @sync test")
		var reader = strings.NewReader("# @sync source=test, interval=5s, target=next")
		// var reader = strings.NewReader("192.168.1.1 alias1 alias2 # comment")
		data, err := ParseHostsFileWithSources(reader, Strict)

		if err != nil {
			t.Errorf("error parsing sample: %s", err)
		}

		fmt.Println(data)
		t.Log(data)
	})
}
