package parser

import (
	"os"
	"strings"
)

func ParseInterfacesFile(path string) ([]InterfaceDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseInterfacesBytes(path, data)
}

func ParseInterfacesBytes(path string, data []byte) ([]InterfaceDefinition, error) {
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	var out []InterfaceDefinition
	i := 0
	for i < len(lines) {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" {
			i++
			continue
		}
		if !strings.HasPrefix(trimmed, "interface") {
			return nil, &ParseError{
				File:    path,
				Line:    i + 1,
				Message: "expected 'interface' block or blank line, got: " + trimmed,
			}
		}
		iface, next, err := parseInterfaceBlock(path, lines, i)
		if err != nil {
			return nil, err
		}
		out = append(out, iface)
		i = next
	}
	return out, nil
}

func parseInterfaceBlock(path string, lines []string, start int) (InterfaceDefinition, int, error) {
	header := strings.TrimSpace(lines[start])
	name, braceOnHeader, err := parseInterfaceHeader(path, start+1, header)
	if err != nil {
		return InterfaceDefinition{}, 0, err
	}

	iface := InterfaceDefinition{
		Name:       name,
		SourceFile: path,
		LineNumber: start + 1,
		Fields:     []Field{},
	}

	i := start + 1
	if !braceOnHeader {
		return InterfaceDefinition{}, 0, &ParseError{
			File:    path,
			Line:    start + 1,
			Message: "expected '{' after interface name",
		}
	}

	for i < len(lines) {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			i++
			continue
		}
		if trimmed == "}" {
			return iface, i + 1, nil
		}

		field, err := parseFieldLine(path, i+1, trimmed)
		if err != nil {
			return InterfaceDefinition{}, 0, err
		}
		iface.Fields = append(iface.Fields, field)
		i++
	}

	return InterfaceDefinition{}, 0, &ParseError{
		File:    path,
		Line:    start + 1,
		Message: "interface '" + name + "' is missing closing '}'",
	}
}

func parseInterfaceHeader(path string, lineNum int, header string) (string, bool, error) {
	rest := strings.TrimSpace(strings.TrimPrefix(header, "interface"))
	if rest == "" || rest == header {
		return "", false, &ParseError{File: path, Line: lineNum, Message: "expected interface name after 'interface'"}
	}

	braceOnHeader := false
	if strings.HasSuffix(rest, "{") {
		braceOnHeader = true
		rest = strings.TrimSpace(strings.TrimSuffix(rest, "{"))
	}

	if rest == "" {
		return "", false, &ParseError{File: path, Line: lineNum, Message: "expected interface name after 'interface'"}
	}
	if strings.ContainsAny(rest, " \t") {
		return "", false, &ParseError{File: path, Line: lineNum, Message: "unexpected tokens in interface header: " + rest}
	}

	return rest, braceOnHeader, nil
}

func parseFieldLine(path string, lineNum int, trimmed string) (Field, error) {
	colon := strings.Index(trimmed, ":")
	if colon == -1 {
		return Field{}, &ParseError{
			File:    path,
			Line:    lineNum,
			Message: "malformed field (expected 'name: Type'): " + trimmed,
		}
	}
	name := strings.TrimSpace(trimmed[:colon])
	typ := strings.TrimSpace(trimmed[colon+1:])
	if name == "" {
		return Field{}, &ParseError{File: path, Line: lineNum, Message: "field name missing before ':'"}
	}
	if typ == "" {
		return Field{}, &ParseError{File: path, Line: lineNum, Message: "field type missing after ':'"}
	}

	optional := false
	if strings.HasSuffix(typ, "?") {
		optional = true
		typ = strings.TrimSpace(strings.TrimSuffix(typ, "?"))
	}

	typ = strings.Join(strings.Fields(typ), " ")

	return Field{Name: name, Type: typ, Optional: optional}, nil
}
