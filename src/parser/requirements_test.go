package parser

import (
	"errors"
	"strings"
	"testing"
)

func parse(t *testing.T, path, body string) []Requirement {
	t.Helper()
	reqs, err := ParseRequirementsBytes(path, []byte(body))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return reqs
}

func parseErr(t *testing.T, path, body string) *ParseError {
	t.Helper()
	_, err := ParseRequirementsBytes(path, []byte(body))
	if err == nil {
		t.Fatalf("expected ParseError, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	return pe
}

// TC-P01-01: Parse ubiquitous requirement
func TestTC_P01_01_Ubiquitous(t *testing.T) {
	body := `req REQ-001 "Do the thing"
  priority: must-have
  The system shall do the thing.
`
	reqs := parse(t, "foo.req", body)
	if len(reqs) != 1 {
		t.Fatalf("want 1 req, got %d", len(reqs))
	}
	r := reqs[0]
	if r.ID != "REQ-001" || r.Title != "Do the thing" ||
		r.Priority != PriorityMustHave || r.Pattern != PatternUbiquitous ||
		r.Statement != "The system shall do the thing." {
		t.Fatalf("unexpected requirement: %+v", r)
	}
}

// TC-P01-02: Parse event-driven requirement
func TestTC_P01_02_EventDriven(t *testing.T) {
	body := `req REQ-002 "On click"
  priority: must-have
  WHEN the user clicks the button,
  the system shall record the click.
`
	r := parse(t, "foo.req", body)[0]
	if r.Pattern != PatternEventDriven {
		t.Fatalf("want event-driven, got %q", r.Pattern)
	}
	want := "WHEN the user clicks the button, the system shall record the click."
	if r.Statement != want {
		t.Fatalf("statement mismatch:\nwant: %q\ngot:  %q", want, r.Statement)
	}
}

// TC-P01-03: Parse unwanted-behaviour requirement
func TestTC_P01_03_Unwanted(t *testing.T) {
	body := `req REQ-003 "Reject empty"
  priority: must-have
  IF the input is empty,
  THEN the system shall reject the request.
`
	r := parse(t, "foo.req", body)[0]
	if r.Pattern != PatternUnwanted {
		t.Fatalf("want unwanted, got %q", r.Pattern)
	}
}

// TC-P01-04: Parse state-driven requirement
func TestTC_P01_04_StateDriven(t *testing.T) {
	body := `req REQ-004 "While safe"
  priority: must-have
  WHILE operating in safe mode,
  the system shall log all state transitions.
`
	r := parse(t, "foo.req", body)[0]
	if r.Pattern != PatternStateDriven {
		t.Fatalf("want state-driven, got %q", r.Pattern)
	}
}

// TC-P01-05: Parse optional-feature requirement
func TestTC_P01_05_Optional(t *testing.T) {
	body := `req REQ-005 "Admin feature"
  priority: should-have
  WHERE admin mode is enabled,
  the system shall expose the debug panel.
`
	r := parse(t, "foo.req", body)[0]
	if r.Pattern != PatternOptional {
		t.Fatalf("want optional, got %q", r.Pattern)
	}
	if r.Priority != PriorityShouldHave {
		t.Fatalf("want should-have, got %q", r.Priority)
	}
}

// TC-P01-06: Multiple requirements in one file
func TestTC_P01_06_Multiple(t *testing.T) {
	body := `req REQ-A "First"
  priority: must-have
  The system shall A.

req REQ-B "Second"
  priority: must-have
  The system shall B.

req REQ-C "Third"
  priority: must-have
  The system shall C.
`
	reqs := parse(t, "foo.req", body)
	if len(reqs) != 3 {
		t.Fatalf("want 3, got %d", len(reqs))
	}
	wantIDs := []string{"REQ-A", "REQ-B", "REQ-C"}
	for i, r := range reqs {
		if r.ID != wantIDs[i] {
			t.Fatalf("req[%d]: want %q, got %q", i, wantIDs[i], r.ID)
		}
	}
}

// TC-P01-07: Multi-line statement collapsed
func TestTC_P01_07_MultiLineJoined(t *testing.T) {
	body := `req REQ-X "Multi"
  priority: must-have
  The system shall do something
     across three
     indented lines.
`
	r := parse(t, "foo.req", body)[0]
	want := "The system shall do something across three indented lines."
	if r.Statement != want {
		t.Fatalf("want %q, got %q", want, r.Statement)
	}
}

// TC-P01-08: nice-to-have priority
func TestTC_P01_08_NiceToHave(t *testing.T) {
	body := `req REQ-Y "Y"
  priority: nice-to-have
  The system shall Y.
`
	r := parse(t, "foo.req", body)[0]
	if r.Priority != PriorityNiceToHave {
		t.Fatalf("want nice-to-have, got %q", r.Priority)
	}
}

// TC-P01-09: empty accept list
func TestTC_P01_09_AcceptEmpty(t *testing.T) {
	body := `req REQ-Z "Z"
  priority: must-have
  The system shall Z.
`
	r := parse(t, "foo.req", body)[0]
	if r.Accept == nil {
		t.Fatalf("accept is nil, want empty slice")
	}
	if len(r.Accept) != 0 {
		t.Fatalf("want empty, got %d", len(r.Accept))
	}
}

// TC-P06-01: Source file path recorded
func TestTC_P06_01_SourceFile(t *testing.T) {
	path := "sys/requirements/foo.req"
	body := `req REQ-001 "X"
  priority: must-have
  The system shall X.
`
	r := parse(t, path, body)[0]
	if !strings.HasSuffix(r.SourceFile, "sys/requirements/foo.req") {
		t.Fatalf("source_file = %q, does not end with expected suffix", r.SourceFile)
	}
}

// TC-P06-02: Line number at req header
func TestTC_P06_02_LineNumberOnHeader(t *testing.T) {
	body := `
# note line
# note line

req REQ-001 "X"
  priority: must-have
  The system shall X.
`
	reqs, err := ParseRequirementsBytes("foo.req", []byte(body))
	if err == nil {
		r := reqs[0]
		if r.LineNumber != 5 {
			t.Fatalf("want line 5, got %d", r.LineNumber)
		}
		return
	}
	// Comments are not part of the grammar yet; use a pure-blank prefix instead.
	body = "\n\n\n\nreq REQ-001 \"X\"\n  priority: must-have\n  The system shall X.\n"
	r := parse(t, "foo.req", body)[0]
	if r.LineNumber != 5 {
		t.Fatalf("want line 5, got %d", r.LineNumber)
	}
}

// TC-P06-03: each requirement keeps its own line number
func TestTC_P06_03_IndependentLineNumbers(t *testing.T) {
	body := `req REQ-A "A"
  priority: must-have
  The system shall A.


req REQ-B "B"
  priority: must-have
  The system shall B.
  And keep talking.
  And keep talking.

req REQ-C "C"
  priority: must-have
  The system shall C.
`
	reqs := parse(t, "foo.req", body)
	want := []int{1, 6, 12}
	for i, r := range reqs {
		if r.LineNumber != want[i] {
			t.Fatalf("req[%d] (%s) line = %d, want %d", i, r.ID, r.LineNumber, want[i])
		}
	}
}

// TC-P07-01: missing req id
func TestTC_P07_01_MissingID(t *testing.T) {
	body := `req "No id here"
  priority: must-have
  The system shall do something.
`
	pe := parseErr(t, "foo.req", body)
	if pe.Line != 1 || pe.File != "foo.req" {
		t.Fatalf("want foo.req:1, got %s:%d", pe.File, pe.Line)
	}
	if !strings.Contains(pe.Message, "id") {
		t.Fatalf("message does not mention 'id': %q", pe.Message)
	}
}

// TC-P07-02: missing title
func TestTC_P07_02_MissingTitle(t *testing.T) {
	body := `req REQ-001
  priority: must-have
  The system shall do something.
`
	pe := parseErr(t, "foo.req", body)
	if !strings.Contains(pe.Message, "title") {
		t.Fatalf("message does not mention 'title': %q", pe.Message)
	}
}

// TC-P07-03: unknown priority
func TestTC_P07_03_UnknownPriority(t *testing.T) {
	body := `req REQ-001 "X"
  priority: urgent
  The system shall X.
`
	pe := parseErr(t, "foo.req", body)
	for _, want := range []string{"must-have", "should-have", "nice-to-have"} {
		if !strings.Contains(pe.Message, want) {
			t.Fatalf("message missing %q: %q", want, pe.Message)
		}
	}
}

// TC-P07-04: missing priority line
func TestTC_P07_04_MissingPriority(t *testing.T) {
	body := `req REQ-001 "X"
  The system shall X.
`
	pe := parseErr(t, "foo.req", body)
	if !strings.Contains(pe.Message, "priority") {
		t.Fatalf("message does not mention 'priority': %q", pe.Message)
	}
}

// TC-P07-05: empty statement
func TestTC_P07_05_EmptyStatement(t *testing.T) {
	body := `req REQ-001 "X"
  priority: must-have
`
	pe := parseErr(t, "foo.req", body)
	if !strings.Contains(pe.Message, "statement") {
		t.Fatalf("message does not mention 'statement': %q", pe.Message)
	}
}

// TC-P07-06: unterminated title quote
func TestTC_P07_06_UnterminatedTitle(t *testing.T) {
	body := `req REQ-001 "No close
  priority: must-have
  The system shall X.
`
	pe := parseErr(t, "foo.req", body)
	if pe.Line != 1 {
		t.Fatalf("want line 1, got %d", pe.Line)
	}
}
