---
status: proposed
date: 2026-07-31
deciders: TBD
---

# Adopt story-level scoping for agentic development workflows

## Context and Problem Statement

Our current development process enforces task-level PR scoping: each change maps to one issue, one branch, one PR, with strict "no unrelated changes" rules. This process was designed to reduce human reviewer cognitive load. Each PR is small, single-concern, and easy to hold in a reviewer's head.

With the shift to agentic development (LLM-driven coding sessions and LLM-assisted review), this granularity creates a behavioral economics problem. When an agent encounters a bug or improvement opportunity while working on a related task, the process forces it to file a ticket and move on. A future session then pays the full cold-start cost to re-derive the codebase understanding the first session already held. The context that makes the fix nearly free is discarded by policy.

If we make no change, we continue to pay redundant context-derivation costs on every deferred fix, and agents spend tokens writing tickets for work they could have completed in-situ at near-zero marginal cost.

## Decision Drivers

- **Context economics:** Input tokens account for ~85% of agentic session cost, at a 25:1 input-to-output ratio ([Vantage][vantage]). Each new session re-reads the codebase, re-derives architecture, and re-loads conversation history. A session that runs 2x as many turns costs 3-4x as much due to quadratic context accumulation. Cold starts are the dominant cost, not generation.
- **Marginal cost of in-situ fixes:** An agent that already holds the relevant context can apply a small fix for near-zero additional tokens. A new session for the same fix pays the full cold-start overhead again.
- **Token variance on re-derivation:** Runs on identical tasks can differ by 30x in token consumption ([arXiv 2604.22750][token-spend]). Re-deriving context is unpredictably expensive.
- **Review model shift:** When reviewers are also LLMs, mixed-concern diffs do not cause the cognitive fatigue that motivated single-concern PRs. An agent reviewer evaluates each logical change independently within the same diff.
- **Constraint decay risk:** Agent performance degrades as structural requirements accumulate in context ([arXiv 2605.06445][constraint-decay]). Unbounded scope increases the chance of constraint violations.

## Considered Options

- **Task-level scoping (status quo):** One issue, one branch, one PR. No unrelated changes permitted.
- **Story-level scoping with commit boundaries:** One story, one branch, one session. Multiple logical commits within the branch. Opportunistic fixes allowed if noted in the PR description. Review at story level.
- **Unbounded session scoping:** No scope limits. Agent works until done or context is exhausted.

## Decision Outcome

Chosen option: **Story-level scoping with commit boundaries**, because it eliminates redundant cold-start costs while preserving reviewability through commit-level structure and avoiding the constraint decay risk of unbounded sessions.

### Consequences

- **Positive:** Agents fix issues in-situ instead of filing tickets, eliminating the deferred cold-start cost. Story-level sessions amortize the ~85% input-token overhead across all related tasks. Review remains structured because reviewers can walk commits sequentially. Aligns with the ["scout rule"](https://lawsofsoftwareengineering.com/laws/boy-scout-rule/) (leave the codebase cleaner than you found it) rather than conflicting with it.
- **Negative:** PRs are larger and require reviewer discipline to evaluate commit-by-commit rather than as a single diff. Opportunistic fixes may occasionally introduce unrelated regressions. Blame granularity is coarser at the PR level (though commit-level blame remains intact). Teams accustomed to task-level PRs will need process adjustment.

## Pros and Cons of the Options

### Task-level scoping (status quo)

- Good, because each PR is small and single-concern, easy for human reviewers
- Good, because blame and bisect map cleanly to one logical change
- Bad, because agents discard context and file tickets instead of fixing in-situ, paying full cold-start cost on every deferred fix
- Bad, because the "no unrelated changes" rule conflicts with the boy scout rule, preventing cleanup that the agent is already positioned to do
- Bad, because verification and refinement (the dominant agentic cost per [MSR 2026][tokenomics]) must be repeated from scratch in each new session

### Story-level scoping with commit boundaries

- Good, because cold-start context overhead is paid once per story instead of once per task
- Good, because in-situ fixes happen at near-zero marginal cost while context is hot
- Good, because commit boundaries preserve blame granularity and allow commit-by-commit review
- Good, because LLM reviewers are not degraded by mixed-concern diffs
- Bad, because larger PRs require reviewers (human or agent) to spend more time per review
- Bad, because requires discipline to enforce the architectural change boundary (see guardrail below)

### Unbounded session scoping

- Good, because maximum context reuse, zero artificial boundaries
- Bad, because constraint decay degrades agent performance as context grows
- Bad, because quadratic context accumulation makes late-session turns disproportionately expensive
- Bad, because review becomes impractical for both humans and agents at mega-diff scale
- Bad, because a single session failure risks losing all uncommitted work

## More Information

### Commit hygiene: rebase before merge

Story-level sessions produce messy commit histories: checkpoint saves, fix-the-fix commits, debug artifacts. During the session this is fine; context is hot and the agent knows what each commit was for. But the git log outlives the session. When something breaks six months from now, a human will be running `git bisect` and `git blame` against this history. If the commits are noise, the audit trail is worthless at the moment it matters most.

End-of-session commit cleanup (interactive rebase to collapse, reorder, and re-message) is cheap: the agent still holds full context, and the operation is mechanical. The cost of _not_ doing it is paid later by a human with no context, under incident pressure, trying to understand why a commit called "wip fix 3" touched four files.

Rules:

- Collapse fix-up and checkpoint commits into their logical parent.
- Each surviving commit should represent one reviewable idea: a feature addition, a bug fix, a migration, a test expansion.
- Commit messages follow conventional commits and explain _why_, not _what_. The diff says what.
- Opportunistic fixes noted in the PR description should also be identifiable as distinct commits so they can be reverted independently if they cause regressions.

This is not about aesthetics. It is about making `git bisect` useful and `git blame` legible when a human needs to untangle the work without the session context that produced it.

### Guardrail: architectural changes are the hard boundary

Scope creep is bounded by a single rule: **no architectural changes within an opportunistic fix.** An architectural change is any modification that alters the system's shape rather than its behavior within the existing shape:

- Introducing a new abstraction or interface
- Changing an existing interface contract
- Adding or removing a dependency
- Restructuring module boundaries or ownership
- Altering data models or storage schemas
- Changing deployment topology

An agent that spots a bug, dead code, a missing guard, or a style inconsistency can fix it in-situ. These are behavioral changes within the existing architecture. The moment the fix requires changing how components relate to each other, it is an architectural decision that deserves its own story and review, regardless of how cheap it would be to do right now.

This is process discipline and a technical necessity. The [constraint decay research][constraint-decay] shows that agent performance degrades as structural requirements accumulate in context. Architectural changes add exactly the kind of structural complexity that triggers this decay. Deferring them to a fresh session with a focused scope produces better outcomes from the agent itself.

### Supporting research

| Reference                                          | Key finding                                                                                                                                                       |
| -------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [How Do AI Agents Spend Your Money?][token-spend]  | 30x token variance across identical tasks; accuracy peaks at intermediate cost                                                                                    |
| [Tokenomics (MSR 2026)][tokenomics]                | Refinement/verification accounts for 59.4% of tokens, not initial generation                                                                                      |
| [Agentic Coding Cost Analysis][vantage]            | ~85% of session cost is input tokens; 25:1 input-to-output ratio; quadratic cost growth with turn count                                                           |
| [Constraint Decay in LLM Agents][constraint-decay] | 40% relative performance loss as structural constraints accumulate (L0 to L3)                                                                                     |
| [Less Context, Better Agents][less-context]        | Selective context retention outperforms full-history retention in enterprise tool-use workflows; applicability to coding agents is inferred, not directly studied |
| [Efficient Context Management][jetbrains-ctx]      | Agent-generated context becomes noise; managed context cuts cost 50%+ without degrading solve rate                                                                |

[token-spend]: https://arxiv.org/abs/2604.22750
[tokenomics]: https://arxiv.org/html/2601.14470v1
[vantage]: https://www.vantage.sh/blog/agentic-coding-costs
[constraint-decay]: https://arxiv.org/html/2605.06445v1
[less-context]: https://arxiv.org/html/2606.10209v1
[jetbrains-ctx]: https://blog.jetbrains.com/research/2025/12/efficient-context-management/
