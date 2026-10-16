# ADR-0007: Ecosystem Compatibility Model

| Field | Value |
|:------|:------|
| Status | proposed |
| Date | 2026-09-29 |
| Decision-makers | @gxmiranda, @jflowers, @jpower432, @marcusburghardt |

## Context and Problem Statement

The complytime ecosystem spans multiple independently-released repositories (complyctl, complytime-providers, complytime-policies, complypack, org-infra). Stakeholder capabilities emerge from the composition of compatible component versions, but no single version number answers "do these work together?" Stakeholders must reverse-engineer compatibility from repository tags, changelogs, and go.mod files. See the [Release Communication](../problems/release-communication.md) problem doc for the full exploration.

The ecosystem needs a simple, low-maintenance mechanism to ensure and communicate component compatibility.

## Decision Drivers

* complyctl defines every API contract other components implement — the `Provider` interface in `pkg/provider`, gRPC plugin protocol in `api/plugin`, and content schema consumed by complypack. This makes complyctl the natural compatibility reference point.
* Compatibility is already mechanically verifiable — Go modules enforce compile-time compatibility, and cross-repo integration CI validates runtime compatibility.
* The primary audience consuming the ecosystem directly is technical. Non-technical stakeholders are concerned with outcomes, not specific version numbers.
* Curated compatibility matrices and release ceremonies carry ongoing maintenance overhead that has not been sustainable to adopt.

## Considered Options

* complyctl as compatibility anchor with CI verification
* Curated compatibility matrix with capability maturity tiers
* No ecosystem-level compatibility mechanism (status quo)

## Decision Outcome

Chosen option: "complyctl as compatibility anchor with CI verification", because it reflects the existing architectural reality — complyctl owns the contracts — and leverages CI for continuous, accurate compatibility verification with zero curation overhead.

### Core principles

**complyctl owns the API contracts.** All ecosystem components either implement or consume contracts defined by complyctl: the `Provider` interface (`pkg/provider`), the gRPC plugin protocol (`api/plugin`), and the content schema for ComplyPacks. Compatibility with the ecosystem reduces to compatibility with complyctl.

**Compatibility is verified continuously in CI, not curated in a published matrix.** complyctl's cross-repo integration CI builds providers against complyctl on every PR. The reverse direction — providers validating against complyctl's latest release — is tracked in [complytime-providers#187](https://github.com/complytime/complytime-providers/issues/187). If CI is green, the combination works.

**Each repository owns its own breaking change detection.** complyctl's v1.x API stability guarantee is enforced by automated CI guardrails as tracked in [complyctl#880](https://github.com/complytime/complyctl/issues/880). Other repositories validate against complyctl's public API in their own CI pipelines.

**Per-provider maturity is communicated in each provider's own documentation.** Provider readiness (e.g., OpenSCAP=stable, OPA=in-development) is a property of each provider, not of the ecosystem. Each provider's documentation or `Describe` RPC metadata is the authoritative source for its maturity status.

### Confirmation

The ADR is confirmed when:
* Cross-repo integration CI exists in both directions — complyctl validates providers (exists) and providers validate against complyctl (tracked in [complytime-providers#187](https://github.com/complytime/complytime-providers/issues/187)).
* Breaking change CI guardrails are active in complyctl (tracked in [complyctl#880](https://github.com/complytime/complyctl/issues/880)).

## Consequences

* Good, because zero maintenance overhead — no compatibility matrix, no release ceremony, no maturity tier taxonomy to maintain.
* Good, because CI results are always accurate — they reflect current reality, not a point-in-time curated snapshot that can lag behind releases.
* Good, because clear ownership — complyctl owns contracts, each repository owns its own enforcement and validation.
* Good, because it aligns with how the architecture already works — complyctl is already the hub in the [component map](../architecture.md), this ADR makes that role explicit for compatibility.
* Bad, because there is no single human-readable page answering "what works with what" — stakeholders must understand that CI green means compatible, or check go.mod dependency versions.
* Bad, because per-provider maturity is decentralized — a stakeholder evaluating multiple providers must consult each provider's documentation rather than one matrix.

## Pros and Cons of the Options

### complyctl as compatibility anchor with CI verification

Designate complyctl as the compatibility reference point. Verify compatibility continuously through cross-repo integration CI. Each repository owns its own breaking change detection.

* Good, because it matches the architectural reality — complyctl already defines every contract.
* Good, because CI verification is continuous, automatic, and cannot drift from reality.
* Good, because it requires no new processes, artifacts, or ceremonies to maintain.
* Neutral, because it requires CI literacy to answer compatibility questions — acceptable for the current technical audience.
* Bad, because if a downstream non-technical audience emerges, a human-readable compatibility artifact would need to be layered on top.

### Curated compatibility matrix with capability maturity tiers

Publish a capability-oriented compatibility matrix mapping stakeholder capabilities to validated component combinations with four maturity tiers (Alpha, Beta, Pre-GA, GA). Maintain an ecosystem release process with tagging and release notes.

* Good, because it provides a single human-readable page showing what works together.
* Good, because maturity tiers communicate per-provider readiness explicitly.
* Bad, because the matrix requires manual curation at each ecosystem release — maintenance overhead that has not been adopted.
* Bad, because maturity tier labels (Alpha, Beta) risk the same semiotic collision with Go semver pre-release suffixes that motivated the original problem.
* Bad, because the ecosystem release ceremony adds process without clear adoption evidence.

### No ecosystem-level compatibility mechanism (status quo)

Each repository releases independently. Compatibility is left to stakeholders to determine from go.mod files and changelogs.

* Good, because zero overhead — no new processes or artifacts.
* Bad, because stakeholders must reverse-engineer compatibility from multiple repository tags and dependency files.
* Bad, because there is no explicit contract ownership — implicit knowledge that complyctl is the hub is not documented.

## More Information

* Problem doc: [Release Communication](../problems/release-communication.md)
* complyctl breaking change CI: [complyctl#880](https://github.com/complytime/complyctl/issues/880)
* Provider-side compatibility CI: [complytime-providers#187](https://github.com/complytime/complytime-providers/issues/187)
* This ADR supersedes the approach proposed in [PR #20](https://github.com/complytime/complytime/pull/20), which bundled a compatibility matrix, breaking change policy, ecosystem release process, and capability maturity model. That broader approach has been simplified to this single architectural principle.
