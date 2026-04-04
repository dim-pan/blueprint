# Index: ears_tutorial.txt — Topic to Line Range Mapping

Source file: `/Users/dimitriospanagiotopoulos/Raylight_Dev/blueprint/resources/ears_tutorial.txt`
Total lines: 1195

Line numbers reference the `cat -n` output of the source file.

---

## Document Structure

| Section | Lines | Notes |
|---|---|---|
| Title / Author / Conference metadata | 1–11 | John Terzakis, Intel Corporation, ICCGI 2013, Nice, France |
| Legal / Copyright notices (front matter) | 12–26 | Repeated throughout; not indexed exhaustively below |
| Agenda | 27–37 | Eight agenda items listed |
| Requirements Overview (section header) | 44–50 | Slide header |
| What is a Requirement? | 52–62 | Three definitions; functional vs. non-functional |
| Functional/Non-Functional Examples (Video over IP) | 68–88 | Conference calling feature list |
| Issues with Requirements (section header) | 94–100 | Slide header |
| Issues with Requirements (content) | 102–115 | Eight root problems listed; copy-paste propagation |
| Examples of Issues with Requirements | 121–149 | Four anti-pattern examples with bullet commentary |
| Overcoming Issues with Requirements (section header) | 151–157 | Slide header |
| Tools to Overcome Requirements Issues | 159–171 | Five tools: training, syntax, checklist, Planguage, EARS |
| Training (detail) | 173–187 | Internal training, communities of practice |
| Consistent Requirements Syntax | 189–207 | Generic syntax template; Trigger/Precondition/Actor/Action/Object |
| Good Requirements Checklist | 209–228 | Ten attributes listed |
| Planguage: A Constrained Natural Language | 234–256 | Background, Tom Gilb, 1988, Competitive Engineering book |
| Planguage Keywords for Any Requirement | 262–296 | Keyword table: Name, Requirement, Rationale, Priority, Status, Contact, Source |
| Planguage Example | 302–317 | Create_Invoice example with all keyword fields filled |
| Express Requirements Using EARS (intro) | 323–337 | EARS definition; Rolls-Royce origin; five pattern names listed |
| EARS Background | 343–367 | Rolls-Royce case study; safety-critical aircraft; eight problems; RE 09 citation |
| Identifying Ubiquitous Requirements (section header) | 374–380 | Slide header |
| Ubiquitous Requirements (definition) | 382–395 | Three properties; contrast with non-ubiquitous |
| About Ubiquitous Requirements | 401–413 | Warning: question apparently ubiquitous requirements; two false-ubiquitous examples |
| Ubiquitous or Not? (quiz slide) | 420–441 | Six candidate requirements presented for classification |
| Ubiquitous or Not? (answers) | 448–481 | Answers with explanations for all six |

---

## EARS Patterns Reference

| Section | Lines | Notes |
|---|---|---|
| EARS and EARS Examples (section header) | 486–492 | Slide header |
| EARS Patterns Table (all five + complex) | 494–528 | Master pattern reference with templates |
| Ubiquitous Requirements (pattern detail) | 535–542 | No keyword; fundamental property; format |
| Ubiquitous Examples (3 examples) | 549–555 | Installer, Java, website availability |
| Event-Driven Requirements (pattern detail) | 562–570 | WHEN keyword; trigger + optional precondition |
| Event-Driven Examples (3 examples) | 578–589 | USB driver, DVD spin-up, water valve |
| Unwanted Behavior Requirements (pattern detail) | 597–606 | IF/THEN keywords; errors, faults, failures |
| Unwanted Behavior Examples (3 examples) | 613–623 | Speed variance, memory checksum, ATM card |
| State-Driven Requirements (pattern detail) | 630–636 | WHILE keyword; "During" also acceptable |
| State-Driven Examples (3 examples) | 644–654 | Low Power Mode, heater, autopilot |
| Optional Feature Requirements (pattern detail) | 661–669 | WHERE keyword; product variants |
| Optional Feature Examples (3 examples) | 677–688 | Thesaurus, hardware encryption, HDMI port |
| Complex Requirements (pattern detail) | 695–706 | Combinations of keywords; multiple conditions |
| Complex Examples (3 examples) | 713–728 | Landing gear alarm, optical drive copy, startup flash card |

---

## Rewriting Requirements Using EARS

| Section | Lines | Notes |
|---|---|---|
| Section header | 734–741 | Slide header |
| Example 1: Installer in Greek (Ubiquitous, no change) | 743–750 | No EARS rewrite needed |
| Example 2: Participant count (Event-Driven) | 757–766 | Adds trigger: "when user selects caller count from menu" |
| Example 3: Phone Alarm Company (Unwanted Behavior) | 774–783 | Adds fault condition: sensor malfunction |
| Example 4: Mute microphone (State-Driven) | 790–797 | Adds state: "while mute button is depressed" |
| Example 5: Download book without charge (Optional Feature) | 805–813 | Adds optional condition: digital format; adds 3-day trial period |
| Example 6: Low battery warning (Complex) | 821–829 | Combines WHILE (battery power) + IF/THEN (below 10%) |

---

## EARS at Intel

| Section | Lines | Notes |
|---|---|---|
| EARS at Intel (section header) | 837–843 | Slide header |
| History of EARS at Intel | 845–858 | Introduced 2010; adopted across diverse project types |
| Intel EARS Examples — Set 1 | 866–878 | Ubiquitous (online help), Event-Driven (J7 jumper), Unwanted Behavior (DRAM config) |
| Intel EARS Examples — Set 2 | 885–899 | State-Driven (Manufacturing Mode), Optional Feature (3G/Wi-Fi), Complex (DC power + error) |

---

## Wrap-up and Summary

| Section | Lines | Notes |
|---|---|---|
| Wrap up (section header) | 906–912 | Slide header |
| Session Summary | 914–928 | Seven bullet recap of session content |
| Final Thoughts | 934–948 | Three key takeaways; closing call to action |
| Contact Information | 955–962 | John Terzakis, john.terzakis@intel.com |

---

## Backup / Appendix

| Section | Lines | Notes |
|---|---|---|
| Backup (section header) | 969–975 | Slide header |
| Papers on EARS (references) | 977–984 | RE 09 and RE 10 citations |
| Complete (attribute definition + example) | 991–1003 | TBD anti-pattern; login attempt example |
| Correct (attribute definition + example) | 1010–1021 | 802.3 Ethernet frame byte count |
| Concise (attribute definition + example) | 1029–1042 | Wordy vs. clean software display example |
| Feasible (attribute definition + example) | 1050–1060 | "Unlimited concurrent users" anti-pattern |
| Necessary (attribute definition + example) | 1068–1082 | "All prior Windows versions" anti-pattern |
| Prioritized (attribute definition + example) | 1089–1104 | High/Medium/Low scale; "all critical" anti-pattern |
| Unambiguous (attribute definition + example) | 1111–1124 | "Install quickly" anti-pattern; timed installation example |
| Verifiable (attribute definition + example) | 1132–1144 | "Easy to find" anti-pattern; path-based example |
| Consistent (attribute definition + example) | 1152–1166 | Conflicting whole-numbers vs. decimal-time requirements |
| Traceable (attribute definition + example) | 1174–1187 | Tag-based identification; Prompt_PIN example |
