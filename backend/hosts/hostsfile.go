package hosts

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/0xcfff/dnspipe/model"
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

func NewHostsFile(content []byte) *HostsFile {
	return &HostsFile{
		content: content,
	}
}

func (f *HostsFile) Dump() {
	s := string(f.content)
	fmt.Println(s)
}

func (f *HostsFile) Save(w io.Writer) (int, error) {
	return w.Write(f.content)
}

func (f *HostsFile) Parse(mode ParseMode) (*HostsFileContent, error) {
	return ParseHostsFileWithSources(bytes.NewReader(f.content), mode)
}

func (f *HostsFile) AppendSource(source model.SourceConfig) error {
	b := &strings.Builder{}

	err := writeOriginalLines(b, f.content, true)
	if err != nil {
		return fmt.Errorf("Error writting original lines %w", err)
	}

	err = writeSourceConfig(b, source)
	if err != nil {
		return fmt.Errorf("Error writting source config %w", err)
	}

	f.content = []byte(b.String())

	return nil
}

func (f *HostsFile) AppendIp(ip *IPRecord) error {
	b := &strings.Builder{}

	err := writeOriginalLines(b, f.content, false)
	if err != nil {
		return fmt.Errorf("Error writting original lines %w", err)
	}

	err = writeIpRecord(b, ip)
	if err != nil {
		return fmt.Errorf("Error writting ip %w", err)
	}

	f.content = []byte(b.String())

	return nil
}

func writeIpRecord(b *strings.Builder, ip *IPRecord) error {
	b.WriteString(ip.IP)
	for _, v := range ip.Aliases {
		b.WriteString(" ")
		b.WriteString(v)
	}
	if ip.Notes != "" {
		b.WriteString(" # ")
		b.WriteString(ip.Notes)
	}
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

func writeOriginalLines(b *strings.Builder, content []byte, ensureNewLine bool) error {
	s := bufio.NewScanner(bytes.NewReader(content))
	lastLineEmpty := true

	for {
		ok := s.Scan()
		if !ok {
			break
		}

		line := s.Text()
		lastLineEmpty = ensureNewLine && strings.TrimSpace(line) == ""

		_, err := b.WriteString(line)
		if err != nil {
			return err
		}
		_, err = b.WriteString(newLine)
		if err != nil {
			return err
		}
	}

	if !lastLineEmpty && ensureNewLine {
		b.WriteString(newLine)
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
