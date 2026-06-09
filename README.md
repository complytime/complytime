# ComplyTime

A living design document for automated compliance assessment in cloud-native systems.

## What is this?

This repo is a living design document exploring how compliance assessment becomes automated, continuous, and trustworthy at scale.

## What's here

- **[docs/vision.md](docs/vision.md)** — What automated compliance assessment should become and why
- **[docs/architecture.md](docs/architecture.md)** — Component vocabulary and relationships
- **[docs/glossary.md](docs/glossary.md)** — Canonical definitions for project-specific terms
- **[docs/problems/](docs/problems/)** — Deep dives into each technical problem domain:
  - [Requirement Fidelity](docs/problems/requirement-fidelity.md) — Preserving meaning and rationale as requirements cross functional boundaries
  - [Evaluator Coupling](docs/problems/evaluator-coupling.md) — Verification logic tied to specific tools and runtimes
  - [Evidence](docs/problems/evidence.md) — Fragmented, manual, and opaque compliance evidence
- **[docs/plans/](docs/plans/)** — Implementation plans for accepted designs
- **[docs/ADRs/](docs/ADRs/)** — Architecture Decision Records for crystallized decisions
- **[docs/guides/](docs/guides/)** — Practical how-to documentation

## How to contribute

Pick a problem area that interests you. Read the existing document. Add your perspective, propose solutions, poke holes in existing proposals. Open a PR.

### Where does my contribution go?

| If you have...                                    | Then...                                                        |
|:--------------------------------------------------|:---------------------------------------------------------------|
| A question, bug, or small suggestion              | **File an issue** — lowest friction, can graduate later.       |
| A new problem area no existing doc covers         | **Create a problem doc** in `docs/problems/` and link it here. |
| More to say about an existing problem area        | **Expand the existing problem doc.**                           |
| A specific decision that needs a yes-or-no answer | **Propose an ADR** in `docs/ADRs/` via a pull request.         |
| An accepted design ready for execution            | **Write a plan** in `docs/plans/`.                             |

When in doubt, start with an issue.

## Browse online

This documentation is browsable at [complytime.github.io/complytime](https://complytime.github.io/complytime) (powered by docsify).

## Acknowledgments

This repo's structure — living design document, problem-driven exploration, ADRs — is adopted from [fullsend](https://github.com/fullsend-ai/fullsend).

## License

This project is licensed under the [Apache License, Version 2.0](LICENSE).
