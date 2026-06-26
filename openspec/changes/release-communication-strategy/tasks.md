## 1. Problem Documentation (complytime repo)

- [x] 1.1 Write problem doc `docs/problems/release-communication.md` covering the version label perception gap, multi-repo coordination challenge, unpredictable release cadence, and the missing downstream translation layer. Reference existing architecture.md and ADRs for context.
- [x] 1.2 [P] Update `docs/_sidebar.md` to include the new problem doc under the Problems section.
- [x] 1.3 [P] Update `README.md` "What's here" section to reference the new problem doc.

## 2. Architecture Decision Record (complytime repo)

- [x] 2.1 Write ADR `docs/ADRs/0007-ecosystem-release-strategy.md` documenting the decision to proceed with complyctl v1.0.0, the risk assessment for API stability, the four-tier capability maturity model (Alpha/Beta/Pre-GA/GA), and the ecosystem release process triggered from the complytime repo. Follow the template in `docs/ADRs/adr-template.md`.
- [x] 2.2 [P] Update `docs/_sidebar.md` to include ADR-0007.
- [x] 2.3 [P] Update `README.md` "What's here" section to reference ADR-0007.

## 3. Compatibility Matrix (complytime repo)

- [x] 3.1 Create `docs/guides/compatibility-matrix.md` with the initial matrix content: capability-oriented rows, maturity levels, validated component versions, last validated date, maturity level definitions, introductory context about release cadence, and footnote about newer versions possibly existing outside the matrix.
- [x] 3.2 [P] Update `docs/_sidebar.md` to include the compatibility matrix under Guides with prominent placement.
- [x] 3.3 [P] Update `README.md` to reference the compatibility matrix.
- [ ] 3.4 Verify the matrix renders correctly on complytime.dev via Docsify.

## 4. Breaking Change Policy (complytime repo)

- [x] 4.1 Create `docs/guides/breaking-change-policy.md` covering: rules for v1.x API stability (interface method additions, struct field changes, function signatures, enum values, proto field numbers), the three-layer CI enforcement model with detection mapping table, existing safeguards inventory with explanations, the maintainer override escape hatch (requiring approval from at least two maintainers), and links to ADR-0004 and ADR-0006.
- [x] 4.2 [P] Update `docs/_sidebar.md` to include the breaking change policy under Guides.
- [x] 4.3 [P] Update `README.md` to reference the breaking change policy.

## 5. Ecosystem Release Process Guide (complytime repo)

- [x] 5.1 Create `docs/guides/ecosystem-release-process.md` covering: the release process steps (trigger, gather highlights, update matrix, compose release notes, publish, tag), the `release-highlight` label convention with examples, the capability-oriented release notes template, the tagging convention for ecosystem releases (date-based, e.g., `ecosystem-2026.07`), minimum validation criteria for matrix entries, and the requirement to reference recent ADRs in release notes.
- [x] 5.2 [P] Update `docs/_sidebar.md` to include the ecosystem release process guide.
- [x] 5.3 [P] Update `README.md` to reference the ecosystem release process guide.

## 6. Release Highlight Label Setup (cross-repo)

- [ ] 6.1 [P] Create the `release-highlight` label in the complytime GitHub org repositories: complyctl, complytime-providers, complytime-policies, complypack, org-infra.
- [ ] 6.2 [P] Add a brief note about the `release-highlight` label to each repository's CONTRIBUTING.md or PR template, linking to the ecosystem release process guide on complytime.dev.

## 7. Breaking Change CI - Layer 2: apidiff (complyctl repo)

- [ ] 7.1 Add `apidiff` as a CI dependency in the complyctl repository (golang.org/x/exp/cmd/apidiff or equivalent; verify current recommended tool at implementation time).
- [ ] 7.2 Create a CI job that runs apidiff comparing `pkg/provider` and `api/plugin` against the latest release tag on every PR.
- [ ] 7.3 Configure the CI job to report breaking changes as a check failure with clear output indicating which exported symbols changed and how.
- [ ] 7.4 Document the maintainer override process for intentional breaking changes (e.g., a specific PR label or review approval that allows merging despite the check failure). Override requires approval from at least two maintainers.
- [ ] 7.5 Validate the apidiff CI job by submitting a test PR with a known non-breaking change (verify pass) and a known breaking change to the Provider interface (verify fail).

## 8. Breaking Change CI - Layer 3: Integration Wire Test (complyctl repo)

- [ ] 8.1 Design the integration wire test job: build provider binaries from the latest complytime-providers release tag, run them against the current complyctl PR branch, verify Describe/Generate/Scan RPCs succeed. Define fallback behavior if no provider release tags exist or the latest tag is unbuildable.
- [ ] 8.2 Evaluate using the existing `test-provider` and `behavioral-report` commands as foundations for the integration test.
- [ ] 8.3 Implement the integration wire test CI job.
- [ ] 8.4 Validate the CI job by testing with a known-compatible provider version and verifying it passes.

## 9. Short-Term: v1.0.0 Release (complyctl repo, by agreed stable release date)

- [ ] 9.1 Cut v1.0.0-rc.1 release of complyctl.
- [ ] 9.2 Coordinate stakeholder testing period with the RC release.
- [ ] 9.3 Review API surface for any remaining concerns before v1.0.0 commitment.
- [ ] 9.4 Cut v1.0.0 stable release of complyctl.
- [ ] 9.5 Trigger the first ecosystem release from the complytime repo: populate the initial compatibility matrix, compose ecosystem release notes, tag the repo, and publish to complytime.dev.
- [ ] 9.6 Communicate the ecosystem release to stakeholders using the capability-oriented format defined in the ecosystem release process guide (Task 5.1).

## 10. Glossary Update (complytime repo)

- [x] 10.1 Update `docs/glossary.md` with a "Process Terms" section defining: ecosystem release, compatibility matrix, capability maturity level (Alpha/Beta/Pre-GA/GA), Pre-GA, and release highlight.
- [x] 10.2 [P] Update `docs/_sidebar.md` if the glossary section needs structural changes.
<!-- spec-review: passed -->
