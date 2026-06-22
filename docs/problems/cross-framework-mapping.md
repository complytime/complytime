# Cross-Framework Mapping

Organizations assessed against multiple compliance frameworks simultaneously — NIST 800-53, CIS Controls, PCI-DSS, SOC 2, ISO 27001 — face overlapping requirements described in different terms. Understanding that NIST AC-2 and CIS Control 5.1 and PCI Requirement 7 address the same concern is manual, subjective, and fragile. It works when one person holds the relationships in their head. It stops working at scale.

No standard way exists to express, maintain, or query these cross-framework relationships in a machine-readable, traceable form. Mapping is typically a spreadsheet exercise repeated per audit cycle, per framework combination.

This problem area is under active exploration. See [requirement-fidelity.md](requirement-fidelity.md) for how requirements lose context in translation, and [evaluator-coupling.md](evaluator-coupling.md) for how the same intent diverges across tools with no shared mapping.
