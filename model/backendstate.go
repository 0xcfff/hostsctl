package model

import "time"

type BackendState struct {
	LastUpdated time.Time
	Sources     []SourceState
}

func (state *BackendState) AppendSource(source SourceState) {
	state.Sources = append(state.Sources, source)
}

func (state *BackendState) SourcesCount() int {
	return len(state.Sources)
}
