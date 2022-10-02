package hosts

import (
	"crypto/sha1"
	"fmt"

	"github.com/0xcfff/dnspipe/model"
)

type hostsBackend struct {
	src *HostsFileSource
}

func (backend *hostsBackend) ReadState() (*model.BackendState, error) {

	file, err := backend.src.openRead()
	if err != nil {
		return nil, fmt.Errorf("Error opening file %w", err)
	}

	defer file.Close()

	stats, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("Error reading stats %w", err)
	}

	hostsContent, err := ParseHostsFileWithSources(file, Strict)
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

func NewBackend(src *HostsFileSource) model.Backend {

	if src == nil {
		src = NewHostsFileSource("", nil)
	}

	backend := hostsBackend{
		src: src,
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
