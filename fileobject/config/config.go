// config parses a config file for rotation of File objects.
package config

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"agerotate"
	"agerotate/fileobject"
)

const (
	CommentChar = "#"
	PathPrefix  = "pathglob"
	RangePrefix = "range"
)

// Parse reads and parses a config.
func Parse(in io.Reader, fieldSep string) (fileobject.Files, []agerotate.Range, error) {
	return newParser(in, fieldSep).parse()
}

type parser struct {
	in       *bufio.Scanner
	lineNo   int
	line     string
	fieldSep string
	path     string
	ranges   []agerotate.Range
}

func newParser(in io.Reader, fieldSep string) *parser {
	return &parser{
		in:       bufio.NewScanner(in),
		fieldSep: fieldSep,
		ranges:   []agerotate.Range{},
	}
}

// parse manages the parser context and performs some sanity checking on the resulting objects before returning them.
func (p *parser) parse() (fileobject.Files, []agerotate.Range, error) {
	for p.in.Scan() {
		p.line = p.in.Text()
		p.lineNo += 1
		err := p.parseLine()
		if err != nil {
			return "", nil, err
		}
	}
	if err := p.in.Err(); err != nil {
		return "", nil, err
	}
	if p.path == "" {
		return "", nil, fmt.Errorf("No file rotation path specified")
	}
	if len(p.ranges) == 0 {
		return "", nil, fmt.Errorf("No ranges specified")
	}
	return fileobject.Files(p.path), p.ranges, nil
}

// parseLine parses the line that's just been read in by parse(), invoking handling functions specified to each line type.
func (p *parser) parseLine() error {
	p.line = clean(p.line)
	if p.line == "" {
		return nil
	}
	fields := strings.Split(p.line, p.fieldSep)
	if len(fields) < 2 {
		return fmt.Errorf("Line %d: Missing values", p.lineNo)
	}
	prefix := strings.ToLower(fields[0])
	switch prefix {
	case PathPrefix:
		return p.setPath(fields[1:])
	case RangePrefix:
		return p.addRange(fields[1:])
	default:
		return fmt.Errorf("Line %d: Invalid prefix %q", p.lineNo, prefix)
	}
	return nil
}

func (p *parser) setPath(values []string) error {
	if p.path != "" {
		return fmt.Errorf("Line %d: Duplicate path specification", p.lineNo)
	}
	if len(values) != 1 {
		return fmt.Errorf("Line %d: Using multiple path values is invalid", p.lineNo)
	}
	if values[0] == "" {
		return fmt.Errorf("Line %d: Must specify path", p.lineNo)
	}
	p.path = values[0]
	return nil
}

func (p *parser) addRange(values []string) error {
	if len(values) != 2 {
		return fmt.Errorf("Line %d: Range lines must have two values", p.lineNo)
	}
	age, err := time.ParseDuration(values[0])
	if err != nil {
		return fmt.Errorf("Line %d: Invalid age: %v", p.lineNo, err.Error())
	}
	interval, err := time.ParseDuration(values[1])
	if err != nil {
		return fmt.Errorf("Line %d: Invalid interval: %v", p.lineNo, err.Error())
	}
	if age < 0 {
		return fmt.Errorf("Line %d: Age values must be positive, got %v", p.lineNo, age)
	}
	if interval < 0 {
		return fmt.Errorf("Line %d: Interval values must be positive, got %v", p.lineNo, interval)
	}
	if len(p.ranges) > 0 && p.ranges[len(p.ranges)-1].Age >= age {
		return fmt.Errorf("Line %d: Age value must be larger than previous age value", p.lineNo)
	}
	p.ranges = append(p.ranges, agerotate.Range{Age: age, Interval: interval})
	return nil
}

// clean performs basic string normalization such as eliminating comments and whitespace.
func clean(s string) string {
	idx := strings.Index(s, CommentChar)
	if idx > -1 {
		s = s[0:idx]
		s = strings.TrimSpace(s)
	}
	return s
}
