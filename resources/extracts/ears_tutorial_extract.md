# EARS: Easy Approach to Requirements Syntax — Curated Extract

Source: `ears_tutorial.txt` (John Terzakis, Intel Corporation, ICCGI Conference, July 21, 2013)
Original work by Alistair Mavin et al., Rolls-Royce PLC (RE 09, 2009)

---

## 1. What is a Requirement?

A requirement is a statement of one of the following:

1. What a system must do (Functional Requirements)
2. A known limitation or constraint on resources or design (Non-Functional)
3. How well the system must do what it does (Non-Functional)

EARS primarily addresses functional requirements (definitions 1 and 2).

---

## 2. Common Problems with Natural Language Requirements

Requirements written in unconstrained natural language suffer from eight major problem types (identified by Rolls-Royce in the EARS origin study):

| Problem | Description |
|---|---|
| Ambiguity | Multiple possible interpretations |
| Vagueness | Imprecise terms ("about five", "quickly", "support") |
| Complexity | Overly convoluted grammar or logic |
| Omission | Missing triggers, preconditions, or "else" branches |
| Duplication | Same requirement stated in multiple places |
| Wordiness | Excess words that obscure meaning |
| Inappropriate implementation | Requirements that prescribe design rather than behavior |
| Untestability | Statements that cannot be verified |

### Anti-Pattern Examples

**Vague verb:**
> "The software shall support a water level sensor."
What does "support" mean? The requirement is unverifiable.

**Imprecise quantity + missing trigger:**
> "The thesaurus software shall display about five alternatives for the requested word."
"About five" is ambiguous; no trigger specified for when alternatives are shown.

**Missing trigger:**
> "The software shall blink the LED on the adapter using a 50% on, 50% off duty cycle."
Does blinking happen at all times? There is no stated trigger.

**Incomplete logic (if without else):**
> "If a boot disk is detected in the system, the software shall boot from it."
What happens when no boot disk is present? The requirement is logically incomplete.

---

## 3. EARS Background and Motivation

EARS was created by Alistair Mavin and colleagues at Rolls-Royce PLC and presented at the 17th IEEE International Requirements Engineering Conference (RE 09) in 2009. The case study involved aircraft engine control software: safety-critical, thousands of components, up to twenty suppliers.

Rewriting requirements using EARS demonstrated a significant reduction in all eight problem types listed above.

EARS was introduced at Intel in 2010. It was praised by developers and validation teams for providing clarity and removing ambiguity, and was rapidly adopted due to its power and simplicity.

---

## 4. The Five EARS Patterns

### Pattern Reference Table

| Pattern | Keyword(s) | Template |
|---|---|---|
| Ubiquitous | (none) | `The <system name> shall <system response>` |
| Event-Driven | WHEN | `WHEN <trigger> [<optional precondition>] the <system name> shall <system response>` |
| Unwanted Behavior | IF / THEN | `IF <unwanted condition or event>, THEN the <system name> shall <system response>` |
| State-Driven | WHILE | `WHILE <system state>, the <system name> shall <system response>` |
| Optional Feature | WHERE | `WHERE <feature is included>, the <system name> shall <system response>` |
| Complex | combinations | `<Multiple Conditions>, the <system name> shall <system response>` |

---

### 4.1 Ubiquitous Requirements

**Definition:** Statements of fundamental system properties that hold at all times. They require no trigger, precondition, or stimulus in order to execute.

**Template:**
```
The <system name> shall <system response>
```

**When to use:** Only when the behavior truly exists at all times, unconditionally. Most requirements that appear ubiquitous are actually event-driven, state-driven, or optional — and must be questioned carefully.

**Legitimate ubiquitous examples:**
- The software package shall include an installer.
- The software shall be written in Java.
- The software shall be available for purchase on the company web site and in retail stores.
- The installer software shall be available in Greek.
- The software shall include an online help file. *(Intel)*

**How to identify genuine ubiquitous requirements:**
Ask: "Is there any state, event, or condition under which this behavior would NOT apply?" If yes, the requirement is not ubiquitous and needs a keyword.

**False ubiquitous examples (and what they really are):**
- "The software shall wake the PC from standby" — missing a trigger (event-driven)
- "The software shall log the date, time and username of failed logins" — missing a trigger (event-driven)
- "The software shall mute the microphone" — occurs during a state (state-driven)
- "The software shall phone the Alarm Company" — only on a fault condition (unwanted behavior)
- "The software shall download the book without charge" — only when an optional condition applies (optional feature)
- "The software shall warn the user of low battery" — requires a state and an event (complex)

---

### 4.2 Event-Driven Requirements

**Definition:** Behaviors that are initiated when and only when a specific trigger occurs or is detected.

**Template:**
```
WHEN <trigger> [<optional precondition>] the <system name> shall <system response>
```

**When to use:** When a discrete, detectable event initiates the behavior. The precondition is optional and further qualifies when the response applies.

**Examples:**
- When an Unregistered Device is plugged into a USB port, the OS shall attempt to locate and load the driver for the device.
- When a DVD is inserted into the DVD player, the OS shall spin up the optical drive.
- When the water level falls below the Low Water Threshold, the software shall open the water valve to fill the tank to the High Water Threshold.
- When the user selects the caller count from the menu, the software shall display a count of the number of participants in the audio call in the UI. *(rewrite of vague ubiquitous)*
- When the software detects the J7 jumper is shorted, it shall clear all stored user names and passwords. *(Intel)*

---

### 4.3 Unwanted Behavior Requirements

**Definition:** Requirements that handle error conditions, failures, faults, disturbances, and other undesired events.

**Template:**
```
IF <unwanted condition or event>, THEN the <system name> shall <system response>
```

**When to use:** When the triggering condition is undesired, abnormal, or a fault — as opposed to a normal operational event. The IF/THEN construction makes the conditional logic explicit and complete.

**Examples:**
- If the measured and calculated speeds vary by more than 10%, then the software shall use the measured speed.
- If the memory checksum is invalid, then the software shall display an error message.
- If the ATM card inserted is reported lost or stolen, then software shall confiscate the card.
- If the alarm software detects that a sensor has malfunctioned, then the alarm software shall phone the Alarm Company to report the malfunction. *(rewrite of trigger-less ubiquitous)*
- If the software detects an invalid DRAM memory configuration, then it shall abort the test and report the error. *(Intel)*

---

### 4.4 State-Driven Requirements

**Definition:** Behaviors that are active continuously while the system is in a particular state.

**Template:**
```
WHILE <system state>, the <system name> shall <system response>
```

**When to use:** When the system occupies a named, persistent mode or condition and a behavior must be maintained throughout that period. The keyword "During" is an acceptable alternative to "While".

**Examples:**
- While in Low Power Mode, the software shall keep the display brightness at the Minimum Level.
- While the heater is on, the software shall close the water intake valve.
- While the autopilot is engaged, the software shall display a visual indication to the pilot.
- While the mute button is depressed, the software shall mute the microphone. *(rewrite of missing-trigger ubiquitous)*
- While in Manufacturing Mode, the software shall boot without user intervention. *(Intel)*

---

### 4.5 Optional Feature Requirements

**Definition:** Behaviors that apply only in systems or configurations that include a specific optional feature or capability.

**Template:**
```
WHERE <feature is included>, the <system name> shall <system response>
```

**When to use:** When a product family has variants and certain behaviors only apply to configurations that include a particular hardware or software feature.

**Examples:**
- Where a thesaurus is part of the software package, the installer shall prompt the user before installing the thesaurus.
- Where hardware encryption is installed, the software shall encrypt data using the hardware instead of using a software algorithm.
- Where a HDMI port is present, the software shall allow the user to select HD content for viewing.
- Where the book is available in digital format, the software shall allow the user to download the book without charge for a trial period of 3 days. *(rewrite of vague ubiquitous)*
- Where both 3G and Wi-Fi radios are available, the software shall prioritize Wi-Fi connections above 3G. *(Intel)*

---

## 5. Complex / Compound Requirements

**Definition:** Requirements that involve multiple conditions — combinations of triggers, states, optional features, and/or unwanted behaviors — that together govern a single system response.

**Template:**
```
<Multiple Conditions>, the <system name> shall <system response>
```

Keywords WHEN, IF/THEN, WHILE, and WHERE are combined as needed.

**Examples:**

*Event inside a state (WHILE + WHEN):*
> While in start up mode, when the software detects an external flash card, the software shall use the external flash card to store photos.

*Event followed by fault detection (WHEN + IF/THEN):*
> When the landing gear button is depressed once, if the software detects that the landing gear does not lock into position, then the software shall sound an alarm.

*Optional feature with event (WHERE + WHEN):*
> Where a second optical drive is installed, when the user selects to copy disks, the software shall display an option to copy directly from one optical drive to the other optical drive.

*State with fault condition (WHILE + IF/THEN):*
> While on battery power, if the battery charge falls below 10% remaining, then the system shall display a warning message to switch to AC power. *(rewrite of vague ubiquitous)*
>
> While on DC power, if the software detects an error, then the software shall cache the error message instead of writing the error message to disk. *(Intel)*

---

## 6. Converting Existing Requirements to EARS Format

The process is:

1. **Read the existing requirement** and determine whether it is truly ubiquitous.
   - Ask: does this behavior occur at all times, unconditionally?
   - If yes: apply the Ubiquitous template (no keyword needed; often no change required).
   - If no: proceed to the next steps.

2. **Identify the governing condition type:**
   - Is there a discrete trigger event? Use WHEN.
   - Is the condition undesired, a fault, or an error? Use IF/THEN.
   - Is the behavior tied to a persistent system mode or state? Use WHILE.
   - Does the behavior only apply when an optional feature is present? Use WHERE.
   - Are multiple condition types involved? Combine keywords (Complex pattern).

3. **Fill in the template precisely:**
   - Name the system explicitly (`the software`, `the OS`, `the alarm software`).
   - State the trigger or condition specifically (avoid vague terms like "about", "support", "quickly").
   - State the system response with a testable, concrete outcome.
   - Add quantities, thresholds, and timeframes where necessary to make the requirement verifiable.

4. **Check completeness of logic:**
   - Does the requirement address what happens in the absence of the stated condition? If incomplete logic is present, add a complementary requirement or extend the current one.

### Before/After Conversion Examples

| Original (Anti-Pattern) | EARS Rewrite | Pattern Applied |
|---|---|---|
| The installer software shall be available in Greek. | The installer software shall be available in Greek. | Ubiquitous (no change) |
| The software shall display a count of the number of participants. | When the user selects the caller count from the menu, the software shall display a count of the number of participants in the audio call in the UI. | Event-Driven |
| The software shall phone the Alarm Company. | If the alarm software detects that a sensor has malfunctioned, then the alarm software shall phone the Alarm Company to report the malfunction. | Unwanted Behavior |
| The software shall mute the microphone. | While the mute button is depressed, the software shall mute the microphone. | State-Driven |
| The software shall download the book without charge. | Where the book is available in digital format, the software shall allow the user to download the book without charge for a trial period of 3 days. | Optional Feature |
| The software shall warn the user of low battery. | While on battery power, if the battery charge falls below 10% remaining, then the system shall display a warning message to switch to AC power. | Complex (State + Unwanted Behavior) |

---

## 7. Good Requirements Attributes (Checklist)

EARS directly supports many of these attributes. Each requirement should be:

| Attribute | Definition |
|---|---|
| Complete | Contains sufficient detail for designers and developers; no TBD values |
| Correct | Error-free; verifiable against source materials and SMEs |
| Concise | No embedded rationale, examples, or compound statements; uses fewest words necessary |
| Feasible | At least one valid design/implementation exists |
| Necessary | Traceable to a customer need, stakeholder, business strategy, or differentiator |
| Prioritized | Ranked or ordered relative to other requirements (e.g., High/Medium/Low) |
| Unambiguous | Single interpretation; terms defined; tested against target audience |
| Verifiable | Provably implemented via demonstration, analysis, inspection, or test |
| Consistent | Does not conflict with any other requirement at any level |
| Traceable | Uniquely and persistently identified with a tag (e.g., `Prompt_PIN: The software shall...`) |

---

## 8. EARS in the Context of a Requirements Framework

EARS integrates naturally with structured requirement frameworks such as Planguage. The `Requirement` field in a Planguage record follows EARS syntax. Example:

```
Name: Create_Invoice
Requirement: When an Order is shipped and Order Terms are not "Prepaid",
             the system shall create an Invoice.
Rationale: Task automation decreases error rate, reduces effort per order.
Priority: High.
Status: Committed
Contact: Hugh P. Essen
Source: I. Send, Shipping
```

---

## 9. Key Takeaways

- EARS is a structured aid, not a grammar checker. It works by forcing authors to identify the governing condition type before writing the requirement.
- Most requirements are NOT ubiquitous. Treating non-ubiquitous requirements as ubiquitous is the most common EARS anti-pattern.
- The keywords (WHEN, IF/THEN, WHILE, WHERE) are the diagnostic tool: selecting the right keyword forces the author to articulate the trigger, state, or feature that was previously implicit or missing.
- Complex requirements are not a separate category to avoid — they are the correct form when multiple conditions legitimately govern a single behavior.
- EARS is beneficial to both developers and testers: removing ambiguity and properly identifying conditions clarifies both implementation intent and test criteria.

---

## 10. Primary References

- Mavin, A. et al. "EARS (Easy Approach to Requirements Syntax)." 17th IEEE International Requirements Engineering Conference (RE 09), 2009, p. 321.
- Mavin, A. et al. "Big EARS: The Return of Easy Approach to Requirements Syntax." 18th IEEE International Requirements Engineering Conference (RE 10), 2010.
- Terzakis, J. "EARS: The Easy Approach to Requirements Syntax." Intel Corporation, ICCGI Conference, Nice, France, July 21, 2013.
