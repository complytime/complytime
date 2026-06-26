# System Modeling and Drift

Compliance assessment assumes you know what system you are assessing. In practice, the system under observation is a moving target. Infrastructure is provisioned dynamically, configurations drift from declared state, and the gap between what was deployed and what is running widens silently between assessment cycles.

Static analysis of infrastructure-as-code (Terraform plans, Kubernetes manifests, CloudFormation templates) captures declared intent, not runtime reality. A compliant-at-deploy cluster that accumulates ad-hoc changes between scans presents a false positive to any tool that reads only the declared state. Drift detection requires comparing a model of expected state against observed state, and that requires having a model in the first place.

The modeling problem has two distinct faces. First, the analytical model: a formal representation of the system's components, their relationships, data flows, and trust boundaries, suitable for reasoning about impact, risk, and attack paths. Second, the runtime inventory: a continuously updated record of what actually exists, discovered from the systems themselves. These serve different purposes but must be reconciled. The analytical model is where you reason, the runtime inventory is ground truth.

Operator burden is the binding constraint. Any system modeling approach that requires engineers to manually maintain a model will fail. The model will rot as soon as operational pressure arrives. Discovery must be automated, incremental, and near-zero friction.

## Current approaches and prior art

### Infrastructure-as-code scanning

Tools like [Checkov][checkov], tfsec, and kube-bench analyze declared configuration. They are effective for pre-deployment gates but blind to runtime drift. They operate on individual resource definitions and do not produce a system model (no relationships, no data flows, no trust boundaries). They answer "does this Terraform file violate this rule?" but not "what depends on this service?" or "what attack paths cross this boundary?"

### Cloud Security Posture Management (CSPM)

Commercial tools such as Prisma Cloud, Wiz, and AWS Security Hub continuously scan cloud environments and provide runtime state visibility. They detect misconfigurations, over-permissioned identities, and exposed resources. However, they are proprietary, cloud-specific, and do not produce a portable system model. The representation is internal and optimized for the vendor's analysis engine. Organizations cannot export the discovered topology to perform custom reasoning or integrate with other tools. CSPM is assessment, not modeling.

### [Cartography][cartography] (CNCF)

Cartography is open-source infrastructure discovery tooling that pulls asset and relationship data from cloud provider APIs and builds a [Neo4j][neo4j] graph. It is strong prior art for automated discovery: modular, extensible, and relationship-focused. It answers questions like "which EC2 instances can assume this IAM role?" or "which S3 buckets are exposed to the internet?" However, Cartography is tightly coupled to Neo4j. The graph is the runtime store, and the query language (Cypher) is the interface. Swapping the backend requires rewriting every adapter.

### [FINOS CALM][finos-calm] (Common Architecture Language Model)

CALM is an emerging open standard for describing system architectures as structured data. It is lighter-weight than formal modeling languages and designed for cross-organizational interchange. The project is early-stage but shows promise for a portable architecture representation that tools can consume. CALM is declarative: it captures what the architect intended, not what is running. It does not address runtime discovery or drift detection.

### [SysML v2][sysml-v2]

SysML (Systems Modeling Language) version 2 is an [OMG][omg] formal specification adopted in July 2025. It is an industry-standard systems modeling language with a formal semantics layer (KerML, the Kernel Modeling Language) and a Systems Modeling API for tool interoperability. SysML v2 is designed for analytical modeling of complex systems: the kind of formal analysis that safety-critical industries (aerospace, automotive, medical devices) require. It supports hierarchical decomposition, allocation of functions to components, requirement traceability, and parametric constraints. It is powerful but heavyweight. The learning curve is steep, the tooling ecosystem is nascent, and the model requires significant upfront investment.

SysML v2 is intended for the kind of rigorous, traceable analysis where you need to prove properties about a system's behavior under all conditions. Compliance assessment rarely needs that level of formality, but the subset that does (verifying cryptographic properties, reasoning about isolation boundaries, proving access control invariants) would benefit from the formal semantics.

### Custom inventory systems

Many organizations build internal Configuration Management Databases (CMDBs) using platforms like ServiceNow, Ralph, or Nautobot. These track assets (servers, applications, contracts, owners) and support operational workflows. However, they typically lack the relationship modeling needed for impact analysis or attack path enumeration. A CMDB knows that Server A exists and belongs to Team B. It does not know that Server A can reach Database C because NetworkPolicy D allows it. CMDBs are catalogs, not models.

## Proposed approaches

The challenge is reconciling two competing needs: automated discovery with near-zero operator burden, and formal modeling for rigorous analysis. These are not the same artifact.

### Two-layer approach

A separation of concerns emerges naturally:

1. **Runtime discovery layer**: automated, continuous, low-burden. Pulls system state from APIs and agents: cloud provider APIs (AWS, Azure, GCP), Kubernetes API server, network traffic analyzers, [OpenTelemetry][opentelemetry] service discovery, configuration management agents. Produces a runtime inventory in a property graph store. The discovery layer should be modular: different discovery adapters for different infrastructure types, extensible to new environments without rewriting the core. Cartography is strong prior art here and we could work with the team to optimize it and make it graph-backend agnostic. The runtime inventory is organic. Its schema is shaped by what actually exists and what queries are actually needed, not by a predefined formal model.

2. **Analytical modeling layer**: formal system model for reasoning about impact, risk, and compliance relationships. SysML v2 is the most complete option for this layer. Its formal semantics enable automated reasoning that lighter approaches cannot support: proving that no data flow crosses a regulatory boundary, verifying that all paths to a sensitive resource require authentication, calculating the blast radius of a compromised component. However, the analytical model need not be the runtime store. SysML v2 (or CALM, or any modeling language) could be a projection from the property graph: an output format for analysis, not the internal representation. This avoids coupling the compliance platform to any single modeling standard while still enabling formal analysis when needed.

The key design question: does the compliance graph need a formal modeling language as its internal schema, or should the internal schema be organic (shaped by actual queries and usage) with formal models as an export or analysis layer? The latter avoids coupling the platform to any single modeling standard while still enabling formal analysis when needed. The cost is translation overhead. Every formal analysis tool needs an adapter that translates from the internal graph to its expected input format.

### Drift detection as a delta operation

Drift detection emerges naturally from the two-layer approach. The analytical model represents expected state: the architecture as designed, possibly blessed by a compliance authority. The runtime inventory represents observed state: what actually exists. Drift is the set-theoretic difference between them.

When a component appears in observed state but not in expected state, that is unplanned infrastructure. When a relationship appears in expected state but not in observed state, that is missing enforcement. When a component's attributes differ between the two, that is configuration drift. Each delta is a drift event that triggers compliance reassessment of affected controls.

This requires a stable identifier scheme that links analytical model entities to runtime inventory entities. An API Gateway declared in a Terraform module and discovered from AWS APIs must map to the same logical component. Without stable identifiers, drift detection devolves into string matching heuristics that fail as soon as naming conventions change.

### Minimal viable discovery

The bootstrapping problem is steep. Full system discovery across every infrastructure type is a multi-year effort. What is the minimum viable discovery that imposes near-zero operator burden and delivers enough value to justify adoption?

Cloud provider APIs are the lowest-friction starting point. They require no agents, no network access to internal systems, and no changes to deployment workflows. Read-only API access can discover compute instances, storage buckets, databases, identity resources, network topology, and access policies. The coverage is broad enough to support meaningful compliance checks. The blind spots (on-premises systems, third-party SaaS, physical infrastructure) can be addressed later.

Kubernetes API discovery is similarly low-friction for containerized workloads. The API server is the source of truth for cluster state. No agents required. The API model is well-documented and stable. Discovery adapters for Kubernetes are well-trodden ground.

Network traffic observation is higher friction. It requires agents, access to network flows (NetFlow, sFlow, packet captures), and sustained compute for analysis. The payoff is runtime behavioral data: what actually talks to what rather than only what is permitted to talk. This is essential for behavioral evaluation (see [Evaluator Coupling](evaluator-coupling.md)) but not necessary for initial discovery.

## Open questions

System modeling for compliance raises questions about schema representation, discovery scope, query architecture, drift semantics, ephemeral infrastructure, and abstraction depth that require architectural decisions. See [ADR-0008](../ADRs/0008-evaluator-interface-contract.md) for the evaluator interface decisions and the broader system modeling ADRs as they crystallize.

See [Evaluator Coupling](evaluator-coupling.md) for how verification logic is tied to specific infrastructure types, creating the fragmentation that a unified system model could address.

[checkov]: https://www.checkov.io/
[cartography]: https://github.com/lyft/cartography
[neo4j]: https://neo4j.com/
[finos-calm]: https://github.com/finos/architecture-as-code
[sysml-v2]: https://www.omg.org/spec/SysML/2.0/
[omg]: https://www.omg.org/
[opentelemetry]: https://opentelemetry.io/
