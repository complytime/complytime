# ADR-0016: Adopt the XDG Base Directory Specification

**Status:** proposed

**Date:** 2026-07-13

**Deciders:** @trevor-vaughan @marcusburghardt @jpower432

## Context

Most `complytime` tools currently store all user-scoped data under a single `$HOME/.complytime` directory. This
convention was an explicit design choice (documented in the complyctl 001 spec, session 2026-02-27) citing alignment
with `.docker`, `.kube`, `.cargo` patterns and avoidance of cross-platform complexity.

Two developments motivate revisiting this decision:

1. **Convention divergence**: `complypack` adopted `$XDG_CACHE_HOME/complypack` in PR complytime/complypack#127. If
   `complyctl` retains `$HOME/.complytime`, complytime ships two incompatible directory conventions at its first stable
   release.

2. **Ecosystem shift**: The Go standard library provides `os.UserCacheDir()` and `os.UserConfigDir()` which implement
   XDG on Linux and native platform conventions on macOS and Windows. Newer Go CLI tools (`gh`, Helm v3) have adopted
   XDG. The older tools that still use dot-directories (`kubectl`, `docker`) predate Go's XDG support and are likely
   legacy holdovers.

The `$HOME/.complytime` directory currently mixes three XDG categories in a single location:

| Current path                                            | Content                                    | XDG category |
|---------------------------------------------------------|--------------------------------------------|--------------|
| `~/.complytime/policies/`, `~/.complytime/complypacks/` | Downloaded OCI artifacts (re-downloadable) | Cache        |
| `~/.complytime/providers/`, `~/.complytime/state.json`  | User-installed binaries, persistent state  | Data         |

## Decision

All `complytime` tools MUST follow the XDG Base Directory Specification for user-scoped directories.

The resolution order for each category is:

| Category | Env var                      | Linux default               | macOS default                              |
|----------|------------------------------|-----------------------------|--------------------------------------------|
| Cache    | `$XDG_CACHE_HOME/complytime` | `~/.cache/complytime`       | `~/Library/Caches/complytime`              |
| Data     | `$XDG_DATA_HOME/complytime`  | `~/.local/share/complytime` | `~/Library/Application Support/complytime` |

Implementation guidelines:

- **Use Go stdlib where available**: `os.UserCacheDir()` for cache, `os.UserConfigDir()` for config. These handle Linux, macOS, and Windows transparently.
- **Implement `XDG_DATA_HOME` manually**: Go has no `os.UserDataDir()`. Check `$XDG_DATA_HOME`; fall back to `$HOME/.local/share` on Linux and `~/Library/Application Support` on macOS.
- **Workspace-local `.complytime/` is unchanged**: The per-project `.complytime/` directory (analogous to `.git/`) is project-scoped, not user-scoped. It is not subject to XDG.
- **Tool-specific overrides**: Each tool MAY provide an env var override (e.g., `COMPLYCTL_CACHE_DIR`) that takes priority over XDG, following the pattern established by Helm (`HELM_CACHE_HOME`) and `gh` (`GH_CONFIG_DIR`).
- **Migration**: On first run, if `~/.complytime/` exists, print a one-time deprecation warning with migration instructions. Do not auto-migrate.

### Scope

Primary Linux support. macOS supported via Go stdlib platform detection. Windows is not a current target, but using `os.UserCacheDir()` and `os.UserConfigDir()` ensures correct behavior if Windows support is added later.

## Consequences

- **Breaking change**: A decision must be made for each application how it will handle the legacy directory. Backwards
  compatible support, or auto-migration, are two mitigating possibilities.
- **Repositories affected**: [complyctl](https://github.com/complytime/complyctl) (primary, ~15 source files),
  [complytime-providers](https://github.com/complytime/complytime-providers) (inherited via vendored constants and
  docs). [complypack](https://github.com/complytime/complypack) already complies. All other org repos are unaffected at
  the time of writing.
- **Supersedes**: The anti-XDG decision in the complyctl 001 spec (session 2026-02-27).
- **Convention alignment**: All complytime tools use the same directory standard. `complypack` cache at `~/.cache/complypack` and `complyctl` cache at `~/.cache/complytime` are natural neighbors.
- **Operator benefit**: Symlinking `~/.cache` to a large disk benefits all complytime tools (and all other XDG-compliant tools) automatically.

## Related

- complytime/complypack#127 — adopted XDG for complypack cache, triggered this ADR
- [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir/latest/)
- Helm v2 to v3 migration — prior art for XDG adoption at a major version boundary
