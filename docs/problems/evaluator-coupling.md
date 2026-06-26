# Evaluator Coupling

Verification logic is tied to specific tools and runtimes. Scanners specialize in a domain, couple data fetching with evaluation, discard raw inputs after assessment, and produce output in tool-specific formats. The compliance requirement and the evaluation logic end up fused into the same artifact. If governance posture is embedded in any single tool, it is locked to that tool's surface. No one tool covers the full assessment surface.

This doc explores *why* that coupling exists and what forms it takes. It does not propose what ComplyTime should build.

## The structural pattern

Every compliance scanner follows the same pattern:

1. **Specialize in a domain.** OpenSCAP handles OS/VM configuration. Kyverno handles Kubernetes admission. Prowler handles cloud APIs. Checkov handles infrastructure-as-code. Each is good at what it does.
2. **Couple fetching with evaluation.** The tool collects system state and evaluates it in one pass. Prowler calls AWS APIs and emits findings. InSpec SSHes into a host, collects attributes, and emits findings. The raw data — the system state actually observed — is consumed internally and not exposed as a reusable artifact.
3. **Discard raw inputs.** After evaluation, the observed state is gone or buried inside tool-specific output. Another tool cannot re-evaluate the same data. Re-running the scan is the only way to reproduce the evidence — and the system state may have changed.
4. **Produce incompatible output.** Each tool emits its own finding format. Prowler JSON does not look like Checkov CLI output does not look like Kyverno admission logs. Aggregating results across tools is manual reconstruction.

The [ComplianceAsCode](https://github.com/ComplianceAsCode/content) community illustrates this at scale. Over a decade, the project built one of the most mature compliance content ecosystems in open source — thousands of XCCDF/OVAL checks maintained across multiple platforms and regulatory frameworks. SCAP is excellent at configuration assessment for operational systems. But OVAL fuses the requirement and the system probe into one artifact: the *what* and the *how* are the same thing. When the assessment surface expanded to Kubernetes and cloud platforms, the content corpus — the community's most valuable asset — could not follow. The intent behind each check cannot be extracted and re-evaluated with a different tool without rewriting it.

This is not unique to SCAP. It is the pattern. Every tool's content corpus is bound to its evaluation model.

## Why this is hard

### Different stacks require different tools

No single evaluation engine spans the full assessment surface:

| Domain | Typical tools | Policy language |
|:---|:---|:---|
| OS/VM configuration | OpenSCAP, InSpec | OVAL/XCCDF, Ruby DSL |
| Kubernetes | Kyverno, OPA Gatekeeper | YAML + CEL, Rego |
| Cloud infrastructure | Prowler, Checkov, Cloud Custodian | Python, HCL-aware |
| Application code | Semgrep, CodeQL | Pattern DSL, QL |

A single compliance requirement — "encrypt data at rest" — needs different checks on each surface. Those checks diverge. Different teams write them independently in different languages against different data models. Nobody maintains a shared mapping between them. The same intent gets expressed N times with no traceability back to a common source.

### Evaluation spans multiple forms

The [Gemara model](https://gemara.openssf.org/model/) distinguishes two evaluation forms at Layer 5:

- **Intent evaluation** examines static resources — code, configuration, infrastructure plans — to determine whether a system is *prepared* in alignment with policy.
- **Behavioral evaluation** observes running systems to determine whether they *behave* as expected under conditions that static analysis cannot predict.

A single requirement often spans both. "Network segmentation between production and staging" needs an intent check (does the NetworkPolicy specify the correct rules?) and a behavioral check (can a staging pod actually reach a production database at runtime?). Most engines optimize for one form. A requirement that spans both needs multiple evaluators — and the coupling multiplies.

### Not all evaluation is technical

Procedural, physical, and operational controls require evaluation methods that automated tooling does not support. This is structural, not a gap waiting to be filled by better scanners.

| Control type | Example | Evaluation method |
|:---|:---|:---|
| Procedural | Incident response plan reviewed annually | Document review, attestation |
| Physical | Server rooms require badge access | Physical inspection, access log review |
| Operational | Backup restoration tested quarterly | Process verification, test records |
| Technical | Containers do not run as root | Configuration scan, admission control |

An organization's compliance posture spans all four types. Automated results (JSON findings) and manual results (a signed attestation) must be unified into a single assessment view. They come from fundamentally different evaluators with different output formats, different freshness semantics, and different trust properties.

### Lock-in compounds over time

Organizations accumulate checks. A small team starts with 20 policies. Three years later: 400 policies, custom helper functions, test fixtures, CI integrations, and a runbook only one engineer understands. The switching cost grows linearly with the corpus. The *risk* of switching grows faster — each policy was validated against a specific engine version and data model. The tool becomes infrastructure. Removing it is a project, not a PR.

## Open questions

- Is full evaluator independence achievable, or is some coupling to the engine's data model unavoidable?
- What is the right boundary between "platform" (routes content to evaluators) and "evaluator" (executes checks)?
- How do manual and automated evaluation results unify without forcing manual assessment into formats designed for machine evaluation?
- When a requirement spans intent and behavioral evaluation, who owns the coupling between the two evaluators — the compliance engineer, the operations team, or an integration layer?

## Resolution status

1. **Is full evaluator independence achievable?** Resolved. Engine diversity is accepted as permanent. Data coupling is broken by the unified property graph. Logic coupling to the engine's native language is accepted and managed by the platform through routing and result unification. See [ADR-0008](../ADRs/0008-evaluator-interface-contract.md).

2. **What is the right boundary between platform and evaluator?** Resolved. The platform sends scoped payloads (message-passing); evaluators return structured findings. Cross-entity graph access is opt-in with strict authorization. See [ADR-0008](../ADRs/0008-evaluator-interface-contract.md).

3. **How do manual and automated evaluation results unify?** Partially resolved. Both map to requirements in the property graph. The unified evidence model is under exploration. See [evidence.md](evidence.md) for the expanded problem exploration.

4. **Who owns the coupling between intent and behavioral evaluators?** Resolved. The requirement declares what evaluation forms are needed; the assessment logic declares how to compose results. The Runtime Client executes but does not define composition. See [ADR-0008](../ADRs/0008-evaluator-interface-contract.md).
