# AGENTS.md

This repository is authored and maintained with LLM assistance.

## Transparency

Design documents, ADRs, problem docs, and plans in this repo are drafted, expanded, and revised with the help of AI coding agents. Human review and editorial judgment govern all merged content.

Commits include an `Assisted-by` trailer when agent-assisted.

## Agent Instructions

This is a **docs-only** repo. No application code lives here — implementation is in separate repositories (see [architecture.md](docs/architecture.md#repositories)).

When contributing to this repo:

- **Follow the escalation model** in [CONTRIBUTING.md](CONTRIBUTING.md): Issue → Problem Doc → ADR → Plan.
- **Write for humans.** Active voice, no filler, bottom-line-up-front. Tables over paragraphs where data is structured.
- **Problem docs explore; ADRs decide.** Do not mix exploration with decisions. If the answer is uncertain, it belongs in `docs/problems/`. If it's settled, write an ADR.
- **ADRs are immutable.** Never edit an accepted ADR. Supersede it with a new one.
- **Keep the sidebar current.** Update `docs/_sidebar.md` when adding new documents.
- **Link new docs in the README.** The "What's here" section in `README.md` is the entry point.
