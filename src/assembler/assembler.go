package assembler

import (
	"github.com/dim-pan/blueprint/src/parser"
)

type AgentContext struct {
	Component    parser.Component              `json:"component"`
	Interfaces   []parser.InterfaceDefinition  `json:"interfaces"`
	Requirements []parser.Requirement          `json:"requirements"`
	TestSpecs    []parser.TestSpec             `json:"test_specs"`
}

func Assemble(m parser.SystemModel) []AgentContext {
	out := make([]AgentContext, 0, len(m.Components))
	for _, c := range m.Components {
		out = append(out, assembleOne(m, c))
	}
	return out
}

func AssembleOnly(m parser.SystemModel, componentIDs []string) []AgentContext {
	wanted := map[string]bool{}
	for _, id := range componentIDs {
		wanted[id] = true
	}
	out := make([]AgentContext, 0, len(wanted))
	for _, c := range m.Components {
		if wanted[c.ID] {
			out = append(out, assembleOne(m, c))
		}
	}
	return out
}

func assembleOne(m parser.SystemModel, c parser.Component) AgentContext {
	ctx := AgentContext{
		Component:    c,
		Interfaces:   append([]parser.InterfaceDefinition{}, m.Interfaces...),
		Requirements: []parser.Requirement{},
		TestSpecs:    []parser.TestSpec{},
	}

	reqByID := map[string]parser.Requirement{}
	for _, r := range m.Requirements {
		reqByID[r.ID] = r
	}

	included := map[string]bool{}
	appendReq := func(id string) {
		if included[id] {
			return
		}
		if r, ok := reqByID[id]; ok {
			ctx.Requirements = append(ctx.Requirements, r)
			included[id] = true
		}
	}

	for _, r := range m.Requirements {
		if r.AllocatedTo == c.ID {
			appendReq(r.ID)
			if r.DerivedFrom != "" {
				appendReq(r.DerivedFrom)
			}
		}
	}
	for _, sat := range c.Satisfies {
		appendReq(sat)
	}

	relevantReqs := included
	for _, t := range m.TestSpecs {
		for _, v := range t.Verifies {
			if relevantReqs[v] {
				ctx.TestSpecs = append(ctx.TestSpecs, t)
				break
			}
		}
	}

	return ctx
}
