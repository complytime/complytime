This spec defines the requirements for the maintainer-triggered ecosystem release process that coordinates communication across all complytime repositories and publishes validated combinations to complytime.dev.

## ADDED Requirements

### Requirement: Ecosystem release is triggered from the complytime repository
The ecosystem release process SHALL be initiated by maintainers through a deliberate action in the complytime repository (e.g., a GitHub Actions workflow dispatch or a documented manual process). Individual repository releases SHALL NOT automatically trigger an ecosystem release. The ecosystem release represents a maintainer sign-off that the referenced component versions have been integration-tested together.

#### Scenario: Maintainer triggers ecosystem release
- **GIVEN** the complytime ecosystem has multiple independently-released repositories
- **WHEN** a maintainer triggers the ecosystem release process in the complytime repository
- **THEN** the process SHALL gather release highlights, update the compatibility matrix, compose ecosystem release notes, and publish updates to complytime.dev

#### Scenario: Individual repo release does not trigger ecosystem release
- **GIVEN** a new release has been tagged in complyctl, complytime-providers, or any other ecosystem repository
- **WHEN** the release event is processed
- **THEN** the compatibility matrix and ecosystem release notes SHALL NOT be automatically updated

### Requirement: Ecosystem release requires minimum validation criteria
Before updating the compatibility matrix with a new validated combination, the ecosystem release process SHALL require confirmation that the referenced component versions have been tested together. At minimum, the validation SHALL include: (1) all applicable breaking change CI layers pass for the referenced complyctl version, (2) the referenced provider versions are compatible with the referenced complyctl version, and (3) the referenced content bundles are loadable by the referenced complyctl version.

#### Scenario: Validation criteria met before matrix update
- **GIVEN** a maintainer has triggered the ecosystem release process
- **WHEN** the process updates the compatibility matrix with new component versions
- **THEN** the maintainer SHALL have confirmed that the minimum validation criteria are satisfied for the referenced combination

#### Scenario: Validation criteria not met
- **GIVEN** a maintainer has triggered the ecosystem release process
- **WHEN** one or more validation criteria cannot be confirmed for the intended component versions
- **THEN** the process SHALL NOT update the compatibility matrix and SHALL report which criteria were not satisfied

### Requirement: Release highlights are gathered from labeled PRs across repos
The ecosystem release process SHALL query all complytime org repositories for merged PRs and issues carrying the `release-highlight` label since the last ecosystem release tag. The "since last release" boundary SHALL be determined by the timestamp of the most recent ecosystem release tag in the complytime repository. The gathered highlights SHALL be used to compose the stakeholder-facing ecosystem release notes.

#### Scenario: Highlights exist across multiple repos
- **GIVEN** the last ecosystem release was tagged as `ecosystem-2026.07` and complyctl has 3 PRs labeled `release-highlight` and complytime-providers has 1 PR labeled `release-highlight` merged since that tag
- **WHEN** the ecosystem release process runs
- **THEN** all 4 highlights SHALL be included in the ecosystem release notes, attributed to their source repository

#### Scenario: No highlights exist since last release
- **GIVEN** the last ecosystem release was tagged and no PRs with `release-highlight` label have been merged since that tag
- **WHEN** the ecosystem release process runs
- **THEN** the ecosystem release notes SHALL indicate that no stakeholder-visible changes occurred and SHALL still include the updated compatibility matrix

#### Scenario: GitHub API is unavailable during release process
- **GIVEN** the ecosystem release process requires querying GitHub for labeled PRs across repositories
- **WHEN** the process cannot query one or more repositories for release highlights
- **THEN** the process SHALL report which repositories could not be queried and SHALL NOT publish ecosystem release notes until all repositories have been successfully queried

### Requirement: Release highlight label convention
Maintainers SHALL apply the `release-highlight` label to PRs in any complytime org repository when the change is visible or relevant to stakeholders. The label SHALL be applied at merge time when context is fresh. The label convention SHALL be documented in each repository's contributing guide.

#### Scenario: Stakeholder-visible feature merged
- **GIVEN** a PR adds a new compliance scanning capability
- **WHEN** a maintainer merges the PR
- **THEN** the maintainer SHALL apply the `release-highlight` label to the PR

#### Scenario: Internal refactoring merged
- **GIVEN** a PR refactors internal code without changing user-facing behavior
- **WHEN** a maintainer merges the PR
- **THEN** the maintainer SHALL NOT apply the `release-highlight` label

#### Scenario: Security patch with stakeholder impact
- **GIVEN** a PR patches a CVE affecting a dependency used in compliance scanning
- **WHEN** a maintainer merges the PR
- **THEN** the maintainer SHALL use judgment to decide whether the `release-highlight` label applies, based on whether stakeholders need to be aware of the change

### Requirement: Ecosystem release notes follow capability-oriented format
The ecosystem release notes SHALL present changes in terms of capabilities delivered, not component versions. Component versions SHALL appear as supporting technical metadata. The format SHALL be reusable across release cycles without reinvention. The release notes SHALL include placeholder sections for: capability summary, highlights organized by capability area, maturity level changes, the updated compatibility matrix, and links to individual repo release notes.

#### Scenario: Release notes structure
- **GIVEN** the ecosystem release process has gathered highlights and updated the compatibility matrix
- **WHEN** ecosystem release notes are composed
- **THEN** the notes SHALL include: a capability-oriented summary of changes, highlights gathered from labeled PRs organized by capability area, any maturity level changes for capabilities, the updated compatibility matrix, and references to individual repo release notes for detailed changelogs

#### Scenario: Maturity level change is communicated
- **GIVEN** a capability's maturity level has changed since the last ecosystem release (e.g., OPA scanning moved from Alpha to Beta)
- **WHEN** ecosystem release notes are composed
- **THEN** the notes SHALL explicitly highlight the maturity transition with context on what changed

### Requirement: Ecosystem release creates a checkpoint tag in the complytime repository
Each ecosystem release SHALL be recorded as a tag in the complytime repository to mark the point-in-time state of ecosystem documentation, matrix, and release notes. The tag SHALL follow the date-based format `ecosystem-YYYY.MM` where `YYYY.MM` is the release year and month. If multiple ecosystem releases occur within the same month, a sequential suffix SHALL be appended (e.g., `ecosystem-2026.07.2`).

#### Scenario: Ecosystem release is tagged
- **GIVEN** the ecosystem release process has completed successfully
- **WHEN** all artifacts have been committed and published
- **THEN** a tag following the `ecosystem-YYYY.MM[.N]` format SHALL be created in the complytime repository

#### Scenario: Tag enables historical comparison
- **GIVEN** multiple ecosystem release tags exist in the complytime repository
- **WHEN** a stakeholder or maintainer wants to compare the current ecosystem state with a prior release
- **THEN** they SHALL be able to diff between ecosystem release tags to see what changed in the matrix and release notes

### Requirement: Ecosystem release includes recent architecture decisions
The ecosystem release notes SHALL reference any ADRs accepted in the complytime repository since the last ecosystem release, providing stakeholders with visibility into architectural direction changes.

#### Scenario: New ADR accepted since last release
- **GIVEN** ADR-0007 was accepted between ecosystem release N and ecosystem release N+1
- **WHEN** ecosystem release N+1 notes are composed
- **THEN** the notes SHALL reference ADR-0007 with a brief summary of the decision and its stakeholder impact

#### Scenario: No new ADRs since last release
- **GIVEN** no ADRs were accepted between ecosystem releases
- **WHEN** ecosystem release notes are composed
- **THEN** the architecture decisions section SHALL be omitted from the release notes
