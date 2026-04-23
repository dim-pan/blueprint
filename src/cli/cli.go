package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dim-pan/blueprint/src/assembler"
	"github.com/dim-pan/blueprint/src/changedetector"
	"github.com/dim-pan/blueprint/src/coverage"
	"github.com/dim-pan/blueprint/src/parser"
	"github.com/dim-pan/blueprint/src/tracer"
	"github.com/dim-pan/blueprint/src/validator"
	"github.com/dim-pan/blueprint/src/visualiser"
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
  assemble [output-dir]
                        Write one AgentContext per component as
                        JSON under output-dir (default: ./build).
  sync [output-dir]     Detect changes since the last sync and
                        write AgentContexts only for affected
                        components. Baseline is stored under
                        .blueprint/baseline.json.
  serve [sys-dir] [--port N]
                        Start a local web server (default :8080)
                        that renders the system model as an
                        interactive canvas. Ctrl+C to stop.

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
	case "assemble":
		return runAssemble(args[1:], stdout, stderr)
	case "sync":
		return runSync(args[1:], stdout, stderr)
	case "serve":
		return runServe(args[1:], stdout, stderr)
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

const defaultBuildDir = "build"
const defaultBaselinePath = ".blueprint/baseline.json"

func runAssemble(args []string, stdout, stderr io.Writer) int {
	sysDir, outDir := defaultSysDir, defaultBuildDir
	if len(args) >= 1 {
		sysDir = args[0]
	}
	if len(args) >= 2 {
		outDir = args[1]
	}
	if len(args) > 2 {
		fmt.Fprintln(stderr, "usage: blueprint assemble [sys-dir] [out-dir]")
		return 2
	}

	model, err := loadModel(sysDir, stderr)
	if err != nil {
		return 1
	}
	result := validator.Validate(model)
	if !result.Valid {
		fmt.Fprintf(stderr, "refusing to assemble: model has %d validation error(s)\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Fprintf(stderr, "  %s:%d  %s\n", e.SourceFile, e.LineNumber, e.Message)
		}
		return 1
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(stderr, "could not create output dir %q: %s\n", outDir, err)
		return 1
	}

	contexts := assembler.Assemble(model)
	for _, ctx := range contexts {
		path := filepath.Join(outDir, ctx.Component.ID+".json")
		data, err := json.MarshalIndent(ctx, "", "  ")
		if err != nil {
			fmt.Fprintf(stderr, "failed to marshal %s: %s\n", ctx.Component.ID, err)
			return 1
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			fmt.Fprintf(stderr, "failed to write %s: %s\n", path, err)
			return 1
		}
		fmt.Fprintf(stdout, "wrote %s  (%d reqs, %d test specs, %d interfaces)\n",
			path, len(ctx.Requirements), len(ctx.TestSpecs), len(ctx.Interfaces))
	}
	fmt.Fprintf(stdout, "assembled %d component context(s)\n", len(contexts))
	return 0
}

func runServe(args []string, stdout, stderr io.Writer) int {
	sysDir := defaultSysDir
	port := visualiser.DefaultPort
	i := 0
	for i < len(args) {
		a := args[i]
		switch a {
		case "--port", "-p":
			if i+1 >= len(args) {
				fmt.Fprintln(stderr, "--port requires a value")
				return 2
			}
			var p int
			if _, err := fmt.Sscanf(args[i+1], "%d", &p); err != nil || p <= 0 || p > 65535 {
				fmt.Fprintf(stderr, "invalid port: %q\n", args[i+1])
				return 2
			}
			port = p
			i += 2
		default:
			sysDir = a
			i++
		}
	}

	model, err := loadModel(sysDir, stderr)
	if err != nil {
		return 1
	}
	if result := validator.Validate(model); !result.Valid {
		fmt.Fprintf(stderr, "warning: model has %d validation error(s); serving anyway\n", len(result.Errors))
	}

	addr := fmt.Sprintf(":%d", port)
	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Fprintf(stdout, "blueprint · %s\n  %d requirements · %d components · %d interfaces · %d test specs\n  serving %s\n  press Ctrl+C to stop\n",
		sysDir, len(model.Requirements), len(model.Components), len(model.Interfaces), len(model.TestSpecs), url)

	if err := visualiser.Serve(model, addr); err != nil {
		fmt.Fprintf(stderr, "server error: %s\n", err)
		return 1
	}
	return 0
}

func runSync(args []string, stdout, stderr io.Writer) int {
	sysDir, outDir := defaultSysDir, defaultBuildDir
	if len(args) >= 1 {
		sysDir = args[0]
	}
	if len(args) >= 2 {
		outDir = args[1]
	}
	if len(args) > 2 {
		fmt.Fprintln(stderr, "usage: blueprint sync [sys-dir] [out-dir]")
		return 2
	}

	model, err := loadModel(sysDir, stderr)
	if err != nil {
		return 1
	}
	if result := validator.Validate(model); !result.Valid {
		fmt.Fprintf(stderr, "refusing to sync: model has %d validation error(s)\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Fprintf(stderr, "  %s:%d  %s\n", e.SourceFile, e.LineNumber, e.Message)
		}
		return 1
	}

	baselinePath := defaultBaselinePath
	changes, err := changedetector.Detect(model, baselinePath)
	if err != nil {
		fmt.Fprintf(stderr, "detect error: %s\n", err)
		return 1
	}

	if len(changes.Changes) == 0 {
		fmt.Fprintln(stdout, "no changes since last sync")
		return 0
	}

	fmt.Fprintf(stdout, "%d change(s):\n", len(changes.Changes))
	for _, c := range changes.Changes {
		fmt.Fprintf(stdout, "  %-10s %s  (%s)\n", c.ChangeType, c.ElementID, c.File)
	}

	if len(changes.AffectedComponents) == 0 {
		fmt.Fprintln(stdout, "no components require regeneration")
	} else {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			fmt.Fprintf(stderr, "could not create output dir %q: %s\n", outDir, err)
			return 1
		}
		contexts := assembler.AssembleOnly(model, changes.AffectedComponents)
		fmt.Fprintf(stdout, "regenerating %d component(s):\n", len(contexts))
		for _, ctx := range contexts {
			path := filepath.Join(outDir, ctx.Component.ID+".json")
			data, err := json.MarshalIndent(ctx, "", "  ")
			if err != nil {
				fmt.Fprintf(stderr, "failed to marshal %s: %s\n", ctx.Component.ID, err)
				return 1
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				fmt.Fprintf(stderr, "failed to write %s: %s\n", path, err)
				return 1
			}
			fmt.Fprintf(stdout, "  wrote %s\n", path)
		}
	}

	if err := os.MkdirAll(filepath.Dir(baselinePath), 0o755); err != nil {
		fmt.Fprintf(stderr, "could not create baseline dir: %s\n", err)
		return 1
	}
	if err := changedetector.SaveBaseline(model, baselinePath); err != nil {
		fmt.Fprintf(stderr, "failed to save baseline: %s\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "baseline saved to %s\n", baselinePath)
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
