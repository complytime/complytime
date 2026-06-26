# ADR-0015: System Modeling Architecture

**Status:** proposed

**Date:** 2026-06-26

**Deciders:** ComplyTime maintainers

## Context

Compliance assessment assumes you know what you are assessing. The [system-modeling-and-drift](../problems/system-modeling-and-drift.md) problem doc explores why this is hard: infrastructure is provisioned dynamically, configurations drift from declared state, and the gap between deployed and running widens between assessment cycles.

The problem doc surveys infrastructure-as-code scanning, CSPM, Cartography, FINOS CALM, SysML v2, and custom CMDBs. It proposes a two-layer approach separating runtime discovery (automated, continuous) from analytical modeling (formal, on-demand). Six architectural decisions have crystallized from that exploration.

## Decision

1. **Formal models are projections from a property graph.** The internal representation is an organic property graph shaped by actual queries and usage. Formal models (SysML v2, CALM) are projections: export formats for automated reasoning, formal validation, and cross-tool interoperability. Translation adapters are regression-tested contracts. Whether the formal model should eventually become the authoritative view with the graph as derived cache remains an open tension, deferred until usage patterns clarify. For now, the graph is primary.

2. **MVP discovery: SCAP results from RHEL, OpenShift/Kubernetes, and cloud provider APIs.** These are the lowest-friction starting points with the broadest compliance coverage. Network traffic observation is explicitly deferred, but the modular adapter architecture must not preclude it. The discovery layer remains extensible: new adapters for new infrastructure types without rewriting the core.

3. **Property graph handles all compliance queries; formal models for specialized analysis.** The property graph is the operational representation for all compliance work. Formal system models are projections generated when specific use cases demand them: automated reasoning about system properties, formal validation, or cross-tool export.

4. **Drift is measured against requirements, not architecture patterns.** Compliance drift means deviating from requirements. If a reference architecture is itself a mandated requirement, deviation is drift. If it is one implementation among alternatives satisfying the same requirements, an alternative that still meets requirements is compliant, not drifted. The system evaluates against requirements, not prescribed topologies.

5. **Ephemeral infrastructure is fully captured with frequency identification and exponential age-off.** No special exemption for ephemeral resources. Full capture with temporal frequency analysis distinguishes legitimate transient workloads (scheduled CI runners) from anomalous patterns (unexpected short-lived instances). Default age-off is exponential: importance decays faster the more ephemeral the resource. Anomalous patterns trigger alerts before age-off completes.

6. **All abstraction levels, down into the OS and service layer.** Individual resources, services, and subsystems are all modeled. The graph extends below the container boundary into operating system and service internals. Queries operate at any altitude without forcing a granularity choice upfront.

## Consequences

- The property graph schema evolves organically from query needs, not from an upfront modeling exercise. Periodic grooming cycles maintain coherence.
- MVP discovery scope is deliberately narrow (SCAP + Kubernetes + cloud APIs). Network traffic observation and other high-friction discovery methods are deferred but architecturally accommodated.
- Drift detection becomes a set-theoretic delta between expected state (requirements) and observed state (runtime inventory), requiring stable identifier schemes to link analytical model entities to runtime discoveries.
- Ephemeral resource handling adds complexity to the graph lifecycle. Age-off policies, frequency analysis, and anomaly detection are operational concerns that must be designed into the discovery layer.
- Full-depth modeling (below the container boundary) increases graph size and query complexity but eliminates blind spots that container-level-only modeling creates.
- Formal model projections require maintained translation adapters, each of which is a regression-tested contract between the graph schema and the target format.
