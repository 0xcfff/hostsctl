package syntax

type Document struct {
	elements []Element
}

func (doc *Document) Elements() []Element {
	elements := make([]Element, len(doc.elements))
	copy(elements, doc.elements)
	return elements
}
