package model

type SourceStateAdd struct {
	NewState SourceState
}
type SourceStateUpdate struct {
	OriginalState SourceState
	NewState      SourceState
}
type SourceStateRemove struct {
	OriginalState SourceState
}

type BackendStateChangeSet struct {
	AddSources    []SourceStateAdd
	UpdateSources []SourceStateUpdate
	RemoveSources []SourceStateRemove
}

func NewBackendStateChangeSet() *BackendStateChangeSet {
	return &BackendStateChangeSet{
		AddSources:    make([]SourceStateAdd, 0),
		UpdateSources: make([]SourceStateUpdate, 0),
		RemoveSources: make([]SourceStateRemove, 0),
	}
}

func (change *BackendStateChangeSet) AppendAdd(newState SourceState) {
	change.AddSources = append(change.AddSources, SourceStateAdd{newState})
}

func (change *BackendStateChangeSet) AppendUpdate(originalState, newState SourceState) {
	change.UpdateSources = append(change.UpdateSources, SourceStateUpdate{originalState, newState})
}

func (change *BackendStateChangeSet) AppendRemove(originalState SourceState) {
	change.RemoveSources = append(change.RemoveSources, SourceStateRemove{originalState})
}
