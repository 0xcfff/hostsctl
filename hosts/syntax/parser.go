package syntax

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/0xcfff/hostsctl/iotools"
	"github.com/0xcfff/hostsctl/iptools"
)

// Main intry point into syntax parsing process
func parse(r io.Reader) (*Document, error) {
	s := bufio.NewScanner(r)
	s.Split(iotools.LinesSplitterRespectEndNewLineFunc())
	els, err := parseLines(s)
	if err != nil {
		return nil, err
	}
	doc := Document{
		elements: els,
	}
	return &doc, nil
}

func parseLines(s *bufio.Scanner) ([]Element, error) {
	elements := make([]Element, 0)
	lineIndex := 0

	for {
		if ok := s.Scan(); !ok {
			break
		}
		rawText := s.Text()
		lineIndex += 1

		element, err := parseLine(lineIndex, rawText)
		if err != nil {
			return nil, fmt.Errorf("error parsing line %v, %w", lineIndex, err)
		}

		elements = append(elements, element)
	}
	return elements, nil
}

func parseLine(idx int, l string) (Element, error) {
	tl := strings.TrimSpace(l)
	elb := elementBase{
		originalLineIndex:    idx,
		preformattedLineText: &l,
	}

	if len(tl) == 0 {
		return &EmptyLine{
			elementBase: elb,
		}, nil
	}

	if strings.HasPrefix(tl, "#") {
		ct := strings.TrimSpace(strings.TrimLeft(tl, "#"))
		return &CommentLine{
			elementBase: elb,
			commentText: ct,
		}, nil
	}

	parts := strings.Fields(tl)
	if len(parts) > 1 && iptools.IsIP(parts[0]) {
		fqdns := make([]string, 0)
		hasComment := false
		for _, s := range parts[1:] {
			if strings.HasPrefix(s, "#") {
				hasComment = true
				break
			}
			fqdns = append(fqdns, s)
		}
		if len(fqdns) > 0 {
			comment := ""
			if hasComment {
				idx := strings.Index(tl, "#")
				comment = strings.TrimSpace(tl[idx+1:])
			}
			return &IPMappingLine{
				elementBase: elb,
				ip:          parts[0],
				domainNames: fqdns,
				commentText: comment,
			}, nil
		}
	}

	return &UnrecognizedLine{
		elementBase: elb,
	}, nil
}
