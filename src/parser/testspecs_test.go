package parser

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func parseTests(t *testing.T, path, body string) []TestSpec {
	t.Helper()
	out, err := ParseTestSpecsBytes(path, []byte(body))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return out
}

func parseTestsErr(t *testing.T, path, body string) *ParseError {
	t.Helper()
	_, err := ParseTestSpecsBytes(path, []byte(body))
	if err == nil {
		t.Fatalf("expected ParseError, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	return pe
}

// TC-P04-01: basic test spec
func TestTC_P04_01_Basic(t *testing.T) {
	body := `test TC-001 "Do the thing"
  verifies: REQ-001
  given: a condition
  expect: a result
`
	s := parseTests(t, "x.testspec", body)[0]
	if s.ID != "TC-001" || s.Title != "Do the thing" {
		t.Fatalf("id/title wrong: %+v", s)
	}
	if !reflect.DeepEqual(s.Verifies, []string{"REQ-001"}) {
		t.Fatalf("verifies = %v", s.Verifies)
	}
	if s.Given != "a condition" {
		t.Fatalf("given = %q", s.Given)
	}
	if s.Expect != "a result" {
		t.Fatalf("expect = %q", s.Expect)
	}
	if s.When != "" {
		t.Fatalf("when should be empty, got %q", s.When)
	}
}

// TC-P04-02: when clause
func TestTC_P04_02_WhenClause(t *testing.T) {
	body := `test TC-002 "X"
  verifies: REQ-002
  given: a starting state
  when: an action is taken
  expect: a new state
`
	s := parseTests(t, "x.testspec", body)[0]
	if s.When != "an action is taken" {
		t.Fatalf("when = %q", s.When)
	}
}

// TC-P04-03: multiple verifies
func TestTC_P04_03_MultipleVerifies(t *testing.T) {
	body := `test TC-003 "X"
  verifies: BP-002, BP-003
  given: x
  expect: y
`
	s := parseTests(t, "x.testspec", body)[0]
	want := []string{"BP-002", "BP-003"}
	if !reflect.DeepEqual(s.Verifies, want) {
		t.Fatalf("verifies = %v", s.Verifies)
	}
}

// TC-P04-04: multi-line given
func TestTC_P04_04_MultiLineGiven(t *testing.T) {
	body := `test TC-004 "X"
  verifies: REQ-004
  given: a sys/ directory
         containing .req, .iface,
         and .component files
  expect: success
`
	s := parseTests(t, "x.testspec", body)[0]
	want := "a sys/ directory containing .req, .iface, and .component files"
	if s.Given != want {
		t.Fatalf("given = %q, want %q", s.Given, want)
	}
}

// TC-P04-05: multi-line expect
func TestTC_P04_05_MultiLineExpect(t *testing.T) {
	body := `test TC-005 "X"
  verifies: REQ-005
  given: x
  expect: a result with
          multiple fields
          set correctly
`
	s := parseTests(t, "x.testspec", body)[0]
	want := "a result with multiple fields set correctly"
	if s.Expect != want {
		t.Fatalf("expect = %q, want %q", s.Expect, want)
	}
}

// TC-P04-06: ## section headers skipped
func TestTC_P04_06_SectionHeadersSkipped(t *testing.T) {
	body := `## BP-001: Something

test TC-A "A"
  verifies: BP-001
  given: x
  expect: y

## BP-002: Something else

test TC-B "B"
  verifies: BP-002
  given: x
  expect: y
`
	specs := parseTests(t, "x.testspec", body)
	if len(specs) != 2 {
		t.Fatalf("want 2, got %d", len(specs))
	}
	if specs[0].ID != "TC-A" || specs[1].ID != "TC-B" {
		t.Fatalf("ids = %q, %q", specs[0].ID, specs[1].ID)
	}
}

// TC-P04-07: multiple tests
func TestTC_P04_07_Multiple(t *testing.T) {
	body := `test TC-A "A"
  verifies: R1
  given: x
  expect: y

test TC-B "B"
  verifies: R2
  given: x
  expect: y

test TC-C "C"
  verifies: R3
  given: x
  expect: y
`
	specs := parseTests(t, "x.testspec", body)
	if len(specs) != 3 {
		t.Fatalf("want 3, got %d", len(specs))
	}
}

// TC-P06-06: source location
func TestTC_P06_06_TestSpecLocation(t *testing.T) {
	body := `## Section

test TC-001 "X"
  verifies: R1
  given: x
  expect: y
`
	// headers ## skipped, blank line, then test on line 3
	s := parseTests(t, "sys/tests/foo.testspec", body)[0]
	if s.LineNumber != 3 {
		t.Fatalf("line_number = %d, want 3", s.LineNumber)
	}
	if !strings.HasSuffix(s.SourceFile, "foo.testspec") {
		t.Fatalf("source_file = %q", s.SourceFile)
	}
}

// TC-P07-12: missing given
func TestTC_P07_12_MissingGiven(t *testing.T) {
	body := `test TC-001 "X"
  verifies: R1
  expect: y
`
	pe := parseTestsErr(t, "x.testspec", body)
	if !strings.Contains(pe.Message, "given") {
		t.Fatalf("message missing 'given': %q", pe.Message)
	}
}

// TC-P07-13: missing expect
func TestTC_P07_13_MissingExpect(t *testing.T) {
	body := `test TC-001 "X"
  verifies: R1
  given: x
`
	pe := parseTestsErr(t, "x.testspec", body)
	if !strings.Contains(pe.Message, "expect") {
		t.Fatalf("message missing 'expect': %q", pe.Message)
	}
}
