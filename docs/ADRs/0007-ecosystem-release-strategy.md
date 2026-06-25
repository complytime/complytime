# ADR-0007: Ecosystem Release Strategy

**Status:** accepted

**Date:** 2026-06-25

**Deciders:** @gvauter, @jflowers, @jpower432, @marcusburghardt

## Context

The complytime ecosystem spans multiple independently-released repositories (complyctl, complytime-providers, complytime-policies, complypack, org-infra, complytime). Three problems motivate this decision:

1. Stakeholders misread Go semver pre-release suffixes (e.g., `v1.0.0-beta.0`) as product immaturity signals. In Go, these suffixes signal API stability commitments — not product readiness — but this distinction is lost on audiences unfamiliar with Go module conventions.
2. No ecosystem-level coordination point exists. Each repository releases independently, and no single artifact tells stakeholders "these versions work together."
3. Breaking change risks lack automated CI enforcement beyond existing protobuf checks.

complyctl's public API surface is narrow and well-guarded: the `Provider` interface in `pkg/provider`, ~15 domain types, and generated gRPC in `api/plugin`. The `internal/` boundary protects most code from external import. Existing safeguards include the optional interface pattern, domain type insulation from protobuf, Buf breaking change detection, and frozen handshake values. The API is stable and the team is ready to proceed to v1.0.0.

## Decision

1. **Proceed with complyctl v1.0.0**, accepting the API stability commitment. The public API surface is narrow, well-guarded, and no breaking changes are expected.
2. **Adopt a four-tier capability maturity model** (Alpha → Beta → Pre-GA → GA) decoupled from individual component version numbers. These terms are standard in cloud-native ecosystems and do not collide with Go semver pre-release suffixes.
3. **Establish an ecosystem release process** triggered from the complytime repository. The compatibility matrix updates only when maintainers deliberately trigger the process after verifying integration compatibility — it is a curated sign-off, not a live dashboard.
4. **Publish a capability-oriented compatibility matrix** on complytime.dev, rendered via Docsify alongside existing documentation.
5. **Implement three-layer breaking change CI enforcement:**

   | Layer | Tool | Catches |
   |---|---|---|
   | Protobuf wire | Buf breaking | Wire protocol schema changes (exists) |
   | Go API surface | apidiff | Exported type/function/interface changes in `pkg/provider` and `api/plugin` (to implement) |
   | Integration wire | Custom CI job | Runtime incompatibilities between compiled provider binaries and complyctl (to implement) |

## Consequences

- Stakeholder communication improves — capability maturity levels convey production readiness without requiring Go semver literacy.
- Ecosystem coordination has a single home — the complytime repository becomes the release coordination point with a validated compatibility matrix.
- Breaking change detection covers three complementary layers, catching failures that any single layer would miss.
- Release highlights are captured at merge time via `release-highlight` PR labels, preserving context while it is fresh.
- Maintainers take on new responsibilities: curating release highlights, triggering ecosystem releases, and maintaining the compatibility matrix. These are manual, deliberate actions — automation is traded for trust that the matrix implies validated combinations.
- The matrix may lag behind individual component releases. The matrix footer will set this expectation explicitly.

**Source:** [Release Communication Strategy proposal](../../openspec/changes/release-communication-strategy/proposal.md)
