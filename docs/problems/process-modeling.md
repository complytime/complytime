# Process Modeling

Many compliance requirements are satisfied by process, not configuration. Access reviews, incident response procedures, change management, training programs: these are procedural and operational controls. Technical assessment tools (scanners, policy engines, configuration auditors) focus on technical controls and cannot evaluate whether a process exists, is followed, or is effective.

Current approaches to procedural compliance rely on manual attestation: interviews, document review, ticket audits. These work for annual assessment cycles but do not scale to continuous assessment. The gap between technical and procedural controls creates two parallel compliance workflows with no shared data model, no shared evidence format, and no shared traceability.

The overlap problem makes this worse. The same procedural requirement appears in multiple frameworks under different identifiers and phrasing. Without a machine-readable process representation, there is no way to recognize or reuse that overlap across audit cycles.

## Current approaches and prior art

### Manual attestation and document review

Manual attestation is the status quo. Auditors interview teams, review runbooks, examine ticket histories, and check whether documented procedures exist and were followed. This approach works for annual audit cycles but has structural limits:

- Evidence is unstructured prose. Runbooks live in wikis, procedures in shared drives, execution records in email threads or ticket comments. No structured query layer exists. Finding whether a specific procedure was followed requires reading through dozens of tickets.
- No traceability. The link between the procedure, its execution, and the compliance requirement it satisfies is manual and implicit. An incident response runbook might satisfy requirements in three different frameworks, but nothing machine-readable records that connection.
- Temporal validity is unclear. A procedure executed six months ago may no longer be valid if the process changed. There is no standard way to version procedures or invalidate old evidence when a process changes.
- No continuous assessment. Manual review happens at audit time. Between audits, there is no visibility into whether procedures are still being followed or whether process drift has occurred.

This approach does not scale. It works for a small number of frameworks and a small number of procedural controls, but the cost grows linearly with both. An organization assessing against five frameworks with 200 procedural controls each faces 1,000 manual attestation points per audit cycle, with significant overlap that cannot be detected or reused.

### GRC platforms

Governance, Risk, and Compliance (GRC) platforms (ServiceNow GRC, RSA Archer, MetricStream, and others) provide workflow tooling for attestation. Teams document controls, assign ownership, track completion, and collect evidence through a structured interface. These platforms improve visibility and workflow coordination over fully manual approaches, but introduce new coupling:

- Proprietary silos. Each platform has its own data model, workflow engine, and evidence format. Migrating between platforms or integrating with external systems requires custom adapters or manual export/import cycles.
- Weak integration with technical assessment. GRC platforms excel at procedural workflow but rarely integrate with technical assessment pipelines. Technical scan results (Kubernetes admission logs, CSPM findings, SAST reports) live in separate systems with different evidence formats. Unifying technical and procedural evidence into a single compliance view requires manual aggregation.
- Platform-specific evidence format. Evidence collected in a GRC platform is structured for that platform's reporting engine, not for cross-platform consumption. Exporting evidence for use in another system often requires flattening it to PDF or CSV, discarding machine-readable structure.

GRC platforms solve the workflow coordination problem for procedural controls but do not solve the representation problem. They provide a container for attestation but not a standard format that technical systems can consume or that multiple platforms can exchange without loss of fidelity.

### [BPMN 2.0][bpmn] (Business Process Model and Notation)

BPMN 2.0 is an [OMG][] standard for business process modeling, published in 2011 and widely adopted in process automation tooling. It provides a graphical notation for process diagrams and an XML serialization format for machine-readable process definitions. BPMN was designed for process definition and execution in workflow engines, not specifically for compliance, but the representation is sufficient for compliance purposes.

BPMN models are composed of flow elements (tasks, events, gateways), sequence flows (connections between elements), and swimlanes (responsibility assignment). A process model can represent who performs each step, under what conditions, with what inputs and outputs, and in what order. This is the same information needed to assess whether a procedural control is defined and followed.

The BPMN ecosystem is mature. Red Hat has a significant investment here. [jBPM][] is the upstream community project that Red Hat productized as Red Hat Process Automation Manager (RHPAM), bundling jBPM with the [Drools][] rule engine and DMN support into an enterprise-supported BPMN 2.0 platform. RHPAM provides integrated authoring (Business Central), execution, and audit trail capabilities across BPMN, DMN, and CMMN from a single vendor stack. Red Hat has since shifted toward [Kogito][], a cloud-native process automation runtime built on [Quarkus][] and designed for Kubernetes and OpenShift deployment. Kogito decomposes BPMN and DMN models into individual microservices rather than deploying them into a monolithic engine, which aligns with how modern compliance infrastructure is typically deployed. The entire ecosystem (jBPM, Drools, Kogito, OptaPlanner, and the newer [SonataFlow][] serverless workflow component) has been consolidated under the [Apache KIE (Incubating)][kie] project, with Apache KIE 10.1.0 released in July 2025. This ecosystem provides the tightest integration across all three OMG standards (BPMN, CMMN, DMN) and the most direct path from process modeling to cloud-native execution with built-in audit trails.

In terms of other tooling, [Camunda][] is widely adopted in the broader market and has strong community tooling ([bpmn.io][], Camunda Modeler). [Flowable][] is another open-source alternative. Enterprise BPMS platforms from IBM, Oracle, and others fill out the rest of the market. LLMs can assist in translating written procedures into BPMN regardless of target engine, improving ingestion consistency for organizations with large runbook corpora in unstructured formats.

BPMN's primary limitation for compliance is that it models process definition, not execution verification. A BPMN model proves a process is documented, but not that it was followed. To assess compliance, you need both the definition (the BPMN model) and execution evidence (records of actual process runs). BPMN does not standardize the execution evidence format. Each workflow engine produces its own logs and audit trails.

### [CMMN][] (Case Management Model and Notation)

CMMN is an OMG standard for case-based work (investigations, audits, incident response, support tickets) where the sequence of steps is not fully predetermined. Unlike BPMN, which models structured workflows with explicit sequence flows, CMMN models semi-structured work where human judgment determines what happens next.

For compliance, CMMN is relevant for requirements that involve case-based investigation rather than repeatable procedures. An incident response process might begin with a BPMN workflow for initial triage, then transition to a CMMN case for investigation and remediation. A security audit might model initial scoping as BPMN and the actual audit work as CMMN.

CMMN adoption is lower than BPMN. Fewer workflow engines support it, and the tooling ecosystem is smaller. For compliance purposes, CMMN represents a secondary concern: most procedural controls are structured workflows that BPMN handles well. Case-based controls exist (investigation procedures, escalation processes), but they are less common than structured procedures.

### [DMN][] (Decision Model and Notation)

DMN is an OMG standard for decision logic that models decision tables and decision graphs: given inputs, what decision should be made? DMN is complementary to BPMN. A BPMN workflow might include a decision task; DMN models the logic inside that decision.

For compliance, DMN is useful for requirements that involve decision-based controls: approval chains (who must approve based on request amount and requester role?), risk classification (what risk tier applies based on system criticality and data sensitivity?), access control decisions (does the user have the required attributes to access this resource?).

DMN has reasonable tooling support. Drools has been a production rule engine for over two decades and provides the most mature open-source DMN runtime, with tight integration into jBPM's BPMN execution. Red Hat Process Automation Manager and its successor Kogito expose DMN alongside BPMN in a single deployment unit, so decision logic and process flow share a common execution context and audit trail. Camunda and other engines also support DMN decision tables natively. The standard is less mature than BPMN but sufficient for representing compliance-relevant decision logic.

## Proposed approaches

The strongest candidate for procedural compliance modeling is **BPMN 2.0**. It has the largest tooling ecosystem, the most mature standard, and the widest industry adoption. LLMs can assist in translating unstructured runbooks into BPMN, reducing ingestion friction for organizations with large existing procedure corpora. BPMN models become first-class evidence when paired with execution records: the model proves the process is defined, the execution log proves it was followed.

The key distinction: modeling a process vs. verifying execution. These are separate concerns that require separate evidence:

- Process modeling answers "Is the procedure documented?" A BPMN XML file, stored in version control, serves as evidence that the process is defined. The model itself is the artifact. Changes to the process produce new versions of the BPMN file, providing an audit trail of process evolution.
- Execution verification answers "Was the procedure followed?" This requires integration with workflow engines, ticket systems, or attestation records. A Kogito microservice deployment on OpenShift with completed process instances in its audit log is evidence of execution. A Camunda process instance with completed tasks is evidence of execution. A Jira ticket tagged with a runbook identifier and marked resolved is evidence of execution. An access review spreadsheet with sign-off timestamps is evidence of execution. The format varies, but the pattern is consistent: a record that the defined process actually ran, with timestamps, actors, inputs, outputs, and outcomes.

Both are needed for compliance. The BPMN model provides the definition. The execution record provides proof. The compliance graph node representing a procedural control links to both.

This creates a unified evidence model across technical and procedural controls. A compliance graph node for "Quarterly access review" points to:

- The requirement (control identifier, framework source, rationale)
- The process model (BPMN file defining review steps, approvers, escalation paths)
- Execution evidence (records of the last four quarterly reviews, with timestamps, reviewers, findings, remediation)
- Related technical controls (RBAC policies, identity provider configuration, access logs)

A node for "Encrypt data at rest" points to:

- The requirement (same structure)
- Technical evidence (CSPM scan showing encrypted volumes, admission policy blocking unencrypted PVCs)
- Related procedural controls (key rotation process, cryptographic policy review)

The graph structure is the same. The evidence types differ, but the traceability model is uniform.

### Composition with technical assessment

Process models compose with technical assessment through shared graph nodes. A compliance requirement node in the graph has outbound edges to both technical evidence and procedural evidence. The assessment engine evaluates both:

- Technical edge: Query the evidence store for the most recent CSPM scan. Does it show all volumes encrypted? If yes, technical evidence satisfied. If no, technical evidence failed.
- Procedural edge: Query the evidence store for the most recent execution record of the quarterly access review process. Does it exist? Is it within the required time window? Were all steps completed? If yes, procedural evidence satisfied. If no, procedural evidence failed.

The compliance posture for that requirement is the logical AND of both edges. A requirement that needs both technical and procedural controls is only satisfied when both are green.

This is the same query model that works for purely technical controls or purely procedural controls. The graph is the unifying abstraction. The evidence types are heterogeneous, but the traceability and query semantics are not.

## Open questions

Process modeling for compliance raises questions about fidelity expectations, evidence acceptance, ingestion architecture, and versioning that require architectural decisions. See [ADR-0010](../ADRs/0010-process-modeling-bpmn.md) for the decisions that have crystallized from this exploration.

## Cross-references

See [Evaluator Coupling](evaluator-coupling.md) for how the same verification gap manifests in technical controls, and [Evidence](evidence.md) for how process evidence fits the broader evidence lifecycle.

<!-- reference links -->
[bpmn]: https://www.omg.org/spec/BPMN/2.0/
[CMMN]: https://www.omg.org/spec/CMMN/
[DMN]: https://www.omg.org/spec/DMN/
[OMG]: https://www.omg.org/
[jBPM]: https://kie.apache.org/components/jbpm/
[Drools]: https://kie.apache.org/components/drools/
[Kogito]: https://kogito.kie.org/
[SonataFlow]: https://kie.apache.org/components/sonataflow/
[kie]: https://kie.apache.org/
[Quarkus]: https://quarkus.io/
[knative]: https://knative.dev/
[Kafka]: https://kafka.apache.org/
[Camunda]: https://camunda.com/
[bpmn.io]: https://bpmn.io/
[Flowable]: https://www.flowable.com/open-source
