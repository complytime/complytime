# Glossary

Terms defined in the [Gemara Lexicon](https://gemara.openssf.org) — including assessment, assessment requirement, catalog, control, evaluation, enforcement, guidance, objective, policy, and risk — are used throughout this project with their canonical Gemara definitions. 
Terms related to governance maturity — including evidence, audit, and policy enforcement — align with the [Automated Governance Maturity Model](https://tag-security.cncf.io/community/resources/automated-governance-maturity-model/). 

This glossary covers terms specific to ComplyTime.

## Domain Terms

| Term | Definition |
|:---|:---|
| **Authority document** | A source standard, regulation, or organizational policy that defines compliance requirements (e.g., [NIST 800-53](https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final), [CIS Benchmarks](https://www.cisecurity.org/cis-benchmarks), PCI-DSS). |
| **Effective policy** | The resolved output of policy composition — all imports flattened, inheritance applied, overrides resolved. The concrete set of requirements an organization is assessed against. |
| **Evidence** | Raw proof that a requirement is met or not — configuration scan output, API responses, attestation documents, logs. The underlying data that supports a finding. See [Evidence problem doc](problems/evidence.md). |
| **Requirement fidelity** | The degree to which a machine-evaluable requirement faithfully represents the intent of its source authority document. See [Requirement Fidelity problem doc](problems/requirement-fidelity.md). |
| **Trust chain** | The end-to-end provenance linking an authority document to a requirement to a check to evidence to a finding. See [architecture](architecture.md#boundaries). |

## Implementation Terms

| Term | Definition |
|:---|:---|
| **ComplyPack** | Uniform OCI distribution envelope for packaged evaluation logic. Content is opaque to the runtime; only the evaluator understands it. See [ComplyPack Phase 0 plan](plans/complypack-phase0.md) and [complypack repo](https://github.com/complytime/complypack). |
| **ComplyPack ID** | Reverse-domain identifier (e.g., `io.complytime.opa.cis-k8s`) that uniquely identifies a ComplyPack. |
| **Evaluator ID** | Identifier (e.g., `opa`) that routes a ComplyPack to the correct provider. |
| **`#EvaluationLog`** | Merged assessment output produced from one or more `#AssessmentLog` entries. The universal output contract of a scan. |
| **`#AssessmentLog`** | Output of a single evaluator's assessment. Contains results for the controls/requirements it assessed. |
| **Gemara** | Schema framework for compliance content. CUE-based. See [gemara.openssf.org](https://gemara.openssf.org). |
| **Native provider** | A gRPC binary that handles data collection and evaluation. Runs with host privileges. See [complytime-providers](https://github.com/complytime/complytime-providers). |
| **OCI layout** | On-disk representation of OCI artifacts (per the [OCI Image Layout Specification](https://github.com/opencontainers/image-spec/blob/main/image-layout.md)). Used for caching and transport. |
| **Provider** | An execution unit that performs data collection and/or evaluation. Discovered by filesystem convention. See [ADR-0004](ADRs/0004-grpc-provider-plugin-architecture.md). |
| **Two-stream model** | Architectural separation between compliance content (what must be true) and assessment logic (how to verify it). Independent lifecycles, independent authorship. See [ADR-0005](ADRs/0005-two-stream-content-model.md). |
