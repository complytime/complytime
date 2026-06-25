# Compatibility Matrix

The complytime ecosystem spans multiple independently-released repositories. Each component follows its own release cadence driven by engineering needs — vulnerability patches, dependency updates, Go version bumps — not by a synchronized schedule.

This matrix represents **validated integration checkpoints**: specific component versions that maintainers have tested together for each capability. Individual components may release new versions between these checkpoints. The matrix shows what has been validated as a working combination, not the latest available versions.

## Maturity Levels

Each capability in the matrix carries a maturity level that describes its readiness for adoption.

| Level       | What It Means                                                                                                                          |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| **Alpha**   | Early stage. The API will change. Suitable for evaluation and feedback. Not recommended for production workloads.                       |
| **Beta**    | Functional and tested. The API may evolve with backward compatibility. Production use is possible with awareness of potential changes.  |
| **Pre-GA**  | Feature-complete and validated by maintainers. Under stakeholder testing. Changes are made only to resolve issues found during validation. |
| **GA**      | Tested across the integrated stack. Stable API. Supported for production use.                                                          |

## Validated Capabilities

> **Note:** The table below is illustrative and shows the expected structure of the matrix once the first ecosystem release is completed. Component versions and maturity levels will be updated with validated data at that time.

| Capability                         | Maturity | Validated Components                                                                                                                  | Last Validated |
| ---------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------- | -------------- |
| CIS Kubernetes scanning (Ampel)    | Pre-GA   | [complyctl] v1.0.0-rc.1 <br> ampel-provider ([complytime-providers] v0.x.y) <br> ampel-granular-rules ([org-infra]) 2026.07 <br> [complytime-policies] 2026.07 | TBD            |
| NIST 800-53 assessment (OpenSCAP)  | Pre-GA   | [complyctl] v1.0.0-rc.1 <br> openscap-provider ([complytime-providers] v0.x.y) <br> [complytime-policies] 2026.07                    | TBD            |
| OPA policy evaluation              | Alpha    | [complyctl] v1.0.0-rc.1 <br> opa-provider ([complytime-providers] v0.x.y) <br> [complytime-policies] 2026.07                         | TBD            |

> **Note:** Individual components may have newer releases that are not yet part of a validated combination. Check individual repository release notes for details.

## Component Release Notes

| Component              | Repository                                                                             |
| ---------------------- | -------------------------------------------------------------------------------------- |
| complyctl              | [complytime/complyctl](https://github.com/complytime/complyctl/releases)               |
| complytime-providers   | [complytime/complytime-providers](https://github.com/complytime/complytime-providers/releases) (Ampel, OpenSCAP, OPA providers) |
| complytime-policies    | [complytime/complytime-policies](https://github.com/complytime/complytime-policies/releases) |
| Ampel Granular Rules   | [complytime/org-infra](https://github.com/complytime/org-infra/releases)               |

---

For component roles and architecture context, see [Architecture](architecture.md).

[complyctl]: https://github.com/complytime/complyctl/releases
[complytime-providers]: https://github.com/complytime/complytime-providers/releases
[complytime-policies]: https://github.com/complytime/complytime-policies/releases
[org-infra]: https://github.com/complytime/org-infra/releases
