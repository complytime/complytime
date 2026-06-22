# ADR-0006: ComplyPack as Content Envelope (Not Bundled Plugin)

**Status:** accepted

**Date:** 2026-06-04 (retroactive — CEP-0001 merged)

**Deciders:** @jpower432
## Context

The original ComplyPack proposal (PR #3 OpenSpec) bundled both the evaluation plugin (Wasm binary) and assessment content (Rego/CUE) in a single OCI artifact. One pull gave you everything needed to evaluate.

Review identified problems:
1. **Authorship conflation** — who owns the plugin vs. who owns the assessment content? Different trust levels.
2. **Update coupling** — fixing a Rego rule requires republishing the entire Wasm plugin.
3. **Complexity** — Wasm runtime (wazero), TinyGo compilation, host function ABI — significant new infrastructure before any assessment logic ships.
4. **Existing patterns** — native gRPC providers already work. The community has providers in production.

## Decision

A ComplyPack is a **content-only** OCI distribution envelope. It contains assessment logic (e.g., OPA bundle) but not the evaluation runtime.

The evaluation runtime is a separate native provider binary (gRPC, go-plugin). complyctl routes ComplyPack content to the correct provider by `evaluator-id` metadata. The provider receives a filesystem path to the extracted content and knows how to consume its own format.

Wasm sandboxing is deferred to a future ADR, triggered when community trust requirements demand execution isolation.

## Consequences

- Ships immediately using existing provider infrastructure — no new runtime to build.
- Assessment content updates independently of provider binary updates.
- Provider authors use familiar Go tooling and full host access.
- No execution isolation between providers — a compromised provider has the same privileges as complyctl.
- The same Rego source can later compile to Wasm (`opa build -t wasm`) without rewriting — forward-compatible.
- Content delivery depends on filesystem access — providers read extracted files from disk. A future Wasm model would compile content into the binary, eliminating this coupling.

**Source:** [CEP-0001](https://github.com/complytime/complytime/blob/main/ceps/cep-0001-complypack-architecture.md), [PR #3 design evolution](https://github.com/complytime/complytime/pull/3)
