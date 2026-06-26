# ADR-0008: Evaluator Interface Contract

**Status:** proposed

**Date:** 2026-06-24

**Deciders:** ComplyTime maintainers

## Context

The [evaluator-coupling](../problems/evaluator-coupling.md) problem doc identifies four open questions about how evaluators interact with the platform. The [system-modeling-and-drift](../problems/system-modeling-and-drift.md) problem doc proposes the property graph as the primary data substrate. [ADR-0004](0004-grpc-provider-plugin-architecture.md) established gRPC as the plugin protocol but did not specify what crosses the wire (what the platform sends to evaluators and what evaluators return. This ADR fills that gap.

## Decision

This decision has three parts.

**Engine diversity is accepted.** Evaluation logic stays engine-native. The platform does not define a portable evaluation language or DSL. Evaluator authors work in whatever language their engine requires: Rego, OVAL, CEL, Ruby DSL, Python, or custom logic. The platform owns routing (discovering and dispatching to evaluators), orchestration (managing evaluation lifecycle), and result unification (normalizing findings into a common evidence format).

**Default interface is message-passing with opt-in graph access.** The platform sends each evaluator a scoped payload containing the system state relevant to its check. Evaluators return structured findings. This is the default and only interface for most evaluators. 

Evaluators that require cross-entity reasoning (checks spanning multiple system components) can request graph access. The platform grants read-only access to an explicitly scoped subgraph. Graph access requires authorization, since evaluators are untrusted by default. The scope, duration, and audit trail of graph access are platform-controlled.

**Composition ownership is split between requirement and content.** When a requirement spans multiple evaluation forms (intent and behavioral, automated and manual), the requirement model (Gemara) declares what evaluation forms are needed. The assessment logic stream ([ADR-0005](0005-two-stream-content-model.md)'s second content stream) declares how to compose the results of specific checks (both-must-pass, either-suffices, behavioral-overrides-intent, or custom logic. The Runtime Client's Scan component executes these composition rules but does not define them.

## Consequences

- Evaluator authors write checks in their native language; no transpilation or abstraction layer to build or maintain.
- The platform must define a scoped-payload schema (what the platform sends) and a findings-return schema (what evaluators return). These schemas are not defined in this ADR and will require their own design work.
- Graph access grants need an authorization model: scope definition, read-only enforcement, and audit trail.
- Composition rules become a first-class artifact in the assessment logic content stream.
- The Gemara requirement model needs a field for declaring required evaluation forms.
- The Runtime Client's merge logic must support pluggable composition rulesbeyond simple aggregation.
