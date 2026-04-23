package parser

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func miniSysDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "requirements/app.req"), `req REQ-001 "X"
  priority: must-have
  The system shall do X.
`)
	writeFile(t, filepath.Join(dir, "interfaces/todo.iface"), `interface Todo {
  id: String
}
`)
	writeFile(t, filepath.Join(dir, "components/svc/svc.component"), `component COMP-001 "Svc"
  responsibility: Do things.
  satisfies: REQ-001
`)
	writeFile(t, filepath.Join(dir, "components/svc/svc.req"), `req REQ-001-01 "Derived"
  derived_from: REQ-001
  allocated_to: COMP-001
  priority: must-have
  The service shall do the thing.
`)
	writeFile(t, filepath.Join(dir, "components/svc/svc.testspec"), `test TC-SVC-001 "Does X"
  verifies: REQ-001-01
  given: input
  expect: output
`)
	writeFile(t, filepath.Join(dir, "tests/app.testspec"), `test TC-SYS-001 "System does X"
  verifies: REQ-001
  given: input
  expect: output
`)
	return dir
}

// TC-P08-01: full sys directory
func TestTC_P08_01_FullDir(t *testing.T) {
	dir := miniSysDir(t)
	m, err := ParseSystemModel(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Requirements) == 0 {
		t.Fatalf("no requirements")
	}
	if len(m.Components) == 0 {
		t.Fatalf("no components")
	}
	if len(m.Interfaces) == 0 {
		t.Fatalf("no interfaces")
	}
	if len(m.TestSpecs) == 0 {
		t.Fatalf("no test specs")
	}
}

// TC-P08-02: derived requirements included in model.requirements
func TestTC_P08_02_DerivedIncluded(t *testing.T) {
	dir := miniSysDir(t)
	m, _ := ParseSystemModel(dir)

	var haveSystem, haveDerived bool
	for _, r := range m.Requirements {
		if r.ID == "REQ-001" && r.DerivedFrom == "" {
			haveSystem = true
		}
		if r.ID == "REQ-001-01" && r.DerivedFrom == "REQ-001" && r.AllocatedTo == "COMP-001" {
			haveDerived = true
		}
	}
	if !haveSystem {
		t.Fatalf("system-level REQ-001 missing from requirements")
	}
	if !haveDerived {
		t.Fatalf("derived REQ-001-01 missing from requirements")
	}
}

// TC-P08-03: missing subdirectory tolerated
func TestTC_P08_03_MissingSubdir(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "requirements/x.req"), `req R "X"
  priority: must-have
  The system shall X.
`)
	m, err := ParseSystemModel(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.TestSpecs) != 0 {
		t.Fatalf("test_specs should be empty: %d", len(m.TestSpecs))
	}
	if len(m.Requirements) != 1 {
		t.Fatalf("want 1 requirement, got %d", len(m.Requirements))
	}
}

// TC-P08-04: empty sys dir
func TestTC_P08_04_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	m, err := ParseSystemModel(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Requirements)+len(m.Interfaces)+len(m.Components)+len(m.TestSpecs) != 0 {
		t.Fatalf("expected empty model, got %+v", m)
	}
}

// TC-P08-05: nonexistent path
func TestTC_P08_05_NonexistentPath(t *testing.T) {
	_, err := ParseSystemModel("/tmp/does/not/exist-blueprint-test-xyz")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TC-P08-06: parse error surfaces with file path
func TestTC_P08_06_ParseErrorSurfaced(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "requirements/broken.req"), `req "No id here"
  priority: must-have
  The system shall X.
`)
	_, err := ParseSystemModel(dir)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *ParseError, got %T: %v", err, err)
	}
	if !strings.HasSuffix(pe.File, "broken.req") {
		t.Fatalf("error file = %q, want suffix broken.req", pe.File)
	}
}
