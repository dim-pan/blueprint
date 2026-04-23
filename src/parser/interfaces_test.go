package parser

import (
	"errors"
	"strings"
	"testing"
)

func parseIfaces(t *testing.T, path, body string) []InterfaceDefinition {
	t.Helper()
	out, err := ParseInterfacesBytes(path, []byte(body))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return out
}

func parseIfacesErr(t *testing.T, path, body string) *ParseError {
	t.Helper()
	_, err := ParseInterfacesBytes(path, []byte(body))
	if err == nil {
		t.Fatalf("expected ParseError, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	return pe
}

// TC-P02-01: Parse interface with primitive fields
func TestTC_P02_01_Primitives(t *testing.T) {
	body := `interface Todo {
  id:    String
  count: Integer
  done:  Boolean
}
`
	ifaces := parseIfaces(t, "foo.iface", body)
	if len(ifaces) != 1 {
		t.Fatalf("want 1, got %d", len(ifaces))
	}
	got := ifaces[0]
	if got.Name != "Todo" {
		t.Fatalf("name = %q", got.Name)
	}
	want := []Field{
		{Name: "id", Type: "String", Optional: false},
		{Name: "count", Type: "Integer", Optional: false},
		{Name: "done", Type: "Boolean", Optional: false},
	}
	if len(got.Fields) != len(want) {
		t.Fatalf("field count = %d, want %d", len(got.Fields), len(want))
	}
	for i, f := range want {
		if got.Fields[i] != f {
			t.Fatalf("field[%d] = %+v, want %+v", i, got.Fields[i], f)
		}
	}
}

// TC-P02-02: Optional field
func TestTC_P02_02_Optional(t *testing.T) {
	body := `interface X {
  note: String?
}
`
	f := parseIfaces(t, "x.iface", body)[0].Fields[0]
	if f.Type != "String" || !f.Optional {
		t.Fatalf("got %+v, want {String true}", f)
	}
}

// TC-P02-03: Array field
func TestTC_P02_03_Array(t *testing.T) {
	body := `interface X {
  tags: String[]
}
`
	f := parseIfaces(t, "x.iface", body)[0].Fields[0]
	if f.Type != "String[]" || f.Optional {
		t.Fatalf("got %+v, want {String[] false}", f)
	}
}

// TC-P02-04: Union type
func TestTC_P02_04_Union(t *testing.T) {
	body := `interface X {
  status: pending | active | done
}
`
	f := parseIfaces(t, "x.iface", body)[0].Fields[0]
	if f.Type != "pending | active | done" {
		t.Fatalf("type = %q", f.Type)
	}
}

// TC-P02-05: Reference to another interface type
func TestTC_P02_05_TypeReference(t *testing.T) {
	body := `interface X {
  owner: Component
}
`
	f := parseIfaces(t, "x.iface", body)[0].Fields[0]
	if f.Type != "Component" {
		t.Fatalf("type = %q", f.Type)
	}
}

// TC-P02-06: Multiple interfaces in one file
func TestTC_P02_06_Multiple(t *testing.T) {
	body := `interface A {
  x: String
}

interface B {
  y: Integer
}
`
	ifaces := parseIfaces(t, "x.iface", body)
	if len(ifaces) != 2 {
		t.Fatalf("want 2, got %d", len(ifaces))
	}
	if ifaces[0].Name != "A" || ifaces[1].Name != "B" {
		t.Fatalf("order wrong: %q then %q", ifaces[0].Name, ifaces[1].Name)
	}
}

// TC-P02-07: Whitespace alignment accepted; name excludes whitespace
func TestTC_P02_07_Alignment(t *testing.T) {
	body := `interface X {
  short:        String
  much_longer:  Integer
}
`
	fields := parseIfaces(t, "x.iface", body)[0].Fields
	if fields[0].Name != "short" || fields[1].Name != "much_longer" {
		t.Fatalf("names = %q, %q", fields[0].Name, fields[1].Name)
	}
	if fields[0].Type != "String" || fields[1].Type != "Integer" {
		t.Fatalf("types = %q, %q", fields[0].Type, fields[1].Type)
	}
}

// TC-P06-04: interface source location recorded
func TestTC_P06_04_InterfaceLocation(t *testing.T) {
	body := `

interface Target {
  a: String
}
`
	iface := parseIfaces(t, "sys/interfaces/foo.iface", body)[0]
	if !strings.HasSuffix(iface.SourceFile, "sys/interfaces/foo.iface") {
		t.Fatalf("source_file = %q", iface.SourceFile)
	}
	if iface.LineNumber != 3 {
		t.Fatalf("line_number = %d, want 3", iface.LineNumber)
	}
}

// TC-P07-07: missing opening brace
func TestTC_P07_07_MissingOpenBrace(t *testing.T) {
	body := `interface Foo
  a: String
}
`
	pe := parseIfacesErr(t, "x.iface", body)
	if !strings.Contains(pe.Message, "{") {
		t.Fatalf("message missing '{': %q", pe.Message)
	}
}

// TC-P07-08: missing closing brace
func TestTC_P07_08_MissingCloseBrace(t *testing.T) {
	body := `interface Foo {
  a: String
`
	pe := parseIfacesErr(t, "x.iface", body)
	if !strings.Contains(pe.Message, "}") {
		t.Fatalf("message missing '}': %q", pe.Message)
	}
}

// TC-P07-09: malformed field (no colon)
func TestTC_P07_09_MalformedField(t *testing.T) {
	body := `interface Foo {
  broken field line
}
`
	pe := parseIfacesErr(t, "x.iface", body)
	if !strings.Contains(pe.Message, "field") {
		t.Fatalf("message missing 'field': %q", pe.Message)
	}
}
