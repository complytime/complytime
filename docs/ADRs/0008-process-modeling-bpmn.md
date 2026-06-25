# ADR-0008: BPMN 2.0 for Process Modeling

**Status:** proposed

**Date:** 2026-06-22

**Deciders:** complytime-dev team

## Context

Many compliance requirements (access reviews, incident response procedures, change management, separation of duties, training programs) are satisfied by process, not by technical configuration. Technical assessment tools (scanners, policy engines, configuration auditors) cannot evaluate these controls. Without a machine-readable representation, procedural controls depend on manual attestation: interviews, document review, and ticket audits that do not scale to continuous assessment.

Several standards exist for modeling business processes in machine-readable form. [BPMN 2.0](https://www.omg.org/spec/BPMN/2.0/) (Business Process Model and Notation) is an [OMG](https://www.omg.org/) standard finalized in 2011 with broad industry adoption. [CMMN](https://www.omg.org/spec/CMMN/) (Case Management Model and Notation) handles less-structured case-based work. [DMN](https://www.omg.org/spec/DMN/) (Decision Model and Notation) models decision logic. All three are mature OMG standards.

## Decisions

Adopt BPMN 2.0 as the standard for modeling compliance-relevant business processes. Process models are first-class artifacts in the compliance graph: a BPMN model paired with execution records proves that a process is both defined and followed.

BPMN 2.0 is chosen for three reasons:
1. **Maturity and stability.** 15 years of industry use, extensive tooling, well-understood semantics. Low risk of breaking changes.
2. **Red Hat ecosystem alignment.** Red Hat has a deep investment in BPMN tooling. [jBPM](https://kie.apache.org/components/jbpm/) is the upstream community project that Red Hat productized as Red Hat Process Automation Manager (RHPAM), bundling jBPM with the [Drools](https://kie.apache.org/components/drools/) rule engine and DMN support. Red Hat has since shifted toward [Kogito](https://kogito.kie.org/), a cloud-native process automation runtime built on [Quarkus](https://quarkus.io/) for Kubernetes and OpenShift deployment. Kogito decomposes BPMN and DMN models into individual microservices with built-in audit trails — the execution evidence model that compliance assessment requires. The full ecosystem (jBPM, Drools, Kogito, [SonataFlow](https://kie.apache.org/components/sonataflow/)) is consolidated under [Apache KIE (Incubating)](https://kie.apache.org/), providing the tightest integration across BPMN, CMMN, and DMN from a single upstream project. [Camunda](https://camunda.com/), [Flowable](https://www.flowable.com/open-source), and other engines are also widely adopted and produce BPMN natively.
3. **LLM compatibility.** Current-generation language models can translate written procedures (runbooks, standard operating procedures, policy documents) into BPMN XML with reasonable fidelity, subject to human review. This provides a scalable ingestion path for organizations with large bodies of undocumented or informally documented processes.

CMMN and DMN are complementary standards that may be adopted in future ADRs as the process modeling domain matures. They are not in scope for this decision.

## Consequences

- Procedural and operational controls become assessable alongside technical controls in a unified compliance graph.
- Existing BPMN tooling ecosystems (process automation platforms, modeling tools) are reusable without format conversion.
- Process evidence has a structured, queryable form. A BPMN model is data, not an opaque document attachment.
- Organizations already using BPMN-compatible tools get native integration. The Red Hat/Apache KIE ecosystem (Kogito on OpenShift) provides the most direct path from process modeling to cloud-native execution with built-in audit trails.
- The standard is mature and stable, with low risk of breaking changes or ecosystem fragmentation.
- BPMN models need a defined ingestion path into the compliance graph. Mapping BPMN elements to graph nodes and edges is implementation work.
- Process execution verification (proving a process was followed, not merely documented) remains an open problem. BPMN defines the process; execution evidence must come from the systems that run it (workflow engines, ticket systems, attestation records).
- BPMN is XML-based and verbose. Tooling exists to manage this, but direct authoring without a modeling tool is impractical.
