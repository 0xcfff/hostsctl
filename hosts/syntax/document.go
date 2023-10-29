package syntax

type Document struct {
	elements []Element
}

func (doc *Document) Elements() []Element {
	elements := make([]Element, len(doc.elements))
	copy(elements, doc.elements)
	return elements
}

// Creates new documen based on the passed data
func NewDocument(elements []Element) *Document {
	content := elements
	if elements == nil {
		content = make([]Element, 0)
	}
	return &Document{
		elements: content,
	}
}
