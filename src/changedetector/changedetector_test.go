package changedetector

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/dim-pan/blueprint/src/parser"
)

func req(id, statement, allocatedTo string) parser.Requirement {
	return parser.Requirement{
		ID: id, Title: id,
		Priority: parser.PriorityMustHave, Pattern: parser.PatternUbiquitous,
		Statement: statement, AllocatedTo: allocatedTo,
		SourceFile: "r.req", LineNumber: 1,
	}
}

func derivedReq(id, parent, allocatedTo, statement string) parser.Requirement {
	r := req(id, statement, allocatedTo)
	r.DerivedFrom = parent
	return r
}

func comp(id string, responsibility string, satisfies, dependsOn []string) parser.Component {
	return parser.Component{
		ID: id, Name: id, Responsibility: responsibility,
		Satisfies: satisfies, DependsOn: dependsOn,
		SourceFile: id + ".component", LineNumber: 1,
	}
}

func iface(name string, fields ...parser.Field) parser.InterfaceDefinition {
	return parser.InterfaceDefinition{
		Name: name, Fields: fields, SourceFile: name + ".iface", LineNumber: 1,
	}
}

func test(id, expect string, verifies ...string) parser.TestSpec {
	return parser.TestSpec{
		ID: id, Title: id, Verifies: verifies,
		Given: "x", Expect: expect,
		SourceFile: id + ".testspec", LineNumber: 1,
	}
}

func hasChange(cs ChangeSet, id string, ct ChangeType) bool {
	for _, c := range cs.Changes {
		if c.ElementID == id && c.ChangeType == ct {
			return true
		}
	}
	return false
}

func affected(cs ChangeSet) []string {
	sort.Strings(cs.AffectedComponents)
	return cs.AffectedComponents
}

// TC-CD01-01
func TestTC_CD01_01_AddedRequirement(t *testing.T) {
	cs := DetectBetween(
		parser.SystemModel{},
		parser.SystemModel{Requirements: []parser.Requirement{req("REQ-001", "X", "")}},
	)
	if !hasChange(cs, "REQ-001", Added) {
		t.Fatalf("missing added REQ-001: %+v", cs.Changes)
	}
}

// TC-CD01-02
func TestTC_CD01_02_AddedComponent(t *testing.T) {
	cs := DetectBetween(
		parser.SystemModel{},
		parser.SystemModel{Components: []parser.Component{comp("COMP-001", "R", nil, nil)}},
	)
	if !hasChange(cs, "COMP-001", Added) {
		t.Fatalf("missing added COMP-001: %+v", cs.Changes)
	}
}

// TC-CD01-03
func TestTC_CD01_03_AddedInterfaceAndTest(t *testing.T) {
	cs := DetectBetween(
		parser.SystemModel{},
		parser.SystemModel{
			Interfaces: []parser.InterfaceDefinition{iface("Todo")},
			TestSpecs:  []parser.TestSpec{test("TC-001", "y")},
		},
	)
	if !hasChange(cs, "Todo", Added) || !hasChange(cs, "TC-001", Added) {
		t.Fatalf("missing additions: %+v", cs.Changes)
	}
}

// TC-CD02-01
func TestTC_CD02_01_ModifiedRequirement(t *testing.T) {
	baseline := parser.SystemModel{Requirements: []parser.Requirement{req("REQ-001", "X", "")}}
	current := parser.SystemModel{Requirements: []parser.Requirement{req("REQ-001", "Y", "")}}
	cs := DetectBetween(baseline, current)
	if !hasChange(cs, "REQ-001", Modified) {
		t.Fatalf("missing modified REQ-001: %+v", cs.Changes)
	}
}

// TC-CD02-02
func TestTC_CD02_02_SourceLocationIgnored(t *testing.T) {
	r1 := req("REQ-001", "X", "")
	r1.SourceFile = "old.req"
	r1.LineNumber = 1
	r2 := req("REQ-001", "X", "")
	r2.SourceFile = "new.req"
	r2.LineNumber = 42
	cs := DetectBetween(
		parser.SystemModel{Requirements: []parser.Requirement{r1}},
		parser.SystemModel{Requirements: []parser.Requirement{r2}},
	)
	if hasChange(cs, "REQ-001", Modified) {
		t.Fatalf("should not be modified from location-only changes: %+v", cs.Changes)
	}
}

// TC-CD02-03
func TestTC_CD02_03_IdenticalEmpty(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("R", "X", "")},
		Components:   []parser.Component{comp("C", "R", nil, nil)},
	}
	cs := DetectBetween(m, m)
	if len(cs.Changes) != 0 || len(cs.AffectedComponents) != 0 {
		t.Fatalf("expected no changes, got %+v", cs)
	}
}

// TC-CD03-01
func TestTC_CD03_01_RemovedRequirement(t *testing.T) {
	cs := DetectBetween(
		parser.SystemModel{Requirements: []parser.Requirement{req("REQ-001", "X", "")}},
		parser.SystemModel{},
	)
	if !hasChange(cs, "REQ-001", Removed) {
		t.Fatalf("missing removed REQ-001: %+v", cs.Changes)
	}
}

// TC-CD03-02
func TestTC_CD03_02_RemovedComponent(t *testing.T) {
	cs := DetectBetween(
		parser.SystemModel{Components: []parser.Component{comp("COMP-001", "R", nil, nil)}},
		parser.SystemModel{},
	)
	if !hasChange(cs, "COMP-001", Removed) {
		t.Fatalf("missing removed COMP-001: %+v", cs.Changes)
	}
}

// TC-CD04-01
func TestTC_CD04_01_DirectComponentChange(t *testing.T) {
	cs := DetectBetween(
		parser.SystemModel{Components: []parser.Component{comp("COMP-001", "R1", nil, nil)}},
		parser.SystemModel{Components: []parser.Component{comp("COMP-001", "R2", nil, nil)}},
	)
	if !reflect.DeepEqual(affected(cs), []string{"COMP-001"}) {
		t.Fatalf("affected = %v", cs.AffectedComponents)
	}
}

// TC-CD04-02
func TestTC_CD04_02_AllocatedReqChange(t *testing.T) {
	baseline := parser.SystemModel{
		Requirements: []parser.Requirement{derivedReq("REQ-001-A", "REQ-001", "COMP-001", "X")},
		Components:   []parser.Component{comp("COMP-001", "R", nil, nil)},
	}
	current := parser.SystemModel{
		Requirements: []parser.Requirement{derivedReq("REQ-001-A", "REQ-001", "COMP-001", "Y")},
		Components:   []parser.Component{comp("COMP-001", "R", nil, nil)},
	}
	cs := DetectBetween(baseline, current)
	if !reflect.DeepEqual(affected(cs), []string{"COMP-001"}) {
		t.Fatalf("affected = %v", cs.AffectedComponents)
	}
}

// TC-CD04-03
func TestTC_CD04_03_SatisfiedSystemReqChange(t *testing.T) {
	baseline := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "X", "")},
		Components:   []parser.Component{comp("COMP-001", "R", []string{"REQ-001"}, nil)},
	}
	current := parser.SystemModel{
		Requirements: []parser.Requirement{req("REQ-001", "Y", "")},
		Components:   []parser.Component{comp("COMP-001", "R", []string{"REQ-001"}, nil)},
	}
	cs := DetectBetween(baseline, current)
	if !reflect.DeepEqual(affected(cs), []string{"COMP-001"}) {
		t.Fatalf("affected = %v", cs.AffectedComponents)
	}
}

// TC-CD04-04
func TestTC_CD04_04_TransitiveDependsOn(t *testing.T) {
	baseline := parser.SystemModel{
		Components: []parser.Component{
			comp("COMP-001", "R1", nil, nil),
			comp("COMP-002", "R2", nil, []string{"COMP-001"}),
		},
	}
	current := parser.SystemModel{
		Components: []parser.Component{
			comp("COMP-001", "R1-changed", nil, nil),
			comp("COMP-002", "R2", nil, []string{"COMP-001"}),
		},
	}
	cs := DetectBetween(baseline, current)
	want := []string{"COMP-001", "COMP-002"}
	if !reflect.DeepEqual(affected(cs), want) {
		t.Fatalf("affected = %v, want %v", cs.AffectedComponents, want)
	}
}

// TC-CD04-05
func TestTC_CD04_05_TestSpecChange(t *testing.T) {
	baseline := parser.SystemModel{
		Requirements: []parser.Requirement{derivedReq("REQ-001-A", "REQ-001", "COMP-001", "X")},
		Components:   []parser.Component{comp("COMP-001", "R", nil, nil)},
		TestSpecs:    []parser.TestSpec{test("TC-001", "y-old", "REQ-001-A")},
	}
	current := parser.SystemModel{
		Requirements: []parser.Requirement{derivedReq("REQ-001-A", "REQ-001", "COMP-001", "X")},
		Components:   []parser.Component{comp("COMP-001", "R", nil, nil)},
		TestSpecs:    []parser.TestSpec{test("TC-001", "y-new", "REQ-001-A")},
	}
	cs := DetectBetween(baseline, current)
	if !reflect.DeepEqual(affected(cs), []string{"COMP-001"}) {
		t.Fatalf("affected = %v", cs.AffectedComponents)
	}
}

// TC-CD04-06
func TestTC_CD04_06_InterfaceChangeFansOut(t *testing.T) {
	baseline := parser.SystemModel{
		Interfaces: []parser.InterfaceDefinition{iface("Todo", parser.Field{Name: "id", Type: "String"})},
		Components: []parser.Component{comp("A", "R", nil, nil), comp("B", "R", nil, nil)},
	}
	current := parser.SystemModel{
		Interfaces: []parser.InterfaceDefinition{iface("Todo", parser.Field{Name: "id", Type: "Integer"})},
		Components: []parser.Component{comp("A", "R", nil, nil), comp("B", "R", nil, nil)},
	}
	cs := DetectBetween(baseline, current)
	want := []string{"A", "B"}
	if !reflect.DeepEqual(affected(cs), want) {
		t.Fatalf("affected = %v, want %v", cs.AffectedComponents, want)
	}
}

// TC-CD05-01
func TestTC_CD05_01_SaveLoadRoundtrip(t *testing.T) {
	m := parser.SystemModel{
		Requirements: []parser.Requirement{req("R", "X", "")},
		Components:   []parser.Component{comp("C", "R", nil, nil)},
		Interfaces:   []parser.InterfaceDefinition{iface("I")},
		TestSpecs:    []parser.TestSpec{test("T", "y", "R")},
	}
	path := filepath.Join(t.TempDir(), "baseline.json")
	if err := SaveBaseline(m, path); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatal(err)
	}
	cs := DetectBetween(loaded, m)
	if len(cs.Changes) != 0 {
		t.Fatalf("roundtrip produced diffs: %+v", cs.Changes)
	}
}

// TC-CD05-02
func TestTC_CD05_02_MissingBaselineTreatedAsEmpty(t *testing.T) {
	current := parser.SystemModel{
		Requirements: []parser.Requirement{req("R", "X", "")},
	}
	cs, err := Detect(current, filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasChange(cs, "R", Added) {
		t.Fatalf("expected R added: %+v", cs.Changes)
	}
}
