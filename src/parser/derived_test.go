package parser

import "testing"

// TC-P05-01: parse derived_from and allocated_to
func TestTC_P05_01_DerivedFields(t *testing.T) {
	body := `req BP-001-P01 "Parse reqs"
  derived_from: BP-001
  allocated_to: COMP-001
  priority: must-have
  The parser shall extract requirements.
`
	r := parse(t, "parser.req", body)[0]
	if r.ID != "BP-001-P01" {
		t.Fatalf("id = %q", r.ID)
	}
	if r.DerivedFrom != "BP-001" {
		t.Fatalf("derived_from = %q, want BP-001", r.DerivedFrom)
	}
	if r.AllocatedTo != "COMP-001" {
		t.Fatalf("allocated_to = %q, want COMP-001", r.AllocatedTo)
	}
	if r.Priority != PriorityMustHave {
		t.Fatalf("priority = %q", r.Priority)
	}
	if r.Statement != "The parser shall extract requirements." {
		t.Fatalf("statement = %q", r.Statement)
	}
}

// TC-P05-02: derived fields empty on system-level requirement
func TestTC_P05_02_EmptyWhenAbsent(t *testing.T) {
	body := `req REQ-001 "X"
  priority: must-have
  The system shall X.
`
	r := parse(t, "blueprint.req", body)[0]
	if r.DerivedFrom != "" || r.AllocatedTo != "" {
		t.Fatalf("expected empty derived fields, got %+v", r)
	}
}

// TC-P05-03: order of derived fields does not matter
func TestTC_P05_03_OrderIndependent(t *testing.T) {
	body := `req R1 "X"
  allocated_to: COMP-002
  derived_from: REQ-001
  priority: must-have
  The system shall X.
`
	r := parse(t, "x.req", body)[0]
	if r.DerivedFrom != "REQ-001" || r.AllocatedTo != "COMP-002" {
		t.Fatalf("wrong: derived_from=%q allocated_to=%q", r.DerivedFrom, r.AllocatedTo)
	}
}
