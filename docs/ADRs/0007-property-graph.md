# ADR-0007: Property Graph for Compliance Relationships

**Status:** proposed

**Date:** 2026-06-22

**Deciders:** complytime-dev team

## Context

Compliance relationships (requirements to requirements, requirements to evidence, systems to controls, controls to risk) are graph problems. A control maps to multiple requirements across frameworks. A piece of evidence satisfies controls in three standards simultaneously. A system change ripples through controls, evidence, and authorization decisions.

Spreadsheet-based mapping is the current norm. It does not scale, cannot be queried programmatically, and loses the reasoning behind relationships. Relational schemas can model some relationships, but traversal queries (impact analysis: "if this control changes, what authorizations are affected?"; traceability chains: "show me the full path from authority document to evidence"; attack path enumeration: "what is the shortest path from this exposure to sensitive data?") are unnatural in SQL, degrade as join depth increases, and become intractable to "desire path" driven change.

Compliance data is also inherently distributed. Different organizations, departments, and programs own different parts of the picture. A single centralized graph cannot hold all of it — sovereignty rules prohibit sending classified or proprietary controls across ownership boundaries, and governance breaks down when a global graph must represent conflicting but jurisdictionally valid interpretations of the same framework. The graph architecture must accommodate federation from the start: sharing generic, framework-level knowledge while structurally preventing leakage of proprietary requirements, classified controls, and evidence.

The domain needs a data store where relationships are first-class, richly attributed (typed, scored, temporally bounded, provenanced), and efficiently traversable at arbitrary depth — and where sovereignty boundaries are a structural property of the schema, not a policy overlay.

See the [Cross-Framework Mapping](../problems/cross-framework-mapping.md), [Evidence](../problems/evidence.md), and [Federated Compliance Graphs](../problems/federated-compliance-graphs.md) problem docs for the domain exploration behind this decision.

## Decisions

Adopt a property graph queried via [GQL](https://www.iso.org/standard/76120.html) (ISO/IEC 39075:2024) as the relationship layer for compliance data. GQL is the ISO standard for property graph queries — a declarative, pattern-matching language that drew from openCypher, PGQL, and G-CORE. The companion standard [SQL/PGQ](https://www.iso.org/standard/84803.html) (ISO 9075-16) embeds graph pattern matching inside SQL; GQL is a superset.

The property graph runs inside [PostgreSQL](https://www.postgresql.org/) using [Apache AGE](https://age.apache.org/), alongside relational data (PostgreSQL native) and vector similarity ([pgvector](https://github.com/pgvector/pgvector)). One database engine serves three data patterns. AGE currently implements openCypher, not GQL. Because openCypher is a semantic predecessor of GQL and the industry is converging (Neo4j's Cypher 10 tracks GQL alignment), queries written against AGE today should migrate to GQL with bounded effort as engine support matures. The target standard is GQL; openCypher via AGE is the pragmatic starting point.

The graph is the system of record for relationships. Projections to external systems (analytics engines, visualization tools, SIEMs, reporting platforms) are optimization, not authority. A projected copy is a cache; the graph is truth.

The graph schema must distinguish between generic framework-level knowledge (how public standards relate to each other) and proprietary or classified data (organization-specific requirements, evidence, system models). Generic mappings are exportable and federatable; proprietary nodes and edges are not serializable across sovereignty boundaries. This distinction is enforced at the schema level — edge and node types carry sovereignty metadata that determines what can cross a boundary — not by runtime access control alone.

The internal graph schema is driven by query needs and actual usage patterns, not by any external modeling language. Schema evolution is organic: new node types, edge types, and properties are added as real queries demand them, with periodic grooming to maintain coherence. Formal modeling languages ([SysML v2](https://www.omg.org/spec/SysML/2.0/), [CALM](https://github.com/finos/architecture-as-code), RDF) may be used as export/analysis formats but do not dictate the internal schema.

RDF is deferred while the data model is experimental. The property graph model (labeled nodes and edges with arbitrary key-value properties) is more natural for a schema that is still forming. If the schema stabilizes and the domain would benefit from formal ontology reasoning, migration to RDF remains possible.

## Consequences

- GQL (ISO/IEC 39075:2024) is the target query language for traversal across requirements, evidence, systems, and risk. As an ISO standard, it provides stronger long-term stability than any single-vendor language.
- openCypher via Apache AGE is the initial implementation. AGE does not yet support GQL. This creates a migration obligation: queries written against openCypher today must be ported when AGE (or a replacement engine) adds GQL support. The migration risk is bounded because openCypher is a semantic subset of GQL and the syntactic differences are incremental.
- Single-engine deployment (PostgreSQL + AGE + pgvector) reduces operational complexity compared to running separate graph, relational, and vector databases.
- Apache AGE is a top-level Apache Software Foundation project with an active community and PostgreSQL-native integration.
- Schema-level sovereignty metadata enables federation without centralization. Each deployment owns its graph; generic framework mappings can be shared across boundaries while proprietary data stays local. The federation design is explored in the [Federated Compliance Graphs](../problems/federated-compliance-graphs.md) problem doc; specific federation protocols are deferred pending implementation experience.
- Schema evolution is organic; no upfront modeling exercise required. This is a deliberate trade-off: faster iteration at the cost of potential inconsistency, mitigated by grooming cycles.
- RDF deferral means no formal ontology reasoning until the schema stabilizes. Acceptable for an experimental data model; revisit if the domain demands it.
- GQL is a graph-specific query language. Teams familiar with SQL but not graph queries face a learning curve. AGE allows mixing Cypher and SQL in the same query, easing transition. SQL/PGQ provides an additional on-ramp for teams that prefer to embed graph patterns inside SQL SELECT statements.
- Vendor optionality: GQL is an ISO standard supported by multiple engines. Migrating from PostgreSQL+AGE to another GQL-compatible engine is feasible if AGE hits scale limits or GQL adoption accelerates elsewhere.
