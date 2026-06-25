## Context

The complytime ecosystem is a multi-repository open-source project delivering automated compliance assessment. The core repositories are:

| Repository | Role | Current Version |
|---|---|---|
| complyctl | Runtime client (CLI + SDK) | v1.0.0-beta.0 |
| complytime-providers | Evaluator plugins (OpenSCAP, Ampel, OPA) | v0.x.y |
| complytime-policies | Gemara policy bundles (OCI artifacts) | content-versioned |
| complypack | Content packaging tool | v0.x.y |
| org-infra | CI/CD, Ampel Granular rules | infrastructure |
| complytime (this repo) | Design docs, architecture decisions | docs-only |

The primary stakeholders today are teams within the same organization, planning production adoption. However, as an open-source project, the complytime ecosystem serves a broader community of potential adopters and contributors. The communication strategy must work for all audiences, not just the current primary stakeholders. A near-term milestone exists for communicating the complyctl v1.0.0 release.

This design originated from an exploratory discussion between @gvauter, @jflowers, @jpower432, and @marcusburghardt, where multiple stakeholder perspectives were examined to identify the core communication and versioning challenges.

The current state has three problems:
1. Go semver pre-release suffixes are misread as product immaturity signals by stakeholders unfamiliar with Go module versioning conventions. In Go specifically, pre-release suffixes signal API stability commitments (a `v1.0.0` locks the import path and any future breaking change forces a `v2.0.0` with import path churn for all consumers), not product readiness. This distinction is less relevant in other language ecosystems where major version bumps carry fewer structural consequences.
2. There is no ecosystem-level coordination point. Each repo releases independently, and no single artifact tells stakeholders "these versions work together for this capability."
3. Breaking change risks, while low and well-mitigated, lack automated CI enforcement beyond protobuf checks.

complyctl is a core project in the ecosystem. Changes to complyctl's API surface directly impact all integrated projects — providers, content tooling, and downstream consumers. This central role makes its stability especially critical. Architectural decisions such as separating providers into their own repository (complytime-providers) were made specifically to reduce the blast radius of changes and keep the core as stable as possible. Further efforts in this direction are being considered to minimize coupling between the core and the broader ecosystem.

complyctl's public API surface is narrow and well-architected: the `Provider` interface in `pkg/provider`, ~15 domain types, and generated gRPC code in `api/plugin`. The `internal/` package boundary protects most code from external import. Existing safeguards include the optional interface pattern for extending RPCs without breaking existing providers, domain type insulation from protobuf, Buf breaking change detection, and frozen handshake values.

## Goals / Non-Goals

**Goals:**

- Decouple capability readiness communication from individual component version numbers so stakeholders can assess production readiness without interpreting Go semver semantics.
- Establish the complytime repository as the single coordination point for ecosystem releases, where maintainers sign off on validated combinations of component versions.
- Define a sustainable, repeatable process for ecosystem releases that works for every release cycle, not just the first stable release milestone.
- Automate breaking change detection across the three layers where incompatibilities can occur (protobuf wire, Go API surface, integration wire).
- Make the compatibility matrix self-explanatory within 30 seconds to stakeholders who may not be familiar with Go module versioning conventions.
- Support per-provider maturity levels within the complytime-providers repository, since providers within a single repo version may be at different readiness stages.
- Account for unpredictable release cadence driven by external factors (CVEs, dependency updates, Go version bumps) in the communication model.

**Non-Goals:**

- Downstream distribution packaging. This proposal focuses on upstream communication and coordination. Downstream consumers may build their own integration and packaging strategies on top of the ecosystem release artifacts.
- Automated integration testing infrastructure. The CI enforcement layers are specified at the detection level; full integration test harnesses are a separate concern.
- External provider compatibility tracking. The matrix covers complytime-org repositories. External providers self-declare compatibility.
- Synchronizing release cadences across repositories. Each repo continues to release independently.
- Defining the detailed implementation of CI workflows. The proposal specifies what each CI layer detects, not the pipeline YAML.
- Content versioning strategy for OCI artifacts (policies, complypacks). These follow their own lifecycle.

## Decisions

### 1. Capability maturity levels use Alpha/Beta/Pre-GA/GA terminology

**Decision**: Use Alpha, Beta, Pre-GA, and GA as the four maturity tiers for capabilities in the compatibility matrix.

**Rationale**: These terms are standard in the cloud-native ecosystem (Kubernetes API groups, CNCF project maturity). They are intuitive to both technical and non-technical audiences. Critically, they do not collide with Go semver pre-release suffixes (-alpha, -beta, -rc), which avoids re-entangling the version label confusion at the capability level. "Pre-GA" was chosen over "RC" (Release Candidate) specifically because RC is developer jargon tied to git tags and would recreate the same semiotic collision the matrix is designed to solve.

**Alternatives considered**:
- Using semver-aligned terms (Alpha/Beta/RC/Stable): Rejected because it re-entangles capability readiness with release labels.
- Using readiness terms (Experimental/Ready for Testing/Production Ready): More descriptive but non-standard and verbose for matrix columns.
- Three tiers without Pre-GA: Rejected because the validation gate between Beta and GA is a distinct, important state for stakeholder testing.

### 2. The compatibility matrix is a maintainer sign-off, not a live dashboard

**Decision**: The matrix updates only when maintainers trigger the ecosystem release process in the complytime repository. It is not automatically updated by individual repo releases.

**Rationale**: A CVE patch in complyctl that bumps the version should not automatically update the matrix, because the new version hasn't been integration-tested with the other components. The matrix's value is that it implies "we tested these together." Automatic updates would undermine that trust signal. Maintainers deliberately trigger the process after verifying integration compatibility, making the matrix a curated checkpoint rather than a version tracker.

**Alternatives considered**:
- Auto-update on any repo release: Rejected because it implies validation that hasn't happened and creates noise from non-capability-related releases (CVE patches, dependency bumps).
- Scheduled periodic updates (monthly): Rejected because release cadence is unpredictable and a fixed schedule doesn't align with engineering reality.

### 3. Release highlights are curated via PR labels at merge time

**Decision**: Maintainers apply a `release-highlight` label to PRs in any complytime org repo when the change is stakeholder-visible. The ecosystem release process queries these labels to compose release notes.

**Rationale**: Curation at merge time captures context when it's fresh. The alternative — reviewing all merged PRs at ecosystem release time — requires maintainers to recall context from weeks or months prior. A single label (`release-highlight`) is the minimal overhead approach; granular label categories (highlight/feature, highlight/security, highlight/breaking) can be added later if needed.

**Alternatives considered**:
- Manual release notes composition at ecosystem release time: Higher cognitive burden, context loss, more error-prone.
- Automated changelog generation from conventional commits: Produces developer-oriented output, not stakeholder-oriented. Lacks curation.
- Granular label categories from the start: More structure but higher adoption friction. Starting simple with one label and evolving is lower risk.

### 4. Three-layer breaking change detection in CI

**Decision**: Implement breaking change detection across three layers, each catching different failure modes.

| Layer | Tool | Catches | Status |
|---|---|---|---|
| 1. Protobuf wire | buf breaking | Proto schema changes that break wire protocol | Exists |
| 2. Go API surface | apidiff | Changes to exported types/functions/interfaces in `pkg/provider` and `api/plugin` | To implement |
| 3. Integration wire | Custom CI job | Runtime incompatibilities between compiled provider binaries and new complyctl | To implement |

**Rationale**: Each layer catches failures the others miss. Buf catches proto breaks but not Go type changes. Apidiff catches Go API breaks but not wire-level serialization issues. Integration tests catch subtle runtime incompatibilities that pass static checks. The existing `test-provider` and `behavioral-report` commands in complyctl provide a foundation for Layer 3.

**Alternatives considered**:
- Buf breaking only (Layer 1): Already in place but insufficient. Doesn't catch Go API changes that don't originate from proto definitions.
- Apidiff only (Layer 2): Catches Go API breaks but misses proto wire issues and runtime incompatibilities.
- Full integration test suite: Desirable long-term but too heavy for the first stable release timeline. Layer 3 is a lightweight version focused specifically on backward compatibility.

### 5. Ecosystem release artifacts live in the complytime repository

**Decision**: The compatibility matrix, ecosystem release notes, and the release process definition all live in this repository (complytime) and are rendered on complytime.dev via Docsify.

**Rationale**: This repository is already the design and architecture home. It uses Docsify for web rendering at complytime.dev. Placing ecosystem release artifacts here makes the complytime repo the single coordination point and avoids creating a new repository for release management. The matrix becomes a markdown file in `docs/guides/` that Docsify renders, maintaining consistency with the existing documentation structure.

**Alternatives considered**:
- Separate release-management repository: More isolation but adds coordination overhead and fragments the documentation.
- In each individual repo: Defeats the purpose of ecosystem-level coordination.
- External platform (wiki, Notion): Breaks the docs-as-code model and complicates automation.

## Risks / Trade-offs

### [Risk] Maturity level disagreements between maintainers
Different maintainers may assess the same capability's readiness differently, leading to inconsistent or contentious maturity assignments.
**Mitigation**: Maturity level changes require PR review by at least two maintainers (matching the existing ADR approval gate). The maturity definitions provide objective criteria. Disagreements are resolved through the standard PR review process.

### [Risk] Label adoption friction
Maintainers may forget to apply `release-highlight` labels, resulting in incomplete ecosystem release notes.
**Mitigation**: Start with a single label to minimize overhead. The ecosystem release process can include a review step where maintainers scan recent merged PRs for missed highlights before publishing. Over time, PR templates can include a reminder.

### [Risk] Matrix staleness
If ecosystem releases happen infrequently, the matrix may lag behind the actual state of the ecosystem, frustrating stakeholders who know newer versions exist.
**Mitigation**: The matrix footer explicitly states "Individual components may have newer releases that are not yet part of a validated combination." This sets the expectation that the matrix is a checkpoint, not real-time tracking. Maintainers commit to triggering ecosystem releases at meaningful capability milestones.

### [Risk] Pre-GA creates another label to misinterpret
Introducing a Pre-GA maturity level adds a new term that stakeholders might confuse with other terminology.
**Mitigation**: The maturity level definitions are published alongside the matrix on complytime.dev. The definitions use plain language focused on what the stakeholder should do (evaluate, adopt with awareness, validate, use in production), not on internal engineering states.

### [Risk] Breaking change CI adds merge friction
Apidiff and integration wire tests may flag changes that are intentional, blocking PRs unnecessarily.
**Mitigation**: Both CI layers include a maintainer override escape hatch. The check warns rather than hard-blocks, and maintainers can merge with documented justification. The breaking change policy documents when overrides are appropriate (e.g., no known external consumers, coordinated migration).

### [Trade-off] Manual ecosystem release trigger vs. automation
The deliberate choice of manual triggering means maintainers must remember to trigger ecosystem releases. This trades convenience for trust (the matrix always implies validation happened).
**Accepted**: The alternative (automatic triggers) undermines the sign-off semantics. The risk of infrequent updates is lower than the risk of publishing untested combinations.

### [Trade-off] Single `release-highlight` label vs. granular categories
Starting with one label means ecosystem release notes lack automatic categorization (features vs. security vs. breaking changes). Maintainers must organize highlights manually during the release process.
**Accepted**: Low adoption friction outweighs categorization convenience. Granular labels can be introduced when the process is mature and the cost of manual categorization becomes demonstrable.

## Resolved Questions

1. **Ecosystem release cadence**: Maintainers define an initial cadence but adapt it to reality and stakeholder demand over time. The process should be flexible and easy to adjust without friction — no rigid schedule that forces releases without meaningful changes or blocks releases when they are needed. The cadence is a guideline, not a constraint.

2. **Tagging convention for ecosystem releases**: Date-based tags (e.g., `ecosystem-2026.07`) are preferred as they are intuitive for the targeted consumers of this information. Like the cadence, the convention should be easy to adapt without friction if a different scheme proves more practical over time.

3. **Provider maturity declaration mechanism**: Multiple sources of maturity information should coexist. Automated signals such as the provider's `Describe` RPC response can report a self-declared maturity level. However, maturity can also change based on factors outside the code — for example, a provider that has been tested over time with no bugs found may be promoted without code changes. The matrix-side declaration in this repo serves as the authoritative source, informed by both automated signals and maintainer judgment. This keeps the mechanism flexible as the ecosystem evolves.

4. **Cross-org label querying permissions**: All complytime repositories are public. Querying labels and PRs across public repositories requires only read access, which is available without special token permissions. The ecosystem release process can use the GitHub API with default public read access.
