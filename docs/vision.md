# Vision

## The Goal

Compliance assessment becomes automated, continuous, and trustworthy because the manual approach does not scale to the pace of modern infrastructure, the breadth of regulatory surface area, and the number of teams that need to participate.

## Why

Compliance practitioners know their domains well. The work is rigorous and the results are real. The problem is throughput. As systems multiply, regulations expand, and teams distribute, the manual processes that work at human scale stop keeping pace.

1. **Requirements lose fidelity at scale** — Translating authority documents (NIST 800-53, CIS Benchmarks, organizational policies) into executable checks is skilled work. It works when one team owns one standard. It breaks down when dozens of teams interpret the same requirement independently, with no shared traceability back to the source.

2. **Evaluation is fragmented — and incomplete** — Different technology stacks require different assessment tools. OpenSCAP handles VMs. Kyverno handles Kubernetes. Cloud platforms have their own compliance tooling. Each brings its own policy language, data model, and runtime. The assessment intent is often the same across stacks, but the expression is completely different — and nothing ties them together. Beyond that, procedural, physical, and operational controls fall outside the reach of automated tooling entirely. Those still depend on manual attestation, interviews, and documentation review — processes that work but add another dimension that scales independently of technical assessment.

3. **Compliance postures are layered but tooling is flat** — Real compliance stacks industry standards, regulatory requirements, organizational baselines, and team-level exceptions. Each layer adds, removes, or overrides requirements. That composition is manageable when one person holds it in their head. It stops working when the number of layers, teams, and exceptions outgrows any single person's context.

4. **Evidence is fragmented across tools and processes** — Every assessment produces proof, but it lands in different tools, different formats, and different workflows. Collecting, normalizing, and tracing evidence back to the requirement it satisfies is feasible for a single audit cycle. It becomes continuous toil when the cadence shifts from annual to continuous.

5. **Trust depends on institutional knowledge** — Assessment content runs, results appear, and experienced practitioners know what to trust. But that trust lives in people's heads — nothing structural connects the authority document to the check to the evidence to the finding. That works until the people rotate, the team scales, or an external auditor asks for proof.

## Principles

Properties any solution in this space should exhibit:

- **Fidelity** — Machine-evaluable requirements must faithfully represent the source authority. Precision cannot come at the cost of accuracy.
- **Decoupling** — Requirements (what must be true) and verification logic (how to check it) have independent lifecycles and independent authorship. Neither should be locked to a single tool.
- **Composability** — Compliance postures compose from reusable components. Layering, overrides, and conflict resolution are explicit and traceable.
- **Evidence as a first-class artifact** — Assessment results, collected data, and audit artifacts are structured, queryable, and traceable to the requirement they address.
- **Provenance** — Every link in the chain — authority document to requirement to check to evidence to finding — is auditable.

