package syntax

type ElementType int

const (
	Unknown ElementType = iota
	IPMapping
	Comment
	Empty
)

// Shared syntax element interface
type Element interface {
	Type() ElementType
	HasOriginalText() bool
	OriginalLineIndex() int
	OriginalLineText() string
}

// Base syntax element functionality
type elementBase struct {
	originalLineText  *string
	originalLineIndex int
}

// Returns true if the element has original text
func (el *elementBase) HasOriginalText() bool {
	return el.originalLineText != nil
}

func (el *elementBase) OriginalLineText() string {
	return *el.originalLineText
}

func (el *elementBase) OriginalLineIndex() int {
	return el.originalLineIndex
}

// Represents a line of comments
type CommentLine struct {
	elementBase
	commentText string
}

func (*CommentLine) Type() ElementType {
	return Comment
}

func (el *CommentLine) CommentText() string {
	return el.commentText
}

type UnrecognizedLine struct {
	elementBase
}

func (*UnrecognizedLine) Type() ElementType {
	return Unknown
}

// Represents an empty line
type EmptyLine struct {
	elementBase
}

func (*EmptyLine) Type() ElementType {
	return Empty
}

// Represents a line of IP to domain name mapping
type IPMappingLine struct {
	elementBase
	ip          string
	domainNames []string
	commentText string
}

func (*IPMappingLine) Type() ElementType {
	return IPMapping
}

func (el *IPMappingLine) IPAddress() string {
	return el.ip
}

func (el *IPMappingLine) DomainNames() []string {
	return el.domainNames
}

func (el *IPMappingLine) CommentText() string {
	return el.commentText
}
