## Why

The complytime ecosystem spans multiple independently-released repositories (complyctl, complytime-providers, complytime-policies, complypack, org-infra). Stakeholders consuming these upstream projects interpret Go semver pre-release suffixes (-alpha, -beta, -rc) as product maturity signals, when they are actually API stability contracts specific to Go's module system. This creates persistent misalignment between sound engineering practices and business expectations, eroding stakeholder confidence despite features being functional and delivered.

The problem is compounded by three factors: (1) there is no ecosystem-level translation layer between upstream release practices and stakeholder expectations, so adopters must interpret individual repository versions directly; (2) release cadence is driven by engineering needs (vulnerability patches, dependency updates, Go version bumps) not stakeholder feature requests, so releases may happen with no stakeholder-visible changes; and (3) individual providers within complytime-providers have different maturity levels, invisible behind a single repository version number. Without a deliberate strategy, this misalignment will recur every release cycle and scale with the ecosystem.

The immediate trigger is an agreed deadline to have the first stable release of complyctl out by July 8th. The API surface is stable, no breaking changes are expected, and the pre-release suffix is precautionary. The team needs both a pragmatic short-term approach for this milestone and a sustainable long-term framework.

## What Changes

- **Ecosystem release process**: Introduce a maintainer-triggered release process in the complytime repository that serves as the single coordination point for ecosystem-wide communication. This process gathers release highlights from all project repos, updates the compatibility matrix, and publishes to complytime.dev.
- **Compatibility matrix**: Create a capability-oriented compatibility matrix that maps stakeholder-meaningful capabilities to validated component combinations with explicit maturity levels (Alpha, Beta, Pre-GA, GA). Hosted on complytime.dev, rendered from this repository, updated only through the ecosystem release process as a maintainer sign-off.
- **Capability maturity levels**: Define a four-tier maturity model (Alpha, Beta, Pre-GA, GA) that describes capability readiness across the integrated stack, deliberately decoupled from individual component version numbers. Per-provider maturity levels within complytime-providers are surfaced independently.
- **Breaking change policy and CI enforcement**: Document the rules for maintaining Go API compatibility within v1.x and implement automated detection. Three layers: protobuf wire compatibility (buf breaking, already exists), Go API surface compatibility (apidiff, to implement), and integration wire testing (old provider binary vs. new complyctl, to implement).
- **Release highlight labeling**: Establish a cross-repo labeling convention (e.g., `release-highlight`) that maintainers apply at PR merge time to flag stakeholder-visible changes. The ecosystem release process queries these labels to compose release notes automatically.
- **Stakeholder communication template**: Define a reusable capability-oriented communication format that presents ecosystem value in stakeholder language, with component versions as supporting technical metadata rather than headlines.

## Capabilities

### New Capabilities

- `ecosystem-release-process`: The maintainer-triggered workflow in the complytime repository that coordinates matrix updates, gathers highlights, composes release notes, and publishes to complytime.dev. Includes the labeling convention for release highlights across repos.
- `compatibility-matrix`: The capability-oriented matrix mapping stakeholder capabilities to validated component combinations with maturity levels. Includes the maturity level definitions (Alpha, Beta, Pre-GA, GA), matrix schema, hosting on complytime.dev, and the principle that the matrix is a maintainer sign-off on integration testing.
- `breaking-change-governance`: The policy document defining rules for Go API compatibility within v1.x, the three-layer CI enforcement model (buf breaking, apidiff, integration wire test), and the escape hatch for maintainer overrides.

### Modified Capabilities

(none - no existing specs to modify)

### Removed Capabilities

(none)

## Constitution Alignment

No org constitution file (`.specify/memory/constitution.md`) exists for this repository. Constitution assessments are not applicable. The change aligns with the project's documented principles in `docs/vision.md` (Fidelity, Decoupling, Composability).

## Impact

- **This repository (complytime)**: New documentation artifacts (problem doc, ADR, guides). Compatibility matrix file added and rendered via Docsify on complytime.dev. Possible GitHub Actions workflow for the ecosystem release process.
- **complyctl**: CI pipeline additions for apidiff and integration wire testing. No code changes to the application itself. Release process for v1.0.0-rc.1 and v1.0.0 by the agreed stable release date.
- **complytime-providers**: Per-provider maturity level declarations. Release highlight labeling convention adopted by maintainers.
- **complytime-policies, complypack, org-infra**: Release highlight labeling convention adopted. No other direct changes.
- **complytime.dev**: New compatibility matrix page rendered from this repo. Ecosystem release notes page.
- **Stakeholder communication**: Shift from per-repo version announcements to capability-oriented ecosystem release notes. Explicit maturity definitions replace implicit version-number interpretation.
- **Cross-repo coordination**: GitHub label convention (`release-highlight`) established across all complytime org repos. Ecosystem release process in this repo queries labels from other repos.
