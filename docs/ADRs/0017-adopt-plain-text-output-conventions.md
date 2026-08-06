# Adopt plain-text output conventions (`NO_COLOR` and emoji suppression)

**Status:** accepted

**Date:** 2026-07-17

**Deciders:** @trevor-vaughan @marcusburghardt @jpower432

## Context

A [proposed feature](https://github.com/complytime/complyctl/pull/744) for `complyctl` added UTF-8 emoji to command
output. That raised a question: how should the tools behave for users who need plain-text output? Screen reader users,
automation pipelines, terminals without emoji fonts, anyone piping output through `grep`.

Two concerns are in play:

1. ANSI color codes make output harder to process in pipes, logs, and screen readers.
2. Emoji characters get read aloud verbatim by screen readers (e.g., "party popper, check mark, cross mark"),
   render as broken glyphs on terminals without emoji fonts, and break text search.

Without a decision, each tool and PR will handle these on its own, and users get inconsistent behavior across the
project.

### Decision Drivers

- Accessibility for screen reader users (WCAG 1.1.1 non-text content).
- Reliable scripting and automation: output must survive pipes, `grep`, and log aggregation.
- Adherence to established conventions rather than inventing project-specific mechanisms.
- Simplicity: one environment variable, one behavior. Users who want plain output should not need to learn
  multiple knobs.

### Considered Options

- Option A: Adopt `NO_COLOR` for color and emoji suppression
- Option B: Adopt `NO_COLOR` for color only, rely on `TERM=dumb` for emoji suppression
- Option C: Adopt `NO_COLOR` with value-based extension (`NO_COLOR=full` for color + emoji)
- Option D: Adopt `NO_COLOR` for color, add a separate `--no-emoji` flag

## Decision

Chosen option: **Option A: `NO_COLOR` for color and emoji suppression**, because users who set `NO_COLOR`
are asking for plain, undecorated output. Splitting color and emoji into separate controls adds complexity
without a clear benefit. The [VMware Tanzu apps-cli-plugin](https://github.com/vmware-tanzu/apps-cli-plugin/issues/301)
reached the same conclusion and shipped it without issue.

When `NO_COLOR` is set (any non-empty value):

- Strip all ANSI color escape sequences from output.
- Replace all emoji characters with text equivalents (e.g., checkmark emoji becomes `[OK]`).

When `TERM=dumb` is set:

- Everything `NO_COLOR` does (color and emoji suppression).
- Suppress spinners, progress bars, and other cursor-movement decorations, since a dumb terminal lacks the
  capabilities those features require.

`NO_COLOR` and `TERM=dumb` are checked independently. Either one triggers color and emoji stripping.
`TERM=dumb` additionally suppresses animations and cursor-dependent output.

## Consequences

- **Positive:** One variable, one behavior. Users who set `NO_COLOR` globally get clean, greppable,
  screen-reader-friendly output from all ComplyTime tools. No extra flags to discover.
- **Negative:** The `NO_COLOR` convention technically covers only ANSI color. Bundling emoji suppression
  goes beyond its defined scope. In practice, this matches user expectations and has prior art (Tanzu), but
  it is a deviation from the letter of the convention.

### Pros and Cons of the Options

#### Option A: `NO_COLOR` for color and emoji suppression

- Good, because it is a single, well-known environment variable. One knob for "give me plain output."
- Good, because it matches what users actually expect when they set `NO_COLOR`, as demonstrated by
  [prior art](https://github.com/vmware-tanzu/apps-cli-plugin/issues/301).
- Good, because `TERM=dumb` is still honored for full plain-text including animation suppression.
- Bad, because it goes beyond the `NO_COLOR` convention's defined scope (ANSI color only) but does reflect intuitive
  understanding.

#### Option B: `NO_COLOR` for color only, `TERM=dumb` for emoji suppression

- Good, because it follows the `NO_COLOR` convention exactly. No scope overload.
- Good, because `TERM=dumb` is an established convention for "strip all decoration" per the
  [Command Line Interface Guidelines](https://clig.dev/).
- Bad, because `TERM=dumb` disables other terminal capabilities (cursor movement, screen clearing) for every
  program in the session, not just ComplyTime tools. Users who only want plain output from `complyctl`
  must accept side effects elsewhere.

#### Option C: `NO_COLOR` with value-based extension (`NO_COLOR=full`)

- Good, because any `NO_COLOR` value still strips color, staying compatible with the convention.
- Good, because `NO_COLOR=full` gives users color + emoji suppression without the session-wide side effects
  of `TERM=dumb`.
- Bad, because the `=full` value is a project-specific extension that other tools will not recognize.
- Bad, because it introduces complexity (two tiers of behavior from one variable) for a distinction most
  users will not need.

#### Option D: `NO_COLOR` for color, separate `--no-emoji` flag

- Good, because each concern has its own explicit control with clear semantics.
- Bad, because `--no-emoji` is not an established convention. It is project-specific, and users must discover
  it per tool.
- Bad, because there is no environment-variable equivalent. Users cannot set it globally across all
  ComplyTime tools.

## More Information

- [NO_COLOR convention](https://no-color.org/). The informal convention for disabling ANSI color output.
- [Command Line Interface Guidelines](https://clig.dev/). Community guidelines covering color, `TERM=dumb`,
  and output conventions.
- [W3C WCAG Technique H86](https://www.w3.org/WAI/WCAG20/Techniques/html/H86). Text alternatives for emoji
  (WCAG 1.1.1 non-text content).
- [FORCE_COLOR](https://force-color.org/). The inverse convention; users can override `NO_COLOR` per-tool.
- [VMware Tanzu apps-cli-plugin #301](https://github.com/vmware-tanzu/apps-cli-plugin/issues/301). Prior art:
  bundled emoji suppression into `--no-color` and shipped without issue.
- [Emoji Accessibility (CUNY OER)](https://guides.cuny.edu/accessibility/memeEmoji). Screen reader behavior
  with emoji content.
- [complyctl PR #744](https://github.com/complytime/complyctl/pull/744). The feature that prompted this decision.
