package parser

import (
	"os"
	"strings"
)

func ParseComponentsFile(path string) ([]Component, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseComponentsBytes(path, data)
}

func ParseComponentsBytes(path string, data []byte) ([]Component, error) {
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	var out []Component
	i := 0
	for i < len(lines) {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			i++
			continue
		}
		if !isComponentHeader(raw, trimmed) {
			return nil, &ParseError{
				File:    path,
				Line:    i + 1,
				Message: "expected 'component' block or blank line, got: " + trimmed,
			}
		}
		comp, next, err := parseComponentBlock(path, lines, i)
		if err != nil {
			return nil, err
		}
		out = append(out, comp)
		i = next
	}
	return out, nil
}

func isComponentHeader(raw, trimmed string) bool {
	if isIndented(raw) {
		return false
	}
	return trimmed == "component" ||
		strings.HasPrefix(trimmed, "component ") ||
		strings.HasPrefix(trimmed, "component\t")
}

func parseComponentBlock(path string, lines []string, start int) (Component, int, error) {
	id, name, err := parseIDTitleHeader(path, start+1, strings.TrimSpace(lines[start]), "component")
	if err != nil {
		return Component{}, 0, err
	}

	comp := Component{
		ID:         id,
		Name:       name,
		SourceFile: path,
		LineNumber: start + 1,
		DependsOn:  []string{},
		Satisfies:  []string{},
	}

	i := start + 1
	var respLines []string
	responsibilityFound := false

	for i < len(lines) {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			i++
			break
		}
		if isComponentHeader(raw, trimmed) {
			break
		}
		if !isIndented(raw) {
			return Component{}, 0, &ParseError{
				File:    path,
				Line:    i + 1,
				Message: "expected indented body line inside component block, got: " + trimmed,
			}
		}

		switch {
		case strings.HasPrefix(trimmed, "responsibility:"):
			responsibilityFound = true
			first := strings.TrimSpace(strings.TrimPrefix(trimmed, "responsibility:"))
			respLines = nil
			if first != "" {
				respLines = append(respLines, first)
			}
			i++
			for i < len(lines) {
				raw2 := lines[i]
				trimmed2 := strings.TrimSpace(raw2)
				if trimmed2 == "" {
					break
				}
				if !isIndented(raw2) {
					break
				}
				if looksLikeKeyLine(trimmed2) {
					break
				}
				respLines = append(respLines, trimmed2)
				i++
			}
		case strings.HasPrefix(trimmed, "depends_on:"):
			comp.DependsOn = parseIDList(strings.TrimPrefix(trimmed, "depends_on:"))
			i++
		case strings.HasPrefix(trimmed, "satisfies:"):
			comp.Satisfies = parseIDList(strings.TrimPrefix(trimmed, "satisfies:"))
			i++
		default:
			key := unknownKeyName(trimmed)
			return Component{}, 0, &ParseError{
				File:    path,
				Line:    i + 1,
				Message: "unknown key '" + key + "' in component block",
			}
		}
	}

	if !responsibilityFound {
		return Component{}, 0, &ParseError{
			File:    path,
			Line:    start + 1,
			Message: "component '" + id + "' is missing a responsibility line",
		}
	}

	comp.Responsibility = joinStatement(respLines)
	return comp, i, nil
}

var componentKeys = []string{"responsibility", "depends_on", "satisfies"}
var testSpecKeys = []string{"verifies", "given", "when", "expect"}

func isKnownKeyLine(trimmed string, keys []string) bool {
	for _, k := range keys {
		if strings.HasPrefix(trimmed, k+":") {
			return true
		}
	}
	return false
}

func looksLikeKeyLine(trimmed string) bool {
	for i, c := range trimmed {
		if c == ':' {
			return i > 0
		}
		if !((c >= 'a' && c <= 'z') || c == '_') {
			return false
		}
	}
	return false
}


func parseIDList(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func unknownKeyName(trimmed string) string {
	if colon := strings.Index(trimmed, ":"); colon != -1 {
		return strings.TrimSpace(trimmed[:colon])
	}
	return trimmed
}
