This spec defines the requirements for the capability-oriented compatibility matrix that maps stakeholder-meaningful capabilities to validated component combinations with explicit maturity levels, hosted on complytime.dev.

## ADDED Requirements

### Requirement: Compatibility matrix maps capabilities to validated component combinations
The compatibility matrix SHALL present ecosystem readiness as a mapping from stakeholder-meaningful capabilities to specific validated component versions. The matrix SHALL NOT be organized by repository or component; it SHALL be organized by what the stakeholder can do.

#### Scenario: Stakeholder looks up a capability
- **WHEN** a stakeholder wants to know if CIS Kubernetes scanning is production-ready
- **THEN** they SHALL find a single row in the matrix showing the capability name, its maturity level, and the exact component versions validated for that capability

#### Scenario: Multiple capabilities share components
- **WHEN** CIS Kubernetes scanning and NIST 800-53 assessment both require complyctl v1.0.0
- **THEN** complyctl v1.0.0 SHALL appear in both capability rows, and the stakeholder SHALL see that each capability has its own independent maturity level

### Requirement: Compatibility matrix uses four-tier maturity levels
The matrix SHALL use exactly four maturity levels: Alpha, Beta, Pre-GA, and GA. Each level SHALL have a published definition focused on stakeholder action, not internal engineering state.

#### Scenario: Maturity level definitions are accessible
- **WHEN** a stakeholder reads the compatibility matrix on complytime.dev
- **THEN** the maturity level definitions SHALL be visible on the same page or immediately linked, with plain-language descriptions of what each level means for adoption decisions

#### Scenario: Alpha capability in the matrix
- **WHEN** a capability is listed with maturity level Alpha
- **THEN** the definition SHALL communicate that the capability is early stage, the API will change, it is suitable for evaluation and feedback only, and it is not recommended for production workloads

#### Scenario: Beta capability in the matrix
- **WHEN** a capability is listed with maturity level Beta
- **THEN** the definition SHALL communicate that the capability is functional and tested, the API may evolve with backward compatibility, and production use is possible with awareness of potential changes

#### Scenario: Pre-GA capability in the matrix
- **WHEN** a capability is listed with maturity level Pre-GA
- **THEN** the definition SHALL communicate that the capability is feature-complete and validated by maintainers, is under stakeholder testing, and changes will only be made to resolve issues found during validation

#### Scenario: GA capability in the matrix
- **WHEN** a capability is listed with maturity level GA
- **THEN** the definition SHALL communicate that the capability is tested across the integrated stack, has a stable API, and is supported for production use

### Requirement: Compatibility matrix surfaces per-provider maturity
Individual providers within the complytime-providers repository SHALL have their maturity levels tracked independently. The repository version number SHALL NOT be used as a proxy for individual provider maturity.

#### Scenario: Providers at different maturity levels
- **WHEN** the OpenSCAP provider is GA, the Ampel provider is Beta, and the OPA provider is Alpha
- **THEN** each provider-backed capability SHALL show its own maturity level in the matrix, regardless of the complytime-providers repository version being the same for all three

#### Scenario: Provider matures without repo version change
- **WHEN** the Ampel provider reaches GA maturity at the same complytime-providers repo version
- **THEN** the matrix SHALL reflect the new maturity level for the Ampel-backed capability without requiring a new repo release

### Requirement: Compatibility matrix is a point-in-time validated snapshot
The matrix SHALL represent a specific validated combination of component versions at a point in time. It SHALL NOT represent the latest available versions of each component. The matrix SHALL explicitly communicate this distinction to readers.

#### Scenario: Newer component versions exist
- **WHEN** complyctl v1.0.1 has been released as a CVE patch but the last ecosystem release validated complyctl v1.0.0
- **THEN** the matrix SHALL still show complyctl v1.0.0, and a footnote SHALL indicate that individual components may have newer releases not yet part of a validated combination

#### Scenario: Validation date is visible
- **WHEN** a stakeholder reads the compatibility matrix
- **THEN** the date of the last ecosystem release (validation checkpoint) SHALL be clearly displayed

### Requirement: Compatibility matrix is hosted on complytime.dev
The compatibility matrix SHALL be maintained as a markdown file in the complytime repository and rendered on complytime.dev via Docsify. The matrix SHALL be accessible from the site navigation.

#### Scenario: Matrix is accessible via website
- **WHEN** a stakeholder navigates to complytime.dev
- **THEN** the compatibility matrix SHALL be reachable from the site navigation (sidebar) within one click

#### Scenario: Matrix source is in the complytime repo
- **WHEN** a maintainer needs to update the matrix
- **THEN** the update SHALL be made via a PR to the complytime repository, following the standard review process

### Requirement: Compatibility matrix has minimal required fields
The matrix SHALL contain exactly four columns to minimize cognitive load: Capability, Maturity, Validated Components, and Last Validated date. Additional detail SHALL be available via links to individual repo release notes, not embedded in the matrix itself.

#### Scenario: Matrix is compact and self-contained
- **GIVEN** a stakeholder unfamiliar with Go module versioning conventions
- **WHEN** they open the compatibility matrix on complytime.dev
- **THEN** the matrix SHALL fit in a single rendered page section without horizontal scrolling, each capability SHALL be assessable without consulting a separate glossary, and maturity level definitions SHALL be on-page rather than linked to a separate document

#### Scenario: Technical details are linked, not inlined
- **WHEN** a stakeholder wants more detail about a specific component version change
- **THEN** the matrix SHALL link to the individual repo's release notes for that version, rather than including detailed changelogs inline

### Requirement: Compatibility matrix accounts for external factors in release cadence
The matrix documentation SHALL acknowledge that component release cadence is driven by engineering needs including vulnerability patches, dependency updates, and Go version bumps, not by stakeholder feature timelines. This context SHALL be visible to matrix readers.

#### Scenario: Context about release cadence is provided
- **WHEN** a stakeholder reads the compatibility matrix page
- **THEN** an introductory section SHALL explain that component versions follow independent release cadences driven by engineering best practices and external factors, and that the matrix represents validated integration checkpoints regardless of individual component release frequency
