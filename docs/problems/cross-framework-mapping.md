# Cross-Framework Mapping

## Why this is hard

Organizations assessed against multiple compliance frameworks simultaneously ([NIST 800-53][nist-800-53], [CIS Controls][cis-controls], [PCI-DSS][pci-dss], SOC 2, [ISO 27001][iso-27001]) face overlapping requirements described in different vocabulary, at different granularity, and with different intent. Understanding that NIST AC-2 (Account Management), CIS Control 5.1 (Establish and Maintain an Inventory of Accounts), and PCI Requirement 7 (Restrict Access to System Components and Cardholder Data) address the same general concern is expert work. It requires reading the full text of each requirement, understanding the context in which each framework uses its terms, and making a judgment about the nature of the relationship.

This mapping is currently a manual, subjective, and fragile exercise. It works when one person holds the relationships in their head. It stops working at scale: when the number of frameworks multiplies, when the expert moves on, and when a framework revises (as NIST 800-53 rev 5 did, reorganizing control families and adding privacy controls).

The reasoning behind a mapping is harder to preserve than the mapping itself. A spreadsheet cell linking AC-2 to CIS 5.1 tells you they are related but not why, not what the boundary of the relationship is, and not whether the relationship is exact or partial. When a framework revises, there is no structured way to determine which existing mappings are affected and which still hold.

Private duplication makes this worse. Multiple organizations map the same frameworks independently, with no practical way to reuse one another's effort.

## Current approaches and prior art

Manual spreadsheet mapping is the dominant approach. Subject matter experts create spreadsheets linking controls across frameworks. No machine-readable format, no provenance, no structured reasoning. Fragile and labor-intensive.

[NIST OLIR][nist-olir] (Online Informative References) is a NIST-maintained repository of framework relationships. Authoritative for what it covers but limited in scope, manually maintained, and provides relationship assertions without structured reasoning.

[NIST IR 8477][nist-ir-8477] defines a set-theory approach to compliance mapping with five relationship types: subset of, superset of, intersects with, equal to, and no relationship. Each relationship has a rationale category: syntactic (textual similarity), semantic (meaning equivalence), or functional (operational equivalence). This is the most rigorous public framework for describing mapping relationships.

[OSCAL][oscal] component definitions provide a machine-readable format for compliance data but do not define cross-framework relationships. Components can reference controls from multiple frameworks but the mapping logic is external.

OSCAL has since added a [Mapping Model][oscal-mapping] that provides a machine-readable format for expressing relationships between controls across different frameworks. The mapping model supports source-to-target control relationships with typed associations. This is a direct step toward standardized cross-framework mapping, though adoption is early and the model does not prescribe confidence scoring or reasoning traces.

Unified Compliance Framework (UCF) is a commercial effort to create a meta-framework mapping. Proprietary, subscription-based, and the mapping methodology is not publicly documented.

[Secure Controls Framework (SCF)][scf] is a meta-framework that maps controls across dozens of regulatory and industry frameworks. Unlike UCF, SCF publishes its mappings openly and provides a downloadable control catalog. SCF defines its own canonical control set and maps each framework's controls to it. The mapping methodology is documented and the catalog is actively maintained. SCF demonstrates that a centralized meta-framework approach can scale to broad framework coverage, though it requires a single authoritative mapping team to maintain consistency.

LLM-assisted analysis uses recent-generation language models (both open-weight and commercial) to read requirement text, compare semantic meaning, and propose relationships with confidence scores. Not reliable enough to operate without human oversight, but capable of handling the bulk of initial mapping work and surfacing uncertain cases for expert review.

## Proposed approaches

Explore an LLM-assisted mapping engine with human-in-the-loop oversight, using IR 8477's relationship taxonomy:

### Relationship model

Each mapping between two requirements is characterized by:
- The IR 8477 relationship type (subset, superset, intersects, equal, none)
- The rationale category (syntactic, semantic, functional)
- A confidence score from the mapping engine
- The reasoning trace (what the engine considered, what it was uncertain about)

### Human-in-the-loop workflow

The confidence score drives escalation:
- High-confidence proposals flow through with spot-check review
- Medium-confidence proposals require expert review before acceptance
- Low-confidence proposals are flagged as requiring investigation
- Thresholds are tunable per deployment; the system should track agreement rates between machine proposals and human decisions over time to calibrate

### Graph representation

Requirements are nodes. Mappings are edges with typed relationships, confidence scores, provenance (who proposed, who ratified, when, under what framework version), and the reasoning behind the relationship. This is a property graph problem: the relationships carry rich attributes that a simple adjacency list cannot represent. [GQL][] provides the query language for traversal (impact analysis, traceability chains, coverage gaps).

### Community compounding

When generic framework-to-framework mappings are published to a commons, the community's mapping work compounds. Each organization refines a shared base rather than rebuilding from scratch. Community corrections become labeled data that improves the mapping engine over time.

## Open questions

Cross-framework mapping raises questions about confidence calibration, reasoning provenance, divergent authority handling, version coexistence, and mapping fidelity that require architectural decisions. See [ADR-0011](../ADRs/0011-cross-framework-mapping.md) for the decisions that have crystallized from this exploration.

## Cross-references

See [Requirement Fidelity](requirement-fidelity.md) for how requirements lose context in translation (the same fidelity problem that mapping must address). See [Evidence](evidence.md) for how mapped requirements connect to the evidence lifecycle.

[nist-800-53]: https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final
[cis-controls]: https://www.cisecurity.org/controls
[pci-dss]: https://www.pcisecuritystandards.org/
[iso-27001]: https://www.iso.org/standard/27001
[nist-olir]: https://csrc.nist.gov/projects/olir
[nist-ir-8477]: https://csrc.nist.gov/pubs/ir/8477/final
[oscal]: https://pages.nist.gov/OSCAL/
[GQL]: https://www.gqlstandards.org/
[scf]: https://securecontrolsframework.com/
[oscal-mapping]: https://pages.nist.gov/OSCAL/learn/concepts/layer/control/mapping/
