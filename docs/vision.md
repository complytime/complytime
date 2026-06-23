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

## How

ComplyTime addresses these problems through architectural choices at different stages of maturity.

**Decided and operational:**

1. **Separate what from how** — Compliance content (what must be true) and assessment logic (how to verify it) are independent streams with independent lifecycles and authorship. Neither is locked to a single tool or vendor. ([ADR-0005](ADRs/0005-two-stream-content-model.md))

2. **Route generically, evaluate specifically** — A single runtime client discovers and orchestrates evaluators without understanding their internal logic. Evaluator plugins handle the mapping between their native results and the common evidence model. ([ADR-0004](ADRs/0004-grpc-provider-plugin-architecture.md))

3. **Distribute as OCI artifacts** — Compliance content and assessment logic are packaged as ComplyPacks and distributed through standard OCI registries. The same infrastructure that moves container images moves compliance content. ([ADR-0003](ADRs/0003-oci-artifact-distribution.md), [ADR-0006](ADRs/0006-complypack-content-envelope.md))

**Under active exploration:**

4. **Structure evidence end-to-end** — Every assessment should produce structured, queryable evidence traceable to the requirement it addresses. Evidence should be a first-class artifact with provenance, not a byproduct buried in tool-specific logs. ([Evidence problem doc](problems/evidence.md))

5. **Compose compliance postures explicitly** — Policy layering, overrides, and exceptions should be expressed as composable operations with deterministic resolution. The effective policy should be derivable, never implicit. ([Requirement Fidelity problem doc](problems/requirement-fidelity.md))

### Evaluator Migration

Decisions 1 and 2 produce a capability worth naming explicitly: existing compliance content communities can migrate between evaluation engines without rewriting their content or losing assessment history.

Because the requirement identity is decoupled from the evaluator, ComplyTime supports **multi-evaluator mapping** — a legacy evaluator and a modern evaluator can run side-by-side against the same requirement, verifying output parity before cutover. Migration becomes gradual and verifiable rather than a risky rewrite. The requirement's identity, evidence chain, and compliance posture remain continuous regardless of which evaluator is replaced underneath. See [Evaluator Coupling](problems/evaluator-coupling.md) for the problem this addresses.

## Principles

Properties any solution in this space should exhibit:

- **Fidelity** — Machine-evaluable requirements must faithfully represent the source authority. Precision cannot come at the cost of accuracy.
- **Decoupling** — Requirements (what must be true) and verification logic (how to check it) have independent lifecycles and independent authorship. Neither should be locked to a single tool.
- **Composability** — Compliance postures compose from reusable components. Layering, overrides, and conflict resolution are explicit and traceable.
- **Evidence as a first-class artifact** — Assessment results, collected data, and audit artifacts are structured, queryable, and traceable to the requirement they address.
- **Provenance** — Every link in the chain — authority document to requirement to check to evidence to finding — is auditable.

