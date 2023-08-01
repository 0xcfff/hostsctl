package syntax

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

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
	OriginalLineIndex() int
	HasPreformattedText() bool
	PreformattedLineText() string
	formatLine() string
}

// Base syntax element functionality
type elementBase struct {
	preformattedLineText *string
	originalLineIndex    int
}

// Returns true if the element has original text
func (el *elementBase) HasPreformattedText() bool {
	return el.preformattedLineText != nil
}

func (el *elementBase) PreformattedLineText() string {
	return *el.preformattedLineText
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

func (el *CommentLine) formatLine() string {
	return fmt.Sprintf("# %v", el.commentText)
}

func NewCommentsLine(commentText string) *CommentLine {
	return &CommentLine{
		commentText: commentText,
	}
}

type UnrecognizedLine struct {
	elementBase
}

func (*UnrecognizedLine) Type() ElementType {
	return Unknown
}

func (*UnrecognizedLine) formatLine() string {
	panic("automated formatting is not supported")
}

// Represents an empty line
type EmptyLine struct {
	elementBase
}

func (*EmptyLine) Type() ElementType {
	return Empty
}

func (*EmptyLine) formatLine() string {
	return ""
}

func NewEmptyLine() *EmptyLine {
	return &EmptyLine{}
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

func (el *IPMappingLine) formatLine() string {
	b := strings.Builder{}
	b.WriteString(el.ip)
	b.WriteString(" ")

	idx := 0

	for _, ip := range el.domainNames {
		if idx > 0 {
			b.WriteString(", ")
		}
		b.WriteString(ip)
	}
	return b.String()
}

func NewIPMappingLine(ip string, domainNames []string, comment string) *IPMappingLine {
	return &IPMappingLine{
		ip:          ip,
		domainNames: slices.Clone(domainNames),
		commentText: comment,
	}
}
