package validator

import (
	"fmt"
	"strings"

	"github.com/dim-pan/blueprint/src/parser"
)

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

type ValidationError struct {
	Message    string
	SourceFile string
	LineNumber int
	Severity   Severity
}

type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

var primitiveTypes = map[string]struct{}{
	"String":    {},
	"Integer":   {},
	"Boolean":   {},
	"Float":     {},
	"Timestamp": {},
}

func Validate(m parser.SystemModel) ValidationResult {
	var errs []ValidationError

	reqByID := map[string][]parser.Requirement{}
	for _, r := range m.Requirements {
		reqByID[r.ID] = append(reqByID[r.ID], r)
	}
	compByID := map[string][]parser.Component{}
	for _, c := range m.Components {
		compByID[c.ID] = append(compByID[c.ID], c)
	}
	ifaceByName := map[string][]parser.InterfaceDefinition{}
	for _, i := range m.Interfaces {
		ifaceByName[i.Name] = append(ifaceByName[i.Name], i)
	}
	testByID := map[string][]parser.TestSpec{}
	for _, t := range m.TestSpecs {
		testByID[t.ID] = append(testByID[t.ID], t)
	}

	errs = append(errs, checkDuplicates(reqByID, compByID, ifaceByName, testByID)...)
	errs = append(errs, checkRequirementLinks(m, reqByID, compByID)...)
	errs = append(errs, checkComponentLinks(m, reqByID, compByID)...)
	errs = append(errs, checkTestSpecLinks(m, reqByID)...)
	errs = append(errs, checkInterfaceTypes(m, ifaceByName)...)

	return ValidationResult{Valid: len(errs) == 0, Errors: errs}
}

func checkDuplicates(
	reqs map[string][]parser.Requirement,
	comps map[string][]parser.Component,
	ifaces map[string][]parser.InterfaceDefinition,
	tests map[string][]parser.TestSpec,
) []ValidationError {
	var out []ValidationError
	for id, rs := range reqs {
		if len(rs) > 1 {
			for _, r := range rs {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("duplicate requirement id %q", id),
					SourceFile: r.SourceFile,
					LineNumber: r.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
	}
	for id, cs := range comps {
		if len(cs) > 1 {
			for _, c := range cs {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("duplicate component id %q", id),
					SourceFile: c.SourceFile,
					LineNumber: c.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
	}
	for name, is := range ifaces {
		if len(is) > 1 {
			for _, i := range is {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("duplicate interface name %q", name),
					SourceFile: i.SourceFile,
					LineNumber: i.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
	}
	for id, ts := range tests {
		if len(ts) > 1 {
			for _, t := range ts {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("duplicate test spec id %q", id),
					SourceFile: t.SourceFile,
					LineNumber: t.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
	}
	return out
}

func checkRequirementLinks(
	m parser.SystemModel,
	reqByID map[string][]parser.Requirement,
	compByID map[string][]parser.Component,
) []ValidationError {
	var out []ValidationError
	for _, r := range m.Requirements {
		if r.DerivedFrom != "" {
			parents, ok := reqByID[r.DerivedFrom]
			if !ok {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("derived_from references unknown requirement %q", r.DerivedFrom),
					SourceFile: r.SourceFile,
					LineNumber: r.LineNumber,
					Severity:   SeverityError,
				})
			} else {
				for _, p := range parents {
					if p.DerivedFrom != "" {
						out = append(out, ValidationError{
							Message: fmt.Sprintf("derived_from %q must point at a system-level requirement (got a derived requirement)", r.DerivedFrom),
							SourceFile: r.SourceFile,
							LineNumber: r.LineNumber,
							Severity:   SeverityError,
						})
						break
					}
				}
			}
		}
		if r.AllocatedTo != "" {
			if _, ok := compByID[r.AllocatedTo]; !ok {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("allocated_to references unknown component %q", r.AllocatedTo),
					SourceFile: r.SourceFile,
					LineNumber: r.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
	}
	return out
}

func checkComponentLinks(
	m parser.SystemModel,
	reqByID map[string][]parser.Requirement,
	compByID map[string][]parser.Component,
) []ValidationError {
	var out []ValidationError
	for _, c := range m.Components {
		for _, sat := range c.Satisfies {
			if _, ok := reqByID[sat]; !ok {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("satisfies references unknown requirement %q", sat),
					SourceFile: c.SourceFile,
					LineNumber: c.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
		for _, dep := range c.DependsOn {
			if _, ok := compByID[dep]; !ok {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("depends_on references unknown component %q", dep),
					SourceFile: c.SourceFile,
					LineNumber: c.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
	}
	return out
}

func checkTestSpecLinks(m parser.SystemModel, reqByID map[string][]parser.Requirement) []ValidationError {
	var out []ValidationError
	for _, t := range m.TestSpecs {
		for _, v := range t.Verifies {
			if _, ok := reqByID[v]; !ok {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("verifies references unknown requirement %q", v),
					SourceFile: t.SourceFile,
					LineNumber: t.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
	}
	return out
}

func checkInterfaceTypes(m parser.SystemModel, ifaceByName map[string][]parser.InterfaceDefinition) []ValidationError {
	var out []ValidationError
	for _, iface := range m.Interfaces {
		for _, f := range iface.Fields {
			if err := resolveType(f.Type, ifaceByName); err != "" {
				out = append(out, ValidationError{
					Message:    fmt.Sprintf("field %q on interface %q: %s", f.Name, iface.Name, err),
					SourceFile: iface.SourceFile,
					LineNumber: iface.LineNumber,
					Severity:   SeverityError,
				})
			}
		}
	}
	return out
}

func resolveType(t string, ifaces map[string][]parser.InterfaceDefinition) string {
	if strings.Contains(t, "|") {
		return ""
	}
	base := strings.TrimSuffix(t, "[]")
	if _, ok := primitiveTypes[base]; ok {
		return ""
	}
	if _, ok := ifaces[base]; ok {
		return ""
	}
	return fmt.Sprintf("unknown type %q", base)
}
