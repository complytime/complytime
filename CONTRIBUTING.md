# Contributing to ComplyTime

## Contribution Model

ComplyTime uses an escalation model for contributions. Start lightweight; formalize only when scope demands it.

### Escalation Ladder

```
Issue → Problem Doc → ADR → Plan → Implementation (in separate repos)
```

| Level       | What it is                        | When to use                                   | Approval |
|:------------|:----------------------------------|:----------------------------------------------|:---------|
| Issue       | Question, observation, suggestion | Always the starting point when unsure         | None     |
| Problem doc | Deep dive into a technical domain | New problem area or significant expansion     | PR merge |
| ADR         | Crystallized yes-or-no decision   | Specific architectural choice needs recording | PR merge |
| Plan        | Implementation breakdown          | Accepted design ready for execution           | PR merge |

### What goes where

- **`docs/problems/`** — Open-ended explorations. No formal approval gate beyond PR review. Multiple approaches can coexist. Anyone can expand.
- **`docs/ADRs/`** — Short, decisive records. State what was decided, why, and what the consequences are. Use the [ADR template](docs/ADRs/adr-template.md).
- **`docs/plans/`** — Phased implementation breakdowns for accepted designs. Reference ADRs and problem docs.
- **`docs/guides/`** — Practical how-to documentation for operators and developers.

### Problem Doc Guidelines

Problem docs explore a technical domain. They are living documents — expect them to evolve.

Structure (suggested, not mandatory):
1. Why this is hard
2. Current approaches / prior art
3. Proposed approaches with trade-offs
4. Open questions

Do not worry about "finishing" a problem doc. Partial explorations with open questions are valuable.

### ADR Guidelines

ADRs record decisions that have been made. They are immutable once accepted — if a decision is reversed, write a new ADR that supersedes the old one.

Use the template at `docs/ADRs/adr-template.md`. Keep them short. One page maximum.

## Process

1. Fork and create a branch
2. Make your changes
3. Open a PR with a descriptive title (`docs: expand evaluator-coupling problem doc`)
4. Address review feedback
5. Merge requires maintainer approval (one for problem docs, two for ADRs/plans)

## Implementation Code

Implementation lives in separate repositories:
- [complyctl](https://github.com/complytime/complyctl) — CLI runtime
- [complytime-providers](https://github.com/complytime/complytime-providers) — scanning providers
- [complypack](https://github.com/complytime/complypack) — pack authoring and packaging
- [complytime-policies](https://github.com/complytime/complytime-policies) — published policy bundles

This repo is for design exploration. Implementation PRs reference design docs here.
