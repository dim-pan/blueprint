# Blueprint Framework

Blueprint is a systems engineering framework where the **system model is the source of truth** and AI agents generate code as disposable artifacts. Full traceability from requirements to components to test specifications to generated code.

## Project Structure

```
sys/                          ← system model (human-authored, source of truth)
  requirements/*.req          ← EARS-formatted requirements
  interfaces/*.iface          ← typed data contracts
  tests/*.testspec            ← test specifications linked to requirements
  components/*.component      ← component definitions (agent-proposed, human-reviewed)

src/                          ← generated code (artifact, NOT source of truth)

example/                      ← example apps built using blueprint
  todo-app/
    sys/                      ← todo app system model
    src/                      ← todo app generated code

resources/                    ← reference materials (PDFs + extracted text)
  extracts/                   ← curated extracts for quick reference
```

## Systems Engineering References

When making design decisions about requirements, decomposition, interfaces, traceability, or verification, consult these curated extracts:

- **EARS requirements syntax**: `resources/extracts/ears_tutorial_extract.md` — all five EARS patterns, templates, examples, anti-patterns
- **NASA SE process**: `resources/extracts/nasa_se_handbook_extract.md` — requirements engineering, functional decomposition, interface management, traceability, V&V
- **SysML/MBSE concepts**: `resources/extracts/sysml_holt_extract.md` — model as source of truth, block decomposition, traceability relationships, ACRE requirements

For deeper dives, use the index files in `resources/extracts/*_index.md` to find specific line ranges in the full text files in `resources/`.

## Key Design Principles

1. **Requirements use EARS format** — WHEN/IF-THEN/WHILE/WHERE keywords. See the EARS extract.
2. **Two-level requirements** — system-level (user-facing, human-authored) and component-level (derived, agent-proposed, human-reviewed)
3. **Humans author**: system requirements, interfaces, test specs
4. **Agents author**: component decomposition, component requirements, wiring, all implementation code
5. **Traceability is mandatory** — every requirement traces to components and test specs; every test spec traces back to a requirement
6. **Code is disposable** — `src/` can be deleted and regenerated from `sys/`
7. **Blueprint dogfoods itself** — the framework's own system model lives in `sys/`

## Development Workflow

1. Write/edit system requirements in `sys/requirements/`
2. Agent proposes component decomposition → human reviews
3. Human co-writes interfaces in `sys/interfaces/`
4. Human writes test specifications in `sys/tests/`
5. Agent generates code in `src/`
6. Tests verify the generated code
7. Failures trace back through the model to identify gaps

## Implementation Language

Go — for CLI tooling, single binary distribution, and HCL ecosystem compatibility.
