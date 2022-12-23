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

// TODO: is this overcomplication needed? Check ip_list.go implementation. Probably implement couple command and then return to this question

// Returns ordered slice of globally known output format names
func OutputFormatNames() []string {
	return []string{TfmtText, TfmtJson, TfmtJson}
}

// Returns ordered slice of globally known output format values
func OutputFormatValues() []OutputFormat {
	return []OutputFormat{FmtText, FmtJson, FmtYaml}
}


// func NewOutputFormatsMap(outputFormatNames []string, outputFormatValues) map[string]OutputFormat {
// 	for it
// }