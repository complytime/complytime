# ADR-0003: OCI Registries for Artifact Distribution

**Status:** accepted

**Date:** 2026-02-01 (retroactive)

**Deciders:** complytime-dev team

## Context

Compliance assessment requires distributing multiple artifact types: policy content (Gemara), assessment logic (ComplyPacks), and potentially provider binaries. 
Each could use its own distribution mechanism (RPM, git clone, manual placement, custom APIs).

OCI registries are already the standard distribution mechanism in cloud-native infrastructure. Organizations running Kubernetes already operate registry infrastructure. 
OCI provides content-addressable storage, signature attachment (via referrers API), and established mirroring tools for airgap environments.

## Decision

OCI registries are the primary distribution mechanism for all ComplyTime artifacts. `complyctl get` pulls artifacts via OCI, caches them locally as OCI layouts, and routes to consumers.

Airgap transport uses standard OCI mirroring tools (skopeo, crane, oras) or filesystem-backed OCI layouts. 
No special architecture for disconnected environments.

Docker credential helpers are reused for registry authentication — no custom auth mechanism.

## Consequences

- One distribution path for all artifact types (Gemara content, ComplyPacks).
- Airgap support is inherited from OCI ecosystem tooling — no special handling.
- Organizations reuse existing registry infrastructure and access controls.
- Content-addressable identity (SHA-256 digests) provides tamper detection without additional trust mechanisms.
- Signature attachment via Sigstore/Cosign follows established OCI patterns.
- Environments without OCI registry access require mirroring setup — adds operational overhead for highly restricted networks.

**Source:** [complyctl spec 001](https://github.com/complytime/complyctl/blob/main/specs/001-gemara-native-workflow/spec.md), CEP-0001
