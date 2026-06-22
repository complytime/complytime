# ADR-0005: Two-Stream Content Model

**Status:** accepted

**Date:** 2026-05-07 (retroactive — formalized in CEP-0001)

**Deciders:** @jpower432

## Context

Compliance assessment involves two distinct concerns with different authorship boundaries:

1. **What must be true** — control objectives, assessment requirements, policies. Authored by compliance/GRC teams against authority documents (NIST, CIS, organizational standards).
2. **How to verify it** — executable checks (Rego rules, CUE constraints, CEL expressions). Authored by engineering teams who understand the target systems.

Bundling both in a single artifact conflates authorship. A compliance engineer updating a control description should not trigger re-testing of assessment logic. An engineer updating a Rego rule should not require compliance team sign-off on content they didn't change.

PR #3 review raised the "fox guarding the henhouse" concern: if one team controls both the definition of compliance and its verification, separation of duties is violated.

## Decision

Two independent OCI artifact streams:
1. **Gemara content** — compliance-authored (`#Policy`, `#ControlCatalog`, `#GuidanceCatalog`). Versioned against Gemara CUE schemas.
2. **ComplyPacks** — engineering-authored assessment logic. Versioned against provider expectations.

No artifact-level dependency exists between the streams. Compatibility is expressed through provenance metadata (source policy ID, requirement IDs in traceability annotations) and validated at runtime by the provider.

Each stream has independent publishing, signing, and lifecycle. `complyctl get` pulls both; the config media type discriminates between them.

## Consequences

- Clear authorship boundaries — compliance and engineering teams work independently.
- Independent release cadences — updating a control catalog does not require republishing assessment logic.
- Provenance metadata enables traceability without coupling.
- Runtime must resolve compatibility (schema version alignment) rather than relying on co-bundled artifacts.
- More artifacts to manage than a monolithic bundle — operational complexity for operators pulling both streams.

**Source:** [CEP-0001](https://github.com/complytime/complytime/blob/main/ceps/cep-0001-complypack-architecture.md), [PR #3 review discussion](https://github.com/complytime/complytime/pull/3)
