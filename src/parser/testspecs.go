package parser

import (
	"os"
	"strings"
)

func ParseTestSpecsFile(path string) ([]TestSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseTestSpecsBytes(path, data)
}

func ParseTestSpecsBytes(path string, data []byte) ([]TestSpec, error) {
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	var out []TestSpec
	i := 0
	for i < len(lines) {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "##") {
			i++
			continue
		}
		if !isTestHeader(raw, trimmed) {
			return nil, &ParseError{
				File:    path,
				Line:    i + 1,
				Message: "expected 'test' block, blank line, or '##' section, got: " + trimmed,
			}
		}
		spec, next, err := parseTestBlock(path, lines, i)
		if err != nil {
			return nil, err
		}
		out = append(out, spec)
		i = next
	}
	return out, nil
}

func isTestHeader(raw, trimmed string) bool {
	if isIndented(raw) {
		return false
	}
	return trimmed == "test" ||
		strings.HasPrefix(trimmed, "test ") ||
		strings.HasPrefix(trimmed, "test\t")
}

func parseTestBlock(path string, lines []string, start int) (TestSpec, int, error) {
	id, title, err := parseIDTitleHeader(path, start+1, strings.TrimSpace(lines[start]), "test")
	if err != nil {
		return TestSpec{}, 0, err
	}

	spec := TestSpec{
		ID:         id,
		Title:      title,
		Verifies:   []string{},
		SourceFile: path,
		LineNumber: start + 1,
	}

	var verifiesFound, givenFound, expectFound bool

	i := start + 1
	for i < len(lines) {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			i++
			break
		}
		if strings.HasPrefix(trimmed, "##") {
			break
		}
		if isTestHeader(raw, trimmed) {
			break
		}
		if !isIndented(raw) {
			return TestSpec{}, 0, &ParseError{
				File:    path,
				Line:    i + 1,
				Message: "expected indented body line inside test block, got: " + trimmed,
			}
		}

		switch {
		case strings.HasPrefix(trimmed, "verifies:"):
			spec.Verifies = parseIDList(strings.TrimPrefix(trimmed, "verifies:"))
			verifiesFound = true
			i++
		case strings.HasPrefix(trimmed, "given:"):
			text, next := collectContinuation(lines, i, "given:", testSpecKeys)
			spec.Given = text
			givenFound = true
			i = next
		case strings.HasPrefix(trimmed, "when:"):
			text, next := collectContinuation(lines, i, "when:", testSpecKeys)
			spec.When = text
			i = next
		case strings.HasPrefix(trimmed, "expect:"):
			text, next := collectContinuation(lines, i, "expect:", testSpecKeys)
			spec.Expect = text
			expectFound = true
			i = next
		default:
			key := unknownKeyName(trimmed)
			return TestSpec{}, 0, &ParseError{
				File:    path,
				Line:    i + 1,
				Message: "unknown key '" + key + "' in test block",
			}
		}
	}

	if !verifiesFound {
		return TestSpec{}, 0, &ParseError{File: path, Line: start + 1, Message: "test '" + id + "' is missing a verifies line"}
	}
	if !givenFound {
		return TestSpec{}, 0, &ParseError{File: path, Line: start + 1, Message: "test '" + id + "' is missing a given clause"}
	}
	if !expectFound {
		return TestSpec{}, 0, &ParseError{File: path, Line: start + 1, Message: "test '" + id + "' is missing an expect clause"}
	}

	return spec, i, nil
}

func collectContinuation(lines []string, start int, keyPrefix string, knownKeys []string) (string, int) {
	first := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[start]), keyPrefix))
	var parts []string
	if first != "" {
		parts = append(parts, first)
	}
	i := start + 1
	for i < len(lines) {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			break
		}
		if !isIndented(raw) {
			break
		}
		if isKnownKeyLine(trimmed, knownKeys) {
			break
		}
		parts = append(parts, trimmed)
		i++
	}
	return joinStatement(parts), i
}
