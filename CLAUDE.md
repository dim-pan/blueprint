# Blueprint Change Protocol

`sys/` is the source of truth. Code in `src/` is a generated artifact. Every
behavior change flows through the model first. See `WORKFLOW.md` for the full
V-model.

## Protocol (mandatory for every behavior change)

1. **Read the model.** Find the files in `sys/` that govern the change:
   - `sys/requirements/*.req` — system requirements
   - `sys/interfaces/*.iface` — data contracts
   - `sys/components/<c>/<c>.component` — component definition
   - `sys/components/<c>/<c>.req` — derived requirements
   - `sys/components/<c>/<c>.testspec` — component test specs
   - `sys/tests/*.testspec` — system / integration test specs

2. **Update the model.** Edit the relevant files first. Per WORKFLOW.md:
   new behavior → new `req` in EARS form; new contract → `.iface` edit;
   new verifiable behavior → new `test` in the matching `.testspec`.

3. **Validate.** `./blueprint validate sys` — fix every error before moving on.

4. **Sync.** `./blueprint sync sys build` — regenerates AgentContext JSON for
   affected components and updates `.blueprint/baseline.json`.

5. **Implement.** Edit code in `src/<component>/`. Keep one Go test function
   per test-case ID (`TestTC_<TC-ID>_<slug>`) for 1:1 traceability.

6. **Test.** `go test ./...`. If the Vite SPA changed, also run
   `npm run build` in `src/visualiser/web` and rebuild the binary with
   `go build -o blueprint ./cmd/blueprint` so the embedded dist is current.

7. **Confirm.** Report what changed in `sys/`, what changed in `src/`, and
   the test outcome.

## When to skip the model step

- Pure refactors that do not change behavior: skip steps 2–4.
- Bug fixes: skip step 2 only if the existing req already unambiguously
  describes the correct behavior. Otherwise update the testspec first so the
  bug becomes a failing test.

## Conventions

- Lightweight custom parser (no HCL).
- One Go test function per TC in the testspec — not table-driven.
- Stdlib-first Go. The visualiser frontend is the only npm tree.

Consult `WORKFLOW.md` for EARS patterns, decomposition guidance, and review
checklists.
