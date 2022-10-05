package hosts

type IPModel struct {
	IP      string   `json:"ip"`
	Aliases []string `json:"aliases"`
	Comment string   `json:"comment,omitempty"`
	Source  string   `json:"source,omitempty"`
}
