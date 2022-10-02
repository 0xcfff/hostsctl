package hosts

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/0xcfff/dnssync/model"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type HostsFile struct {
	content []byte
}

var (
	ErrNoSource     = errors.New("No source property")
	newLine         = osDefaultNewLine()
	singleLineProps = []string{"enabled", "interval"}
)

func NewHostsFileFromContent(content []byte) *HostsFile {
	return &HostsFile{
		content: content,
	}
}

func (f *HostsFile) Parse(mode ParseMode) (*HostsFileContent, error) {
	return ParseHostsFile(bytes.NewReader(f.content), mode)
}

func (f *HostsFile) AddSource(source model.SourceConfig) error {
	b := &strings.Builder{}
	s := bufio.NewScanner(bytes.NewReader(f.content))
	lastLineEmpty := true

	for {
		ok := s.Scan()
		if !ok {
			break
		}

		line := s.Text()
		lastLineEmpty = strings.TrimSpace(line) == ""

		_, err := b.WriteString(line)
		if err != nil {
			return fmt.Errorf("Error writting config %w", err)
		}
		_, err = b.WriteString(newLine)
		if err != nil {
			return fmt.Errorf("Error writting config %w", err)
		}
	}

	if !lastLineEmpty {
		b.WriteString(newLine)
	}

	err := writeSourceConfig(b, source)
	if err != nil {
		return fmt.Errorf("Error writting source config %w", err)
	}

	f.content = []byte(b.String())

	return nil
}

func writeSourceConfig(b *strings.Builder, source model.SourceConfig) error {
	src, ok := source.Property(model.CFGPROP_SOURCE)
	if !ok {
		return ErrNoSource
	}

	_, err := b.WriteString(fmt.Sprintf("# @sync %s%s", src, newLine))
	if err != nil {
		return err
	}

	props := source.Properties(nil)
	if len(props) > 1 {

		keys := maps.Keys(props)
		slices.Sort(keys)

		// write multiprop lines
		_, err = b.WriteString("# @props")
		if err != nil {
			return err
		}
		for _, k := range keys {
			if k != model.CFGPROP_SOURCE && !slices.Contains(singleLineProps, k) {
				_, err = b.WriteString(fmt.Sprintf(" %s=%s", k, props[k]))
				if err != nil {
					return err
				}
			}
		}
		b.WriteString(newLine)
		if err != nil {
			return err
		}

		// write singleprop lines
		for _, k := range singleLineProps {
			v, ok := props[k]
			if ok {
				_, err := b.WriteString(fmt.Sprintf("# @%s %s%s", k, v, newLine))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func osDefaultNewLine() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}
