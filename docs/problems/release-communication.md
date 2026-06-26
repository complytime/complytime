# Release Communication

Stakeholders consuming the complytime ecosystem interpret version numbers as readiness signals. They are not — at least not in the way stakeholders expect. Go semver pre-release suffixes (`-alpha`, `-beta`, `-rc`) are API stability contracts specific to Go's module system. A fully functional capability labeled `v1.0.0-beta` is perceived as "not ready for production" by anyone unfamiliar with Go conventions. The version label says one thing to the Go toolchain and something entirely different to the person deciding whether to adopt.

This problem is compounded by the ecosystem's multi-repo structure: value emerges from the composition of compatible components across repositories, but no single version number answers "can I do X?" The result is a persistent gap between sound engineering practices and stakeholder confidence — one that recurs every release cycle and scales with the ecosystem.

This doc explores *why* release communication is hard in this ecosystem and what structural factors make it resistant to simple fixes. It does not propose what ComplyTime should build.

This problem was identified through discussions between @gvauter, @jflowers, @jpower432, and @marcusburghardt.

## The version label perception gap

Go's module system assigns structural meaning to version numbers that other ecosystems do not. The distinction matters:

| Version | What Go means | What stakeholders hear |
|:---|:---|:---|
| `v1.0.0-beta.0` | The API surface may still change. Import path is not yet locked. | "This is beta software. Not production-ready." |
| `v1.0.0` | The API surface is stable. The import path is locked. Breaking changes require `v2.0.0`. | "This is the first stable release." |
| `v0.x.y` | No stability guarantees. The API may change freely. | "This is pre-release." |

The first and second columns are not contradictory — they describe orthogonal concerns. A component at `v1.0.0-beta.0` can be fully functional, well-tested, and used in production workflows. The pre-release suffix is a statement about API contract stability, not about whether the software works. But the distinction is invisible to anyone outside the Go ecosystem, and most stakeholders are outside the Go ecosystem.

The perception gap is not a misunderstanding that documentation can fix. It is a semiotic collision: the same label carries different meanings in different contexts, and there is no mechanism to signal which meaning applies.

## The Go v1.0.0 commitment

In most language ecosystems, bumping from `v1.x` to `v2.0.0` is a routine major release. In Go, it is a structural event. Go modules encode the major version in the import path:

```
github.com/complytime/complyctl          // v0.x or v1.x
github.com/complytime/complyctl/v2       // v2.x — different import path
```

Every consumer that imports `complyctl` must update their import paths when `v2.0.0` is released. For the complytime ecosystem, where [providers](../ADRs/0004-grpc-provider-plugin-architecture.md) and content tooling depend on `complyctl`'s `pkg/provider` interface and `api/plugin` gRPC definitions, a major version bump propagates across the entire [component map](../architecture.md). This makes `v1.0.0` a heavyweight decision — it is a commitment to API stability that carries real structural consequences if broken.

The practical effect: maintainers are rationally cautious about dropping the pre-release suffix, even when the software is functionally complete. That caution reads as hesitation to stakeholders. The engineering prudence and the communication signal point in opposite directions.

## Multi-repo coordination

The complytime ecosystem spans multiple independently-released repositories:

| Repository | Role | Versioning |
|:---|:---|:---|
| complyctl | Runtime client (CLI + SDK) | Go semver |
| complytime-providers | Evaluator plugins (OpenSCAP, Ampel, OPA) | Go semver |
| complytime-policies | Gemara policy bundles (OCI artifacts) | Content-versioned |
| complypack | Content packaging tool | Go semver |
| org-infra | CI/CD, Ampel Granular rules | Infrastructure |

See [architecture.md](../architecture.md) for the component vocabulary and [ADR-0005](../ADRs/0005-two-stream-content-model.md) for the two-stream content model that governs how assessment logic and compliance content are separated.

A stakeholder capability — "scan my RHEL system against NIST SP 800-53" — requires a specific combination: a `complyctl` version, a `complytime-providers` version with a functional OpenSCAP provider, and a `complytime-policies` bundle containing the relevant profile. No single repository version answers the stakeholder's question. The answer lives in the intersection of compatible versions across repositories, and that intersection is not published anywhere.

Each repository releases on its own schedule. A new `complyctl` release does not imply that providers have been tested against it. A provider release does not imply that new policy bundles exist. The absence of an ecosystem-level coordination point forces stakeholders to reverse-engineer compatibility from repository tags, changelogs, and go.mod files.

## Unpredictable release cadence

Releases in this ecosystem are driven by engineering events, not stakeholder milestones:

- A CVE in a transitive dependency triggers a patch release.
- A new Go version requires a compatibility update.
- A dependency bumps its minimum supported version.

Any of these can produce a release tomorrow with zero new stakeholder-visible features. Conversely, a feature that stakeholders are waiting for may land in a repository without triggering an ecosystem-level announcement, because no such announcement mechanism exists.

This cadence is correct for an open-source project — you release when there is something to release, not on a schedule. But it means that stakeholders cannot predict when capabilities they care about will be available in a validated combination. Tying communication to individual releases creates noise: most releases are not stakeholder-relevant, and the ones that are get lost in the stream.

## Per-provider maturity asymmetry

The [provider plugin architecture](../ADRs/0004-grpc-provider-plugin-architecture.md) enables multiple evaluators to coexist within `complytime-providers`. Each provider wraps a different assessment tool — OpenSCAP for OS configuration, OPA for policy-as-code, Ampel for granular checks. These providers are at different stages of maturity:

| Provider | Maturity | Notes |
|:---|:---|:---|
| OpenSCAP | Stable | Production-tested, well-understood failure modes |
| Ampel | Stable | Actively used, clear scope |
| OPA | In development | Functional but interface still evolving |

All three ship under a single repository version. When `complytime-providers` releases `v0.4.0`, the version tells you nothing about which providers are ready for production use. A stakeholder evaluating the OPA provider sees the same version as a stakeholder relying on OpenSCAP, but their confidence should be different.

The repository version hides provider-level maturity behind a single number. There is no mechanism to surface "OpenSCAP is stable, OPA is alpha" within the versioning scheme itself — Go semver applies to the module, not to subcomponents within it.

## The missing translation layer

In most production software ecosystems, a downstream layer absorbs these tensions. A distribution, a managed service, or a product team selects validated combinations, assigns their own version numbers, writes stakeholder-facing release notes, and insulates adopters from upstream versioning semantics. The downstream layer translates engineering artifacts into adoption artifacts.

The complytime ecosystem does not have this layer. Stakeholders consume upstream directly. They read GitHub release pages. They interpret Go version tags. They correlate repositories manually. Every tension described above — the perception gap, the multi-repo coordination problem, the cadence mismatch, the provider asymmetry — lands directly on the stakeholder without mediation.

This is not unusual for open-source projects at this stage. But it means the upstream project must solve communication problems that would otherwise be absorbed by downstream productization. The upstream versioning scheme, release process, and documentation must serve double duty: they must be correct for Go tooling and comprehensible to non-Go stakeholders simultaneously.

## Open questions

- What is the right unit of communication — individual component releases, validated combinations, or capability milestones?
- Can a compatibility matrix remain trustworthy if release cadence is unpredictable and integration testing is not fully automated?
- How should per-provider maturity be declared — in code (e.g., the provider's `Describe` RPC), in documentation, or both?
- Is a four-tier maturity model (Alpha, Beta, Pre-GA, GA) the right granularity, or does it risk recreating the same label confusion it aims to solve?
- What mechanism triggers ecosystem-level communication when the significant event is the *combination* of releases, not any single release?
