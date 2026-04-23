package validator

import (
	"strings"
	"testing"

	"github.com/dim-pan/blueprint/src/parser"
)

func req(id, derivedFrom, allocatedTo string) parser.Requirement {
	return parser.Requirement{
		ID:          id,
		Title:       id,
		Priority:    parser.PriorityMustHave,
		Pattern:     parser.PatternUbiquitous,
		Statement:   "The system shall " + id + ".",
		DerivedFrom: derivedFrom,
		AllocatedTo: allocatedTo,
		SourceFile:  "req.req",
		LineNumber:  1,
	}
}

func comp(id, name string, satisfies, dependsOn []string) parser.Component {
	return parser.Component{
		ID:             id,
		Name:           name,
		Responsibility: "Do " + name,
		Satisfies:      satisfies,
		DependsOn:      dependsOn,
		SourceFile:     "comp.component",
		LineNumber:     1,
	}
}

func iface(name string, fields ...parser.Field) parser.InterfaceDefinition {
	return parser.InterfaceDefinition{
		Name:       name,
		Fields:     fields,
		SourceFile: name + ".iface",
		LineNumber: 1,
	}
}

func field(name, typ string) parser.Field {
	return parser.Field{Name: name, Type: typ}
}

func test(id string, verifies ...string) parser.TestSpec {
	return parser.TestSpec{
		ID: id, Title: id, Verifies: verifies,
		Given: "x", Expect: "y",
		SourceFile: "test.testspec", LineNumber: 1,
	}
}

func assertValid(t *testing.T, r ValidationResult) {
	t.Helper()
	if !r.Valid {
		t.Fatalf("expected valid, got errors: %+v", r.Errors)
	}
}

func assertInvalid(t *testing.T, r ValidationResult, substrings ...string) {
	t.Helper()
	if r.Valid {
		t.Fatalf("expected invalid, got valid")
	}
	for _, want := range substrings {
		found := false
		for _, e := range r.Errors {
			if strings.Contains(e.Message, want) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("no error message containing %q; errors: %+v", want, r.Errors)
		}
	}
}

// TC-V01-01
func TestTC_V01_01_ValidSatisfies(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", "")},
		Components:   []parser.Component{comp("COMP-001", "X", []string{"REQ-001"}, nil)},
	}
	assertValid(t, Validate(m))
}

// TC-V01-02
func TestTC_V01_02_BrokenSatisfies(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("COMP-001", "X", []string{"REQ-NONEXISTENT"}, nil)},
	}
	assertInvalid(t, Validate(m), "satisfies", "REQ-NONEXISTENT")
}

// TC-V02-01
func TestTC_V02_01_ValidVerifies(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", "")},
		TestSpecs:    []parser.TestSpec{test("TC-001", "REQ-001")},
	}
	assertValid(t, Validate(m))
}

// TC-V02-02
func TestTC_V02_02_BrokenVerifies(t *testing.T) {
	m := parser.SystemModel{
		TestSpecs: []parser.TestSpec{test("TC-001", "REQ-NONEXISTENT")},
	}
	assertInvalid(t, Validate(m), "verifies", "REQ-NONEXISTENT")
}

// TC-V03-01
func TestTC_V03_01_DerivedFromResolves(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-001", "", ""),
			req("REQ-001-A", "REQ-001", "COMP-001"),
		},
		Components: []parser.Component{comp("COMP-001", "X", nil, nil)},
	}
	assertValid(t, Validate(m))
}

// TC-V03-02
func TestTC_V03_02_DerivedFromMustBeSystemLevel(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-001", "", ""),
			req("REQ-001-A", "REQ-001", "COMP-001"),
			req("REQ-001-A-x", "REQ-001-A", "COMP-001"),
		},
		Components: []parser.Component{comp("COMP-001", "X", nil, nil)},
	}
	assertInvalid(t, Validate(m), "derived_from", "system-level")
}

// TC-V03-03
func TestTC_V03_03_DerivedFromMissing(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-A", "REQ-MISSING", "")},
	}
	assertInvalid(t, Validate(m), "REQ-MISSING")
}

// TC-V04-01
func TestTC_V04_01_AllocatedToResolves(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("R", "", "COMP-001")},
		Components:   []parser.Component{comp("COMP-001", "X", nil, nil)},
	}
	assertValid(t, Validate(m))
}

// TC-V04-02
func TestTC_V04_02_AllocatedToBroken(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("R", "", "COMP-MISSING")},
	}
	assertInvalid(t, Validate(m), "allocated_to", "COMP-MISSING")
}

// TC-V05-01
func TestTC_V05_01_DependsOnResolves(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{
			comp("COMP-001", "A", nil, nil),
			comp("COMP-002", "B", nil, []string{"COMP-001"}),
		},
	}
	assertValid(t, Validate(m))
}

// TC-V05-02
func TestTC_V05_02_DependsOnBroken(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("COMP-001", "A", nil, []string{"COMP-MISSING"})},
	}
	assertInvalid(t, Validate(m), "depends_on", "COMP-MISSING")
}

// TC-V06-01
func TestTC_V06_01_Primitives(t *testing.T) {
	m := parser.SystemModel{
		Interfaces: []parser.InterfaceDefinition{
			iface("X",
				field("s", "String"),
				field("n", "Integer"),
				field("b", "Boolean"),
				field("f", "Float"),
				field("t", "Timestamp"),
			),
		},
	}
	assertValid(t, Validate(m))
}

// TC-V06-02
func TestTC_V06_02_InterfaceRef(t *testing.T) {
	m := parser.SystemModel{
		Interfaces: []parser.InterfaceDefinition{
			iface("A", field("b", "B")),
			iface("B", field("x", "String")),
		},
	}
	assertValid(t, Validate(m))
}

// TC-V06-03
func TestTC_V06_03_ArrayBaseChecked(t *testing.T) {
	m := parser.SystemModel{
		Interfaces: []parser.InterfaceDefinition{
			iface("A", field("xs", "Unknown[]")),
		},
	}
	assertInvalid(t, Validate(m), "Unknown")
}

// TC-V06-04
func TestTC_V06_04_UnionNotResolved(t *testing.T) {
	m := parser.SystemModel{
		Interfaces: []parser.InterfaceDefinition{
			iface("A", field("color", "red | green | blue")),
		},
	}
	assertValid(t, Validate(m))
}

// TC-V06-05
func TestTC_V06_05_UnknownScalar(t *testing.T) {
	m := parser.SystemModel{
		Interfaces: []parser.InterfaceDefinition{
			iface("A", field("w", "Widget")),
		},
	}
	assertInvalid(t, Validate(m), "Widget")
}

// TC-V07-01
func TestTC_V07_01_DuplicateRequirements(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", ""), req("REQ-001", "", "")},
	}
	assertInvalid(t, Validate(m), "duplicate", "REQ-001")
}

// TC-V07-02
func TestTC_V07_02_DuplicateComponents(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("COMP-001", "A", nil, nil), comp("COMP-001", "B", nil, nil)},
	}
	assertInvalid(t, Validate(m), "duplicate", "COMP-001")
}

// TC-V07-03
func TestTC_V07_03_DuplicateInterfaces(t *testing.T) {
	m := parser.SystemModel{
		Interfaces: []parser.InterfaceDefinition{iface("Todo"), iface("Todo")},
	}
	assertInvalid(t, Validate(m), "duplicate", "Todo")
}

// TC-V07-04
func TestTC_V07_04_DuplicateTestSpecs(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("R", "", "")},
		TestSpecs:    []parser.TestSpec{test("TC-001", "R"), test("TC-001", "R")},
	}
	assertInvalid(t, Validate(m), "duplicate", "TC-001")
}

// TC-V08-01 — error carries location of offending element
func TestTC_V08_01_ErrorLocation(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{{
			ID: "COMP-001", Name: "X", Responsibility: "R",
			Satisfies:  []string{"REQ-MISSING"},
			SourceFile: "components/foo/foo.component",
			LineNumber: 3,
		}},
	}
	r := Validate(m)
	if r.Valid {
		t.Fatalf("expected invalid")
	}
	found := false
	for _, e := range r.Errors {
		if strings.HasSuffix(e.SourceFile, "foo.component") && e.LineNumber == 3 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("no error at foo.component:3; errors: %+v", r.Errors)
	}
}

// TC-V09-01 — any error means invalid
func TestTC_V09_01_AnyErrorInvalid(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("C", "X", []string{"R-MISSING"}, nil)},
	}
	r := Validate(m)
	if r.Valid {
		t.Fatalf("expected invalid")
	}
}

// TC-V09-02 — no errors means valid with empty errors list
func TestTC_V09_02_NoErrorsValid(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("R", "", "")},
	}
	r := Validate(m)
	if !r.Valid {
		t.Fatalf("expected valid, got %+v", r.Errors)
	}
	if len(r.Errors) != 0 {
		t.Fatalf("errors list should be empty, got %+v", r.Errors)
	}
}
