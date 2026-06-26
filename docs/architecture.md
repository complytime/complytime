# Architecture

Component vocabulary for ComplyTime. This document names the functional roles and describes how they relate (not how they work internally or which repositories implement them.

## Component Map

```
┌──────────────────────────────────────┐  ┌─────────────────────────┐
│          Content Registry            │  │  Cross-Framework        │
│  ┌──────────────┐  ┌──────────────┐  │  │  Mapping                │
│  │ Compliance   │  │ Assessment   │  │  │                         │
│  │ Content      │  │ Logic        │  │  │  Relationship graph,    │
│  │ (policies,   │  │ (evaluation  │  │  │  standard comparison,   │
│  │  catalogs)   │  │  packages)   │  │  │  traceability           │
│  └──────┬───────┘  └──────┬───────┘  │  └─────────────────────────┘
└─────────┼─────────────────┼──────────┘
          │                 │
          ▼                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Runtime Client                              │
│                                                                 │
│  ┌──────────┐  ┌──────────────┐  ┌───────────────────────────┐  │
│  │ Content  │  │ Policy       │  │ Scan                      │  │
│  │ Pull     │  │ Resolution   │  │ (orchestrate evaluators,  │  │
│  │          │  │ (compose,    │  │  merge results)           │  │
│  │          │  │  flatten)    │  │                           │  │
│  └──────────┘  └──────────────┘  └────────────┬──────────────┘  │
└────────────────────────────────────────────────┼────────────────┘
                                                 │
                         ┌───────────────────────┼───────────────┐
                         │                       │               │
                         ▼                       ▼               ▼
                ┌──────────────┐      ┌──────────────┐   ┌────────────┐
                │ Evaluator    │      │ Evaluator    │   │ Evaluator  │
                │ (filesystem  │      │ (API-based   │   │ (policy-   │
                │  scanning)   │      │  assessment) │   │  as-code)  │
                └──────┬───────┘      └──────┬───────┘   └─────┬──────┘
                       │                     │                 │
                       └─────────────────────┼─────────────────┘
                                             │ evidence
                                             ▼
                              ┌──────────────────────────┐
                              │   Evidence Platform      │
                              │                          │
                              │   Ingestion, storage,    │
                              │   verification, posture  │
                              └──────────┬───────────────┘
                                         │
                                         ▼
                              ┌──────────────────────────┐
                              │   Audit Preparation      │
                              │                          │
                              │   Analysis, artifact     │
                              │   drafting, reporting    │
                              └──────────────────────────┘
```

## Roles

### Content Registry

Distributes two independent streams: compliance content (what must be true) and assessment logic (how to verify it). Both streams use the same transport mechanism and have independent lifecycles.

### Runtime Client

Pulls content, resolves policies, discovers evaluators, orchestrates scans, merges results. Never interprets evaluator-specific content — routes by metadata only.

### Evaluators

Standalone processes that receive scoped system-state payloads from the Runtime Client and return structured findings. Each evaluator owns its evaluation logic in its native engine language. The default interface is message-passing; evaluators requiring cross-entity reasoning can request read-only access to a scoped subgraph of the property graph. The Runtime Client never participates in evaluation. See [ADR-0008](ADRs/0008-evaluator-interface-contract.md) for the interface contract.

### Evidence Platform

Ingests, stores, and verifies compliance evidence produced by evaluators. Provides posture analytics and evidence traceability. See the [Evidence](problems/evidence.md) problem doc for the domain exploration.

### Cross-Framework Mapping

Compares compliance standards, maps relationships between requirements across frameworks, and stores the complete graph with traceability. Enables organizations assessed against multiple frameworks to understand overlapping requirements without manual spreadsheet exercises. See the [Cross-Framework Mapping](problems/cross-framework-mapping.md) problem doc for the domain exploration.

### Audit Preparation

Consumes stored evidence to support audit activities — analysis, artifact drafting, and reporting. Downstream of the evidence platform.

## Current State

| Role | Status | Implementation |
|:---|:---|:---|
| Content Registry | Operational | OCI registries via oras-go |
| Runtime Client | Operational | [complyctl](https://github.com/complytime/complyctl) |
| Evaluators | Operational | [complytime-providers](https://github.com/complytime/complytime-providers) (OpenSCAP, AMPEL); OPA in development |
| Cross-Framework Mapping | Experimental | [crosscodex](https://github.com/complytime-labs/crosscodex) |
| Evidence Platform | Experimental | [complytime-core](https://github.com/complytime/complytime-core) |
| Audit Preparation | Experimental | [complytime-studio](https://github.com/complytime-labs/complytime-studio) |

Supporting repositories: [complypack](https://github.com/complytime/complypack) (pack authoring), [complytime-policies](https://github.com/complytime/complytime-policies) (published bundles), [org-infra](https://github.com/complytime/org-infra) (CI/CD), [community](https://github.com/complytime/community) (governance).

## Boundaries

| Boundary      | Left              | Right             | Interface           |
|:--------------|:------------------|:------------------|:--------------------|
| Distribution  | Content registry  | Runtime client    | Content pull        |
| Orchestration | Runtime client    | Evaluators        | Plugin interface    |
| Graph access  | Runtime client    | Evaluators        | Scoped subgraph query (opt-in) |
| Evidence      | Evaluators        | Evidence platform | Evidence submission |
| Audit         | Evidence platform | Audit preparation | Evidence query      |
