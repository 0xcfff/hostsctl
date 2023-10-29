package common

type OutputFormat int

const (
	FmtText OutputFormat = iota
	FmtJson OutputFormat = iota
	FmtYaml OutputFormat = iota
	FmtNext OutputFormat = iota
)

const (
	TfmtText = "text"
	TfmtJson = "json"
	TfmtYaml = "yaml"
)
