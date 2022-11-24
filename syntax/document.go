package syntax

type Document interface {
	Elements() []Element
}

type document struct {
	elements []Element
}

func (doc *document) Elements() []Element {
	return doc.elements
}
