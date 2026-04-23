# Blueprint

### NASA doesn't vibe-code spacecraft. Why are you vibe-coding production software?

**Blueprint takes the systems engineering discipline that lands rovers on Mars — requirements, decomposition, interfaces, traceability, verification — and retools it for the way AI agents actually build software.** The result is plan mode you can commit: a structured **system model** that lives in your repo, survives past the context window, and is the thing every agent reads before it touches a single line of code.

You've used plan mode in Claude or Cursor. You've felt how good it is when the agent actually understands what you're trying to build — and how it evaporates the moment the conversation ends. Blueprint is that plan, made durable, version-controlled, diffed in PRs, and grounded in 60+ years of systems engineering practice from NASA, INCOSE, and the EARS / SysML communities.

---

## The Problems Blueprint Solves

**1. Your agent loses the plot past a few thousand lines.**
Vibe coding works on toy apps. On real codebases, the agent can't fit the system in its head, drifts from your intent, and starts confidently building the wrong thing. Blueprint bounds what each agent sees to the slice of the model relevant to one component — never the whole codebase.

**2. Specs and code drift apart the moment you ship.**
Documentation rots. Comments lie. The README describes a system that hasn't existed for six months. In Blueprint the model *is* the spec, and code is built from it — drift is structurally impossible because there's nothing to drift between.

**3. Six months later, nobody remembers why the code exists.**
Not the human, not the next agent, not you. Every component, function, and test in Blueprint traces back through a chain to a requirement a human wrote in plain language and explicitly approved. "Why is this here?" always has an answer.

**4. Agents guess instead of asking the right questions.**
Vague prompts produce plausible-looking code that misses the actual intent — and you only find out when something breaks in production. EARS-formatted requirements force ambiguity to the surface *during planning*, when it's cheap to fix, not at 2am during debugging.

**5. Refactor terror.**
Change one thing, break five others you didn't know depended on it. Blueprint's traceability + change detection tells you exactly which components a model edit affects, before you regenerate anything.

---

## The Aha Moments

### Plan mode you can commit, diff, and reuse

Plan mode is the best part of working with AI agents. Blueprint makes it a first-class artifact: version it, review it in PRs, hand it to the next agent (or the next teammate), and have them pick up exactly where you left off. The plan stops being a throwaway side-conversation and becomes the thing you maintain.

### The model is what you'd never want to lose — code is the build output

`sys/` is the artifact you protect. `src/` is what gets built from it. If you lost `src/` you'd rebuild it (component by component, slowly). If you lost `sys/` you'd have lost the actual product. This inverts how most projects treat their artifacts and it's surprisingly clarifying.

### Components are the unit of regeneration, not the codebase

Code is **not** rebuilt from scratch on every change — that would burn tokens, introduce LLM nondeterminism, and re-break things you'd already fixed. Instead, when you edit a requirement or interface, change detection identifies exactly which components are affected, and only those get regenerated. Once a component's tests pass, its code is stable until its slice of the model changes.

### You actually understand the system you built with AI

Reading your own `sys/` is faster than reading generated code, because the structure forces you to think about decomposition *before* implementation. By the time the agent writes any code, you already have a mental model of what it's building and why.

### Failing tests point at a human decision, not a black box

When something breaks, the failure traces back to one of: an implementation bug (agent fixes), an interface mismatch (review the contract), a decomposition gap (the components don't fully cover the requirement), or an ambiguous requirement (human clarifies). You always know which.

---

## How It Works

Blueprint borrows directly from the **NASA Systems Engineering Handbook**, the **EARS** requirements notation (Mavin et al., used widely in aerospace and automotive), and **SysML/MBSE** modelling practice. These aren't decorative references — agents are instructed to consult them at every step of the workflow, so the model you end up with looks like one a systems engineer would recognize, not a pile of bullet points an LLM hallucinated.

The model lives in `sys/` and is made up of four artifact types — all plain text, version-controlled, diff-friendly:

| File | Purpose | Authored by |
|---|---|---|
| `.req` | Requirements in [EARS](https://alistairmavin.com/ears/) format (WHEN/IF-THEN/WHILE/WHERE) | Human (system) / Agent (derived) |
| `.iface` | Typed data contracts between components | Human / Agent |
| `.component` | Component definitions with responsibilities and dependencies | Agent, human reviews |
| `.testspec` | Test specifications linked to the requirements they verify | Human / Agent |

### Two levels of requirements

- **System requirements** — user-facing, human-authored, describe what the system does
- **Component requirements** — derived by the agent from system requirements, allocated to specific components, reviewed by the human

Each derived requirement carries a `derived_from` link back to its parent. Coverage gaps are detectable by querying the model.

### The traceability chain

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

You can trace forward (requirement → code) or backward (failing test → requirement).

### Who does what

| Role | Owns |
|---|---|
| **Human** | Direction, system requirements, interfaces, test specs, final approval on everything |
| **Agent** | Formalizing direction into EARS, proposing component decomposition, deriving component requirements, generating code, running tests, tracing failures |

The human can give input as loose as "I want a todo app" or as tight as fully-written EARS requirements. The agent fills the gaps and asks for review.

---

## A Requirement in EARS Format

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

That's the whole vocabulary. No new framework to learn — just structured English with linked IDs.

---

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

src/                          ← generated code (built from sys/, regenerated per component)

resources/extracts/           ← curated systems engineering references
                                (EARS, NASA SE handbook, SysML/MBSE)

WORKFLOW.md                   ← step-by-step process guide
CLAUDE.md                     ← agent operating instructions
```

---

## Workflow at a Glance

Blueprint follows the **V-model** from systems engineering: every requirement gets its verification defined at the same level and same time, before any design work begins. Then you work in **vertical slices** — never big-design-up-front. Pick 2–3 core requirements and take them all the way down to verified code before starting the next slice.

1. **Define system requirements** — human directs, agent formalizes in EARS, human reviews
2. **Write system-level test specs** — define verification at the same time as requirements (V-model)
3. **Propose component decomposition** — agent proposes, human reviews
4. **Define interfaces** — shared data contracts at component boundaries
5. **Derive component requirements** — flow down from system requirements with traceability links
6. **Write component and integration test specs**
7. **Generate code** — agent reads the assembled context per component and writes `src/`
8. **Verify** — run tests, trace failures back through the model
9. **Sync** — when the model changes, regenerate only the affected components

The full process is in [WORKFLOW.md](WORKFLOW.md).

---

## Self-Hosted

Blueprint dogfoods itself. The framework's own system model lives in `sys/`:

- **9 system requirements** — model parsing, validation, context assembly, traceability, change detection, incremental sync, coverage analysis, visual representation
- **7 components** — parser, validator, assembler, tracer, coverage-analyser, change-detector, cli
- **7 interfaces** defining the data contracts between them
- **30 system-level test specifications** covering every requirement

You can read the model as a worked example of how to use the framework.

---

## Status

Blueprint is in active design. The system model is complete; the implementation is being built from it as the first real test of the workflow. Implementation language is **Go** — single binary distribution, plays nicely with the HCL ecosystem.

---

## Standing on the Shoulders of Systems Engineering

Blueprint isn't inventing a new methodology. It's taking the parts of systems engineering that have worked for 60+ years on safety-critical hardware and software and adapting them to a workflow where AI agents do the implementation. Curated extracts of the source material live in `resources/extracts/` and are loaded into agent context at the relevant step:

- **NASA Systems Engineering Handbook (NASA/SP-2016-6105)** — requirements engineering, functional decomposition, interface management, traceability, verification & validation. The backbone of the workflow.
- **EARS — Easy Approach to Requirements Syntax (Mavin et al.)** — the five constrained-natural-language patterns (Ubiquitous, Event-driven, Unwanted, State-driven, Optional) used throughout aerospace and automotive. Forces requirements to be unambiguous and testable.
- **SysML / MBSE concepts (Holt & Perry)** — model as the single source of truth, block decomposition, ports and interface blocks, the satisfy/verify/derive/refine traceability relationships.

When the agent makes a decision about how to phrase a requirement, where to draw a component boundary, or what belongs in an interface, it's consulting these sources rather than guessing. That's the difference between "AI wrote a spec" and "AI wrote a spec a systems engineer would sign off on."

---

## License

TBD
