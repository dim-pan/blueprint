# Blueprint Workflow Guide

This guide describes the step-by-step process for building software using the Blueprint framework. The system model is the source of truth. Code is a generated artifact.

The agent follows this workflow top-down, consulting the systems engineering references in `resources/extracts/` at every step. The references provide guidance on EARS requirements syntax, NASA SE process (decomposition, interfaces, traceability, V&V), and SysML/MBSE modelling concepts.

## Who Does What

| Role | Responsibility |
|---|---|
| **Human** | Gives direction, reviews and approves all model artifacts, has final say on requirements, interfaces, and test specs |
| **Agent** | Helps formalize direction into EARS requirements, proposes decomposition, derives component requirements, generates code, runs tests, traces failures back to the model |

The human can provide as much or as little detail as they want — from a one-line direction ("I want a todo app") to fully written EARS requirements. The agent fills in the gaps, always consulting the systems engineering references for best practice.

## Iterative Slices, Not Big Design Up Front

The steps below describe the full top-to-bottom process. You do **not** need to complete every step for the entire system before moving forward. Instead, work in **vertical slices**:

1. Pick 2–3 core requirements
2. Follow Steps 1–8 all the way down to working, verified code for those requirements
3. Pick the next slice of requirements and repeat

Each slice follows the V-model structure (requirements → test specs → decomposition → interfaces → derived requirements → component tests → code → verify), but scoped to a subset of the system. This gives you:

- **Working software early** — you don't wait for a complete spec to start building
- **Fast feedback** — each slice reveals gaps in the model while the context is fresh
- **Full traceability at every point** — even a partial system has complete traceability for the requirements that exist
- **Safe iteration** — adding new requirements in later slices triggers the sync process (Step 9), which identifies affected components and regenerates only what changed

The key discipline: **never skip levels within a slice**. Every requirement in a slice gets test specs, component allocation, derived requirements, and verification before moving to the next slice. The traceability chain must be complete for each slice, even if the system isn't complete yet.

## Prerequisites

Your project needs a `sys/` directory at the root. This is the system model.

```
your-project/
  sys/
    requirements/          ← system-level requirements
    interfaces/            ← shared data contracts
    tests/                 ← system-level and integration test specs
    components/            ← one folder per component
      component-name/
        name.component     ← component definition
        name.req           ← derived requirements (agent-generated)
        name.testspec      ← component test specs (agent-generated)
  src/                     ← generated code (artifact, not source of truth)
```

---

## Step 1: Define System Requirements

**Who:** Human gives direction, agent formalizes, human reviews

The human describes what they want. The agent then:

1. Consults the EARS extract (`resources/extracts/ears_tutorial_extract.md`) for pattern selection
2. Consults the NASA SE extract (`resources/extracts/nasa_se_handbook_extract.md`) for requirement writing rules
3. Identifies the distinct behavioral requirements in the direction
4. Splits compound statements into atomic requirements (one behavior per requirement)
5. Rewrites each in EARS format using the appropriate pattern
6. Assigns IDs and priorities
7. Presents the draft for human review

The human reviews, edits, and approves. The approved requirements are written to `sys/requirements/<domain>.req`.

### EARS Patterns

| Pattern | Template | Use when |
|---|---|---|
| Ubiquitous | The \<system\> shall \<response\>. | Always true, no trigger needed |
| Event-driven | WHEN \<trigger\>, the \<system\> shall \<response\>. | Something happens |
| Unwanted | IF \<condition\>, THEN the \<system\> shall \<response\>. | Error/edge case handling |
| State-driven | WHILE \<state\>, the \<system\> shall \<response\>. | Behavior during a state |
| Optional | WHERE \<feature\>, the \<system\> shall \<response\>. | Feature-dependent behavior |

### Requirement Quality Rules

From NASA SE Handbook (Appendix C) and EARS best practice:

- One requirement per statement — do not combine multiple behaviors
- Each requirement gets a unique ID (e.g., `REQ-001`)
- Use `shall` for binding requirements, `should` for goals, `will` for facts
- Avoid vague terms: "fast", "user-friendly", "flexible", "robust", "efficient", "adequate"
- Every requirement must be verifiable — if you can't test it, rewrite it
- Set priority: `must-have`, `should-have`, or `nice-to-have`
- Include acceptance criteria where the requirement has quantitative aspects
- Requirements state WHAT, not HOW — no implementation details
- Each requirement must be traceable bidirectionally (to parent need and to future components/tests)

### Example

Human says: "I want a todo app where I can add and complete tasks."

Agent formalizes:

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

req REQ-003 "Complete todo"
  priority: must-have
  WHEN the user marks a to-do item as complete,
  the system shall update its status to complete and keep
  it visible in the list.
```

Human reviews: "Looks good, but I also need persistence." Agent adds the requirement and re-presents.

---

## Step 2: Write System-Level Test Specifications

**Who:** Agent proposes, human reviews

Immediately after requirements are approved, define how each will be verified. This follows the V-model: verification is defined at the same level and same time as requirements, before any design work begins.

The agent consults the NASA SE extract on verification methods and proposes test specs. The human reviews and approves.

### Rules

- Every system requirement must have at least one test specification
- `verifies` links to the requirement ID being tested
- Use `given`/`when`/`expect` clauses — structured enough for agents to generate test code from
- Test specs are written to `sys/tests/<domain>.testspec`
- System-level tests treat the system as a black box — no knowledge of components

### Example

```
test TC-SYS-001 "Create with valid title"
  verifies: REQ-001
  given: CreateTodoRequest{title: "Buy milk"}
  expect: TodoResult{ok: true, item.title: "Buy milk", item.completed: false}

test TC-SYS-002 "Create with empty title fails"
  verifies: REQ-002
  given: CreateTodoRequest{title: ""}
  expect: TodoResult{ok: false, error: "title cannot be empty"}

test TC-SYS-003 "Complete a todo"
  verifies: REQ-003
  given: existing Todo{id: "1", completed: false}
  when: UpdateTodoRequest{id: "1", completed: true}
  expect: TodoResult{ok: true, item.completed: true}
```

### Review Checklist (Human)

- Does every system requirement have at least one test spec?
- Are the test specs concrete enough to generate test code from?
- Do they cover both happy paths and edge cases?

---

## Step 3: Propose Component Decomposition

**Who:** Agent proposes, human reviews

The agent reads all system requirements and consults the NASA SE extract on functional decomposition and the SysML extract on block decomposition. It proposes a set of components that together satisfy all requirements.

Each component gets:

- A unique ID (e.g., `COMP-001`)
- A human-readable name
- A clear responsibility (one sentence describing what it does)
- `satisfies` links to the system requirements it addresses
- `depends_on` links to other components it needs (by ID)

### Rules

- Each system requirement must be satisfied by at least one component
- Each component should have a single, focused responsibility
- Dependencies should flow in one direction — avoid circular dependencies
- Write components to `sys/components/<name>/<name>.component`

### Example

```
component COMP-001 "TodoService"
  responsibility: Enforce business rules for todo creation,
                  completion, and deletion.
  satisfies: REQ-001, REQ-002, REQ-003
  depends_on: COMP-002

component COMP-002 "TodoStore"
  responsibility: Persist and retrieve todo items.
  satisfies: REQ-004
```

### Review Checklist (Human)

- Does every system requirement have at least one component satisfying it?
- Are responsibilities clear and non-overlapping?
- Are the dependencies reasonable?
- Is the decomposition at the right granularity — not too coarse, not too fine?

---

## Step 4: Define Interfaces

**Who:** Agent proposes, human reviews

Define the data contracts that flow between components. The agent consults the NASA SE extract on interface management and the SysML extract on ports and interface blocks.

Interfaces are shared — they belong to the boundary, not to either component.

### Rules

- One interface per logical data type
- Every field has a name, type, and optionality
- Use union types for enums (e.g., `ok | failed | pending`)
- Use `?` for optional fields
- Use `[]` for arrays
- Write interfaces to `sys/interfaces/<name>.iface`

### Example

```
interface Todo {
  id:         String
  title:      String
  completed:  Boolean
  created_at: Timestamp
}

interface CreateTodoRequest {
  title: String
}

interface TodoResult {
  ok:    Boolean
  item:  Todo?
  error: String?
}
```

---

## Step 5: Derive Component Requirements

**Who:** Agent proposes, human reviews

For each component, the agent derives concrete requirements from the system requirements it satisfies. The agent consults the NASA SE extract on requirements flowdown — understanding the difference between allocated requirements (flowed down from parent) and derived requirements (arising from design decisions).

### Rules

- Every derived requirement must link back via `derived_from` to a system requirement
- Every derived requirement must link via `allocated_to` to a component ID
- Use EARS format for the statement
- Write to `sys/components/<name>/<name>.req`

### Example

```
req REQ-001-SVC-01 "Validate title not empty"
  derived_from: REQ-001
  allocated_to: COMP-001
  priority: must-have
  IF a CreateTodoRequest has an empty title,
  THEN TodoService shall return TodoResult with ok=false
  and error="title cannot be empty".
```

### Review Checklist (Human)

- Does the full set of derived requirements completely cover the parent system requirement?
- Are any derived requirements contradictory?
- Are they all verifiable?

---

## Step 6: Write Component and Integration Test Specifications

**Who:** Agent derives, human reviews

The agent derives test specifications for each component (unit level) and for component interactions (integration level).

### Three Levels

| Level | What it tests | Written by | Location |
|---|---|---|---|
| System | End-to-end user-facing behavior | Step 2 (already done) | `sys/tests/<name>.testspec` |
| Integration | Components working together at interfaces | Agent | `sys/tests/integration.testspec` |
| Component | Individual component against derived reqs | Agent | `sys/components/<name>/<name>.testspec` |

### Integration Tests

Integration tests verify that components work correctly together at their boundaries. They reference two or more components and test the interface contracts between them.

```
test TC-INT-001 "Parser output feeds Validator"
  verifies: BP-002-V01
  components: COMP-001, COMP-002
  given: a sys/ directory with a valid .req and .component file
  when: Parser output is passed to Validator
  expect: ValidationResult{valid: true, errors: []}
```

### Component Tests

Component tests verify individual derived requirements in isolation.

```
test TC-P01 "Parse requirement with all fields"
  verifies: BP-001-P01
  given: a .req file containing one valid requirement
  expect: Requirement with id, title, priority, pattern, and statement populated
```

### Rules

- Every derived requirement must have at least one component test spec
- Integration tests cover the interfaces between components
- `verifies` links to the requirement ID being tested
- Component test specs go in `sys/components/<name>/<name>.testspec`
- Integration test specs go in `sys/tests/integration.testspec`

---

## Step 7: Generate Code

**Who:** Agent

The agent reads the assembled context for each component and generates the implementation code.

### What the Agent Receives (per component)

- The component definition (responsibility, dependencies)
- All derived requirements allocated to it
- The interfaces it consumes and produces
- The test specifications it must satisfy
- The parent system requirements (for intent)

### Rules

- Generated code goes in `src/<component-name>/`
- One directory per component — mirrors the system model structure
- The agent must generate test code from the test specifications
- Generated code must conform to the interfaces exactly — no deviations

### Structural Traceability

```
sys/components/todo-service/     →     src/todo-service/
sys/components/todo-store/       →     src/todo-store/
```

Every component in the model maps to a directory in `src/`. This mapping is the traceability link between model and code.

---

## Step 8: Verify

**Who:** Agent runs, human reviews failures

Run all generated tests at every level. Trace failures back through the model.

### Failure Diagnosis

| Symptom | Likely Cause |
|---|---|
| Component test fails | Implementation bug — agent fixes the code |
| Integration test fails | Interface mismatch — review the interface definition |
| System test fails, all component tests pass | Decomposition gap — the component requirements don't fully cover the system requirement |
| Test spec is untestable | Requirement is ambiguous — human clarifies |

Every failure traces back to a human-authored artifact. Nothing disappears into a black box.

---

## Step 9: Sync (Ongoing)

**Who:** Human edits model, agent syncs code

After the initial generation, the model evolves. When requirements or interfaces change:

1. Identify what changed in the model
2. Determine which components are affected (direct changes + transitive dependencies)
3. Regenerate only the affected components
4. Re-run tests for affected components

This is not a full regeneration — it is a targeted sync of the delta.

---

## Traceability at Every Level

The complete chain:

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
  ↓ produces
Test Result (pass/fail)
```

Every element links to every other. You can trace forward (requirement → code) or backward (failing test → requirement). Gaps in coverage are visible by querying the model.

---

## Quick Reference: File Types

| Extension | Purpose | Authored by |
|---|---|---|
| `.req` | Requirements (system or component level) | Human / Agent |
| `.iface` | Interface definitions (data contracts) | Human / Agent |
| `.component` | Component definitions | Agent (human reviews) |
| `.testspec` | Test specifications | Human / Agent |

---

## Systems Engineering References

The agent should consult these extracts when making decisions at each step:

| Step | Primary Reference |
|---|---|
| Requirements | `resources/extracts/ears_tutorial_extract.md` — EARS patterns and quality rules |
| Requirements | `resources/extracts/nasa_se_handbook_extract.md` — requirement writing, validation, metadata |
| Decomposition | `resources/extracts/nasa_se_handbook_extract.md` — functional decomposition, system architecture |
| Decomposition | `resources/extracts/sysml_holt_extract.md` — block decomposition, composition relationships |
| Interfaces | `resources/extracts/nasa_se_handbook_extract.md` — interface management |
| Interfaces | `resources/extracts/sysml_holt_extract.md` — ports, interface blocks, item flows |
| Traceability | `resources/extracts/sysml_holt_extract.md` — satisfy, verify, derive, refine relationships |
| Verification | `resources/extracts/nasa_se_handbook_extract.md` — V&V methods, test planning |

For deeper dives, use the index files (`*_index.md`) to find specific line ranges in the full text files.
