# ADR-0009: Provenance Triad for Compliance Operations

**Status:** proposed

**Date:** 2026-06-22

**Deciders:** complytime-dev team

## Context

Compliance decisions carry consequences: authorization to operate, risk acceptance, audit findings, legal liability. The provenance behind these decisions must be independently verifiable, tamper-evident, and recoverable long after the people who made them have moved on.

A single audit log is a single point of falsification. An actor with write access to the log can alter the record of a decision with no independent counterpart to detect the change. This is insufficient for compliance operations where the record itself is the product.

The problem extends to LLM-assisted operations (mapping proposals, risk scoring, evidence analysis). When a machine contributes to a compliance decision, the provenance chain must capture what model was used, what inputs were provided, and what reasoning path was followed. Without this, there is no way to audit the machine's contribution to a compliance decision.

Three independent technologies provide the building blocks: [OpenTelemetry][opentelemetry] (distributed tracing), persistent audit streams (ordered, durable event logs), and [in-toto][in-toto] (cryptographic supply-chain attestations). Each solves one dimension of the provenance problem. None is sufficient alone.

## Decisions

Every compliance-critical operation produces three independent records joined by a single correlation ID (OpenTelemetry trace ID):

| Record | What it provides | Technology pattern |
|:---|:---|:---|
| Telemetry trace | What happened and when | OpenTelemetry spans with compliance-specific attributes |
| Persistent audit stream | Durable, ordered, replayable record with integrity verification | Append-only stream with mandatory provenance headers and content hashing |
| Cryptographic attestation | Mathematical proof of inputs, outputs, and authorization | in-toto layouts (authorized steps) and links (per-step digests) |

An auditor can start at any one record and reach the other two via the shared trace ID. Tamper evidence comes from the triad: altering one record leaves two independently held counterparts intact.

The write path is fail-closed: operations that cannot be recorded across all three channels are rejected. This prioritizes integrity over availability. An audit infrastructure outage blocks compliance operations rather than allowing unrecorded decisions.

Transparency logs provide append-only immutability for the audit stream. Witness cosigning provides non-equivocation (the producer cannot show different versions of the log to different parties). The transparency log is the source of truth; queryable stores built from it are rebuildable caches.

The specific transparency log implementation ([Tessera][tessera], [Rekor][rekor], self-hosted, or other), witness policies, and operational deployment topology are deferred to implementation.

## Consequences

- Three-way cross-verification raises the bar for undetected tampering from "compromise one system" to "compromise three independent systems simultaneously."
- Fail-closed write path means audit infrastructure availability is critical; outages block compliance operations. This is a deliberate trade-off: integrity over availability for compliance-critical paths. Non-critical operations may use a degraded-provenance mode if future requirements demand it.
- OpenTelemetry integration allows existing observability infrastructure to participate in compliance provenance without new tooling.
- in-toto attestations extend supply-chain provenance patterns to compliance operations, a novel application of an established CNCF-graduated framework.
- The triad imposes overhead: three writes per compliance-critical operation. Batch operations (bulk ingestion, large-scale remapping) need efficient batching strategies to avoid write amplification.
- LLM-assisted operations get auditable provenance: the in-toto link captures model ID and I/O hashes; the telemetry trace captures the reasoning chain.
- Privacy constraint: raw sensitive content must not leak into transparency logs that may be broadly readable. A privacy guard must hard-fail if sensitive content would be anchored.

[opentelemetry]: https://opentelemetry.io/
[in-toto]: https://github.com/in-toto/in-toto
[tessera]: https://github.com/transparency-dev/tessera
[rekor]: https://sigstore.dev/
