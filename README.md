# Blueprint

**A systems engineering framework for building software with AI agents — where the system model is the source of truth and code is a disposable artifact.**

Blueprint flips the usual relationship between specs and code. Instead of code being the thing you maintain and docs going stale, you maintain a structured **system model** — requirements, interfaces, components, test specs — and AI agents generate the code from it. When the model changes, you regenerate. When a test fails, you trace it back through the model to the human-authored intent that produced it.

## Why

Vibe coding with AI agents works until the codebase grows past what fits in a single context window. Then the agent loses the thread, drifts from the original intent, and you end up debugging code nobody — human or AI — fully understands.

Traditional systems engineering solved this problem decades ago for hardware and safety-critical software: capture intent in a structured model, decompose it, trace every artifact back to a requirement, verify at every level. Blueprint brings that discipline to AI-assisted software development without the heavyweight tooling.

The result:

- **No black boxes.** Every line of generated code traces back through a component, a derived requirement, and a system requirement that a human wrote and approved.
- **Regenerable code.** `src/` can be deleted and rebuilt from `sys/` at any time.
- **Bounded agent context.** Agents only see the slice of the model relevant to the component they're working on — not the whole codebase.
- **Failures point at gaps in understanding, not mystery bugs.** A failing test traces back to either a buggy implementation, an ambiguous requirement, or a missing decomposition — and you can tell which.

## Core Concepts

### The System Model

The model lives in `sys/` and is made up of four artifact types:

| File | Purpose | Authored by |
|---|---|---|
| `.req` | Requirements in [EARS](https://alistairmavin.com/ears/) format (WHEN/IF-THEN/WHILE/WHERE) | Human (system) / Agent (derived) |
| `.iface` | Typed data contracts between components | Human / Agent |
| `.component` | Component definitions with responsibilities and dependencies | Agent, human reviews |
| `.testspec` | Test specifications linked to the requirements they verify | Human / Agent |

Everything is plain text, version-controlled, diff-friendly.

### Two Levels of Requirements

- **System requirements** — user-facing, human-authored, describe what the system does
- **Component requirements** — derived by the agent from system requirements, allocated to specific components, reviewed by the human

Each derived requirement carries a `derived_from` link back to its parent. Coverage gaps are detectable by querying the model.

### Traceability Chain

```
System Requirement
  ↓ verified by
System Test Specification
  ↓ satisfied by
Component (with derived requirements)
  ↓ verified by
Component + Integration Test Specifications
  ↓ implemented as
Generated Code + Generated Tests
```

You can trace forward (requirement → code) or backward (failing test → requirement). Nothing exists in the model without a link to something else.

### Who Does What

| Role | Owns |
|---|---|
| **Human** | Direction, system requirements, interfaces, test specs, final approval on everything |
| **Agent** | Formalizing direction into EARS, proposing component decomposition, deriving component requirements, generating code, running tests, tracing failures |

The human can give input as loose as "I want a todo app" or as tight as fully-written EARS requirements. The agent fills the gaps and asks for review.

## Project Layout

```
sys/                          ← system model (source of truth)
  requirements/*.req          ← EARS-formatted requirements
  interfaces/*.iface          ← typed data contracts
  tests/*.testspec            ← test specifications linked to requirements
  components/<name>/
    <name>.component          ← component definition
    <name>.req                ← derived requirements
    <name>.testspec           ← component test specs

src/                          ← generated code (artifact)

resources/extracts/           ← curated systems engineering references
                                (EARS, NASA SE handbook, SysML/MBSE)

WORKFLOW.md                   ← step-by-step process guide
CLAUDE.md                     ← agent operating instructions
```

## Workflow at a Glance

Work in **vertical slices** — never big-design-up-front. Pick 2–3 core requirements and take them all the way down to verified code before starting the next slice.

1. **Define system requirements** — human directs, agent formalizes in EARS, human reviews
2. **Write system-level test specs** — define verification at the same time as requirements (V-model)
3. **Propose component decomposition** — agent proposes, human reviews
4. **Define interfaces** — shared data contracts at component boundaries
5. **Derive component requirements** — flow down from system requirements with traceability links
6. **Write component and integration test specs**
7. **Generate code** — agent reads the assembled context per component and generates `src/`
8. **Verify** — run tests, trace failures back through the model
9. **Sync** — when the model changes, regenerate only the affected components

The full process is documented in [WORKFLOW.md](WORKFLOW.md).

## Example: A Requirement in EARS Format

```
req REQ-001 "Create todo"
  priority: must-have
  WHEN the user submits a new to-do item with a valid title,
  the system shall add the item to the list and display it.

req REQ-002 "Empty title rejected"
  priority: must-have
  IF the user submits a to-do item with an empty title,
  THEN the system shall reject the request and display a
  validation error.
```

…and the test spec that verifies it:

```
test TC-SYS-001 "Create with valid title"
  verifies: REQ-001
  given: CreateTodoRequest{title: "Buy milk"}
  expect: TodoResult{ok: true, item.title: "Buy milk", item.completed: false}
```

## Self-Hosted

Blueprint dogfoods itself. The framework's own system model lives in `sys/`:

- **9 system requirements** covering model parsing, validation, context assembly, traceability, change detection, incremental sync, coverage analysis, and visual representation
- **7 components**: parser, validator, assembler, tracer, coverage-analyser, change-detector, cli
- **7 interfaces** defining the data contracts between them
- **30 system-level test specifications** covering every requirement

You can read the model as a worked example of how to use the framework.

## Status

Blueprint is in active design. The system model is complete; the implementation is being built from it as the first real test of the workflow. Implementation language is **Go** — single binary distribution, plays nicely with the HCL ecosystem.

## Systems Engineering References

Blueprint stands on the shoulders of decades of systems engineering practice. Curated extracts live in `resources/extracts/`:

- **EARS requirements syntax** — all five patterns with templates and anti-patterns
- **NASA Systems Engineering Handbook** — requirements engineering, decomposition, interface management, traceability, V&V
- **SysML/MBSE concepts (Holt)** — model as source of truth, block decomposition, traceability relationships

Agents consult these at each step of the workflow rather than guessing at best practice.

## License

TBD
