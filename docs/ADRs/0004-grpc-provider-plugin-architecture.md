# ADR-0004: gRPC Provider Plugin Architecture

**Status:** accepted

**Date:** 2026-02-01 (retroactive)

**Deciders:** complytime-dev team

## Context

complyctl needs to orchestrate multiple evaluation engines (OpenSCAP for filesystem scanning, AMPEL for API-based assessment, OPA for policy-as-code). 
Each engine has different runtime requirements, dependencies, and release cadences.

Options considered:
1. **Embedded libraries** — compile all engines into complyctl. Creates massive binary, CGo dependencies, and forces synchronized releases.
2. **Shell exec** — invoke engines as CLI commands. Fragile argument passing, no structured communication, error handling is parsing stdout.
3. **gRPC subprocess (go-plugin)** — each provider is a standalone binary communicating via gRPC over a local socket. HashiCorp's go-plugin provides the multiplexing, health checking, and lifecycle management.

## Decision

Providers are standalone gRPC binaries using HashiCorp go-plugin. The protocol defines four RPCs: `Describe`, `Generate`, `Scan`, and `Export`. Each provider is an independent binary with its own release cycle.

Provider discovery is by filesystem convention: executables named `complyctl-provider-*` in `~/.complytime/providers/` (user) with fallback to `/usr/libexec/complytime/providers/` (system). The runtime routes ComplyPacks to providers by evaluator-id.

Wire handshake values (`COMPLYCTL_PLUGIN` magic cookie, UUID) are frozen for backward compatibility.

## Consequences

- Providers release independently — no complyctl rebuild needed for provider updates.
- Each provider can have different dependencies (CGo, system libraries) without affecting complyctl.
- Adding a new evaluation engine requires only implementing the gRPC interface — no complyctl changes.
- Providers run with full host privileges (same as complyctl). No isolation between providers.
- Discovery by convention means no manifest files or configuration — but providers must follow the naming convention exactly.
- The go-plugin dependency is a de-facto standard (Terraform, Vault, Packer use the same pattern).

**Source:** [complyctl spec 001](https://github.com/complytime/complyctl/blob/main/specs/001-gemara-native-workflow/spec.md), [spec 004](https://github.com/complytime/complyctl/blob/main/specs/004-providers-repository-split/spec.md)
