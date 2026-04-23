package parser

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func parseComps(t *testing.T, path, body string) []Component {
	t.Helper()
	out, err := ParseComponentsBytes(path, []byte(body))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return out
}

func parseCompsErr(t *testing.T, path, body string) *ParseError {
	t.Helper()
	_, err := ParseComponentsBytes(path, []byte(body))
	if err == nil {
		t.Fatalf("expected ParseError, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	return pe
}

// TC-P03-01: single-line responsibility
func TestTC_P03_01_Simple(t *testing.T) {
	body := `component COMP-001 "Parser"
  responsibility: Read model files.
  satisfies: BP-001
`
	c := parseComps(t, "x.component", body)[0]
	if c.ID != "COMP-001" || c.Name != "Parser" {
		t.Fatalf("id/name wrong: %+v", c)
	}
	if c.Responsibility != "Read model files." {
		t.Fatalf("responsibility = %q", c.Responsibility)
	}
	if !reflect.DeepEqual(c.Satisfies, []string{"BP-001"}) {
		t.Fatalf("satisfies = %v", c.Satisfies)
	}
	if len(c.DependsOn) != 0 {
		t.Fatalf("depends_on should be empty: %v", c.DependsOn)
	}
}

// TC-P03-02: multi-line responsibility
func TestTC_P03_02_WrappedResponsibility(t *testing.T) {
	body := `component COMP-001 "X"
  responsibility: Read the files
                  and return a
                  structured model.
  satisfies: BP-001
`
	c := parseComps(t, "x.component", body)[0]
	want := "Read the files and return a structured model."
	if c.Responsibility != want {
		t.Fatalf("responsibility = %q, want %q", c.Responsibility, want)
	}
}

// TC-P03-03: depends_on list
func TestTC_P03_03_DependsOn(t *testing.T) {
	body := `component COMP-003 "A"
  responsibility: X.
  depends_on: COMP-001, COMP-002, COMP-003
`
	c := parseComps(t, "x.component", body)[0]
	want := []string{"COMP-001", "COMP-002", "COMP-003"}
	if !reflect.DeepEqual(c.DependsOn, want) {
		t.Fatalf("depends_on = %v, want %v", c.DependsOn, want)
	}
}

// TC-P03-04: satisfies list
func TestTC_P03_04_Satisfies(t *testing.T) {
	body := `component COMP-002 "V"
  responsibility: X.
  satisfies: BP-002, BP-003
`
	c := parseComps(t, "x.component", body)[0]
	want := []string{"BP-002", "BP-003"}
	if !reflect.DeepEqual(c.Satisfies, want) {
		t.Fatalf("satisfies = %v", c.Satisfies)
	}
}

// TC-P03-05: no satisfies is OK
func TestTC_P03_05_NoSatisfies(t *testing.T) {
	body := `component COMP-007 "CLI"
  responsibility: X.
  depends_on: COMP-001
`
	c := parseComps(t, "x.component", body)[0]
	if len(c.Satisfies) != 0 {
		t.Fatalf("expected empty satisfies, got %v", c.Satisfies)
	}
}

// TC-P03-06: multiple components back-to-back
func TestTC_P03_06_MultipleNoSeparator(t *testing.T) {
	body := `component COMP-A "A"
  responsibility: A.
  satisfies: BP-001
component COMP-B "B"
  responsibility: B.
  satisfies: BP-002
`
	comps := parseComps(t, "x.component", body)
	if len(comps) != 2 || comps[0].ID != "COMP-A" || comps[1].ID != "COMP-B" {
		t.Fatalf("wrong parse: %+v", comps)
	}
}

// TC-P06-05: source location on component
func TestTC_P06_05_ComponentLocation(t *testing.T) {
	body := `
component COMP-001 "Parser"
  responsibility: X.
  satisfies: BP-001
`
	c := parseComps(t, "sys/components/parser/parser.component", body)[0]
	if c.LineNumber != 2 {
		t.Fatalf("line_number = %d, want 2", c.LineNumber)
	}
	if !strings.HasSuffix(c.SourceFile, "parser.component") {
		t.Fatalf("source_file = %q", c.SourceFile)
	}
}

// TC-P07-10: missing responsibility
func TestTC_P07_10_MissingResponsibility(t *testing.T) {
	body := `component COMP-001 "Parser"
  satisfies: BP-001
`
	pe := parseCompsErr(t, "x.component", body)
	if !strings.Contains(pe.Message, "responsibility") {
		t.Fatalf("message missing 'responsibility': %q", pe.Message)
	}
}

// TC-P07-11: unknown body key
func TestTC_P07_11_UnknownKey(t *testing.T) {
	body := `component COMP-001 "Parser"
  responsibility: X.
  foo: bar
`
	pe := parseCompsErr(t, "x.component", body)
	if !strings.Contains(pe.Message, "foo") {
		t.Fatalf("message missing 'foo': %q", pe.Message)
	}
}
