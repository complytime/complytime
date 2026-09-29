# Release Communication

The complytime ecosystem spans multiple independently-released repositories. Stakeholder capabilities — "scan my system against NIST SP 800-53" — emerge from the composition of compatible components, not from any single release. No version number answers "can I do X?" because the answer lives in the intersection of compatible versions across repositories, and that intersection is not published anywhere.

This problem is compounded by asymmetric provider maturity, unpredictable release cadence, and the absence of a downstream productization layer. The result is a persistent gap between engineering reality and stakeholder understanding — one that recurs every release cycle and scales with the ecosystem.

This doc explores *why* release communication is hard in this ecosystem and what structural factors make it resistant to simple fixes. It does not propose what ComplyTime should build.

This problem was identified through discussions between @gvauter, @jflowers, @jpower432, and @marcusburghardt.

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

## The version label perception gap

Go's module system assigns structural meaning to version numbers that other ecosystems do not:

| Version | What Go means | What stakeholders hear |
|:---|:---|:---|
| `v1.0.0-beta.0` | Pre-release. The API surface is not yet stable. Import path is not yet locked. | "This is beta software. Not production-ready." |
| `v1.0.0` | The API surface is stable. The import path is locked. Breaking changes require `v2.0.0`. | "This is the first stable release." |
| `v0.x.y` | No stability guarantees. The API may change freely. | "This is pre-release." |

In Go, the major version is encoded in the import path (`github.com/complytime/complyctl` vs `github.com/complytime/complyctl/v2`). Every consumer must update import paths when `v2.0.0` is released. For the complytime ecosystem, where [providers](../ADRs/0004-grpc-provider-plugin-architecture.md) and content tooling depend on `complyctl`'s `pkg/provider` interface and `api/plugin` gRPC definitions, a major version bump propagates across the entire [component map](../architecture.md). This makes `v1.0.0` a heavyweight commitment — it is a commitment to API stability that carries real structural consequences if broken.

A component at `v1.0.0-beta.0` can be fully functional and well-tested, but the pre-release suffix — which signals API instability to Go tooling — is read as product immaturity by stakeholders unfamiliar with Go conventions. The engineering prudence of maintaining a pre-release suffix and the communication signal it sends point in opposite directions.

With complyctl now at v1.0.0, this particular tension has been resolved. The concern remains documented here because it shaped the ecosystem's versioning decisions and may recur if new components adopt pre-release versions.

## Unpredictable release cadence

Releases in this ecosystem are driven by engineering events, not stakeholder milestones:

- A CVE in a transitive dependency triggers a patch release.
- A new Go version requires a compatibility update.
- A dependency bumps its minimum supported version.

Any of these can produce a release with zero new stakeholder-visible features. Conversely, a feature that stakeholders are waiting for may land in a repository without triggering an ecosystem-level announcement, because no such announcement mechanism exists.

This cadence is correct for an open-source project — you release when there is something to release, not on a schedule. But it means that stakeholders cannot predict when capabilities they care about will be available in a validated combination. Tying communication to individual releases creates noise: most releases are not stakeholder-relevant, and the ones that are get lost in the stream.

## Per-provider maturity asymmetry

The [provider plugin architecture](../ADRs/0004-grpc-provider-plugin-architecture.md) enables multiple evaluators to coexist within `complytime-providers`. Each provider wraps a different assessment tool — OpenSCAP for OS configuration, OPA for policy-as-code, Ampel for granular checks. These providers are at different stages of maturity:

| Provider | Maturity | Notes |
|:---|:---|:---|
| OpenSCAP | Stable | Production-tested, well-understood failure modes |
| Ampel | Stable | Actively used, clear scope |
| OPA | In development | Functional but interface still evolving |

All three ship under a single repository version. When `complytime-providers` releases a new version, that version tells you nothing about which providers are ready for production use. The repository version hides provider-level maturity behind a single number. There is no mechanism to surface "OpenSCAP is stable, OPA is in development" within the versioning scheme itself — Go semver applies to the module, not to subcomponents within it.

## The missing translation layer

In most production software ecosystems, a downstream layer absorbs these tensions. A distribution, a managed service, or a product team selects validated combinations, assigns their own version numbers, and insulates adopters from upstream versioning semantics. The complytime ecosystem does not have this layer. Stakeholders consume upstream directly — they read GitHub release pages, interpret Go version tags, and correlate repositories manually. The upstream versioning scheme, release process, and documentation must serve double duty: correct for Go tooling and comprehensible to stakeholders simultaneously.

## Open questions

- What is the simplest mechanism to ensure ecosystem compatibility given that complyctl defines the API contracts all other components implement?
- How should per-provider maturity be declared — in code (e.g., the provider's `Describe` RPC), in documentation, or both?
