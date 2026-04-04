# NASA Systems Engineering Handbook — Curated Extract
## Relevant to: Blueprint Framework (System Models as Source of Truth, AI-Generated Code)

Source: NASA SP-2016-6105 Rev2 — NASA Systems Engineering Handbook
Full file: `/resources/nasa_se_handbook.txt` (25,273 lines)

---

## 1. Systems Engineering Fundamentals (Lines 1017–1199)

### Definition of Systems Engineering

> "At NASA, 'systems engineering' is defined as a methodical, multi-disciplinary approach for the design, realization, technical management, operations, and retirement of a system. A 'system' is the combination of elements that function together to produce the capability required to meet a need."

The value of a system is primarily created by the **relationships among its parts** — how they are interconnected — not by the parts themselves.

### The 17 SE Processes (SE Engine)

NASA organises all SE work into three groups of common technical processes:

**System Design Processes** (what to build):
1. Stakeholder Expectations Definition
2. Technical Requirements Definition
3. Logical Decomposition
4. Design Solution Definition

**Product Realization Processes** (building it):
5. Product Implementation
6. Product Integration
7. Product Verification
8. Product Validation
9. Product Transition

**Technical Management Processes** (crosscutting control):
10. Technical Planning
11. Requirements Management
12. Interface Management
13. Technical Risk Management
14. Configuration Management
15. Technical Data Management
16. Technical Assessment
17. Decision Analysis

The SE engine is applied **recursively** (repeating the same process at successive lower levels of the system hierarchy) and **iteratively** (reapplying a process to correct discovered discrepancies).

> "Iterative: the application of a process to the same product or set of products to correct a discovered discrepancy."
> "Recursive: adding value to the system by the repeated application of processes to design next lower layer system products or to realize next upper layer end products within the system structure."

---

## 2. Requirements Engineering (Lines 6054–6787)

### 2.1 What Requirements Are

The Technical Requirements Definition Process "transforms the stakeholder expectations into a definition of the problem and then into a complete set of validated technical requirements expressed as 'shall' statements that can be used for defining a design solution."

Requirements types:
- **Functional requirements**: define *what* functions need to be performed
- **Performance requirements**: define *how well* those functions must be performed
- **Interface requirements**: product-to-product interaction requirements
- **Crosscutting requirements**: environmental, safety, human factors, and "-ilities"

### 2.2 Requirement Writing Rules (Appendix C, Lines 18194–18410)

**Terms of art:**
- `shall` = a binding requirement
- `will` = a fact or declaration of purpose
- `should` = a goal (not mandatory)

**Active voice form:**
"The [product] shall [verb] [description]."
Example: "The software shall acquire data from the..."

**Requirements must be:**
- **Clear and unambiguous** — free of terms like "as appropriate", "etc.", "and/or", "but not limited to"
- **Concise** — one thought, one subject, one predicate per statement
- **Complete** — with tolerances for all quantitative values
- **Implementation-free** — state WHAT is needed, not HOW to provide it
- **Consistent** — no self-contradiction
- **Correct** — technically accurate and feasible
- **Verifiable** — measurable and testable
- **Traceable** — bidirectional link to parent and children

**Unverifiable terms to avoid:** flexible, easy, sufficient, safe, adequate, user-friendly, fast, maximize, minimize, robust, quickly, easily, clearly, other "-ly" and "-ize" words.

**TBD vs TBR:**
Avoid TBD (To Be Determined). Prefer TBR (To Be Resolved) with rationale, owner, and resolution date.

### 2.3 Requirements Metadata (Lines 6553–6591)

Each requirement should carry structured metadata:

| Field | Purpose |
|---|---|
| Requirement ID | Unique identifier for sorting and tracking |
| Rationale | Clarifies intent; records reason, assumptions, constraints |
| Traced from | Bidirectional links to parent and child requirements |
| Owner | Person/group responsible for managing the requirement |
| Verification method | Test / inspection / analysis / demonstration |
| Verification lead | Person/group assigned to verify |
| Verification level | System / subsystem / element level |

### 2.4 Rationale Content (Lines 6598–6611)

Rationale must capture:
- **Reason for the requirement** — constraint, ConOps reference, or parent requirement
- **Assumptions** — technology dependencies, programme conditions
- **Relationships** — links to ConOps and operational scenarios
- **Design constraints** — why a particular implementation method was chosen

### 2.5 Requirement Validation (Lines 6616–6686)

Six validation questions for each requirement:
1. Is it written correctly (format/editorial)?
2. Does it satisfy stakeholders?
3. Is it technically correct (bidirectional traceability, valid assumptions, consistent with success criteria)?
4. Have all relevant stakeholders removed defects?
5. Is it feasible?
6. Is it non-redundant and non-overlapping?

### 2.6 Benefits of Well-Written Requirements (Lines 6491–6551)

| Benefit | Rationale |
|---|---|
| Stakeholder agreement | Complete description of what the product must do |
| Reduced rework | Forces rigorous consideration before design begins |
| Cost/schedule basis | Realistic estimation baseline |
| V&V baseline | Both test plans and acceptance criteria derive from requirements |
| Transferability | Eases handover and reuse |
| Enhancement basis | Requirements serve as a baseline for future changes |

### 2.7 Requirements Flowdown (Lines 6288–6465)

Requirements flow from mission authority down through the system hierarchy:

```
Mission Authority
  └── Mission Requirements
        └── System Functional/Performance Requirements
              └── Subsystem A Requirements (allocated)
              └── Subsystem B Requirements (allocated + derived)
              └── Subsystem N Requirements
```

- **Allocated requirements**: flowed down from a parent document
- **Derived requirements**: arise from design decisions not explicitly in the parent

> "In the next level of decomposition, the baselined derived (and allocated) requirements become the set of high-level requirements for the decomposed elements and the process begins again."

---

## 3. Functional Decomposition (Lines 6789–7062)

### 3.1 Purpose of Logical Decomposition

Logical decomposition "identifies the 'what' that should be achieved by the system at each level." It uses **functional analysis** to:
1. Translate top-level requirements into functions
2. Decompose and allocate functions to lower levels of the PBS
3. Identify and describe functional and subsystem interfaces

### 3.2 System Architecture

The system architecture "is the strategic organisation of the functional elements of the system, laid out to enable the roles, relationships, dependencies, and interfaces between elements to be clearly defined and understood."

It is "strategic in its focus on the overarching structure of the system and how its elements fit together to contribute to the whole, instead of on the particular workings of the elements themselves."

The architecture enables elements to be **developed separately** while ensuring they **work together** to achieve top-level requirements.

### 3.3 Functional Analysis Process (Lines 6962–7021)

Three key steps:
1. Translate top-level requirements into functions
2. Decompose and allocate functions to lower levels of the PBS
3. Identify and describe functional and subsystem interfaces

Each function is described in terms of:
- Inputs
- Outputs
- Failure modes
- Consequence of failure
- Interface requirements

> "Functions are arranged in a logical sequence so that any specified operational usage of the system can be traced in an end-to-end path."

The process is recursive and iterative until all desired levels are "analysed, defined, and baselined." The systems engineer must "keep an open mind and a willingness to go back and change previously established architecture and system requirements."

### 3.4 Outputs of Logical Decomposition (Lines 7029–7053)

- **Logical Decomposition Models**: define relationships of requirements, functions, and behaviours; include system architecture models
- **Derived Technical Requirements**: requirements that arise from architectural choices not explicitly stated in the parent requirements
- **Logical Decomposition Work Products**: other products generated during the process

---

## 4. Interface Definition and Management (Lines 12273–12497)

### 4.1 Why Interface Management Matters

> "The definition, management, and control of interfaces are crucial to successful programmes or projects. Interface management is a process to assist in controlling product development when efforts are divided among parties and/or to define and maintain compliance among the products that should interoperate."

### 4.2 Basic Interface Management Tasks (Lines 12237–12297)

- Define interfaces
- Identify characteristics (physical, electrical, mechanical, human, etc.)
- Ensure interface compatibility at all defined interfaces
- Control all interface processes during design and construction
- Identify assembly documentation and integration instructions

### 4.3 Interface Management Activities (Lines 12383–12443)

**During system design:** Analyse the ConOps to identify both external and internal interfaces. Establish origin, destination, stimuli, and special characteristics.

**During product integration:** Review integration and assembly procedures to ensure interfaces are properly marked and compatible with ICDs. Interface control documentation is an input to both Product Verification and Product Validation.

### 4.4 Interface Control

An **Interface Working Group (IWG)** establishes communication between those responsible for interfacing systems. The IWG ensures planning, scheduling, and execution of all interface activities. It may work independently or as part of a larger change control board.

> "Interface requirements verification is a critical aspect of the overall system verification."

### 4.5 Interface Documentation Types (Lines 12446–12479)

- **Interface Requirements Document (IRD)**: top-level interface requirements
- **Interface Control Document/Drawing (ICD)**: detailed interface definitions and drawings
- **Interface Definition Document (IDD)**: interface parameter definitions
- **Interface Control Plan (ICP)**: process for managing interface control

All interface documentation feeds into the Configuration Management Process.

> "For interfaces that require approval from all sides, unanimous approval is required. Changing interface requirements late in the design or implementation life cycle is more likely to have a significant impact on the cost, schedule, or technical design/operations."

---

## 5. Traceability (Lines 11860–11879, 12014–12053)

### 5.1 Definitions

> "Traceability: A discernible association between two or more logical entities such as requirements, system elements, verifications, or tasks."

> "Bidirectional traceability: The ability to trace any given requirement/expectation to its parent requirement/expectation and to its allocated children requirements/expectations."

### 5.2 Requirements Management Traceability (Lines 12014–12053)

The Requirements Management Process maintains bidirectional traceability between:
- Stakeholder expectations
- Customer requirements
- Technical product requirements
- Product component requirements
- Design documents
- Test plans and procedures

> "Requirements traceability is usually recorded in a requirements matrix or through the use of a requirements modelling application."

Traceability integrity check: "Ensure that all top-level parent document requirements have been allocated to the lower level requirements. If there is no parent for a particular requirement and it is not an acceptable self-derived requirement, it should be assumed either that the traceability process is flawed and should be redone or that the requirement is 'gold plating' and should be eliminated."

### 5.3 Requirements Verification Matrix (Appendix D, Lines 18413–18492)

Every "shall" statement requires a corresponding verification record:

| Field | Content |
|---|---|
| Document number | Source specification |
| Requirement ID | Unique identifier |
| Shall statement | The exact requirement text |
| Verification success criteria | Measurable pass/fail criteria |
| Verification method | Test / Analysis / Inspection / Demonstration |
| Verification level | System / subsystem / element |
| Verification lead | Responsible person/team |
| Results | Pass/fail and date |

---

## 6. Verification and Validation (Lines 1844–1878, 8854–10120)

### 6.1 Fundamental Distinction

> "Verification of a product shows proof of compliance with requirements — that the product can meet each 'shall' statement as proven through performance of a test, analysis, inspection, or demonstration."

> "Validation of a product shows that the product accomplishes the intended purpose in the intended environment — that it meets the expectations of the customer and other stakeholders."

Verification relates back to the **approved requirements set**.
Validation relates back to the **ConOps document**.

> "Are we building the product right? (verification) Are we building the right product? (validation)"

### 6.2 Methods of Verification (Lines 9151–9176)

- **Analysis**: mathematical modelling and analytical techniques; includes modelling and simulation; used when a fabricated product is not available
- **Demonstration**: showing that use of the product achieves the specified requirement; basic confirmation without detailed data gathering; uses mock-ups or simulators
- **Inspection**: visual examination; used for physical design features
- **Test**: use of the end product to obtain detailed data under controlled conditions; most resource-intensive; "Test as you fly, and fly as you test"

### 6.3 Methods of Validation (Lines 9804–9826)

Same four methods as verification but applied against stakeholder expectations (ConOps, MOEs) rather than "shall" statements.

### 6.4 End-to-End Testing (Lines 9296–9438)

> "End-to-end testing verifies that the data flows throughout the multisystem environment are correct, that the system provides the required functionality, and that the outputs at the eventual end points correspond to expected results."

End-to-end tests:
- Execute complete operational scenarios across multiple configuration items
- Focus on **external interfaces** (hardware, software, or human-based)
- Demonstrate interface compatibility and desired total functionality

### 6.5 Verification Closure Requirements (Lines 9480–9499)

> "Verification results should be recorded in a requirements compliance or verification matrix or other method developed during the Technical Requirements Definition Process to trace compliance for each product requirement."

Discrepancies trigger a discrepancy report. Nonconformances may require re-engineering and re-verification.

### 6.6 Validation Deficiencies (Lines 9964–10027)

> "Validation should be performed as early and as iteratively as possible in the SE process since the earlier re-engineering needs are discovered, the less expensive they are to resolve."

A system can pass verification (meets its stated requirements) but fail validation (does not satisfy stakeholder expectations) when:
- The ConOps was poorly defined in early phases
- User community was not adequately involved
- Requirements did not adequately capture operational needs

---

## 7. Requirements Management (Lines 11860–12263)

### 7.1 Process Purpose

Requirements Management is used to:
- Identify, control, decompose, and allocate requirements across all levels of the WBS
- Provide bidirectional traceability
- Maintain consistency between requirements, ConOps, and architecture/design
- Evaluate all change requests over the life of the project

### 7.2 Change Control (Lines 12063–12155)

Once requirements are validated at the System Requirements Review (SRR), they go under **formal configuration control**. All subsequent changes require Configuration Control Board (CCB) approval.

Change impact assessment must evaluate:
- Cost and schedule impact
- Performance margins
- Interfaces
- ConOps consistency
- Higher and lower level requirements

### 7.3 Requirements Creep (Lines 12159–12197)

> "Requirements creep is the term used to describe the subtle way that requirements grow imperceptibly during the course of a project."

Prevention techniques:
- Develop a thorough ConOps agreed to by all stakeholders
- Flush out conscious, unconscious, and undreamed-of requirements early
- Establish strict CCB-based change request channels
- Measure the functionality of each change request against system impact
- Determine if the change can be accommodated within resource margins

---

## 8. Configuration Management (Lines 12912–13551)

### 8.1 Five Elements of CM (Lines 12912 area)

NASA CM comprises five elements:
1. **Configuration identification** — selecting and documenting CI attributes with unique identifiers
2. **Configuration change management** — systematic proposal, justification, evaluation, and incorporation of changes
3. **Configuration status accounting (CSA)** — recording and reporting of configuration data
4. **Configuration verification and audits** — inspecting documents, products, and records to confirm attributes
5. **CM planning** — strategy document governing all CM activities

### 8.2 The Four Technical Baselines (Lines 13135–13300)

| Baseline | Established At | Contains |
|---|---|---|
| **Functional Baseline** | System Definition Review (SDR) | System-level performance, interface, and verification requirements |
| **Allocated Baseline** | Preliminary Design Review (PDR) | Functional and performance requirements allocated to CIs |
| **Product Baseline** | Critical Design Review (CDR) | Detailed form/fit/function characteristics; production acceptance test requirements |
| **As-Deployed Baseline** | Operational Readiness Review (ORR) | Final deployed configuration; all changes incorporated |

### 8.3 Change Types (Lines 13421–13431)

- **Engineering Change (Major)**: significant impact — requires retrofit of delivered products, affects cost/safety/compatibility
- **Engineering Change (Minor)**: modifies documentation without interchangeability impact
- **Waiver**: documented relief from a requirement; does not constitute a change to the baseline

### 8.4 CM Key Considerations for Software Tools (Lines 13457–13492)

Critical attributes for CM tooling:
- Version control and comparison (track history of an object or product)
- Secure check-out/check-in
- Real-time data sharing with internal and external stakeholders
- Tracking capabilities (time, date, who, time in phases)
- Integration with drafting and modelling programs
- Workflow and lifecycle management
- Capable of acting as the one and only source for released information

### 8.5 CM Outputs (Lines 13509–13531)

- **Configuration Status Accounting (CSA) reports**: complete list of items under control, updated throughout lifecycle
- **Current baselines**: available to all technical teams and stakeholders
- **CM reports**: periodic status reports at agreed frequency and key reviews
- **Other work products**: strategy, procedures, descriptions/drawings/models, change requests and dispositions, audit results

---

## 9. Technical Reviews (Lines 13907–14979)

### 9.1 Purpose

Technical assessment monitors technical progress through:
- Periodic technical reviews
- Technical performance indicators (MOEs, MOPs, KPPs, TPMs)

> "Technical assessment is focused on providing a periodic assessment of the programme/project's technical and programmatic status and health at key points in the life cycle."

### 9.2 Life-Cycle Reviews for Spaceflight Projects (Lines 14395–14978)

| Review | Phase | Key Question |
|---|---|---|
| Mission Concept Review (MCR) | Pre-Phase A | Does the mission concept meet the need? |
| System Requirements Review (SRR) | Phase A | Are requirements complete and consistent with the mission? |
| System Definition Review (SDR) | Phase A | Is the architecture responsive to requirements? |
| Preliminary Design Review (PDR) | Phase B | Does the preliminary design meet all requirements with acceptable risk? |
| Critical Design Review (CDR) | Phase C | Is the design mature enough to support fabrication? |
| Production Readiness Review (PRR) | Phase C | Are production plans and facilities ready? |
| System Integration Review (SIR) | Phase C/D boundary | Are all components ready for integration? |
| System Acceptance Review (SAR) | Phase D | Does the system meet maturity and compliance criteria? |
| Operational Readiness Review (ORR) | Phase D/E boundary | Is the system ready for operations? |
| Flight Readiness Review (FRR) | Phase D | Is everything ready for launch? |

### 9.3 Peer Reviews (Lines 9139–9141)

> "Peer reviews are additional reviews that may be conducted formally or informally to ensure readiness for verification (as well as the results of the verification process)."

Test Readiness Reviews (TRRs) assess readiness of test ranges, facilities, instrumentation, integration labs, trained testers, and support equipment before each major test.

### 9.4 Status Reporting Principles (Lines 14250–14267)

- Use an agreed-upon set of well-defined technical measures
- Report measures consistently at all project levels
- Maintain historical data for trend identification
- Use colour-coded (red/yellow/green) alert zones for all technical measures
- Support assessments with quantitative risk measures

---

## 10. Key Design Principles Directly Applicable to Blueprint

The following quotes and principles from the handbook are directly relevant to the Blueprint framework design.

**Requirements as the single source of truth:**
> "Complete and thorough requirements traceability is a critical factor in successful validation of requirements."

> "Document all decisions made during the development of the original design concept in the technical data package. This will make the original design philosophy and negotiation results available to assess future proposed changes."

**Models over documents:**
> "References to 'documents' are intended to include not only paper or digital files but also models, graphics, drawings, or other appropriate forms that capture the intended information."

**Early validation saves cost:**
> "The life cycle costs of a programme or project tend to get 'locked in' early in design and development. The cost curves clearly show that late identification of and fixes to problems cost considerably more later in the life cycle. During design, only about 15% of the costs might be expended, but the design itself will commit about 75% of the life cycle costs."

**Iterative decomposition:**
> "The technical processes are applied recursively and iteratively to break down the initialising concepts of the system to a level of detail concrete enough that the technical team can implement a product from the information."

**Verification method defined at requirements time:**
> "Verification method: Captures the method of verification (test, inspection, analysis, demonstration) and should be determined as the requirements are developed." (Requirements Metadata table)

**Traceability is bidirectional:**
> "Bidirectional traceability: The ability to trace any given requirement/expectation to its parent requirement/expectation and to its allocated children requirements/expectations."

---

## 11. Appendix Reference: Requirements Verification Matrix Structure

The NASA Requirements Verification Matrix (Appendix D) maps each requirement to its verification:

```
Document | Req ID | Shall Statement | Success Criteria | Method | Level | Lead | Result
```

This directly maps to Blueprint's need for a traceable requirements-to-tests matrix, where AI agents can:
1. Receive a "shall" statement as input
2. Identify the verification method (test/analysis/inspection/demonstration)
3. Generate the verification artefact (test code, analysis script, etc.)
4. Record pass/fail evidence against the requirement ID

---

*End of curated extract. Total source lines: 25,273. Key line ranges indexed in `nasa_se_handbook_index.md`.*
