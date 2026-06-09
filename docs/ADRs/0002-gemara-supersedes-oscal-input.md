# ADR-0002: Gemara Supersedes OSCAL as Input Model

**Status:** accepted

**Date:** 2026-02-01 (retroactive)

**Deciders:** complytime-dev team

## Context

The original compliance pipeline (complyscribe) used OSCAL as both the input authoring format and the runtime data model. 
This required compliance engineers to work directly with OSCAL's verbose JSON/XML structures and created tight coupling between authoring tooling and the assessment runtime.

Gemara provides a CUE-based schema framework purpose-built for compliance content — simpler authoring, schema validation at the type level, and OCI-native distribution. 
OSCAL remains valuable as an output/interchange format but is not the right authoring surface.

## Decision

Gemara is the input model for compliance content authoring and policy definition. complyctl consumes Gemara artifacts (`#Policy`, `#ControlCatalog`) natively.

OSCAL is supported as an output format only — `complyctl scan --format oscal` produces OSCAL Assessment Results for downstream consumers that require it.

The OSCAL-input pipeline (complyscribe) is archived.

## Consequences

- Compliance engineers author in Gemara's CUE-validated YAML — simpler than raw OSCAL.
- Schema evolution is managed via CUE module versioning, not OSCAL's monolithic release cycle.
- Organizations with existing OSCAL tooling can still consume scan results in OSCAL format.
- The compliance-trestle dependency is eliminated from the runtime path.
- complyscribe is archived; its CaC-to-OSCAL sync capabilities are no longer maintained.

**Source:** [complyctl PR #381](https://github.com/complytime/complyctl/pull/381), [complyscribe PR #854](https://github.com/complytime/complyscribe/pull/854)
