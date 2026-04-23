package changedetector

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/dim-pan/blueprint/src/parser"
)

type ChangeType string

const (
	Added    ChangeType = "added"
	Modified ChangeType = "modified"
	Removed  ChangeType = "removed"
)

type Change struct {
	File       string     `json:"file"`
	ElementID  string     `json:"element_id"`
	ChangeType ChangeType `json:"change_type"`
}

type ChangeSet struct {
	AffectedComponents []string `json:"affected_components"`
	Changes            []Change `json:"changes"`
}

func DetectBetween(baseline, current parser.SystemModel) ChangeSet {
	cs := ChangeSet{
		AffectedComponents: []string{},
		Changes:            []Change{},
	}

	cs.Changes = append(cs.Changes, diffRequirements(baseline, current)...)
	cs.Changes = append(cs.Changes, diffComponents(baseline, current)...)
	cs.Changes = append(cs.Changes, diffInterfaces(baseline, current)...)
	cs.Changes = append(cs.Changes, diffTestSpecs(baseline, current)...)

	cs.AffectedComponents = computeAffected(current, baseline, cs.Changes)
	return cs
}

func Detect(current parser.SystemModel, baselinePath string) (ChangeSet, error) {
	baseline, err := LoadBaseline(baselinePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DetectBetween(parser.SystemModel{}, current), nil
		}
		return ChangeSet{}, err
	}
	return DetectBetween(baseline, current), nil
}

func LoadBaseline(path string) (parser.SystemModel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return parser.SystemModel{}, err
	}
	var m parser.SystemModel
	if err := json.Unmarshal(data, &m); err != nil {
		return parser.SystemModel{}, err
	}
	return m, nil
}

func SaveBaseline(m parser.SystemModel, path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func diffRequirements(baseline, current parser.SystemModel) []Change {
	before := map[string]parser.Requirement{}
	for _, r := range baseline.Requirements {
		before[r.ID] = r
	}
	after := map[string]parser.Requirement{}
	for _, r := range current.Requirements {
		after[r.ID] = r
	}
	var out []Change
	for id, r := range after {
		prev, ok := before[id]
		if !ok {
			out = append(out, Change{File: r.SourceFile, ElementID: id, ChangeType: Added})
		} else if reqFingerprint(prev) != reqFingerprint(r) {
			out = append(out, Change{File: r.SourceFile, ElementID: id, ChangeType: Modified})
		}
	}
	for id, r := range before {
		if _, ok := after[id]; !ok {
			out = append(out, Change{File: r.SourceFile, ElementID: id, ChangeType: Removed})
		}
	}
	return out
}

func diffComponents(baseline, current parser.SystemModel) []Change {
	before := map[string]parser.Component{}
	for _, c := range baseline.Components {
		before[c.ID] = c
	}
	after := map[string]parser.Component{}
	for _, c := range current.Components {
		after[c.ID] = c
	}
	var out []Change
	for id, c := range after {
		prev, ok := before[id]
		if !ok {
			out = append(out, Change{File: c.SourceFile, ElementID: id, ChangeType: Added})
		} else if compFingerprint(prev) != compFingerprint(c) {
			out = append(out, Change{File: c.SourceFile, ElementID: id, ChangeType: Modified})
		}
	}
	for id, c := range before {
		if _, ok := after[id]; !ok {
			out = append(out, Change{File: c.SourceFile, ElementID: id, ChangeType: Removed})
		}
	}
	return out
}

func diffInterfaces(baseline, current parser.SystemModel) []Change {
	before := map[string]parser.InterfaceDefinition{}
	for _, i := range baseline.Interfaces {
		before[i.Name] = i
	}
	after := map[string]parser.InterfaceDefinition{}
	for _, i := range current.Interfaces {
		after[i.Name] = i
	}
	var out []Change
	for name, i := range after {
		prev, ok := before[name]
		if !ok {
			out = append(out, Change{File: i.SourceFile, ElementID: name, ChangeType: Added})
		} else if ifaceFingerprint(prev) != ifaceFingerprint(i) {
			out = append(out, Change{File: i.SourceFile, ElementID: name, ChangeType: Modified})
		}
	}
	for name, i := range before {
		if _, ok := after[name]; !ok {
			out = append(out, Change{File: i.SourceFile, ElementID: name, ChangeType: Removed})
		}
	}
	return out
}

func diffTestSpecs(baseline, current parser.SystemModel) []Change {
	before := map[string]parser.TestSpec{}
	for _, t := range baseline.TestSpecs {
		before[t.ID] = t
	}
	after := map[string]parser.TestSpec{}
	for _, t := range current.TestSpecs {
		after[t.ID] = t
	}
	var out []Change
	for id, t := range after {
		prev, ok := before[id]
		if !ok {
			out = append(out, Change{File: t.SourceFile, ElementID: id, ChangeType: Added})
		} else if testFingerprint(prev) != testFingerprint(t) {
			out = append(out, Change{File: t.SourceFile, ElementID: id, ChangeType: Modified})
		}
	}
	for id, t := range before {
		if _, ok := after[id]; !ok {
			out = append(out, Change{File: t.SourceFile, ElementID: id, ChangeType: Removed})
		}
	}
	return out
}

func reqFingerprint(r parser.Requirement) string {
	r.SourceFile = ""
	r.LineNumber = 0
	data, _ := json.Marshal(r)
	return string(data)
}

func compFingerprint(c parser.Component) string {
	c.SourceFile = ""
	c.LineNumber = 0
	data, _ := json.Marshal(c)
	return string(data)
}

func ifaceFingerprint(i parser.InterfaceDefinition) string {
	i.SourceFile = ""
	i.LineNumber = 0
	data, _ := json.Marshal(i)
	return string(data)
}

func testFingerprint(t parser.TestSpec) string {
	t.SourceFile = ""
	t.LineNumber = 0
	data, _ := json.Marshal(t)
	return string(data)
}

func computeAffected(current, baseline parser.SystemModel, changes []Change) []string {
	affected := map[string]bool{}
	changedIDs := map[string]bool{}
	interfaceChanged := false

	for _, ch := range changes {
		changedIDs[ch.ElementID] = true
	}

	currentCompIDs := map[string]bool{}
	for _, c := range current.Components {
		currentCompIDs[c.ID] = true
	}

	for _, ch := range changes {
		if currentCompIDs[ch.ElementID] {
			affected[ch.ElementID] = true
		}
	}
	for _, c := range baseline.Components {
		if !currentCompIDs[c.ID] {
			// removed components are no longer affected (they don't exist)
			continue
		}
	}

	reqToComponent := map[string]map[string]bool{}
	addReqOwner := func(reqID, compID string) {
		if reqToComponent[reqID] == nil {
			reqToComponent[reqID] = map[string]bool{}
		}
		reqToComponent[reqID][compID] = true
	}
	for _, r := range current.Requirements {
		if r.AllocatedTo != "" {
			addReqOwner(r.ID, r.AllocatedTo)
		}
	}
	for _, c := range current.Components {
		for _, sat := range c.Satisfies {
			addReqOwner(sat, c.ID)
		}
	}

	testToComponents := map[string]map[string]bool{}
	for _, t := range current.TestSpecs {
		for _, v := range t.Verifies {
			for compID := range reqToComponent[v] {
				if testToComponents[t.ID] == nil {
					testToComponents[t.ID] = map[string]bool{}
				}
				testToComponents[t.ID][compID] = true
			}
		}
	}

	for _, ch := range changes {
		switch {
		case len(reqToComponent[ch.ElementID]) > 0:
			for compID := range reqToComponent[ch.ElementID] {
				affected[compID] = true
			}
		case len(testToComponents[ch.ElementID]) > 0:
			for compID := range testToComponents[ch.ElementID] {
				affected[compID] = true
			}
		}
		if isInterfaceChange(ch.ElementID, baseline, current) {
			interfaceChanged = true
		}
	}

	if interfaceChanged {
		for id := range currentCompIDs {
			affected[id] = true
		}
	}

	for changed := true; changed; {
		changed = false
		for _, c := range current.Components {
			if affected[c.ID] {
				continue
			}
			for _, dep := range c.DependsOn {
				if affected[dep] {
					affected[c.ID] = true
					changed = true
					break
				}
			}
		}
	}

	out := make([]string, 0, len(affected))
	for _, c := range current.Components {
		if affected[c.ID] {
			out = append(out, c.ID)
		}
	}
	return out
}

func isInterfaceChange(elementID string, baseline, current parser.SystemModel) bool {
	for _, i := range current.Interfaces {
		if i.Name == elementID {
			return true
		}
	}
	for _, i := range baseline.Interfaces {
		if i.Name == elementID {
			return true
		}
	}
	return false
}
