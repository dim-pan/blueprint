# NASA SE Handbook — Topic Index
## Source file: `/resources/nasa_se_handbook.txt` (25,273 lines)

This index maps topics to approximate line ranges in the original text file.
Use these ranges with `Read` (offset + limit) for targeted extraction.

---

## Front Matter

| Topic | Lines | Notes |
|---|---|---|
| Title page, cover photos, contact info | 1–34 | Skip for content |
| Table of Contents | 35–363 | Chapter/section overview |
| Table of Figures | 366–531 | Figure list |
| Table of Tables | 534–628 | Table list |
| Table of Boxes | 632–695 | Sidebar/box list |
| Preface | 698–750 | Mentions Model-Based SE |
| Acknowledgments | 755–948 | Skip for content |

---

## Chapter 1: Introduction

| Topic | Lines | Notes |
|---|---|---|
| Purpose and scope | 950–1013 | Includes "documents = models" statement |

---

## Chapter 2: Fundamentals of Systems Engineering

| Topic | Lines | Notes |
|---|---|---|
| SE definition | 1017–1068 | Core definition of SE at NASA |
| Role of the systems engineer | 1069–1133 | Responsibilities overview |
| SE Engine diagram and 17 processes | 1210–1365 | Key: iterative/recursive definitions |
| Project life cycle phase overview | 1452–1764 | SE engine applied per phase |
| Using the SE Engine (Pre-Phase A through D) | 1770–1836 | Narrative walkthrough |
| Verification vs Validation distinction | 1843–1878 | CRITICAL: core V&V definitions |
| Cost effectiveness / life cycle cost curves | 1880–1951 | 75% of LCC locked in at design |
| Human Systems Integration in SE | 1953–2100 | HSI; less relevant to Blueprint |

---

## Chapter 3: NASA Program/Project Life Cycle

| Topic | Lines | Notes |
|---|---|---|
| SE Product Maturity table (by phase) | 2730–2999 | Requirements/interfaces/designs matured per KDP |
| Phase A: Concept & Technology Development | ~3100–3300 | Requirements definition begins |
| Phase B: Preliminary Design | ~3300–3450 | PDR; requirements baselined |
| Phase C: Final Design & Fabrication | ~3450–3550 | CDR; build-to baseline |
| Phase D: System Assembly, Integration & Test | 3552–3629 | V&V activities |
| Phase E: Operations | 3632–3668 | Sustaining engineering |
| Phase F: Closeout | 3670–3743 | Archiving |
| Tailoring and Customisation of NPR 7123.1 | 3783–5050 | Compliance matrix examples |

---

## Chapter 4: System Design Processes

| Topic | Lines | Notes |
|---|---|---|
| System Design overview / interrelationship diagram | 5195–5343 | FIGURE 4.0-1; key flowdown diagram |
| System Design Keys (sidebar) | 5321–5334 | CRITICAL: traceability, clarity, documentation |
| **4.1 Stakeholder Expectations Definition** | 5368–6050 | ConOps, NGOs, MOEs |
| Stakeholder identification | 5400–5540 | Stakeholder table by lifecycle phase |
| ConOps vs Operations Concept (sidebar) | 5871–5884 | Important distinction |
| **4.2 Technical Requirements Definition** | 6054–6788 | CORE SECTION |
| Requirements process description | 6076–6121 | Inputs, activities, outputs |
| Functional vs performance requirements | 6207–6270 | With TVC example |
| Requirements flowdown figure (FIGURE 4.2-3) | 6400–6467 | Allocated vs derived |
| Benefits of well-written requirements (TABLE 4.2-1) | 6491–6551 | Six key benefits |
| Requirements metadata (TABLE 4.2-2) | 6553–6591 | ID, rationale, traceability, owner, V method |
| Rationale box | 6598–6612 | What to capture in rationale |
| Requirements validation (six steps) | 6616–6686 | Completeness check |
| MOPs and TPMs | 6705–6764 | Measures of Performance / Technical Performance |
| **4.3 Logical Decomposition** | 6789–7062 | Functional decomposition |
| Logical decomposition process description | 6773–6999 | Three steps of functional analysis |
| System architecture development | 6824–6913 | Architecture = strategic structure |
| Functional analysis: three key steps | 6972–6999 | Translate → Decompose → Interface |
| Outputs: derived requirements & models | 7029–7053 | Key artefacts |
| **4.4 Design Solution Definition** | 7063–7500 | Design alternatives and trade studies |
| Doctrine of successive refinement | 7203–7347 | FIGURE 4.4-2 |

---

## Chapter 5: Product Realization

| Topic | Lines | Notes |
|---|---|---|
| Product Realization Keys (sidebar) | 8650 area | Brief summary |
| **5.1 Product Implementation** | ~8300–8700 | Build/code/buy/reuse |
| **5.2 Product Integration** | ~8700–8900 | Integration activities |
| **5.3 Product Verification** | 8854–9499 | CORE SECTION |
| Verification process description | 8850–9120 | Inputs, activities, outputs |
| Differences: Verification vs Testing (sidebar) | ~8880–8900 | Key distinctions |
| Methods of Verification (sidebar) | 9151–9176 | Analysis / Demo / Inspection / Test |
| Verification procedure content (TABLE 5.3-1) | 9195–9285 | What a procedure and report must contain |
| End-to-end testing | 9296–9438 | Interface compatibility; operational scenarios |
| Verification result recording | 9480–9499 | Requirements compliance matrix |
| **5.4 Product Validation** | 9540–10128 | CORE SECTION |
| Methods of Validation (sidebar) | 9804–9826 | Same four methods, different purpose |
| Validation deficiencies | 9964–10027 | Pass V&V but fail validation; regression testing |
| **5.5 Product Transition** | 10130–10200 | Delivery to next level |

---

## Chapter 6: Crosscutting Technical Management

| Topic | Lines | Notes |
|---|---|---|
| **6.1 Technical Planning** | ~10500–11858 | SEMP, WBS, plans |
| Verification plan | 11410–11530 | V plan content |
| Validation plan | 11519–11651 | Val plan content |
| **6.2 Requirements Management** | 11860–12263 | CORE SECTION |
| Traceability definitions (sidebar) | 11875–11879 | Bidirectional definition |
| Requirements management activities | 12000–12053 | Hierarchical tree, bidirectional traceability |
| Requirements change management | 12063–12155 | CCB, impact assessment tools |
| Requirements creep | 12159–12197 | Prevention techniques |
| Requirements Management outputs | 12203–12263 | Documents, baselines, work products |
| **6.3 Interface Management** | 12273–12497 | CORE SECTION |
| Interface management purpose | 12273–12280 | Critical to successful projects |
| Interface management basic tasks | 12237–12297 | Define, identify, control |
| Interface management process activities | 12383–12479 | Design → Integration → Control |
| Interface documentation types | 12446–12479 | IRD, ICD, IDD, ICP |
| **6.4 Technical Risk Management** | ~12500–13100 | Risk process |
| **6.5 Configuration Management** | 13100–13551 | CORE SECTION |
| CM five elements | ~13100 area | Identification, change, CSA, audits, planning |
| CM plan structure | 13107–13133 | Internal and external uses |
| Configuration identification | 13127–13133 | CI selection and baselines |
| Four technical baselines | 13135–13300 | Functional/Allocated/Product/As-Deployed |
| Baseline evolution figure (FIGURE 6.5-3) | 13172–13265 | Timelines and review gates |
| Change control process | 13303–13413 | FIGURE 6.5-4; CCB workflow |
| Types of CM changes (sidebar) | 13421–13431 | Engineering change, waiver |
| Configuration status accounting | 13434–13492 | CSA; software tool criteria |
| CM outputs | 13509–13531 | CSA reports, current baselines |
| **6.6 Technical Data Management** | 13556–13905 | Data lifecycle |
| **6.7 Technical Assessment (Reviews)** | 13907–14392 | Reviews and metrics |
| Technical assessment purpose | 13907–14095 | Six criteria assessed |
| Life-cycle reviews table (TABLE 6.7-1) | 14395–14978 | All reviews: MCR through DRR |
| MCR | 14409–14437 | Mission Concept Review |
| SRR | 14439–14476 | System Requirements Review |
| MDR/SDR | 14477–14526 | Architecture validation |
| PDR | 14528–14567 | Preliminary design; interfaces identified |
| CDR | 14589–14627 | Final design; coding authorised |
| PRR | 14628–14669 | Production readiness |
| SIR | 14671–14712 | Integration readiness |
| SAR | 14714–14742 | System acceptance |
| ORR | 14757–14786 | Operational readiness |
| FRR | 14788–14812 | Flight readiness |
| **6.8 Decision Analysis** | ~14980–15800 | Trade studies and decision frameworks |

---

## Appendices

| Appendix | Topic | Lines |
|---|---|---|
| A | Acronyms | ~16100–16400 |
| B | Glossary (full) | 16405–18193 |
| B (selected) | Bidirectional traceability definition | 16493–16496 |
| B (selected) | Architecture (System) definition | 16498–16511 |
| B (selected) | Baseline definition | 16489–16492 |
| B (selected) | Configuration Items definition | 16588–16595 |
| B (selected) | ConOps definition | 16556–16577 |
| B (selected) | Requirements Management Process def | 17701–17704 |
| B (selected) | Stakeholder Expectations definition | 17798–17815 |
| B (selected) | SEMP definition | 17914–17921 |
| B (selected) | Technical Requirements Definition def | 17996–18002 |
| **C** | **How to Write a Good Requirement (Checklist)** | **18194–18410** |
| C.1 | shall/will/should terminology | 18196–18199 |
| C.2 | Editorial checklist | 18201–18233 |
| C.3 | General goodness checklist | 18237–18260 |
| C.4 | Requirements validation checklist | 18261–18410 |
| C.4 sections | Clarity, Completeness, Compliance, Consistency, Traceability, Correctness, Interfaces, Verifiability | 18261–18410 |
| **D** | **Requirements Verification Matrix** | **18413–18492** |
| E | Creating Validation Plan with Validation Requirements Matrix | ~18493–18600 |
| F | Functional, Timing, and State Analysis | ~19000 area |
| G | Technology Assessment/Insertion | ~19200 area |
| H | Integration Plan Outline | ~19800 area |
| I | Verification and Validation Plan Outline | ~20000 area |
| J | SEMP Content Outline | ~21000 area |
| K | Technical Plans | ~21500 area |
| L | Interface Requirements Document Outline | ~21700 area |
| M | CM Plan Outline | ~22000 area |
| N | Guidance on Technical Peer Reviews/Inspections | ~22200 area |
| R | HSI Plan Content Outline | ~23500 area |
| S | Concept of Operations Annotated Outline | ~24000 area |
| T | Systems Engineering in Phase E | ~24500 area |

---

## Quick Reference: Lines for Blueprint-Specific Topics

| Blueprint Need | Handbook Section | Line Range |
|---|---|---|
| Writing EARS-style requirements | Appendix C editorial checklist | 18196–18233 |
| Requirement metadata schema | TABLE 4.2-2 | 6553–6591 |
| Rationale field content | Rationale box | 6598–6612 |
| Bidirectional traceability | §6.2; definitions | 11860–11879, 12014–12053 |
| Requirement validation rules | §4.2.1.2.4; Appendix C.4 | 6616–6686, 18261–18410 |
| Functional decomposition steps | §4.3 | 6789–7062 |
| Interface definition and ICD structure | §6.3; App L | 12273–12497 |
| V&V method selection | §5.3, §5.4 methods sidebars | 9151–9176, 9804–9826 |
| Requirements-to-test traceability matrix | Appendix D | 18413–18492 |
| Configuration baseline types | §6.5.1.2.2 | 13135–13300 |
| CM tooling requirements | §6.5.1.2.4 | 13457–13492 |
| Technical reviews and gate criteria | TABLE 6.7-1 | 14395–14978 |
| Requirements creep prevention | §6.2.1.2.5 | 12159–12197 |
| ConOps purpose and content | §4.1.1.2.4 | 5800–5884 |
| System architecture definition | §4.3.1.2.1 | 6824–6913 |
| Cost of late changes | §2.5 FIGURE 2.5-1 | 1912–1935 |
| MOE/MOP/TPM definitions | §4.2.1.2.5; Appendix B | 6705–6764, ~17961–17983 |

---

*Index generated from full read of nasa_se_handbook.txt (25,273 lines). Line numbers are approximate due to text extraction artefacts; use ±50 lines when targeting a specific passage.*
