# Encryption Service Agents Guide

This service manages encryption-related server-side data and key workflows.

## Scope

- encrypted vault and recipe key storage flows
- key exchange or sharing support exposed by backend contracts
- persistence and transport around encrypted metadata

## Working Rules

- Treat every change as security-sensitive.
- Preserve the boundary between client-side cryptographic operations and server-side key management responsibilities.
- Do not weaken assumptions around private key protection, passphrase-derived data, or recipe sharing flows.
- Validate any contract changes against the mobile client expectations.
