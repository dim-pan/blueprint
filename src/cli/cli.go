package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/dim-pan/blueprint/src/coverage"
	"github.com/dim-pan/blueprint/src/parser"
	"github.com/dim-pan/blueprint/src/tracer"
	"github.com/dim-pan/blueprint/src/validator"
)

const defaultSysDir = "sys"

const usage = `blueprint - systems-engineering workflow for AI-assisted development

Usage:
  blueprint <command> [args...] [model-dir]

Commands:
  validate              Parse and validate a system model directory.
  trace <id> [forward|backward]
                        Show the traceability chain for a
                        requirement, component, or test spec.
  verify                Report test-coverage and allocation gaps.

If model-dir is omitted, the default is "sys" relative to the
current working directory.
`

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}

	switch args[0] {
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "trace":
		return runTrace(args[1:], stdout, stderr)
	case "verify":
		return runVerify(args[1:], stdout, stderr)
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}

func runValidate(args []string, stdout, stderr io.Writer) int {
	dir := defaultSysDir
	if len(args) > 0 {
		dir = args[0]
	}

	model, err := parser.ParseSystemModel(dir)
	if err != nil {
		var pe *parser.ParseError
		if errors.As(err, &pe) {
			fmt.Fprintf(stderr, "parse error: %s:%d: %s\n", pe.File, pe.Line, pe.Message)
		} else {
			fmt.Fprintf(stderr, "parse error: %s\n", err)
		}
		return 1
	}

	result := validator.Validate(model)
	if result.Valid {
		fmt.Fprintf(stdout,
			"model is valid: %d requirements, %d components, %d interfaces, %d test specs\n",
			len(model.Requirements), len(model.Components),
			len(model.Interfaces), len(model.TestSpecs))
		return 0
	}

	fmt.Fprintf(stderr, "model is invalid: %d error(s)\n", len(result.Errors))
	for _, e := range result.Errors {
		fmt.Fprintf(stderr, "  %s:%d  %s\n", e.SourceFile, e.LineNumber, e.Message)
	}
	return 1
}

func runTrace(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "trace requires an id (requirement, component, or test spec)")
		return 2
	}
	id := args[0]

	direction := tracer.Forward
	dir := defaultSysDir
	for _, a := range args[1:] {
		switch a {
		case "forward":
			direction = tracer.Forward
		case "backward":
			direction = tracer.Backward
		default:
			dir = a
		}
	}

	model, err := loadModel(dir, stderr)
	if err != nil {
		return 1
	}

	tr := tracer.New(model)
	result := tr.Query(tracer.TraceQuery{FromID: id, Direction: direction})

	if len(result.Chains) == 0 {
		fmt.Fprintf(stderr, "no chains found for %q (direction=%s)\n", id, direction)
		return 1
	}

	fmt.Fprintf(stdout, "trace (%s) for %s:\n", direction, id)
	for _, c := range result.Chains {
		fmt.Fprintf(stdout, "  %s\n", c.RequirementID)
		if len(c.SatisfiedBy) > 0 {
			fmt.Fprintf(stdout, "    satisfied_by: %s\n", joinCSV(c.SatisfiedBy))
		}
		if len(c.VerifiedBy) > 0 {
			fmt.Fprintf(stdout, "    verified_by:  %s\n", joinCSV(c.VerifiedBy))
		}
		ancestors := tr.AncestorChain(c.RequirementID)
		if len(ancestors) > 1 {
			fmt.Fprintf(stdout, "    chain:        %s\n", joinArrow(ancestors))
		}
	}
	return 0
}

func runVerify(args []string, stdout, stderr io.Writer) int {
	dir := defaultSysDir
	if len(args) > 0 {
		dir = args[0]
	}
	model, err := loadModel(dir, stderr)
	if err != nil {
		return 1
	}

	r := coverage.Analyse(model)
	fmt.Fprintf(stdout,
		"coverage: %d/%d requirements verified (%.1f%%)\n",
		r.CoveredRequirements, r.TotalRequirements, r.CoveragePercentage)
	if len(r.UncoveredRequirements) > 0 {
		fmt.Fprintf(stdout, "uncovered (%d):\n", len(r.UncoveredRequirements))
		for _, id := range r.UncoveredRequirements {
			fmt.Fprintf(stdout, "  %s\n", id)
		}
	}
	if len(r.UnallocatedRequirements) > 0 {
		fmt.Fprintf(stdout, "unallocated system requirements (%d):\n", len(r.UnallocatedRequirements))
		for _, id := range r.UnallocatedRequirements {
			fmt.Fprintf(stdout, "  %s\n", id)
		}
	}
	if len(r.UncoveredRequirements)+len(r.UnallocatedRequirements) > 0 {
		return 1
	}
	return 0
}

func loadModel(dir string, stderr io.Writer) (parser.SystemModel, error) {
	m, err := parser.ParseSystemModel(dir)
	if err != nil {
		var pe *parser.ParseError
		if errors.As(err, &pe) {
			fmt.Fprintf(stderr, "parse error: %s:%d: %s\n", pe.File, pe.Line, pe.Message)
		} else {
			fmt.Fprintf(stderr, "parse error: %s\n", err)
		}
		return parser.SystemModel{}, err
	}
	return m, nil
}

func joinCSV(xs []string) string {
	out := ""
	for i, x := range xs {
		if i > 0 {
			out += ", "
		}
		out += x
	}
	return out
}

func joinArrow(xs []string) string {
	out := ""
	for i, x := range xs {
		if i > 0 {
			out += " -> "
		}
		out += x
	}
	return out
}
