# Breaking Change Policy

complyctl is a core project in the complytime ecosystem. Its public API surface — the `Provider` interface in `pkg/provider` and the generated gRPC code in `api/plugin` — is imported by every provider implementation. A breaking change to this surface forces a `v2.0.0` release, which changes the Go import path for all consumers and triggers a cascade of updates across complytime-providers, downstream tooling, and any external integrations. This policy defines the rules that keep the API stable within v1.x.

## Rules for v1.x API Stability

These rules apply to all packages importable by external consumers: `pkg/provider` and `api/plugin`.

### Never add methods to the `Provider` interface

Adding a method to an interface breaks every existing implementation. Any type that previously satisfied `Provider` would fail to compile after the addition. Introduce new RPCs via separate optional interfaces (e.g., `Validator`) that providers can adopt incrementally. complyctl already uses this pattern.

### Never remove or change the type of struct fields in public types

Removing a field or changing its type breaks consumers that reference the field by name or rely on its type. Only add new fields. Additive changes are backward-compatible because existing code ignores fields it does not reference.

### Never add parameters to existing exported functions

Adding a parameter changes the function signature and breaks every call site. Use functional options (a variadic `...Option` parameter) to extend behavior without modifying the signature. This pattern lets callers opt in to new behavior without changing existing calls.

### Never remove, rename, or reorder enum values

Removing or renaming an enum value breaks consumers that match on it. Reordering changes the iota-assigned integer values, which breaks serialization and any stored data that references the old value. Only append new values at the end.

### Never change protobuf field numbers or field types in existing messages

Protobuf field numbers are the wire identity of each field. Changing a number or type breaks serialization compatibility with every compiled binary that uses the old definition. Assign new field numbers for new fields. Treat assigned numbers as immutable.

## Existing Safeguards

complyctl's architecture already includes several safeguards against accidental breaking changes.

| Safeguard | What it protects | How |
|---|---|---|
| Optional interface pattern | Provider interface stability | New RPCs are added as separate interfaces (e.g., `Validator`), not as methods on `Provider`. Existing providers compile without changes. |
| Domain type insulation | Public Go types from protobuf churn | `pkg/provider` defines its own domain types and converts to/from protobuf. Proto field additions do not change the public Go API. |
| `internal/` package boundary | Implementation details from external import | Go's compiler enforces that `internal/` packages cannot be imported outside the module. Most of complyctl's code lives here. |
| Buf FILE-level breaking change detection | Protobuf wire compatibility | `buf breaking` runs against the latest release and catches field number changes, field type changes, and RPC removals in `.proto` files. |
| Frozen handshake values in `plugin.go` | gRPC plugin protocol compatibility | The magic cookie and handshake UUID are frozen. Changing them would prevent existing compiled providers from connecting. |

## Three-Layer CI Enforcement

Each category of breaking change maps to a CI detection layer. No single tool catches everything — the three layers complement each other.

| Category | Detection Layer | Tool |
|---|---|---|
| Interface method additions | Layer 2: Go API | `apidiff` |
| Struct field removal / type changes | Layer 2: Go API | `apidiff` |
| Function signature changes | Layer 2: Go API | `apidiff` |
| Exported type / constant removal | Layer 2: Go API | `apidiff` |
| Proto field number / type changes | Layer 1: Protobuf | `buf breaking` |
| Proto RPC removal | Layer 1: Protobuf | `buf breaking` |
| Enum value removal / reordering | Layer 2: Go API | `apidiff` |
| Wire serialization incompatibilities | Layer 3: Integration | Custom CI job |

**Layer 1 (Protobuf)** detects schema-level changes that break the wire protocol. This layer already exists in complyctl's CI.

**Layer 2 (Go API)** detects changes to exported types, functions, and interfaces in `pkg/provider` and `api/plugin`. It compares the current branch against the latest release tag using `apidiff` (`golang.org/x/exp/cmd/apidiff`).

**Layer 3 (Integration)** detects runtime incompatibilities that pass static checks. It builds provider binaries from the latest complytime-providers release and runs them against the current complyctl branch, verifying that the gRPC plugin protocol remains functional across versions.

## Maintainer Override

When all three layers agree a change is breaking but the change is intentional, maintainers may override the CI failure. This applies to situations such as coordinated migrations with no known external consumers.

An override requires:

- Approval from at least two maintainers.
- Documented justification in the PR description or a dedicated comment explaining why the breaking change is acceptable and what migration path (if any) exists for consumers.

Overrides are the exception, not the norm. The CI layers exist to catch accidental breakage — they should block by default.

## References

- [ADR-0004: gRPC Provider Plugin Architecture](../ADRs/0004-grpc-provider-plugin-architecture.md) — provider binary model, go-plugin protocol, frozen handshake values
- [ADR-0006: ComplyPack as Content Envelope](../ADRs/0006-complypack-content-envelope.md) — content-only OCI envelope, evaluator-id routing, provider/content separation
