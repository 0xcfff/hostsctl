package model

type Backend interface {
	ReadState() (*BackendState, error)
	UpdateState(changeSet BackendStateChangeSet) (*BackendState, error)
}
