# Ecosystem Release Process

The ecosystem release process is the coordination point for communicating complytime's capabilities to stakeholders. Maintainers trigger it from the complytime repository — individual repo releases do not trigger it automatically. Each ecosystem release represents a maintainer sign-off that the referenced component versions have been integration-tested together.

## When to trigger an ecosystem release

Trigger an ecosystem release when a meaningful capability milestone has been reached and the component versions have been validated together. There is no fixed cadence. Release when there is something worth communicating to stakeholders.

Examples of good triggers:
- A new scanning capability is available end-to-end
- A provider moves from Alpha to Beta maturity
- A significant compatibility expansion (new content bundle support)

Do not trigger an ecosystem release for routine maintenance (dependency bumps, internal refactoring) unless it changes what stakeholders can do.

## Release process steps

### Step 1: Confirm validation criteria

Before updating the compatibility matrix, confirm that the component versions you intend to reference pass the minimum validation criteria:

| Criterion | What to verify |
|:---|:---|
| Breaking change CI | All applicable breaking change CI layers pass for the referenced complyctl version |
| Provider compatibility | Referenced provider versions are compatible with the referenced complyctl version |
| Content bundle loading | Referenced content bundles are loadable by the referenced complyctl version |

If any criterion cannot be confirmed, stop. Do not update the matrix. Record which criteria failed and resolve them before proceeding.

### Step 2: Gather release highlights

Query all complytime org repositories for merged PRs carrying the `release-highlight` label since the last ecosystem release tag.

The "since last release" boundary is the timestamp of the most recent `ecosystem-YYYY.MM` tag in the complytime repository.

Repositories to query:

| Repository | Role |
|:---|:---|
| complyctl | Runtime client (CLI + SDK) |
| complytime-providers | Evaluator plugins (OpenSCAP, Ampel, OPA) |
| complytime-policies | Gemara policy bundles (OCI artifacts) |
| complypack | Content packaging tool |
| org-infra | CI/CD, Ampel Granular rules |

Before publishing, scan recent merged PRs for missed highlights. Maintainers may have forgotten to apply the label at merge time.

If no highlights exist since the last release, proceed with the release but note in the release notes that no stakeholder-visible changes occurred. The updated compatibility matrix is still valuable on its own.

If the GitHub API is unavailable for any repository, do not publish. Report which repositories could not be queried and retry later.

### Step 3: Update the compatibility matrix

Update the compatibility matrix with the new validated component versions and maturity level assessments.

Maturity level changes require PR review by at least two maintainers. Use these definitions:

| Level | Meaning for stakeholders |
|:---|:---|
| Alpha | Evaluate — the capability exists but may change significantly |
| Beta | Adopt with awareness — the capability is functional but not yet hardened |
| Pre-GA | Validate — the capability is ready for stakeholder testing and feedback |
| GA | Use in production — the capability is stable and supported |

The matrix footer should state: "Individual components may have newer releases that are not yet part of a validated combination."

### Step 4: Compose ecosystem release notes

Write release notes using the capability-oriented format. Present changes in terms of what stakeholders can now do, not which component versions changed. Component versions appear as supporting technical metadata.

Use the [release notes template](#release-notes-template) below.

### Step 5: Reference recent ADRs

Check the `docs/ADRs/` directory for any ADRs accepted since the last ecosystem release. Include each with a brief summary of the decision and its stakeholder impact.

If no ADRs were accepted since the last release, omit the architecture decisions section from the release notes.

### Step 6: Commit, tag, and publish

1. Commit the updated compatibility matrix and release notes to the complytime repository.
2. Create an ecosystem release tag following the [tagging convention](#tagging-convention).
3. The website auto-updates via Docsify — no separate deployment step is needed.

## Release highlight label convention

Apply the `release-highlight` label to PRs at merge time when a change is visible or relevant to stakeholders. Apply it while context is fresh.

**When to apply:**

- New scanning capability
- New provider support
- New content format or content bundle
- Security patch that stakeholders need to be aware of

**When NOT to apply:**

- Internal refactoring with no user-facing behavior change
- Dependency bumps (unless security-relevant and stakeholder-facing)
- CI/CD infrastructure changes
- Documentation-only changes (unless they reflect a capability change)

For security patches, use judgment: apply the label if stakeholders need to take action or be aware of the change.

## Tagging convention

Use date-based tags:

```
ecosystem-YYYY.MM
```

If multiple ecosystem releases occur within the same month, append a sequential suffix:

```
ecosystem-2026.07
ecosystem-2026.07.2
ecosystem-2026.07.3
```

The first release in a month has no suffix. The second release appends `.2`, the third `.3`, and so on.

Tags enable historical comparison — maintainers and stakeholders can diff between ecosystem release tags to see what changed in the matrix and release notes.

## Release notes template

```markdown
# Ecosystem Release ecosystem-YYYY.MM

## Capability Summary

<!-- What can stakeholders now do that they could not do before? -->
<!-- Write 2-3 sentences summarizing the release in capability terms. -->

## Highlights

<!-- Gathered from PRs labeled `release-highlight`, organized by capability area. -->
<!-- Attribute each highlight to its source repository. -->

### Scanning & Evaluation

- [complyctl] Description of the change
- [complytime-providers] Description of the change

### Content & Packaging

- [complypack] Description of the change
- [complytime-policies] Description of the change

### Infrastructure & CI

- [org-infra] Description of the change

## Maturity Level Changes

<!-- Include only if maturity levels changed since the last release. -->
<!-- Example: -->
| Capability | Previous | Current | Context |
|:---|:---|:---|:---|
| OPA scanning | Alpha | Beta | Passed integration testing with 3 policy bundles |

## Compatibility Matrix

<!-- Updated validated combination. Link to the full matrix if maintained separately. -->
| Component | Version | Maturity |
|:---|:---|:---|
| complyctl | vX.Y.Z | — |
| complytime-providers | vX.Y.Z | Per-provider (see below) |
| complytime-policies | content-version | — |
| complypack | vX.Y.Z | — |

Individual components may have newer releases that are not yet part of
a validated combination.

## Architecture Decisions

<!-- Omit this section if no ADRs were accepted since the last release. -->
- [ADR-NNNN: Title](../ADRs/NNNN-slug.md) — Brief summary of the decision
  and its stakeholder impact.

## Individual Repository Release Notes

- [complyctl vX.Y.Z](https://github.com/complytime/complyctl/releases/tag/vX.Y.Z)
- [complytime-providers vX.Y.Z](https://github.com/complytime/complytime-providers/releases/tag/vX.Y.Z)
- [complypack vX.Y.Z](https://github.com/complytime/complypack/releases/tag/vX.Y.Z)
```
