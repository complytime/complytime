# Federated Compliance Graphs

Compliance data is inherently distributed. Different organizations, departments, and programs own different parts of the picture, and a single centralized graph cannot hold all of it. Sovereignty rules prohibit sending classified or proprietary controls across ownership boundaries. Governance breaks down when a global graph must represent conflicting but jurisdictionally valid interpretations of the same framework. Scale limits emerge when a single store spans organizations that trust each other only conditionally.

At the same time, the value of compliance mapping increases with network effects. Organizations that map the same frameworks independently waste collective effort. The architectural problem: share generic, framework-level knowledge (how public standards relate to each other) while structurally preventing leakage of proprietary requirements, classified controls, and evidence. This is a federation problem, not an access control problem.

## Current approaches and prior art

Teams share mapping spreadsheets ad-hoc between auditors, consultants, and compliance staff. No structure, no versioning, no provenance, no way to reconcile conflicting versions. A mapping spreadsheet that works for one organization's interpretation of a framework may be subtly wrong for another's. There is no trace of who originated the mapping, under what interpretation, or when it was last validated. Spreadsheets scale to tens of mappings but collapse under hundreds.

Commercial GRC tools provide tenant isolation but within a single vendor's platform. Data sovereignty depends entirely on vendor trust. The vendor hosts your proprietary requirements alongside everyone else's. The platform promises logical separation, but the data lives in the vendor's infrastructure, subject to the vendor's jurisdiction, subpoena risk, and operational policies. No federation between instances. If two organizations both use the same GRC platform but want to share generic framework mappings without exposing proprietary data, they cannot. The platform was not designed for peer-to-peer federation.

The Open Security Architecture (OSA) project maintains a repository of security architecture patterns and control mappings. Publicly available but manually curated, no machine-readable API, and the project maintains that human curation is superior to automation. OSA provides valuable reference mappings but cannot scale to the full combinatorial space of framework-to-framework relationships. It covers the most common pairings (NIST to ISO, CIS to NIST) but does not address how an organization extends those mappings with proprietary controls or how to federate proprietary extensions without centralizing them.

NIST maintains the National Cybersecurity Online Informative References ([OLIR](https://csrc.nist.gov/projects/olir)), a mapping repository between frameworks. Authoritative for the mappings it covers but limited in scope, manually maintained, and not designed for bidirectional community contribution. OLIR is read-only for most users. Organizations can consume the mappings but cannot contribute corrections or extensions back into the commons without going through NIST's publication process. This works for stable, slowly-evolving frameworks but does not support the continuous feedback loop that a living compliance mapping ecosystem requires.

Academic and industry work on federating knowledge graphs (SPARQL federation, linked data, semantic web standards) is technically mature but not applied to compliance specifically. Relevant patterns: query federation (route queries across boundaries without copying data), data virtualization (expose a unified view without centralizing storage), and selective replication (sync chosen subsets of a graph between federated nodes). The challenge is adapting these patterns to compliance's unique sovereignty and trust requirements. Not all data is equal, and not all boundaries are symmetrical.

## Proposed approaches

Explore a federated meta-graph architecture where sovereignty is a structural property, not a policy overlay:

Each area of responsibility owns its own property graph, its own security boundary, and its own rules about what is visible and to whom. A classified government program maintains an air-gapped graph with zero external connectivity. A regulated financial institution maintains a private graph behind firewall and VPN. An open-source community maintains a public graph accessible to all. The same architecture serves all three scenarios without special-casing any of them.

A federation layer knits these graphs into a navigable meta-graph. The meta-graph carries only generic, framework-level knowledge: how [NIST 800-53](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final) maps to [ISO 27001](https://www.iso.org/standard/27001), how [CIS Controls](https://www.cisecurity.org/controls) relate to SOC 2, how [PCI-DSS](https://www.pcisecuritystandards.org/) overlaps with GDPR. Proprietary requirements, classified controls, and evidence never cross the owner's boundary. The boundary is enforced by the graph schema itself: edges that represent generic framework mappings are exportable; edges that represent proprietary implementations are not. The architecture makes it impossible to accidentally leak proprietary data because the schema does not permit those relationships to traverse federation boundaries.

The same architecture serves diverse deployment scenarios without forking the design:

- Air-gapped deployment: fully isolated, zero external dependencies. The graph contains classified controls and proprietary system models. The deployment syncs with the commons only if and when its owner chooses, via one-way import of generic framework mappings on removable media or during scheduled sync windows. The air-gapped node can contribute mappings back to the commons through the same mechanism, but only after declassification review.

- Regulated organization: controlled connectivity, selective sharing. The graph contains proprietary risk models and audit evidence that must stay within the organization's boundary. The deployment federates with the commons for generic framework mappings but exports nothing proprietary. The organization benefits from community-vetted mappings and contributes corrections back to the commons when it discovers errors, but its internal requirements never leave its boundary.

- Open community: full participation in the commons. The graph contains only public framework definitions and community-contributed mappings. All nodes are readable by anyone. The community reviews, corrects, flags, and votes on mappings. Every interaction is labeled data that improves the mapping engine over time.

The community commons model: organizations publish the generic mappings they choose to share. The community reviews, corrects, flags, and votes on them. An arbiter (domain expert, standards body, or delegated authority) blesses accepted mappings. The commons does not replace authoritative sources like NIST OLIR; it complements them by covering the long tail of framework pairings that no single authority can maintain. The commons is a feedback loop: organizations consume vetted mappings, discover errors or gaps in practice, and contribute corrections back.

Sovereignty is a structural property, not a policy. The architecture must enforce boundary rules by construction, not by trusting participants to obey access control lists. What crosses a boundary is always a generic mapping (how two public frameworks relate), never a proprietary requirement or piece of evidence. The schema enforces this: proprietary nodes and edges are not serializable across federation boundaries. An attempt to export them fails at the type level, not at runtime after the data has already leaked.

## Open questions

Federation introduces questions about trust governance, sovereignty enforcement, sync mechanisms, and query performance that require architectural decisions. See [ADR-0009](../ADRs/0009-property-graph.md) for the property graph decision and [ADR-0012](../ADRs/0012-federated-graph-sovereignty.md) for the federation governance decisions that have crystallized from this exploration.

## See also

- [Cross-Framework Mapping](cross-framework-mapping.md): the mapping problem federation distributes
- [Requirement Fidelity](requirement-fidelity.md): how meaning is preserved across boundaries
- [Evaluator Coupling](evaluator-coupling.md): how tool-specific logic affects evidence collection
- [Evidence](evidence.md): fragmented, manual, and opaque compliance evidence
