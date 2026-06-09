# Requirement Fidelity

Requirements cross functional boundaries. A compliance engineer needs to understand what applies to a product, its security characteristics, and whether evidence supports those characteristics. The development team governs that same product through engineering practices, business processes, and operational standards. The hard part is not writing a check. The hard part is preserving the full meaning of a requirement — including *why* it exists — as it moves between these groups through translation steps that were never designed to carry context.

This doc explores *why* fidelity degrades and where it degrades. It does not propose what ComplyTime should build.

## Where fidelity gets lost

The typical flow looks like this: During an audit, the compliance engineer receives a request for information from an auditor and relays that request to the development or operations team. The team responds with whatever they have — architecture diagrams, configuration screenshots, ticket histories, verbal explanations. The compliance engineer assembles this into a response.

The problem is making that information readily available in a consumable format. It arrives in fragments, across different tools and conversations, and none of it is structured for the question being asked. The compliance engineer becomes the human judge of the posture — interpreting, translating, and attesting based on what they were able to collect. Fidelity depends on that person's ability to bridge two contexts that do not naturally communicate.

The chain from authority document to assessment result has multiple translation steps, and each one can drop context:

```
Authority document
    → structured requirement
        → resolved policy (composed layers)
            → executable check
                → assessment result
```

### Compliance is not governance

Compliance demonstrates adherence to a set of requirements. Governance is the broader organizational management of decisions, risk, and operations. These are different functions — and conflating them makes both worse.

The [Gemara lexicon](https://gemara.openssf.org) draws this distinction explicitly: governance expresses intent; compliance verifies adherence. The [Automated Governance Maturity Model (AGMM)](https://tag-security.cncf.io/community/resources/automated-governance-maturity-model/) frames progressive maturity across policy, evaluation, enforcement, and audit. The distinction is emerging — GRC teams understand it, but the broader industry often conflates compliance tooling with governance posture.

### Requirements are multi-layered

Requirements flow from abstract to concrete:

| Layer                   | Example                                                                   | What it contributes                                 |
|:------------------------|:--------------------------------------------------------------------------|:----------------------------------------------------|
| Risk vectors            | Credential stuffing, lateral movement                                     | *Why* controls exist                                |
| Industry guidance       | OWASP Top 10, NIST CSF, CIS Controls                                      | Generalized recommendations                         |
| Technology controls     | "Require MFA for privileged access"                                       | Testable objectives                                 |
| Assessment requirements | "Verify OIDC MFA for admin ServiceAccounts"                               | Executable conditions against a specific data model |
| Organizational policy   | "Payments services require MFA; batch jobs may use compensating controls" | Risk appetite and scoping                           |
| Team overrides          | "Legacy batch job exempted per ticket SEC-1847"                           | Local context                                       |

Each layer carries rationale. The difference between "we implemented MFA to prevent credential stuffing against admin interfaces" and "we implemented MFA because we thought we should" is the difference between a requirement that teaches and one that merely constrains. That rationale is what gets lost when it does not travel with the requirement.

[OSCAL](https://pages.nist.gov/OSCAL/) has made significant strides in structuring the compliance artifact chain — catalogs, profiles with imports and modifications, component definitions, and assessment results. OSCAL profiles capture tailoring decisions and support programmatic composition. What remains harder is the cross-functional communication that produces those artifacts: the scoping conversation between a compliance engineer and a developer, the risk acceptance documented in a ticket comment, the compensating control approved by email.

### Natural language does not compile

Authority documents describe requirements in natural language. Assessment tooling needs precise, executable conditions. The translation is lossy — and that is a property of the domain, not a tooling bug.

Different authors interpret the same requirement differently. "Encrypt data at rest" becomes an S3 bucket check for one author, an RDS check for another, and an EBS volume check for a third. All defensible. None covers the full requirement. "Maintain an incident response plan" becomes a check for document existence, document currency, or tabletop exercise completion depending on who writes it.

### Non-technical controls are usually excluded

Most automated evaluation tooling focuses on technical controls — the subset directly observable from system state. Procedural, physical, and operational controls are not second-class. They are usually excluded entirely.

| Control category | Example                             | Typical evaluation                    |
|:-----------------|:------------------------------------|:--------------------------------------|
| **Technical**    | Containers do not run as root       | Configuration scan, admission control |
| **Procedural**   | Change management process exists    | Document review, attestation          |
| **Physical**     | Badge access with audit logging     | Physical inspection                   |
| **Operational**  | Backup restoration tested quarterly | Process verification, test records    |

Self-attestation handles some of this, but self-attestation across trust boundaries is increasingly weak as evidence. It is also taxing to maintain. The scaling problem does not end with technical controls. GRC platforms typically handle attestation workflows or allow uploading opaque documents, but the result is a PDF that proves something was uploaded, not that a process was followed.

### Policy-as-code is assessment logic, not governance

Two distinct activities under the word "policy" are often conflated:

1. **Governance policy** — organizational decisions about what to achieve, what risk to accept, and why.
2. **Assessment policy** — executable logic that checks whether observed state aligns with a requirement.

Policy-as-code is assessment logic. It answers "does this resource violate this rule?" It does not answer "why does this rule exist?" or "should this service be in scope?"
Organizations express intent through governance. They do not write Rego. The [Gemara model](https://gemara.openssf.org/model/) separates these layers explicitly: Layer 3 (Organizational Policy) expresses intent; Layer 5 (Intent & Behavior Evaluation) executes checks. Fidelity requires a traceable link between them.

## Open questions

- What is the minimum representation that preserves both precision and rationale as requirements cross functional boundaries?
- How should requirements that resist full automation be represented alongside automatable ones?
- Can rationale be structured without becoming boilerplate? "We implemented MFA because credential stuffing" is useful. "We implemented MFA because control AC-2 requires it" is circular.
- Does the distinction between governance policy and assessment policy need to be enforced architecturally, or can a single artifact carry both?
