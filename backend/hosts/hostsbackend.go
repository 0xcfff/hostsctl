package hosts

import (
	"crypto/sha1"
	"fmt"
	"runtime"

	"github.com/0xcfff/dnssync/model"
	"github.com/spf13/afero"
)

const (
	PathHostsFile = "/etc/hosts"
)

type hostsBackend struct {
	etcHostsPath string
	fs           afero.Fs
}

func (backend *hostsBackend) ReadState() (*model.BackendState, error) {

	file, err := backend.fs.Open(backend.etcHostsPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %w", err)
	}

	defer file.Close()

	stats, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("Error reading stats %w", err)
	}

	hostsContent, err := ParseHostsFile(file, Strict)
	if err != nil {
		return nil, err
	}

	state := convertToBackendFormat(hostsContent)
	state.LastUpdated = stats.ModTime()

	return state, nil
}

func (backend *hostsBackend) UpdateState(changeSet model.BackendStateChangeSet) (*model.BackendState, error) {
	return &model.BackendState{}, nil
}

func DefaultBackend() model.Backend {
	return NewBackend("", nil)
}

func NewBackend(hostsFilePath string, fs afero.Fs) model.Backend {

	backend := hostsBackend{}

	osName := runtime.GOOS

	if osName == "windows" {
		panic("not implemented")

	} else {
		if hostsFilePath == "" {
			backend.etcHostsPath = PathHostsFile
		} else {
			backend.etcHostsPath = hostsFilePath
		}

		backend.fs = fs
		if backend.fs == nil {
			backend.fs = afero.NewOsFs()
		}
	}

	return &backend
}

func convertToBackendFormat(parsedContent *HostsFileContent) *model.BackendState {
	state := model.BackendState{
		ContentHash: parsedContent.ContentHash,
	}

	for _, sourceSync := range parsedContent.SyncBlocks {
		props := make(map[string]string)
		for _, p := range sourceSync.InlineProps {
			props[p.Name] = p.Value
		}
		sourceConfig := model.NewSourceConfig(props)
		source := model.SourceState{
			Config: sourceConfig,
			Data: &sourceStateWrapper{
				sourceBlock: sourceSync,
			},
		}
		state.Sources = append(state.Sources, source)
	}
	return &state
}

type sourceStateWrapper struct {
	sourceBlock *SyncBlock
}

func (w *sourceStateWrapper) DataHash() []byte {
	hasher := sha1.New()

	for _, ip := range w.sourceBlock.Data.IPRecords {
		hasher.Write([]byte(ip.IP))
		for _, al := range ip.Aliases {
			hasher.Write([]byte(al))
		}
	}

	return hasher.Sum(nil)
}
