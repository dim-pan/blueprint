# SysML for Systems Engineering — Curated Extract for Blueprint

Source: *SysML for Systems Engineering: A Model-Based Approach*, 3rd Edition
Authors: Jon Holt and Simon Perry (IET, 2019)

This extract is curated for the Blueprint framework project. It captures the **concepts and principles** relevant to a system model as source of truth, with AI agents generating code as artefacts. Diagram notation details are omitted; the focus is on modelling philosophy, structural concepts, requirements engineering, traceability, and the relationship between model and implementation.

---

## 1. The Core Proposition: Model as Single Source of Truth

When MBSE is applied properly, the model becomes the collected knowledge associated with the project or system and, ideally, should be considered the **single source of truth**.

Benefits of a coherent model-based approach include:

- **Automatic generation and maintenance of system documents.** All system documents may be generated automatically from the model, resulting in simpler document maintenance, more consistent content, and drastically reduced documentation effort.
- **Complexity control and management.** Models may be measured, and hence controlled. This measurement may be automated and used to control the complexity of the model.
- **Consistency of information.** A true model produces consistent and coherent views across the whole system architecture.
- **Inherent traceability.** When the model is correct, traceability between all system artefacts, across all life cycle stages, is contained within the model.
- **Simpler access to information.** Without a coherent model, knowledge of the system is potentially spread across multiple sources — heterogeneous models, spreadsheets, and documents.
- **Improved communication.** When a model is in place and defined using an established notation, that notation becomes a common language.
- **Automation.** The MBSE Framework provides the basis for automating the approach. One of the main benefits is that it saves a lot of time and effort, as many of the process artefacts may be automatically generated.

### The MBSE Definition Used in This Book

> Model-based Systems Engineering is an approach to realising successful systems that is driven by a **model that comprises a coherent and consistent set of views that reflect multiple viewpoints of the system**.

The key shift from traditional systems engineering: modelling should **drive** engineering activities, not merely support them. The model is model-centric, not document-centric.

---

## 2. Fundamental MBSE Concepts: Ontology, Framework, Views

The core structural approach underpinning everything is: **Ontology → Framework → Views**.

### 2.1 The Goal

One or more **Model** abstracts a **System**. The Model is made up of a set of one or more **Views**. These Views all interact with one another to make up the Model.

### 2.2 Ontology

The **Ontology** is the domain-specific language. It identifies and defines the concepts, terminology, and inter-relationships that form the foundation for all modelling activities. The Ontology:

- Defines concepts and terms with their interrelationships.
- Provides the basis for all Framework definitions.
- Enforces consistency and rigour across Viewpoints.
- Can serve as the basis for process automation.

The Ontology is composed of **Ontology Elements** of two kinds:
- **Core Elements** — specific MBSE concepts with explicit relationships.
- **Cross-cutting Elements** — concepts that apply across the entire Ontology (e.g., Context, Traceable Element, Testable Element, Interface, Process).

Key Cross-cutting Elements relevant to Blueprint:
- **Traceable Element** — any Ontology Element may be traced to another. Traceability is inherent.
- **Testable Element** — any Ontology Element may be tested.
- **Interface Element** — many elements have interfaces (Process, Stage, View, etc.).
- **Context** — many elements have their own Context (Need, Project, Stakeholder, Life Cycle, etc.).

### 2.3 Framework

The **MBSE Framework** comprises the Ontology plus a set of **Viewpoints**. Each Viewpoint focuses on a subset of the Ontology and has a specific purpose. Viewpoints serve as templates for Views.

Properties of a good Framework:
- **Coverage** — the totality of the Viewpoints covers the whole Ontology.
- **Rigour** — generating all Views based on defined Viewpoints produces a true model.
- **Flexibility of scale** — not all Viewpoints need be realised for every project.
- **Flexibility of realisation** — Views may be visualised in any suitable notation.
- **Automation** — the Framework provides the basis for automating process artefact generation.

A **Pattern** is a special type of Framework: a defined set of Viewpoints and associated Views that may be re-used for multiple applications. A Framework is bespoke; a Pattern is reusable.

### 2.4 Views and Viewpoints

Critical distinction:
- A **Viewpoint** is a definition — a template for one or more Views, defined as part of the MBSE Framework.
- A **View** is a realisation — an actual artefact (usually a diagram) produced as part of a project. Each View must conform to its defining Viewpoint.

A View forms part of the Model and has a purpose, a set of interested Stakeholder Roles, and yields some benefit. A single View may be visualised using multiple notations (graphical, textual, tabular). The View is the same; the visualisation may change.

Creating Views without ensuring they are consistent does not produce a Model. **Consistency is king in modelling.**

### 2.5 Rules

Rules constrain the Framework, enforcing correctness and consistency at the Viewpoint or Viewpoint Element level. Rules are the mechanism by which model consistency may be automated — a critical property for agent-driven code generation.

---

## 3. System Decomposition Concepts (Block Definition Diagrams)

The block definition diagram is the primary structural modelling tool. It shows **what conceptual things exist in a System** and what relationships exist between them.

### 3.1 Block

A **Block** is the fundamental unit of structure. It describes a type of thing that exists in the system. A Block is made up of:

- **Part Properties** — things owned by the Block (intrinsic, with their own identity). Shown as composition or aggregation.
- **Reference Properties** — things referenced but not owned by the Block. Shown as associations.
- **Value Properties** — quantities without individual identity (numbers, colours, etc.).
- **Flow Properties** — elements that can flow to or from a block (used with ports and item flows).
- **Operations** — behaviours the Block supports.
- **Constraints** — rules that constrain properties.

A Block may also have **Features** — properties or operations it provides to others (Provided Feature), requires from others (Required Feature), or both.

### 3.2 Relationships Between Blocks

Three main types:
- **Composition** — the owning block wholly contains the part (the part's identity is subordinate to the whole).
- **Aggregation** — the part may be shared between multiple owning blocks.
- **Association** — a general relationship; the related block remains independent.
- **Generalisation** — one block specialises another (inheritance of type).
- **Dependency** — the weakest relationship; shows a general, often unspecified dependency.

### 3.3 Ports and Interfaces

**Ports** are interaction points defined on a block — the mechanism by which blocks interact with the external world.

- **Full Port** — a separate element of the model with its own internal parts and behaviour. Represents a distinct interaction point.
- **Proxy Port** — makes features of its owning block available to external blocks. It has no separate internal parts; it exposes the owner's features. A Proxy Port may only be typed by an Interface Block.
- **Nested Ports** — ports may contain other ports, supporting hierarchical interface decomposition.

**Interface Block** — a special type of block used specifically to define interfaces. It defines the contract between interacting blocks.

**Item Flows** — flows between two ports. An item flow conveys a Flow Property (the type of thing flowing). Item flows make explicit what passes across an interface.

Practical guidance: if uncertain whether a port should be full or proxy, leave it as a plain port and decide later as the model evolves.

### 3.4 Levels of Abstraction

Any System may be considered at many levels of abstraction — from the overall system as a single entity down to individual components. Effective modelling requires views at multiple levels: high (single entity), intermediate (subsystem), and low (component/interface detail).

Complexity of a system manifests through **relationships between elements**, not just from the elements themselves. The whole is more complex than the sum of its parts.

---

## 4. Internal Structure (Internal Block Diagrams)

Where block definition diagrams show types and type-level relationships, **internal block diagrams** show the actual internal assembly of a specific block — how its parts are wired together.

Key elements:
- **Parts** — instances of blocks that make up the containing block.
- **Binding Connectors** — connections between ports of parts, showing the actual wiring of the internal assembly.
- **Shared Parts** — parts that are shared with external context (shown as participants).

The internal block diagram makes explicit the logical relationships between elements within a block, not just that they are elements of the block. This is where interface connections, ports, and item flows become concrete and visible.

---

## 5. Requirements Modelling (the ACRE Approach)

The ACRE approach (Approach to Context-based Requirements Engineering) provides a structured, model-based framework for capturing, organising, and tracing requirements.

### 5.1 The Need Taxonomy

The fundamental concept is the **Need** — an abstract concept that, when put into a Context, represents something necessary or desirable. Four specialisations:

- **Requirement** — a property of a System needed or wanted by a Stakeholder Role. One or more Requirements are needed to deliver each Capability.
  - *Functional Requirement* — yields an observable result to something using the System.
  - *Non-functional Requirement* — constrains how a Functional Requirement may be realised.
  - *Business Requirement* — states the needs of the business; drives all projects.
- **Capability** — the ability of an Organisation or Organisational Unit; meets one or more Goal.
- **Goal** — a high-level organisational Need; met by one or more Capability.
- **Concern** — a Need whose Context is an Architecture, Architectural Framework, or Viewpoint.

A **Need Description** is a tangible description of an abstract Need, with defined attributes. It is not the Need itself but a structured representation. Key attributes: id, name, description text, origin, priority, ownership, verification criteria, validation criteria.

**Critical principle:** Need Descriptions should be kept as a flat list, grouped only by Context (Use Case) — not by type or functionality. Grouping by function before contextualisation destroys meaning.

### 5.2 Source Elements

A **Source Element** is the ultimate origin of a Need. It can be almost anything: requirements documents, emails, conversations, standards, existing systems, workshops, specifications. The key requirement: every Source Element must be a configurable item (identifiable by version and location).

Every Need Description must be traceable to at least one Source Element, and every Source Element must trace to at least one Need Description. This bidirectional traceability is mandatory.

### 5.3 Rules

**Rules** constrain the attributes of Need Descriptions. Forms include:
- Vocabulary restrictions (forbidden words: "quick", "reasonable", "minimum", "maximum").
- Complexity measures (e.g., Flesch-Kincaid readability score must fall within a defined range).
- Attribute value constraints (e.g., Priority must be one of: Essential, Desirable, Bells-and-whistles, Unknown).
- Structural constraints on the set of Need Descriptions.

Rules are best realised as automated checks rather than manual review.

### 5.4 Context and Use Cases

A **Context** is a specific point of view. Types include: System Context (based on hierarchy level), Stakeholder Context (based on Stakeholder Role), Project Context, Process Context.

A **Use Case** is a Need that has been given meaning by putting it into Context. A single Need Description may give rise to multiple Use Cases — one per Context. Different Contexts may conflict (e.g., the passenger wanting a cheap fare vs. the airline owner wanting to cut costs).

This is the heart of the ACRE approach: **Needs are given meaning only in Context.** A Need Description without Context is incomplete.

### 5.5 Scenarios (Validation)

A **Scenario** is an exploration of a "what if" for a Use Case — an ordered set of interactions with a specific outcome that validates a Use Case. Two types:

- **Semi-formal Scenario** — demonstrable through visual interaction diagrams (e.g., sequence diagrams). Used at two levels: Stakeholder-level (system as black box) and System-level (internal element interactions).
- **Formal Scenario** — mathematically provable, using constraint networks (parametric models). Particularly powerful for trade-off analysis and safety-critical systems.

Scenarios serve multiple purposes:
- Understanding — analysing each Use Case.
- Calibrating the right level of abstraction.
- Validation — demonstrating Use Cases can be satisfied.
- Providing a tangible link between the Needs Model and the rest of the System model.

### 5.6 The ACRE Framework Viewpoints

The ACRE Framework defines six core Viewpoints:

1. **Source Element Viewpoint** — all relevant source information. Primary purpose: establishing traceability. Every Need Description must trace back to at least one Source Element here.

2. **Requirement Description Viewpoint** — structured Need Descriptions, each with full attributes (id, name, text, origin, priority, ownership, verification/validation criteria). Visualised as a flat list.

3. **Definition Rule Set Viewpoint** — Rules constraining the attributes of Need Descriptions. May be simple text rules or formal mathematical constraint blocks.

4. **Requirement Context Viewpoint** — Needs put into Context, expressed as Use Cases. This is where Needs acquire meaning. Multiple diagrams, one per Context.

5. **Context Definition Viewpoint** — defines the Contexts (points of view) explored in the Requirement Context Viewpoint. Identifies Stakeholder Roles and System hierarchy levels.

6. **Validation Viewpoint** — Scenarios that validate the Use Cases. Semi-formal (sequence diagrams) or Formal (parametric constraints). Every Use Case must have at least one validating Scenario.

Plus a cross-cutting **Traceability Viewpoint** — explicitly links between Views and between View Elements across the whole model.

---

## 6. Traceability Relationships

The SysML requirement diagram supports six types of relationship. These are the formal vocabulary of traceability:

### 6.1 Satisfy

A model element (block, component, subsystem) **satisfies** a Requirement. Used to show that a design or implementation element meets the intent of a requirement. This is the primary link between requirements and design/implementation.

Direction: `DesignElement --«satisfy»--> Requirement`

### 6.2 Verify

A test case **verifies** a Requirement. The test case demonstrates that the requirement has been met. Any SysML element (sequence diagram, parametric diagram) may be stereotyped as `«testCase»` and connected to a requirement via this relationship.

Direction: `«testCase» TestScenario --«verify»--> Requirement`

### 6.3 Derive

A Requirement is **derived from** another Requirement. Used when systems engineers derive implicit requirements from explicit ones during analysis. Derived requirements are not directly stated by stakeholders but follow necessarily from stated requirements.

Direction: `DerivedReq --«deriveReqt»--> SourceReq`

### 6.4 Refine

A model element or requirement **refines** another. Used to show how a use case, design element, or more detailed requirement adds precision or detail to a higher-level element.

Direction: `DetailElement --«refine»--> AbstractElement`

### 6.5 Trace

A general-purpose relationship showing that one model element can be traced to another. Weaker semantics than the above — says nothing about the nature of the relationship. Should be used only when a more specific type is not applicable.

Direction: bidirectional, typically `ModelElement --«trace»--> SourceElement`

### 6.6 Nesting (Decomposition)

A Requirement may be decomposed into sub-Requirements via nesting. Used when a requirement is not atomic and must be broken down into related atomic sub-requirements.

### 6.7 Principle

Where possible, use the most specific relationship type (satisfy, verify, derive, refine) rather than the generic trace. The generic trace has weakly defined semantics and makes automated reasoning harder.

---

## 7. Use Cases and Test Cases

### 7.1 Use Cases

A Use Case is a Need in Context — the central mechanism for giving requirements meaning. Key properties:
- Belongs to a specific Context (Stakeholder, System hierarchy level, etc.).
- Is validated by one or more Scenarios.
- May include or extend other Use Cases.
- Stakeholder Roles interact with Use Cases (they have an interest in them).
- A Need Description must relate to at least one Use Case.
- A Use Case must have at least one validating Scenario.

Use cases that appear related may conflict when the contexts differ. Identifying these conflicts is one of the key analytical benefits of the use case model.

### 7.2 Test Cases

A **test case** in SysML is the `«testCase»` stereotype applied to a SysML element (typically a sequence diagram, activity diagram, or parametric diagram). It represents a specific execution path that verifies a requirement.

The `«verify»` relationship connects the test case to the requirement it verifies. This creates an explicit, machine-traversable link from implementation validation back to the requirement.

The ACRE approach extends this: Scenarios (both semi-formal and formal) serve as the primary validation mechanism. A formal Scenario (parametric constraint network) is particularly powerful for mathematically provable verification.

---

## 8. The Relationship Between Model and Implementation

### 8.1 Model as Source of Truth

The model is not a description of the implementation — it is the authoritative source from which implementations are generated. When the model is correct, traceability between all artefacts exists within the model itself.

Key corollary: a diagram is not a model. Creating diagrams even with standard notation does not necessarily produce a model. A Model requires all Views to be consistent, and all Views to conform to their Viewpoints.

### 8.2 Code as a Model Artefact

The MBSE perspective treats generated documents, specifications, and — by extension — generated code as artefacts produced from model elements. The `«satisfy»` relationship is the formal bridge: a design or implementation element satisfies a requirement. An AI agent generating code is performing the same function as a `«satisfy»` resolution — it is traversing the model graph from requirements to implementation.

### 8.3 Automation

The Framework provides the basis for automating MBSE:
- Documents may be generated automatically from the model.
- Rules may be enforced automatically (consistency checking, completeness checking).
- Model correctness can be verified through scripting against the model's meta-model.
- Artefacts (including code) can be generated from model elements.

The Ontology, Framework, and Viewpoints define what should be generated and the constraints it must satisfy. The generation process is guided by the traceability graph.

### 8.4 Correctness Through Scripting

Industry-grade modelling tools support scripting to enhance model management in three ways:
- **Consistency checking** — verify that model elements conform to rules (e.g., every Need Description traces to a Source Element).
- **Correctness checking** — verify structural and semantic correctness of the model.
- **Generation** — automatically generate documents, reports, or other artefacts from the model.

These scripts operate on the model's internal representation, traversing the graph of model elements and their relationships. This is directly analogous to how Blueprint agents would traverse the system model to generate code.

---

## 9. Model Structure

### 9.1 Structuring Principles

A model must be well structured to support navigability and usability. Common structuring approaches:
- By Life Cycle Stage
- By engineering process or activity
- By System and sub-system (structural vs. behavioural split)
- By Framework Perspective and Viewpoint

The structure should be defined to meet the needs of the project and the users of the model. The structure may change during the model's lifetime as understanding grows.

### 9.2 Packages

Packages are the primary organisational unit. Packages organise model elements similarly to how directories organise files. Packages may contain other packages, blocks, diagrams, and any other model elements. Packages may be imported (public) or accessed (private) by other packages, with visibility rules analogous to module visibility in code.

### 9.3 Version Management

Version management for models mirrors version management for code:
- The model should be stored in a version control system.
- Baselining allows snapshots at specific points in time.
- Rolling back one branch may break relationships in another branch — model integrity requires coordination.
- Process discipline around access, locking, and change management is essential.

### 9.4 Model Access and Sandboxing

- **Access control** — different users have different permissions (editor, reviewer, administrator). Limits accidental corruption.
- **Sandboxing** — multiple copies of a model at a point in time, each evolved independently, then merged. Enables parallel exploration of design alternatives while preserving the original model's integrity. Directly analogous to git branching.

---

## 10. Key Principles for Blueprint

Synthesised from the above, these are the principles most relevant to Blueprint's design:

1. **The model is the source of truth.** Code and documents are artefacts generated from the model, not the model itself.

2. **Needs must be given Context to have meaning.** A requirement without a Use Case (Context) is incomplete. The Use Case is what the agent must satisfy.

3. **Every Need Description must trace to its source.** Source Elements are first-class model objects. Traceability is bidirectional and mandatory.

4. **Satisfy is the key relationship.** The formal statement that a design or implementation element meets a requirement. Every generated artefact should be linked to the requirements it satisfies.

5. **Verify closes the loop.** Test cases (Scenarios) verify requirements. The `«verify»` link connects validation back to the requirement, completing the traceability chain: Source → Requirement → Use Case → Implementation → Test.

6. **Consistency is enforced by the Ontology.** The Ontology defines the domain language. All Views must use elements from the Ontology. Cross-view consistency is the hallmark of a true model.

7. **Rules should be automated.** Consistency checks, completeness checks, and validation rules should be executable, not just documented. This is the direct entry point for agents.

8. **Decomposition is hierarchical and typed.** Components are composed via typed relationships (composition, aggregation, association). Interfaces are explicit (ports, interface blocks, item flows). Type information enables automated checking and code skeleton generation.

9. **The model structure mirrors the engineering structure.** Packages correspond to architectural boundaries. This alignment enables agents to navigate model structure directly.

10. **Views are not Diagrams.** A View is a semantic artefact with purpose and stakeholder relevance. A Diagram is merely its visualisation. Agents reason about Views, not Diagrams.
