# ADR-0014: Signed Evidence Attestation Pipeline (Engine-Agnostic)

**Status:** proposed

**Date:** 2026-06-24

**Deciders:** ComplyTime developers

> 🤖 LLM WARNING 🤖
>
> This was written with LLM (AI) assistance and human review.
>
> 🤖 LLM WARNING 🤖

> **TL;DR:** Stop trusting the pipe; trust the payload. Today complyctl resolves a shared OIDC secret and has its plugins push evidence over OTLP, which puts one replayable secret on every host and ties delivery to a controllable IdP. The new model has complyctl sign and encrypt each evidence batch (DSSE-wrapped in-toto, claim-check to a content-addressed store); an untrusted ETL or streaming engine of the operator's choice moves it; the consumer verifies signature, identity, freshness, and integrity at its edge. Trust rides in the payload, so the transport is interchangeable and runs air-gapped, and both the shared secret and the external-IdP dependency go away. Trust anchors and recipient encryption keys ship in a signed, self-sovereign per-operator bundle.
>
> **Know before you adopt:** this is the portable delivery layer, not the whole story. Ordering, non-equivocation, and global completeness still need the transparency-log sink, which is mandatory for the compliance profile. Verification is non-bypassable only when you control or attest the consumer. FIPS mode and the compliance profile cannot both hold end to end today, because the transparency-layer crypto is Ed25519-only, so fail closed. The new moving parts: an ingester component, a signed trust bundle with enrollment, and mandatory payload encryption.

## Related Tickets

- [chore: remove collector and export infrastructure](https://github.com/complytime/complyctl/issues/606)
  - Created to stabilize `complyctl` while the transport mechanism (this document) is sorted out
- [Remove Compass enrichment pipeline (truthbeam client, mock-compass, compose profiles)](https://github.com/complytime/complytime-collector-components/issues/326)
  - Removes the portions of the collector that are arguably more than generic pipeline tooling can accomplish
  - A post-review pushes for abandoning the project in favor of reusing existing processing pipeline technologies

## Context and Problem Statement

complyctl produces compliance evidence (OSCAL, SARIF, Gemara) that must reach a downstream consumer (a collector or ingester). Today complyctl resolves a shared OIDC client-credentials secret (`collector.auth.client-secret`, resolved in [`resolveOIDCToken`](https://github.com/complytime/complyctl/blob/main/cmd/complyctl/cli/scan.go); configured via [`AuthConfig`](https://github.com/complytime/complyctl/blob/main/internal/complytime/config.go)) and hands the resulting bearer token to its scan plugins, which push the evidence over OTLP (via the ProofWatch library). That design couples two separate concerns into one channel (whether the sender is *allowed to connect*, and whether the evidence is *authentic*) and pins trust to a long-lived, replayable secret that is distributed to every host that runs a scan. , The model also assumes a reachable, controllable identity provider (IdP). In some target deployments the operator does not control the IdP and it offers no self-service onboarding for new clients, and some deployments are air-gapped. The downstream verifier accepts only a JWKS-discoverable Bearer JWT today, with no mTLS or static-key path (see [`deployment-profiles-and-constraints.md`](../deployment-profiles-and-constraints.md)). Together, these constraints limit the value of further hardening the shared push credential and motivate reconsidering how evidence is trusted in transit.

Making no decision leaves a shared bearer secret on every host, no per-workload identity, no payload-level provenance, and no path to air-gapped or vendor-neutral operation, all of which are limitations for a compliance-evidence tool.

## Decision Drivers

- Run anywhere, vendor-neutral: hosts, containers, cloud, and fully air-gapped ("anyone, anywhere").
- Decouple evidence delivery from transport authentication; remove the shared bearer secret.
- Reuse maintained FOSS / cloud-native supply-chain tooling; minimize bespoke transport and cryptography code.
- Meet compliance integrity needs (authenticity, trusted time, tamper-evidence) and remain FIPS-capable.
- Operate without control over an external IdP and without self-service onboarding.
- Let producer and consumer scale and operate independently, including offline/buffered delivery.

## Considered Options

- A. Signed, self-attesting evidence over an untrusted, engine-agnostic transport ("trust the payload, not the pipe").
- B. Keep the push model, replace the shared secret with stronger client authentication (pluggable `private_key_jwt` / Vault-brokered / mTLS).
- C. Adopt a single-vendor managed pipeline (for example AWS EventBridge/Kinesis/SQS) as the delivery mechanism.
- D. Write directly into complytime-core's existing transparency log / ingest as the only delivery path.

## Decision Outcome

Chosen option: **A: signed, self-attesting evidence over an untrusted, engine-agnostic transport**, because it decouples delivery from trust and is vendor-neutral and air-gap-capable while reusing established FOSS supply-chain tooling. This satisfies the run-anywhere, remove-the-secret, and reuse-FOSS drivers.

In this model complyctl signs evidence and encrypts it to the consumer. Next, an untrusted ETL/streaming engine of the operator's choosing moves it. Finally, the consumer verifies signature, identity, freshness, and integrity at its edge before use. The wire contract is a [DSSE](https://github.com/secure-systems-lab/dsse)-wrapped [in-toto](https://github.com/in-toto/attestation) statement using a claim-check with two digests per item: the signed `subject` carries the **plaintext** digest (the end-to-end integrity check, validated after decryption), while a **ciphertext digest** in the predicate content-addresses the encrypted blob in the store (the in-transit integrity check, confirmable without decrypting). Verifying the signature, matching the ciphertext digest on fetch, then matching the plaintext digest after decryption chains integrity end to end. Verification trust anchors are pluggable ([Sigstore](https://www.sigstore.dev/) keyless, X.509/[SPIFFE](https://spiffe.io/), or pinned keys). Trust anchors and recipient encryption keys are distributed self-sovereignly per operator via a signed, versioned bundle over a pluggable source ([TUF](https://theupdateframework.io/) / OCI / SPIFFE-federation / file), bootstrapped from a pinned root.

Scope note: this pipeline is the portable delivery layer that feeds complytime-core's transparency-log and witness stack; it does not replace it. The transparency-log sink is required for the compliance profile.

Beyond the transport pivot, choosing Option A also commits this ADR to several sub-decisions, each detailed in Appendix A and surfaced here so they are reviewed as decisions rather than buried:

- the self-sovereign trust- and key-distribution model (signed `TrustBundle` + pluggable `BundleSource`, manual/gitops enrollment baseline, fail-closed-with-grace);
- mandatory end-to-end payload encryption and its recipient-key distribution;
- the claim-check wire format with two digests per item and a frozen `evidence/v1` predicate (including an envelope nonce);
- topic + time-based retention (not ack-GC), with the transparency-log sink as the durable system of record;
- a new consumer/ingester component whose home is still open (see the design's Open Decisions).

### Consequences

- **Positive:** removes the shared bearer secret and the dependency on a controllable external IdP; trust travels with the payload, so the transport can be any FOSS engine or cloud service and can run air-gapped; producer and consumer decouple and scale independently; reuses DSSE/in-toto/Sigstore/TUF rather than bespoke crypto or transport; the producer-identity problem becomes a signing-identity problem with mature protection options (TPM, Vault, SPIFFE, Sigstore); payload confidentiality is end-to-end, though envelope metadata (producer identity, timing, blob digests/sizes) is intentionally cleartext so the transport can route — an acknowledged residual leak.
- **Negative:** integrity guarantees beyond a single delivered item (tamper-evident ordering, non-equivocation, and global completeness/non-suppression) hold only when paired with the transparency-log sink; in a bring-your-own-engine deployment the project can enable but not force verification; the design adds new operational surface (a signed trust bundle, its distribution, and an enrollment workflow) plus mandatory payload encryption with its own key distribution; keyless signing is FIPS-friendly by default (cosign/Fulcio default to ECDSA P-256), but the mandatory transparency-log sink's checkpoint and cosignature crypto is Ed25519-only today, so FIPS mode and the compliance profile cannot both hold end-to-end until that layer gains a FIPS-validated algorithm; a new consumer/ingester component must be built and operated; and it is more moving parts than continuing to push to a single endpoint.

## Pros and Cons of the Options

### A. Signed evidence over untrusted, engine-agnostic transport

**Pros**

- Trust rides in the payload, so the transport is interchangeable (a FOSS engine, a cloud pipeline, or a plain loop) and works air-gapped.
- Removes the shared secret and the external-IdP dependency; identity becomes a signing key with mature protection options.
- Reuses standard supply-chain tooling and complyctl's existing OCI machinery, and introduces one new pluggable signing-key provider rather than bespoke cryptography.

**Cons**

- Ordering, non-equivocation, and completeness need the additional transparency-log sink; the pipeline alone does not provide them.
- Adds trust-bundle distribution, enrollment, and mandatory encryption as new operational concerns.
- Verification cannot be forced on operator-owned consumers.

### B. Keep push, strengthen client authentication (`private_key_jwt` / Vault / mTLS)

**Pros**

- Smallest change; preserves the existing OTLP path.
- `private_key_jwt` ([RFC 7523](https://datatracker.ietf.org/doc/html/rfc7523)) and mTLS ([RFC 8705](https://datatracker.ietf.org/doc/html/rfc8705)) are asymmetric and non-replayable, an improvement over a shared secret.

**Cons**

- Trust stays channel-bound rather than payload-bound; evidence is not independently verifiable later, and delivery still needs a synchronous, online verifier (no offline/air-gap buffering).
- mTLS or static-JWKS at the collector require an upstream complytime-core change, and several IdPs cannot self-onboard the required client registrations.

### C. Single-vendor managed pipeline (for example AWS)

**Pros**

- Durability, retries, and fan-out are turnkey and provider-operated.

**Cons**

- Does not satisfy the vendor-neutral / run-anywhere driver; cannot serve air-gapped or non-AWS adopters.
- Still needs a payload-trust model on top; the managed pipe does not establish provenance.

### D. Write directly into complytime-core's transparency log / ingest only

**Pros**

- Reuses the existing system of record and its ordering and non-equivocation properties directly.

**Cons**

- Re-creates the tight coupling and the ingest-authentication problem this decision removes.
- Forces every producer to reach Core synchronously: no vendor-neutral or air-gapped delivery, and no decoupled buffering.

## More Information

- Companion design document (full architecture, threat model, open decisions): [Appendix A](#appendices), embedded below.
- Deployment profiles and cross-cutting constraints (grounds the JWT-only verifier claim and the air-gap, FIPS, and transparency-log context): [Appendix B](#appendices), embedded below.
- Superseded predecessor design: the push + OIDC `client_secret` approach is summarized under Considered Options B above and retained in the project's design-spec history (outside this repo); not embedded here.
- Current implementation being replaced: complyctl [`cmd/complyctl/cli/scan.go`](https://github.com/complytime/complyctl/blob/main/cmd/complyctl/cli/scan.go), [`internal/complytime/config.go`](https://github.com/complytime/complyctl/blob/main/internal/complytime/config.go), [`api/plugin/plugin.proto`](https://github.com/complytime/complyctl/blob/main/api/plugin/plugin.proto).
- Standards and tooling: [DSSE](https://github.com/secure-systems-lab/dsse), [in-toto attestations](https://github.com/in-toto/attestation), [Sigstore](https://www.sigstore.dev/) ([Fulcio](https://github.com/sigstore/fulcio), [Rekor](https://github.com/sigstore/rekor), [cosign](https://github.com/sigstore/cosign)), [The Update Framework (TUF)](https://theupdateframework.io/), [SPIFFE/SPIRE](https://spiffe.io/), [age](https://github.com/FiloSottile/age), [RFC 3161 trusted timestamping](https://datatracker.ietf.org/doc/html/rfc3161).
- Candidate transport engines (choice deferred): [Benthos](https://github.com/redpanda-data/benthos) (MIT), [Redpanda Connect](https://github.com/redpanda-data/connect) (Redpanda Community License, with an Apache-2.0 free bundle), [Bento](https://github.com/warpstreamlabs/bento) (community fork).
- Template: [MADR 4.0](https://adr.github.io/madr/).

## Appendices

The supporting material below is embedded so this ADR is self-contained. It is reference detail; the decision itself is stated above.

<details>
<summary>Appendix A: design of the signed evidence attestation pipeline (full architecture, threat model, open decisions)</summary>

# Signed Evidence Attestation Pipeline (engine-agnostic)

- **Date:** 2026-06-23
- **Target repos:** `github.com/complytime/complyctl` (producer) + a new consumer component (the ingester); a shared security-core module
- **Status:** Proposed; hardened after adversarial review + trust/key distribution brainstorm (2026-06-23); spec only (no implementation plan requested)
- **Supersedes:** `2026-06-23-complyctl-pluggable-collector-auth-design.md`

## Problem

complyctl produces compliance evidence (OSCAL / SARIF / Gemara) that must reach a downstream consumer. The current implementation **pushes** it over OTLP and authenticates the push with a shared OIDC `client-secret` (`cmd/complyctl/cli/scan.go` `resolveOIDCToken`). That couples *"is the sender allowed to talk to me"* with *"is this evidence authentic,"* forces a synchronous tight coupling, and pins trust to a shared bearer secret.

## Pivot: trust the payload, not the pipe

Decouple delivery from trust. complyctl **signs** the evidence; an untrusted transport moves it; the consumer **verifies** signature, identity, and freshness at its edge. Self-attesting evidence makes the transport commodity plumbing — safe to delegate to any ETL/streaming engine. **This is an ETL pipeline — extract, transform, load — *except* the trust boundary**, which stays ours.

## Hard requirements

1. **Anyone, anywhere — vendor-neutral.** Any substrate: cloud object stores, brokers, a filesystem drop, fully air-gapped. A cloud-native pipeline (AWS EventBridge/Kinesis/SQS) is a valid *deployment*, never *the design*.
1. **Verification is non-bypassable only when the consumer deployment is controlled or attested**; otherwise the property is *verifiable by anyone holding the trust anchors* (see §Non-bypassable gate). This is an enabled-and-default property, not a forced one.
1. **The wire format is the contract**, not any engine — cross-language interop.
1. **Confidentiality is mandatory.** Evidence is sensitive; the transport is untrusted, therefore it must not be able to read payloads (see §Confidentiality).
1. **Trusted time is required** (see §Trusted time) — self-asserted timestamps are not sufficient.
1. **No hand-rolled crypto.** Maintained FOSS supply-chain libraries only.
1. **FIPS mode (optional)** must be available and correct end-to-end when enabled.

## Threat model & guarantees

**What the pipeline guarantees (per delivered item):** integrity (tamper-evidence), authenticity (who signed), confidentiality (E2E encryption), trusted time, and **intra-batch completeness**.

**What it does NOT guarantee on its own — explicit non-goals:**

- **Evidence *truth*.** Signing proves origin, not correctness. A compromised scanner signs "compliant" with a valid key and the consumer accepts it. Detection is out of band (host integrity, the transparency log).
- **Forced verification.** In a bring-your-own-engine model the operator owns the consumer; we cannot *force* verification, only *enable* it. The property is **"verifiable by anyone holding the trust anchors,"** not "always verified." Forcing it requires controlling/attesting the consumer deployment.
- **Global completeness / non-suppression.** A whole batch dropped before the consumer learns it exists, or a sender that never enrolled, is undetectable by the pipeline. Only the transparency-log sink + an external reporting-cadence expectation closes this.
- **Tamper-evident ordering / non-equivocation.** Provided by the log + witnesses, not by the queue.

**Guarantees with vs without the transparency-log sink:**

| Guarantee                             | Pipeline alone | + log sink |
| ------------------------------------- | :------------: | :--------: |
| Integrity (delivered item)            |       ✓        |     ✓      |
| Authenticity                          |       ✓        |     ✓      |
| Confidentiality (mandatory E2E)       |       ✓        |     ✓      |
| Trusted time (RFC 3161 TSA + Rekor)   |       ✓        |     ✓      |
| Intra-batch completeness              |       ✓        |     ✓      |
| Global completeness / non-suppression |       ✗        |     ✓      |
| Tamper-evident ordering               |       ✗        |     ✓      |
| Non-equivocation                      |       ✗        |     ✓      |

**Structural decision:** this pipeline is the **portable delivery layer that feeds** complytime's transparency-log / witness stack — it does not replace it. The **transparency-log sink is mandatory for the compliance profile** and optional only for non-compliance uses.

## Architecture: engine-agnostic security core + two edges

The shared module owns *only* the trust boundary. Transport (topics, blob stores, connectors, retries, fan-out, dead-letter, backpressure) belongs to the ETL engine.

```
┌────────────┐  Seal (sign+encrypt) ┌──────────────────────────┐  Open (verify+decrypt) ┌──────────────┐
│  complyctl │ ───────────────────► │  ETL engine + connectors  │ ─────────────────────► │  ingester    │
│ (producer) │   envelope + blob    │  (Benthos/Bento/AWS/Go)   │   verified plaintext   │ (consumer)   │
└────────────┘                      │  — untrusted transport —  │                        └──────┬───────┘
  key: TPM/Vault/SPIFFE/KMS         └──────────────────────────┘        sinks: log (mandatory, compliance) + OTLP/DB
```

```go
// Producer edge
type Signer interface {
    Seal(ctx context.Context, batch []Evidence, meta Predicate) (SignedEnvelope, []Blob, error)
}
// Consumer edge — engine injects how to fetch claim-check blobs
type Verifier interface {
    Open(ctx context.Context, env SignedEnvelope, fetch BlobFetcher) ([]Evidence, VerifiedIdentity, error)
}
type BlobFetcher func(ctx context.Context, d Digest) ([]byte, error)
```

We define **no** `Queue`/`BlobStore` interface — the engine's connectors own transport. The only injected port is `BlobFetcher`. Two edges + one callback = the whole trusted surface.

## Wire contract (the reusable artifact)

A **claim-check batch**: one small signed envelope per scan run travels the pipeline; bulk evidence lives content-addressed in a store.

- Envelope = **DSSE** wrapping an **in-toto Statement**.
- `payloadType`: `application/vnd.in-toto+json`
- `predicateType`: `https://complytime.org/evidence/v1` *(versioned interop contract)*
- `subject`: **`[{ name, digest: { sha256: <hex> } }, …]` — one entry per evidence item in the batch, where the digest is over the item's *plaintext*.** This *is* the signed batch manifest: the consumer learns exactly which N items to expect and detects any missing one (**intra-batch completeness, §below**). The plaintext digest is the **end-to-end** integrity check, confirmed *after* decryption (§Confidentiality).
- `predicate`: `{ batch-id, nonce, producer, signed-time-source, schema-version, encryption: { recipients/alg }, per-item: { content-type, ciphertext-digest: { sha256 } } }` — kept **minimal** (metadata is plaintext; see §Confidentiality). The `nonce` makes every envelope identity unique (§Dedup); the per-item **ciphertext digest** content-addresses the encrypted blob in the store and is the **in-transit** integrity check — fetch blob, recompute sha256, confirm match — *without* decrypting.
- `signatures`: DSSE sigs; cert chain / SPIFFE SVID / Sigstore bundle (incl. Rekor inclusion) carried cosign-style.

Versioned + documented for cross-language interop; proven in tests against the `cosign` CLI. **Unknown `predicateType` → reject-safe** (no crash, never trust-by-default).

## Completeness (refined — batch-level, sender-agnostic)

- **Intra-batch completeness (primary):** the multi-subject envelope above is a signed manifest of the batch. Detects any item dropped from a *received* batch. Needs no foreknowledge of senders. ✅
- **Per-enrolled-sender sequencing (optional):** for deployments that register producers, a monotonic sequence per signing identity + an expected cadence enables cross-batch gap detection and liveness alerting. Not viable in a fully open model (unknown senders) — hence optional.
- **Whole-batch / unknown-sender suppression (out of scope for the pipeline):** closed only by the transparency-log sink + external reporting-cadence expectations.

## Trusted time & transparency (required)

Self-asserted predicate timestamps are not trusted. The **default/compliance profile requires Sigstore** (Fulcio short-lived certs + **Rekor** for transparency/inclusion). A plain Rekor *inclusion* timestamp comes from Rekor's own clock and is **not** an independently trusted time anchor, so trusted time is anchored by an **RFC 3161 TSA** (the default in Rekor v2) in *both* the connected and air-gapped profiles. Short-lived Fulcio certs plus a TSA-anchored "valid at signing time" proof also largely sidestep long-lived-key revocation (§Identity).

- **Air-gap:** run Sigstore (Fulcio/Rekor) **in-enclave**, or substitute an **RFC 3161 TSA**. `x509`/`pinned` policies remain available there but carry the weaker time story unless paired with a TSA.

## Identity, keys, revocation

- **Producer signing key** (a new pluggable signing-key provider this design introduces): file / TPM-PKCS#11 / **Vault Transit** (no key egress) / **SPIFFE SVID** / cloud KMS / **Sigstore keyless** (Fulcio). This is where the producer-identity question lands — as the *signing* identity, not a transport credential.
- **Revocation / rotation / trust-anchor & recipient-key distribution** are handled by the **§Trust & key distribution** design below (self-sovereign signed bundle + pluggable source + operator-owned enrollment). Sigstore's short-lived certs materially reduce the revocation burden.

## Trust & key distribution

Solves the dual problem in one mechanism: consumers need producers' **verifier anchors**; producers need consumers' **recipient encryption keys**. Both are *public* material, so distribution needs only **integrity + freshness**, never confidentiality.

**Authority: self-sovereign.** Each operator (a trust-domain) runs its own root and decides who it trusts. complytime ships the mechanism, never a global trust list.

**The `TrustBundle`** — one signed, versioned artifact per trust-domain: `version · issued-at · expires-at (freshness) · epoch (monotonic, rollback) · trust-domain · verifier-anchors[] (Fulcio+Rekor roots / CA & SPIFFE bundles / pinned keys + identity→scope policy) · recipient-keys[] (consumer age/KMS public keys + scope) · revoked[]`. Signed by the operator's bundle-signing key, chaining to the **pinned root**. Carries both directions in one document. **Contains no secrets.**

**Pluggable `BundleSource`** — `Fetch(ctx) (SignedBundle, FreshnessProof, error)`:

| Backend             | Freshness/rollback strength                                                    | Notes                                   |
| ------------------- | ------------------------------------------------------------------------------ | --------------------------------------- |
| `tuf`               | strongest (rotation, freshness, rollback, threshold, repo-compromise recovery) | reuses the Sigstore-mandated TUF client |
| `oci`               | medium (cosign-signed artifact; expiry via tags/metadata)                      | reuses complyctl OCI                    |
| `spiffe-federation` | medium (anchors only; no recipient keys)                                       | SPIFFE-centric fleets                   |
| `file` / `gitops`   | basic (signed file; freshness = git + expiry)                                  | air-gap floor                           |

Authenticity always = the bundle's signature to the pinned root; the source only adds freshness/rollback (tiered, documented like the engine). **Bootstrap = the pinned root, set out-of-band at install** — the one irreducible trust-on-install (TUF `root.json`, or a pinned bundle-signing key for `oci`/`file`). Air-gap = mirror metadata/targets as files; the bundle crosses the diode like any signed artifact.

**Enrollment (self-sovereign, pluggable proofing)** — the onboarding the external IdP denied us, now operator-owned:

1. Subject generates material (signing keypair / SPIFFE SVID / Sigstore identity; or encryption keypair).
1. Submits it with a **proof** — `manual`/`gitops` PR (baseline, always available) · `join-token` · platform attestation (TPM/cloud/k8s via SPIRE) · `oidc-at-enrollment` (the operator's own IdP, one-time, not per-message).
1. Operator approves (human or policy) → entry added → bundle re-signed, republished.
1. Propagates on the next freshness refresh.

The proofing method's strength **is** the assurance level (tiers documented); pluggable `EnrollmentProof` on the operator side, `manual`/`gitops` always present.

**Revocation / rotation:** prefer **short-lived everything** (Sigstore keyless / SPIRE SVIDs) so revocation ≈ stop re-issuing + Rekor proves validity-at-time; explicit `revoked[]` + epoch bump for pinned/long-lived keys, with freshness preventing a stale trusting bundle from being force-fed; recipient-key rotation = publish-new + overlap window (old key retained until the retention window lapses for in-flight decrypt).

**Freshness failure (decided): fail-closed with grace + alert.** On inability to refresh before `expires-at`, operate on last-known-good for a bounded, per-profile **grace** window while alarming; after grace, refuse to verify/produce. Fail-open is opt-in only. Tuned looser for air-gap, tighter for connected/compliance.

**Failure modes:** pinned-root compromise = root-level game-over → protect offline, TUF threshold + rotation; bundle-source DoS → cached last-known-good + grace + alert; weak enrollment proofing = the real weak link → assurance tiers documented.

**Integration:** consumer `Verifier` loads anchors from the source (refresh on the freshness window); producer `Signer` loads recipient keys to encrypt to; both pinned-root-bootstrapped; the same crypto core verifies the bundle.

## Verifier policy (pluggable, composable — read the warning)

OR-combined branches, each with an identity constraint; returns `VerifiedIdentity` for audit:

| Policy     | Trust root                      | Identity check                          | Library                     |
| ---------- | ------------------------------- | --------------------------------------- | --------------------------- |
| `sigstore` | Fulcio roots + Rekor            | cert OIDC issuer + SAN regex + log incl | `sigstore/sigstore`, cosign |
| `x509`     | CA bundle / SPIFFE trust domain | chain valid + SAN matches pattern       | `crypto/x509`, `go-spiffe`  |
| `pinned`   | pinned public keys              | keyid ∈ allowlist                       | stdlib                      |

> **WARNING — `OR` is the union of attack surface.** Enabling multiple anchors means the *weakest* one can sign everything; one compromised CA forges all. **Recommended: enable exactly one anchor.** If combining, scope per-producer / per-evidence-type (least trust). **SAN-regex must be anchored** (`^…$`) — an unanchored `.*@corp\.com` is identity confusion.

- **Algorithm allowlist (enforced):** the verifier pins acceptable **signature** and **digest** algorithms and rejects anything else — no downgrade, no attacker-chosen-alg.

## Confidentiality (mandatory)

- **Payload E2E encryption:** evidence blobs are encrypted to the consumer recipient(s) — `age` (X25519) or KMS envelope; **FIPS mode → RSA-OAEP/ECDH-ES + AES-GCM**. The untrusted transport cannot read payloads.
- **Recipient key distribution** is handled by the same mechanism as trust-anchor distribution: the consumer's recipient keys travel in the signed `TrustBundle` (§Trust & key distribution), so producers learn who to encrypt to from the same freshness-checked artifact — no separate channel.
- **Metadata minimization:** the envelope/predicate is plaintext (the engine needs it to route). Keep it minimal; encrypt or omit sensitive predicate fields. **Residual leak** (producer identity, timing, blob sizes) is acknowledged; add padding/batching only if traffic analysis is in scope.

## FIPS mode (optional)

When enabled, restrict the **entire path** — signatures, content digests, and payload encryption — to FIPS-approved algorithms (ECDSA P-256/384/521 or RSA-PSS; SHA-2; AES-GCM). Signing is the easy leg: **cosign and Fulcio already default to ECDSA P-256**, so keyless and `x509` signing are FIPS-friendly out of the box (no switch away from a default needed). The real FIPS cliff is downstream and must not be hidden:

- the **transparency-log sink** the compliance profile mandates signs checkpoints with `golang.org/x/mod/sumdb/note` and witnesses with `cosignature/v1`, both **Ed25519-only** today; the FIPS-capable path is `cosignature/v2` (ML-DSA-44 / FIPS 204), not yet in the pinned libraries.
- **payload encryption** with `age` uses X25519 + ChaCha20-Poly1305, neither FIPS-approved — FIPS mode must select **RSA-OAEP / ECDH-ES + AES-GCM** instead.

Consequently **FIPS mode and the compliance profile cannot both hold end-to-end with today's pinned transparency-layer crypto**; wire the whole path under one switch and **fail closed** if any leg can't comply rather than emit false assurance.

## Delivery: topic + time-based retention (not ack-GC)

Model the transport as a **topic/stream with a time-based retention window**, a deployment parameter sized to the slowest expected consumer. Ack-based GC is unsolvable with unknown, independent fan-out consumers, so we do not attempt it.

- **CAS blob retention ≥ topic retention** so claim-checks never dangle.
- **Durability beyond the window = the transparency-log sink** (system of record).
- Dedup (next §) and retention together replace per-consumer ack tracking.

## Dedup (fixed — on envelope identity, never the blob)

Idempotency keys on the **envelope identity** (hash of the DSSE payload, which includes `batch-id` + trusted time + **the `nonce`**, guaranteeing uniqueness even across scans with identical evidence). **Never** dedup on a subject/plaintext digest — two distinct scans with identical evidence bytes share a plaintext digest, and digest-keyed dedup would silently drop the second scan's proof. Encrypted **blobs** dedup for **storage** on their *ciphertext* digest only (identical plaintext encrypted to different recipients yields distinct ciphertext digests, so cross-recipient storage dedup does not occur — an accepted cost of mandatory encryption).

## Failure modes, ordering, hardening

- Untrusted signer / digest mismatch / missing blob / unknown version / stale time → **reject**, route to dead-letter, record `VerifiedIdentity` + reason; never emit.
- **Verify signature before parsing** the predicate; then **match the ciphertext digest on fetch** (in-transit integrity, no decryption), **decrypt inside the verified boundary**, and **match the plaintext subject digest before parsing** the blob (OSCAL/SARIF parsers are complex — keep them behind verification *and* decryption). The DSSE/JSON envelope parser itself eats unverified input — keep it minimal/hardened. *(Deserialization boundary — validate here regardless of "boundary" rules.)*
- At-least-once delivery → sinks **idempotent by envelope identity**.
- **Deployment hardening (environment, not core capability — noted):** write-authz + rate limits on the topic (untrusted-for-*trust* ≠ unauthenticated-for-*access*), **max blob size + streaming verification**, decompression-bomb guards. These bound DoS; they are operator responsibilities.

## Observability & alerting (security signal)

Rejections **are** the security signal. Required: alert on verification-failure rate, dead-letter growth, **intra-batch gaps**, per-enrolled-sender sequence gaps, and **liveness silence** from expected producers. An unmonitored verify boundary is blind to both attack and suppression.

## Non-bypassable gate — honestly scoped

Wrapping `Open` inside the engine's **input** makes verification non-bypassable *only when you control/attest the consumer deployment*. In bring-your-own-engine deployments the operator can skip it. Ship a correct reference consumer; document that **reading the store without verifying = trusting the pipe**, which the design explicitly does not.

## Engine adapters (thin, deferred)

Engine choice is deferred (engine-agnostic). Core ships three ways over one contract: a **standalone Go runner** (reference + tests + small/air-gap), a **Benthos/Bento plugin pack** (input=`Open`, output=`Seal`; MIT engine + Apache free bundle if chosen), and a **cloud-pipeline** wrapper (e.g. AWS Lambda at the verify edge).

## Testing (TDD)

- **Wire-format golden vectors** + **`cosign` CLI interop**.
- `Signer`/`Verifier` per key source (file/TPM-sim/Vault-dev/SPIFFE/Sigstore) and per policy: accept **and** reject (forged sig, untrusted CA, wrong SPIFFE domain, expired cert, **digest mismatch**, **stale/forged time**, **unknown version**, **disallowed algorithm**, **missing batch item**, **decryption failure**).
- E2E via standalone runner: seal→encrypt→transport→verify→decrypt→sink; tamper, untrusted signer, duplicate (envelope-id dedup), and a dropped-batch-item all rejected/detected.
- FIPS-mode test: Ed25519 path rejected when the switch is on.
- Match complyctl's std-`testing` conventions; no new framework.

## Dependencies (FOSS, maintained; pin versions)

- `secure-systems-lab/go-securesystemslib/dsse`, `in-toto/in-toto-golang`, `sigstore/sigstore` (+ cosign, Rekor client), `spiffe/go-spiffe/v2`, `filippo.io/age` (or KMS) for encryption. Introduce one new pluggable signing-key provider; reuse complyctl's existing OCI libs. **Engine deferred** — no Benthos/Bento dependency taken.

### Licensing note (verified 2026-06-23)

`redpanda-data/benthos` (engine) = **MIT**; `redpanda-data/connect` (connectors) = **Redpanda Community License** (source-available), split into `public/bundle/free/v4` (Apache-2.0) and `public/bundle/enterprise/v4`. If adopted, restrict to MIT engine + Apache free bundle and verify each connector's tier; `warpstreamlabs/bento` is an all-FOSS fork pending a maintenance review. Not incurred now (engine deferred).

## What this discards / supersedes

The OTLP-push + OIDC `Authenticator` design is **superseded** (no shared secret, no `resolveOIDCToken`, no collector-side mTLS/static-JWKS verifier). Producer-identity work survives as the **signing** identity. **OTLP-forward becomes a sink** downstream of verification — OTel users are not abandoned.

## Open decisions (owed before/at implementation)

1. **Trust & key distribution — RESOLVED** (see §Trust & key distribution): self-sovereign, pluggable `BundleSource` over a signed public bundle, manual/gitops enrollment baseline, fail-closed-with-grace. Residual knobs: default grace windows per profile; which `BundleSource` + `EnrollmentProof` backends ship first.
1. **Transparency-log sink integration** — confirm the mandatory-for-compliance log sink against complytime-core's anchoring API; air-gap Sigstore vs RFC 3161 TSA.
1. **Consumer component home** — new repo vs complyctl binary vs complytime-adjacent.
1. **Retention window defaults** and CAS-vs-topic retention coupling per profile.
1. **Freeze `…/evidence/v1`** predicate fields (incl. the envelope `nonce`, per-item ciphertext digest, encryption, and batch metadata).
1. **Confidentiality vs traffic-analysis** — is metadata leakage in scope (padding)?

## References

- `cmd/complyctl/cli/scan.go`, `internal/complytime/config.go` — discarded path.
- New pluggable signing-key provider — the producer signing identity.
- `docs/dev/deployment-profiles-and-constraints.md` — transparency-log / witness / FIPS context this design now explicitly feeds.
- `2026-06-23-complyctl-pluggable-collector-auth-design.md` — superseded predecessor.

</details>

<details>
<summary>Appendix B: deployment profiles and cross-cutting constraints (governing principle, touchpoints, signing-key and FIPS ladder, identity per profile, anti-equivocation; cross-repo issue-planning section omitted)</summary>

# Deployment Profiles & Cross-Cutting Constraints (evidence pipeline + ComplyTime Core)

> **Shared design reference** for this evidence pipeline's signing/anchoring edge and the ComplyTime Core witness work. Grounded against source on 2026-06-17 and re-grounded against `complytime-core@main` on 2026-06-24 (noted inline where Core has since moved).

## The governing principle

**Air-gapped is the floor; cloud is a relaxation.** Every external touchpoint sits behind an interface with ≥2 implementations, and the **offline implementation is built and tested first**. If a thing works air-gapped, it works in cloud; the reverse is never guaranteed, so a cloud-convenient dependency must never become load-bearing in the core path.

This is a runtime **profile selector**, not a code fork: one codebase, a `profile: airgap | hybrid | cloud` configuration that wires implementations.

### The load-bearing invariant (make it a test)

**The `airgap` profile MUST have zero runtime egress.** Add a conformance test: run the full pipeline (sign → emit Gemara → anchor → witness → checkpoint export → verify) under the `airgap` profile with **all network egress blocked**, and assert it completes. This test is what prevents a dependency from smuggling an external call into the hot path. Treat any egress under `airgap` as a release blocker.

## Touchpoint matrix

| Concern           | `airgap` (floor)                                                                                                          | `cloud` (relaxation)                              | Interface / seam                                                                               |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| Signing key       | PKCS#11 (HSM / TPM-via-p11-kit / SoftHSM) → file-based fallback                                                           | Cloud KMS / Managed HSM; public Sigstore optional | signing-key provider (new)                                                                     |
| Ingest identity   | **in-enclave OIDC IdP** (issuer + JWKS resolve internally), no Core change; mTLS/static-JWKS need an upstream Core change | Cloud OIDC w/ JWKS discovery                      | `anchor.Authenticator` (client) + Core JWT verifier (discovery-only today)                     |
| Anti-equivocation | in-enclave multi-host witnesses **+ cosigned-checkpoint export diode (REQUIRED)**                                         | public witness network / Omniwitness              | [Tessera `WithWitnesses`](https://pkg.go.dev/github.com/transparency-dev/tessera) + export job |
| Trust roots       | provisioned via sealed image/media                                                                                        | fetched from known endpoints                      | config (key material in)                                                                       |
| Time              | internal trusted time source (sealed RTC / internal NTP)                                                                  | public NTP                                        | platform                                                                                       |
| Telemetry         | in-enclave OTLP collector only                                                                                            | cloud collector                                   | OTLP endpoint config                                                                           |
| Build / deps      | vendored, in-enclave mirror, **pin-by-digest**                                                                            | public module proxy / registry                    | go.mod + build pipeline                                                                        |
| Updates / CVE     | in-enclave update pipeline                                                                                                | normal patching                                   | ops                                                                                            |

## Signing key assurance ladder (decision: support all tiers)

One signing-key provider interface (this design introduces it; file + ephemeral-for-tests as the software tiers). Add a **PKCS#11 provider**, PKCS#11 is the universal interface that covers a real HSM, a **TPM via p11-kit**, and **SoftHSM2 for tests**, so a single provider serves hardware, platform, and CI custody.

| Tier              | Provider                         | Key custody                    | FIPS posture                                                                         |
| ----------------- | -------------------------------- | ------------------------------ | ------------------------------------------------------------------------------------ |
| test only         | ephemeral (in-memory)            | none                           | n/a                                                                                  |
| fallback (no HSM) | file-based (PEM)                 | software                       | **algorithms** FIPS-approved if built against a FIPS module; **no** hardware custody |
| platform          | PKCS#11 → TPM via p11-kit        | hardware-backed (platform TPM) | partial, TPM, not CMVP HSM                                                           |
| recommended       | PKCS#11 → FIPS 140-2/3 L2–L3 HSM | hardware                       | algorithms **and** custody                                                           |
| cloud profile     | Cloud KMS / Managed HSM          | managed                        | managed FIPS HSM                                                                     |

**FIPS nuance to document, not hide:** keeping a file-based provider does **not** break FIPS *algorithm* compliance, a FIPS Go build (BoringCrypto / Go 1.24+ native FIPS module) performing ECDSA P-256/384/521 is algorithm-compliant wherever the key lives. What the file provider gives up is *key-custody assurance*, not algorithm validation. So "won't supply an HSM" is a **custody-level downgrade, not a FIPS cliff.** This design's signing edge restricts attestation keys to ECDSA P-256/384/521 and rejects anything else; avoid Ed25519 for attestation keys (only approvable under FIPS 186-5 and module-dependent).

### FIPS is one knob, but it must hold all the way down

FIPS *is* a single config switch. The catch is that the switch's guarantee is **transitive**: when it is on, *every* signing operation in the chain must be FIPS-correct, or the system must refuse to run. A FIPS knob that reads "on" while an Ed25519 signature hides downstream is **false assurance, worse than off**.

Where the chain stands today:

- **Attestation keys, covered.** This design's signing edge restricts to ECDSA P-256/384/521 and rejects anything else.
- **Log checkpoint signer, not covered.** Core signs checkpoints with [`golang.org/x/mod/sumdb/note`](https://pkg.go.dev/golang.org/x/mod/sumdb/note), which is **Ed25519-only by design**, and it sits in Core, below this design's FIPS switch.
- **Witness cosignatures, not covered.** Core's witnesses use `cosignature/v1` (the [`transparency-dev/formats@v0.1.0`](https://pkg.go.dev/github.com/transparency-dev/formats/note) format), **also Ed25519-only**. The FIPS-capable successor is `cosignature/v2` (ML-DSA-44 / FIPS 204), not yet wired.

So honoring "FIPS on" end-to-end needs two things, both tracked in the Core witness work: (1) make the transparency-layer crypto FIPS-capable — a `note` / cosignature swap to a FIPS-validated algorithm (ML-DSA-44 / FIPS 204 is the target, not in the pinned libs today); and (2) wire it under the *same* switch and **fail closed**, with FIPS on, refuse to anchor against a non-FIPS checkpoint signer rather than emit a non-compliant proof. (Ed25519 is approvable only under FIPS 186-5 and is module-dependent, hence the gap.) Core now **persists** the checkpoint signer across restarts, so the remaining gap is purely the algorithm, not key stability.

## Identity per profile (decision: air-gap-native)

**Verified (2026-06-17):** Core `/api/ingest` **hard-requires a Bearer JWT** ([`handlers_ingest.go`](https://github.com/complytime-labs/complytime-core/blob/0d5496bd5e046339a383dea0ea57ae19da22016a/internal/store/handlers_ingest.go) → `extractBearerToken`, 401 if absent). The verifier ([`internal/auth/jwt.go`](https://github.com/complytime-labs/complytime-core/blob/0d5496bd5e046339a383dea0ea57ae19da22016a/internal/auth/jwt.go)) is **JWKS-discovery-only**, it fetches `<issuer>/.well-known/jwks.json` and has **no** static-key path, and there is **no mTLS client-cert auth** anywhere in the gateway. JWT algorithms are restricted to `ES256/384/512` + `RS256/384/512` (FIPS-clean; no Ed25519).

- **`airgap`, no Core change:** run an **in-enclave OIDC IdP** whose issuer URL resolves *inside* the enclave and serves `/.well-known/jwks.json` internally; Core's discovery cache fetches it with zero external egress. Lowest-friction air-gap path; needs no upstream change.
- **`airgap`, IdP-free (mTLS client-cert or static pre-provisioned JWKS):** an **upstream Core change**, neither exists today. Track in the Core issue.
- **`cloud`:** Cloud OIDC with JWKS discovery, as originally specced.

The `anchor.Authenticator` seam selects how the client *obtains/presents* the token, but Core's server side is JWT-only regardless. Net: **"OIDC fine for adopter #1" still holds even air-gapped, via an in-enclave IdP, and is in fact *less* work than mTLS, which would require changing Core.**

## Anti-equivocation under air-gap (decision: insider is in scope = hard requirement)

> **Status (re-grounded 2026-06-24):** as of `complytime-core@main`, Core now wires Tessera witnesses (`WithWitnesses`, gated on a Sigsum-format witness policy), persists its checkpoint signer across restarts (`internal/tessera/signerkey.go`), and serves checkpoints/inclusion proofs. The remaining target-state items are the cross-domain cosigned-checkpoint **export diode** and a **FIPS-capable** checkpoint/cosignature algorithm; the witness itself is the external SHA-pinned omniwitness, with `transparency-dev/formats@v0.1.0` providing the cosignature format.

A witness's power comes from being in a trust domain the log operator does not control. A single air-gapped enclave cannot provide that. Therefore:

- **External-attacker equivocation is PREVENTED** by in-enclave multi-host witnesses under separation of duties (different hosts, keys, HSMs, admin teams): an attacker must now subvert N independent hosts, not one.
- **Operator/insider equivocation is DETECTED, not prevented**, via a **required one-way export of cosigned checkpoints** (signed tree heads only, tiny, no sensitive content) across the air-gap boundary (data diode) to an external trust domain that reconciles them. A split view becomes *provable* on the next export.
- **The export cadence is a security parameter:** it bounds the maximum time a split view can persist before detection. Shorter cadence ⇒ tighter bound ⇒ more transfers. This dial is what an accreditor evaluates.
- A **separate-custody hardware witness** (sealed, different administering authority) inside the enclave raises the insider's cost toward partial prevention, but a physical-access insider can still attack it. Use it to raise the bar; rely on the diode for detection.

### Guarantee statement (use verbatim; do not oversell)

> Tamper-evident and append-only always. Equivocation is **prevented** against external attackers (in-enclave multi-host witnesses). Equivocation is **detectable within one export interval** against a malicious operator (cross-domain cosigned-checkpoint export). Real-time *prevention* of insider equivocation is not achievable inside a single air gap.

</details>
