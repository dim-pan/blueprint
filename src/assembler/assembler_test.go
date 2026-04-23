package assembler

import (
	"testing"

	"github.com/dim-pan/blueprint/src/parser"
)

func req(id, derivedFrom, allocatedTo string) parser.Requirement {
	return parser.Requirement{
		ID: id, Title: id,
		Priority: parser.PriorityMustHave, Pattern: parser.PatternUbiquitous,
		Statement: "X.", DerivedFrom: derivedFrom, AllocatedTo: allocatedTo,
		SourceFile: "r.req", LineNumber: 1,
	}
}

func comp(id string, satisfies ...string) parser.Component {
	return parser.Component{
		ID: id, Name: id, Responsibility: "R",
		Satisfies: satisfies, SourceFile: "c.component", LineNumber: 1,
	}
}

func iface(name string, fields ...parser.Field) parser.InterfaceDefinition {
	return parser.InterfaceDefinition{
		Name: name, Fields: fields, SourceFile: name + ".iface", LineNumber: 1,
	}
}

func test(id string, verifies ...string) parser.TestSpec {
	return parser.TestSpec{
		ID: id, Title: id, Verifies: verifies,
		Given: "x", Expect: "y",
		SourceFile: "t.testspec", LineNumber: 1,
	}
}

func ctxFor(ctxs []AgentContext, id string) *AgentContext {
	for i := range ctxs {
		if ctxs[i].Component.ID == id {
			return &ctxs[i]
		}
	}
	return nil
}

func hasReq(c *AgentContext, id string) bool {
	for _, r := range c.Requirements {
		if r.ID == id {
			return true
		}
	}
	return false
}

func hasTest(c *AgentContext, id string) bool {
	for _, t := range c.TestSpecs {
		if t.ID == id {
			return true
		}
	}
	return false
}

func hasInterface(c *AgentContext, name string) bool {
	for _, i := range c.Interfaces {
		if i.Name == name {
			return true
		}
	}
	return false
}

// TC-A01-01
func TestTC_A01_01_ContextIncludesComponent(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("COMP-001")},
	}
	c := ctxFor(Assemble(m), "COMP-001")
	if c == nil || c.Component.ID != "COMP-001" {
		t.Fatalf("missing or wrong component: %+v", c)
	}
}

// TC-A01-02
func TestTC_A01_02_OnlyAllocatedReqs(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-A-01", "", "COMP-001"),
			req("REQ-B-01", "", "COMP-002"),
		},
		Components: []parser.Component{comp("COMP-001"), comp("COMP-002")},
	}
	c := ctxFor(Assemble(m), "COMP-001")
	if !hasReq(c, "REQ-A-01") || hasReq(c, "REQ-B-01") {
		t.Fatalf("wrong requirements: %+v", c.Requirements)
	}
}

// TC-A01-03
func TestTC_A01_03_OnlyVerifyingTestSpecs(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-A-01", "", "COMP-001"),
			req("REQ-B-01", "", "COMP-002"),
		},
		Components: []parser.Component{comp("COMP-001"), comp("COMP-002")},
		TestSpecs: []parser.TestSpec{
			test("TC-001", "REQ-A-01"),
			test("TC-002", "REQ-B-01"),
		},
	}
	c := ctxFor(Assemble(m), "COMP-001")
	if !hasTest(c, "TC-001") || hasTest(c, "TC-002") {
		t.Fatalf("wrong test_specs: %+v", c.TestSpecs)
	}
}

// TC-A01-04
func TestTC_A01_04_OnePerComponent(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("A"), comp("B"), comp("C")},
	}
	ctxs := Assemble(m)
	if len(ctxs) != 3 {
		t.Fatalf("want 3 contexts, got %d", len(ctxs))
	}
}

// TC-A02-01
func TestTC_A02_01_InterfacesIncluded(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("COMP-001")},
		Interfaces: []parser.InterfaceDefinition{iface("Todo"), iface("TodoResult")},
	}
	c := ctxFor(Assemble(m), "COMP-001")
	if !hasInterface(c, "Todo") || !hasInterface(c, "TodoResult") {
		t.Fatalf("missing interfaces: %+v", c.Interfaces)
	}
}

// TC-A02-02
func TestTC_A02_02_ReferencedInterfacesResolve(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("COMP-001")},
		Interfaces: []parser.InterfaceDefinition{
			iface("AgentContext", parser.Field{Name: "component", Type: "Component"}),
			iface("Component", parser.Field{Name: "id", Type: "String"}),
		},
	}
	c := ctxFor(Assemble(m), "COMP-001")
	names := map[string]bool{}
	for _, i := range c.Interfaces {
		names[i.Name] = true
	}
	for _, i := range c.Interfaces {
		for _, f := range i.Fields {
			base := stripArray(f.Type)
			if isPrimitive(base) || isUnion(f.Type) {
				continue
			}
			if !names[base] {
				t.Fatalf("field %q.%q references %q which is absent from context interfaces", i.Name, f.Name, base)
			}
		}
	}
}

// TC-A03-01
func TestTC_A03_01_ParentReqIncluded(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-001", "", ""),
			req("REQ-001-A", "REQ-001", "COMP-001"),
		},
		Components: []parser.Component{comp("COMP-001")},
	}
	c := ctxFor(Assemble(m), "COMP-001")
	if !hasReq(c, "REQ-001-A") || !hasReq(c, "REQ-001") {
		t.Fatalf("missing reqs: %+v", c.Requirements)
	}
}

// TC-A03-02
func TestTC_A03_02_SatisfiedSystemReqIncluded(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "", "")},
		Components:   []parser.Component{comp("COMP-001", "REQ-001")},
	}
	c := ctxFor(Assemble(m), "COMP-001")
	if !hasReq(c, "REQ-001") {
		t.Fatalf("missing REQ-001: %+v", c.Requirements)
	}
}

// TC-A03-03
func TestTC_A03_03_ParentIncludedOnce(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{
			req("REQ-001", "", ""),
			req("REQ-001-A", "REQ-001", "COMP-001"),
			req("REQ-001-B", "REQ-001", "COMP-001"),
		},
		Components: []parser.Component{comp("COMP-001")},
	}
	c := ctxFor(Assemble(m), "COMP-001")
	count := 0
	for _, r := range c.Requirements {
		if r.ID == "REQ-001" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("REQ-001 appears %d times, want 1: %+v", count, c.Requirements)
	}
}

// TC-A08-01
func TestTC_A08_01_AssembleOnlyFilters(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("A"), comp("B"), comp("C")},
	}
	ctxs := AssembleOnly(m, []string{"A", "C"})
	if len(ctxs) != 2 {
		t.Fatalf("want 2, got %d", len(ctxs))
	}
	ids := []string{ctxs[0].Component.ID, ctxs[1].Component.ID}
	for _, want := range []string{"A", "C"} {
		found := false
		for _, g := range ids {
			if g == want {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing %s in %v", want, ids)
		}
	}
}

// TC-A08-02
func TestTC_A08_02_AssembleOnlyEmpty(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("A"), comp("B"), comp("C")},
	}
	ctxs := AssembleOnly(m, nil)
	if len(ctxs) != 0 {
		t.Fatalf("want empty, got %+v", ctxs)
	}
}

// TC-A08-03
func TestTC_A08_03_AssembleOnlyIgnoresUnknown(t *testing.T) {
	m := parser.SystemModel{
		Components: []parser.Component{comp("A"), comp("B")},
	}
	ctxs := AssembleOnly(m, []string{"A", "UNKNOWN"})
	if len(ctxs) != 1 || ctxs[0].Component.ID != "A" {
		t.Fatalf("got %+v", ctxs)
	}
}

// helpers used only by the assertion in TC-A02-02
func stripArray(t string) string {
	if len(t) > 2 && t[len(t)-2:] == "[]" {
		return t[:len(t)-2]
	}
	return t
}

func isPrimitive(t string) bool {
	switch t {
	case "String", "Integer", "Boolean", "Float", "Timestamp":
		return true
	}
	return false
}

func isUnion(t string) bool {
	for _, c := range t {
		if c == '|' {
			return true
		}
	}
	return false
}
