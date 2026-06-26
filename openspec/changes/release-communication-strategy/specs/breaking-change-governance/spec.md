This spec defines the requirements for the breaking change policy, the three-layer CI enforcement model protecting complyctl's Go API surface, and the maintainer override escape hatch for intentional changes.

## ADDED Requirements

### Requirement: Breaking change policy documents rules for v1.x API stability
A breaking change policy document SHALL be published as a guide in the complytime repository. It SHALL define what constitutes a breaking change in the Go API surface of complyctl and SHALL provide actionable rules for contributors to maintain backward compatibility within the v1.x major version line.

#### Scenario: Policy covers interface changes
- **WHEN** a contributor reviews the breaking change policy
- **THEN** the policy SHALL state that methods MUST NOT be added to the `Provider` interface and that new RPCs MUST be introduced via separate optional interfaces

#### Scenario: Policy covers struct field changes
- **WHEN** a contributor reviews the breaking change policy
- **THEN** the policy SHALL state that struct fields in public types (`pkg/provider`, `api/plugin`) MUST NOT be removed or have their types changed, and that new fields MUST only be added (additive changes)

#### Scenario: Policy covers function signature changes
- **WHEN** a contributor reviews the breaking change policy
- **THEN** the policy SHALL state that existing exported function signatures MUST NOT have parameters added or removed, and that extensibility MUST be achieved via functional options patterns (variadic parameters)

#### Scenario: Policy covers enum value changes
- **WHEN** a contributor reviews the breaking change policy
- **THEN** the policy SHALL state that enum values MUST NOT be removed, renamed, or reordered, and that new values MUST only be appended

#### Scenario: Policy covers protobuf field numbers
- **WHEN** a contributor reviews the breaking change policy
- **THEN** the policy SHALL state that protobuf field numbers are immutable once assigned, and that field types MUST NOT be changed in existing messages

### Requirement: Go API surface compatibility is checked in CI
complyctl CI SHALL include an automated check that compares the current branch's exported Go API surface against the latest release tag. The check SHALL use `apidiff` (golang.org/x/exp/cmd/apidiff) or an equivalent tool. The check SHALL cover all packages importable by external consumers: `pkg/provider` and `api/plugin`.

#### Scenario: Non-breaking change passes CI
- **WHEN** a PR adds a new exported type `ValidateRequest` to `pkg/provider` without modifying existing types
- **THEN** the apidiff CI check SHALL pass, reporting the change as compatible

#### Scenario: Breaking change detected in CI
- **WHEN** a PR adds a `Validate` method to the `Provider` interface in `pkg/provider`
- **THEN** the apidiff CI check SHALL fail, reporting that the change breaks the `Provider` interface for existing implementors

#### Scenario: Breaking change with maintainer override
- **GIVEN** a PR contains a breaking change that is intentional and justified (e.g., no known external consumers, coordinated migration planned)
- **WHEN** a maintainer invokes the override for a breaking change CI failure
- **THEN** the PR SHALL require approval from at least two maintainers, and the override justification SHALL be recorded in the PR description or a dedicated comment

### Requirement: Protobuf wire compatibility is checked in CI
complyctl CI SHALL include an automated check that detects breaking changes in protobuf definitions. This check SHALL use Buf with `FILE`-level breaking change detection, comparing the current branch against the latest release.

#### Scenario: Proto field number change detected
- **WHEN** a PR changes the field number of an existing field in `plugin.proto`
- **THEN** the buf breaking CI check SHALL fail

#### Scenario: New proto field added
- **WHEN** a PR adds a new field to an existing message in `plugin.proto` with a new field number
- **THEN** the buf breaking CI check SHALL pass

#### Scenario: New RPC added to service
- **WHEN** a PR adds a new RPC to the Plugin service in `plugin.proto`
- **THEN** the buf breaking CI check SHALL pass (additive change)

### Requirement: Integration wire compatibility is tested in CI
complyctl CI SHALL include a job that tests backward compatibility at the wire level by building provider binaries from the latest released version of complytime-providers and running them against the current branch of complyctl. This test SHALL verify that the gRPC plugin protocol remains functional across versions.

#### Scenario: Old provider works with new complyctl
- **WHEN** the integration wire test runs with provider binaries built from the latest complytime-providers release tag against the current complyctl PR branch
- **THEN** the Describe, Generate, and Scan RPCs SHALL all complete without transport errors

#### Scenario: Wire incompatibility detected
- **WHEN** a change in complyctl breaks the gRPC serialization or handshake with older provider binaries
- **THEN** the integration wire test SHALL fail, indicating that the change would break existing compiled providers

#### Scenario: Test uses existing test infrastructure
- **WHEN** the integration wire test is implemented
- **THEN** it SHALL leverage the existing `test-provider` and `behavioral-report` commands in complyctl's `cmd/` directory as foundations, minimizing new test infrastructure

#### Scenario: Provider release artifacts are unavailable
- **GIVEN** the integration wire test depends on building provider binaries from the latest complytime-providers release tag
- **WHEN** no release tags exist in complytime-providers or the latest tag fails to build
- **THEN** the integration wire test SHALL skip with an informational warning rather than failing the CI check, and SHALL report which provider artifact could not be resolved

### Requirement: Breaking change categories are mapped to detection layers
The breaking change policy SHALL include a mapping table that connects each category of breaking change to the CI layer responsible for detecting it. This mapping SHALL serve as a reference for both contributors (understanding what is checked) and maintainers (verifying CI coverage).

#### Scenario: Contributor looks up what CI catches
- **WHEN** a contributor wants to know if their change to a struct field type in `pkg/provider` will be caught by CI
- **THEN** the mapping table SHALL show that struct field type changes are detected by the apidiff CI layer (Layer 2)

#### Scenario: Mapping covers all breaking change categories
- **WHEN** the mapping table is reviewed
- **THEN** it SHALL cover at minimum: interface method additions, struct field removal/type changes, function signature changes, exported type/constant removal, proto field number/type changes, proto RPC removal, enum value removal/reordering, and wire-level serialization incompatibilities

### Requirement: Breaking change policy references existing safeguards
The policy document SHALL acknowledge and reference the safeguards already in place in complyctl's codebase, including the optional interface pattern for extending RPCs, domain type insulation from protobuf, internal package boundaries, Buf breaking detection configuration, and frozen handshake values.

#### Scenario: Existing safeguards are documented
- **WHEN** a new contributor reads the breaking change policy
- **THEN** they SHALL find a section listing the existing architectural safeguards with explanations of how each one protects against specific categories of breaking changes

#### Scenario: Policy connects to architecture decisions
- **WHEN** the breaking change policy references the provider plugin architecture
- **THEN** it SHALL link to ADR-0004 (gRPC Provider Plugin Architecture) and ADR-0006 (ComplyPack as Content Envelope) for architectural context
