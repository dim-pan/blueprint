package tracer

import (
	"reflect"
	"testing"

	"github.com/dim-pan/blueprint/src/parser"
)

func req(id, derivedFrom, allocatedTo string) parser.Requirement {
	return parser.Requirement{
		ID: id, Title: id,
		Priority: parser.PriorityMustHave, Pattern: parser.PatternUbiquitous,
		Statement: "X.", DerivedFrom: derivedFrom, AllocatedTo: allocatedTo,
		SourceFile: "req.req", LineNumber: 1,
	}
}

func comp(id string, satisfies ...string) parser.Component {
	return parser.Component{
		ID: id, Name: id, Responsibility: "R",
		Satisfies:  satisfies,
		SourceFile: "comp.component", LineNumber: 1,
	}
}

func test(id string, verifies ...string) parser.TestSpec {
	return parser.TestSpec{
		ID: id, Title: id, Verifies: verifies,
		Given: "x", Expect: "y",
		SourceFile: "t.testspec", LineNumber: 1,
	}
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// TC-T01-01
func TestTC_T01_01_ForwardSatisfied(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", "")},
		Components:   []parser.Component{comp("COMP-A", "REQ-001"), comp("COMP-B")},
	}
	tr := New(m)
	res := tr.Query(TraceQuery{FromID: "REQ-001", Direction: Forward})
	if len(res.Chains) != 1 {
		t.Fatalf("want 1 chain, got %d", len(res.Chains))
	}
	if !contains(res.Chains[0].SatisfiedBy, "COMP-A") {
		t.Fatalf("missing COMP-A: %v", res.Chains[0].SatisfiedBy)
	}
	if contains(res.Chains[0].SatisfiedBy, "COMP-B") {
		t.Fatalf("unexpected COMP-B: %v", res.Chains[0].SatisfiedBy)
	}
}

// TC-T01-02
func TestTC_T01_02_ForwardVerified(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", "")},
		TestSpecs:    []parser.TestSpec{test("TC-001", "REQ-001")},
	}
	tr := New(m)
	res := tr.Query(TraceQuery{FromID: "REQ-001", Direction: Forward})
	if !contains(res.Chains[0].VerifiedBy, "TC-001") {
		t.Fatalf("missing TC-001: %v", res.Chains[0].VerifiedBy)
	}
}

// TC-T01-03
func TestTC_T01_03_AllocatedCounts(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-001", "", ""),
			req("REQ-001-A", "REQ-001", "COMP-X"),
		},
		Components: []parser.Component{comp("COMP-X")}, // does NOT list REQ-001-A in satisfies
	}
	tr := New(m)
	res := tr.Query(TraceQuery{FromID: "REQ-001-A", Direction: Forward})
	if !contains(res.Chains[0].SatisfiedBy, "COMP-X") {
		t.Fatalf("missing COMP-X via allocated_to: %v", res.Chains[0].SatisfiedBy)
	}
}

// TC-T02-01
func TestTC_T02_01_BackwardFromTest(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", ""), req("REQ-002", "", "")},
		TestSpecs:    []parser.TestSpec{test("TC-001", "REQ-001", "REQ-002")},
	}
	tr := New(m)
	res := tr.Query(TraceQuery{FromID: "TC-001", Direction: Backward})
	if len(res.Chains) != 2 {
		t.Fatalf("want 2 chains, got %d", len(res.Chains))
	}
	ids := []string{res.Chains[0].RequirementID, res.Chains[1].RequirementID}
	if !contains(ids, "REQ-001") || !contains(ids, "REQ-002") {
		t.Fatalf("want REQ-001 and REQ-002, got %v", ids)
	}
}

// TC-T02-02
func TestTC_T02_02_BackwardFromComponent(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-001", "", ""),
			req("REQ-002-A", "REQ-002", "COMP-X"),
		},
		Components: []parser.Component{comp("COMP-X", "REQ-001")},
	}
	tr := New(m)
	res := tr.Query(TraceQuery{FromID: "COMP-X", Direction: Backward})
	var ids []string
	for _, c := range res.Chains {
		ids = append(ids, c.RequirementID)
	}
	if !contains(ids, "REQ-001") || !contains(ids, "REQ-002-A") {
		t.Fatalf("want REQ-001 and REQ-002-A, got %v", ids)
	}
}

// TC-T03-01
func TestTC_T03_01_AncestorChain(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-001", "", ""),
			req("REQ-001-A", "REQ-001", ""),
		},
	}
	tr := New(m)
	got := tr.AncestorChain("REQ-001-A")
	want := []string{"REQ-001-A", "REQ-001"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// TC-T03-02
func TestTC_T03_02_SystemLevelChainOne(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", "")},
	}
	tr := New(m)
	got := tr.AncestorChain("REQ-001")
	if !reflect.DeepEqual(got, []string{"REQ-001"}) {
		t.Fatalf("got %v", got)
	}
}

// TC-T03-03
func TestTC_T03_03_BrokenDerivedFromStops(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-X", "REQ-MISSING", "")},
	}
	tr := New(m)
	got := tr.AncestorChain("REQ-X")
	if !reflect.DeepEqual(got, []string{"REQ-X"}) {
		t.Fatalf("expected chain to stop at REQ-X, got %v", got)
	}
}

// TC-T04-01
func TestTC_T04_01_MatrixAllRequirements(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("A", "", ""), req("B", "", ""), req("C", "", ""),
		},
	}
	tr := New(m)
	res := tr.Matrix()
	if len(res.Chains) != 3 {
		t.Fatalf("want 3, got %d", len(res.Chains))
	}
	for i, want := range []string{"A", "B", "C"} {
		if res.Chains[i].RequirementID != want {
			t.Fatalf("chain[%d] = %q, want %q", i, res.Chains[i].RequirementID, want)
		}
	}
}

// TC-T04-02
func TestTC_T04_02_MatrixPopulatesLinks(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", "")},
		Components:   []parser.Component{comp("COMP-A", "REQ-001")},
		TestSpecs:    []parser.TestSpec{test("TC-001", "REQ-001")},
	}
	tr := New(m)
	res := tr.Matrix()
	if len(res.Chains) != 1 {
		t.Fatalf("want 1 chain, got %d", len(res.Chains))
	}
	c := res.Chains[0]
	if !reflect.DeepEqual(c.SatisfiedBy, []string{"COMP-A"}) {
		t.Fatalf("satisfied_by = %v", c.SatisfiedBy)
	}
	if !reflect.DeepEqual(c.VerifiedBy, []string{"TC-001"}) {
		t.Fatalf("verified_by = %v", c.VerifiedBy)
	}
}
