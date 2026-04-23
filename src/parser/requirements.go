package parser

import (
	"bufio"
	"os"
	"strings"
)

var validPriorities = map[string]Priority{
	"must-have":    PriorityMustHave,
	"should-have":  PriorityShouldHave,
	"nice-to-have": PriorityNiceToHave,
}

func ParseRequirementsFile(path string) ([]Requirement, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parseRequirements(path, lines)
}

func ParseRequirementsBytes(path string, data []byte) ([]Requirement, error) {
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return parseRequirements(path, lines)
}

func parseRequirements(path string, lines []string) ([]Requirement, error) {
	var out []Requirement
	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			i++
			continue
		}
		if strings.HasPrefix(trimmed, "req") && isReqHeader(trimmed) {
			req, next, err := parseReqBlock(path, lines, i)
			if err != nil {
				return nil, err
			}
			out = append(out, req)
			i = next
			continue
		}
		return nil, &ParseError{
			File:    path,
			Line:    i + 1,
			Message: "expected 'req' block or blank line, got: " + trimmed,
		}
	}
	return out, nil
}

func isReqHeader(trimmed string) bool {
	if trimmed == "req" {
		return true
	}
	return strings.HasPrefix(trimmed, "req ") || strings.HasPrefix(trimmed, "req\t")
}

func parseReqBlock(path string, lines []string, start int) (Requirement, int, error) {
	header := strings.TrimSpace(lines[start])
	id, title, err := parseReqHeader(path, start+1, header)
	if err != nil {
		return Requirement{}, 0, err
	}

	req := Requirement{
		ID:         id,
		Title:      title,
		SourceFile: path,
		LineNumber: start + 1,
		Accept:     []AcceptanceCriterion{},
	}

	i := start + 1
	priorityFound := false
	var statementLines []string
	var statementFirstLine int

	for i < len(lines) {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			i++
			break
		}
		if isReqHeader(trimmed) && !strings.HasPrefix(raw, " ") && !strings.HasPrefix(raw, "\t") {
			break
		}

		if !isIndented(raw) {
			return Requirement{}, 0, &ParseError{
				File:    path,
				Line:    i + 1,
				Message: "expected indented body line inside req block, got: " + trimmed,
			}
		}

		if strings.HasPrefix(trimmed, "priority:") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "priority:"))
			p, ok := validPriorities[value]
			if !ok {
				return Requirement{}, 0, &ParseError{
					File:    path,
					Line:    i + 1,
					Message: "unknown priority '" + value + "' (allowed: must-have, should-have, nice-to-have)",
				}
			}
			req.Priority = p
			priorityFound = true
			i++
			continue
		}
		if strings.HasPrefix(trimmed, "derived_from:") {
			req.DerivedFrom = strings.TrimSpace(strings.TrimPrefix(trimmed, "derived_from:"))
			i++
			continue
		}
		if strings.HasPrefix(trimmed, "allocated_to:") {
			req.AllocatedTo = strings.TrimSpace(strings.TrimPrefix(trimmed, "allocated_to:"))
			i++
			continue
		}

		if len(statementLines) == 0 {
			statementFirstLine = i + 1
		}
		statementLines = append(statementLines, trimmed)
		i++
	}

	if !priorityFound {
		return Requirement{}, 0, &ParseError{
			File:    path,
			Line:    start + 1,
			Message: "req '" + id + "' is missing a priority line",
		}
	}
	if len(statementLines) == 0 {
		return Requirement{}, 0, &ParseError{
			File:    path,
			Line:    start + 1,
			Message: "req '" + id + "' has an empty statement",
		}
	}

	req.Statement = joinStatement(statementLines)
	req.Pattern = detectPattern(req.Statement)
	_ = statementFirstLine

	return req, i, nil
}

func parseReqHeader(path string, lineNum int, header string) (string, string, error) {
	return parseIDTitleHeader(path, lineNum, header, "req")
}

func parseIDTitleHeader(path string, lineNum int, header, keyword string) (string, string, error) {
	rest := strings.TrimSpace(strings.TrimPrefix(header, keyword))
	if rest == "" {
		return "", "", &ParseError{File: path, Line: lineNum, Message: "expected id after '" + keyword + "'"}
	}
	if strings.HasPrefix(rest, "\"") {
		return "", "", &ParseError{File: path, Line: lineNum, Message: "expected id after '" + keyword + "' before the title"}
	}

	idEnd := strings.IndexAny(rest, " \t")
	if idEnd == -1 {
		return "", "", &ParseError{File: path, Line: lineNum, Message: "expected quoted title after " + keyword + " id"}
	}
	id := rest[:idEnd]
	titleSection := strings.TrimSpace(rest[idEnd:])

	if !strings.HasPrefix(titleSection, "\"") {
		return "", "", &ParseError{File: path, Line: lineNum, Message: "expected quoted title after " + keyword + " id"}
	}
	closing := strings.Index(titleSection[1:], "\"")
	if closing == -1 {
		return "", "", &ParseError{File: path, Line: lineNum, Message: "unterminated title quote"}
	}
	title := titleSection[1 : 1+closing]
	return id, title, nil
}

func isIndented(raw string) bool {
	return strings.HasPrefix(raw, " ") || strings.HasPrefix(raw, "\t")
}

func joinStatement(lines []string) string {
	parts := make([]string, 0, len(lines))
	for _, l := range lines {
		parts = append(parts, strings.Join(strings.Fields(l), " "))
	}
	return strings.Join(parts, " ")
}

func detectPattern(statement string) Pattern {
	switch {
	case startsWithWord(statement, "WHEN"):
		return PatternEventDriven
	case startsWithWord(statement, "IF"):
		return PatternUnwanted
	case startsWithWord(statement, "WHILE"):
		return PatternStateDriven
	case startsWithWord(statement, "WHERE"):
		return PatternOptional
	default:
		return PatternUbiquitous
	}
}

func startsWithWord(s, word string) bool {
	if !strings.HasPrefix(s, word) {
		return false
	}
	if len(s) == len(word) {
		return true
	}
	next := s[len(word)]
	return next == ' ' || next == '\t'
}
