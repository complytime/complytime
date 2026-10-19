# Provenance and Trust Architecture

Compliance decisions have consequences: authorization to operate, risk acceptance, audit findings, legal liability. The provenance behind these decisions must be independently verifiable, tamper-evident, and recoverable after the people who made them move on. Today, provenance is an afterthought: audit logs exist but are single points of falsification, reasoning is captured in prose (if at all), and the link between a decision and its evidence is maintained by institutional memory.

The trust problem has three dimensions:

1. **Authenticity**: did the claimed actor actually produce this record? (Solved by cryptographic signatures, but rarely applied to compliance operations.)
2. **Integrity**: has the record been altered since it was produced? (Solved by content-addressable storage and hash chains, but compliance records are typically stored in mutable databases.)
3. **Non-equivocation**: can the producer show one version of the record to one party and a different version to another? (Requires independent witnesses or transparency logs; this is the hardest property to achieve.)

Each dimension has established solutions in other domains (supply chain security, certificate transparency, distributed systems). The challenge is composing them into a coherent architecture that serves compliance without imposing prohibitive operational overhead.

A further complication: compliance operations increasingly involve LLM-assisted analysis (mapping proposals, risk scoring, evidence summarization). The provenance chain must extend through these non-deterministic steps. It must capture what was decided, what inputs were provided, what model produced the analysis, and what reasoning path was followed. Without this, there is no way to audit the machine's contribution to a compliance decision.

## Current approaches and prior art

### Traditional audit logs

Application-level logging to databases or files. Mutable, single point of failure, no cryptographic integrity. Standard practice but insufficient for adversarial audit scenarios. An administrator with database access can alter history without detection. Logs are often retained based on storage capacity rather than compliance retention requirements, and there is no mechanism to prove a log was not selectively deleted.

### Transparency logs

Certificate Transparency (CT), [Sigstore Rekor][sigstore], and [Tessera][tessera] implement append-only, cryptographically verifiable logs based on Merkle trees. These systems provide integrity and non-equivocation. CT is production-proven at internet scale. Every public TLS certificate issued since 2018 must be logged to multiple independent transparency logs before browsers will trust it. The model works: CAs cannot issue rogue certificates without public detection.

Tessera (Google/open-source) provides a transparency log framework that generalizes the CT design for arbitrary record types. Sigstore Rekor applies the pattern to software supply chain: every artifact signed via Sigstore produces a Rekor log entry with a cryptographic proof of inclusion. The transparency log becomes the trust anchor. Queryable stores are rebuildable caches; if they diverge from the log, the log wins.

### [in-toto][in-toto]

A CNCF graduated framework for securing software supply chains through layouts (authorized steps) and links (per-step attestations with cryptographic digests of inputs and outputs). Originally designed for build pipelines: a layout defines "the build must execute these steps, in this order, by these actors," and each step produces a signed link recording what it consumed and produced. At verification time, the full chain is reconstructed and validated.

The attestation model generalizes to any multi-step process with defined inputs and outputs. A compliance assessment is exactly this: a sequence of steps (requirement interpretation, evidence collection, evaluation, decision) performed by specific actors (human, automated, or hybrid) with defined inputs (requirements, policies, system state) and outputs (assessment results, attestations, risk scores).

### [OpenTelemetry][opentelemetry]

A CNCF graduated observability framework providing distributed tracing, metrics, and structured logging. Not designed for compliance provenance, but the trace ID propagation and span model provide the "what happened when" layer that complements cryptographic proof. Massive ecosystem adoption means instrumentation and correlation infrastructure already exists in many environments.

A trace captures the operational narrative: what operations were performed, in what order, with what latency, and where they failed. A span records context (the request ID that triggered the operation, the service that performed it, the downstream calls it made). For compliance, this becomes the audit trail of "how did we reach this decision," especially valuable when that decision involved multiple services, human input, and LLM-assisted reasoning steps.

### [NATS JetStream][nats-jetstream]

A distributed messaging system with persistent streams. Provides durable, ordered, replayable event sequences. Not a transparency log (no Merkle proofs), but provides the persistent audit stream layer. Streams can be configured with retention policies (time, message count, storage limits) and replication for durability. Consumers can replay from any offset, enabling after-the-fact audit reconstruction.

The stream becomes a durable record of compliance events. An auditor examining a decision made six months ago can replay the stream to see every input, every intermediate step, and every output (assuming the retention policy kept it).

### Witness cosigning

Independent parties countersign transparency log checkpoints to prevent equivocation. Omniwitness and witness.dev extend the transparency log model with distributed trust. A single transparency log operator could equivocate by showing different views to different clients. Witness cosigning solves this: witnesses independently verify the log's consistency and countersign checkpoints. A client trusts the checkpoint only if a quorum of independent witnesses attests to it.

This model requires operational independence: witnesses must be run by different organizations, in different jurisdictions, with different incentive structures. For compliance, this could mean industry consortia, regulatory bodies, or independent auditors operating witness nodes that countersign compliance log checkpoints.

## Proposed approaches

A provenance triad: three independent records of every compliance-critical operation, joined by a single correlation ID (OpenTelemetry trace ID).

| Record                       | Purpose                                  | Technology pattern                                                                         |
|:-----------------------------|:-----------------------------------------|:-------------------------------------------------------------------------------------------|
| Telemetry trace              | What happened and when                   | OpenTelemetry spans with compliance-specific attributes                                    |
| Persistent audit stream      | Durable, ordered, replayable record      | Append-only stream with mandatory provenance headers and content hashing                   |
| Cryptographic attestation    | Mathematical proof of inputs and outputs | in-toto layouts and links with embedded trace ID                                           |

Tamper evidence comes from the triad as a whole: altering one record leaves two independently held counterparts intact. An auditor can start at any one record and reach the other two via the shared trace ID. If the telemetry trace says "requirement R-042 was evaluated at 2025-01-15T14:32:00Z with result PASS," the audit stream must contain a corresponding event with the same trace ID, timestamp, and result, and the cryptographic attestation must include a signed link recording the inputs (requirement definition hash, system state hash) and outputs (result, evidence identifiers).

The write path should be fail-closed: operations that cannot be recorded across all three channels are rejected. This trades availability for integrity. Audit infrastructure outages block compliance operations rather than allowing unrecorded decisions. In practice, this means the compliance engine refuses to commit a decision until it has successfully written to the telemetry backend, the audit stream, and the transparency log. If any write fails, the operation is rolled back and the decision is not recorded.

This is operationally demanding. Most systems are designed to tolerate logging failures (log and continue; if the log is unavailable, drop the log entry and proceed). Fail-closed provenance inverts this: the log becomes critical path infrastructure. Outages in the transparency log or audit stream block compliance operations until service is restored or a compensating mechanism (local buffering, degraded-mode operation) is activated.

### Immutable store selection

Transparency logs (Tessera or equivalent) provide append-only guarantees with Merkle inclusion proofs. Witness cosigning adds non-equivocation. The transparency log is the source of truth; queryable stores are rebuildable caches. If a database claims a decision was made but the transparency log has no corresponding entry, the database is wrong (either through corruption, tampering, or desynchronization).

The critical property: anyone with the transparency log can independently rebuild the queryable store and verify it matches. This enables after-the-fact detection of tampering. An auditor who suspects a compliance record was altered can fetch the transparency log checkpoint from multiple witnesses, reconstruct the log locally, and compare it against the presented records. Divergence is proof of tampering.

### LLM-assisted operations

The in-toto attestation captures the model identifier, input hash, and output hash. The telemetry trace captures the full reasoning chain (prompts, intermediate steps, final output). Together they provide the "how did the machine reach this conclusion" audit trail.

For example: an LLM-assisted requirement mapping operation produces an in-toto link recording:

- Inputs: requirement definition (hash), existing mappings (hash), model configuration (hash)
- Outputs: proposed mapping (hash), confidence score, reasoning summary (hash)
- Metadata: model identifier (e.g., `claude-sonnet-4-5@20250929`), temperature, token count
- Trace ID: links to the OpenTelemetry trace capturing the full prompt/response exchange

The attestation proves "this model, with these inputs, produced this output." The trace provides the narrative. An auditor examining the mapping six months later can verify the attestation signature, recompute the input and output hashes to confirm they match, retrieve the full prompt/response exchange from the telemetry backend via the trace ID, and reconstruct the reasoning path.

This does not solve the non-determinism problem. Running the same inputs through the same model will likely produce a different output. But it does provide auditability: the reasoning path that was actually followed is cryptographically attested and independently recoverable.

## Open questions

The provenance triad architecture proposed above raises several questions that remain under discussion. See [ADR-0013](../ADRs/0013-provenance-triad.md) for the decisions that have crystallized from this exploration.

See [Evidence](evidence.md) for the broader evidence lifecycle that provenance serves, and [Requirement Fidelity](requirement-fidelity.md) for why both the decisions and the reasoning behind them must survive.

[sigstore]: https://sigstore.dev/
[tessera]: https://github.com/transparency-dev/tessera
[in-toto]: https://github.com/in-toto/in-toto
[opentelemetry]: https://opentelemetry.io/
[nats-jetstream]: https://docs.nats.io/nats-concepts/jetstream
