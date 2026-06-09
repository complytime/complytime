# ComplyPack Phase 0: Distribution Format

Implementation plan for the ComplyPack OCI distribution envelope. Extracted from CEP-0001.

**Status:** in-progress

**Related ADRs:** [0003](../ADRs/0003-oci-artifact-distribution.md), [0005](../ADRs/0005-two-stream-content-model.md), [0006](../ADRs/0006-complypack-content-envelope.md)

## Goal

Define and ship the ComplyPack as a uniform OCI distribution envelope that complyctl can pull and route to providers generically.

## Deliverables

| Step | Deliverable | Repo |
|:---|:---|:---|
| Define OCI manifest structure (config with ComplyPack ID, evaluator-id + opaque content layers) | ComplyPack OCI artifact spec | complypack |
| Extend `complyctl get` to pull ComplyPack artifacts and route to provider by evaluator-id | `complyctl get` handles ComplyPacks | complyctl |
| Build OPA ComplyPack: standard OPA bundle inside envelope | OPA ComplyPack published to OCI registry | complypack |

## OCI Artifact Structure

```
OCI Manifest
├── Config: application/vnd.complypack.config.v1+json
│   (ComplyPack ID, evaluator-id, version, provenance)
└── Layer: application/vnd.complypack.content.v1.tar+gzip
    (provider-specific assessment content)
```

**Config fields:**
- `id` — reverse-domain identifier for this ComplyPack (e.g., `io.complytime.opa.cis-k8s`)
- `evaluator-id` — which provider consumes this pack (e.g., `opa`)
- `version` — pack version (semver)
- `source.gemara_content` — provenance: Gemara content version this was generated from
- `source.policy_id` — provenance: which policy this implements

## Verification requirements

| Test | Criteria |
|:---|:---|
| OCI delivery | `complyctl get` pulls ComplyPack, routes to provider by evaluator-id |
| Signature verification | Unsigned pack rejected; signed pack accepted; tampered pack rejected |
| Content safety | Rejects anomalously large or malformed content, validates paths |

## Next phase

[Phase 1: OPA Native Provider](complypack-phase1.md) — ships the OPA evaluation provider that consumes ComplyPack content.
