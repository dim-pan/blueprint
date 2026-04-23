package tracer

import (
	"github.com/dim-pan/blueprint/src/parser"
)

type Direction string

const (
	Forward  Direction = "forward"
	Backward Direction = "backward"
)

type TraceChain struct {
	RequirementID string
	SatisfiedBy   []string
	VerifiedBy    []string
}

type TraceQuery struct {
	FromID    string
	Direction Direction
}

type TraceResult struct {
	Chains []TraceChain
}

type Tracer struct {
	model         parser.SystemModel
	reqByID       map[string]parser.Requirement
	componentByID map[string]parser.Component
	testByID      map[string]parser.TestSpec

	satisfiedBy map[string][]string
	verifiedBy  map[string][]string
}

func New(m parser.SystemModel) *Tracer {
	tr := &Tracer{
		model:         m,
		reqByID:       map[string]parser.Requirement{},
		componentByID: map[string]parser.Component{},
		testByID:      map[string]parser.TestSpec{},
		satisfiedBy:   map[string][]string{},
		verifiedBy:    map[string][]string{},
	}
	for _, r := range m.Requirements {
		tr.reqByID[r.ID] = r
		if r.AllocatedTo != "" {
			tr.satisfiedBy[r.ID] = appendUnique(tr.satisfiedBy[r.ID], r.AllocatedTo)
		}
	}
	for _, c := range m.Components {
		tr.componentByID[c.ID] = c
		for _, reqID := range c.Satisfies {
			tr.satisfiedBy[reqID] = appendUnique(tr.satisfiedBy[reqID], c.ID)
		}
	}
	for _, t := range m.TestSpecs {
		tr.testByID[t.ID] = t
		for _, reqID := range t.Verifies {
			tr.verifiedBy[reqID] = appendUnique(tr.verifiedBy[reqID], t.ID)
		}
	}
	return tr
}

func (tr *Tracer) Query(q TraceQuery) TraceResult {
	switch q.Direction {
	case Backward:
		return tr.queryBackward(q.FromID)
	default:
		return tr.queryForward(q.FromID)
	}
}

func (tr *Tracer) queryForward(fromID string) TraceResult {
	if _, ok := tr.reqByID[fromID]; !ok {
		return TraceResult{Chains: []TraceChain{}}
	}
	return TraceResult{Chains: []TraceChain{tr.chainFor(fromID)}}
}

func (tr *Tracer) queryBackward(fromID string) TraceResult {
	var reqIDs []string
	if t, ok := tr.testByID[fromID]; ok {
		reqIDs = append(reqIDs, t.Verifies...)
	}
	if _, ok := tr.componentByID[fromID]; ok {
		for _, r := range tr.model.Requirements {
			if r.AllocatedTo == fromID {
				reqIDs = appendUnique(reqIDs, r.ID)
			}
		}
		for _, c := range tr.model.Components {
			if c.ID == fromID {
				for _, s := range c.Satisfies {
					reqIDs = appendUnique(reqIDs, s)
				}
			}
		}
	}

	chains := make([]TraceChain, 0, len(reqIDs))
	for _, id := range reqIDs {
		chains = append(chains, tr.chainFor(id))
	}
	return TraceResult{Chains: chains}
}

func (tr *Tracer) Matrix() TraceResult {
	chains := make([]TraceChain, 0, len(tr.model.Requirements))
	for _, r := range tr.model.Requirements {
		chains = append(chains, tr.chainFor(r.ID))
	}
	return TraceResult{Chains: chains}
}

func (tr *Tracer) AncestorChain(reqID string) []string {
	chain := []string{}
	seen := map[string]bool{}
	id := reqID
	for id != "" && !seen[id] {
		seen[id] = true
		chain = append(chain, id)
		r, ok := tr.reqByID[id]
		if !ok {
			break
		}
		if r.DerivedFrom == "" {
			break
		}
		if _, ok := tr.reqByID[r.DerivedFrom]; !ok {
			break
		}
		id = r.DerivedFrom
	}
	return chain
}

func (tr *Tracer) chainFor(reqID string) TraceChain {
	return TraceChain{
		RequirementID: reqID,
		SatisfiedBy:   nonNil(tr.satisfiedBy[reqID]),
		VerifiedBy:    nonNil(tr.verifiedBy[reqID]),
	}
}

func appendUnique(xs []string, v string) []string {
	for _, x := range xs {
		if x == v {
			return xs
		}
	}
	return append(xs, v)
}

func nonNil(xs []string) []string {
	if xs == nil {
		return []string{}
	}
	return xs
}
