package hosts

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"regexp"
	"strings"

	"github.com/0xcfff/dnspipe/model"
	log "github.com/sirupsen/logrus"
)

type ParseMode int

const (
	// Any error in file leads to parsing failure
	Strict ParseMode = iota
	// Only errors in sync blocks lead to errors
	Moderate
	// Any errors which can be ignored are ignored
	Safe
)

type Position struct {
	Line int
}

type IPRecord struct {
	Pos Position

	IP      string
	Aliases []string
	Notes   string
}

type InlineProperty struct {
	Pos Position

	Name  string
	Value string
}

type SyncBlock struct {
	Pos Position

	Text        string
	InlineProps []*InlineProperty

	PosEndHeader Position

	Data *SyncDataBlock
}

type SyncDataBlock struct {
	Pos Position

	IPRecords []*IPRecord

	PosEndData Position
}

type HostsFileContent struct {
	IPRecords   []*IPRecord
	SyncBlocks  []*SyncBlock
	ContentHash []byte
}

type hostsParseContext struct {
	scanner      bufio.Scanner
	hasher       hash.Hash
	mode         ParseMode
	target       *HostsFileContent
	curLine      string
	lineNum      int
	lineReturned bool
	position     int
	parseSources bool
}

var (
	rxSyncBlockBegin = regexp.MustCompile(`^\s*#\s+@sync\s+`)
	rxSyncBlockProps = regexp.MustCompile(`^\s*#\s+@(params?|props?)\s+`)
	rxSyncBlockLine  = regexp.MustCompile(`^\s*#\s+@`)
	rxDataBlockBegin = regexp.MustCompile(`^\s*#\s+@begin_sync\s*`)
	rxDataBlockEnd   = regexp.MustCompile(`^\s*#\s+@end_sync\s*`)
	rxCommentLine    = regexp.MustCompile(`^\s*#\s*`)

	rxEmpty                = regexp.MustCompile(`^\s*$`)
	rxIpAddress            = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})|([\da-fA-F]{0,}:[\da-fA-F]{0,}:[\da-fA-F]{0,}(:[\da-fA-F]{0,}){0,5})`)
	rxSingleLineProps      = regexp.MustCompile(`\s*(\w[\w\d]+)\s*=\s*(\S+)\s*,?`)
	rxSingleLinePropsCheck = regexp.MustCompile(`^(\s*(\w[\w\d]+)\s*=\s*(\S+)\s*,?)+$`)
)

// Parses passed file in /etc/hosts format
// The fuction extracts list of IP <=> Domain Name mappings
// as well as additional synchronization blocks descrioption
func ParseHostsFileWithSources(r io.Reader, parseMode ParseMode) (*HostsFileContent, error) {
	ctx := newHostsParseContext(r, parseMode, true)
	result, err := ctx.parse()
	return result, err
}

// Parses passed file in /etc/hosts format
// The fuction extracts list of IP <=> Domain Name mappings
// this method does not parse synchronization blocks descrioption
func ParseHostsFileWithoutSources(r io.Reader, parseMode ParseMode) (*HostsFileContent, error) {
	ctx := newHostsParseContext(r, parseMode, false)
	result, err := ctx.parse()
	return result, err
}

func newHostsParseContext(r io.Reader, parseMode ParseMode, parseSources bool) *hostsParseContext {
	hasher := sha1.New()
	rr := io.TeeReader(r, hasher)
	scanner := bufio.NewScanner(rr)
	result := HostsFileContent{}
	parser := hostsParseContext{
		scanner:      *scanner,
		target:       &result,
		hasher:       hasher,
		mode:         parseMode,
		parseSources: parseSources,
	}

	return &parser
}

func (ctx *hostsParseContext) parse() (*HostsFileContent, error) {

	result := ctx.target

	for {
		line, ok := ctx.readLine()
		if !ok {
			if err := ctx.error(); err != nil {
				return nil, err
			}
			break // parsing is finished
		}

		// line analysis
		if rxEmpty.MatchString(line) {
			continue
		} else if ctx.parseSources && rxSyncBlockBegin.MatchString(line) {
			record, err := ctx.parseSyncBlock()
			if err != nil {
				return nil, fmt.Errorf("Error parsing sync block at line %d, error: %w", ctx.lineNum, err)
			}
			result.SyncBlocks = append(result.SyncBlocks, record)
		} else if rxCommentLine.MatchString(line) {
			continue // skip comments
		} else {
			record, err := ctx.parseIpLine()
			if err != nil {
				if ctx.mode != Strict {
					log.Warnf("Error parsing content line of hosts file, line: %d, error: %s", ctx.lineNum, err)
					continue
				}
				return nil, fmt.Errorf("Error parsing content line (IPs) of hosts file, line: %d, error: %w", ctx.lineNum, err)
			}
			result.IPRecords = append(result.IPRecords, record)
		}
	}

	result.ContentHash = ctx.hasher.Sum(nil)
	return result, nil
}

func (ctx *hostsParseContext) parseSyncBlock() (*SyncBlock, error) {
	line, ok := ctx.currentLine()
	if !ok {
		panic("should never be called in this state")
	}

	if !rxSyncBlockBegin.MatchString(line) {
		return nil, fmt.Errorf("Line %d is not a start sync line", ctx.lineNum)
	}

	rest := trimPrefixRegex(line, rxSyncBlockBegin)

	if rxEmpty.MatchString(rest) {
		return nil, fmt.Errorf("Line %d is not a start sync line (malformed)", ctx.lineNum)
	}

	lastHeadLine := ctx.lineNum
	record := SyncBlock{
		Pos: Position{
			Line: ctx.lineNum,
		},
	}

	props := ctx.parseInlineProps(rest)
	if len(props) == 0 {
		prop := InlineProperty{
			Pos:   Position{Line: ctx.lineNum},
			Name:  model.CFGPROP_SOURCE,
			Value: rest,
		}
		record.InlineProps = append(record.InlineProps, &prop)
	} else {
		record.InlineProps = append(record.InlineProps, props...)
	}

	for {
		line, ok = ctx.readLine()
		if !ok {
			break
		}

		if rxEmpty.MatchString(line) {
			continue
		} else if rxSyncBlockBegin.MatchString(line) {
			ctx.returnLine()
			break
		} else if rxDataBlockBegin.MatchString(line) {
			origLine := ctx.lineNum
			data, err := ctx.parseDataBlock()
			if err != nil {
				log.Warnf("error parsing sync data block, line %d, error: %s", ctx.lineNum, err)
				if ctx.mode != Safe {
					return nil, fmt.Errorf("error pasring sync data block, line %d, error: %w", ctx.lineNum, err)
				}
			}
			if record.Data != nil {
				log.Warnf("several data blocks foud for one sync block, line %d", origLine)
				if ctx.mode == Strict {
					return nil, fmt.Errorf("several data blocks found for one sync block, line %d", origLine)
				}
				log.Warnln("skipping data block")
			} else {
				record.Data = data
			}
		} else if rxSyncBlockProps.MatchString(line) {
			rest = trimPrefixRegex(line, rxSyncBlockProps)
			props = ctx.parseInlineProps(rest)
			if len(props) == 0 {
				log.Warnf("properties block does not contain properties, line %d", ctx.lineNum)
				if ctx.mode == Strict {
					return nil, fmt.Errorf("properties block does not contain properties, line %d", ctx.lineNum)
				}
			} else {
				record.InlineProps = append(record.InlineProps, props...)
			}
			lastHeadLine = ctx.lineNum
		} else if rxSyncBlockLine.MatchString(line) {
			rest = trimPrefixRegex(line, rxSyncBlockLine)
			parts := strings.SplitN(rest, " ", 2) //todo: this needs to be replaced with more advanced implementation which takes into accout tabs
			if len(parts) != 2 {
				log.Warnf("property block does not contain correct property value, line %d", ctx.lineNum)
				if ctx.mode == Strict {
					return nil, fmt.Errorf("property block does not contain correct property value, line %d", ctx.lineNum)
				}
			} else {
				prop := InlineProperty{
					Pos: Position{
						Line: ctx.lineNum,
					},
					Name:  parts[0],
					Value: parts[1],
				}
				record.InlineProps = append(record.InlineProps, &prop)
			}
			lastHeadLine = ctx.lineNum
		} else {
			ctx.returnLine()
			break
		}

	}

	record.PosEndHeader = Position{Line: lastHeadLine}
	return &record, nil
}
func (ctx *hostsParseContext) parseDataBlock() (*SyncDataBlock, error) {
	line, ok := ctx.currentLine()
	if !ok {
		panic("should never be called in this state")
	}

	if !rxDataBlockBegin.MatchString(line) {
		return nil, fmt.Errorf("Line %d is not a start sync line", ctx.lineNum)
	}

	lastDataLine := ctx.lineNum
	record := SyncDataBlock{
		Pos: Position{Line: ctx.lineNum},
	}

	for {
		line, ok = ctx.readLine()

		if !ok {
			log.Warnf("Incomplete data sync block, line: %d", ctx.lineNum)
			if ctx.mode != Safe {
				return nil, fmt.Errorf("Incomplete data sync block, line: %d", ctx.lineNum)
			}
			break
		} else if rxEmpty.MatchString(line) {
			continue
		} else if rxDataBlockEnd.MatchString(line) {
			lastDataLine = ctx.lineNum
			break
		} else if rxCommentLine.MatchString(line) {
			continue
		} else {
			ip, err := ctx.parseIpLine()
			if err != nil {
				if ctx.mode != Strict {
					log.Warnf("Error parsing sync content line of hosts file, line: %d, error: %s", ctx.lineNum, err)
					continue
				}
				return nil, fmt.Errorf("Error parsing sync content line (IPs) of hosts file, line: %d, error: %w", ctx.lineNum, err)
			}
			record.IPRecords = append(record.IPRecords, ip)
			// this may lead to file corruption if ctx.Mode != strict was used to parse file
			// and then the data was used for updating the hosts file
			lastDataLine = ctx.lineNum
		}
	}

	record.PosEndData = Position{Line: lastDataLine}
	return &record, nil
}

func (ctx *hostsParseContext) parseIpLine() (*IPRecord, error) {
	line, ok := ctx.currentLine()
	if !ok {
		panic("should never be called in this state")
	}

	parts := strings.Fields(line)
	if len(parts) < 2 || !rxIpAddress.MatchString(parts[0]) {
		return nil, fmt.Errorf("wrong format of line which starts as IP alias line, line: %d", ctx.lineNum)
	}

	record := IPRecord{
		Pos: Position{
			Line: ctx.lineNum,
		},
		IP: parts[0],
	}

	aliases := parts[1:]

	for _, val := range aliases {
		if strings.HasPrefix(val, "#") {
			idx := strings.Index(line, val)
			record.Notes = strings.TrimSpace(line[idx+1:])
			break
		}
		record.Aliases = append(record.Aliases, val)
	}

	if len(record.Aliases) == 0 {
		return nil, fmt.Errorf("wrong format of line which starts as IP alias line, line: %d", ctx.lineNum)
	}

	return &record, nil
}

func (ctx *hostsParseContext) parseInlineProps(trimmedLine string) []*InlineProperty {
	var result []*InlineProperty

	if !rxSingleLinePropsCheck.MatchString(trimmedLine) {
		return nil
	}

	syncProps := rxSingleLineProps.FindAllStringSubmatch(trimmedLine, -1)
	for _, match := range syncProps {
		prop := InlineProperty{
			Pos: Position{
				Line: ctx.lineNum,
			},
			Name:  match[1],
			Value: strings.TrimSuffix(match[2], ","),
		}
		result = append(result, &prop)
	}
	return result
}

func (ctx *hostsParseContext) readLine() (string, bool) {
	if ctx.lineReturned {
		ctx.lineReturned = false
		ctx.lineNum++
		return ctx.curLine, true
	}

	if ctx.scanner.Scan() {
		ctx.curLine = ctx.scanner.Text()
		ctx.lineNum++
		return ctx.curLine, true
	}

	ctx.curLine = ""

	return "", false
}

func (ctx *hostsParseContext) currentLine() (string, bool) {
	if ctx.lineReturned {
		panic("this method should never be called with such a state")
	}
	return ctx.curLine, true
}

func (ctx *hostsParseContext) returnLine() {
	if ctx.lineReturned {
		panic("this method should never be called with such a state")
	}
	ctx.lineReturned = true
	ctx.lineNum--
}

func (ctx *hostsParseContext) error() error {
	return ctx.scanner.Err()
}

func trimPrefixRegex(src string, rx *regexp.Regexp) string {
	prefix := rx.FindString(src)
	rest := strings.TrimPrefix(src, prefix)
	return rest
}
