# SysML for Systems Engineering (Holt & Perry) — Line Index

Source file: `/Users/dimitriospanagiotopoulos/Raylight_Dev/blueprint/resources/sysml_for_se_holt.txt`
Total lines: 62,190

This index maps topics to approximate line ranges in the source file. Line numbers reference the
raw text file. Use these as offsets when calling `Read` with `offset` and `limit` parameters.

---

## Front Matter and Table of Contents

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Book title, publication info | 1 | 76 | ISBN, copyright, edition details |
| Table of contents (chapter list) | 77 | 975 | Full TOC with page numbers |
| Author biographies | 977 | 998 | |
| Preface to Third Edition | 1002 | 1023 | Summary of changes from 2nd ed. |

---

## Part I — Introduction

### Chapter 1: Introduction to Model-Based Systems Engineering

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Ch1 intro — why MBSE | 1084 | 1123 | Three evils: complexity, understanding, communication |
| Systems engineering definitions (Ramo, Eisner, INCOSE) | 1162 | 1228 | Four definitions compared |
| MBSE definition — INCOSE and book's own definition | 1232 | 1337 | Model as single source of truth (line 1148) |
| Benefits of effective MBSE (bullet list) | 1364 | 1415 | Automatic generation, traceability, consistency, etc. |
| Common language — spoken vs. domain-specific | 1417 | 1486 | SysML as notation; Ontology as domain language |
| Book structure overview | 1514 | 1598 | Six parts described |

### Chapter 2: Approach

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| MBSE Mantra — People, Process, Tools | 1743 | 1791 | Figure 2.1 discussion |
| MBSE Fundamentals — Goal, Approach, Visualisation | 1797 | 2044 | Core principle: Model abstracts System |
| View vs. Diagram distinction | 1878 | 1901 | Critical: View is semantic, Diagram is visual |
| MBSE Ontology concept | 2068 | 2120 | What the Ontology is used for |
| MBSE Framework concept | 2221 | 2275 | Coverage, rigour, automation, flexibility |
| Framework vs. Pattern | 2293 | 2310 | Pattern = reusable Framework |
| View and Viewpoint definitions | 2312 | 2437 | Viewpoint is template; View is realisation |
| Rules constraining the Framework | 2438 | 2499 | Rule = constraint on Viewpoint/element |
| MBSE Meta-model (full combined diagram) | 2565 | 2595 | Figure 2.12 |
| Chapter summary | 2574 | 2611 | Three main considerations |

### Chapter 3: MBSE Concepts

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Ch3 intro — why concepts matter | 2616 | 2641 | |
| Sources used for the Ontology (ISO, INCOSE, CMMI, etc.) | 2645 | 2744 | |
| The MBSE Ontology — high-level diagram | 2829 | 3145 | Figure 3.1 |
| System concept — ISO 15288 definition | 3192 | 3290 | System, System Element, System of Interest, Enabling System |
| Need/Requirement concepts — multiple sources | 3600 | 3999 | ISO, INCOSE, model-based RE book |
| MBSE Ontology — Need-related section | 3824 | 4001 | Figure 3.7; definitions of Need, Goal, Capability, Requirement, Use Case, Scenario |
| Architecture concept | 4003 | 4200 | ISO 42010 definitions |
| Life Cycle concept | 4200 | 4500 | Stages, Life Cycle Models |
| Process concept | 4500 | 4900 | Process, Activity, Artefact, Stakeholder Role |
| Competence concept | 4900 | 5300 | Person, Competency, Stakeholder Role |
| Project concept | 5300 | 5700 | Project, Programme, Schedule |

---

## Part II — Modelling

### Chapter 4: Introduction to SysML and Systems Modelling

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Why we model (kennel/house/office block analogy) | 6900 | 7050 | Complexity through relationships — Figure 4.4 |
| The Brontosaurus of Complexity | 7055 | 7099 | Complexity grows then shrinks through life cycle |
| Lack of understanding (cascade effect) | 7101 | 7138 | Requirements mistakes cost more to fix later |
| Communication failures | 7140 | 7207 | Tower of Babel; levels of communication failure |
| What is SysML? — OMG definition | 7209 | 7350 | SysML vs. UML relationship |
| Defining modelling — model as simplification | 7354 | 7411 | Mathematical, physical, visual, text models |
| Choice of model; level of abstraction | 7412 | 7491 | Four requirements for a modelling language |
| Independent views of the same system | 7493 | 7600 | Multiple viewpoints |
| Structural modelling introduction | 7600 | 7900 | Blocks and relationships introduced |
| Behavioural modelling introduction | 7900 | 8200 | State machines, activity diagrams |
| Levels of abstraction in systems | 8200 | 8500 | System → subsystem → component |

### Chapter 5: The SysML Notation

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Diagram structure and frames | 9509 | 9556 | Every diagram has same underlying structure |
| Stereotypes — definition and usage | 9558 | 9700 | How to extend SysML; tags |
| SysML meta-model overview | 9702 | 9717 | SysML defined in terms of UML |
| **Block Definition Diagrams** | | | |
| BDD — meta-model, blocks and relationships | 9722 | 9840 | Figure 5.10 |
| BDD — Ports (Full Port, Proxy Port) | 9853 | 9875 | Port types and their meaning |
| BDD — Property types (Part, Reference, Value, Flow) | 9875 | 9966 | Figure 5.12 examples |
| BDD — Relationship types (Association, Aggregation, Composition, Generalisation, Dependency) | 10070 | 10200 | Figure 5.13 |
| BDD — Interface Blocks and interfaces | 11350 | 11530 | Required/provided interfaces; Figures 5.30–5.32 |
| **Internal Block Diagrams** | | | |
| IBD — meta-model | 11050 | 11200 | Figure 5.22 |
| IBD — ports, connectors, item flows | 11200 | 11530 | Figures 5.23–5.34 |
| IBD — nested ports | 11587 | 11626 | Figure 5.35 |
| IBD — summary | 11628 | 11646 | |
| **Package Diagrams** | | | |
| Package diagrams | 11648 | 11879 | Structuring a model; import/access |
| **Parametric Diagrams** | | | |
| Parametric diagrams — constraint blocks | 11893 | 12082 | Figures 5.42–5.47 |
| Parametric diagrams — examples (go/no-go decision) | 12400 | 12882 | Figures 5.48–5.52 |
| **Requirement Diagrams** — MOST RELEVANT | | | |
| Requirement diagram — overview | 12884 | 12958 | Figure 5.53 meta-model |
| Requirement diagram — notation | 12959 | 13084 | Figure 5.54; all six relationship types |
| Satisfy relationship — definition | 13044 | 13048 | Design element satisfies requirement |
| Trace relationship — definition | 13049 | 13053 | General traceability |
| Refine relationship — definition | 13054 | 13056 | Adds detail or precision |
| Verify relationship — definition | 13057 | 13064 | Test case verifies requirement |
| Derive relationship — definition | 13029 | 13033 | Derived requirement from stated requirement |
| Nesting — decomposition | 13019 | 13028 | Requirement broken into sub-requirements |
| Requirement diagram — examples | 13097 | 13360 | Figures 5.55–5.59; trace, derive, refine examples |
| Requirement diagram — summary | 13356 | 13371 | Use most specific relationship type |
| **State Machine Diagrams** | | | |
| State machines | 13377 | 13800 | States, transitions, events, composite states |
| **Sequence Diagrams** | | | |
| Sequence diagrams | 14400 | 14900 | Interactions, lifelines, messages |
| **Activity Diagrams** | | | |
| Activity diagrams | 15000 | 15570 | Activities, flows, decision nodes |
| **Use Case Diagrams** | | | |
| Use case diagrams | 14990 | 15570 | Actors, use cases, include/extend; ACRE connection |
| Use case diagram — summary | 15568 | 15580 | Central to ACRE |

### Chapter 6: Diagramming Guidelines

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Naming conventions for structural/behavioural diagrams | 15600 | 15800 | Best practice naming |
| Diagram frame labels | 15800 | 15900 | How to label frames |
| Interface and item flow conventions | 15900 | 16050 | Showing interfaces; item flows |

---

## Part III — Applications

### Chapter 7: Process Modelling with MBSE

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Process modelling overview — MBSE Ontology for processes | 16600 | 16900 | |
| The seven views approach — framework | 16900 | 17200 | Seven viewpoints for process modelling |
| Process Structure Viewpoint | 17200 | 17500 | |
| Process Content Viewpoint | 17500 | 18000 | Processes, Activities, Artefacts |
| Stakeholder Viewpoint | 18000 | 18300 | |
| Requirement Context Viewpoint (for processes) | 18300 | 18600 | Use cases for the process |
| Process Behaviour Viewpoint | 18600 | 19000 | Activity diagrams for process behaviour |
| Information Viewpoint | 19000 | 19083 | Artefacts and relationships; process automation basis |
| Process Instance Viewpoint | 19085 | 19290 | Validates processes against use cases |
| Complete Process Modelling Framework | 19291 | 19387 | Figure 7.25 |
| Using the framework — various scenarios | 19388 | 19699 | Analysing, creating, improving processes |

### Chapter 8: Expanded Process Modelling

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Standards modelling | 19700 | 20200 | Compliance mapping |
| Competency frameworks | 20200 | 22200 | Framework, Applicable Competency, Scope, Profile Views |

### Chapter 9: Requirements Modelling with MBSE (ACRE) — MOST RELEVANT

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Ch9 introduction — ACRE overview | 24086 | 24220 | ACRE = Approach to Context-based Requirements Engineering |
| ACRE Framework Context (use case diagram) | 24095 | 24220 | Figure 9.1; needs for the approach |
| MBSE Ontology — Need-related elements | 24221 | 24340 | Figure 9.2 |
| Need concept — four types | 24341 | 24467 | Goal, Capability, Requirement, Concern |
| Requirement taxonomy (Functional, Non-functional, Business) | 24445 | 24467 | |
| Need Description concept | 24468 | 24493 | Description vs. Need itself |
| Source Element concept | 24495 | 24507 | Any origin is valid; must be configurable |
| Rule concept | 24509 | 24553 | Constraining Need Descriptions; automatable |
| Scenario concept (two types) | 24555 | 24623 | Semi-formal (sequence diagrams); Formal (parametric) |
| Context concept — multiple types | 24625 | 24756 | System Context; Stakeholder Context |
| Use Case concept — detailed | 24758 | 24820 | Need in Context; examples (airline, car) |
| ACRE Framework — six Viewpoints overview | 24822 | 24924 | Figure 9.6 |
| Source Element Viewpoint — detail | 24934 | 25243 | Traceability; configurable items; Figures 9.7–9.10 |
| Requirement Description Viewpoint — detail | 25244 | 25855 | Attributes; flat structure principle; Figures 9.11–9.14 |
| Definition Rule Set Viewpoint — detail | 25857 | 26295 | Rules forms; Flesch-Kincaid example; Figures 9.15–9.20 |
| Requirement Context Viewpoint — detail | 26296 | 27000 | Use cases; context diagrams; Figures 9.21–9.22 |
| Validation Viewpoint — detail | 27000 | 27500 | Semi-formal and formal scenarios |
| Traceability Viewpoint | 27500 | 27900 | Cross-view traceability links |
| Complete ACRE Framework | 27900 | 28200 | Full framework diagram |
| The ACRE Process | 28200 | 28800 | Process for applying the framework |

### Chapter 10: Expanded Requirements Modelling — Systems of Systems

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Systems of Systems types | 28800 | 29200 | Directed, acknowledged, collaborative, virtual |
| SoS requirements approach | 29200 | 29800 | Multi-context requirements |

### Chapter 11: Architectures and Architectural Frameworks

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Architecture concepts | 29800 | 30500 | Architecture, architecture description, viewpoints |
| Architectural Frameworks (DoDAF, MODAF, NAF, TOGAF) | 30500 | 31500 | |
| Framework for Architectural Frameworks (FAF) | 31500 | 32500 | Meta-framework |

---

## Part IV — Case Study

### Chapter 14: The Case Study

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Case study intro | 36000 | 36200 | Martian invasion scenario |
| Need Perspective — Source Element View | 36200 | 36600 | Figure 14.1 |
| Need Perspective — Requirement Description View | 36600 | 37000 | |
| Need Perspective — Context Definition View | 37000 | 37400 | |
| Need Perspective — Requirement Context View | 37400 | 37900 | Use cases in context |
| Need Perspective — Validation View | 37900 | 38400 | Scenarios validating use cases |
| Need Perspective — Traceability View | 38400 | 38900 | End-to-end traceability shown |
| System Perspective — System Structure View | 39000 | 39500 | Block definition; components |
| System Perspective — Interface Definition View | 39500 | 40000 | Ports, interface blocks, item flows |
| System Perspective — System Configuration View | 40000 | 40500 | Internal block diagrams; assembly |
| System Perspective — System State View | 40500 | 41000 | State machines |
| System Perspective — System Behaviour View | 41000 | 41500 | Activity diagrams |

---

## Part V — Deploying MBSE

### Chapter 15: Benefits of MBSE

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| MBSE benefits narrative | 42100 | 43000 | "Old Lady" analogy; incremental complexity |

### Chapter 17: The Process

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| The ACRE Process — Process Content View | 44500 | 45500 | How to apply ACRE in practice |
| Quick/dirty, semi-formal, formal process variants | 45500 | 46500 | Three levels of rigour |
| Automated compliance | 46500 | 47000 | Process automation via model |

### Chapter 18: The Tool

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Individual Tool vs. Tool Chain | 47500 | 48000 | |
| Tool selection criteria | 48000 | 49000 | Modelling capability, interoperability, process compatibility |
| MonTE process (tool evaluation) | 49000 | 49500 | |

### Chapter 19: Model Structure and Management — HIGHLY RELEVANT

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Model structure approaches | 49984 | 50109 | Life cycle, engineering activity, framework, system hierarchy |
| Model management introduction | 50110 | 50121 | Version management, access, sandboxing, correctness |
| Version management | 50122 | 50179 | Baselining; rolling back; dependency integrity issues |
| Model access control | 50181 | 50220 | Permissions; locking; data integrity |
| Sandboxing | 50228 | 50264 | Branching model; merge back; parallel exploration |
| Correctness through scripting | 50266 | 50400 | Scripts for consistency, correctness, generation |

### Chapter 20: Model Maturity

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Technology/process/individual maturity | 50700 | 51000 | Three dimensions of maturity |
| Technology Readiness Levels for models | 51000 | 51400 | Readiness levels defined |
| Model Maturity assessment | 51400 | 51800 | How to assess and apply |

---

## Part VI — Annex

| Topic | Start | End | Notes |
|-------|-------|-----|-------|
| Appendix A — Ontology and Glossary | 52000 | 54000 | Full glossary of all terms |
| Appendix B — Summary of SysML Notation | 54000 | 58000 | Crib sheets for all nine diagram types |
| Appendix C — Process Model for ISO15288:2015 | 58000 | 60000 | Full model of the standard |
| Appendix D — Competency Framework | 60000 | 61000 | MBSE competency scopes |
| Appendix E — The MBSE Memory Palace | 61000 | 62000 | Mnemonics for concepts |

---

## Quick Reference: Key Passages

| Concept | Line | Exact note |
|---------|------|------------|
| "single source of truth" | 1148 | The model should be considered the single source of truth |
| MBSE definition (book's own) | 1328–1331 | "driven by a model that comprises a coherent and consistent set of views" |
| Benefit: automatic document generation | 1367–1369 | "All system documents may be generated automatically from the model" |
| Benefit: inherent traceability | 1393–1394 | "traceability...is contained within the model" |
| View vs. Diagram distinction | 1878–1891 | View is semantic; Diagram is just visualisation |
| "Consistency is king" | 1453 | Exact phrase |
| Automation benefit in Framework | 2270–2273 | "many of the Process Artefacts may be automatically generated" |
| Block = fundamental structural unit | 9732–9750 | Shows what things exist and their relationships |
| Ports and interfaces | 9853–9875 | Full Port vs. Proxy Port |
| Satisfy relationship definition | 13044–13048 | Design element satisfies requirement |
| Verify relationship definition | 13057–13064 | Test case verifies requirement |
| Derive relationship definition | 13029–13033 | Derived requirement tracing |
| Need in Context = Use Case | 24758–24765 | Core ACRE principle |
| Need Description is not the Need | 24472–24483 | Critical distinction |
| Source Element must be configurable | 24505–24507 | Identifiable by version and location |
| Rules should be automated | 25131–25132 | "In order to maximise benefits...automated rather than manually applied" |
| Flat list principle for needs | 25735–25851 | Never group by type; always group by Context |
| Correctness through scripting | 50266–50282 | Three ways scripts manage model integrity |
