package model

type SourceData interface {
	// method calculates hash based on the source data
	// if data changes, the hash should change as well
	DataHash() []byte
}
