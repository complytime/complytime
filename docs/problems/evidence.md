# Evidence

Compliance evidence is fragmented across automated and manual sources with no unified model for collection, normalization, freshness, or trust. Every assessment tool produces its own artifact format. Manual evidence (attestations, policy review records, physical inspection reports) lives outside automated systems entirely. No standard way exists to trace any piece of evidence back to the requirement it satisfies.

This doc explores *why* evidence unification is hard and what a solution needs to address. It does not propose what ComplyTime should build. See [evaluator-coupling.md](evaluator-coupling.md) for how collection and evaluation conflation affects evidence: evaluators currently discard raw system state after assessment, making evidence reproduction dependent on re-scanning. See [requirement-fidelity.md](requirement-fidelity.md) for how fidelity loss weakens the link between evidence and requirements. See [system-modeling-and-drift.md](system-modeling-and-drift.md) for how drift measurement against requirements interacts with evidence freshness.

## Why this is hard

### Automated and manual evidence are fundamentally different

Automated evidence is continuous, reproducible, and machine-generated. A configuration scan runs hourly and produces structured findings. Manual evidence is different in kind: cycle-based, asserted by humans, and often document-form (screenshots, signed attestations, meeting minutes, physical inspection reports). Forcing manual evidence into automated schemas loses context. Keeping them separate breaks unified posture reporting.

### Freshness semantics differ

A scan result expires in hours or days because the system state it describes changes constantly. An annual policy review is valid for 12 months. A quarterly backup restoration test, 90 days. Each evidence type has its own validity window, and those windows differ by orders of magnitude. Any system that treats all evidence as equally current or equally stale will misrepresent compliance posture.

### Trust properties differ

Automated findings can be re-run for verification. The same scan on the same system state produces the same result. Manual attestations rely on the authority and process of the attestor. Re-running is not meaningful in the same way. An attestation's trust comes from who signed it, under what authority, and following what process. These are fundamentally different trust models that a unified evidence system must surface without collapsing the distinction.

### Compensating controls are evidence

A compensating control satisfies a requirement through alternative means. It is not a separate concept from evidence; it produces evidence that links back to the original requirement. The evidence envelope must capture the compensating nature and the original requirement linkage without inventing a parallel evidence system.

## Current approaches / prior art

- **OSCAL assessment results:** Single schema for assessment output, but manual evidence fits awkwardly. The observations-vs-findings distinction captures some of the manual/automated difference, but the schema was designed for machine-generated results.
- **GRC platforms:** Separate manual evidence workflows with poor integration with automated scanning. Evidence lives in two systems with manual reconciliation.
- **SOC 2 Type II:** Auditors distinguish "test of one" (automated, reproducible) from "inquiry and observation" (manual, point-in-time). The audit profession has established vocabulary for evidence types that predates automated compliance tooling.

## Open questions

1. What is the minimum evidence envelope schema that spans automated findings, manual attestations, and compensating controls without over-generalizing?
2. How does evidence expiration interact with continuous monitoring? Does stale evidence trigger re-assessment, or just a warning?
3. What provenance metadata is needed for manual attestations to satisfy auditor expectations (SOC 2, ISO 27001)?
4. How should the evidence platform handle conflicting evidence, where two sources disagree about the same requirement?
