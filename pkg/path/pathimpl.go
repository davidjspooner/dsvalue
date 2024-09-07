package dspath

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/davidjspooner/dsvalue/pkg/key"
	"github.com/davidjspooner/dsvalue/pkg/value"
)

type Path []key.Interface

func EvaluateFieldFor(obj value.Value, field key.Interface) (value.Value, error) {
	kind := obj.Kind()
	switch kind {
	case value.MapKind:
		mapValue, ok := obj.(value.Map)
		if !ok {
			return nil, fmt.Errorf("expected map, but got %s", kind)
		}
		return mapValue.Field(field)
	case value.ArrayKind:
		arrayValue, ok := obj.(value.Array)
		if !ok {
			return nil, fmt.Errorf("expected array, but got %s", kind)
		}
		return arrayValue.Index(field)
	default:
		return nil, fmt.Errorf("expected map or array, but got %s", kind)
	}
}

func (p *Path) EvaluateFor(obj value.Value) (value.Value, error) {
	var err error
	for n, segment := range *p {
		obj, err = EvaluateFieldFor(obj, segment)
		if err != nil {
			sb := strings.Builder{}
			for i := 0; i <= n; i++ {
				sb.WriteString((*p)[i].String())
			}
			return nil, &ErrInvalidPath{Path: sb.String(), Inner: err}
		}
	}
	return obj, nil
}

func (p *Path) String() string {
	if len(*p) == 0 {
		return "."
	}
	sb := strings.Builder{}
	for _, segment := range *p {
		s := segment.String()
		sb.WriteString(s)
	}
	return sb.String()
}

func failedExpectation(text, expected string, actual *scanner.Scanner) error {
	actualText := actual.TokenText()
	if actualText == "" {
		return &ErrInvalidPath{
			Path:  text,
			Inner: fmt.Errorf("expected '%s', but got <EOF>", expected),
		}
	}
	return &ErrInvalidPath{
		Path:  text,
		Inner: fmt.Errorf("expected '%s', but got '%s'", expected, actualText),
	}
}

func CompilePath(text string) (Path, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(text))
	s.Mode ^= scanner.SkipComments // don't skip comments
	s.Whitespace = 0
	var path Path
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch tok {
		case '.':
			if tok = s.Scan(); tok != scanner.Ident {
				if tok == scanner.EOF && len(path) == 0 {
					return path, nil
				}

				return nil, failedExpectation(text, "identifier", &s)
			}
			path = append(path, key.Value[string]{X: s.TokenText()})
		case '[':
			tok = s.Scan()
			prefix := ""
			if tok == '-' {
				prefix = "-"
				tok = s.Scan()
			}
			switch tok {
			case scanner.Int:
				firstIndex := prefix + s.TokenText()
				tok = s.Scan()
				switch tok {
				case ']':
					index, err := strconv.Atoi(firstIndex)
					if err != nil {
						return nil, &ErrInvalidPath{Path: text, Inner: err}
					}
					path = append(path, key.Value[int]{X: index})
				case ':':
					prefix = ""
					if s.Scan() == '-' {
						prefix = "-"
						s.Scan()
					}
					secondIndex := prefix + s.TokenText()
					if s.Scan() != ']' {
						return nil, failedExpectation(text, "]", &s)
					}
					start, err := strconv.Atoi(firstIndex)
					if err != nil {
						return nil, &ErrInvalidPath{Path: text, Inner: err}
					}
					end, err := strconv.Atoi(secondIndex)
					if err != nil {
						return nil, &ErrInvalidPath{Path: text, Inner: err}
					}
					path = append(path, &key.Range{Start: start, End: end})
				default:
					return nil, failedExpectation(text, ": or ]", &s)
				}
			case ':':
				start := 0
				tok = s.Scan()
				prefix = ""
				if tok == '-' {
					prefix = "-"
					tok = s.Scan()
				}
				switch tok {
				case scanner.Int:
					secondIndex := prefix + s.TokenText()
					if s.Scan() != ']' {
						return nil, failedExpectation(text, "]", &s)
					}
					end, err := strconv.Atoi(secondIndex)
					if err != nil {
						return nil, &ErrInvalidPath{Path: text, Inner: err}
					}
					path = append(path, &key.Range{Start: start, End: end})
				case ']':
					path = append(path, &key.Range{Start: start, Tail: true})
				default:
					return nil, failedExpectation(text, "end index or ]", &s)
				}

			default:
				return nil, failedExpectation(text, "an index or range", &s)
			}
		default:
			return nil, failedExpectation(text, ". or [", &s)
		}
	}
	return path, nil
}
