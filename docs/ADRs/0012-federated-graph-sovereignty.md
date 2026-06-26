# ADR-0012: Federated Graph Sovereignty

**Status:** proposed

**Date:** 2026-06-26

**Deciders:** ComplyTime maintainers

## Context

The [federated compliance graphs](../problems/federated-compliance-graphs.md) problem doc establishes that compliance data is inherently distributed and that a single centralized graph cannot serve all deployment scenarios (air-gapped, regulated, open community). Federation must be structural: enforced by schema, not by runtime policy.

Two governance decisions have crystallized. Operational details (sync mechanisms, query performance) are deferred.

## Decision

1. **Conflicting authority mappings are both valid; federation reinforces jurisdiction.** Divergent mappings from different authorities are correct within their respective jurisdictions. Federation graph boundaries delineate areas of responsibility. This aligns with [ADR-0011](0011-cross-framework-mapping.md) Decision 3.

2. **All three trust mechanisms are available; governance is configurable and auditable.** Pure voting, authority-weighted voting, and designated arbiters are valid governance mechanisms for the community mapping commons. The mechanism is configurable per commons instance. Every acceptance decision records who decided, under what authority, and with what mechanism.

## Consequences

- The commons governance model must support pluggable trust mechanisms with full audit trails.
- Federated nodes expose minimum metadata required for participation (least privilege). Specific metadata exposure rules are deferred until the schema solidifies.
- Air-gapped sync mechanisms (snapshot export/import, selective replication, integrity verification) are deferred pending implementation experience.
- Federation query performance (latency, caching, query planning across sovereign boundaries) is deferred to implementation.
- The architecture must serve three deployment modes (air-gapped, regulated, open community) without special-casing any of them.
